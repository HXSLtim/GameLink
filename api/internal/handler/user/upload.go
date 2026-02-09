package user

import (
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	uploadservice "gamelink/internal/service/upload"
	"gamelink/pkg/apierr"
)

// UploadHandler 用户端上传处理器
type UploadHandler struct {
	svc *uploadservice.Service
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler(svc *uploadservice.Service) *UploadHandler {
	return &UploadHandler{svc: svc}
}

// RegisterUploadRoutes 注册用户端上传相关路由（与前端路径对齐）
func RegisterUploadRoutes(api *gin.RouterGroup, authMiddleware gin.HandlerFunc, svc *uploadservice.Service) {
	handler := NewUploadHandler(svc)

	userGroup := api.Group("/user")
	userGroup.Use(authMiddleware)
	userGroup.POST("/avatar", handler.UploadAvatar)

	certGroup := api.Group("/certification")
	certGroup.Use(authMiddleware)
	certGroup.POST("/upload/image", handler.UploadCertificationImage)

	chatGroup := api.Group("/chat")
	chatGroup.Use(authMiddleware)
	chatGroup.POST("/upload/image", handler.UploadChatImage)

	reviewGroup := api.Group("/review")
	reviewGroup.Use(authMiddleware)
	reviewGroup.POST("/upload/images", handler.UploadReviewImages)
}

// UploadAvatar 上传用户头像
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	userID := getUserIDFromContext(c)
	file, err := c.FormFile("avatar")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("missing avatar file"))
		return
	}

	res, err := h.svc.UploadImage(c, userID, file, model.UploadTypeAvatar, "./uploads/images/avatar")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondSuccess(c, "OK", gin.H{"url": res.FileURL})
}

// UploadCertificationImage 上传认证材料
func (h *UploadHandler) UploadCertificationImage(c *gin.Context) {
	userID := getUserIDFromContext(c)
	file, err := c.FormFile("image")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("missing image file"))
		return
	}

	uploadType := strings.ToLower(c.PostForm("type"))
	var mappedType model.UploadType
	var basePath string
	switch uploadType {
	case "id-card":
		mappedType = model.UploadTypeIDCard
		basePath = "./uploads/images/certification/id-card"
	case "skill-proof":
		mappedType = model.UploadTypeGameScreenshot
		basePath = "./uploads/images/certification/skill-proof"
	default:
		respondAPIError(c, apierr.BadRequest("invalid type, expected id-card or skill-proof"))
		return
	}

	res, err := h.svc.UploadImage(c, userID, file, mappedType, basePath)
	if err != nil {
		respondAPIError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondSuccess(c, "OK", gin.H{"url": res.FileURL})
}

// UploadChatImage 上传聊天图片
func (h *UploadHandler) UploadChatImage(c *gin.Context) {
	userID := getUserIDFromContext(c)
	file, err := c.FormFile("image")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("missing image file"))
		return
	}

	res, err := h.svc.UploadImage(c, userID, file, model.UploadTypeChatImage, "./uploads/images/chat")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondSuccess(c, "OK", gin.H{"url": res.FileURL})
}

// UploadReviewImages 上传评价图片（多图）
func (h *UploadHandler) UploadReviewImages(c *gin.Context) {
	userID := getUserIDFromContext(c)
	form, err := c.MultipartForm()
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid multipart form"))
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		respondAPIError(c, apierr.BadRequest("missing images"))
		return
	}

	results, err := h.svc.UploadImages(c, userID, files, model.UploadTypeReviewImage, "./uploads/images/review")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(err.Error()))
		return
	}

	urls := make([]string, 0, len(results))
	for _, res := range results {
		urls = append(urls, res.FileURL)
	}
	respondSuccess(c, "OK", gin.H{"urls": urls})
}
