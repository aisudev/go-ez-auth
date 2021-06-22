package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username     string         `gorm:"varchar(50);not null;unique" json:"username"`
	Password     string         `gorm:"varchar(512);not null" json:"password" `
	AccessToken  string         `gorm:"varchar(512);unique" json:"accesstoken"  `
	RefreshToken string         `gorm:"varchar(512);unique" json:"refreshtoken"`
	CreateAt     *time.Time     `gorm:"autoCreateTime" json"-"`
	DeleteAt     gorm.DeletedAt `json"-"`
}
