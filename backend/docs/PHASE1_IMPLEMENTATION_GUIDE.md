# Phase 1 实施指南：抽成机制 & 服务分类

## 🎯 目标

完成GameLink平台的核心商业功能：
1. **抽成机制** - 平台收入来源
2. **服务分类** - 业务差异化

预计时间：**3周**

---

## 📅 Week 1: 抽成机制

### Day 1-2: Repository层

#### 1. 创建 CommissionRepository

```go
// backend/internal/repository/commission_repository.go
package repository

import (
	"context"
	"time"
	"gamelink/internal/model"
	"gorm.io/gorm"
)

type CommissionRepository interface {
	// 抽成规则
	CreateRule(ctx context.Context, rule *model.CommissionRule) error
	GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error)
	GetDefaultRule(ctx context.Context) (*model.CommissionRule, error)
	GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error)
	ListRules(ctx context.Context, opts CommissionRuleListOptions) ([]model.CommissionRule, int64, error)
	UpdateRule(ctx context.Context, rule *model.CommissionRule) error
	DeleteRule(ctx context.Context, id uint64) error
	
	// 抽成记录
	CreateRecord(ctx context.Context, record *model.CommissionRecord) error
	GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error)
	GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error)
	ListRecords(ctx context.Context, opts CommissionRecordListOptions) ([]model.CommissionRecord, int64, error)
	UpdateRecord(ctx context.Context, record *model.CommissionRecord) error
	
	// 月度结算
	CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error
	GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error)
	GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error)
	ListSettlements(ctx context.Context, opts SettlementListOptions) ([]model.MonthlySettlement, int64, error)
	UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error
	
	// 统计查询
	GetMonthlyStats(ctx context.Context, month string) (*MonthlyStats, error)
	GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error)
}

type CommissionRuleListOptions struct {
	Type     *string
	GameID   *uint64
	PlayerID *uint64
	IsActive *bool
	Page     int
	PageSize int
}

type CommissionRecordListOptions struct {
	OrderID          *uint64
	PlayerID         *uint64
	SettlementStatus *string
	SettlementMonth  *string
	DateFrom         *time.Time
	DateTo           *time.Time
	Page             int
	PageSize         int
}

type SettlementListOptions struct {
	PlayerID        *uint64
	SettlementMonth *string
	Status          *string
	Page            int
	PageSize        int
}

type MonthlyStats struct {
	TotalOrders      int64
	TotalIncome      int64
	TotalCommission  int64
	TotalPlayerIncome int64
}
```

**任务清单：**
- [ ] 实现 `commissionRepository` 结构体
- [ ] 实现所有接口方法
- [ ] 编写单元测试

---

### Day 3-4: Service层

#### 2. 创建 CommissionService

