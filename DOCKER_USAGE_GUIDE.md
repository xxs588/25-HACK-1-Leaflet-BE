# Docker 使用指南

## 前提条件

确保您的系统已安装：
- Docker Desktop (Windows/Mac)
- 或 Docker Engine (Linux)

## 基本使用步骤

### 1. 准备环境配置
```bash
# 复制配置模板
cp .env.example .env

# 编辑.env文件，设置您的数据库密码
# 使用文本编辑器打开 .env 文件，修改以下内容：
DB_PASSWORD=您的真实密码
MYSQL_ROOT_PASSWORD=您的真实密码
```

### 2. 启动开发环境
```bash
# 构建并启动所有服务
docker-compose up -d

# 或者构建并启动（显示日志）
docker-compose up
```

## 详细操作步骤

### 1. 环境设置
```bash
# 1. 复制配置模板
cp .env.example .env

# 2. 编辑.env文件，设置您的数据库密码

# 3. 启动服务
docker-compose up -d
```

### 2. 查看服务状态
```bash
# 查看运行中的容器
docker-compose ps

# 或者查看所有Docker容器
docker ps
```

### 3. 访问应用
- **应用地址**: http://localhost:8080
- **数据库地址**: localhost:3306

### 4. 停止服务
```bash
# 停止并移除容器
docker-compose down

# 或者只停止不删除
docker-compose stop
```

## 常用命令

### 构建相关
```bash
# 仅构建镜像（不启动）
docker-compose build

# 重新构建并启动
docker-compose up -d --build
```

### 5. 查看日志
```bash
# 查看应用日志
docker-compose logs app

# 查看数据库日志
docker-compose logs mysql

# 实时查看日志
docker-compose logs -f app
```

### 6. 清理资源
```bash
# 停止并删除所有容器、网络、卷
docker-compose down -v
```

## 开发工作流

### 1. 本地开发（推荐）
```bash
# 直接运行Go应用
go run main.go
```

## 生产环境部署

### 1. 构建生产镜像
```bash
# 构建镜像
docker build -t your-app-name .

# 运行容器
docker run -p 8080:8080 your-app-name
```

## 故障排除

### 1. 端口冲突
如果8080端口被占用，可以修改端口映射：
```yaml
# 在docker-compose.yml中修改
ports:
  - "8081:8080"  # 将外部端口改为8081
```

### 2. 数据库连接问题
确保MySQL服务在容器中正常运行：
```bash
# 检查MySQL容器状态
docker-compose ps mysql

# 进入MySQL容器
docker-compose exec mysql mysql -u root -p
```

## 验证部署

### 1. 测试API
```bash
# 测试注册接口
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'
```

## 协作开发注意事项

### 1. 配置管理
- **提交到版本控制**:
  - Dockerfile
  - docker-compose.yml
  - .env.example

# 测试登录接口
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'
```

### 2. 版本同步
- 定期拉取最新代码
- 更新依赖包
- 运行数据库迁移