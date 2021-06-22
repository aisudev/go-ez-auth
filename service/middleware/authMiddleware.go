package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			token := c.Request().Header.Get("Authorization")

			tokenSplits := strings.Split(token, " ")

			if err := VerifyAccessToken(tokenSplits[1]); err != nil {
				return c.String(http.StatusUnauthorized, err.Error())
			}

			return next(c)
		}
	}
}

// *Verify Method
func VerifyAccessToken(accessToken string) error {

	resp, err := PostUrl("http://localhost:9898/auth/access", map[string]interface{}{"accessToken": accessToken})

	if err != nil {
		return err
	}

	result := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]interface{})

	if !data["isValid"].(bool) {
		return errors.New("permission denied.")
	}

	return nil
}

// *FETCH
func PostUrl(url string, data map[string]interface{}) (*http.Response, error) {

	data_json, _ := json.Marshal(data)
	return http.Post(url, "application/json", bytes.NewBuffer(data_json))

}