```go
// backend/internal/service/commission/commission_service.go
package commission

import (
	"context"
	"fmt"
	"time"
	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type CommissionService struct {
	commissions repository.CommissionRepository
	orders      repository.OrderRepository
	players     repository.PlayerRepository
}

func NewCommissionService(
	commissions repository.CommissionRepository,
	orders repository.OrderRepository,
	players repository.PlayerRepository,
) *CommissionService {
	return &CommissionService{
		commissions: commissions,
		orders:      orders,
		players:     players,
	}
}

// CalculateCommission 计算订单抽成
func (s *CommissionService) CalculateCommission(ctx context.Context, orderID uint64) (*CommissionCalculation, error) {
	// 1. 获取订单
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}
	
	// 2. 查找适用的抽成规则
	rule, err := s.commissions.GetRuleForOrder(ctx, &order.GameID, &order.PlayerID, nil)
	if err != nil || rule == nil {
		// 使用默认规则
		rule, err = s.commissions.GetDefaultRule(ctx)
		if err != nil {
			// 如果没有默认规则，使用20%
			rule = &model.CommissionRule{Rate: 20}
		}
	}
	
	// 3. 计算抽成
	totalAmount := order.PriceCents
	commissionRate := rule.Rate
	commissionAmount := totalAmount * int64(commissionRate) / 100
	playerIncome := totalAmount - commissionAmount
	
	return &CommissionCalculation{
		OrderID:            orderID,
		TotalAmountCents:   totalAmount,
		CommissionRate:     commissionRate,
		CommissionCents:    commissionAmount,
		PlayerIncomeCents:  playerIncome,
	}, nil
}

// RecordCommission 记录订单抽成
func (s *CommissionService) RecordCommission(ctx context.Context, orderID uint64) error {
	// 1. 检查是否已记录
	existing, _ := s.commissions.GetRecordByOrderID(ctx, orderID)
	if existing != nil {
		return fmt.Errorf("commission already recorded for order %d", orderID)
	}
	
	// 2. 计算抽成
	calc, err := s.CalculateCommission(ctx, orderID)
	if err != nil {
		return err
	}
	
	// 3. 获取订单信息
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return err
	}
	
	// 4. 创建抽成记录
	now := time.Now()
	record := &model.CommissionRecord{
		OrderID:            orderID,
		PlayerID:           order.PlayerID,
		TotalAmountCents:   calc.TotalAmountCents,
		CommissionRate:     calc.CommissionRate,
		CommissionCents:    calc.CommissionCents,
		PlayerIncomeCents:  calc.PlayerIncomeCents,
		SettlementStatus:   "pending",
		SettlementMonth:    now.Format("2006-01"),
	}
	
	return s.commissions.CreateRecord(ctx, record)
}

// SettleMonth 月度结算
func (s *CommissionService) SettleMonth(ctx context.Context, month string) error {
	// 1. 获取该月所有待结算记录
	status := "pending"
	records, _, err := s.commissions.ListRecords(ctx, repository.CommissionRecordListOptions{
		SettlementMonth:  &month,
		SettlementStatus: &status,
		Page:             1,
		PageSize:         10000,
	})
	if err != nil {
		return err
	}
	
	// 2. 按陪玩师分组统计
	playerStats := make(map[uint64]*PlayerMonthStats)
	for _, record := range records {
		stats, exists := playerStats[record.PlayerID]
		if !exists {
			stats = &PlayerMonthStats{PlayerID: record.PlayerID}
			playerStats[record.PlayerID] = stats
		}
		stats.OrderCount++
		stats.TotalAmountCents += record.TotalAmountCents
		stats.TotalCommissionCents += record.CommissionCents
		stats.TotalIncomeCents += record.PlayerIncomeCents
	}
	
	// 3. 为每个陪玩师创建月度结算记录
	for _, stats := range playerStats {
		settlement := &model.MonthlySettlement{
			PlayerID:             stats.PlayerID,
			SettlementMonth:      month,
			TotalOrderCount:      stats.OrderCount,
			TotalAmountCents:     stats.TotalAmountCents,
			TotalCommissionCents: stats.TotalCommissionCents,
			TotalIncomeCents:     stats.TotalIncomeCents,
			BonusCents:           0, // 奖金在排名系统中计算
			FinalIncomeCents:     stats.TotalIncomeCents,
			Status:               "pending",
		}
		
		err := s.commissions.CreateSettlement(ctx, settlement)
		if err != nil {
			return err
		}
	}
	
	// 4. 更新抽成记录状态
	for _, record := range records {
		record.SettlementStatus = "settled"
		now := time.Now()
		record.SettledAt = &now
		s.commissions.UpdateRecord(ctx, &record)
	}
	
	return nil
}

type CommissionCalculation struct {
	OrderID            uint64
	TotalAmountCents   int64
	CommissionRate     int
	CommissionCents    int64
	PlayerIncomeCents  int64
}

type PlayerMonthStats struct {
	PlayerID             uint64
	OrderCount           int64
	TotalAmountCents     int64
	TotalCommissionCents int64
	TotalIncomeCents     int64
}
```

