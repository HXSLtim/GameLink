package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

func seedOperationLogs(tx *gorm.DB, users map[string]*model.User, players map[string]*model.Player, orders map[string]*model.Order) error {
	var actorID *uint64
	if admin := users["adminA"]; admin != nil {
		actorID = &admin.ID
	} else if superAdmin := users["superAdminA"]; superAdmin != nil {
		actorID = &superAdmin.ID
	}

	var anyPayment model.Payment
	_ = tx.Order("id ASC").First(&anyPayment).Error

	var anyDispute model.OrderDispute
	_ = tx.Order("id ASC").First(&anyDispute).Error

	var anyReview model.Review
	_ = tx.Order("id ASC").First(&anyReview).Error

	var anyFeed model.Feed
	_ = tx.Order("id ASC").First(&anyFeed).Error

	var anyMessage model.ChatMessage
	_ = tx.Order("id ASC").First(&anyMessage).Error

	getOrderID := func(key string) uint64 {
		if o := orders[key]; o != nil {
			return o.ID
		}
		for _, o := range orders {
			if o != nil {
				return o.ID
			}
		}
		return 0
	}
	getPlayerID := func(key string) uint64 {
		if p := players[key]; p != nil {
			return p.ID
		}
		for _, p := range players {
			if p != nil {
				return p.ID
			}
		}
		return 0
	}
	getUserID := func(key string) uint64 {
		if u := users[key]; u != nil {
			return u.ID
		}
		for _, u := range users {
			if u != nil {
				return u.ID
			}
		}
		return 0
	}

	type logSpec struct {
		EntityType model.OperationEntityType
		Action     model.OperationAction
		EntityID   uint64
		Reason     string
		Metadata   map[string]interface{}
		Offset     time.Duration
	}

	specs := []logSpec{
		{model.OpEntityOrder, model.OpActionCreate, getOrderID("orderPending1"), "后台创建测试订单", map[string]interface{}{"source": "seed", "scenario": "order_create"}, -36 * time.Hour},
		{model.OpEntityOrder, model.OpActionAssignPlayer, getOrderID("orderPending1"), "系统自动派单", map[string]interface{}{"source": "seed", "scenario": "order_assign"}, -35 * time.Hour},
		{model.OpEntityOrder, model.OpActionConfirm, getOrderID("orderConfirmed1"), "客服确认订单", map[string]interface{}{"source": "seed", "scenario": "order_confirm"}, -34 * time.Hour},
		{model.OpEntityOrder, model.OpActionCancel, getOrderID("orderCanceled1"), "用户申请取消订单", map[string]interface{}{"source": "seed", "scenario": "order_cancel"}, -33 * time.Hour},
		{model.OpEntityOrder, model.OpActionComplete, getOrderID("orderCompleted1"), "服务完成，系统结单", map[string]interface{}{"source": "seed", "scenario": "order_complete"}, -30 * time.Hour},

		{model.OpEntityPayment, model.OpActionCapture, anyPayment.ID, "第三方支付成功回调", map[string]interface{}{"source": "seed", "scenario": "payment_capture"}, -29 * time.Hour},
		{model.OpEntityPayment, model.OpActionRefund, anyPayment.ID, "客服执行退款", map[string]interface{}{"source": "seed", "scenario": "payment_refund"}, -28 * time.Hour},

		{model.OpEntityPlayer, model.OpActionUpdateStatus, getPlayerID("playerA"), "陪玩师认证通过", map[string]interface{}{"source": "seed", "status": "verified"}, -27 * time.Hour},
		{model.OpEntityPlayer, model.OpActionUpdateStatus, getPlayerID("playerH"), "陪玩师认证待补充材料", map[string]interface{}{"source": "seed", "status": "pending"}, -26 * time.Hour},

		{model.OpEntityUser, model.OpActionUpdateRole, getUserID("customerA"), "调整用户角色", map[string]interface{}{"source": "seed", "to": "user"}, -25 * time.Hour},
		{model.OpEntityUser, model.OpActionUpdateStatus, getUserID("customerB"), "用户临时封禁", map[string]interface{}{"source": "seed", "status": "banned"}, -24 * time.Hour},

		{model.OpEntityDispute, model.OpActionInitiateDispute, anyDispute.ID, "用户发起纠纷", map[string]interface{}{"source": "seed", "scenario": "dispute_create"}, -23 * time.Hour},
		{model.OpEntityDispute, model.OpActionAssignDispute, anyDispute.ID, "指派客服处理纠纷", map[string]interface{}{"source": "seed", "scenario": "dispute_assign"}, -22 * time.Hour},
		{model.OpEntityDispute, model.OpActionResolveDispute, anyDispute.ID, "纠纷已处理完毕", map[string]interface{}{"source": "seed", "scenario": "dispute_resolve"}, -21 * time.Hour},

		{model.OpEntityReview, model.OpActionApprove, anyReview.ID, "审核通过评价", map[string]interface{}{"source": "seed", "scenario": "review_approve"}, -20 * time.Hour},
		{model.OpEntityReview, model.OpActionHandleReport, anyReview.ID, "处理评价举报", map[string]interface{}{"source": "seed", "scenario": "review_report"}, -19 * time.Hour},

		{model.OpEntityFeed, model.OpActionBatchApprove, anyFeed.ID, "批量审核动态通过", map[string]interface{}{"source": "seed", "scenario": "feed_batch_approve"}, -18 * time.Hour},
		{model.OpEntityChatMessage, model.OpActionDeleteMessage, anyMessage.ID, "删除违规聊天消息", map[string]interface{}{"source": "seed", "scenario": "chat_delete"}, -17 * time.Hour},
		{model.OpEntityChatMessage, model.OpActionMuteUser, anyMessage.ID, "临时禁言用户", map[string]interface{}{"source": "seed", "scenario": "chat_mute"}, -16 * time.Hour},
	}

	for i, spec := range specs {
		metaBytes, err := json.Marshal(spec.Metadata)
		if err != nil {
			return err
		}
		traceID := fmt.Sprintf("seed-op-%03d", i+1)

		var existing model.OperationLog
		err = tx.Where("trace_id = ?", traceID).First(&existing).Error
		if err == nil {
			updates := map[string]interface{}{
				"entity_type":   spec.EntityType,
				"entity_id":     spec.EntityID,
				"actor_user_id": actorID,
				"action":        spec.Action,
				"reason":        spec.Reason,
				"metadata_json": json.RawMessage(metaBytes),
				"created_at":    time.Now().Add(spec.Offset),
			}
			if err := tx.Model(&model.OperationLog{}).Where("id = ?", existing.ID).Updates(updates).Error; err != nil {
				return err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		row := &model.OperationLog{
			EntityType:   string(spec.EntityType),
			EntityID:     spec.EntityID,
			ActorUserID:  actorID,
			Action:       string(spec.Action),
			Reason:       spec.Reason,
			TraceID:      traceID,
			MetadataJSON: json.RawMessage(metaBytes),
		}
		row.CreatedAt = time.Now().Add(spec.Offset)
		row.UpdatedAt = row.CreatedAt
		if err := tx.Create(row).Error; err != nil {
			return err
		}
	}

	return nil
}
