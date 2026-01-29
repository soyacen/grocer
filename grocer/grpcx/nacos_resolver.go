package grpcx

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/soyacen/grocer/grocer/nacosx"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(&nacosResolverBuilder{})
}

type nacosResolverBuilder struct{}

type nacosResolver struct {
	cancelFunc context.CancelFunc
	client     naming_client.INamingClient
}

type nacosResolverTarget struct {
	Timeout      time.Duration
	PollInterval time.Duration
}

func (b *nacosResolverBuilder) Build(tgt resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	dsn := strings.Join([]string{b.Scheme() + ":/", tgt.URL.Host, tgt.URL.Path + "?" + tgt.URL.RawQuery}, "/")
	rawURL, err := url.Parse(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "nacosx: malformed URL")
	}
	target, err := nacosx.ParseDSN(rawURL)
	if err != nil {
		return nil, errors.Wrap(err, "Wrong nacos URL")
	}
	cli, err := nacosx.NewNamingClient(target)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't connect to the Nacos API")
	}
	ctx, cancel := context.WithCancel(context.Background())
	pipe := make(chan []model.Instance, 1)
	go watchNacosService(ctx, cli, target, rawURL.Path, pipe)
	go populateEndpoints(ctx, cc, pipe)

	return &nacosResolver{cancelFunc: cancel, client: cli}, nil
}

func (b *nacosResolverBuilder) Scheme() string {
	return nacosx.Scheme
}

func (r *nacosResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *nacosResolver) Close() {
	r.cancelFunc()
	if r.client != nil {
		r.client.CloseClient()
	}
}

func watchNacosService(ctx context.Context, client naming_client.INamingClient, t *nacosx.Options, service string, output chan<- []model.Instance) {
	// Subscribe to service updates
	subscribeParam := &vo.SubscribeParam{
		ServiceName: service,
		GroupName:   t.GetGroupName().GetValue(),
		Clusters:    []string{t.GetClusterName().GetValue()},
		SubscribeCallback: func(services []model.Instance, err error) {
			if err != nil {
				grpclog.Errorf("[Nacos resolver] Error in subscription callback for service=%s; error=%v", service, err)
				return
			}
			grpclog.Infof("[Nacos resolver] %d endpoints received for service=%s", len(services), service)
			select {
			case output <- services:
			case <-ctx.Done():
				return
			}
		},
	}

	// Start subscription in a goroutine
	go func() {
		if err := client.Subscribe(subscribeParam); err != nil {
			grpclog.Errorf("[Nacos resolver] Couldn't subscribe to service=%s; error=%v", service, err)
		}
	}()

	// Handle unsubscribe when context is cancelled
	go func() {
		<-ctx.Done()
		if err := client.Unsubscribe(subscribeParam); err != nil {
			grpclog.Errorf("[Nacos resolver] Couldn't unsubscribe to service=%s; error=%v", service, err)
		}
	}()
}

func populateEndpoints(ctx context.Context, clientConn resolver.ClientConn, input <-chan []model.Instance) {
	for {
		select {
		case instances, ok := <-input:
			if !ok {
				// Channel was closed, exit the loop
				grpclog.Info("[Nacos resolver] Input channel closed, ending population")
				return
			}

			instanceSet := make(map[string]model.Instance, len(instances))
			for _, instance := range instances {
				// Filter out unhealthy or disabled instances
				if !instance.Enable || !instance.Healthy || instance.Weight <= 0 {
					continue
				}
				addr := fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
				instanceSet[addr] = instance
			}

			addresses := make([]resolver.Address, 0, len(instanceSet))
			for addr, instance := range instanceSet {
				address := resolver.Address{
					Addr: addr,
				}
				address.Attributes = attributes.New("instance_id", instance.InstanceId)
				address.Attributes = address.Attributes.WithValue("cluster_name", instance.ClusterName)
				address.Attributes = address.Attributes.WithValue("service_name", instance.ServiceName)
				for key, value := range instance.Metadata {
					address.Attributes = address.Attributes.WithValue(key, value)
				}
				addresses = append(addresses, address)
			}
			slices.SortFunc(addresses, func(a, b resolver.Address) int { return strings.Compare(a.Addr, b.Addr) })
			if err := clientConn.UpdateState(resolver.State{Addresses: addresses}); err != nil {
				grpclog.Errorf("[Nacos resolver] Couldn't update client connection. error=%v", err)
				continue
			}
			grpclog.Infof("[Nacos resolver] Updated state with %d healthy endpoints", len(addresses))
		case <-ctx.Done():
			grpclog.Info("[Nacos resolver] Watch has been finished")
			return
		}
	}
}
