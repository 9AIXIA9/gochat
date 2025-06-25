FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/cloudclip cloudclip.go

# 使用轻量级的基础镜像
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/cloudclip /app/
COPY --from=builder /app/etc/cloudclip-api.yaml /app/etc/

# 设置时区环境变量
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8888

# 运行应用
CMD ["/app/cloudclip", "-f", "/app/etc/cloudclip-api.yaml"]
