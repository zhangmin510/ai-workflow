package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"xiaohongshu-auto-poster/config"
)

// 火山引擎豆包模型API响应结构体
type DoubaoResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
	} `json:"error"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 获取火山引擎豆包模型的access_token（火山引擎API使用API Key直接认证，无需获取token）
func getAccessToken(cfg config.DoubaoConfig) (string, error) {
	// 火山引擎API直接使用API Key进行认证，无需获取access_token
	// 此函数保持兼容，直接返回空字符串和nil错误
	return "", nil
}

// 调用火山引擎豆包模型API
type DoubaoChatRequest struct {
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

func callDoubaoAPI(cfg config.DoubaoConfig, messages []Message, temperature float64, maxTokens int) (string, error) {
	// 构建请求体
	requestBody := map[string]interface{}{
		"model":       cfg.Model,
		"messages":    messages,
		"temperature": temperature,
		"max_tokens":  maxTokens,
		"top_p":       1.0,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequest("POST", "https://ark.cn-beijing.volces.com/api/v3/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.APIKey))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var doubaoResponse DoubaoResponse
	err = json.Unmarshal(body, &doubaoResponse)
	if err != nil {
		return "", err
	}

	if doubaoResponse.Error.Code != 0 {
		return "", fmt.Errorf("火山引擎豆包API错误: %d - %s", doubaoResponse.Error.Code, doubaoResponse.Error.Message)
	}

	if len(doubaoResponse.Choices) == 0 || doubaoResponse.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("未获取到生成内容")
	}

	return doubaoResponse.Choices[0].Message.Content, nil
}

func GenerateContent(topic string, topicType string) (string, error) {
	cfg := config.GetDoubaoConfig()
	promptTemplate := config.GetPromptTemplate(topicType)

	// 替换模板中的变量
	prompt := fmt.Sprintf(promptTemplate, topic)

	messages := []Message{
		{
			Role:    "user",
			Content: "你是一个小红书文案专家，擅长创作各种类型的笔记内容。" + prompt,
		},
	}

	return callDoubaoAPI(cfg, messages, cfg.Temperature, cfg.MaxTokens)
}

func GenerateTitle(topic string) (string, error) {
	cfg := config.GetDoubaoConfig()

	prompt := fmt.Sprintf(`你是一个小红书标题专家，擅长创作高点击率的标题。请为关于"%s"的小红书笔记生成5个吸引人的标题，要求：
1. 每个标题不超过20字
2. 使用emoji增强吸引力
3. 突出核心亮点
4. 符合小红书用户喜好

例子：
"🍜上海最正宗的兰州拉面！一口穿越到大西北"`, topic)

	messages := []Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return callDoubaoAPI(cfg, messages, 0.9, 200)
}

func GenerateTags(content string, topic string) ([]string, error) {
	cfg := config.GetDoubaoConfig()

	prompt := fmt.Sprintf(`你是一个小红书标签专家，擅长选择高流量的话题标签。请为以下小红书笔记内容生成8-10个相关的话题标签：
内容：%s
主题：%s

要求：
1. 使用#开头
2. 包含1-2个大话题（如#美食）
3. 包含3-4个中话题（如#上海美食）
4. 包含3-4个精准话题（如#上海正宗拉面）
5. 符合小红书话题标签规范`, content, topic)

	messages := []Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	result, err := callDoubaoAPI(cfg, messages, 0.7, 200)
	if err != nil {
		return nil, err
	}

	return parseTags(result), nil
}

func parseTags(tagContent string) []string {
	var tags []string
	lines := strings.Split(tagContent, "\n")

	for _, line := range lines {
		line := strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			tag := strings.TrimPrefix(line, "#")
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, "#"+tag)
			}
		}
	}

	return tags
}

func OptimizeContent(content string) (string, error) {
	cfg := config.GetDoubaoConfig()

	prompt := fmt.Sprintf(`你是一个小红书文案优化专家，擅长提升笔记的吸引力。请优化以下小红书笔记内容，使其更符合小红书平台风格：
%s

优化要求：
1. 语言更生动活泼，使用网络流行语
2. 增加表情符号，增强可读性
3. 结构更清晰，分段合理
4. 突出重点内容，吸引读者
5. 保持原意不变`, content)

	messages := []Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return callDoubaoAPI(cfg, messages, 0.8, 1500)
}
