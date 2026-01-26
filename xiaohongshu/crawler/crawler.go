package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"xiaohongshu-auto-poster/config"
)

type Topic struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	ImageURL     string    `json:"image_url"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	ShareCount   int       `json:"share_count"`
	PostTime     time.Time `json:"post_time"`
	Author       string    `json:"author"`
	Tags         []string  `json:"tags"`
}

type HotTopicResponse struct {
	Code    int     `json:"code"`
	Data    []Topic `json:"data"`
	Message string  `json:"message"`
}

func GetHotTopics() ([]Topic, error) {
	cfg := config.GetXiaohongshuConfig()
	var allTopics []Topic

	for _, topic := range cfg.Topics {
		topics, err := crawlTopic(topic)
		if err != nil {
			log.Printf("爬取话题%s失败: %v", topic, err)
			continue
		}
		allTopics = append(allTopics, topics...)
	}

	return allTopics, nil
}

func crawlTopic(keyword string) ([]Topic, error) {
	cfg := config.GetXiaohongshuConfig()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 构造请求URL
	baseURL := "https://www.xiaohongshu.com/api/sns/v3/search/notes"
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("page_size", "20")
	params.Set("sort", "hot")
	params.Set("page", "1")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Cookie", cfg.Cookie)
	req.Header.Set("User-Agent", cfg.UserAgent)
	req.Header.Set("Referer", "https://www.xiaohongshu.com/")
	req.Header.Set("Origin", "https://www.xiaohongshu.com")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	var response HotTopicResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("API返回错误: %s", response.Message)
	}

	return response.Data, nil
}

func GetTopicDetail(noteID string) (*Topic, error) {
	cfg := config.GetXiaohongshuConfig()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	url := fmt.Sprintf("https://www.xiaohongshu.com/api/sns/v3/notes/%s", noteID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Cookie", cfg.Cookie)
	req.Header.Set("User-Agent", cfg.UserAgent)
	req.Header.Set("Referer", "https://www.xiaohongshu.com/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	var topic Topic
	err = json.NewDecoder(resp.Body).Decode(&topic)
	if err != nil {
		return nil, err
	}

	return &topic, nil
}

// 使用网页解析方式爬取热门话题（备用方案）
func crawlHotTopicsWeb() ([]Topic, error) {
	cfg := config.GetXiaohongshuConfig()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://www.xiaohongshu.com/discovery", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", cfg.Cookie)
	req.Header.Set("User-Agent", cfg.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var topics []Topic

	// 解析热门话题
	doc.Find(".note-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".note-title").Text()
		content := s.Find(".note-content").Text()
		imageURL, _ := s.Find(".note-image").Attr("src")
		author := s.Find(".author-name").Text()

		// 解析点赞数等信息
		likeStr := s.Find(".like-count").Text()
		likeCount := parseNumber(likeStr)

		topic := Topic{
			Title:     strings.TrimSpace(title),
			Content:   strings.TrimSpace(content),
			ImageURL:  imageURL,
			LikeCount: likeCount,
			Author:    strings.TrimSpace(author),
		}

		topics = append(topics, topic)
	})

	return topics, nil
}

func parseNumber(numStr string) int {
	numStr = strings.ReplaceAll(numStr, "万", "0000")
	numStr = strings.ReplaceAll(numStr, "千", "000")
	numStr = strings.ReplaceAll(numStr, ".", "")
	numStr = strings.TrimSpace(numStr)

	var num int
	fmt.Sscanf(numStr, "%d", &num)
	return num
}

func GetTopicTrends() ([]string, error) {
	// 获取话题趋势
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://www.xiaohongshu.com/api/sns/v3/trending/topics", nil)
	if err != nil {
		return nil, err
	}

	cfg := config.GetXiaohongshuConfig()
	req.Header.Set("Cookie", cfg.Cookie)
	req.Header.Set("User-Agent", cfg.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	var result struct {
		Code int      `json:"code"`
		Data []string `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
