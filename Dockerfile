FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apk add --no-cache git make ca-certificates tzdata

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建（静态链接）
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=1.0.0 -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
    -o es-monitor \
    ./cmd/monitor

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 创建非 root 用户（生产环境安全）
RUN addgroup -g 1000 esmonitor && \
    adduser -D -u 1000 -G esmonitor esmonitor

WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/es-monitor .

# 修改所有权
RUN chown -R esmonitor:esmonitor /app

# 切换到非 root 用户
USER esmonitor

# 设置时区
ENV TZ=Asia/Shanghai

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pgrep -f es-monitor || exit 1

# 入口点
ENTRYPOINT ["./es-monitor"]

# 默认参数（只读模式）
CMD ["-host", "elasticsearch", "-port", "9200", "-interval", "2", "-readonly"]
