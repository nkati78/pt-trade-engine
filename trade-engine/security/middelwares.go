package security

import (
	"github.com/gin-gonic/gin"
)

const UserCtxKey = "user"

func (a Auth) TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userID, err := a.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set(UserCtxKey, userID)
	}
}
