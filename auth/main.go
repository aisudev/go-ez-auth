package main

import (
	"auth-api/domain"
	"auth-api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	Handler "auth-api/controller"
)

var SqliteDB *gorm.DB

//Initial
func init() {
	ConnectingDatabase()
	AutoMigrate()
}

//Main
func main() {
	e := echo.New()

	handler := Handler.NewHandler(SqliteDB)

	// *SignUp
	e.POST("/signup", func(c echo.Context) error {

		var user domain.User

		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, utils.Response(false, nil, nil, err))
		}

		if err := handler.CreateUser(&user); err != nil {
			return c.JSON(http.StatusBadRequest, utils.Response(false, nil, nil, err))
		}

		return c.JSON(http.StatusBadRequest, utils.Response(true, nil, nil, nil))
	})

	// *SignIn
	e.POST("/signin", func(c echo.Context) error {

		var reqMap map[string]interface{}
		var resMap map[string]interface{}
		var err error

		if err = c.Bind(&reqMap); err != nil {
			return c.JSON(http.StatusBadRequest, utils.Response(false, nil, nil, err))
		}

		if resMap, err = handler.GetUser(reqMap["username"].(string), reqMap["password"].(string)); err != nil {
			return c.JSON(http.StatusBadRequest, utils.Response(false, nil, nil, err))
		}

		return c.JSON(http.StatusOK, utils.Response(true, nil, resMap, nil))

	})

	// *Verify AccessToken
	e.POST("/auth/access", func(c echo.Context) error {
		var reqMap map[string]interface{}
		var resMap map[string]interface{}
		var err error

		if err = c.Bind(&reqMap); err != nil {
			return c.JSON(http.StatusBadRequest, utils.Response(false, nil, nil, err))
		}

		if resMap, err = handler.AccessUser(reqMap["accessToken"].(string)); err != nil {
			return c.JSON(http.StatusUnauthorized, utils.Response(false, nil, resMap, err))
		}

		return c.JSON(http.StatusOK, utils.Response(true, nil, resMap, nil))

	})

	// *Refresh Token
	e.POST("/auth/refresh", func(c echo.Context) error {
		var reqMap map[string]interface{}
		var resMap map[string]interface{}
		var err error

		if err = c.Bind(&reqMap); err != nil {
			return c.JSON(http.StatusBadRequest, utils.Response(false, nil, nil, err))
		}

		if resMap, err = handler.RefreshUser(reqMap["refreshToken"].(string)); err != nil {
			return c.JSON(http.StatusUnauthorized, utils.Response(false, nil, nil, err))
		}

		return c.JSON(http.StatusOK, utils.Response(true, nil, resMap, nil))

	})

	e.Start(":9898")
}

//Function
func ConnectingDatabase() {

	var err error

	SqliteDB, err = gorm.Open(sqlite.Open("db/auth.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

}

func AutoMigrate() {
	SqliteDB.AutoMigrate(&domain.User{})
}
