package auth

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	gorm.Model

	Username    string `gorm:"unique;not null"`
	Password    string
	Active      bool
	LastLoginAt time.Time
}
