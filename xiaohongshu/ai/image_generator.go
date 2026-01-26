package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"xiaohongshu-auto-poster/config"
)

type ImageGenerator interface {
	GenerateImage(prompt string) (string, error)
}

type DALL_EGenerator struct {
	apiKey string
}

type StableDiffusionGenerator struct {
	apiKey string
}

type MidjourneyGenerator struct {
	apiKey string
}

type ImageResponse struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func NewImageGenerator() (ImageGenerator, error) {
	cfg := config.GlobalConfig
	if cfg == nil {
		return nil, fmt.Errorf("配置未初始化")
	}

	switch cfg.ImageGenerator.Provider {
	case "dall-e":
		return &DALL_EGenerator{apiKey: cfg.ImageGenerator.APIKey}, nil
	case "stable-diffusion":
		return &StableDiffusionGenerator{apiKey: cfg.ImageGenerator.APIKey}, nil
	case "midjourney":
		return &MidjourneyGenerator{apiKey: cfg.ImageGenerator.APIKey}, nil
	default:
		return nil, fmt.Errorf("不支持的图片生成器: %s", cfg.ImageGenerator.Provider)
	}
}

func (d *DALL_EGenerator) GenerateImage(prompt string) (string, error) {
	cfg := config.GlobalConfig
	if cfg == nil {
		return "", fmt.Errorf("配置未初始化")
	}

	// 添加小红书风格前缀
	fullPrompt := fmt.Sprintf("%s, %s", prompt, cfg.ImageGenerator.Style)

	requestBody := map[string]interface{}{
		"model":  "dall-e-3",
		"prompt": fullPrompt,
		"n":      1,
		"size":   "1024x1792", // 小红书推荐比例 9:16
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	authHeader := fmt.Sprintf("Bearer %s", d.apiKey)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var imageResponse ImageResponse
	err = json.Unmarshal(body, &imageResponse)
	if err != nil {
		return "", err
	}

	if imageResponse.Error.Message != "" {
		return "", fmt.Errorf("DALL-E API错误: %s", imageResponse.Error.Message)
	}

	if len(imageResponse.Data) > 0 {
		return imageResponse.Data[0].URL, nil
	}

	return "", fmt.Errorf("未生成图片")
}

func (s *StableDiffusionGenerator) GenerateImage(prompt string) (string, error) {
	cfg := config.GlobalConfig
	if cfg == nil {
		return "", fmt.Errorf("配置未初始化")
	}

	fullPrompt := fmt.Sprintf("%s, %s, 8k, 高清, 细节丰富", prompt, cfg.ImageGenerator.Style)

	requestBody := map[string]interface{}{
		"prompt":          fullPrompt,
		"negative_prompt": "低质量, 模糊, 变形, 丑陋, 不自然",
		"width":           1024,
		"height":          1792,
		"steps":           30,
		"cfg_scale":       7.5,
		"sampler_index":   "DPM++ 2M Karras",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 假设使用Stable Diffusion WebUI API
	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	req, err := http.NewRequest("POST", "http://localhost:7860/sdapi/v1/txt2img", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("请求失败: %s", string(body))
	}

	var response struct {
		Images []string `json:"images"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	if len(response.Images) > 0 {
		// 返回base64编码的图片
		return response.Images[0], nil
	}

	return "", fmt.Errorf("未生成图片")
}

func (m *MidjourneyGenerator) GenerateImage(prompt string) (string, error) {
	// Midjourney通常通过Discord机器人调用
	// 这里简化处理，实际需要集成Discord API
	cfg := config.GlobalConfig
	if cfg == nil {
		return "", fmt.Errorf("配置未初始化")
	}

	fullPrompt := fmt.Sprintf("/imagine prompt: %s, %s --ar 9:16 --style raw --q 2", prompt, cfg.ImageGenerator.Style)

	// 这里应该调用Discord API发送消息
	log.Printf("Midjourney生成命令: %s", fullPrompt)

	// 实际场景中需要等待图片生成完成并获取图片URL
	// 这里返回一个示例URL
	return "https://example.com/midjourney-image.jpg", nil
}

func GenerateImagesForContent(content string, topic string) ([]string, error) {
	generator, err := NewImageGenerator()
	if err != nil {
		return nil, err
	}

	// 从内容中提取关键元素生成图片提示词
	prompts, err := generateImagePrompts(content, topic)
	if err != nil {
		return nil, err
	}

	var imageURLs []string
	for _, prompt := range prompts {
		imageURL, err := generator.GenerateImage(prompt)
		if err != nil {
			log.Printf("生成图片失败: %v", err)
			continue
		}
		imageURLs = append(imageURLs, imageURL)

		// 避免API调用过于频繁
		time.Sleep(2 * time.Second)
	}

	return imageURLs, nil
}

func generateImagePrompts(content string, topic string) ([]string, error) {
	cfg := config.GetDoubaoConfig()

	prompt := fmt.Sprintf(`你是一个图片提示词专家，擅长创作适合AI生成的高质量提示词。请为以下小红书笔记内容生成3-5个图片生成提示词：
内容：%s
主题：%s

要求：
1. 每个提示词20-50字
2. 具体描述画面内容
3. 符合小红书审美风格
4. 包含细节和氛围描述
5. 适合生成9:16比例的图片`, content, topic)

	messages := []Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	result, err := callDoubaoAPI(cfg, messages, 0.8, 500)
	if err != nil {
		return nil, err
	}

	return parsePrompts(result), nil
}

func parsePrompts(promptContent string) []string {
	var prompts []string
	lines := strings.Split(promptContent, "\n")

	for _, line := range lines {
		line := strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "-") {
			// 如果不是以-开头，可能是编号格式
			if idx := strings.Index(line, "."); idx != -1 {
				line = strings.TrimSpace(line[idx+1:])
			}
			prompts = append(prompts, line)
		} else if strings.HasPrefix(line, "-") {
			prompt := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			if prompt != "" {
				prompts = append(prompts, prompt)
			}
		}
	}

	return prompts
}
