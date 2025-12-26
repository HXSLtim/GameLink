package wallet

import (
	"context"
	"errors"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/interfaces"
	walletrepo "gamelink/internal/repository/wallet"
)

var ErrInvalidAmount = errors.New("invalid amount")

// WalletService 提供钱包余额与充值能力
type WalletService struct {
	wallets  walletrepo.Repository
	payments repository.PaymentRepository
	orders   interfaces.OrderRepository
}

func NewWalletService(wallets walletrepo.Repository, payments repository.PaymentRepository, orders interfaces.OrderRepository) *WalletService {
	return &WalletService{wallets: wallets, payments: payments, orders: orders}
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

	now := time.Now().UTC()
	order := &model.Order{
		UserID:          userID,
		ItemID:          0,
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
		Method:      req.Method,
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
