package middleware

import (
	"Todo-Server/database/dbHelper"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(401, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		userID, err := dbHelper.GetUserIDBySession(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid session"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("sessionId", token)
		c.Next()
	}
}
