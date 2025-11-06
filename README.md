# 25-HACK-1-Leaflet-BE

一个基于 Go 和 Gin 框架的后端项目，支持用户注册和登录功能。

## 项目结构

```
├── main.go              # 应用入口
├── config/              # 配置相关
│   └── db.go           # 数据库连接配置
├── controller/          # 控制器层
│   └── user.go         # 用户相关控制器
├── middlewares/         # 中间件
│   └── auth.go         # JWT认证中间件
├── model/               # 数据模型
│   └── user.go         # 用户模型
├── routes/              # 路由定义
│   └── user.go         # 用户路由
└── .env                 # 环境配置（不提交）
```

## Docker 化部署

### 开发环境

````bash
# 复制配置模板
cp .env.example .env

# 编辑.env文件，设置您的数据库密码

## 快速开始

1. 设置环境配置：
   ```bash
   cp .env.example .env
   # 编辑.env文件，设置您的数据库密码

# 启动开发环境
docker-compose up -d

# 或者直接使用
go run main.go
````

## 协作开发说明

### 1. 环境配置

-   使用`.env.example`作为配置模板
-   个人`.env`文件不提交到版本控制

### 2. 启动应用

```bash
# 方式一：使用Docker
docker-compose up -d

# 方式二：本地运行
go run main.go
```

## API 接口

-   `POST /register` - 用户注册
-   `POST /login` - 用户登录
