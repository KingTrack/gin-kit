package middleware

import "github.com/gin-gonic/gin"

type Middleware struct{}

func (m *Middleware) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
