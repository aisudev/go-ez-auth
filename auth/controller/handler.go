package controller

import (
	"auth-api/domain"

	"gorm.io/gorm"
)

type UserHandler interface {
	CreateUser(*domain.User) error
	GetUser(string, string) map[string]interface{}
	AccessUser(string) map[string]interface{}
	RefreshUser(string) map[string]interface{}
}

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) UserHandler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) CreateUser(user *domain.User) error {
	return nil
}

func (h *Handler) GetUser(username, password string) map[string]interface{} {
	return nil
}

func (h *Handler) AccessUser(accessToken string) map[string]interface{} {
	return nil
}

func (h *Handler) RefreshUser(refreshToken string) map[string]interface{} {
	return nil
}
