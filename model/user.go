package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户
type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"unique;not null"`
	PasswordHash string `json:"password_hash" gorm:"not null"`
}

// 注册
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 校园场景设置和标签选择
type Status struct {
	gorm.Model
	UserID    int `gorm:"index:idx_user_time"` // 用户ID uint类型非负整数，索引
	TagID     int
	MoodType  string //情绪气泡类型
	Content   string `gorm:"type:text;size: 100"` //限制100字
	LeafColor string
}

type Tag struct {
	gorm.Model
	TagName string `gorm:"unique;not null"` //标签名称
}

// 创建状态请求
type CreateStatusRequest struct {
	MoodType  string `json:"mood_type" binding:"required"`
	TagID     int    `json:"tag_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
	LeafColor string `json:"leaf_color" binding:"required"`
}

// 更新状态请求
type UpdateStatusRequest struct {
	MoodType  string `json:"mood_type" binding:"required"`
	TagID     int    `json:"tag_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
	LeafColor string `json:"leaf_color" binding:"required"`
}

// 密码加密
func (u *User) HashPassword(password string) (err error) {
	// 生成哈希值
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)

	return nil
}

// 密码验证
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
