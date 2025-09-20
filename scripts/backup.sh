#!/bin/bash

# 数据备份脚本

set -e

BACKUP_DIR="/var/backups/reading-app"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="reading_app_backup_${DATE}.tar.gz"

echo "🗄️ Creating backup of Reading Microservices data..."

# 创建备份目录
sudo mkdir -p "$BACKUP_DIR"

cd "$(dirname "$0")/../deployments/docker"

# 导出数据库
echo "📊 Backing up MySQL database..."
docker-compose exec -T mysql mysqldump -u root -ppassword reading_app > "${BACKUP_DIR}/database_${DATE}.sql"

# 备份Redis数据
echo "🔴 Backing up Redis data..."
docker-compose exec -T redis redis-cli --rdb /data/dump_${DATE}.rdb
docker cp reading-redis:/data/dump_${DATE}.rdb "${BACKUP_DIR}/"

# 备份下载文件
echo "📁 Backing up download files..."
if [ -d "/var/lib/reading-app/downloads" ]; then
    sudo cp -r /var/lib/reading-app/downloads "${BACKUP_DIR}/downloads_${DATE}"
fi

# 压缩备份
echo "🗜️ Compressing backup..."
cd "$BACKUP_DIR"
sudo tar -czf "$BACKUP_FILE" database_${DATE}.sql dump_${DATE}.rdb downloads_${DATE}/ 2>/dev/null || true

# 清理临时文件
sudo rm -f database_${DATE}.sql dump_${DATE}.rdb
sudo rm -rf downloads_${DATE}/

# 保留最近7天的备份
echo "🧹 Cleaning old backups (keeping last 7 days)..."
sudo find "$BACKUP_DIR" -name "reading_app_backup_*.tar.gz" -mtime +7 -delete

echo "✅ Backup completed: ${BACKUP_DIR}/${BACKUP_FILE}"
echo "📦 Backup size: $(sudo du -h ${BACKUP_DIR}/${BACKUP_FILE} | cut -f1)"