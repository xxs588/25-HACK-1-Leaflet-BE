package controller

import (
	"net/http"

	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/config"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/consts"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/middlewares"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/model"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// 注册用户
func RegisterUser(c *gin.Context) {
	var req model.RegisterRequest
	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效请求"})
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "用户名已存在"})
		return
	}

	// 创建用户信息
	user := model.User{
		Username: req.Username,
	}
	err = user.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "密码加密失败"})
		return
	}

	// 保存到数据库
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功", "user": req.Username})
	consts.Logger.WithFields(logrus.Fields{
		"username": req.Username,
		"user_id":  user.ID,
		"action":   "user_register_success",
	}).Info("用户注册成功")
}

// 用户登录
func LoginUser(c *gin.Context) {
	var req model.LoginRequest
	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效请求"})
		return
	}

	// 从数据库查找用户
	var user model.User
	err := config.DB.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := middlewares.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": token})
	// 记录成功事件
	consts.Logger.WithFields(logrus.Fields{
		"username": req.Username,
		"user_id":  user.ID,
		"action":   "user_login_success",
	}).Info("用户登录成功")
}
