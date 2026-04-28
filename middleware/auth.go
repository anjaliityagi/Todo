//	func AuthMiddleware() gin.HandlerFunc {
//		return func(c *gin.Context) {
//
//			token := c.GetHeader("Authorization")
//
//			if token == "" {
//				c.JSON(401, gin.H{"error": "missing token"})
//				c.Abort()
//				return
//			}
//
//			userID, err := dbHelper.GetUserIDBySession(token)
//			if err != nil {
//				c.JSON(401, gin.H{"error": "invalid session"})
//				c.Abort()
//				return
//			}
//
//			c.Set("userID", userID)
//	       c.Set("sessionId",token)
//			c.Next()
//		}
//	}
package middleware

import (
	"Todo-Server/database/dbHelper"
	"strings"

	"Todo-Server/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token format"})
			return
		}

		tokenString := parts[1]

		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		sessionID := claims.SessionID

		valid, err := dbHelper.IsSessionActive(sessionID)
		if err != nil || !valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "session expired"})
			return
		}

		suspended, err := dbHelper.IsUserSuspended(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "internal error"})
			return
		}
		if suspended {
			c.AbortWithStatusJSON(403, gin.H{"error": "user suspended"})
			return
		}

		c.Set("sessionId", sessionID)
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		role := c.GetString("role")

		for _, r := range allowedRoles {
			if role == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
	}
}
