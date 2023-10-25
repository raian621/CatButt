package server

import "github.com/gin-gonic/gin"

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
