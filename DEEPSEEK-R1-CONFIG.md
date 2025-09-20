# DeepSeek R1 配置指南

## ✅ 配置完成状态

- ✅ Ollama服务已配置到Docker Compose
- ✅ DeepSeek R1 (1.5B) 模型已下载成功
- ✅ 模型测试正常运行，支持中文对话

## 📋 当前配置

### Docker Compose服务配置

```yaml
ollama:
  image: ollama/ollama:latest
  container_name: reading-ollama
  ports:
    - "11434:11434"
  volumes:
    - ollama_data:/root/.ollama
  environment:
    - OLLAMA_HOST=0.0.0.0
    - OLLAMA_ORIGINS=*
    - OLLAMA_NUM_PARALLEL=2
    - OLLAMA_MAX_LOADED_MODELS=1
    - OLLAMA_KEEP_ALIVE=5m
  networks:
    - reading-network
  restart: unless-stopped
  healthcheck:
    test: ["CMD", "curl", "-f", "http://localhost:11434/api/version"]
    interval: 30s
    timeout: 15s
    retries: 5
    start_period: 30s
```

### 模型信息

```
名称: deepseek-r1:1.5b
ID: e0979632db5a
大小: 1.1 GB
类型: 推理优化模型
特点: DeepSeek R1蒸馏版本，支持思维链推理
```

## 🚀 使用指南

### 1. 启动服务

```bash
cd deployments/docker
docker compose up -d ollama
```

### 2. 验证服务状态

```bash
# 检查容器状态
docker compose ps | grep ollama

# 检查API是否可用
curl http://localhost:11434/api/version

# 查看已下载的模型
docker exec reading-ollama ollama list
```

### 3. 直接使用模型

```bash
# 交互式使用
docker exec -it reading-ollama ollama run deepseek-r1:1.5b

# 单次查询
docker exec reading-ollama ollama run deepseek-r1:1.5b "你好，介绍一下自己"
```

### 4. 在Dify中配置使用

#### 步骤1：添加模型供应商
1. 访问 Dify Web界面：http://localhost:8093
2. 登录管理控制台
3. 进入 **设置** → **模型供应商**
4. 点击 **+ 添加供应商**
5. 选择 **OpenAI兼容**

#### 步骤2：配置DeepSeek R1
```
供应商名称: DeepSeek-R1-Local
API Base URL: http://ollama:11434/v1
API Key: deepseek-r1
```

#### 步骤3：添加模型
```
模型名称: deepseek-r1:1.5b
模型类型: 文本生成
上下文长度: 4096
最大输出token: 2048
支持视觉: 否
```

## 🎯 模型特性

### DeepSeek R1 特点
- **推理能力强**：基于DeepSeek-R1架构，支持复杂推理
- **中文友好**：对中文理解和生成能力优秀
- **思维链推理**：支持step-by-step思维过程
- **轻量化**：1.5B参数版本，适合本地部署
- **高效率**：经过蒸馏优化，推理速度快

### 适用场景
- 逻辑推理和数学问题
- 代码分析和生成
- 中文自然语言处理
- 问答系统
- 文档理解和总结

## 🔧 高级配置

### GPU加速（可选）

如果你有NVIDIA GPU，可以启用GPU加速：

1. **取消注释GPU配置**：
   ```yaml
   deploy:
     resources:
       reservations:
         devices:
           - driver: nvidia
             count: all
             capabilities: [gpu]
   ```

2. **重启服务**：
   ```bash
   docker compose down ollama
   docker compose up -d ollama
   ```

### 性能调优

```bash
# 设置模型保持时间（减少重复加载）
docker exec reading-ollama ollama run deepseek-r1:1.5b --keep-alive 10m

# 预加载模型到内存
docker exec reading-ollama curl -X POST http://localhost:11434/api/generate \
  -d '{"model": "deepseek-r1:1.5b", "prompt": "", "keep_alive": -1}'
```

## 📡 API使用示例

### 直接API调用

```bash
# 基础对话
curl -X POST http://localhost:11434/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-r1:1.5b",
    "messages": [
      {"role": "user", "content": "解释一下量子计算的基本原理"}
    ],
    "temperature": 0.7,
    "max_tokens": 1000
  }'
```

### 通过Dify API

```bash
curl -X POST 'http://localhost:8091/v1/chat-messages' \
  --header 'Authorization: Bearer app-xxxxxxxxxx' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "inputs": {},
    "query": "帮我分析一下这个算法的复杂度",
    "response_mode": "blocking",
    "conversation_id": "",
    "user": "reading-user"
  }'
```

## 🛠️ 管理维护

### 模型管理

```bash
# 查看模型信息
docker exec reading-ollama ollama show deepseek-r1:1.5b

# 更新模型
docker exec reading-ollama ollama pull deepseek-r1:1.5b

# 删除模型（如需要）
docker exec reading-ollama ollama rm deepseek-r1:1.5b
```

### 监控和日志

```bash
# 查看容器状态
docker compose ps ollama

# 查看服务日志
docker logs reading-ollama -f

# 监控资源使用
docker stats reading-ollama
```

### 数据备份

```bash
# 备份模型数据
docker run --rm -v reading-microservices_ollama_data:/data \
  -v $(pwd):/backup ubuntu \
  tar czf /backup/deepseek-r1-backup.tar.gz /data
```

## 🚨 故障排除

### 常见问题

1. **模型响应慢**
   - 检查内存使用情况
   - 考虑启用GPU加速
   - 调整OLLAMA_KEEP_ALIVE参数

2. **内存不足**
   - 确保至少有3GB可用内存
   - 关闭其他占用内存的应用
   - 使用swap分区作为补充

3. **网络连接问题**
   - 检查端口11434是否被占用
   - 验证Docker网络配置
   - 测试API连通性

### 性能优化建议

```yaml
# 针对DeepSeek R1的优化环境变量
environment:
  - OLLAMA_HOST=0.0.0.0
  - OLLAMA_ORIGINS=*
  - OLLAMA_NUM_PARALLEL=1      # 1.5B模型建议单线程
  - OLLAMA_MAX_LOADED_MODELS=1 # 只保持一个模型在内存
  - OLLAMA_KEEP_ALIVE=10m      # 延长模型保持时间
```

## 📚 相关资源

- **DeepSeek官网**: https://www.deepseek.com/
- **DeepSeek R1论文**: https://arxiv.org/abs/2501.12948
- **Ollama文档**: https://ollama.com/docs
- **模型详情**: https://ollama.com/deepseek/deepseek-r1

## 🎉 测试验证

模型已通过以下测试：
- ✅ 基础对话功能
- ✅ 中文理解和生成
- ✅ 推理思维链输出
- ✅ API接口调用
- ✅ Docker容器稳定性

DeepSeek R1模型现已完全配置就绪，可以在你的阅读微服务系统中使用！