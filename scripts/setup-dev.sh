#!/bin/bash

# CloudClip 开发环境设置脚本

set -e

echo "🚀 正在设置 CloudClip 开发环境..."

# 检查 Go 版本
echo "检查 Go 版本..."
go version

# 检查 Docker 和 Docker Compose
echo "检查 Docker 和 Docker Compose..."
docker --version
docker-compose --version

# 安装 Go 工具
echo "安装 Go 开发工具..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 下载 Go 依赖
echo "下载 Go 模块依赖..."
go mod download
go mod tidy

# 创建必要的目录
echo "创建必要的目录..."
mkdir -p bin
mkdir -p logs

# 设置环境变量文件
if [ ! -f .env ]; then
    echo "创建环境变量配置文件..."
    cp .env.example .env
    echo "请编辑 .env 文件以配置您的环境变量"
fi

echo "✅ 开发环境设置完成!"
echo ""
echo "接下来的步骤:"
echo "1. 编辑 .env 文件以配置环境变量"
echo "2. 运行 'make run' 启动应用"
echo "3. 运行 'make test' 执行测试"
