package middleware

import (
	"github.com/gin-gonic/gin"
)

// ResponseMiddleware 响应中间件，统一处理响应格式
type ResponseMiddleware struct {
	enableResponseWrapper bool
}

// NewResponseMiddleware 创建响应中间件
func NewResponseMiddleware(enableWrapper bool) *ResponseMiddleware {
	return &ResponseMiddleware{
		enableResponseWrapper: enableWrapper,
	}
}

// Handle 处理响应包装
func (m *ResponseMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.enableResponseWrapper {
			c.Next()
			return
		}

		// 创建自定义的响应写入器
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           []byte{},
			statusCode:     200,
		}
		c.Writer = writer

		// 继续处理请求
		c.Next()

		// 处理响应包装逻辑
		m.processResponse(c, writer)
	}
}

// processResponse 处理响应包装
func (m *ResponseMiddleware) processResponse(c *gin.Context, writer *responseWriter) {
	// 如果已经有错误处理，不包装
	if c.IsAborted() {
		return
	}

	// 只包装 JSON 响应
	contentType := c.Writer.Header().Get("Content-Type")
	if contentType != "application/json" {
		return
	}

	// 如果状态码表示错误，包装为错误响应
	if writer.statusCode >= 400 {
		// 这里可以根据需要添加错误响应包装逻辑
		return
	}

	// 如果状态码表示成功，包装为成功响应
	if writer.statusCode >= 200 && writer.statusCode < 300 {
		// 这里可以根据需要添加成功响应包装逻辑
		return
	}
}

// responseWriter 自定义响应写入器
type responseWriter struct {
	gin.ResponseWriter
	body       []byte
	statusCode int
}

// Write 写入响应数据
func (w *responseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return w.ResponseWriter.Write(data)
}

// WriteHeader 写入状态码
func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// GetBody 获取响应体
func (w *responseWriter) GetBody() []byte {
	return w.body
}

// GetStatusCode 获取状态码
func (w *responseWriter) GetStatusCode() int {
	return w.statusCode
}
