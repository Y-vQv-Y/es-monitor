package display

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/Y-vQv-Y/es-monitor/internal/model"
)

// DisplayHealthIssues 显示健康问题
func (t *Terminal) DisplayHealthIssues(issues []model.HealthIssue) {
	if len(issues) == 0 {
		StatusGreen.Println("[集群健康检查] 未发现问题")
		fmt.Println()
		return
	}

	// 按严重程度排序
	sort.Slice(issues, func(i, j int) bool {
		levelOrder := map[string]int{"critical": 0, "warning": 1, "info": 2}
		return levelOrder[issues[i].Level] < levelOrder[issues[j].Level]
	})

	SectionColor.Println("[集群健康检查 - 发现问题]")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))
	fmt.Println()

	criticalCount := 0
	warningCount := 0
	infoCount := 0

	for _, issue := range issues {
		var levelColor *color.Color
		var levelText string

		switch issue.Level {
		case "critical":
			levelColor = StatusRed
			levelText = "[严重]"
			criticalCount++
		case "warning":
			levelColor = StatusYellow
			levelText = "[警告]"
			warningCount++
		case "info":
			levelColor = InfoColor
			levelText = "[提示]"
			infoCount++
		}

		levelColor.Printf("%s ", levelText)
		fmt.Printf("[%s] %s", issue.Component, issue.Message)
		if issue.NodeName != "" {
			fmt.Printf(" (节点: %s)", issue.NodeName)
		}
		if issue.IndexName != "" {
			fmt.Printf(" (索引: %s)", issue.IndexName)
		}
		fmt.Println()

		if issue.Suggestion != "" {
			fmt.Printf("       建议: %s\n", issue.Suggestion)
		}
		fmt.Println()
	}

	// 统计摘要
	fmt.Println(DrawSeparator(DisplayWidth, "-"))
	if criticalCount > 0 {
		StatusRed.Printf("严重问题: %d  ", criticalCount)
	}
	if warningCount > 0 {
		StatusYellow.Printf("警告: %d  ", warningCount)
	}
	if infoCount > 0 {
		InfoColor.Printf("提示: %d", infoCount)
	}
	fmt.Println()
	fmt.Println()
}

// DisplayThreadPools 显示线程池状态
func (t *Terminal) DisplayThreadPools(pools map[string]model.ThreadPoolStats) {
	if len(pools) == 0 {
		return
	}

	SectionColor.Println("[线程池状态]")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	fmt.Printf("%-15s %8s %8s %10s %12s %12s\n",
		"线程池", "活跃", "队列", "队列大小", "拒绝数", "完成数")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	// 关键线程池顺序
	keyPools := []string{"search", "write", "get", "bulk", "management"}

	for _, poolName := range keyPools {
		if pool, ok := pools[poolName]; ok {
			fmt.Printf("%-15s %8d %8d %10d %12d %12d",
				pool.PoolName,
				pool.Active,
				pool.Queue,
				pool.QueueSize,
				pool.Rejected,
				pool.Completed)

			// 警告判断
			if pool.Rejected > 0 {
				StatusRed.Print(" [有拒绝]")
			}
			queuePercent := float64(pool.Queue) / float64(pool.QueueSize) * 100
			if queuePercent > 80 {
				StatusYellow.Printf(" [队列满: %.1f%%]", queuePercent)
			}
			fmt.Println()
		}
	}

	fmt.Println()
}

// DisplayCircuitBreakers 显示断路器状态
func (t *Terminal) DisplayCircuitBreakers(breakers map[string]model.CircuitBreakerStats) {
	if len(breakers) == 0 {
		return
	}

	SectionColor.Println("[断路器状态]")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	fmt.Printf("%-20s %12s %12s %10s %10s\n",
		"断路器", "限制", "使用", "触发次数", "使用率")
	fmt.Println(DrawSeparator(DisplayWidth, "-"))

	for _, breaker := range breakers {
		fmt.Printf("%-20s %10.2fMB %10.2fMB %10d %9.1f%%",
			breaker.Name,
			breaker.LimitSizeMB,
			breaker.EstimatedMB,
			breaker.Tripped,
			breaker.UsedPercent)

		// 警告判断
		if breaker.Tripped > 0 {
			StatusRed.Print(" [已触发]")
		} else if breaker.UsedPercent > 80 {
			StatusYellow.Print(" [接近限制]")
		}
		fmt.Println()
	}

	fmt.Println()
}
