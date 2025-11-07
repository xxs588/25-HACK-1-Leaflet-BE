# Docker配置简化说明

## 关于三个容器的解释

在当前的`docker-compose.yml`中确实定义了三个服务，但这是为了提供灵活性：

### 1. **app** - 您的Go应用（必须）
- 这是您的主要应用服务

### 2. **mysql** - 主数据库（必须）
- 这是应用依赖的数据库

### 3. **mysql-dev** - 开发数据库（可选）

## 实际使用建议

### 1. **最小化配置**（推荐）
您只需要运行前两个容器：
- **app** - 您的Go应用
- **mysql** - 数据库服务

**mysql-dev**是**可选的**，用于：
- 为不同的协作者提供独立的开发数据库
- 避免数据库冲突

## 简化后的docker-compose.yml

```yaml
version: '3.8'

services:
  # Go应用服务
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
  - DB_USER=root
  - DB_PASSWORD=${DB_PASSWORD}
  - DB_NAME=hackweek_db
  - DB_CHARSET=utf8mb4
  - DB_PARSE_TIME=True
  - DB_LOC=Local
  - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - mysql
    restart: unless-stopped

  # MySQL数据库服务
  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql_data:
```

## 实际运行方式

### 1. **只运行必要服务**：
```bash
# 只运行app和mysql
docker-compose up app mysql
```

## 为什么这样设计

### 1. **开发环境灵活性**
- 协作者可以选择使用本地数据库或容器数据库

### 2. **生产环境简化**
在生产环境中，您可能只需要：
- **app**容器（连接到外部数据库）

## 总结

**您只需要运行**：
- `app` - 您的Go应用
- `mysql` - 数据库服务

**mysql-dev**是**完全可选的**，可以根据团队需求决定是否使用。