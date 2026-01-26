package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type DoubaoConfig struct {
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
}

type XiaohongshuConfig struct {
	Cookie    string   `yaml:"cookie"`
	UserAgent string   `yaml:"user_agent"`
	Topics    []string `yaml:"topics"`
}

type ImageGeneratorConfig struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
	Style    string `yaml:"style"`
}

type WorkflowConfig struct {
	Schedule          string `yaml:"schedule"`
	MaxPostsPerDay    int    `yaml:"max_posts_per_day"`
	EnableAutoPublish bool   `yaml:"enable_auto_publish"`
	EnableImageGen    bool   `yaml:"enable_image_generation"`
	EnableTopicCrawl  bool   `yaml:"enable_topic_crawling"`
}

type PromptTemplate struct {
	Food   string `yaml:"food"`
	Travel string `yaml:"travel"`
}

type KeywordsConfig struct {
	Food   []string `yaml:"food"`
	Travel []string `yaml:"travel"`
}

type OutputConfig struct {
	SaveDir       string `yaml:"save_dir"`
	Format        string `yaml:"format"`
	IncludeImages bool   `yaml:"include_images"`
}

type Config struct {
	Doubao          DoubaoConfig         `yaml:"doubao"`
	Xiaohongshu     XiaohongshuConfig    `yaml:"xiaohongshu"`
	ImageGenerator  ImageGeneratorConfig `yaml:"image_generator"`
	Workflow        WorkflowConfig       `yaml:"workflow"`
	PromptTemplates PromptTemplate       `yaml:"prompt_templates"`
	Keywords        KeywordsConfig       `yaml:"keywords"`
	Output          OutputConfig         `yaml:"output"`
}

var GlobalConfig *Config

func LoadConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		log.Printf("无法打开配置文件: %v", err)
		return err
	}
	defer file.Close()

	var config Config
	yamlDecoder := yaml.NewDecoder(file)
	if err := yamlDecoder.Decode(&config); err != nil {
		log.Printf("解析配置文件失败: %v", err)
		return err
	}

	GlobalConfig = &config
	log.Println("配置文件加载成功")
	return nil
}

func GetDoubaoConfig() DoubaoConfig {
	if GlobalConfig == nil {
		return DoubaoConfig{}
	}
	return GlobalConfig.Doubao
}

func GetXiaohongshuConfig() XiaohongshuConfig {
	if GlobalConfig == nil {
		return XiaohongshuConfig{}
	}
	return GlobalConfig.Xiaohongshu
}

func GetWorkflowConfig() WorkflowConfig {
	if GlobalConfig == nil {
		return WorkflowConfig{}
	}
	return GlobalConfig.Workflow
}

func GetPromptTemplate(topicType string) string {
	if GlobalConfig == nil {
		return ""
	}

	switch topicType {
	case "food":
		return GlobalConfig.PromptTemplates.Food
	case "travel":
		return GlobalConfig.PromptTemplates.Travel
	default:
		return GlobalConfig.PromptTemplates.Food
	}
}

func GetKeywords(topicType string) []string {
	if GlobalConfig == nil {
		return []string{}
	}

	switch topicType {
	case "food":
		return GlobalConfig.Keywords.Food
	case "travel":
		return GlobalConfig.Keywords.Travel
	default:
		return GlobalConfig.Keywords.Food
	}
}
