@echo off
rem GoChat 开发环境设置脚本 (Windows)

echo 🚀 正在设置 GoChat 开发环境...

rem 检查 Go 版本
echo 检查 Go 版本...
go version
if %errorlevel% neq 0 (
    echo 错误: 未找到 Go，请先安装 Go 1.21+
    exit /b 1
)

rem 检查 Docker 和 Docker Compose
echo 检查 Docker 和 Docker Compose...
docker --version
if %errorlevel% neq 0 (
    echo 错误: 未找到 Docker，请先安装 Docker Desktop
    exit /b 1
)

docker-compose --version
if %errorlevel% neq 0 (
    echo 错误: 未找到 Docker Compose，请确保 Docker Desktop 已正确安装
    exit /b 1
)

rem 安装 Go 工具
echo 安装 Go 开发工具...
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/zeromicro/go-zero/tools/goctl@latest

rem 下载 Go 依赖
echo 下载 Go 模块依赖...
go mod download
go mod tidy

rem 创建必要的目录
echo 创建必要的目录...
if not exist bin mkdir bin
if not exist logs mkdir logs

rem 设置环境变量文件
if not exist .env (
    echo 创建环境变量配置文件...
    copy .env.example .env
    echo 请编辑 .env 文件以配置您的环境变量
)

echo ✅ 开发环境设置完成!
echo.
echo 接下来的步骤:
echo 1. 编辑 .env 文件以配置环境变量
echo 2. 运行 'make run' 启动应用
echo 3. 运行 'make test' 执行测试

pause
