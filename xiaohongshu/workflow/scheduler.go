package workflow

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"xiaohongshu-auto-poster/ai"
	"xiaohongshu-auto-poster/config"
	"xiaohongshu-auto-poster/crawler"
	"xiaohongshu-auto-poster/monitor"
	"xiaohongshu-auto-poster/poster"
)

type Workflow struct {
	cron    *cron.Cron
	running bool
}

func NewWorkflow() *Workflow {
	return &Workflow{
		cron: cron.New(),
	}
}

func (w *Workflow) Start() error {
	if w.running {
		return nil
	}

	cfg := config.GetWorkflowConfig()

	// 添加定时任务
	if cfg.Schedule != "" {
		_, err := w.cron.AddFunc(cfg.Schedule, w.DailyTask)
		if err != nil {
			return err
		}
	}

	// 添加每周报告任务
	_, err := w.cron.AddFunc("0 0 * * 0", w.WeeklyReportTask)
	if err != nil {
		return err
	}

	w.cron.Start()
	w.running = true
	log.Println("工作流已启动")

	// 立即执行一次（可选）
	// go w.DailyTask()

	return nil
}

func (w *Workflow) Stop() {
	if w.running {
		w.cron.Stop()
		w.running = false
		log.Println("工作流已停止")
	}
}

func (w *Workflow) DailyTask() {
	log.Println("开始执行每日任务")

	cfg := config.GetWorkflowConfig()
	maxPosts := cfg.MaxPostsPerDay

	// 1. 获取热门话题
	var topics []string
	if cfg.EnableTopicCrawl {
		log.Println("开始爬取热门话题")
		hotTopics, err := crawler.GetHotTopics()
		if err != nil {
			log.Printf("爬取热门话题失败: %v", err)
		} else {
			// 提取话题标题
			for _, topic := range hotTopics {
				topics = append(topics, topic.Title)
			}
			log.Printf("成功获取%d个热门话题", len(topics))
		}
	}

	// 如果没有获取到热门话题，使用配置中的默认话题
	if len(topics) == 0 {
		topics = config.GetXiaohongshuConfig().Topics
		log.Printf("使用默认话题: %v", topics)
	}

	// 2. 生成内容
	for i := 0; i < maxPosts; i++ {
		// 随机选择话题
		if len(topics) == 0 {
			log.Println("没有可用话题，停止生成")
			break
		}

		topicIndex := rand.Intn(len(topics))
		topic := topics[topicIndex]

		// 移除已使用的话题，避免重复
		topics = append(topics[:topicIndex], topics[topicIndex+1:]...)

		err := generateAndPublishContent(topic)
		if err != nil {
			log.Printf("生成和发布内容失败: %v", err)
			continue
		}

		// 避免频率过高
		time.Sleep(5 * time.Minute)
	}

	log.Println("每日任务执行完成")
}

func generateAndPublishContent(topic string) error {
	log.Printf("开始为话题'%s'生成内容", topic)

	// 1. 生成文案
	log.Println("正在生成文案...")
	content, err := ai.GenerateContent(topic, detectTopicType(topic))
	if err != nil {
		return err
	}

	// 2. 生成标题
	log.Println("正在生成标题...")
	title, err := ai.GenerateTitle(topic)
	if err != nil {
		return err
	}

	// 3. 生成标签
	log.Println("正在生成标签...")
	tags, err := ai.GenerateTags(content, topic)
	if err != nil {
		return err
	}

	// 4. 生成图片
	var images []string
	if config.GetWorkflowConfig().EnableImageGen {
		log.Println("正在生成图片...")
		images, err = ai.GenerateImagesForContent(content, topic)
		if err != nil {
			log.Printf("生成图片失败: %v", err)
			// 图片生成失败仍继续发布
		}
	}

	// 5. 创建帖子
	post := poster.CreateDraftPost(title, content, images, tags)

	// 6. 保存草稿
	err = poster.SaveDraft(post)
	if err != nil {
		log.Printf("保存草稿失败: %v", err)
	}

	// 7. 发布帖子
	if config.GetWorkflowConfig().EnableAutoPublish {
		log.Println("正在发布帖子...")
		resp, err := poster.PublishPost(post)
		if err != nil {
			return err
		}

		if resp.Success {
			log.Printf("帖子发布成功，ID: %s", resp.PostID)
		} else {
			log.Printf("帖子发布失败: %s", resp.Message)
		}
	} else {
		log.Println("自动发布已禁用，仅保存草稿")
	}

	return nil
}

func detectTopicType(topic string) string {
	topic = strings.ToLower(topic)

	if strings.Contains(topic, "美食") || strings.Contains(topic, "吃") || strings.Contains(topic, "餐厅") || strings.Contains(topic, "烹饪") {
		return "food"
	} else if strings.Contains(topic, "旅行") || strings.Contains(topic, "旅游") || strings.Contains(topic, "景点") || strings.Contains(topic, "攻略") {
		return "travel"
	} else if strings.Contains(topic, "时尚") || strings.Contains(topic, "穿搭") || strings.Contains(topic, "服装") || strings.Contains(topic, "搭配") {
		return "fashion"
	} else if strings.Contains(topic, "美妆") || strings.Contains(topic, "护肤") || strings.Contains(topic, "化妆品") || strings.Contains(topic, "彩妆") {
		return "beauty"
	} else {
		return "food" // 默认美食模板
	}
}

func (w *Workflow) WeeklyReportTask() {
	log.Println("生成每周报告")

	err := monitor.SaveReport()
	if err != nil {
		log.Printf("生成周报失败: %v", err)
	} else {
		log.Println("周报生成成功")
	}
}

// 手动执行单次任务
func (w *Workflow) RunOnce(topic string) error {
	if topic == "" {
		// 随机选择话题
		topics := config.GetXiaohongshuConfig().Topics
		if len(topics) == 0 {
			return nil
		}
		topic = topics[rand.Intn(len(topics))]
	}

	return generateAndPublishContent(topic)
}

// 立即生成指定数量的内容
func (w *Workflow) GenerateBatch(count int, topics []string) error {
	if count <= 0 {
		count = 1
	}

	if len(topics) == 0 {
		topics = config.GetXiaohongshuConfig().Topics
	}

	if len(topics) == 0 {
		return nil
	}

	for i := 0; i < count; i++ {
		// 循环使用话题
		topicIndex := i % len(topics)
		topic := topics[topicIndex]

		err := generateAndPublishContent(topic)
		if err != nil {
			log.Printf("生成内容失败: %v", err)
			// 继续生成下一个
		}

		time.Sleep(2 * time.Minute)
	}

	return nil
}

// 获取工作流状态
func (w *Workflow) IsRunning() bool {
	return w.running
}

// 紧急停止所有任务
func (w *Workflow) EmergencyStop() {
	if w.running {
		w.cron.Stop()
		w.running = false
		log.Println("工作流已紧急停止")
	}
}

// 重新加载配置
func (w *Workflow) ReloadConfig() error {
	// 这里可以实现配置热重载
	log.Println("重新加载配置")
	return nil
}

// 获取统计信息
func (w *Workflow) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["running"] = w.running
	stats["entries"] = len(w.cron.Entries())

	if w.running {
		nextRun := w.cron.Entries()[0].Next
		stats["next_run"] = nextRun.Format("2006-01-02 15:04:05")
	}

	return stats
}
