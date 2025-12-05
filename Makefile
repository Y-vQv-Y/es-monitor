.PHONY: build clean run install test lint fmt help build-all docker-build release

# 项目信息
BINARY_NAME=es-monitor
VERSION=1.0.0
BUILD_DIR=build
MAIN_PATH=./cmd/monitor

# Go 命令
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# 构建标志
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -s -w"

# 默认目标
all: clean build

# 构建
build:
	@echo "==> 构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)
	@echo "==> 构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 跨平台编译
build-all: build-linux build-darwin build-windows
	@echo "==> 所有平台构建完成"

build-linux:
	@echo "==> 构建 Linux 版本..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "==> Linux 版本构建完成"

build-darwin:
	@echo "==> 构建 macOS 版本..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "==> macOS 版本构建完成"

build-windows:
	@echo "==> 构建 Windows 版本..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "==> Windows 版本构建完成"

# 运行
run: build
	@echo "==> 运行 $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# 开发模式运行
run-dev: build
	@echo "==> 开发模式运行..."
	./$(BUILD_DIR)/$(BINARY_NAME) -host localhost -port 9200 -interval 2

# 清理
clean:
	@echo "==> 清理构建文件..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "==> 清理完成"

# 安装依赖
install:
	@echo "==> 安装依赖..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "==> 依赖安装完成"

# 测试
test:
	@echo "==> 运行测试..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "==> 测试完成"

# 查看测试覆盖率
coverage: test
	@echo "==> 生成覆盖率报告..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "==> 覆盖率报告已生成: coverage.html"

# 代码检查
lint:
	@echo "==> 运行代码检查..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "请先安装 golangci-lint"; exit 1; }
	golangci-lint run ./...
	@echo "==> 代码检查完成"

# 格式化代码
fmt:
	@echo "==> 格式化代码..."
	$(GOFMT) ./...
	@echo "==> 代码格式化完成"

# Docker 构建
docker-build:
	@echo "==> 构建 Docker 镜像..."
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest
	@echo "==> Docker 镜像构建完成"

# 创建发布包
release: build-all
	@echo "==> 创建发布包..."
	@mkdir -p $(BUILD_DIR)/release
	cd $(BUILD_DIR) && \
	tar czf release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
	tar czf release/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64 && \
	tar czf release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
	tar czf release/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 && \
	zip -q release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "==> 发布包已创建在 $(BUILD_DIR)/release/"
	@ls -lh $(BUILD_DIR)/release/

# 安装到系统
install-system: build
	@echo "==> 安装到系统..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "==> 已安装到 /usr/local/bin/$(BINARY_NAME)"

# 卸载
uninstall:
	@echo "==> 卸载..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "==> 卸载完成"

# 显示版本
version:
	@echo "版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"

# 帮助
help:
	@echo "可用命令:"
	@echo "  make build           - 构建项目（当前平台）"
	@echo "  make build-all       - 跨平台构建（Linux/Mac/Windows）"
	@echo "  make run             - 构建并运行"
	@echo "  make run-dev         - 开发模式运行"
	@echo "  make clean           - 清理构建文件"
	@echo "  make install         - 安装依赖"
	@echo "  make test            - 运行测试"
	@echo "  make coverage        - 生成测试覆盖率报告"
	@echo "  make lint            - 代码检查"
	@echo "  make fmt             - 格式化代码"
	@echo "  make docker-build    - 构建 Docker 镜像"
	@echo "  make release         - 创建发布包"
	@echo "  make install-system  - 安装到系统路径"
	@echo "  make uninstall       - 从系统卸载"
	@echo "  make version         - 显示版本信息"
