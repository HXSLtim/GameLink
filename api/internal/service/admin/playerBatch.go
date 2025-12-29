package admin

import (
	"context"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/pkg/apierr"
)

// BatchUpdatePlayerVerificationStatus 批量更新陪玩师认证状态
// 支持状态转换: pending -> verified, pending -> rejected, verified -> rejected, rejected -> pending
func (s *AdminService) BatchUpdatePlayerVerificationStatus(ctx context.Context, playerIDs []uint64, status model.VerificationStatus, verifiedBy *uint64, remark string) (*BatchOperationResponse, error) {
	response := &BatchOperationResponse{
		TotalCount:   len(playerIDs),
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchErrorItem, 0),
	}

	remark = normalizeNote(remark)

	// Validate status
	if status != model.VerificationPending && status != model.VerificationVerified && status != model.VerificationRejected {
		response.FailedCount = len(playerIDs)
		for _, playerID := range playerIDs {
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      playerID,
				Message: "invalid verification status",
			})
		}
		return response, apierr.BadRequest("invalid verification status")
	}

	for _, playerID := range playerIDs {
		// Get player
		player, err := s.players.Get(ctx, playerID)
		if err != nil {
			if apierr.IsNotFound(err) {
				response.FailedItems = append(response.FailedItems, BatchErrorItem{
					ID:      playerID,
					Message: "player not found",
				})
				response.FailedCount++
				continue
			}
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      playerID,
				Message: fmt.Sprintf("get player failed: %v", err),
			})
			response.FailedCount++
			continue
		}

		// Validate status transition
		if !isValidVerificationStatusTransition(player.VerificationStatus, status) {
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      playerID,
				Message: fmt.Sprintf("cannot transition from %s to %s", player.VerificationStatus, status),
			})
			response.FailedCount++
			continue
		}

		// Update player verification status
		player.VerificationStatus = status

		// Set verification metadata
		now := time.Now()
		if status == model.VerificationVerified {
			player.VerifiedAt = &now
			player.VerifiedBy = verifiedBy
			player.VerifyRemark = remark
		} else if status == model.VerificationRejected {
			player.RejectReason = remark
		}

		// Clear verified fields if moving back to pending
		if status == model.VerificationPending {
			player.VerifiedAt = nil
			player.VerifiedBy = nil
			player.VerifyRemark = ""
			player.RejectReason = ""
		}

		err = s.players.Update(ctx, player)
		if err != nil {
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      playerID,
				Message: fmt.Sprintf("update failed: %v", err),
			})
			response.FailedCount++
			continue
		}

		response.SuccessItems = append(response.SuccessItems, playerID)
		response.SuccessCount++

		// Log operation
		s.appendLogAsync(ctx, string(model.OpEntityPlayer), playerID, string(model.OpActionUpdate), map[string]any{
			"verification_status":   status,
			"verified_by":           verifiedBy,
			"previous_status":       player.VerificationStatus,
			"remark":                remark,
		})
	}

	if response.SuccessCount > 0 {
		s.invalidateCache(ctx, cacheKeyPlayers)
	}

	return response, nil
}

// BatchRevokePlayerCertification 批量撤销陪玩师认证
// 将认证状态从 verified 改为 pending，清除认证信息
func (s *AdminService) BatchRevokePlayerCertification(ctx context.Context, playerIDs []uint64, reason string, revokedBy *uint64) (*BatchOperationResponse, error) {
	response := &BatchOperationResponse{
		TotalCount:   len(playerIDs),
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchErrorItem, 0),
	}

	reason = normalizeNote(reason)

	for _, playerID := range playerIDs {
		// Get player
		player, err := s.players.Get(ctx, playerID)
		if err != nil {
			if apierr.IsNotFound(err) {
				response.FailedItems = append(response.FailedItems, BatchErrorItem{
					ID:      playerID,
					Message: "player not found",
				})
				response.FailedCount++
				continue
			}
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      playerID,
				Message: fmt.Sprintf("get player failed: %v", err),
			})
			response.FailedCount++
			continue
		}

		// Check if player is verified (only verified players can be revoked)
		if player.VerificationStatus != model.VerificationVerified {
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      playerID,
				Message: fmt.Sprintf("cannot revoke player with status: %s", player.VerificationStatus),
			})
			response.FailedCount++
			continue
		}

		// Store previous status for logging
		previousStatus := player.VerificationStatus

		// Revoke certification - set to pending and clear verification fields
		player.VerificationStatus = model.VerificationPending
		player.VerifiedAt = nil
		player.VerifiedBy = nil
		player.VerifyRemark = reason
		player.RejectReason = fmt.Sprintf("认证已撤销: %s", reason)

		err = s.players.Update(ctx, player)
		if err != nil {
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      playerID,
				Message: fmt.Sprintf("revoke failed: %v", err),
			})
			response.FailedCount++
			continue
		}

		response.SuccessItems = append(response.SuccessItems, playerID)
		response.SuccessCount++

		// Log operation
		s.appendLogAsync(ctx, string(model.OpEntityPlayer), playerID, "revoke_certification", map[string]any{
			"previous_status":     previousStatus,
			"new_status":          player.VerificationStatus,
			"revoked_by":          revokedBy,
			"revoke_reason":       reason,
		})
	}

	if response.SuccessCount > 0 {
		s.invalidateCache(ctx, cacheKeyPlayers)
	}

	return response, nil
}

// isValidVerificationStatusTransition 检查认证状态转换是否合法
func isValidVerificationStatusTransition(current, new model.VerificationStatus) bool {
	// Define valid transitions
	validTransitions := map[model.VerificationStatus][]model.VerificationStatus{
		model.VerificationPending: {
			model.VerificationVerified,
			model.VerificationRejected,
		},
		model.VerificationVerified: {
			model.VerificationRejected,
			model.VerificationPending, // Allow revoking
		},
		model.VerificationRejected: {
			model.VerificationPending,
			model.VerificationVerified, // Allow re-verification
		},
	}

	allowedStatuses, ok := validTransitions[current]
	if !ok {
		return false
	}

	for _, allowed := range allowedStatuses {
		if allowed == new {
			return true
		}
	}

	return false
}
