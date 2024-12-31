package jwt

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaim struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateToken(userID int, expirationTime time.Time) (string, error) {
	claims := CustomClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(r *http.Request) error {
	tokenString, err := GetToken(r)
	if err != nil {
		return errors.New("token is missing")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if token.Valid {
		return nil
	}
	return err
}

func GetClaims(r *http.Request) (*CustomClaim, error) {
	tokenString, err := GetToken(r)
	if err != nil {
		return nil, errors.New("token is missing")
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if claims, ok := token.Claims.(*CustomClaim); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func GetToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("authorization header is not in 'Bearer <token>' format")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return "", errors.New("token is missing")
	}

	return tokenString, nil
}
