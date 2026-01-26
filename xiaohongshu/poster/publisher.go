package poster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"xiaohongshu-auto-poster/config"
)

type Post struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Images      []string  `json:"images"`
	Tags        []string  `json:"tags"`
	Status      string    `json:"status"` // draft, published, failed
	PublishTime time.Time `json:"publish_time"`
	Platform    string    `json:"platform"` // xiaohongshu
}

type PublishResponse struct {
	Success bool   `json:"success"`
	PostID  string `json:"post_id"`
	Message string `json:"message"`
}

func PublishPost(post *Post) (*PublishResponse, error) {
	if !config.GetWorkflowConfig().EnableAutoPublish {
		return &PublishResponse{
			Success: false,
			Message: "自动发布功能已禁用",
		}, nil
	}

	// 上传图片
	var imageIDs []string
	for _, imageURL := range post.Images {
		imageID, err := uploadImage(imageURL)
		if err != nil {
			log.Printf("上传图片失败: %v", err)
			continue
		}
		imageIDs = append(imageIDs, imageID)
	}

	// 发布笔记
	return createNote(post, imageIDs)
}

func uploadImage(imageURL string) (string, error) {
	cfg := config.GetXiaohongshuConfig()
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// 下载图片
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	// 准备上传请求
	uploadURL := "https://www.xiaohongshu.com/api/sns/v1/upload/image"
	data := url.Values{}
	data.Set("source", "web")
	data.Set("type", "note")

	req, err := http.NewRequest("POST", uploadURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Cookie", cfg.Cookie)
	req.Header.Set("User-Agent", cfg.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 实际上传需要multipart/form-data格式
	// 这里简化处理，实际实现需要正确处理文件上传
	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("上传图片失败: %s", string(body))
	}

	var result struct {
		Success bool   `json:"success"`
		ImageID string `json:"image_id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if result.Success {
		return result.ImageID, nil
	}

	return "", fmt.Errorf("上传图片失败")
}

func createNote(post *Post, imageIDs []string) (*PublishResponse, error) {
	cfg := config.GetXiaohongshuConfig()
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	noteURL := "https://www.xiaohongshu.com/api/sns/v3/notes"

	tagStr := strings.Join(post.Tags, ",")
	tagStr = strings.ReplaceAll(tagStr, "#", "")

	requestBody := map[string]interface{}{
		"title":            post.Title,
		"content":          post.Content,
		"image_ids":        imageIDs,
		"tags":             strings.Split(tagStr, ","),
		"note_type":        "normal",
		"public_type":      0, // 公开
		"sync_to_facebook": false,
		"sync_to_weibo":    false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", noteURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", cfg.Cookie)
	req.Header.Set("User-Agent", cfg.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://www.xiaohongshu.com/publish")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Success bool   `json:"success"`
		NoteID  string `json:"note_id"`
		Message string `json:"message"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Success {
		post.ID = result.NoteID
		post.Status = "published"
		post.PublishTime = time.Now()

		// 保存发布记录
		savePostRecord(post)

		return &PublishResponse{
			Success: true,
			PostID:  result.NoteID,
			Message: "发布成功",
		}, nil
	}

	return &PublishResponse{
		Success: false,
		Message: result.Message,
	}, fmt.Errorf("发布失败: %s", result.Message)
}

func savePostRecord(post *Post) error {
	cfg := config.GlobalConfig
	if cfg == nil {
		return fmt.Errorf("配置未初始化")
	}

	saveDir := cfg.Output.SaveDir
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("post_%s_%s.json", post.Platform, timestamp)
	filePath := filepath.Join(saveDir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(post)
}

func SaveDraft(post *Post) error {
	cfg := config.GlobalConfig
	if cfg == nil {
		return fmt.Errorf("配置未初始化")
	}

	saveDir := filepath.Join(cfg.Output.SaveDir, "drafts")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("draft_%s.json", timestamp)
	filePath := filepath.Join(saveDir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(post)
}

func GetPostStatus(postID string) (string, error) {
	cfg := config.GetXiaohongshuConfig()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	statusURL := fmt.Sprintf("https://www.xiaohongshu.com/api/sns/v3/notes/%s/status", postID)

	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Cookie", cfg.Cookie)
	req.Header.Set("User-Agent", cfg.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Status, nil
}

func GetPublishedPosts() ([]Post, error) {
	cfg := config.GetXiaohongshuConfig()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	postsURL := "https://www.xiaohongshu.com/api/sns/v3/user/notes"

	req, err := http.NewRequest("GET", postsURL, nil)
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

	var result struct {
		Data []Post `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func CreateDraftPost(title, content string, images []string, tags []string) *Post {
	return &Post{
		Title:    title,
		Content:  content,
		Images:   images,
		Tags:     tags,
		Status:   "draft",
		Platform: "xiaohongshu",
	}
}
