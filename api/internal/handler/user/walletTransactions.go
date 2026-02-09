package user

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// WalletTransactionsHandler 钱包交易记录处理器
type WalletTransactionsHandler struct {
	paymentRepo repository.PaymentRepository
}

// NewWalletTransactionsHandler 创建钱包交易记录处理器
func NewWalletTransactionsHandler(paymentRepo repository.PaymentRepository) *WalletTransactionsHandler {
	return &WalletTransactionsHandler{paymentRepo: paymentRepo}
}

// TransactionItem 交易记录项
type TransactionItem struct {
	ID                uint64 `json:"id"`
	Type              string `json:"type"` // recharge, payment, refund, withdrawal
	AmountCents       int64  `json:"amountCents"`
	Status            string `json:"status"` // pending, success, failed
	Description       string `json:"description"`
	OrderID           uint64 `json:"orderId,omitempty"`
	BalanceAfterCents int64  `json:"balanceAfterCents,omitempty"`
	CreatedAt         string `json:"createdAt"`
}

// TransactionsResponse 交易记录响应
type TransactionsResponse struct {
	Items      []TransactionItem `json:"items"`
	Pagination PaginationInfo    `json:"pagination"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

// GetTransactions 获取交易记录
// @Summary      获取交易记录
// @Description  获取当前用户的钱包交易记录
// @Tags         User - Wallet
// @Security     BearerAuth
// @Produce      json
// @Param        page      query     int     false  "页码"      default(1)
// @Param        pageSize  query     int     false  "每页数量"  default(20)
// @Param        type      query     string  false  "交易类型"  Enums(recharge, payment, refund, withdrawal)
// @Param        status    query     string  false  "状态"      Enums(pending, success, failed)
// @Success      200       {object}  TransactionsResponse
// @Failure      401       {object}  apierr.APIError
// @Failure      500       {object}  apierr.APIError
// @Router       /user/wallet/transactions [get]
func (h *WalletTransactionsHandler) GetTransactions(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	// 解析分页参数
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	// 解析过滤参数
	typeFilter := c.Query("type")
	statusFilter := c.Query("status")

	ctx := c.Request.Context()

	// 构建查询选项
	opts := repository.PaymentListOptions{
		UserID:   &userID,
		Page:     page,
		PageSize: pageSize,
	}

	// 状态过滤
	if statusFilter != "" {
		var status model.PaymentStatus
		switch statusFilter {
		case "pending":
			status = model.PaymentStatusPending
		case "success":
			status = model.PaymentStatusPaid
		case "failed":
			status = model.PaymentStatusFailed
		}
		if status != "" {
			opts.Status = &status
		}
	}

	// 查询支付记录
	payments, total, err := h.paymentRepo.List(ctx, opts)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取交易记录失败").WithDetails(err.Error()))
		return
	}

	// 转换为交易记录
	items := make([]TransactionItem, 0, len(payments))
	for _, p := range payments {
		item := h.paymentToTransaction(&p)

		// 类型过滤
		if typeFilter != "" && item.Type != typeFilter {
			continue
		}

		items = append(items, item)
	}

	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	resp.OK(c, TransactionsResponse{
		Items: items,
		Pagination: PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// paymentToTransaction 将支付记录转换为交易记录
func (h *WalletTransactionsHandler) paymentToTransaction(p *model.Payment) TransactionItem {
	item := TransactionItem{
		ID:        p.ID,
		OrderID:   p.OrderID,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
	}

	// 确定交易类型
	switch {
	case p.Status == model.PaymentStatusRefunded || p.RefundedAmountCents > 0:
		item.Type = "refund"
		item.AmountCents = p.RefundedAmountCents
		item.Description = "订单退款"
	case p.Method == model.PaymentMethodWallet && p.AmountCents > 0:
		// 钱包充值
		item.Type = "recharge"
		item.AmountCents = p.AmountCents
		item.Description = "钱包充值"
	default:
		// 订单支付
		item.Type = "payment"
		item.AmountCents = -p.AmountCents // 支出为负数
		item.Description = "订单支付"
	}

	// 确定状态
	switch p.Status {
	case model.PaymentStatusPending:
		item.Status = "pending"
	case model.PaymentStatusPaid:
		item.Status = "success"
	case model.PaymentStatusFailed:
		item.Status = "failed"
	case model.PaymentStatusRefunded:
		item.Status = "success"
	default:
		item.Status = "pending"
	}

	return item
}

// RegisterWalletTransactionsRoutes 注册钱包交易记录路由
func RegisterWalletTransactionsRoutes(rg *gin.RouterGroup, paymentRepo repository.PaymentRepository, _ gin.HandlerFunc) {
	h := NewWalletTransactionsHandler(paymentRepo)
	walletGroup := rg.Group("/wallet")
	walletGroup.GET("/transactions", h.GetTransactions)
}
