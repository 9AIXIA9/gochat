# 第一阶段：构建器
FROM golang:1.25-alpine AS builder
WORKDIR /app

# 配置国内镜像源
RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories

# 配置Go环境
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=1

# 安装构建依赖
RUN apk add --no-cache git build-base pkgconfig librdkafka-dev

# 下载依赖并清理缓存
COPY go.mod go.sum ./
RUN go mod download && \
    go mod verify && \
    rm -rf "$(go env GOPATH)/pkg/mod"

# 复制源码并构建
COPY . .

# 生成代码
RUN go generate ./...

# 构建静态链接的二进制文件
RUN GOOS=linux GOARCH=amd64 go build -tags musl -ldflags="-s -w" -o /app/server ./cmd;

# 第二阶段：运行时
FROM alpine:3.20
WORKDIR /app

# 配置国内镜像源并安装运行时依赖（tzdata 以支持时区、curl 用于健康检查、证书）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache tzdata curl ca-certificates && \
    update-ca-certificates

# 默认时区可被外部 TZ 覆盖（compose .env 中设置）
ENV TZ=Asia/Shanghai

# 设置容器本地时区（使 Go/系统日志等均使用本地时间）
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN mkdir -p /app/config /app/logs /app/data /app/certs

# 复制必需的库文件（librdkafka 运行时依赖）
COPY --from=builder /usr/lib/librdkafka.so.* /usr/lib/

# 创建非 root 用户
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app
USER appuser

# 复制二进制文件
COPY --from=builder --chown=appuser:appgroup /app/server /app/server

# 暴露端口和健康检查
EXPOSE 8888
HEALTHCHECK --interval=30s --timeout=3s \
CMD curl -f http://localhost:8888/health_check || exit 1

CMD ["/app/server"]