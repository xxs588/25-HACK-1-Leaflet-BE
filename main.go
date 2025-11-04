package main

import (
	"log"
	"us/config"
	"us/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 首先加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到 .env 文件")
	}
	// 连接数据库
	config.ConnectDatabase()

	// 创建 Gin 路由器
	r := gin.Default()
	// 设置路由组
	routes.Routes(r)

	// 启动服务器
	log.Println("服务器启动在端口 8080")
	r.Run(":8080")

}