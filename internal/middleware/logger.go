package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func Logger(log *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        c.Next()

        latency := time.Since(start)
        statusCode := c.Writer.Status()

        log.Info("request completed",
            zap.String("path", path),
            zap.Int("status", statusCode),
            zap.Duration("latency", latency),
            zap.String("ip", c.ClientIP()),
            zap.String("method", c.Request.Method),
        )
    }
}