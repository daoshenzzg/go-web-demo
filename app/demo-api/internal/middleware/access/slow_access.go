package access

import (
	"github.com/gin-gonic/gin"
	"go-web-demo/library/log"
	"time"
)

const (
	_family          = "access"
	_slowLogDuration = time.Second
)

// SlowAccess handler record slow access
func SlowAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		latency := time.Since(start)
		if latency > _slowLogDuration {
			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			log.Warn("%s|%v|%d|%v|%s|%s|%s|%s|",
				_family,
				start.Format("2006-01-02 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
				raw)
		}
	}
}
