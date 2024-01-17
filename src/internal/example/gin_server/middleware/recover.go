package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoz194/FlyXGo/src/pkg/utils/logutil"
	"runtime/debug"
)

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logutil.LogrusObj.Errorf("Panic recovered:\n%s\nStack trace:\n%s", r, debug.Stack())
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
