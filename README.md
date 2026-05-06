
# 🍃 Leaflet -一页（叶）倾心后端

一个基于 Go + Gin + GORM 的心情记录与社交互动平台后端服务。

## 📋 项目简介

Leaflet 是一个帮助用户记录日常心情、培养记录习惯的应用。用户可以：

-   📝 创建带标签的心情记录
-   🌳 查看连续记录天数和成长等级
-   💬 分享困扰并获得他人的鼓励
-   🎨 自定义个人头像
-   📊 查看同状态的用户数量

## 🏗️ 项目结构

```
25-HACK-1-Leaflet-BE/
├── main.go                 # 应用入口
├── config/                 # 配置相关
│   └── db.go              # 数据库连接配置
├── controller/             # 控制器层
│   ├── user.go            # 用户认证控制器
│   ├── status.go          # 心情状态控制器
│   ├── communication.go   # 情绪互动控制器
│   ├── encouragements.go  # 鼓励话语控制器
│   └── myself.go          # 个人信息控制器
├── middlewares/            # 中间件
│   └── auth.go            # JWT 认证中间件
├── model/                  # 数据模型
│   └── user.go            # 数据库模型定义
├── routes/                 # 路由定义
│   └── user.go            # 路由配置
├── consts/                 # 常量定义
│   ├── image.go           # 预设头像列表
│   └── log.go             # 日志配置
├── docker-compose.yml      # Docker Compose 配置
├── Dockerfile              # Docker 构建文件
├── .env.example            # 环境变量模板
└── README.md              # 项目文档
```

## 🚀 快速开始

### 环境要求

-   **Docker & Docker Compose**（推荐）
-   或者：Go 1.22+ & MariaDB/MySQL 5.7+

### 方式一：Docker 部署（推荐）

#### 1. 克隆项目

```bash
git clone https://github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE.git
cd 25-HACK-1-Leaflet-BE
```

#### 2. 配置环境变量

```bash
cp .env.example .env
# 根据需要修改 .env 中的配置
```

#### 3. 启动服务

```bash
# 构建并启动所有服务
docker-compose up --build

# 后台运行
docker-compose up -d --build

# 查看日志
docker-compose logs -f app
```

#### 4. 访问服务

-   API 服务：http://localhost:8080
-   MariaDB：localhost:3339

#### 5. 停止服务

```bash
docker-compose down

# 同时删除数据卷（清空数据库）
docker-compose down -v
```

### 方式二：本地开发

#### 1. 安装依赖

```bash
go mod download
```

#### 2. 配置数据库

```bash
# 启动 MySQL/MariaDB
# 创建数据库
mysql -u root -p
CREATE DATABASE leaflet_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 3. 配置环境变量

```bash
cp .env.example .env
# 修改 .env 中的数据库连接信息
# DB_HOST=127.0.0.1
# DB_PORT=3306
# DB_USER=root
# DB_PASSWORD=your_password
# DB_NAME=leaflet_db
```

#### 4. 运行应用

```bash
go run main.go
```

## 🐛 常见问题

### Docker 相关

#### 问题：镜像拉取失败

**解决方案**：配置 Docker 镜像加速

创建或编辑 `/etc/docker/daemon.json`（Linux）或 Docker Desktop 设置（Windows/Mac）：

```json
{
    "registry-mirrors": [
        "https://docker.mirrors.ustc.edu.cn",
        "https://hub-mirror.c.163.com",
        "https://mirror.baidubce.com"
    ]
}
```

重启 Docker：

```bash
# Linux
sudo systemctl restart docker

# Windows/Mac：重启 Docker Desktop
```

#### 问题：端口被占用

```bash
# 修改 docker-compose.yml 中的端口映射
ports:
  - "8081:8080"  # 将 8080 改为其他端口
```

#### 问题：数据库连接失败

```bash
# 等待数据库完全启动（约 10-20 秒）
docker-compose logs mariadb