**任务清单：**
- [ ] 实现 CommissionService
- [ ] 集成到订单完成流程
- [ ] 编写单元测试

---

### Day 5: 定时任务

#### 3. 月度结算定时任务

```go
// backend/internal/scheduler/settlement_job.go
package scheduler

import (
	"context"
	"log"
	"time"
	"gamelink/internal/service/commission"
	"github.com/robfig/cron/v3"
)

type SettlementScheduler struct {
	commissionSvc *commission.CommissionService
	cron          *cron.Cron
}

func NewSettlementScheduler(commissionSvc *commission.CommissionService) *SettlementScheduler {
	return &SettlementScheduler{
		commissionSvc: commissionSvc,
		cron:          cron.New(),
	}
}

func (s *SettlementScheduler) Start() {
	// 每月1号凌晨2点执行月度结算
	s.cron.AddFunc("0 2 1 * *", s.monthlySettlement)
	s.cron.Start()
	log.Println("Settlement scheduler started")
}

func (s *SettlementScheduler) Stop() {
	s.cron.Stop()
	log.Println("Settlement scheduler stopped")
}

func (s *SettlementScheduler) monthlySettlement() {
	ctx := context.Background()
	
	// 结算上个月
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	month := lastMonth.Format("2006-01")
	
	log.Printf("Starting monthly settlement for %s", month)
	
	err := s.commissionSvc.SettleMonth(ctx, month)
	if err != nil {
		log.Printf("Monthly settlement failed: %v", err)
		return
	}
	
	log.Printf("Monthly settlement completed for %s", month)
}

// 手动触发结算（用于测试）
func (s *SettlementScheduler) TriggerSettlement(month string) error {
	ctx := context.Background()
	return s.commissionSvc.SettleMonth(ctx, month)
}
```

**任务清单：**
- [ ] 实现定时任务
- [ ] 集成到 main.go
- [ ] 测试定时任务

---

## 📅 Week 2: 服务分类系统

### Day 1-2: Repository层

#### 4. 创建 ServiceRepository

```go
// backend/internal/repository/service_repository.go
package repository

import (
	"context"
	"gamelink/internal/model"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	// 服务管理
	Create(ctx context.Context, service *model.Service) error
	Get(ctx context.Context, id uint64) (*model.Service, error)
	List(ctx context.Context, opts ServiceListOptions) ([]model.Service, int64, error)
	Update(ctx context.Context, service *model.Service) error
	Delete(ctx context.Context, id uint64) error
	
	// 批量操作
	BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error
	BatchUpdatePrice(ctx context.Context, ids []uint64, pricePerHour int64) error
}

type ServiceListOptions struct {
	GameID   *uint64
	Type     *model.ServiceType
	IsActive *bool
	Page     int
	PageSize int
}
```

#### 5. 创建 GiftRepository

```go
// backend/internal/repository/gift_repository.go
package repository

import (
	"context"
	"gamelink/internal/model"
)

type GiftRepository interface {
	// 礼物管理
	Create(ctx context.Context, gift *model.Gift) error
	Get(ctx context.Context, id uint64) (*model.Gift, error)
	List(ctx context.Context, opts GiftListOptions) ([]model.Gift, int64, error)
	Update(ctx context.Context, gift *model.Gift) error
	Delete(ctx context.Context, id uint64) error
	
	// 礼物赠送记录
	CreateRecord(ctx context.Context, record *model.GiftRecord) error
	GetRecord(ctx context.Context, id uint64) (*model.GiftRecord, error)
	ListRecords(ctx context.Context, opts GiftRecordListOptions) ([]model.GiftRecord, int64, error)
	
	// 统计
	GetPlayerGiftStats(ctx context.Context, playerID uint64) (*GiftStats, error)
}

type GiftListOptions struct {
	Category *string
	IsActive *bool
	Page     int
	PageSize int
}

type GiftRecordListOptions struct {
	UserID   *uint64
	PlayerID *uint64
	GiftID   *uint64
	OrderID  *uint64
	Page     int
	PageSize int
}

type GiftStats struct {
	TotalReceived int64
	TotalIncome   int64
	TotalCount    int64
}
```

