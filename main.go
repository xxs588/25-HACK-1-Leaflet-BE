package main

import (
	"log"

	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/config"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/consts"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 初始化日志
	consts.InitLogger()
	consts.Logger.Info("应用程序启动")
	// 首先加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		consts.Logger.Warn("未找到 .env 文件")
	}
	// 连接数据库
	config.ConnectDatabase()

	// 创建 Gin 路由器
	r := gin.Default()
	// 设置路由组
	routes.Routes(r)
	// 配置 CORS 跨域
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	// 启动服务器
	log.Println("服务器启动在端口 8080")
	r.Run(":8080")

}
