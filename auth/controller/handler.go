package controller

import (
	"auth-api/domain"
	"auth-api/utils"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtKey = []byte("secret")

type Claims struct {
	Username string
	Password string
	ExpireAt int64
	jwt.StandardClaims
}

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
		utils.Logger("", err)
		return err
	}

	user.Password = hashPassword

	if err := h.db.Create(&user).Error; err != nil {
		utils.Logger("", err)
		return err
	}

	return nil
}

func (h *Handler) GetUser(username, password string) (map[string]interface{}, error) {

	var user domain.User
	var err error

	if err = h.db.Where("username = ?", username).Find(&user).Error; err != nil {
		utils.Logger("", err)
		return nil, err
	}

	if err = ComparePassword(user.Password, password); err != nil {
		utils.Logger("", err)
		return nil, err
	}

	var accessToken *string
	if accessToken, err = GenerateToken(user.Username, user.Password, time.Now().Add(time.Minute*5).Unix()); err != nil {
		utils.Logger("", err)
		return nil, err
	}
	user.AccessToken = *accessToken

	var refreshToken *string
	if refreshToken, err = GenerateToken(user.Username, user.Password, time.Now().Add(time.Minute*15).Unix()); err != nil {
		utils.Logger("", err)
		return nil, err
	}
	user.RefreshToken = *refreshToken

	if err = h.db.Where("username = ?", user.Username).Updates(&user).Error; err != nil {
		utils.Logger("", err)
		return nil, err
	}

	return map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToekn": refreshToken,
	}, nil
}

func (h *Handler) AccessUser(accessToken string) (map[string]interface{}, error) {

	if _, err := VerifyToken(accessToken); err != nil {
		utils.Logger("", err)
		return map[string]interface{}{"isValid": false}, err
	}

	if !h.IsExist("access_token = ?", accessToken) {
		return map[string]interface{}{"isValid": false}, errors.New("not valid.")
	}

	return map[string]interface{}{"isValid": true}, nil
}

func (h *Handler) RefreshUser(refreshToken string) (map[string]interface{}, error) {

	var claims *Claims
	var err error

	if claims, err = VerifyToken(refreshToken); err != nil {
		utils.Logger("", err)
		return nil, err
	}

	if !h.IsExist("refresh_token = ?", refreshToken) {
		return map[string]interface{}{"isValid": false}, errors.New("not valid.")
	}

	var accToken *string
	if accToken, err = GenerateToken(claims.Username, claims.Password, time.Now().Add(time.Minute*5).Unix()); err != nil {
		utils.Logger("", err)
		return nil, err
	}

	var rfToken *string
	if rfToken, err = GenerateToken(claims.Username, claims.Password, time.Now().Add(time.Minute*15).Unix()); err != nil {
		utils.Logger("", err)
		return nil, err
	}

	var user domain.User
	if err = h.db.Model(&user).
		Where("username = ?", claims.Username).
		Updates(map[string]interface{}{"access_token": accToken, "refresh_token": rfToken}).
		Error; err != nil {

		utils.Logger("", err)

		return nil, err
	}

	return map[string]interface{}{"access_token": accToken, "refresh_token": rfToken}, nil
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

func GenerateToken(username, password string, expire int64) (*string, error) {

	token := jwt.New(jwt.SigningMethodHS384)

	claims := token.Claims.(jwt.MapClaims)
	claims["Username"] = username
	claims["Password"] = password
	claims["ExpireAt"] = expire

	tkn, err := token.SignedString(jwtKey)

	if err != nil {
		return nil, err
	}

	return &tkn, nil
}

func VerifyToken(token string) (*Claims, error) {

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("token is not valid.")
	}

	subTime := claims.ExpireAt - time.Now().Unix()

	if subTime <= 0 {
		return nil, errors.New("permission denied.")
	}

	return claims, nil
}

// *DB
func (h *Handler) IsExist(filter string, value interface{}) bool {

	var user domain.User
	count := int64(0)

	if err := h.db.Model(&user).Where(filter, value).Count(&count).Error; err != nil {
		return false
	}

	return count > 0
}