**任务清单：**
- [ ] 实现 ServiceRepository
- [ ] 实现 GiftRepository
- [ ] 编写单元测试

---

### Day 3-4: Service层

#### 6. 创建 ServiceManagementService

```go
// backend/internal/service/servicemanagement/service_management.go
package servicemanagement

import (
	"context"
	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type ServiceManagementService struct {
	services repository.ServiceRepository
	games    repository.GameRepository
}

func NewServiceManagementService(
	services repository.ServiceRepository,
	games repository.GameRepository,
) *ServiceManagementService {
	return &ServiceManagementService{
		services: services,
		games:    games,
	}
}

// CreateService 创建服务
func (s *ServiceManagementService) CreateService(ctx context.Context, req CreateServiceRequest) (*model.Service, error) {
	// 验证游戏存在
	_, err := s.games.Get(ctx, req.GameID)
	if err != nil {
		return nil, err
	}
	
	service := &model.Service{
		GameID:         req.GameID,
		Name:           req.Name,
		Description:    req.Description,
		Type:           req.Type,
		PricePerHour:   req.PricePerHour,
		MinDuration:    req.MinDuration,
		MaxDuration:    req.MaxDuration,
		RequiredRank:   req.RequiredRank,
		CommissionRate: req.CommissionRate,
		IsActive:       true,
		SortOrder:      req.SortOrder,
		Icon:           req.Icon,
		Tags:           req.Tags,
	}
	
	err = s.services.Create(ctx, service)
	if err != nil {
		return nil, err
	}
	
	return service, nil
}

type CreateServiceRequest struct {
	GameID         uint64                `json:"gameId" binding:"required"`
	Name           string                `json:"name" binding:"required,max=128"`
	Description    string                `json:"description"`
	Type           model.ServiceType     `json:"type" binding:"required"`
	PricePerHour   int64                 `json:"pricePerHour" binding:"required,min=1"`
	MinDuration    float32               `json:"minDuration" binding:"required,min=1"`
	MaxDuration    float32               `json:"maxDuration" binding:"required,min=1"`
	RequiredRank   string                `json:"requiredRank"`
	CommissionRate int                   `json:"commissionRate" binding:"required,min=0,max=100"`
	SortOrder      int                   `json:"sortOrder"`
	Icon           string                `json:"icon"`
	Tags           string                `json:"tags"`
}
```

#### 7. 创建 GiftService

