# Dify AI 平台部署指南

## 🚀 快速开始

### 1. 启动服务
```bash
cd reading-microservices
start.bat
```

### 2. 访问Dify平台
- **Dify Web界面**: http://localhost:8093
- **Dify API**: http://localhost:8091

## 📋 服务组件

Dify平台包含以下组件：

| 服务 | 端口 | 说明 |
|------|------|------|
| dify-web | 8093 | Web管理界面 |
| dify-api | 8091 | API服务 |
| dify-worker | - | 后台任务处理 |
| dify-db | - | PostgreSQL数据库 |
| dify-redis | - | Redis缓存 |
| dify-weaviate | - | 向量数据库 |

## 🔧 初始配置

### 1. 首次访问
1. 打开浏览器访问 http://localhost:8093
2. 创建管理员账户
3. 登录Dify管理控制台

### 2. 配置DeepSeek模型（可选）
1. 在Dify中进入"设置" → "模型供应商"
2. 添加"自定义模型"或"OpenAI兼容"
3. 配置DeepSeek API：
   - **API Base URL**: `https://api.deepseek.com/v1`
   - **API Key**: 你的DeepSeek API密钥
   - **模型名称**: `deepseek-chat`, `deepseek-coder`

### 3. 创建应用
1. 点击"创建应用"
2. 选择应用类型（聊天助手、Agent等）
3. 配置你的AI应用

## 🔗 集成到微服务

### 获取API密钥
1. 在Dify应用详情页面
2. 复制"API密钥"
3. 设置环境变量：`DIFY_API_KEY=app-xxxxxxxxxx`

### API调用示例

#### 聊天消息API
```bash
curl -X POST 'http://localhost:8091/v1/chat-messages' \
--header 'Authorization: Bearer app-xxxxxxxxxx' \
--header 'Content-Type: application/json' \
--data-raw '{
    "inputs": {},
    "query": "Hello, how are you?",
    "response_mode": "blocking",
    "conversation_id": "",
    "user": "reading-user"
}'
```

#### 在Go代码中集成
```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type DifyRequest struct {
    Inputs         map[string]interface{} `json:"inputs"`
    Query          string                 `json:"query"`
    ResponseMode   string                 `json:"response_mode"`
    ConversationID string                 `json:"conversation_id"`
    User           string                 `json:"user"`
}

type DifyResponse struct {
    Answer         string `json:"answer"`
    ConversationID string `json:"conversation_id"`
}

func callDify(query string, apiKey string) (string, error) {
    req := DifyRequest{
        Inputs:       map[string]interface{}{},
        Query:        query,
        ResponseMode: "blocking",
        User:         "reading-user",
    }

    jsonData, _ := json.Marshal(req)
    httpReq, _ := http.NewRequest("POST", "http://localhost:8091/v1/chat-messages",
        bytes.NewBuffer(jsonData))

    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer " + apiKey)

    client := &http.Client{}
    resp, err := client.Do(httpReq)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result DifyResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return result.Answer, nil
}
```

## 🛠️ 管理和维护

### 查看日志
```bash
# 查看所有Dify服务日志
docker compose logs -f dify-api dify-worker dify-web

# 查看特定服务日志
docker logs reading-dify-api
```

### 重启服务
```bash
# 重启所有Dify服务
docker compose restart dify-api dify-worker dify-web dify-db dify-redis dify-weaviate

# 重启特定服务
docker compose restart dify-api
```

### 数据备份
```bash
# 备份PostgreSQL数据库
docker exec reading-dify-db pg_dump -U postgres dify > dify_backup.sql

# 备份向量数据库
docker exec reading-dify-weaviate tar czf - /var/lib/weaviate > weaviate_backup.tar.gz
```

## 🔒 安全配置

### 生产环境建议
1. **修改默认密码**：
   ```yaml
   # 在docker-compose.yml中修改
   POSTGRES_PASSWORD: your-secure-password
   REDIS_PASSWORD: your-redis-password
   SECRET_KEY: your-secret-key
   ```

2. **配置HTTPS**：使用nginx反向代理配置SSL

3. **网络隔离**：限制对内部端口的访问

### 环境变量
```bash
# 复制环境变量模板
cp .env.example .env

# 编辑配置
vim .env
```

## 📚 更多资源

- [Dify官方文档](https://docs.dify.ai/)
- [Dify GitHub](https://github.com/langgenius/dify)
- [DeepSeek API文档](https://platform.deepseek.com/api-docs/)

## ❓ 常见问题

### Q: Dify启动失败怎么办？
A: 检查Docker资源分配，确保至少有4GB内存可用

### Q: 如何更换模型？
A: 在Dify控制台的"模型供应商"中添加和切换模型

### Q: 数据存储在哪里？
A: 数据存储在Docker volumes中，可以通过 `docker volume ls` 查看

### Q: 如何升级Dify？
A: 修改docker-compose.yml中的镜像版本号，然后执行 `docker compose pull && docker compose up -d`