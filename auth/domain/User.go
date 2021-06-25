package domain

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	UUID         uuid.UUID      `gorm:"varchar(128):not null;unique;primaryKey" json:"uuid"`
	Username     string         `gorm:"varchar(50);not null;unique" json:"username"`
	Password     string         `gorm:"varchar(512);not null" json:"password" `
	AccessToken  string         `gorm:"varchar(512);null" json:"accesstoken"  `
	RefreshToken string         `gorm:"varchar(512);null" json:"refreshtoken"`
	CreateAt     *time.Time     `gorm:"autoCreateTime" json"-"`
	DeleteAt     gorm.DeletedAt `json"-"`
}
