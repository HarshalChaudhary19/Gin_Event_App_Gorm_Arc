package utils

import (
	"Gin_Event_App_Arc/models"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var PrivateKey *rsa.PrivateKey // Declaring variable and inserting value via function
var PublicKey *rsa.PublicKey   // load at startup

func GenerateJWTRSA(userId int, roleList []*models.Role) (string, error) {
	err := InitKeys("private.pem", "public.pem")
	if err != nil {
		fmt.Println("Error is this btw", err.Error())
		return string(jwtKey), err
	}
	claims := &Claims{
		UserId: strconv.Itoa(userId),
		Roles:  roleList,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "key-1"
	return token.SignedString(PrivateKey)
}

func ValidateJWTRSA(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		// Verify "alg" is RS256
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return PublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// Using Function
// Initiallizing the private and public keys
func InitKeys(privatePath, publicPath string) error {
	// --- Load Private Key ---
	// fmt.Println("Path Private", privatePath)
	privBytes, err := os.ReadFile("private.pem")
	if err != nil {
		return fmt.Errorf("error reading private key file: %w", err)
	}

	privBlock, _ := pem.Decode(privBytes)
	if privBlock == nil || privBlock.Type != "PRIVATE KEY" {
		return fmt.Errorf("failed to decode PEM block containing private key")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
	if err != nil {
		return fmt.Errorf("error parsing private key: %w", err)
	}
	PrivateKey = privKey.(*rsa.PrivateKey)
	// --- Load Public Key ---
	pubBytes, err := os.ReadFile("public.pem")
	if err != nil {
		return fmt.Errorf("error reading public key file: %w", err)
	}

	pubBlock, _ := pem.Decode(pubBytes)
	if pubBlock == nil || pubBlock.Type != "PUBLIC KEY" {
		return fmt.Errorf("failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return fmt.Errorf("error parsing public key: %w", err)
	}

	switch pubKeyTyped := pubKey.(type) {
	case *rsa.PublicKey:
		PublicKey = pubKeyTyped
	default:
		return fmt.Errorf("not RSA public key")
	}

	return nil
}
