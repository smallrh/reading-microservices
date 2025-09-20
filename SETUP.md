# 🛠️ 开发环境安装指南

这个指南将帮助您在Windows环境下安装和配置开发环境。

## 📋 系统要求

- Windows 10/11 (64位)
- 至少 4GB RAM
- 至少 10GB 可用磁盘空间

## 🔧 安装步骤

### 1. 检查当前环境

运行环境检查脚本：
```cmd
check-requirements.bat
```

这个脚本会检查您系统上是否已安装必要的工具。

### 2. 安装 Docker Desktop

如果您还没有安装Docker：

1. **下载Docker Desktop**
   - 访问：https://www.docker.com/products/docker-desktop/
   - 点击 "Download Docker Desktop for Windows"

2. **安装Docker Desktop**
   - 运行下载的 `Docker Desktop Installer.exe`
   - 勾选 "Install required Windows components for WSL 2"
   - 完成安装后重启计算机

3. **启动Docker Desktop**
   - 从开始菜单启动 Docker Desktop
   - 等待Docker启动完成（系统托盘图标变为绿色）
   - 可能需要登录或创建Docker Hub账户

4. **验证安装**
   ```cmd
   docker --version
   docker run hello-world
   ```

### 3. 安装 Go (可选，用于开发)

如果您计划修改源码：

1. **下载Go**
   - 访问：https://golang.org/dl/
   - 下载最新的Go 1.21+版本 (例：go1.21.x.windows-amd64.msi)

2. **安装Go**
   - 运行下载的MSI文件
   - 使用默认安装路径：`C:\Program Files\Go`

3. **验证安装**
   ```cmd
   go version
   ```

### 4. 安装 Make (可选)

Make工具是可选的，您可以直接使用批处理文件：

#### 选项1: 使用 Chocolatey (推荐)
```cmd
# 首先安装Chocolatey (如果未安装)
# 以管理员身份运行PowerShell，然后执行：
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# 安装Make
choco install make
```

#### 选项2: 不安装Make，直接使用批处理文件
```cmd
start.bat   # 代替 make start
dev.bat     # 代替 make dev
stop.bat    # 代替 make stop
```

## 🚀 启动服务

完成安装后，您可以使用以下方式启动服务：

### 方式1: 使用批处理文件 (推荐)
```cmd
# 启动所有服务
start.bat

# 仅启动基础设施 (开发模式)
dev.bat

# 停止所有服务
stop.bat
```

### 方式2: 使用Make命令 (如果已安装Make)
```cmd
make start
make dev
make stop
```

### 方式3: 直接使用Docker命令
```cmd
cd deployments\docker
docker compose up -d --build
```

## 🔍 验证安装

1. **运行环境检查**
   ```cmd
   check-requirements.bat
   ```

2. **启动服务**
   ```cmd
   start.bat
   ```

3. **检查服务状态**
   - API Gateway: http://localhost:8080/health
   - Consul UI: http://localhost:8500

## ❓ 常见问题

### Docker相关问题

**Q: Docker Desktop启动失败**
A:
- 确保Windows功能中启用了"适用于Linux的Windows子系统"和"虚拟机平台"
- 重启计算机后再试
- 检查BIOS中是否启用了虚拟化

**Q: "docker: command not found"**
A:
- 确保Docker Desktop正在运行
- 重启命令行窗口
- 检查环境变量PATH中是否包含Docker路径

**Q: 容器启动失败，端口被占用**
A:
```cmd
# 检查端口占用
netstat -ano | findstr :8080

# 停止占用端口的进程 (替换PID)
taskkill /PID <进程ID> /F
```

### Go相关问题

**Q: "go: command not found"**
A:
- 确保Go安装成功
- 重启命令行窗口
- 检查环境变量GOPATH和GOROOT

### 性能问题

**Q: 服务启动很慢**
A:
- 确保Docker Desktop分配足够内存 (建议4GB+)
- 关闭不必要的后台程序
- 使用SSD硬盘可显著提升性能

## 📞 获取帮助

如果遇到问题：
1. 运行 `check-requirements.bat` 检查环境
2. 查看Docker Desktop的日志和错误信息
3. 检查防火墙和杀毒软件设置
4. 确保以管理员权限运行相关命令

## 🎯 下一步

环境配置完成后，您可以：
- 查看 `README.md` 了解API使用方法
- 修改各服务的配置文件
- 开始开发和测试您的应用