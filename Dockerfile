# 多阶段构建
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apk add --no-cache git make

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN make build

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 复制二进制文件
COPY --from=builder /app/build/es-monitor .

# 设置时区
ENV TZ=Asia/Shanghai

# 入口点
ENTRYPOINT ["./es-monitor"]

# 默认参数
CMD ["-host", "elasticsearch", "-port", "9200"]
