package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ximofam/chess-online/services/auth/token"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, "Authorization format must be Bearer <token>")
			return
		}

		accessToken := parts[1]

		claims, err := token.ParseAccessToken(accessToken)
		if err != nil {
			c.JSON(401, err.Error())
			return
		}

		c.Set(userIDKey, claims.UserID)

		c.Next()
	}
}

func AuthenticateByPathValue() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Query("token")

		claims, err := token.ParseAccessToken(accessToken)
		if err != nil {
			c.JSON(401, err.Error())
			return
		}

		c.Set(userIDKey, claims.UserID)

		c.Next()
	}
}
