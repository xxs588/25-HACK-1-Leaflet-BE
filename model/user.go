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
