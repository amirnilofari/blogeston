package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/amirnilofari/hash-go-mysql/models"
	"github.com/amirnilofari/hash-go-mysql/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided!"})
			c.Abort()
			return
		}

		tokenString = strings.TrimSpace(strings.Replace(tokenString, "Bearer", "", 1))

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return utils.JwtSecretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invaild token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["user_id"].(float64)
			if ok {
				c.Set("user_id", int(userID))
				var user models.User
				err := utils.DB.QueryRowContext(c, "SELECT user_id, role FROM users WHERE user_id = $1", int(userID)).Scan(&user.ID, &user.Role)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user" + err.Error()})
					return
				}
				fmt.Println("user role:", user.Role)
				if user.Role == "admin" || user.Role == "user" {
					c.Set("user_role", user.Role)
				} else {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user role in token"})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				return
			}

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
