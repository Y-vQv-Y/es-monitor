# Elasticsearch 实时监控工具

一个功能全面的 Elasticsearch 命令行实时监控工具，支持物理机、Docker、Kubernetes 等多种部署环境。

## 功能特性

### 集群监控
- 集群健康状态（Green/Yellow/Red）
- 节点数量统计
- 分片状态（活跃/迁移/未分配）
- 待处理任务监控

### 节点监控
- **JVM 指标**
  - 堆内存使用率和容量
  - 非堆内存统计
  - 线程数和 GC 统计
  
- **系统资源**
  - CPU 使用率和负载
  - 系统内存使用情况
  - Swap 内存统计
  - 磁盘空间和 IO
  - 文件描述符使用率
  
- **网络统计**
  - 传输层流量统计
  - HTTP 连接数
  
- **读写性能**
  - 实时索引写入速率
  - 实时查询速率
  - 文档统计
  - Merge、Refresh、Flush 操作统计

### 索引监控
- 索引健康状态
- 文档数量和大小
- 分片配置
- 实时写入和查询速率
- 分段信息（数量、内存使用）（新增）

### 系统监控
- **CPU 详细信息**
  - 总使用率
  - 用户态/内核态/空闲/IO等待百分比
  
- **内存详细信息**
  - 总容量/已使用/可用
  - 缓冲区和缓存
  - Swap 使用情况
  
- **磁盘详细信息**
  - 实时读写吞吐量（字节/秒）
  - 实时读写操作数（ops/秒）
  - IO 利用率（新增）
  - 分区使用情况及分区 IO 详细（新增）
  
- **网络详细信息**
  - 实时发送/接收吞吐量
  - 实时数据包速率
  - 各网卡统计及瞬时速率（新增）

### 异常告警
自动识别并标记异常：
- 集群状态异常
- JVM 堆内存过高
- CPU 使用率过高
- 系统内存不足
- 磁盘空间不足
- Full GC 频繁
- 未分配分片

## 快速开始

### 前置要求
- Go 1.21+
- Elasticsearch 7.x/8.x

### 编译安装

```bash
# 克隆项目
git clone https://github.com/Y-vQv-Y/es-monitor.git
cd es-monitor

# 安装依赖
make install

# 编译
make build

# 运行
./build/es-monitor -host localhost -port 9200

# 编译所有平台
make build-all

# 单独编译
make build-linux    # Linux
make build-mac      # macOS
make build-windows  # Windows

# 构建镜像
make docker-build

# 运行容器
docker run --rm es-monitor:1.0.0 -host es-host -port 9200

# 使用 docker-compose
docker-compose up


# 创建 ConfigMap
kubectl create configmap es-monitor-config \
  --from-literal=ES_HOST=elasticsearch.default.svc.cluster.local \
  --from-literal=ES_PORT=9200

# 部署 Pod
kubectl apply -f k8s/deployment.yaml


./es-monitor [选项]

选项:
  -host string
        Elasticsearch 主机地址 (默认 "localhost")
  -port string
        Elasticsearch 端口 (默认 "9200")
  -interval int
        刷新间隔，单位秒 (默认 2)
  -user string
        用户名（可选）
  -pass string
        密码（可选）

示例:
  # 默认连接
  ./es-monitor

  # 指定地址
  ./es-monitor -host 192.168.1.100 -port 9200

  # 快捷方式
  ./es-monitor 192.168.1.100:9200

  # 带认证
  ./es-monitor -host es-host -port 9200 -user elastic -pass password

  # 自定义刷新间隔（5秒）
  ./es-monitor -host es-host -port 9200 -interval 5

  # docker 运行
  docker run -d \
    --name es-monitor \
    --restart unless-stopped \
    --network host \
    --cpus="0.5" \
    --memory="256m" \
    yvqvy/es-monitor:master \
    -host localhost -port 9200 -interval 5
```
### 监控阈值
默认告警阈值

JVM 堆内存
警告: 75%
严重: 85%

CPU 使用率
警告: 60%
严重: 80%

系统内存
警告: 80%
严重: 90%

磁盘使用
警告: 80%
严重: 90%
可在 internal/config/config.go 中自定义阈值。


### 项目结构
```text
es-monitor/
├── cmd/monitor/          # 程序入口
├── internal/
│   ├── client/          # ES 客户端
│   ├── collector/       # 指标采集器
│   ├── config/          # 配置管理
│   ├── display/         # 终端显示
│   ├── model/           # 数据模型
│   └── monitor/         # 监控核心
├── pkg/util/            # 工具函数
├── build/               # 构建输出
├── Makefile             # 构建脚本
├── Dockerfile           # Docker 配置
├── docker-compose.yml   # Docker Compose 配置
└── README.md            # 项目文档
```

### 扩展开发
添加新的采集器

在 internal/collector/ 创建新文件
实现采集接口
在 internal/monitor/monitor.go 中集成

添加新的显示模块

在 internal/display/terminal.go 添加显示方法
在 internal/monitor/monitor.go 中调用

添加新的数据模型

在 internal/model/ 定义结构体
在 collector 中使用


### 贡献
欢迎提交 Issue 和 Pull Request！

### 许可证
MIT License

### 作者
Y-vQv-Y

#### 更新日志
v1.0.0 (2025-12-05)

初始版本发布
支持集群、节点、索引监控
支持系统资源详细监控
支持异常自动告警
