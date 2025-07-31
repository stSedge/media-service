package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"media-service/pkg/config"
	"time"
)

func GenerateTokens(email string) (string, string, uuid.UUID, error) {
	payload := jwt.MapClaims{
		"sub":  email,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
		"type": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(config.Cnfg.JwtSecret))
	if err != nil {
		return "", "", uuid.Nil, err
	}

	jti := uuid.New()
	payloadRefreshToken := jwt.MapClaims{
		"sub":  email,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
		"type": "refresh",
		"jti":  jti.String(),
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, payloadRefreshToken)
	tRefreshToken, err := token.SignedString([]byte(config.Cnfg.JwtSecret))
	if err != nil {
		return "", "", uuid.Nil, err
	}

	return t, tRefreshToken, jti, nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Cnfg.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