# 检查应用日志
docker-compose logs app
```

## 📡 API 接口文档

### 公开接口（无需认证）

#### 用户认证

-   `POST /register` - 用户注册
-   `POST /login` - 用户登录
-   `GET /encouragements` - 获取鼓励话语

### 需要认证的接口（需要 JWT Token）

> 在请求头中添加：`Authorization: Bearer <token>`

#### 心情状态管理

-   `POST /status` - 创建心情记录
-   `GET /status/mine` - 获取个人所有记录
-   `GET /status/by_tag/:tag_id` - 获取指定标签的用户数
-   `GET /status/level` - 获取用户心情树等级
-   `PUT /status/:id` - 更新心情记录
-   `DELETE /status/:id` - 删除心情记录

#### 情绪互动

-   `POST /mind` - 发布困扰
-   `GET /mind` - 获取所有困扰
-   `PUT /mind/:id` - 修改困扰
-   `DELETE /mind/:id` - 删除困扰
-   `POST /solve/:id` - 回复解决方案
-   `GET /solve` - 获取所有解决方案

#### 个人信息

-   `GET /image` - 获取头像列表和当前头像
-   `PUT /image` - 更新头像
-   `PUT /user/name` - 更新用户名

### 标签 ID 对照表

| Tag ID | 标签名称   | 树叶颜色 |
| ------ | ---------- | -------- |
| 1      | 困倦的早八 | 绿       |
| 2      | 自习室刷题 | 蓝       |
| 3      | 图书馆阅读 | 红       |
| 4      | 食堂干饭   | 紫       |
| 5      | 备考冲刺   | 橙       |
| 6      | 社团活动   | 粉       |
| 7      | 情绪波动时 | 灰       |
| 其他   | 默认       | 黄       |

## 🛠️ 技术栈

-   **框架**：Gin Web Framework
-   **ORM**：GORM
-   **数据库**：MariaDB 10.11
-   **认证**：JWT (golang-jwt/jwt/v5)
-   **日志**：Logrus
-   **密码加密**：bcrypt
-   **CORS**：gin-contrib/cors
-   **环境变量**：godotenv

## 📝 开发规范

### 代码提交

```bash
# 提交格式
git commit -m "feat: 添加新功能"
git commit -m "fix: 修复bug"
git commit -m "docs: 更新文档"
git commit -m "refactor: 重构代码"
```

### 环境变量说明

| 变量名                  | 说明              | 示例值                                   |
| ----------------------- | ----------------- | ---------------------------------------- |
| `DB_HOST`               | 数据库主机        | `mariadb`（Docker）/ `127.0.0.1`（本地） |
| `DB_PORT`               | 数据库端口        | `3306`                                   |
| `DB_USER`               | 数据库用户        | `root`                                   |
| `DB_PASSWORD`           | 数据库密码        | `your_password`                          |
| `DB_NAME`               | 数据库名称        | `leaflet_db`                             |
| `JWT_SECRET`            | JWT 密钥          | `your_secret_key`                        |
| `MARIADB_ROOT_PASSWORD` | MariaDB root 密码 | `your_password`                          |
| `MARIADB_DATABASE`      | MariaDB 数据库名  | `leaflet_db`                             |

## 📊 数据库表结构

### users - 用户表

-   `id` - 主键
-   `username` - 用户名（唯一）
-   `password_hash` - 密码哈希
-   `created_at` - 创建时间
-   `updated_at` - 更新时间

### statuses - 心情状态表

-   `id` - 主键
-   `user_id` - 用户 ID
-   `tag_id` - 标签 ID
-   `content` - 内容
-   `leaf_color` - 树叶颜色
-   `count` - 连续天数
-   `all_record_count` - 总记录数
-   `created_at` - 创建时间

### problems - 困扰表

-   `id` - 主键
-   `sender_name` - 发送者名称
-   `user_id` - 用户 ID
-   `context` - 问题内容
-   `response` - 回应次数

### solves - 解决方案表

-   `id` - 主键
-   `user_id` - 用户 ID
-   `solution` - 解决方案
-   `problem_id` - 问题 ID

### myselfs - 个人信息表

-   `id` - 主键
-   `user_id` - 用户 ID（唯一）
-   `url` - 头像 URL
-   `profile_picture_id` - 头像 ID

## 🔒 安全说明

-   ⚠️ 生产环境请务必修改 `JWT_SECRET`
-   ⚠️ 不要将 `.env` 文件提交到版本控制
-   ⚠️ 数据库密码使用强密码
-   ⚠️ 生产环境建议使用 HTTPS

## 📄 License

本项目仅供学习交流使用是 25 级家园 HACKWEEK 第一组小登的努力。

## 👥 团队协作

小组 7 人密切合作，积极讨论，一起开会商议确定产品功能，老大朱延枫抗压，设计邱雨纳努力画图，前端周之杰快乐赶工，嘉哥李咏嘉高效输出，小学生甘宇强苦苦练习，运营 uu 万诗琴 卢艺文头脑风暴，配合默契，每个人都非常努力，朝共同目标携手共进

### 克隆并启动项目

```bash
# 1. 克隆仓库
git clone https://github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE.git
cd 25-HACK-1-Leaflet-BE

# 2. 直接启动（Docker 会自动下载所需镜像）
docker-compose up --build

# 首次启动会下载以下镜像（自动完成）：
# - golang:1.23-alpine (~300MB)（我们用的1.24.5）
# - alpine:3.20 (~8MB)
# - mariadb:10.11 (~400MB)
```

### 开发工作流

```bash
# 1. 拉取最新代码
git pull origin develop

# 2. 创建功能分支
git checkout -b feature/your-feature

# 3. 开发并测试
docker-compose up --build

# 4. 提交代码
git add .
git commit -m "feat: 添加新功能"
git push origin feature/your-feature

# 5. 创建 Pull Request
```

## 📞 联系方式

如有问题，请提交 Issue 或联系项目维护者。

---

**Happy Coding! 🎉**
