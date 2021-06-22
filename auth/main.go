package main

import (
	"auth-api/domain"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
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
