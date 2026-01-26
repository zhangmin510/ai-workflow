# 火山引擎豆包模型配置说明

## 1. 获取API Key
- 登录火山引擎控制台 (https://console.volcengine.com/)
- 进入"火山方舟"产品
- 创建API密钥
- 复制API Key

## 2. 配置文件修改
在`config/config.yaml`中修改以下内容：
```yaml
doubao:
  api_key: "your_volcengine_api_key_here"  # 替换为你的火山引擎API Key
  model: "doubao-1.8"  # 指定使用豆包1.8模型
  temperature: 0.8      # 温度参数，控制生成结果的随机性
  max_tokens: 1000      # 最大生成令牌数
```

## 3. 支持的模型版本
- `doubao-1.8`: 豆包大模型1.8版本（推荐）
- `doubao-1.5`: 豆包大模型1.5版本
- `doubao-lite-4k`: 豆包轻量版（4k上下文窗口）

## 4. 注意事项
- 火山引擎豆包API使用Bearer Token认证，直接使用API Key即可
- 无需获取access_token，API调用时会自动处理
- 请确保你的API Key有足够的权限调用豆包模型API
- 具体API文档请参考火山引擎官方文档: https://www.volcengine.com/docs/82379/2123228
