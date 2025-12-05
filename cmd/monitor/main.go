package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/es-monitor/internal/client"
	"github.com/yourusername/es-monitor/internal/config"
	"github.com/yourusername/es-monitor/internal/monitor"
)

func main() {
	// 解析命令行参数
	var (
		host     = flag.String("host", "localhost", "Elasticsearch 主机地址")
		port     = flag.String("port", "9200", "Elasticsearch 端口")
		interval = flag.Int("interval", 2, "刷新间隔（秒）")
		username = flag.String("user", "", "用户名（可选）")
		password = flag.String("pass", "", "密码（可选）")
	)
	flag.Parse()

	// 如果第一个参数是 host:port 格式
	if flag.NArg() > 0 {
		addr := flag.Arg(0)
		if h, p := parseAddress(addr); h != "" {
			*host = h
			if p != "" {
				*port = p
			}
		}
	}

	// 创建配置
	cfg := &config.Config{
		Host:     *host,
		Port:     *port,
		Username: *username,
		Password: *password,
		Interval: time.Duration(*interval) * time.Second,
	}

	// 创建 ES 客户端
	esClient := client.NewElasticsearchClient(cfg)

	// 测试连接
	fmt.Printf("正在连接 Elasticsearch: %s:%s ...\n", cfg.Host, cfg.Port)
	ctx := context.Background()
	if err := esClient.Ping(ctx); err != nil {
		fmt.Printf("[错误] 连接失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("[成功] 连接成功!")
	time.Sleep(1 * time.Second)

	// 创建监控器
	mon := monitor.NewMonitor(esClient, cfg)

	// 处理退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动监控
	go mon.Start(ctx)

	// 等待退出信号
	<-sigChan
	fmt.Println("\n正在退出...")
	mon.Stop()
}

func parseAddress(addr string) (host, port string) {
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[:i], addr[i+1:]
		}
	}
	return addr, ""
}
