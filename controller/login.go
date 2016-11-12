package controller

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	isGranted, err := checkLoginCredencials(username, password)
	if err != nil {
		return err
	}
	if isGranted {
		// Set custom claims
		claims := &jwtCustomClaims{
			username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(viper.GetString("Secret")))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

func checkLoginCredencials(username string, password string) (bool, error) {
	if !viper.IsSet("ACL") {
		return false, errors.New("Missing ACL config")
	}
	acl := viper.Sub("ACL")
	if acl.IsSet(username) && acl.IsSet(username+".password") {
		if acl.GetString(username+".password") == password {
			return true, nil
		}
	}
	return false, nil
}
