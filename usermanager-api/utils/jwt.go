package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID string `json:"userId,omitempty"`
	Role   string `json:"role,omitempty"`
	jwt.StandardClaims
}

func GenerateToken(userId, role string, secret string) (string, error) {
	claims := Claims{
		UserID: userId,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(secret, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, err
	}
	return nil, err
}
