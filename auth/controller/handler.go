package controller

import (
	"auth-api/domain"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler interface {
	CreateUser(*domain.User) error
	GetUser(string, string) (map[string]interface{}, error)
	AccessUser(string) (map[string]interface{}, error)
	RefreshUser(string) (map[string]interface{}, error)
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

	hashPassword, err := HashPassword(user.Password)

	if err != nil {
		fmt.Println("Error")
		return err
	}

	user.Password = hashPassword

	if err := h.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetUser(username, password string) (map[string]interface{}, error) {

	var user domain.User
	var err error

	if err = h.db.Where("username = ?", username).Find(&user).Error; err != nil {
		return nil, err
	}

	if err = ComparePassword(user.Password, password); err != nil {
		return nil, err
	}

	var accessToken *string
	if accessToken, err = GenerateToken(user.Username, user.Password, time.Now().Add(time.Minute*5).Unix()); err != nil {
		return nil, err
	}
	user.AccessToken = *accessToken

	var refreshToken *string
	if refreshToken, err = GenerateToken(user.Username, user.Password, time.Now().Add(time.Minute*15).Unix()); err != nil {
		return nil, err
	}
	user.RefreshToken = *refreshToken

	if err = h.db.Where("username = ?", user.Username).Updates(&user).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToekn": refreshToken,
	}, nil
}

func (h *Handler) AccessUser(accessToken string) (map[string]interface{}, error) {
	return nil, nil
}

func (h *Handler) RefreshUser(refreshToken string) (map[string]interface{}, error) {
	return nil, nil
}

// *FUNCTION HASH
func HashPassword(password string) (string, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)

	return string(hashPassword), err
}

func ComparePassword(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}

// *FUNCTION JWT

type Claims struct {
	Username string
	Password string
	Expire   string
	jwt.StandardClaims
}

func GenerateToken(username, password string, expire int64) (*string, error) {

	token := jwt.New(jwt.SigningMethodHS384)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["password"] = password
	claims["expire"] = expire

	tkn, err := token.SignedString([]byte("secret"))

	if err != nil {
		return nil, err
	}

	return &tkn, nil
}
