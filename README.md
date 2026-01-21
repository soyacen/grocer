# grocer

`grocer`是一个Go语言微服务基础设施库，旨在为分布式应用提供统一的中间件和服务集成能力，提升开发效率和系统可维护性。

## 特性

- **dbx**: 统一MySQL/PostgreSQL数据库访问层
- **esx**: Elasticsearch集成
- **goseex**: Goose数据库迁移工具集成
- **grpcx**: gRPC客户端与服务端封装
- **idx**: ULID唯一ID生成器
- **jeagerx**: Jaeger分布式追踪集成
- **kafkax**: Kafka生产者/消费者支持
- **mongox**: MongoDB连接与查询封装
- **nacosx**: Nacos配置中心与服务发现
- **otelx**: OpenTelemetry资源初始化
- **promx**: Prometheus指标收集
- **protobufx**: TLS等通用Protobuf消息定义
- **redisx**: Redis多模式（单机/集群/哨兵）支持
- **registryx**: 服务注册与定义抽象

## 快速开始

### 使用Helm部署

我们提供了Helm charts来方便地部署grocer相关的服务。

#### 安装前提

- Kubernetes 1.19+
- Helm 3.0+

#### 安装Chart

要安装名为`my-grocer`的chart，请执行以下命令：

```bash
# 进入charts目录
cd charts/grocer

# 安装chart
helm install my-grocer .
```

#### 自定义配置

如果需要自定义配置，请创建一个values.yaml文件，然后使用以下命令安装：

```bash
helm install my-grocer . -f my-values.yaml
```

#### 卸载Chart

要卸载名为`my-grocer`的部署，请执行以下命令：

```bash
helm delete my-grocer
```

更多关于Helm chart的配置选项，请参阅 [charts/grocer/README.md](charts/grocer/README.md)。

### 安装

``bash
go get github.com/soyacen/grocer
```

### 使用示例

``go
package main

import (
    "github.com/soyacen/grocer/redisx"
    "github.com/soyacen/grocer/dbx"
)

func main() {
    // 使用 redisx 连接 Redis
    redisClient := redisx.NewClient(&redisx.Config{
        // 配置参数
    })
    
    // 使用 dbx 连接数据库
    db := dbx.Connect(&dbx.Config{
        // 配置参数
    })
    
    // ... 其他业务逻辑
}
```

## 架构设计

`grocer` 采用模块化设计，每个子包都是一个独立的功能模块。通过 Google Wire 实现依赖注入，使得组件间的耦合度更低，更易于测试和维护。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进项目。

## 许可证

MIT License