```go
// backend/internal/service/gift/gift_service.go
package gift

import (
	"context"
	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type GiftService struct {
	gifts       repository.GiftRepository
	commissions repository.CommissionRepository
	players     repository.PlayerRepository
}

func NewGiftService(
	gifts repository.GiftRepository,
	commissions repository.CommissionRepository,
	players repository.PlayerRepository,
) *GiftService {
	return &GiftService{
		gifts:       gifts,
		commissions: commissions,
		players:     players,
	}
}

// SendGift 赠送礼物
func (s *GiftService) SendGift(ctx context.Context, userID uint64, req SendGiftRequest) (*model.GiftRecord, error) {
	// 1. 验证礼物
	gift, err := s.gifts.Get(ctx, req.GiftID)
	if err != nil {
		return nil, err
	}
	
	// 2. 验证陪玩师
	player, err := s.players.Get(ctx, req.PlayerID)
	if err != nil {
		return nil, err
	}
	
	// 3. 计算金额
	totalPrice := gift.PriceCents * int64(req.Quantity)
	commissionAmount := totalPrice * int64(gift.CommissionRate) / 100
	playerIncome := totalPrice - commissionAmount
	
	// 4. 创建礼物记录
	record := &model.GiftRecord{
		UserID:            userID,
		PlayerID:          req.PlayerID,
		GiftID:            req.GiftID,
		Quantity:          req.Quantity,
		TotalPriceCents:   totalPrice,
		CommissionCents:   commissionAmount,
		PlayerIncomeCents: playerIncome,
		Message:           req.Message,
		IsAnonymous:       req.IsAnonymous,
		OrderID:           req.OrderID,
	}
	
	err = s.gifts.CreateRecord(ctx, record)
	if err != nil {
		return nil, err
	}
	
	// TODO: 发送通知给陪玩师
	
	return record, nil
}

type SendGiftRequest struct {
	PlayerID    uint64  `json:"playerId" binding:"required"`
	GiftID      uint64  `json:"giftId" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,min=1,max=99"`
	Message     string  `json:"message"`
	IsAnonymous bool    `json:"isAnonymous"`
	OrderID     *uint64 `json:"orderId"`
}
```

**任务清单：**
- [ ] 实现 ServiceManagementService
- [ ] 实现 GiftService
- [ ] 编写单元测试

---

### Day 5: Handler层

#### 8. 管理端API

```go
// backend/internal/handler/admin_service.go
func RegisterAdminServiceRoutes(router gin.IRouter, svc *servicemanagement.ServiceManagementService) {
	group := router.Group("/admin/services")
	{
		group.POST("", createServiceHandler(svc))
		group.GET("", listServicesHandler(svc))
		group.GET("/:id", getServiceHandler(svc))
		group.PUT("/:id", updateServiceHandler(svc))
		group.DELETE("/:id", deleteServiceHandler(svc))
		group.POST("/batch-update-status", batchUpdateStatusHandler(svc))
		group.POST("/batch-update-price", batchUpdatePriceHandler(svc))
	}
}
```

#### 9. 用户端API

```go
// backend/internal/handler/user_gift.go
func RegisterUserGiftRoutes(router gin.IRouter, svc *gift.GiftService) {
	group := router.Group("/user/gifts")
	{
		group.GET("", listGiftsHandler(svc))
		group.POST("/send", sendGiftHandler(svc))
		group.GET("/records", getMyGiftRecordsHandler(svc))
	}
}
```

**任务清单：**
- [ ] 实现管理端API
- [ ] 实现用户端API
- [ ] 编写API文档
- [ ] 集成测试

---

## 📅 Week 3: 订单改造 & 集成

### Day 1-2: 订单改造

#### 10. 修改Order模型

```go
// backend/internal/model/order.go
type Order struct {
	// ... 现有字段
	
	// 新增字段
	ServiceID         *uint64     `gorm:"index" json:"serviceId"`
	ServiceType       string      `gorm:"type:varchar(32)" json:"serviceType"`
	CommissionRate    int         `gorm:"default:20" json:"commissionRate"`
	CommissionCents   int64       `gorm:"default:0" json:"commissionCents"`
	PlayerIncomeCents int64       `gorm:"default:0" json:"playerIncomeCents"`
}
```

#### 11. 更新OrderService

```go
// backend/internal/service/order/order_service.go

