package auth

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// VerifyTokenFromRequest verifies token from the request
func VerifyTokenFromRequest(c echo.Context, signingKey string) (*jwt.Token, error) {
	tokenString := extractToken(c)
	return VerifyToken(tokenString, signingKey)
}

// VerifyToken verifies the given token
func VerifyToken(tokenString, signingKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
}

func extractToken(c echo.Context) string {

	bearToken := c.Request().Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
