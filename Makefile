.PHONY: build clean run install test lint

# 项目名称
BINARY_NAME=es-monitor
VERSION=1.0.0

# Go 参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 构建目录
BUILD_DIR=build

# 默认目标
all: clean build

# 构建
build:
	@echo "构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $$ (BUILD_DIR)/ $$(BINARY_NAME) -v ./cmd/monitor
	@echo "构建完成: $$ (BUILD_DIR)/ $$(BINARY_NAME)"

# 跨平台编译
build-all: build-linux build-mac build-windows

build-linux:
	@echo "构建 Linux 版本..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $$ (BUILD_DIR)/ $$(BINARY_NAME)-linux-amd64 ./cmd/monitor
	@echo "构建完成: $$ (BUILD_DIR)/ $$(BINARY_NAME)-linux-amd64"

build-mac:
	@echo "构建 macOS 版本..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $$ (BUILD_DIR)/ $$(BINARY_NAME)-darwin-amd64 ./cmd/monitor
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $$ (BUILD_DIR)/ $$(BINARY_NAME)-darwin-arm64 ./cmd/monitor
	@echo "构建完成: $$ (BUILD_DIR)/ $$(BINARY_NAME)-darwin-*"

build-windows:
	@echo "构建 Windows 版本..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $$ (BUILD_DIR)/ $$(BINARY_NAME)-windows-amd64.exe ./cmd/monitor
	@echo "构建完成: $$ (BUILD_DIR)/ $$(BINARY_NAME)-windows-amd64.exe"

# 运行
run:
	$(GOBUILD) -o $$ (BUILD_DIR)/ $$(BINARY_NAME) -v ./cmd/monitor
	./$$ (BUILD_DIR)/ $$(BINARY_NAME)

# 运行并指定参数
run-dev:
	$(GOBUILD) -o $$ (BUILD_DIR)/ $$(BINARY_NAME) -v ./cmd/monitor
	./$$ (BUILD_DIR)/ $$(BINARY_NAME) -host localhost -port 9200 -interval 2

# 清理
clean:
	@echo "清理构建文件..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "清理完成"

# 安装依赖
install:
	@echo "安装依赖..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "依赖安装完成"

# 测试
test:
	@echo "运行测试..."
	$(GOTEST) -v ./...

# 代码检查
lint:
	@echo "运行代码检查..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "请先安装 golangci-lint"; exit 1; }
	golangci-lint run ./...

# 格式化代码
fmt:
	@echo "格式化代码..."
	$(GOCMD) fmt ./...

# 生成版本信息
version:
	@echo "版本: $(VERSION)"

# Docker 构建
docker-build:
	@echo "构建 Docker 镜像..."
	docker build -t $$ (BINARY_NAME): $$(VERSION) .
	@echo "Docker 镜像构建完成"

# 帮助
help:
	@echo "可用命令:"
	@echo "  make build        - 构建项目"
	@echo "  make build-all    - 跨平台构建（Linux/Mac/Windows）"
	@echo "  make run          - 运行项目"
	@echo "  make run-dev      - 开发模式运行"
	@echo "  make clean        - 清理构建文件"
	@echo "  make install      - 安装依赖"
	@echo "  make test         - 运行测试"
	@echo "  make lint         - 代码检查"
	@echo "  make fmt          - 格式化代码"
	@echo "  make docker-build - 构建 Docker 镜像"
