.PHONY: build clean test run help docker-build docker-push docker-run

APP_NAME=cloudclip
BUILD_DIR=./bin
MAIN_FILE=cloudclip.go
CONFIG_FILE=etc/cloudclip-api.yaml

# 获取当前git commit id
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
VERSION ?= $(shell git describe --tags --always --dirty --match "v*" 2> /dev/null || echo "v0.0.0")

help:
	@echo "使用说明:"
	@echo "make build    - 编译应用程序"
	@echo "make clean    - 清理编译产物"
	@echo "make test     - 运行测试"
	@echo "make run      - 本地运行应用"
	@echo "make lint     - 运行代码检查"
	@echo "make docker-build - 构建Docker镜像"
	@echo "make docker-push  - 推送Docker镜像"
	@echo "make docker-run   - 运行Docker容器"

build:
	@echo "编译应用..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.GitCommit=$(GIT_COMMIT)" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "编译完成: $(BUILD_DIR)/$(APP_NAME)"

clean:
	@echo "清理编译产物..."
	@rm -rf $(BUILD_DIR)
	@echo "清理完成"

test:
	@echo "运行测试..."
	go test -v ./...

run:
	@echo "运行应用..."
	go run $(MAIN_FILE) -f $(CONFIG_FILE)

lint:
	@echo "运行代码检查..."
	golangci-lint run

docker-build:
	@echo "构建Docker镜像..."
	docker build -t $(APP_NAME):$(VERSION) -t $(APP_NAME):latest .
	@echo "Docker镜像构建完成: $(APP_NAME):$(VERSION)"

docker-push:
	@echo "推送Docker镜像..."
	docker push $(APP_NAME):$(VERSION)
	docker push $(APP_NAME):latest
	@echo "Docker镜像推送完成"

docker-run:
	@echo "运行Docker容器..."
	docker run --name $(APP_NAME) -p 8888:8888 -d $(APP_NAME):latest
