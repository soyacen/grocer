package temporalx

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/soyacen/gox/conc/lazyload"
	"go.temporal.io/sdk/client"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type IConfig interface {
	GetConfigs() map[string]IOptions
}

type IOptions interface {
	GetHostPort() *wrapperspb.StringValue
	GetNamespace() *wrapperspb.StringValue
}

type ISDKOptions interface {
	GetSDKOptions() client.Options
}

func NewClients(config IConfig) *lazyload.Group[client.Client] {
	return &lazyload.Group[client.Client]{
		New: func(key string) (client.Client, error) {
			options, ok := config.GetConfigs()[key]
			if !ok {
				return nil, errors.Errorf("temporalx: config %s not found", key)
			}
			return NewClient(options)
		},
		Finalize: func(ctx context.Context, obj client.Client) error {
			obj.Close()
			return nil
		},
	}
}

func NewClient(options IOptions) (client.Client, error) {
	var sdkOption client.Options
	if getter, ok := options.(ISDKOptions); ok {
		sdkOption = getter.GetSDKOptions()
	}
	if options.GetHostPort() != nil {
		sdkOption.HostPort = options.GetHostPort().GetValue()
	}
	if options.GetNamespace() != nil {
		sdkOption.Namespace = options.GetNamespace().GetValue()
	}
	return client.Dial(sdkOption)
}
