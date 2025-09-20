# Dify免费大模型配置指南

## 🎯 推荐配置顺序

### 方案一：本地Ollama模型（推荐，完全免费）

#### 1. 安装Ollama
```bash
# Windows: 下载安装包
# https://ollama.com/download/windows

# 或使用包管理器
winget install Ollama.Ollama
```

#### 2. 下载推荐模型
```bash
# 轻量级模型（适合CPU运行）
ollama pull qwen2.5:7b
ollama pull llama3.1:8b

# 如果有GPU，可以用更大模型
ollama pull qwen2.5:14b
ollama pull llama3.1:70b
```

#### 3. 在Dify中配置Ollama
1. 访问 http://localhost:8093
2. 进入 **设置** → **模型供应商**
3. 添加 **OpenAI兼容** 模型
4. 配置参数：
   ```
   供应商名称: Ollama
   API Base URL: http://host.docker.internal:11434/v1
   API Key: ollama
   模型名称: qwen2.5:7b
   ```

---

### 方案二：智谱AI GLM-4-Flash

#### 1. 获取API Key
1. 访问：https://open.bigmodel.cn/
2. 注册账号并完成实名认证
3. 在控制台获取API Key

#### 2. 在Dify中配置
1. 添加 **自定义模型**
2. 配置参数：
   ```
   供应商名称: ZhipuAI
   API Base URL: https://open.bigmodel.cn/api/paas/v4
   API Key: 你的智谱AI API Key
   模型名称: glm-4-flash
   模型类型: 文本生成
   ```

---

### 方案三：阿里云通义千问

#### 1. 获取API Key
1. 访问：https://dashscope.aliyun.com/
2. 注册阿里云账号
3. 开通DashScope服务
4. 创建API-KEY

#### 2. 在Dify中配置
1. 添加 **OpenAI兼容** 模型
2. 配置参数：
   ```
   供应商名称: Qwen
   API Base URL: https://dashscope.aliyuncs.com/compatible-mode/v1
   API Key: 你的DashScope API-KEY
   模型名称: qwen-turbo
   ```

---

### 方案四：讯飞星火（Spark）

#### 1. 获取API凭证
1. 访问：https://xinghuo.xfyun.cn/
2. 注册讯飞开发者账号
3. 创建应用获取APPID、APISecret、APIKey

#### 2. 在Dify中配置
1. 添加 **讯飞星火** 官方供应商
2. 填入相应的凭证信息

---

## 🔧 配置步骤详解

### 步骤1：启动Dify
```bash
cd reading-microservices
start.bat
```

### 步骤2：访问Dify控制台
打开浏览器访问：http://localhost:8093

### 步骤3：模型供应商配置
1. 登录后点击右上角用户头像
2. 选择 **设置**
3. 点击左侧 **模型供应商**
4. 点击 **+ 添加供应商**

### 步骤4：测试模型
1. 创建新应用或编辑现有应用
2. 在模型配置中选择刚添加的模型
3. 发送测试消息验证模型是否正常工作

---

## 📊 免费额度对比

| 模型供应商 | 免费额度 | 模型质量 | 响应速度 | 推荐指数 |
|-----------|---------|---------|---------|---------|
| Ollama本地 | 无限制 | 中等 | 快（本地） | ⭐⭐⭐⭐⭐ |
| 智谱GLM-4-Flash | 1000万tokens/日 | 高 | 快 | ⭐⭐⭐⭐⭐ |
| 通义千问 | 100万tokens/月 | 高 | 中等 | ⭐⭐⭐⭐ |
| 讯飞星火 | 200万tokens/月 | 中等 | 中等 | ⭐⭐⭐ |

---

## 🚨 注意事项

### Ollama配置注意
- **Docker网络问题**：使用 `host.docker.internal:11434` 而不是 `localhost:11434`
- **内存要求**：7B模型至少需要8GB RAM，14B模型需要16GB+
- **GPU加速**：如果有NVIDIA GPU，Ollama会自动使用CUDA加速

### API配置注意
- **防火墙**：确保Dify容器能访问外部API
- **代理设置**：如果使用代理，需要在Docker环境中配置
- **API限制**：注意各供应商的QPS限制

### 多模型策略
建议同时配置多个模型：
1. **主力模型**：智谱GLM-4-Flash（日常使用）
2. **备用模型**：Ollama本地模型（无网络时使用）
3. **专用模型**：针对特定任务选择最适合的模型

---

## 🛠️ 故障排除

### 连接问题
```bash
# 测试Ollama连接
curl http://localhost:11434/api/tags

# 测试外部API
curl -H "Authorization: Bearer YOUR_API_KEY" https://open.bigmodel.cn/api/paas/v4/models
```

### Docker网络问题
```bash
# 检查Docker网络
docker network ls
docker network inspect reading-microservices_default
```

### 日志查看
```bash
# 查看Dify API日志
docker logs reading-dify-api -f

# 查看Ollama日志
ollama logs
```