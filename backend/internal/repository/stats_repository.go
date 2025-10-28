package repository

import (
    "context"
)

// DateValue 表示按天聚合的数值。
type DateValue struct {
    Date  string
    Value int64
}

// PlayerTop 表示排行榜条目。
type PlayerTop struct {
    PlayerID      uint64
    Nickname      string
    RatingAverage float32
    RatingCount   uint32
}

// Dashboard 返回首页所需汇总数据。
type Dashboard struct {
    TotalUsers          int64
    TotalPlayers        int64
    TotalGames          int64
    TotalOrders         int64
    OrdersByStatus      map[string]int64
    PaymentsByStatus    map[string]int64
    TotalPaidAmountCents int64
}

// StatsRepository 提供统计查询能力。
type StatsRepository interface {
    Dashboard(ctx context.Context) (Dashboard, error)
    RevenueTrend(ctx context.Context, days int) ([]DateValue, error)
    UserGrowth(ctx context.Context, days int) ([]DateValue, error)
    OrdersByStatus(ctx context.Context) (map[string]int64, error)
    TopPlayers(ctx context.Context, limit int) ([]PlayerTop, error)
}

