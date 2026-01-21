# 使用特定版本的基础镜像
FROM golang:1.25-alpine AS builder

# 设置必要的环境变量
ENV CGO_ENABLED=0 \
    GOOS=linux

# 设置工作目录
WORKDIR /app

# 首先复制 go 相关文件以利用 Docker 缓存
COPY go.mod go.sum ./

# 下载项目依赖
RUN go mod download

# 将源代码复制到容器中
COPY . .

# 构建应用，使用 ldflags 减少二进制大小
RUN go build -ldflags="-w -s" -o app .

# 使用更小的 alpine 镜像作为最终镜像
FROM alpine:3.23 AS runtime

# 安装 ca-certificates，如果应用需要连接 HTTPS 服务
RUN apk --no-cache add ca-certificates

# 设置时区（可选）
RUN apk add --no-cache tzdata

# 创建非 root 用户运行应用
RUN addgroup -g 65532 -S nonroot &&\
    adduser -S -u 65532 -G nonroot nonroot

# 设置工作目录
WORKDIR /app

# 从 builder 镜像中复制编译好的二进制文件
COPY --from=builder /app/app .

# 更改文件权限给 nonroot 用户
RUN chown -R nonroot:nonroot /app && \
    chmod +x ./app
USER nonroot

# 启动应用
ENTRYPOINT ["/app/app"]