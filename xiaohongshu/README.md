# 小红书自动创作工作流

一个基于Go语言开发的小红书内容自动创作与发布系统，集成AI文案生成、图片生成、热门话题爬取和自动发布功能。

## ✨ 功能特性

### 🎯 核心功能
- **AI智能文案生成**: 基于OpenAI GPT模型生成小红书风格文案
- **热门话题自动爬取**: 实时获取小红书平台热门话题，紧跟流量趋势
- **AI图片生成**: 集成DALL-E、Stable Diffusion、Midjourney等图片生成器
- **自动发布**: 支持一键发布或定时自动发布到小红书
- **数据分析**: 生成账号表现报告，分析内容效果和用户互动

### 📊 管理功能
- **多模式运行**: 支持服务模式、单次运行、批量生成三种模式
- **RESTful API**: 提供完整的API接口，方便集成到其他系统
- **优雅关闭**: 支持系统信号监听，实现优雅关闭
- **配置热重载**: 支持在线修改配置，无需重启服务

## 🚀 快速开始

### 环境要求
- Go >= 1.25.0
- 有效的OpenAI API Key
- 小红书账号Cookie（可选，用于发布）

### 安装依赖

```bash
# 安装Go依赖
go get github.com/gin-gonic/gin
go get gopkg.in/yaml.v3
go get github.com/robfig/cron/v3
go get github.com/PuerkitoBio/goquery
```

### 配置文件

复制并编辑配置文件：

```bash
cp config/config.yaml.example config/config.yaml
```

主要配置项：

```yaml
# OpenAI配置
openai:
  api_key: "your_openai_api_key_here"
  model: "gpt-4o-mini"
  temperature: 0.8

# 小红书配置
xiaohongshu:
  cookie: "your_xiaohongshu_cookie_here"
  topics: ["美食", "旅行", "时尚", "美妆", "健身"]

# 工作流配置
workflow:
  schedule: "0 9 * * *"  # 每天早上9点执行
  max_posts_per_day: 3
  enable_auto_publish: false  # 先设为false预览效果
```

### 运行方式

#### 1. 服务模式（推荐）

```bash
go run main.go --mode server
```

服务启动后，访问 http://localhost:8080 查看API文档

#### 2. 单次运行模式

```bash
# 使用随机话题生成一篇内容
go run main.go --mode once

# 指定话题生成
go run main.go --mode once --topic "上海美食攻略"
```

#### 3. 批量生成模式

```bash
# 批量生成5篇内容
go run main.go --mode batch --count 5
```

## 📖 API文档

### 基础信息
- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`

### 主要接口

#### 获取工作流状态
```http
GET /api/status
```

#### 启动工作流
```http
POST /api/start
```

#### 停止工作流
```http
POST /api/stop
```

#### 手动执行一次任务
```http
POST /api/run-once
Content-Type: application/json

{
  "topic": "上海美食攻略"
}
```

#### 批量生成内容
```http
POST /api/generate-batch
Content-Type: application/json

{
  "count": 5,
  "topics": ["美食", "旅行", "时尚"]
}
```

## 🔧 架构设计

### 项目结构

```
xiaohongshu/
├── config/          # 配置管理模块
│   ├── config.go    # 配置加载和解析
│   └── config.yaml  # YAML配置文件
├── crawler/         # 爬虫模块
│   └── crawler.go   # 小红书话题爬取
├── ai/              # AI相关模块
│   ├── content_generator.go  # 文案生成
│   └── image_generator.go    # 图片生成
├── poster/          # 发布模块
│   └── publisher.go # 小红书发布功能
├── monitor/         # 监控和分析模块
│   └── analytics.go # 数据统计和报告
├── workflow/        # 工作流调度
│   └── scheduler.go # 定时任务管理
└── main.go          # 主入口文件
```

### 工作流程

```
┌───────────────────┐    ┌───────────────────┐    ┌───────────────────┐
│   热门话题爬取    │───▶│   AI文案生成      │───▶│   AI图片生成      │
└───────────────────┘    └───────────────────┘    └───────────────────┘
          │                       │                       │
          ▼                       ▼                       ▼
┌───────────────────┐    ┌───────────────────┐    ┌───────────────────┐
│   数据清洗处理    │───▶│   内容优化审核    │───▶│   自动发布到小红书 │
└───────────────────┘    └───────────────────┘    └───────────────────┘
          │
          ▼
┌───────────────────┐
│   数据分析报告    │
└───────────────────┘
```

## 📈 使用指南

### 最佳实践

1. **内容策略**：
   - 先启用`enable_auto_publish: false`预览效果
   - 调整`temperature`参数控制文案创造性
   - 使用不同的提示词模板生成多样化内容

2. **图片优化**：
   - 使用9:16比例适配小红书移动端
   - 加入"小红书风格, 明亮色调, 高清画质"等关键词
   - 可同时生成多张图片选择最佳效果

3. **发布时机**：
   - 早上9点、中午12点、晚上8点为最佳发布时间
   - 避免深夜发布，影响内容曝光

### 注意事项

1. **API调用限制**：
   - OpenAI API有调用频率限制，避免短时间大量生成
   - 建议设置合理的生成间隔，避免触发反爬机制

2. **账号安全**：
   - 自动发布功能请谨慎使用，避免账号被限制
   - 定期更换Cookie，避免失效

3. **内容合规**：
   - 遵守小红书社区规范，避免违规内容
   - 避免过于营销化的文案，保持内容自然

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

### 开发规范
- 遵循Go语言编码规范
- 提交信息使用Conventional Commits格式
- 新增功能请添加相应的测试

## 📄 许可证

MIT License

## 📞 支持

如有问题，请提交Issue或联系开发团队。

---

**免责声明**：本项目仅供学习和研究使用，请遵守相关平台的使用条款和法律法规。