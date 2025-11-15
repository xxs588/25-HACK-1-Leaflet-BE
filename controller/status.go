package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/config"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/consts"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 根据标签ID决定树叶颜色(现在好像不用这个功能了)
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

// 创建心情状态记录
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
	location, err := time.LoadLocation("Asia/Shanghai")

	// 检查是否加载失败
	if err != nil {
		// 💡 最佳实践：如果加载失败，打印错误日志，并使用 time.Local 或 time.UTC 作为备用，防止程序崩溃
		fmt.Printf("Error loading location 'Asia/Shanghai': %v. Using time.Local instead.\n", err)
		location = time.Local // 或者 time.UTC
	}
	now := time.Now().In(location)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)

	var sumcont int64
	config.DB.Model(&model.Status{}).Where("user_id = ?", currentUserID).Count(&sumcont)

	leafColor := determineLeafColor(req.TagID) //根据标签ID决定树叶颜色（现在好像不用这个功能了）

	// 计算连续记录天数
	var consecutiveDays int64 = 1 // 默认是第1天
	yesterdayStart := todayStart.Add(-24 * time.Hour)

	var yesterdayStatus model.Status
	// 查昨天的记录
	err = config.DB.Where("user_id = ? AND created_at >= ? AND created_at < ?", currentUserID, yesterdayStart, todayStart).Order("created_at DESC").First(&yesterdayStatus).Error

	if err == nil {
		// 如果昨天有记录，连续天数+1
		consecutiveDays = int64(yesterdayStatus.Count) + 1
	}
	// 如果没有昨天的记录重新开始计数

	status := config.DB.Create(&model.Status{
		UserID:         currentUserID,
		TagID:          req.TagID,
		LeafColor:      leafColor,
		Content:        req.Content,
		Count:          consecutiveDays,   // 连续连续天数
		AllRecordCount: uint(sumcont) + 1, //加一是因为这个记录正在创建，还没存到数据库，所以后面数不到，所以加一
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
	c.JSON(http.StatusOK, gin.H{
		"message":          "状态提交成功",
		"all_record_count": uint(sumcont) + 1,
	}) //返回成功

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
	//倒序返回
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

	statusID := c.Param("id") //这个是标签表中的id
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

	//添加权限判断，避免用户修改他人的状态（由于讨论也不用弄这个了，因为看不了别人的记录所以删改都不行）
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
