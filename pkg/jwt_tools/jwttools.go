package jwt_tools

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userGUID string, ipAdderss string) (string, error) {
	claims := jwt.MapClaims{
		"user_guid": userGUID,
		"ip":        ipAdderss,
		"exp":       jwt.NewNumericDate(time.Now().Add(time.Hour * 30)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte("secret"))
	return tokenString, err
}