// CreateOrder 创建订单（新版本）
func (s *OrderService) CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) (*CreateOrderResponse, error) {
	var priceCents int64
	var commissionRate int = 20
	
	// 如果指定了服务，从服务获取价格
	if req.ServiceID != nil {
		service, err := s.services.Get(ctx, *req.ServiceID)
		if err != nil {
			return nil, err
		}
		priceCents = int64(float32(service.PricePerHour) * req.DurationHours)
		commissionRate = service.CommissionRate
	} else {
		// 兼容旧版本，从陪玩师时薪计算
		player, err := s.players.Get(ctx, req.PlayerID)
		if err != nil {
			return nil, err
		}
		priceCents = int64(float32(player.HourlyRateCents) * req.DurationHours)
	}
	
	// 计算抽成
	commissionCents := priceCents * int64(commissionRate) / 100
	playerIncomeCents := priceCents - commissionCents
	
	order := &model.Order{
		// ... 现有字段
		ServiceID:         req.ServiceID,
		ServiceType:       req.ServiceType,
		CommissionRate:    commissionRate,
		CommissionCents:   commissionCents,
		PlayerIncomeCents: playerIncomeCents,
	}
	
	// ...
}
```

**任务清单：**
- [ ] 数据库迁移脚本
- [ ] 更新OrderService
- [ ] 更新相关Handler
- [ ] 数据迁移测试

---

### Day 3: 集成抽成记录

#### 12. 订单完成时自动记录抽成

```go
// backend/internal/service/order/order_service.go

func (s *OrderService) CompleteOrder(ctx context.Context, userID uint64, orderID uint64) error {
	// ... 现有逻辑
	
	// 订单完成后，自动记录抽成
	if err := s.commissionSvc.RecordCommission(ctx, orderID); err != nil {
		log.Printf("Failed to record commission for order %d: %v", orderID, err)
		// 不影响订单完成，异步处理即可
	}
	
	return nil
}
```

**任务清单：**
- [ ] 集成抽成记录
- [ ] 测试完整流程
- [ ] 处理边界情况

---

### Day 4-5: 测试 & 文档

#### 13. 集成测试

```bash
# 测试场景
1. 创建服务 → 下单 → 完成 → 查看抽成记录
2. 赠送礼物 → 查看礼物记录 → 查看陪玩师收入
3. 月度结算 → 查看结算记录
4. 管理员配置抽成规则 → 验证生效
```

#### 14. API文档更新

```markdown
## 新增API

### 管理端
- POST   /admin/services              - 创建服务
- GET    /admin/services              - 服务列表
- PUT    /admin/services/:id          - 更新服务
- DELETE /admin/services/:id          - 删除服务
- POST   /admin/commission-rules      - 创建抽成规则
- GET    /admin/settlements           - 月度结算列表

### 用户端
- GET    /user/services               - 浏览服务
- GET    /user/gifts                  - 礼物列表
- POST   /user/gifts/send             - 赠送礼物
- GET    /user/gifts/records          - 我的礼物记录

### 陪玩师端
- GET    /player/commission/records   - 抽成记录
- GET    /player/commission/summary   - 收入汇总
- GET    /player/gifts/received       - 收到的礼物
```

**任务清单：**
- [ ] 编写集成测试
- [ ] 更新API文档
- [ ] 更新用户手册

---

## ✅ 验收标准

### 功能验收
- [ ] 可以创建和管理服务
- [ ] 可以赠送和查看礼物
- [ ] 订单完成时自动记录抽成
- [ ] 月度自动结算正常运行
- [ ] 管理员可以配置抽成规则
- [ ] 陪玩师可以查看收入明细

### 质量验收
- [ ] 单元测试覆盖率 > 80%
- [ ] 集成测试通过
- [ ] API文档完整
- [ ] 无严重Bug

### 性能验收
- [ ] 月度结算在10分钟内完成
- [ ] API响应时间 < 500ms
- [ ] 数据库查询优化

---

## 🚀 部署清单

### 数据库
- [ ] 执行数据库迁移
- [ ] 创建索引
- [ ] 数据验证

### 应用
- [ ] 更新依赖包
- [ ] 编译部署
- [ ] 启动定时任务

### 配置
- [ ] 配置定时任务
- [ ] 配置默认抽成规则
- [ ] 配置礼物列表

---

## 📞 联系与支持

如有问题，请联系：
- 技术负责人：[姓名]
- 项目经理：[姓名]

---

**文档版本**: 1.0  
**最后更新**: 2025-11-02  
**状态**: ✅ 就绪

