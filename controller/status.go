package controller

import (
	"net/http"
	"time"

	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/config"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/model"
	"github.com/gin-gonic/gin"
)

func determineLeafColor(moodType string) string {
	switch moodType {
	case "happy":
		return "绿"
	case "sad":
		return "蓝"
	case "angry":
		return "红"
	default:
		return "黄"
	}
}

func CreatStatusEntry(c *gin.Context) {
	userID, exists := c.Get("user_ID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	} //判断用户是否存在
	currentUserID := userID.(uint)

	var req model.CreateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误或内容过长", "details": err.Error()})
		return
	} //绑定请求参数
	todayStart := time.Now().Truncate(24 * time.Hour)
	var existingStatus model.Status //声明一个状态变量

	err := config.DB.Where("user_id = ? AND created_at >= ?", currentUserID, todayStart).First(&existingStatus).Error //数据库里面查找存到真实用户id

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "今日状态已提交"}) //如果查到了就说明今天提交过了就返回冲突
		return
	}
	leafColor := determineLeafColor(req.MoodType)

	status := config.DB.Create(&model.Status{
		UserID:    int(currentUserID),
		MoodType:  req.MoodType,
		LeafColor: leafColor,
		TagID:     req.TagID,
		Content:   req.Content,
	}) //创建一个新的状态，把他存进数据库
	if status.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "心情状态保存失败", "details": status.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "状态提交成功"}) //返回成功
}

func GetStatusesByTag(c *gin.Context) {
	tagID := c.Param("tag_id") //从路径参数获取tagID
	var Statuses []model.Status
	if err := config.DB.Model(&model.Status{}). /*指定要查询的表的空的结构体指针*/ Where("tag_id = ?", tagID).Distinct().Find(&Statuses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取状态", "details": err.Error()})
		return
	}
	length := len(Statuses)
	c.JSON(http.StatusOK, gin.H{"data": Statuses, "count": length}) //返回状态和数量
}
