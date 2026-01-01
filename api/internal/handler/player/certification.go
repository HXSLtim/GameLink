package player

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "gamelink/internal/model" // Imported for Swagger annotations
	"gamelink/internal/handler/resp"
	"gamelink/internal/service/playercertification"
	"gamelink/internal/service/playerrank"
	"gamelink/pkg/apierr"
)

// CertificationHandler 陪玩师认证接口（陪玩师端）
type CertificationHandler struct {
	rankSvc *playerrank.PlayerRankService
	certSvc *playercertification.PlayerCertificationService
}

// NewCertificationHandler 创建Handler
func NewCertificationHandler(
	rankSvc *playerrank.PlayerRankService,
	certSvc *playercertification.PlayerCertificationService,
) *CertificationHandler {
	return &CertificationHandler{
		rankSvc: rankSvc,
		certSvc: certSvc,
	}
}

// ApplyRankRequest 申请段位认证请求
type ApplyRankRequest struct {
	GameID         uint64 `json:"gameId" binding:"required"`
	RankID         uint64 `json:"rankId" binding:"required"`
	ScreenshotURLs string `json:"screenshotUrls"` // JSON数组
	Remark         string `json:"remark"`
}

// ApplyCertificationRequest 申请实名认证请求
type ApplyCertificationRequest struct {
	RealName       string `json:"realName" binding:"required,max=64"`
	IDCardNo       string `json:"idCardNo" binding:"required"`
	IDCardFrontURL string `json:"idCardFrontUrl" binding:"required"`
	IDCardBackURL  string `json:"idCardBackUrl" binding:"required"`
	PhotoURL       string `json:"photoUrl"`
	VoiceURL       string `json:"voiceUrl"`
}

// ApplyRank
// @Summary      申请段位认证
// @Tags         Player/Certification
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  ApplyRankRequest  true  "申请信息"
// @Success      201  {object}  model.APIResponse[model.PlayerRankRecord]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /player/certification/rank [post]
func (h *CertificationHandler) ApplyRank(c *gin.Context) {
	playerID, ok := resp.GetPlayerIDOrFail(c)
	if !ok {
		return
	}

	var req ApplyRankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	record, err := h.rankSvc.Apply(c.Request.Context(), playerrank.ApplyInput{
		PlayerID:       playerID,
		GameID:         req.GameID,
		RankID:         req.RankID,
		ScreenshotURLs: req.ScreenshotURLs,
		Remark:         req.Remark,
	})
	if err != nil {
		respondAPIError(c, err)
		return
	}

	resp.Created(c, record)
}

// GetMyRanks
// @Summary      获取我的段位认证列表
// @Tags         Player/Certification
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.PlayerRankRecord]
// @Router       /player/certification/ranks [get]
func (h *CertificationHandler) GetMyRanks(c *gin.Context) {
	playerID, ok := resp.GetPlayerIDOrFail(c)
	if !ok {
		return
	}

	records, err := h.rankSvc.ListByPlayerID(c.Request.Context(), playerID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	resp.OK(c, records)
}

// GetRankDetail
// @Summary      获取段位认证详情
// @Tags         Player/Certification
// @Security     BearerAuth
// @Param        id   path  int  true  "认证记录ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.PlayerRankRecord]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /player/certification/rank/{id} [get]
func (h *CertificationHandler) GetRankDetail(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	record, err := h.rankSvc.Get(c.Request.Context(), id)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	resp.OK(c, record)
}

// ApplyCertification
// @Summary      申请实名认证
// @Tags         Player/Certification
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  ApplyCertificationRequest  true  "申请信息"
// @Success      201  {object}  model.APIResponse[model.PlayerCertification]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /player/certification/identity [post]
func (h *CertificationHandler) ApplyCertification(c *gin.Context) {
	playerID, ok := resp.GetPlayerIDOrFail(c)
	if !ok {
		return
	}

	var req ApplyCertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	cert, err := h.certSvc.Apply(c.Request.Context(), playercertification.ApplyInput{
		PlayerID:       playerID,
		RealName:       req.RealName,
		IDCardNo:       req.IDCardNo,
		IDCardFrontURL: req.IDCardFrontURL,
		IDCardBackURL:  req.IDCardBackURL,
		PhotoURL:       req.PhotoURL,
		VoiceURL:       req.VoiceURL,
	})
	if err != nil {
		respondAPIError(c, err)
		return
	}

	resp.Created(c, cert)
}

// GetMyCertification
// @Summary      获取我的实名认证
// @Tags         Player/Certification
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.PlayerCertification]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /player/certification/identity [get]
func (h *CertificationHandler) GetMyCertification(c *gin.Context) {
	playerID, ok := resp.GetPlayerIDOrFail(c)
	if !ok {
		return
	}

	cert, err := h.certSvc.GetByPlayerID(c.Request.Context(), playerID)
	if err != nil {
		respondAPIError(c, err)
		return
	}

	resp.OK(c, cert)
}

// RegisterCertificationRoutes 注册陪玩师认证路由（陪玩师端）
func RegisterCertificationRoutes(
	router gin.IRouter,
	rankSvc *playerrank.PlayerRankService,
	certSvc *playercertification.PlayerCertificationService,
	authMiddleware gin.HandlerFunc,
) {
	h := NewCertificationHandler(rankSvc, certSvc)

	group := router.Group("/certification")
	group.Use(authMiddleware)
	{
		// 段位认证
		group.POST("/rank", h.ApplyRank)
		group.GET("/ranks", h.GetMyRanks)
		group.GET("/rank/:id", h.GetRankDetail)

		// 实名认证
		group.POST("/identity", h.ApplyCertification)
		group.GET("/identity", h.GetMyCertification)
	}
}
