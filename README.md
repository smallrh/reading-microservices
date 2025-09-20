# Reading Microservices

基于Go的阅读应用微服务架构 - 完整版本

## 架构概述

```
┌─────────────────┐
│   Client Apps   │
│ (iOS/Android)   │
└─────────────────┘
         │
┌─────────────────┐
│   API Gateway   │
│      (Gin)      │
└─────────────────┘
         │
┌─────────────────────────────────────────────────────────┐
│                 Microservices Architecture              │
├─────────────────┬─────────────────┬─────────────────────┤
│   User Service  │ Content Service │  Reading Service    │
│     :8081       │     :8082       │     :8083           │
├─────────────────┼─────────────────┼─────────────────────┤
│ • 用户认证        │ • 小说管理       │ • 阅读记录          │
│ • 第三方登录      │ • 章节管理       │ • 书架管理          │
│ • 用户信息        │ • 分类标签       │ • 收藏评论          │
├─────────────────┼─────────────────┼─────────────────────┤
│  Payment Service│ Notification Svc│ Download Service    │
│     :8084       │     :8085       │     :8086           │
├─────────────────┼─────────────────┼─────────────────────┤
│ • VIP管理        │ • 推送通知       │ • 离线下载          │
│ • 积分阅读币      │ • 通知设置       │ • 任务管理          │
│ • 签到礼品        │ • 消息队列       │ • 文件存储          │
└─────────────────┴─────────────────┴─────────────────────┘
         │
┌─────────────────────────────────────────────────────────┐
│                Infrastructure                           │
├─────────────────┬─────────────────┬─────────────────────┤
│     MySQL       │      Redis      │      Consul         │
│     :3306       │      :6379      │      :8500          │
└─────────────────┴─────────────────┴─────────────────────┘
```

## 🎯 已实现的功能

### 1. 用户服务 (User Service) - 端口 8081
- ✅ 用户注册/登录
- ✅ JWT Token认证
- ✅ 第三方登录支持
- ✅ 用户信息管理
- ✅ 会话管理

### 2. 内容服务 (Content Service) - 端口 8082
- ✅ 小说管理 (CRUD)
- ✅ 章节管理 (CRUD)
- ✅ 分类/标签管理
- ✅ 内容搜索
- ✅ 统计数据

### 3. 阅读服务 (Reading Service) - 端口 8083
- ✅ 阅读记录管理
- ✅ 书架功能 (在读/收藏/下载)
- ✅ 评论系统
- ✅ 搜索历史
- ✅ 阅读统计

### 4. 支付会员服务 (Payment Service) - 端口 8084
- ✅ VIP会员管理
- ✅ 积分系统
- ✅ 阅读币管理
- ✅ 签到系统
- ✅ 礼品/兑换码系统
- ✅ 钱包管理

### 5. 通知服务 (Notification Service) - 端口 8085
- ✅ 通知管理
- ✅ 通知设置
- ✅ 推送Token管理
- ✅ 批量通知

### 6. 下载服务 (Download Service) - 端口 8086
- ✅ 下载任务管理
- ✅ 离线下载
- ✅ 文件格式支持 (TXT/EPUB/PDF)
- ✅ 下载进度跟踪

### 7. API网关 (API Gateway) - 端口 8080
- ✅ 统一入口
- ✅ 路由分发
- ✅ 认证授权
- ✅ 限流保护
- ✅ 健康检查

## 🛠️ 技术栈

- **Language**: Go 1.21+
- **Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL 8.0
- **Cache**: Redis 6.0
- **Service Discovery**: Consul
- **Config**: Viper + YAML
- **Logging**: Logrus
- **Container**: Docker + Docker Compose
- **Authentication**: JWT

## 📁 项目结构

```
reading-microservices/
├── api-gateway/          # API网关服务
├── user-service/         # 用户服务
├── content-service/      # 内容服务
├── reading-service/      # 阅读服务
├── payment-service/      # 支付会员服务
├── notification-service/ # 通知服务
├── download-service/     # 下载服务
├── shared/              # 共享组件
│   ├── config/          # 配置管理
│   ├── utils/           # 工具函数
│   ├── middleware/      # 中间件
│   └── go.mod
├── deployments/         # 部署配置
│   └── docker/          # Docker配置
├── scripts/             # 部署脚本
├── go.mod              # 根模块
└── README.md
```

## 🚀 快速开始

### 1. 环境要求
- **开发环境 (Windows)**: Docker Desktop, Go 1.21+
- **生产环境 (Linux)**: Docker, Docker Compose, 2GB+ RAM

### 2. 开发环境启动 (Windows)

#### 方式1: 使用批处理文件 (推荐)
```cmd
# 启动所有服务
start.bat

# 开发模式 (仅启动基础设施)
dev.bat

# 停止所有服务
stop.bat
```

#### 方式2: 使用Make命令
```cmd
make start   # 启动所有服务
make dev     # 开发模式
make stop    # 停止服务
```

#### 方式3: 手动开发模式
```cmd
# 启动基础设施
dev.bat

# 然后在不同终端启动各个服务
cd user-service && go run main.go
cd content-service && go run main.go
cd reading-service && go run main.go
cd payment-service && go run main.go
cd notification-service && go run main.go
cd download-service && go run main.go
cd api-gateway && go run main.go
```

