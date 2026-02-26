package admin

import (
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/pkg/apierr"
)

// BatchUpdateVerificationStatusRequest 批量更新陪玩师认证状态请求
type BatchUpdateVerificationStatusRequest struct {
	PlayerIDs []uint64 `json:"playerIds" binding:"required,min=1,max=100"`
	Status    string   `json:"status" binding:"required,oneof=pending verified rejected"`
	Remark    string   `json:"remark"`
}

// BatchRevokeCertificationRequest 批量撤销陪玩师认证请求
type BatchRevokeCertificationRequest struct {
	PlayerIDs []uint64 `json:"playerIds" binding:"required,min=1,max=100"`
	Reason    string   `json:"reason"`
}

// BatchUpdateVerificationStatus 批量更新陪玩师认证状态
func (h *PlayerHandler) BatchUpdateVerificationStatus(c *gin.Context) {
	var req BatchUpdateVerificationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	var verifiedBy *uint64
	if actorID := getUserIDFromContext(c); actorID != 0 {
		verifiedBy = &actorID
	}

	result, err := h.svc.BatchUpdatePlayerVerificationStatus(
		contextWithActor(c),
		req.PlayerIDs,
		model.VerificationStatus(strings.ToLower(strings.TrimSpace(req.Status))),
		verifiedBy,
		req.Remark,
	)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("batch update player verification failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchRevokeCertification 批量撤销陪玩师认证
func (h *PlayerHandler) BatchRevokeCertification(c *gin.Context) {
	var req BatchRevokeCertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	var revokedBy *uint64
	if actorID := getUserIDFromContext(c); actorID != 0 {
		revokedBy = &actorID
	}

	result, err := h.svc.BatchRevokePlayerCertification(contextWithActor(c), req.PlayerIDs, req.Reason, revokedBy)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("batch revoke player certification failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}
