package main

import (
	"auth-api/domain"
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
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := handler.CreateUser(&user); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.String(http.StatusOK, "ok")
	})

	// *SignIn
	e.POST("/signin", func(c echo.Context) error {

		var reqMap map[string]interface{}
		var resMap map[string]interface{}
		var err error

		if err = c.Bind(&reqMap); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if resMap, err = handler.GetUser(reqMap["username"].(string), reqMap["password"].(string)); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, resMap)

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
