package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckRole(minRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")

		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "roles not found in context"})
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid roles format"})
			return
		}

		if !hasRequiredRole(userRoles, minRole) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}

func hasRequiredRole(userRoles []string, minRole string) bool {
	rolesTop := map[string]int{
		"ROLE_ADMIN": 3,
		"ROLE_PM":    2,
		"ROLE_USER":  1,
	}

	for _, role := range userRoles {
		if rolesTop[role] >= rolesTop[minRole] {
			return true
		}
	}

	return false
}
