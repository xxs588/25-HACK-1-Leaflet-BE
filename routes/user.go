package routes

import (
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/controller"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/middlewares"

	"github.com/gin-gonic/gin"
)

func Routes(rg *gin.Engine) {
	// 公开路由
	rg.POST("/register", controller.RegisterUser)
	rg.POST("/login", controller.LoginUser)
	rg.GET("/encouragements", controller.GetEncouragements)

	// 需要JWT认证的路由组
	auth := rg.Group("/")
	auth.Use(middlewares.JWTAuthMiddleware())
	{
		// Mind 相关路由
		mind := auth.Group("/mind")
		{
			mind.POST("/", controller.UploadProblem)
			mind.GET("", controller.GetProblems)
			mind.PUT("/:id", controller.ChangeProblem)
			mind.DELETE("/:id", controller.DeleteProblem)
		}

		// Status 相关路由
		status := auth.Group("/status")
		{
			status.POST("", controller.CreateStatusEntry)
			status.GET("/by_tag/:tag_id", controller.GetStatusesByTag)
			status.GET("/by_user/:user_id", controller.GetStatus)
			status.PUT("/:id", controller.UpdateStatus)
			status.DELETE("/:id", controller.DeleteStatus)
		}

		// Solve 相关路由
		solve := auth.Group("/solve")
		{
			solve.POST("/:id", controller.UploadSolve)
			solve.GET("", controller.GetSolves)
		}
	}
}
