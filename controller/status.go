package controller

import (
	"net/http"
	"time"

	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/config"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/consts"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 根据标签ID决定树叶颜色(后续可替换颜色)
func determineLeafColor(tagID uint) string {
	switch tagID {
	case 1: // 困倦的早八
		return "绿"
	case 2: // 自习室刷题
		return "蓝"
	case 3: // 图书馆阅读
		return "红"
	case 4: // 食堂干饭
		return "紫"
	case 5: // 备考冲刺
		return "橙"
	case 6: // 社团活动
		return "粉"
	case 7: // 情绪波动时
		return "灰"
	default:
		return "黄"
	}
}

// 创建状态条目
func CreateStatusEntry(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	} //判断用户是否存在
	currentUserID := userID.(uint)

	var user model.User
	if err := config.DB.First(&user, currentUserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}

	var req model.Status
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误或内容过长", "details": err.Error()})
		return
	} //绑定请求参数
	location, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(location)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location) //获取今天的开始时间，避免说要距上一次24小时才能发状态

	var existingStatus model.Status //声明一个状态变量

	err := config.DB.Where("user_id = ? AND created_at >= ?", currentUserID, todayStart).First(&existingStatus).Error //数据库里面查找存到真实用户id

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "今日状态已提交"}) //如果查到了就说明今天提交过了就返回冲突
		return
	}
	leafColor := determineLeafColor(req.TagID) //根据标签ID决定树叶颜色
	var count int64
	err = config.DB.Model(&model.Status{}).Where("user_id = ?", currentUserID).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取状态计数", "details": err.Error()})
		consts.Logger.WithFields(logrus.Fields{
			"username": user.Username,
			"user_id":  user.ID,
			"action":   "create_status_count_error",
			"error":    err.Error(),
		}).Error("获取状态计数失败")
		return
	}

	status := config.DB.Create(&model.Status{
		UserID:    currentUserID,
		TagID:     req.TagID,
		LeafColor: leafColor,
		Content:   req.Content,
		Count:     count + 1,
	}) //创建一个新的状态，把他存进数据库
	if status.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "心情状态保存失败", "details": status.Error.Error()})
		consts.Logger.WithFields(logrus.Fields{
			"username": user.Username,
			"user_id":  user.ID,
			"action":   "create_status",
			"error":    status.Error.Error(),
		}).Error("心情状态保存失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "状态提交成功"}) //返回成功

	// 记录成功心情记录事件
	consts.Logger.WithFields(logrus.Fields{
		"username":   user.Username,
		"user_id":    user.ID,
		"tag_id":     req.TagID,
		"leaf_color": leafColor,
		"action":     "create_status",
	}).Info("用户创建心情状态成功")
}

// 根据标签获取相同状态人数
func GetStatusesByTag(c *gin.Context) {
	tagID := c.Param("tag_id") //从路径参数获取tagID
	var uniqueUsersCount int64
	err := config.DB.Model(&model.Status{}).Where("tag_id = ?", tagID).Distinct("user_id").Count(&uniqueUsersCount).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取状态", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": uniqueUsersCount}) //返回状态和数量
}

// 查询个人所有记录
func GetStatus(c *gin.Context) {
	currentUserID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或者令牌无效"})
		return
	}
	var status []model.Status
	if err := config.DB.Where("user_id = ?", currentUserID).Order("created_at desc").Find(&status).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "[]model.Status{}"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}

// 删除
func DeleteStatus(c *gin.Context) {
	currentUserID := c.GetUint("user_id")
	var user model.User
	if err := config.DB.First(&user, currentUserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}

	statusID := c.Param("id")
	var status model.Status
	if err := config.DB.Where("id = ?", statusID).First(&status).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "个人心情记录状态未找到"})
		return
	}
	//添加权限判断，避免用户删除他人的状态
	if status.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此个人心情记录状态"})
		consts.Logger.WithFields(logrus.Fields{
			"username":  user.Username,
			"user_id":   user.ID,
			"status_id": statusID,
			"action":    "unauthorized_delete_attempt",
		}).Warn("用户尝试删除不属于自己的心情状态")
		return
	}
	if err := config.DB.Delete(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "个人心情记录状态删除失败", "details": err.Error()})
		consts.Logger.WithFields(logrus.Fields{
			"username":  user.Username,
			"user_id":   user.ID,
			"status_id": statusID,
			"action":    "delete_status",
			"error":     err.Error(),
		}).Error("心情状态删除失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "个人心情记录状态删除成功"})

	// 记录成功事件
	consts.Logger.WithFields(logrus.Fields{
		"username":  user.Username,
		"user_id":   user.ID,
		"status_id": statusID,
		"action":    "delete_status",
	}).Info("用户删除心情状态成功")
}

// 编辑状态（这个id还是表里面这个记录的ID）
func UpdateStatus(c *gin.Context) {
	currentUserID := c.GetUint("user_id")
	var user model.User
	if err := config.DB.First(&user, currentUserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}

	statusID := c.Param("id")
	var req model.Status
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数ID错误或内容过长", "details": err.Error()})
		return
	}
	var status model.Status
	if err := config.DB.Where("id = ?", statusID).First(&status).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没找到这条记录"})
		return
	}

	//添加权限判断，避免用户修改他人的状态
	if status.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限编辑此个人心情记录状态"})
		consts.Logger.WithFields(logrus.Fields{
			"username":  user.Username,
			"user_id":   user.ID,
			"status_id": statusID,
			"action":    "unauthorized_update_attempt",
		}).Warn("用户尝试修改不属于自己的心情状态")
		return
	}

	status.Content = req.Content
	status.TagID = req.TagID
	status.LeafColor = determineLeafColor(req.TagID)
	if err := config.DB.Save(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新状态失败", "details": err.Error()})
		consts.Logger.WithFields(logrus.Fields{
			"username":  user.Username,
			"user_id":   user.ID,
			"status_id": statusID,
			"action":    "update_status",
			"error":     err.Error(),
		}).Error("心情状态更新失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新状态成功"})

	// 记录成功事件
	consts.Logger.WithFields(logrus.Fields{
		"username":  user.Username,
		"user_id":   user.ID,
		"status_id": statusID,
		"action":    "update_status",
	}).Info("用户更新心情状态成功")
}
