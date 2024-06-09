package models

import (
	"database/sql"
	"time"
)

// 登录用的数据模型
type LoginUser struct {
	Email     string `json:"email" binding:"required"`
	Encrypted string `json:"encrypted" binding:"required"`
	Nanoid    string `json:"nanoid" binding:"required"`
}
type EmailUpdater struct {
	LoginUser
	NewEmail   string `json:"newEmail" binding:"required"`
	VerifyCode string `json:"verifyCode" binding:"required"`
}

// 返回安全的字段 过滤敏感字段
type ResponseUser struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	Email     string `json:"email" binding:"required"`
	Encrypted string `json:"-"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime"`
}
type User struct {
	ID        uint         `json:"id" gorm:"primarykey"`
	CreatedAt int64        `json:"-" gorm:"autoCreateTime"`
	UpdatedAt int64        `json:"-" gorm:"autoUpdateTime"`
	DeletedAt sql.NullTime `json:"-" gorm:"index"`

	Email     string `json:"email" binding:"required" gorm:"uniqueIndex;not null"` //gorm索引标签
	Encrypted string `json:"encrypted" binding:"required" gorm:"not null"`
}
type PersonalInfo struct {
	Nickname string    `json:"nickname,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
	IconfontLink
}
type RegisterUser struct {
	ID         uint   `json:"_" gorm:"primarykey"`
	VerifyCode string `json:"verifyCode" binding:"required"`
	LoginUser
	PersonalInfo
}

type PasswordResetUser struct {
	VerifyCode string `json:"verifyCode" binding:"required"`
	LoginUser
}

type Token struct {
	Token string `json:"token" binding:"required"`
}
