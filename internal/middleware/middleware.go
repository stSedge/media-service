package middleware

import (
	"github.com/gin-gonic/gin"
	"media-service/internal/services"
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

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		tokenType, ok := claims["type"].(string)

		if !ok || tokenType != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token type: expected access token"})
			return
		}

		email, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "subject not found in token"})
			return
		}

		c.Set("user_email", email)

		user, err := services.GetUserByMail(email)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found or unauthorized"})
		}
		c.Set("roles", user.Roles)

		c.Next()
	}
}
