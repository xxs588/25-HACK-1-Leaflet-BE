package controller

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"strconv"

	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/config"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/consts"
	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/model"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// 上传问题
func UploadProblem(c *gin.Context) {

	// 从上下文中获取用户信息
	claims, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户信息未授权"})
		return
	}
	var user model.User
	if err := config.DB.First(&user, claims.(uint)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}
	req := model.Problem{
		SenderName: user.Username,
		Response:   0,
		UserID:     claims.(uint),
	}
	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效请求"})
		return
	}
	// 保存问题到数据库
	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		consts.Logger.WithFields(logrus.Fields{
			"username": req.SenderName,
			"user_id":  user.ID,
			"action":   "user_upload_problem",
		}).Error("用户上传问题失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "上传成功"})

	// 记录成功事件
	consts.Logger.WithFields(logrus.Fields{
		"username": req.SenderName,
		"user_id":  user.ID,
		"action":   "user_upload_problem",
	}).Info("用户上传问题成功")
}

// 修改问题
func UpdateProblem(c *gin.Context) {
	claims, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户信息未授权"})
		return
	}
	var user model.User
	if err := config.DB.First(&user, claims.(uint)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}

	// 从URL参数获取问题ID
	problemIDStr := c.Param("id")
	if problemIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "问题ID不能为空"})
		return
	}

	// 将字符串转换为uint
	problemID, err := strconv.ParseUint(problemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问题ID"})
		return
	}

	// 定义请求体结构，只需要包含要更新的内容
	var updateReq struct {
		Context string `json:"context" binding:"required"`
	}

	// 绑定请求体
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效请求"})
		return
	}

	// 验证问题是否存在且属于当前用户
	var existingProblem model.Problem
	if err := config.DB.First(&existingProblem, uint(problemID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "问题不存在"})
		return
	}

	// 检查问题是否属于当前用户
	if existingProblem.UserID != claims.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改此问题"})
		consts.Logger.WithFields(logrus.Fields{
			"username":   user.Username,
			"user_id":    user.ID,
			"problem_id": problemID,
			"action":     "unauthorized_update_at+-tempt",
		}).Warn("用户尝试修改不属于自己的问题")
		return
	}

	// 更新问题到数据库
	if err := config.DB.Model(&existingProblem).Updates(model.Problem{Context: updateReq.Context}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		consts.Logger.WithFields(logrus.Fields{
			"username": user.Username,
			"user_id":  user.ID,
			"action":   "user_update_problem",
		}).Error("用户更新问题失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})

	// 记录成功事件
	consts.Logger.WithFields(logrus.Fields{
		"username": user.Username,
		"user_id":  user.ID,
		"action":   "user_update_problem",
	}).Info("用户更新问题成功")
}

// 删除问题
func DeleteProblem(c *gin.Context) {
	claims, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户信息未授权"})
		return
	}
	var user model.User
	if err := config.DB.First(&user, claims.(uint)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}

	// 从URL参数获取问题ID
	problemIDStr := c.Param("id")
	consts.Logger.WithFields(logrus.Fields{
		"action":         "delete_problem_debug",
		"problem_id_str": problemIDStr,
	}).Info("获取到的问题ID字符串")

	if problemIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "问题ID不能为空"})
		return
	}

	// 将字符串转换为uint
	problemID, err := strconv.ParseUint(problemIDStr, 10, 32)
	if err != nil {
		consts.Logger.WithFields(logrus.Fields{
			"action":         "delete_problem_debug",
			"problem_id_str": problemIDStr,
			"error":          err.Error(),
		}).Error("问题ID转换失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问题ID"})
		return
	}

	consts.Logger.WithFields(logrus.Fields{
		"action":     "delete_problem_debug",
		"problem_id": problemID,
		"user_id":    claims.(uint),
	}).Info("开始查找问题")

	// 验证问题是否存在且属于当前用户
	var existingProblem model.Problem
	if err := config.DB.First(&existingProblem, uint(problemID)).Error; err != nil {
		consts.Logger.WithFields(logrus.Fields{
			"action":     "delete_problem_debug",
			"problem_id": problemID,
			"error":      err.Error(),
		}).Error("查找问题失败")
		c.JSON(http.StatusNotFound, gin.H{"error": "问题不存在"})
		return
	}

	consts.Logger.WithFields(logrus.Fields{
		"action":           "delete_problem_debug",
		"problem_id":       problemID,
		"problem_owner_id": existingProblem.UserID,
		"current_user_id":  claims.(uint),
	}).Info("检查问题权限")

	// 检查问题是否属于当前用户
	if existingProblem.UserID != claims.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此问题"})
		consts.Logger.WithFields(logrus.Fields{
			"username":   user.Username,
			"user_id":    user.ID,
			"problem_id": problemID,
			"action":     "unauthorized_delete_attempt",
		}).Warn("用户尝试删除不属于自己的问题")
		return
	}

	consts.Logger.WithFields(logrus.Fields{
		"action":     "delete_problem_debug",
		"problem_id": problemID,
	}).Info("开始删除问题")

	// 删除问题从数据库
	if err := config.DB.Delete(&existingProblem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		consts.Logger.WithFields(logrus.Fields{
			"username": user.Username,
			"user_id":  user.ID,
			"action":   "user_delete_problem",
			"error":    err.Error(),
		}).Error("用户删除问题失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})

	// 记录成功事件
	consts.Logger.WithFields(logrus.Fields{
		"username": user.Username,
		"user_id":  user.ID,
		"action":   "user_delete_problem",
	}).Info("用户删除问题成功")
}

// 获取问题
func GetProblems(c *gin.Context) {
	claims, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户信息未授权"})
		return
	}
	var user model.User
	if err := config.DB.First(&user, claims.(uint)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}
	var allProblems []model.Problem
	// 获取所有问题
	if err := config.DB.Find(&allProblems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果问题数量少于3条，返回所有问题
	if len(allProblems) <= 3 {
		c.JSON(http.StatusOK, gin.H{"problems": allProblems})
		return
	}

	// 使用 crypto/rand 随机选择3个不重复的索引
	indices := make([]int, 3)
	for i := 0; i < 3; i++ {
		for {
			// 生成安全的随机数
			nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(allProblems))))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "随机数生成失败"})
				return
			}
			idx := int(nBig.Int64())

			// 检查索引是否已经使用过
			duplicate := false
			for j := 0; j < i; j++ {
				if indices[j] == idx {
					duplicate = true
					break
				}
			}
			if !duplicate {
				indices[i] = idx
				break
			}
		}
	}

	// 根据随机索引获取问题
	problems := make([]model.Problem, 3)
	for i, idx := range indices {
		problems[i] = allProblems[idx]
	}

	c.JSON(http.StatusOK, gin.H{"problems": problems})
	consts.Logger.WithFields(logrus.Fields{
		"action":   "get_problems",
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("获取问题成功")
}

// 修改指定ID的问题内容
func ChangeProblem(c *gin.Context) {
	claims, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户信息未授权"})
		return
	}
	var user model.User
	if err := config.DB.First(&user, claims.(uint)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}

	// 从URL参数获取问题ID
	problemIDStr := c.Param("id")
	if problemIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "问题ID不能为空"})
		return
	}

	// 将字符串转换为uint
	problemID, err := strconv.ParseUint(problemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问题ID"})
		return
	}

	// 定义请求体结构，只需要包含要更新的内容
	var changeReq struct {
		Content string `json:"content" binding:"required"`
	}

	// 绑定请求体
	if err := c.ShouldBindJSON(&changeReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效请求"})
		return
	}

	// 验证问题是否存在且属于当前用户
	var existingProblem model.Problem
	if err := config.DB.First(&existingProblem, uint(problemID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "问题不存在"})
		return
	}

	// 检查问题是否属于当前用户
	if existingProblem.UserID != claims.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改此问题"})
		consts.Logger.WithFields(logrus.Fields{
			"username":   user.Username,
			"user_id":    user.ID,
			"problem_id": problemID,
			"action":     "unauthorized_change_attempt",
		}).Warn("用户尝试修改不属于自己的问题")
		return
	}

	// 更新问题内容到数据库
	if err := config.DB.Model(&existingProblem).Updates(model.Problem{Context: changeReq.Content}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		consts.Logger.WithFields(logrus.Fields{
			"username": user.Username,
			"user_id":  user.ID,
			"action":   "user_change_problem",
		}).Error("用户修改问题内容失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "修改成功"})

	// 记录成功事件
	consts.Logger.WithFields(logrus.Fields{
		"username": user.Username,
		"user_id":  user.ID,
		"action":   "user_change_problem",
	}).Info("用户修改问题内容成功")
}

func UploadSolve(c *gin.Context) {
	// 从上下文中获取用户信息
	claims, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户信息未授权"})
		return
	}
	var user model.User
	if err := config.DB.First(&user, claims.(uint)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}
	
	// 从URL参数获取问题ID
	problemIDStr := c.Param("id")
	if problemIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "问题ID不能为空"})
		return
	}
	
	// 验证问题是否存在
	problemID, err := strconv.ParseUint(problemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问题ID"})
		return
	}
	
	var problem model.Problem
	if err := config.DB.First(&problem, uint(problemID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "问题不存在"})
		return
	}
	
	req := model.Solve{
		ProblemID:    problemIDStr,
		UserID:        claims.(uint),
	}
	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效请求"})
		return
	}
	
	// 开始数据库事务
	tx := config.DB.Begin()
	
	// 保存解决方案到数据库
	if err := tx.Create(&req).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		consts.Logger.WithFields(logrus.Fields{
			"username": user.Username,
			"user_id":  user.ID,
			"action":   "user_upload_solve",
		}).Error("用户上传解决方案失败")
		return
	}
	
	// 更新问题的回应次数
	if err := tx.Model(&problem).Update("response", problem.Response+1).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新问题回应次数失败"})
		consts.Logger.WithFields(logrus.Fields{
			"username":   user.Username,
			"user_id":    user.ID,
			"problem_id": problemID,
			"action":     "update_problem_response",
		}).Error("更新问题回应次数失败")
		return
	}
	
	// 提交事务
	tx.Commit()
	
	c.JSON(http.StatusOK, gin.H{"message": "上传成功"})

	// 记录成功事件
	consts.Logger.WithFields(logrus.Fields{
		"username":    user.Username,
		"user_id":     user.ID,
		"problem_id":  problemID,
		"action":      "user_upload_solve",
		"new_response": problem.Response + 1,
	}).Info("用户上传解决方案成功")
}

func GetSolves(c *gin.Context) {
	claims, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户信息未授权"})
		return
	}
	var user model.User
	if err := config.DB.First(&user, claims.(uint)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}
	
	// 首先获取当前用户上传的所有问题
	var userProblems []model.Problem
	if err := config.DB.Where("user_id = ?", claims.(uint)).Find(&userProblems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// 如果用户没有上传任何问题，返回空数组
	if len(userProblems) == 0 {
		c.JSON(http.StatusOK, gin.H{"solves": []model.Solve{}})
		return
	}
	
	// 收集所有问题ID
	problemIDs := make([]string, len(userProblems))
	for i, problem := range userProblems {
		problemIDs[i] = strconv.FormatUint(uint64(problem.ID), 10)
	}
	
	// 获取这些问题对应的所有解决方案
	var solves []model.Solve
	if err := config.DB.Where("problem_id IN ?", problemIDs).Find(&solves).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"solves": solves})
	consts.Logger.WithFields(logrus.Fields{
		"action":     "get_solves",
		"user_id":    user.ID,
		"username":   user.Username,
		"problem_count": len(userProblems),
		"solve_count":  len(solves),
	}).Info("获取用户问题的解决方案成功")
}