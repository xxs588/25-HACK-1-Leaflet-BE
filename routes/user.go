package routes

import (
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/controller"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/middlewares"

	"github.com/gin-gonic/gin"
)

func Routes(rg *gin.Engine) {

	rg.POST("/register", controller.RegisterUser)
	rg.POST("/login", controller.LoginUser)
	// 受保护的路由示例（要鉴权的路由）
	// rg.GET("/xxxx", middlewares.JWTAuthMiddleware(), controller.xxxx)
	// 注意：带参数的路由应该在无参数路由之前定义，避免路由冲突
	rg.PUT("/mind/:id", middlewares.JWTAuthMiddleware(), controller.ChangeProblem)
	rg.DELETE("/mind/:id", middlewares.JWTAuthMiddleware(), controller.DeleteProblem)
	rg.POST("/mind/", middlewares.JWTAuthMiddleware(), controller.UploadProblem)
	rg.GET("/mind", middlewares.JWTAuthMiddleware(), controller.GetProblems)
}