### 3. 生产环境部署 (Linux)

```bash
# 克隆项目到Linux服务器
git clone <repository-url>
cd reading-microservices

# 一键部署到生产环境
chmod +x scripts/*.sh
sudo ./scripts/deploy.sh

# 检查服务状态
./scripts/status.sh

# 数据备份
./scripts/backup.sh
```

## 📊 服务端点

| 服务 | 端口 | 健康检查 | 描述 |
|------|------|----------|------|
| API Gateway | 8080 | `/health` | 统一入口 |
| User Service | 8081 | `/health` | 用户管理 |
| Content Service | 8082 | `/health` | 内容管理 |
| Reading Service | 8083 | `/health` | 阅读功能 |
| Payment Service | 8084 | `/health` | 支付会员 |
| Notification Service | 8085 | `/health` | 通知推送 |
| Download Service | 8086 | `/health` | 下载管理 |

### 基础设施
- **Consul UI**: http://localhost:8500
- **MySQL**: localhost:3306
- **Redis**: localhost:6379

## 🔌 API文档

### 认证流程
```bash
# 1. 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "platform": "web"
  }'

# 2. 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "platform": "web"
  }'

# 3. 使用Token访问受保护的API
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <token>"
```

### 主要API端点

#### 用户相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/user/profile` - 获取用户信息
- `PUT /api/v1/user/profile` - 更新用户信息

#### 内容相关
- `GET /api/v1/content/novels/search` - 搜索小说
- `GET /api/v1/content/novels/:id` - 获取小说详情
- `GET /api/v1/content/novels/:novel_id/chapters` - 获取章节列表

#### 阅读相关
- `POST /api/v1/reading/progress` - 更新阅读进度
- `GET /api/v1/reading/bookshelf` - 获取书架
- `POST /api/v1/reading/favorites/:novel_id` - 添加收藏

#### 支付相关
- `POST /api/v1/payment/checkin` - 每日签到
- `GET /api/v1/payment/wallet` - 获取钱包信息
- `POST /api/v1/payment/vip` - 购买VIP

## 📋 命令参考

### Windows开发环境
```cmd
# 批处理文件命令
start.bat          # 启动所有服务
dev.bat           # 开发模式 (仅基础设施)
stop.bat          # 停止所有服务

# Make命令 (需要安装make)
make help          # 显示帮助信息
make build         # 构建所有服务
make start         # 启动所有服务 (自动检测Windows)
make stop          # 停止所有服务
make dev           # 启动开发环境
make clean         # 清理Docker资源
make test          # 运行测试
```

### Linux生产环境
```bash
# 部署脚本
./scripts/deploy.sh   # 一键部署生产环境
./scripts/status.sh   # 检查服务状态
./scripts/backup.sh   # 数据备份
./scripts/start.sh    # 启动服务 (带健康检查)
./scripts/stop.sh     # 停止服务
./scripts/dev.sh      # 开发模式

# Docker Compose命令
docker-compose up -d --build    # 启动所有服务
docker-compose ps               # 查看服务状态
docker-compose logs -f          # 查看日志
docker-compose down             # 停止服务
```

## 🔧 配置说明

每个服务都有独立的配置文件 `config.yaml`，主要配置项：

```yaml
server:
  name: "service-name"
  host: "localhost"
  port: 8081

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "reading_app"

redis:
  host: "localhost"
  port: 6379
  db: 0

jwt:
  secret: "reading-app-secret-key"
  expires_in: 86400
```

## 🐛 故障排除

### 常见问题

1. **服务启动失败**
   ```bash
   # 检查日志
   docker-compose logs <service-name>

   # 检查端口占用
   netstat -tlnp | grep :8080
   ```

2. **数据库连接失败**
   ```bash
   # 检查MySQL服务
   docker-compose logs mysql

   # 手动连接测试
   mysql -h localhost -P 3306 -u root -p
   ```

3. **Redis连接失败**
   ```bash
   # 检查Redis服务
   docker-compose logs redis

   # 测试连接
   redis-cli ping
   ```

### 性能优化建议

1. **数据库优化**
   - 合理使用索引
   - 优化查询语句
   - 配置连接池

2. **缓存策略**
   - 热点数据缓存
   - 查询结果缓存
   - 分布式锁

3. **服务优化**
   - 异步处理
   - 批量操作
   - 限流保护

## 🔒 安全说明

- JWT Token过期时间: 24小时
- 密码采用bcrypt加密
- API限流保护
- CORS跨域配置
- 输入参数验证

## 📝 开发规范

1. **代码规范**
   - 遵循Go官方规范
   - 使用gofmt格式化
   - 添加必要注释

2. **Git规范**
   - feat: 新功能
   - fix: 修复bug
   - docs: 文档更新
   - refactor: 重构

3. **测试规范**
   - 单元测试覆盖率 > 70%
   - 集成测试
   - API测试

## 🤝 贡献指南

1. Fork项目
2. 创建特性分支
3. 提交代码
4. 发起Pull Request

## 📄 许可证

MIT License

---

**项目状态**: ✅ 完成
**版本**: v1.0.0
**最后更新**: 2025-01-21