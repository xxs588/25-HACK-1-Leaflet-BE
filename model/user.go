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

// 标签表（校园场景标签）
type Tag struct {
	gorm.Model
	TagName string `json:"tag_name" gorm:"unique;not null"` // 标签名称
}

// 心情状态表
type Status struct {
	gorm.Model
	UserID    uint   `json:"user_id" gorm:"index:idx_user_time;not null"` // 用户ID
	TagID     uint   `json:"tag_id" binding:"required" gorm:"not null"`   // 标签ID (1=困倦的早八, 2=自习室刷题, 等等)
	Content   string `json:"content" binding:"required" gorm:"type:text"` // 内容
	LeafColor string `json:"leaf_color"`                                  // 树叶颜色（自动生成）
}

// 以下为情绪互动模块
type Problem struct {
	gorm.Model
	SenderName string `json:"sender_name" gorm:"not null"` // 外键，关联用户发送人
	UserID     uint   `json:"user_id" gorm:"not null"`     // 鉴权用户ID
	Context    string `json:"context" gorm:"not null"`     // 问题
	Response   uint   `json:"response" gorm:"not null"`    // 回应次数
}

type Solve struct {
	gorm.Model
	UserID    uint   `json:"user_id" gorm:"not null"`    // 外键，关联用户解决者
	Solution  string `json:"solution" gorm:"not null"`   //解决方案
	ProblemID string `json:"problem_id" gorm:"not null"` //问题ID
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
