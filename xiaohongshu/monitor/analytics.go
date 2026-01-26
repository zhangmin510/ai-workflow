package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"xiaohongshu-auto-poster/config"
	"xiaohongshu-auto-poster/poster"
)

type PostStats struct {
	PostID         string    `json:"post_id"`
	Title          string    `json:"title"`
	PublishTime    time.Time `json:"publish_time"`
	LikeCount      int       `json:"like_count"`
	CommentCount   int       `json:"comment_count"`
	ShareCount     int       `json:"share_count"`
	CollectCount   int       `json:"collect_count"`
	ViewCount      int       `json:"view_count"`
	EngagementRate float64   `json:"engagement_rate"`
}

type DailyStats struct {
	Date          time.Time `json:"date"`
	TotalPosts    int       `json:"total_posts"`
	TotalLikes    int       `json:"total_likes"`
	TotalComments int       `json:"total_comments"`
	TotalShares   int       `json:"total_shares"`
	AvgEngagement float64   `json:"avg_engagement"`
}

type TopicPerformance struct {
	Topic          string  `json:"topic"`
	PostCount      int     `json:"post_count"`
	AvgLikes       float64 `json:"avg_likes"`
	AvgComments    float64 `json:"avg_comments"`
	EngagementRate float64 `json:"engagement_rate"`
}

func GetPostStats(postID string) (*PostStats, error) {
	// 模拟获取帖子数据
	// 实际需要调用小红书API获取真实数据
	stats := &PostStats{
		PostID:       postID,
		PublishTime:  time.Now().Add(-24 * time.Hour),
		LikeCount:    156,
		CommentCount: 23,
		ShareCount:   15,
		CollectCount: 42,
		ViewCount:    2580,
	}

	// 计算互动率
	if stats.ViewCount > 0 {
		stats.EngagementRate = float64(stats.LikeCount+stats.CommentCount+stats.ShareCount) / float64(stats.ViewCount) * 100
	}

	return stats, nil
}

func GetAllPostStats() ([]PostStats, error) {
	// 从保存的帖子记录中获取统计数据
	cfg := config.GlobalConfig
	if cfg == nil {
		return nil, fmt.Errorf("配置未初始化")
	}

	saveDir := cfg.Output.SaveDir
	files, err := filepath.Glob(filepath.Join(saveDir, "post_xiaohongshu_*.json"))
	if err != nil {
		return nil, err
	}

	var allStats []PostStats
	for _, file := range files {
		post, err := loadPostFromFile(file)
		if err != nil {
			log.Printf("加载帖子文件失败: %v", err)
			continue
		}

		stats, err := GetPostStats(post.ID)
		if err != nil {
			log.Printf("获取帖子统计失败: %v", err)
			continue
		}

		stats.Title = post.Title
		stats.PublishTime = post.PublishTime
		allStats = append(allStats, *stats)
	}

	// 按发布时间排序
	sort.Slice(allStats, func(i, j int) bool {
		return allStats[i].PublishTime.After(allStats[j].PublishTime)
	})

	return allStats, nil
}

