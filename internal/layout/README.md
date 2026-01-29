# project
## 项目结构
```
.
├── Dockerfile
├── Makefile
├── README.md
├── api
│   ├── doc.go
│   ├── grpc
│   │   ├── grpc.pb.go
│   │   ├── grpc.proto
│   │   ├── grpc_goose.pb.go
│   │   └── grpc_grpc.pb.go
│   └── http
│       ├── http.pb.go
│       ├── http.proto
│       ├── http_goose.pb.go
│       └── http_grpc.pb.go
├── cmd
│   ├── cronjob.go
│   ├── fx.go
│   ├── grpc.go
│   ├── http.go
│   ├── job.go
│   └── root.go
├── config
│   ├── config.go
│   ├── config.pb.go
│   ├── config.proto
│   ├── config_gonfig.pb.go
│   ├── config_template.yaml
│   ├── config_test.go
│   └── fx.go
├── deploy
│   ├── Chart.yaml
│   ├── templates
│   │   ├── cronjob.yaml
│   │   ├── grpc.yaml
│   │   ├── http.yaml
│   │   ├── job.yaml
│   │   └── nacos-configmap.yaml
│   └── values
│       ├── common.yaml
│       ├── cronjob.yaml
│       ├── grpc.yaml
│       ├── http.yaml
│       └── job.yaml
├── go.mod
├── go.sum
├── internal
│   ├── cronjob
│   │   ├── fx.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── doc.go
│   ├── grpc
│   │   ├── index.go
│   │   ├── model.go
│   │   └── service.go
│   ├── http
│   │   ├── index.go
│   │   ├── model.go
│   │   └── service.go
│   └── job
│       ├── fx.go
│       ├── model.go
│       ├── repository.go
│       └── service.go
├── main.go
├── pkg
│   ├── doc.go
│   └── logx
│       └── log.go
├── scripts
│   ├── common_proto_lib.sh
│   ├── config.sh
│   ├── format.sh
│   ├── go_gen.sh
│   ├── goose.sh
│   ├── grpc.sh
│   └── lint.sh
├── third_party
│   ├── google
│   │   ├── api
│   │   │   ├── annotations.proto
│   │   │   ├── http.proto
│   │   │   └── httpbody.proto
│   │   ├── rpc
│   │   │   ├── code.proto
│   │   │   ├── error_details.proto
│   │   │   ├── http.proto
│   │   │   └── status.proto
│   │   └── type
│   │       ├── calendar_period.proto
│   │       ├── color.proto
│   │       ├── date.proto
│   │       ├── datetime.proto
│   │       ├── dayofweek.proto
│   │       ├── decimal.proto
│   │       ├── expr.proto
│   │       ├── fraction.proto
│   │       ├── interval.proto
│   │       ├── latlng.proto
│   │       ├── localized_text.proto
│   │       ├── money.proto
│   │       ├── month.proto
│   │       ├── phone_number.proto
│   │       ├── postal_address.proto
│   │       ├── quaternion.proto
│   │       └── timeofday.proto
│   ├── grocer
│   │   ├── dbx
│   │   │   └── config.proto
│   │   ├── esx
│   │   │   └── config.proto
│   │   ├── jeagerx
│   │   │   └── config.proto
│   │   ├── kafkax
│   │   │   └── config.proto
│   │   ├── mongox
│   │   │   └── config.proto
│   │   ├── nacosx
│   │   │   └── config.proto
│   │   ├── protobufx
│   │   │   └── tls.proto
│   │   ├── pyroscopex
│   │   │   └── config.proto
│   │   ├── redisx
│   │   │   └── config.proto
│   │   └── s3x
│   │       └── config.proto
│   └── validate
│       └── validate.proto
└── tools
    └── tools.go

```