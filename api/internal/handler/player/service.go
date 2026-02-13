package player

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	playerservice "gamelink/internal/service/player"
	"gamelink/pkg/apierr"
)

type playerServiceDTO struct {
	ID          uint64 `json:"id"`
	PlayerID    uint64 `json:"playerId"`
	GameID      uint64 `json:"gameId"`
	RankID      uint64 `json:"rankId"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"isActive"`
	IsOnline    bool   `json:"isOnline"`
	GameName    string `json:"gameName,omitempty"`
	GameIcon    string `json:"gameIcon,omitempty"`
	RankName    string `json:"rankName,omitempty"`
	PriceCents  int64  `json:"priceCents"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// RegisterServiceRoutes registers player service routes.
func RegisterServiceRoutes(router gin.IRouter, svc *playerservice.ServiceManagement, authMiddleware gin.HandlerFunc) {
	group := router.Group("/services")
	group.Use(authMiddleware)
	group.GET("", func(c *gin.Context) { listServicesHandler(c, svc) })
	group.POST("", func(c *gin.Context) { createServiceHandler(c, svc) })
	group.PUT("/:id", func(c *gin.Context) { updateServiceHandler(c, svc) })
	group.DELETE("/:id", func(c *gin.Context) { deleteServiceHandler(c, svc) })
	group.PUT("/:id/status", func(c *gin.Context) { updateServiceStatusHandler(c, svc) })
}

// listServicesHandler 获取我的服务列表
// @Summary 获取我的服务列表
// @Tags 陪玩师-服务
// @Produce json
// @Security BearerAuth
// @Success 200 {array} playerServiceDTO
// @Router /player/services [get]
func listServicesHandler(c *gin.Context, svc *playerservice.ServiceManagement) {
	userID := resp.GetUserID(c)
	services, err := svc.ListServices(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	items := make([]playerServiceDTO, 0, len(services))
	for i := range services {
		items = append(items, toPlayerServiceDTO(&services[i]))
	}
	resp.OK(c, items)
}

// createServiceHandler 添加服务
// @Summary 添加服务
// @Tags 陪玩师-服务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body playerservice.CreateServiceRequest true "create payload"
// @Success 200 {object} playerServiceDTO
// @Router /player/services [post]
func createServiceHandler(c *gin.Context, svc *playerservice.ServiceManagement) {
	userID := resp.GetUserID(c)
	var req playerservice.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}
	created, err := svc.CreateService(c.Request.Context(), userID, req)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.JSON(c, http.StatusOK, model.APIResponse[playerServiceDTO]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    toPlayerServiceDTO(created),
	})
}

// updateServiceHandler 编辑服务
// @Summary 编辑服务
// @Tags 陪玩师-服务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "service ID"
// @Param body body playerservice.UpdateServiceRequest true "update payload"
// @Success 200 {object} playerServiceDTO
// @Router /player/services/{id} [put]
func updateServiceHandler(c *gin.Context, svc *playerservice.ServiceManagement) {
	userID := resp.GetUserID(c)
	serviceID, err := parseUintParam(c, "id")
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的服务ID"))
		return
	}
	var req playerservice.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}
	updated, err := svc.UpdateService(c.Request.Context(), userID, serviceID, req)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.OK(c, toPlayerServiceDTO(updated))
}

// deleteServiceHandler 删除服务
// @Summary 删除服务
// @Tags 陪玩师-服务
// @Produce json
// @Security BearerAuth
// @Param id path int true "service ID"
// @Success 200 {object} model.SuccessResponse
// @Router /player/services/{id} [delete]
func deleteServiceHandler(c *gin.Context, svc *playerservice.ServiceManagement) {
	userID := resp.GetUserID(c)
	serviceID, err := parseUintParam(c, "id")
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的服务ID"))
		return
	}
	if err := svc.DeleteService(c.Request.Context(), userID, serviceID); err != nil {
		resp.Error(c, err)
		return
	}
	resp.OK(c, gin.H{"deleted": true})
}

// updateServiceStatusHandler 上下架服务
// @Summary 上下架服务
// @Tags 陪玩师-服务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "service ID"
// @Param body body playerservice.UpdateServiceStatusRequest true "status payload"
// @Success 200 {object} playerServiceDTO
// @Router /player/services/{id}/status [put]
func updateServiceStatusHandler(c *gin.Context, svc *playerservice.ServiceManagement) {
	userID := resp.GetUserID(c)
	serviceID, err := parseUintParam(c, "id")
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的服务ID"))
		return
	}
	var req playerservice.UpdateServiceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}
	updated, err := svc.UpdateServiceStatus(c.Request.Context(), userID, serviceID, req.IsActive)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.OK(c, toPlayerServiceDTO(updated))
}

func toPlayerServiceDTO(service *model.PlayerService) playerServiceDTO {
	if service == nil {
		return playerServiceDTO{}
	}

	gameName := ""
	gameIcon := ""
	if service.Game != nil {
		gameName = strings.TrimSpace(service.Game.Name)
		gameIcon = strings.TrimSpace(service.Game.IconURL)
	}

	rankName := ""
	priceCents := int64(0)
	if service.Rank != nil {
		rankName = strings.TrimSpace(service.Rank.Name)
		priceCents = service.Rank.PriceCents
	}

	return playerServiceDTO{
		ID:          service.ID,
		PlayerID:    service.PlayerID,
		GameID:      service.GameID,
		RankID:      service.RankID,
		Description: service.Description,
		IsActive:    service.IsActive,
		IsOnline:    service.IsActive,
		GameName:    gameName,
		GameIcon:    gameIcon,
		RankName:    rankName,
		PriceCents:  priceCents,
		CreatedAt:   service.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   service.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
