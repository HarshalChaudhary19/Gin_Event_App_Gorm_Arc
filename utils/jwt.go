package utils

import (
	"Gin_Event_App_Arc/models"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("SecretKey")

type Claims struct {
	UserId string         `json:"user"`
	Roles  []*models.Role `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateJWT(userId int, roleList []*models.Role) (string, error) {
	fmt.Println("Public Key ye hai", PublicKey)
	claims := &Claims{
		UserId: strconv.Itoa(userId),
		Roles:  roleList,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
