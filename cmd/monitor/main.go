package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Y-vQv-Y/es-monitor/internal/client"
	"github.com/Y-vQv-Y/es-monitor/internal/config"
	"github.com/Y-vQv-Y/es-monitor/internal/monitor"
)

var (
	Version   = "1.0.0"
	BuildTime = "unknown"
)

func main() {
	// 解析命令行参数
	var (
		host     = flag.String("host", "localhost", "Elasticsearch 主机地址")
		port     = flag.String("port", "9200", "Elasticsearch 端口")
		interval = flag.Int("interval", 2, "刷新间隔（秒）")
		username = flag.String("user", "", "用户名（可选）")
		password = flag.String("pass", "", "密码（可选）")
		version  = flag.Bool("version", false, "显示版本信息")
		readonly = flag.Bool("readonly", true, "只读模式（生产环境必须开启）")
	)
	flag.Parse()

	// 显示版本
	if *version {
		fmt.Printf("ES Monitor Version: %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		os.Exit(0)
	}

	// 生产环境安全检查
	if !*readonly {
		fmt.Println("[警告] 非只读模式在生产环境中不安全，强制启用只读模式")
		*readonly = true
	}

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
		ReadOnly: *readonly,
	}

	// 显示启动信息
	fmt.Println("===========================================")
	fmt.Printf("ES Monitor v%s\n", Version)
	fmt.Println("生产环境安全监控工具 - 只读模式")
	fmt.Println("===========================================")

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
	fmt.Println("\n正在安全退出...")
	mon.Stop()
	fmt.Println("已安全退出")
}

func parseAddress(addr string) (host, port string) {
	// 使用 strings 包
	idx := strings.LastIndex(addr, ":")
	if idx == -1 {
		return addr, ""
	}
	return addr[:idx], addr[idx+1:]
}
