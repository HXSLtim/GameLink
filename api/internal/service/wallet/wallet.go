package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/interfaces"
	walletrepo "gamelink/internal/repository/wallet"
)

var ErrInvalidAmount = errors.New("invalid amount")

// WalletService 提供钱包余额与充值能力
type WalletService struct {
	wallets      walletrepo.Repository
	payments     repository.PaymentRepository
	orders       interfaces.OrderRepository
	serviceItems serviceItemLister
}

type serviceItemLister interface {
	List(ctx context.Context, opts repository.ServiceItemListOptions) ([]model.ServiceItem, int64, error)
}

func NewWalletService(
	wallets walletrepo.Repository,
	payments repository.PaymentRepository,
	orders interfaces.OrderRepository,
	serviceItems ...serviceItemLister,
) *WalletService {
	var lister serviceItemLister
	if len(serviceItems) > 0 {
		lister = serviceItems[0]
	}
	return &WalletService{
		wallets:      wallets,
		payments:     payments,
		orders:       orders,
		serviceItems: lister,
	}
}

type RechargeRequest struct {
	AmountCents int64
	Method      model.PaymentMethod
}

type RechargeResponse struct {
	OrderID   uint64 `json:"orderId"`
	PaymentID uint64 `json:"paymentId"`
	Balance   int64  `json:"balanceCents"`
}

// Recharge 创建充值订单并直接记为已支付，累加余额（简化版，无第三方回调）
func (s *WalletService) Recharge(ctx context.Context, userID uint64, req RechargeRequest) (*RechargeResponse, error) {
	if req.AmountCents <= 0 {
		return nil, ErrInvalidAmount
	}

	// Reuse an existing valid service item as recharge order template to satisfy order FK constraints.
	template, err := s.findRechargeOrderTemplate(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	order := &model.Order{
		UserID:          userID,
		ItemID:          template.ItemID,
		GameID:          template.GameID,
		Status:          model.OrderStatusCompleted,
		Title:           "Wallet Recharge",
		UnitPriceCents:  req.AmountCents,
		TotalPriceCents: req.AmountCents,
		CompletedAt:     &now,
	}
	if err := s.orders.Create(ctx, order); err != nil {
		return nil, err
	}

	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      userID,
		Method:      model.PaymentMethodWallet,
		AmountCents: req.AmountCents,
		Currency:    model.CurrencyCNY,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &now,
	}
	if err := s.payments.Create(ctx, payment); err != nil {
		return nil, err
	}

	w, err := s.wallets.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w = &model.Wallet{UserID: userID}
		} else {
			return nil, err
		}
	}
	w.BalanceCents += req.AmountCents
	if err := s.wallets.Save(ctx, w); err != nil {
		return nil, err
	}

	return &RechargeResponse{
		OrderID:   order.ID,
		PaymentID: payment.ID,
		Balance:   w.BalanceCents,
	}, nil
}

func (s *WalletService) findRechargeOrderTemplate(ctx context.Context, userID uint64) (*model.Order, error) {
	// Prefer user's own order to keep tenant/data affinity.
	opts := interfaces.OrderListOptions{Page: 1, PageSize: 1, UserID: &userID}
	if orders, _, err := s.orders.List(ctx, opts); err == nil && len(orders) > 0 {
		return &orders[0], nil
	}

	// Fallback to any existing order if current user has no history.
	opts = interfaces.OrderListOptions{Page: 1, PageSize: 1}
	orders, _, err := s.orders.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		if s.serviceItems != nil {
			isActive := true
			itemOpts := repository.ServiceItemListOptions{
				Page:     1,
				PageSize: 1,
				IsActive: &isActive,
			}
			if items, _, itemErr := s.serviceItems.List(ctx, itemOpts); itemErr == nil && len(items) > 0 {
				return &model.Order{
					ItemID: items[0].ID,
					GameID: items[0].GameID,
				}, nil
			}
		}
		return nil, fmt.Errorf("no order template available for wallet recharge (missing orders and service items)")
	}
	return &orders[0], nil
}

func (s *WalletService) GetBalance(ctx context.Context, userID uint64) (*model.Wallet, error) {
	w, err := s.wallets.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &model.Wallet{UserID: userID, BalanceCents: 0, FrozenCents: 0}, nil
		}
		return nil, err
	}
	return w, nil
}
