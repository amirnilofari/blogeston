package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve user role from context (set during authentication)
		role, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			return
		}

		// Check if the user's role is in the allowed roles
		userRole := role.(string)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		// User does not have the required role
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access forbidden: insufficient permissions"})
	}
}
