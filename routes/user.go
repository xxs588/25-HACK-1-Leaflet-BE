package routes

import (
	"us/controller"

	"github.com/gin-gonic/gin"
)


func Routes(rg *gin.Engine) {

		rg.POST("/register", controller.RegisterUser)
		rg.POST("/login", controller.LoginUser)
		// 受保护的路由示例（要鉴权的路由）
		// rg.GET("/xxxx", middlewares.JWTAuthMiddleware(), controller.xxxx)
}

	