package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	jwtSecretKey = "your_access_secret"
)

func GenerateTokens(email string) (string, string, error) {
	payload := jwt.MapClaims{
		"sub":  email,
		"exp":  time.Now().Add(30 * time.Minute).Unix(),
		"type": "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString([]byte(jwtSecretKey))

	if err != nil {
		return "", "", err
	}

	payloadRefreshToken := jwt.MapClaims{
		"sub":  email,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
		"type": "refresh",
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, payloadRefreshToken)

	tRefreshToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", "", err
	}

	return t, tRefreshToken, nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
