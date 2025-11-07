# Docker配置文件详细解释

## Dockerfile 逐行解释

```dockerfile
# 多阶段构建 Dockerfile for Go应用

# 第一阶段：构建阶段
FROM golang:1.24.5-alpine AS builder
```
**作用**: 使用轻量级的Alpine Linux作为基础镜像来构建Go应用

```dockerfile
# 设置工作目录
WORKDIR /app
```
**作用**: 在容器内创建并切换到/app目录

```dockerfile
# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./
```
**作用**: 只复制依赖管理文件，利用Docker缓存

```dockerfile
# 下载依赖
RUN go mod download
```
**作用**: 下载Go模块依赖，如果go.mod和go.sum没有变化，可以复用缓存

```dockerfile
# 复制源代码
COPY . .
```
**作用**: 将项目所有源代码复制到容器中

```dockerfile
# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
```
**作用**: 编译Go应用，生成可执行文件main

```dockerfile
# 第二阶段：运行阶段
FROM alpine:latest
```
**作用**: 使用更小的Alpine镜像作为运行环境

```dockerfile
# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates
```
**作用**: 安装CA证书，用于HTTPS连接

```dockerfile
# 创建非root用户
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
```
**作用**: 创建非root用户，提高安全性

```dockerfile
# 切换到非root用户
USER appuser
```
**作用**: 以非特权用户运行应用

```dockerfile
WORKDIR /app
```
**作用**: 设置运行时的工作目录

```dockerfile
# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env.example
```
**作用**: 从构建阶段复制编译好的应用和配置模板

```dockerfile
# 暴露端口
EXPOSE 8080
```
**作用**: 声明容器监听的端口

```dockerfile
# 运行应用
CMD ["./main"]
```
**作用**: 容器启动时执行的命令

## docker-compose.yml 逐行解释

```yaml
version: '3.8'
```
**作用**: 指定Docker Compose文件格式版本

```yaml
services:
  # Go应用服务
  app:
    build:
      context: .
      dockerfile: Dockerfile
```
**作用**: 定义应用服务，从当前目录构建

```yaml
    ports:
      - "8080:8080"
```
**作用**: 将容器内的8080端口映射到宿主机的8080端口
```

```yaml
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=${DB_PASSWORD}
```
**作用**: 设置环境变量，${DB_PASSWORD}会从.env文件读取

```yaml
    depends_on:
      - mysql
```
**作用**: 确保MySQL服务先启动

```yaml
  # MySQL数据库服务
  mysql:
    image: mysql:8.0
```
**作用**: 使用MySQL 8.0官方镜像

```yaml
    environment:
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
```
**作用**: 设置MySQL的root密码

## 协作者需要修改的内容

### 1. 环境配置文件 (.env)
```bash
# 必须修改的配置
DB_PASSWORD=您的真实数据库密码
MYSQL_ROOT_PASSWORD=您的MySQL root密码
```

## 协作开发具体操作

### 协作者需要：
1. **复制配置模板**:
   ```bash
   cp .env.example .env
   ```

2. **编辑.env文件**:
   ```bash
   # 修改以下配置项
DB_PASSWORD=您的真实密码
MYSQL_ROOT_PASSWORD=您的真实密码
```

### 2. 数据库连接配置
- 如果使用本地数据库：`DB_HOST=127.0.0.1`
- 如果使用容器数据库：`DB_HOST=mysql`

## 后续开发修改指南

### 1. 添加新功能
- 在相应目录添加新的Go文件
- 更新路由配置
- 重新构建镜像

## 部署到服务器的修改

### 1. 生产环境配置
```yaml
# 在docker-compose.yml中修改
environment:
  - DB_HOST=生产环境数据库地址
```

### 2. 环境变量调整
- 修改`.env`文件中的数据库连接信息

### 3. 安全配置
- 使用更强的JWT密钥
- 配置HTTPS证书
- 设置防火墙规则

### 4. 性能优化
- 调整容器资源限制
- 配置健康检查
- 设置日志收集

## 服务器部署关键修改点

### 1. 数据库连接
```yaml
# 生产环境
environment:
  - DB_HOST=生产数据库IP或域名
```

## 总结

**协作者只需**：
1. 复制`.env.example`为`.env`
2. 在`.env`中设置生产环境密码
```

现在您和协作者都有了清晰的Docker使用指南和协作开发流程。