func loadPostFromFile(filePath string) (*poster.Post, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var post poster.Post
	err = json.NewDecoder(file).Decode(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func GetDailyStats(days int) ([]DailyStats, error) {
	allStats, err := GetAllPostStats()
	if err != nil {
		return nil, err
	}

	dailyMap := make(map[string]*DailyStats)

	for _, stats := range allStats {
		// 按日期分组
		dateStr := stats.PublishTime.Format("2006-01-02")

		if dailyMap[dateStr] == nil {
			dailyMap[dateStr] = &DailyStats{
				Date: stats.PublishTime,
			}
		}

		dailyStats := dailyMap[dateStr]
		dailyStats.TotalPosts++
		dailyStats.TotalLikes += stats.LikeCount
		dailyStats.TotalComments += stats.CommentCount
		dailyStats.TotalShares += stats.ShareCount
	}

	// 计算平均值
	for _, stats := range dailyMap {
		if stats.TotalPosts > 0 {
			stats.AvgEngagement = float64(stats.TotalLikes+stats.TotalComments+stats.TotalShares) / float64(stats.TotalPosts) / 10
		}
	}

	// 转换为切片并排序
	var dailyStatsList []DailyStats
	for _, stats := range dailyMap {
		dailyStatsList = append(dailyStatsList, *stats)
	}

	sort.Slice(dailyStatsList, func(i, j int) bool {
		return dailyStatsList[i].Date.After(dailyStatsList[j].Date)
	})

	// 限制返回天数
	if days > 0 && len(dailyStatsList) > days {
		dailyStatsList = dailyStatsList[:days]
	}

	return dailyStatsList, nil
}

func GetTopicPerformance() ([]TopicPerformance, error) {
	allStats, err := GetAllPostStats()
	if err != nil {
		return nil, err
	}

	topicMap := make(map[string]*TopicPerformance)

	// 简单通过标题关键词判断话题类型
	for _, stats := range allStats {
		var topic string
		if strings.Contains(stats.Title, "美食") || strings.Contains(stats.Title, "吃") || strings.Contains(stats.Title, "餐厅") {
			topic = "美食"
		} else if strings.Contains(stats.Title, "旅行") || strings.Contains(stats.Title, "旅游") || strings.Contains(stats.Title, "景点") {
			topic = "旅行"
		} else if strings.Contains(stats.Title, "时尚") || strings.Contains(stats.Title, "穿搭") || strings.Contains(stats.Title, "服装") {
			topic = "时尚"
		} else if strings.Contains(stats.Title, "美妆") || strings.Contains(stats.Title, "护肤") || strings.Contains(stats.Title, "化妆品") {
			topic = "美妆"
		} else {
			topic = "其他"
		}

		if topicMap[topic] == nil {
			topicMap[topic] = &TopicPerformance{
				Topic: topic,
			}
		}

		tp := topicMap[topic]
		tp.PostCount++
		tp.AvgLikes += float64(stats.LikeCount)
		tp.AvgComments += float64(stats.CommentCount)
		tp.EngagementRate += stats.EngagementRate
	}

	// 计算平均值
	for _, tp := range topicMap {
		if tp.PostCount > 0 {
			tp.AvgLikes /= float64(tp.PostCount)
			tp.AvgComments /= float64(tp.PostCount)
			tp.EngagementRate /= float64(tp.PostCount)
		}
	}

	// 转换为切片并排序
	var topicList []TopicPerformance
	for _, tp := range topicMap {
		topicList = append(topicList, *tp)
	}

	sort.Slice(topicList, func(i, j int) bool {
		return topicList[i].EngagementRate > topicList[j].EngagementRate
	})

	return topicList, nil
}

func GenerateReport() (string, error) {
	var report strings.Builder

	report.WriteString("📊 小红书账号表现报告\n")
	report.WriteString("=========================\n\n")

	// 概览
	report.WriteString("📈 整体概览\n")
	report.WriteString("----------\n")

	allStats, err := GetAllPostStats()
	if err != nil {
		return "", err
	}

	if len(allStats) == 0 {
		report.WriteString("暂无数据\n")
		return report.String(), nil
	}

	totalLikes := 0
	totalComments := 0
	totalShares := 0
	totalViews := 0

	for _, stats := range allStats {
		totalLikes += stats.LikeCount
		totalComments += stats.CommentCount
		totalShares += stats.ShareCount
		totalViews += stats.ViewCount
	}

	avgEngagement := float64(totalLikes+totalComments+totalShares) / float64(totalViews) * 100

	report.WriteString(fmt.Sprintf("总帖子数: %d\n", len(allStats)))
	report.WriteString(fmt.Sprintf("总点赞数: %d\n", totalLikes))
	report.WriteString(fmt.Sprintf("总评论数: %d\n", totalComments))
	report.WriteString(fmt.Sprintf("总分享数: %d\n", totalShares))
	report.WriteString(fmt.Sprintf("总浏览数: %d\n", totalViews))
	report.WriteString(fmt.Sprintf("平均互动率: %.2f%%\n\n", avgEngagement))

	// 每日表现
	report.WriteString("📅 近7天表现\n")
	report.WriteString("------------\n")

	dailyStats, err := GetDailyStats(7)
	if err != nil {
		return "", err
	}

	for _, stats := range dailyStats {
		report.WriteString(fmt.Sprintf("%s: ", stats.Date.Format("01-02")))
		report.WriteString(fmt.Sprintf("发帖%d篇, ", stats.TotalPosts))
		report.WriteString(fmt.Sprintf("互动%.2f\n", stats.AvgEngagement))
	}

	// 话题表现
	report.WriteString("\n🏷️ 话题表现\n")
	report.WriteString("------------\n")

	topicStats, err := GetTopicPerformance()
	if err != nil {
		return "", err
	}

	for _, tp := range topicStats {
		report.WriteString(fmt.Sprintf("%s: ", tp.Topic))
		report.WriteString(fmt.Sprintf("发帖%d篇, ", tp.PostCount))
		report.WriteString(fmt.Sprintf("平均点赞%.0f, ", tp.AvgLikes))
		report.WriteString(fmt.Sprintf("互动率%.2f%%\n", tp.EngagementRate))
	}

	// 最佳帖子
	report.WriteString("\n🏆 最佳表现\n")
	report.WriteString("------------\n")

	// 按互动率排序
	sort.Slice(allStats, func(i, j int) bool {
		return allStats[i].EngagementRate > allStats[j].EngagementRate
	})

	for i, stats := range allStats[:min(3, len(allStats))] {
		report.WriteString(fmt.Sprintf("第%d名: %s\n", i+1, stats.Title))
		report.WriteString(fmt.Sprintf("  互动率: %.2f%%, 点赞: %d, 评论: %d\n",
			stats.EngagementRate, stats.LikeCount, stats.CommentCount))
	}

	return report.String(), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SaveReport() error {
	report, err := GenerateReport()
	if err != nil {
		return err
	}

	cfg := config.GlobalConfig
	if cfg == nil {
		return fmt.Errorf("配置未初始化")
	}

	saveDir := cfg.Output.SaveDir
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("report_%s.txt", timestamp)
	filePath := filepath.Join(saveDir, filename)

	err = os.WriteFile(filePath, []byte(report), 0644)
	if err != nil {
		return err
	}

	log.Printf("报告已保存至: %s", filePath)
	return nil
}

func GetEngagementTrend() ([]float64, error) {
	dailyStats, err := GetDailyStats(30)
	if err != nil {
		return nil, err
	}

	// 按日期升序排列
	sort.Slice(dailyStats, func(i, j int) bool {
		return dailyStats[i].Date.Before(dailyStats[j].Date)
	})

	var trend []float64
	for _, stats := range dailyStats {
		trend = append(trend, stats.AvgEngagement)
	}

	return trend, nil
}
