package response

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSON 发送 JSON 格式的成功响应
func JSON(c *gin.Context, data interface{}) {
	JSONWithStatus(c, http.StatusOK, data)
}

// JSONWithStatus 发送指定状态码的 JSON 成功响应
func JSONWithStatus(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Success(data))
}

// OK 发送标准的成功响应
func OK(c *gin.Context) {
	JSON(c, nil)
}

// Created 发送创建成功的响应
func Created(c *gin.Context, data interface{}) {
	JSONWithStatus(c, http.StatusCreated, data)
}

// Accepted 发送已接受请求的响应
func Accepted(c *gin.Context, data interface{}) {
	JSONWithStatus(c, http.StatusAccepted, data)
}

// NoContent 发送无内容响应
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Paginated 发送分页数据响应
func Paginated(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	paginationData := map[string]interface{}{
		"items": data,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      page,
			"pageSize":  pageSize,
			"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}
	JSON(c, paginationData)
}

// List 发送列表数据响应
func List(c *gin.Context, data interface{}) {
	listData := map[string]interface{}{
		"items": data,
		"count": len(data.([]interface{})),
	}
	JSON(c, listData)
}

// Message 发送自定义消息的成功响应
func Message(c *gin.Context, message string) {
	JSON(c, map[string]string{"message": message})
}

// MessageWithData 发送带数据的自定义消息响应
func MessageWithData(c *gin.Context, message string, data interface{}) {
	resp := Success(data)
	resp.Message = message
	c.JSON(http.StatusOK, resp)
}

// Token 发送认证令牌响应
func Token(c *gin.Context, token string, expiresAt interface{}) {
	JSON(c, map[string]interface{}{
		"token":     token,
		"expiresAt": expiresAt,
	})
}

// File 发送文件下载响应
func File(c *gin.Context, filename string, filepath string) {
	c.FileAttachment(filepath, filename)
}

// Stream 发送流式响应
func Stream(c *gin.Context, status int, contentType string, reader interface{}) {
	c.Header("Content-Type", contentType)
	c.Stream(func(w io.Writer) bool {
		// 这里可以根据具体的 reader 类型实现流式写入
		return false
	})
}
