package admin

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
)

// RegisterUploadRoutes 注册上传路由（图片）
func RegisterUploadRoutes(router gin.IRouter) {
	group := router.Group("/upload")
	group.POST("/image", UploadImageHandler)
}

// UploadImageHandler 处理图片上传
// @Summary      图片上传
// @Description  上传图片文件，返回保存路径
// @Tags         Admin - Upload
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file  formData  file  true  "Image file"
// @Success      200   {object}  model.APIResponse[map[string]string]
// @Router       /admin/upload/image [post]
func UploadImageHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: "missing file"})
		return
	}

	// 基础安全校验
	const maxSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: "file too large (max 5MB)"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	defer f.Close()

	head := make([]byte, 512)
	n, _ := f.Read(head)
	contentType := http.DetectContentType(head[:n])
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowed[contentType] {
		c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: "unsupported file type"})
		return
	}
	// reset read pointer for saver
	if _, err := f.Seek(0, 0); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	cfg := middleware.GetImageConfig()
	// ensure path exists
	if err := os.MkdirAll(cfg.UploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse[any]{Success: false, Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	res, err := middleware.SaveFile(c, file, cfg)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse[any]{Success: false, Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[map[string]string]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]string{
			"filePath": res.FilePath,
			"hash":     res.Hash,
		},
	})
}
