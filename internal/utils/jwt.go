package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessTokenSecret = []byte(os.Getenv("JWT_ACCESS_SECRET"))
var refreshTokenSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))

type JWTClaims struct {
	IDUser   string  `json:"id_user"`
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	Role     *string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

type CreateJWTPayload struct {
	IDUser   string
	Name     string
	Email    string
	Role     string
	IsActive bool
}

func CreateAccessToken(payload *CreateJWTPayload) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &JWTClaims{
		IDUser:   payload.IDUser,
		Name:     &payload.Name,
		Email:    &payload.Email,
		Role:     &payload.Role,
		IsActive: &payload.IsActive,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "rizz",
		},
	}

	fmt.Printf("Claims: %+v\n", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(accessTokenSecret)

	if err != nil {
		fmt.Printf("ERROR creating token: %v\n", err)
		return "", err
	}

	fmt.Printf("Token created successfully: %s...\n", tokenString[:min(len(tokenString), 50)])
	fmt.Printf("===============================\n")

	return tokenString, nil
}

func CreateRefreshToken(idUser string) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour)

	claims := &JWTClaims{
		IDUser: idUser,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "rizz-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(refreshTokenSecret)
}

func ValidateToken(tokenString string, secretKey []byte) (*JWTClaims, error) {

	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		fmt.Printf("ERROR validating token: %v\n", err)
		fmt.Printf("============================\n")
		return nil, err
	}

	if !token.Valid {
		fmt.Printf("ERROR: token is not valid\n")
		fmt.Printf("============================\n")
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return ValidateToken(tokenString, refreshTokenSecret)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
