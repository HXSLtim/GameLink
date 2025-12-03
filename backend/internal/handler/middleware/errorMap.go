package middleware

import (
	"gamelink/internal/handler"
	"github.com/gin-gonic/gin"
)

// ErrorMap inspects gin errors and maps known errors to standard envelope responses.
// Only applies when handler hasn't already written a response.
func ErrorMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Written() {
			return
		}
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors[0].Err

		// 使用统一的错误映射函数，传入路径信息以获得特定的404消息
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		statusCode, response := handler.MapServiceErrorWithPath(err, path)

		// 添加请求ID到响应中
		if requestID := c.GetString("request_id"); requestID != "" && response.TraceID == "" {
			response.TraceID = requestID
		}

		c.JSON(statusCode, response)
	}
}
