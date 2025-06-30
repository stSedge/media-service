package middleware

import (
	"github.com/gin-gonic/gin"
	"media-service/pkg/jwt"
	"net/http"
	"strings"
)

//tokenString := extractTokenFromRequest(c)

//claims, err := jwt.ParseToken(tokenString)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		email, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Сохраняем claims в контекст
		c.Set("user_email", email)
		c.Next()
	}
}
