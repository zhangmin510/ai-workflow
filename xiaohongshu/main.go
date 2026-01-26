package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"xiaohongshu-auto-poster/config"
	"xiaohongshu-auto-poster/workflow"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "./config/config.yaml", "配置文件路径")
	mode := flag.String("mode", "server", "运行模式: server, once, batch")
	topic := flag.String("topic", "", "指定话题（仅once模式有效）")
	count := flag.Int("count", 3, "生成数量（仅batch模式有效）")

	flag.Parse()

	// 加载配置
	log.Printf("加载配置文件: %s", *configPath)
	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化工作流
	wf := workflow.NewWorkflow()

	// 处理不同运行模式
	switch *mode {
	case "server":
		startServerMode(wf)
	case "once":
		startOnceMode(wf, *topic)
	case "batch":
		startBatchMode(wf, *count)
	default:
		log.Fatalf("不支持的运行模式: %s", *mode)
	}
}

func startServerMode(wf *workflow.Workflow) {
	log.Println("启动服务器模式")

	// 启动工作流
	err := wf.Start()
	if err != nil {
		log.Fatalf("启动工作流失败: %v", err)
	}

	// 设置HTTP服务器
	r := gin.Default()

	// API路由
	api := r.Group("/api")
	{
		// 工作流状态
		api.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"running": wf.IsRunning(),
				"stats":   wf.GetStats(),
			})
		})

		// 启动工作流
		api.POST("/start", func(c *gin.Context) {
			if !wf.IsRunning() {
				err := wf.Start()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
			}
			c.JSON(200, gin.H{"message": "工作流已启动"})
		})

		// 停止工作流
		api.POST("/stop", func(c *gin.Context) {
			if wf.IsRunning() {
				wf.Stop()
			}
			c.JSON(200, gin.H{"message": "工作流已停止"})
		})

		// 紧急停止
		api.POST("/emergency-stop", func(c *gin.Context) {
			wf.EmergencyStop()
			c.JSON(200, gin.H{"message": "工作流已紧急停止"})
		})

		// 手动执行一次
		api.POST("/run-once", func(c *gin.Context) {
			var req struct {
				Topic string `json:"topic"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			err := wf.RunOnce(req.Topic)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"message": "任务执行完成"})
		})

		// 批量生成
		api.POST("/generate-batch", func(c *gin.Context) {
			var req struct {
				Count  int      `json:"count"`
				Topics []string `json:"topics"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			err := wf.GenerateBatch(req.Count, req.Topics)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"message": "批量生成完成"})
		})
	}

	// 优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("正在优雅关闭...")
		wf.Stop()
		os.Exit(0)
	}()

	// 启动HTTP服务
	log.Println("HTTP服务器启动在 :8080")
	log.Fatal(r.Run(":8080"))
}

func startOnceMode(wf *workflow.Workflow, topic string) {
	log.Println("启动单次运行模式")

	err := wf.RunOnce(topic)
	if err != nil {
		log.Fatalf("执行任务失败: %v", err)
	}

	log.Println("单次任务执行完成")
}

func startBatchMode(wf *workflow.Workflow, count int) {
	log.Printf("启动批量生成模式，将生成%d篇内容", count)

	// 这里可以从命令行参数获取话题列表
	// 目前使用配置中的默认话题
	err := wf.GenerateBatch(count, nil)
	if err != nil {
		log.Fatalf("批量生成失败: %v", err)
	}

	log.Println("批量生成任务完成")
}
