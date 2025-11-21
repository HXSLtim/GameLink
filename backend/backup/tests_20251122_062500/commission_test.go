package commission

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"gamelink/internal/model"
	commissionrepo "gamelink/internal/repository/commission"
	rankingrepo "gamelink/internal/repository/ranking"
	repoiface "gamelink/internal/repository/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type mockCommissionRepository struct {
	mock.Mock
}

func (m *mockCommissionRepository) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRecord), args.Error(1)
}

func (m *mockCommissionRepository) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *mockCommissionRepository) ListRecords(ctx context.Context, opts commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.CommissionRecord), args.Get(1).(int64), args.Error(2)
}

func (m *mockCommissionRepository) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *mockCommissionRepository) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	args := m.Called(ctx, playerID, month)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockCommissionRepository) ListSettlements(ctx context.Context, opts commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.MonthlySettlement), args.Get(1).(int64), args.Error(2)
}

func (m *mockCommissionRepository) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	args := m.Called(ctx, settlement)
	return args.Error(0)
}

func (m *mockCommissionRepository) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *mockCommissionRepository) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	args := m.Called(ctx, gameID, playerID, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *mockCommissionRepository) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *mockCommissionRepository) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *mockCommissionRepository) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *mockCommissionRepository) GetMonthlyStats(ctx context.Context, month string) (*commissionrepo.MonthlyStats, error) {
	args := m.Called(ctx, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commissionrepo.MonthlyStats), args.Error(1)
}

type mockOrderReader struct {
	mock.Mock
}

func (m *mockOrderReader) Get(ctx context.Context, id uint64) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

type mockPlayerRepository struct {
	mock.Mock
}

func (m *mockPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

type mockServiceItemRepository struct {
	mock.Mock
}

func (m *mockServiceItemRepository) Get(ctx context.Context, id uint64) (*model.ServiceItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ServiceItem), args.Error(1)
}

type mockRankingRepository struct {
	mock.Mock
}

func (m *mockRankingRepository) GetPlayerRanking(ctx context.Context, playerID uint64, rankingType model.RankingType, period string, month string) (*model.PlayerRanking, error) {
	args := m.Called(ctx, playerID, rankingType, period, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PlayerRanking), args.Error(1)
}

type mockRankingCommissionRepository struct {
	mock.Mock
}

func (m *mockRankingCommissionRepository) GetActiveConfigForMonth(ctx context.Context, rankingType model.RankingType, month string) (*model.RankingCommissionConfig, error) {
	args := m.Called(ctx, rankingType, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RankingCommissionConfig), args.Error(1)
}

// TestGetRankingCommissionRate tests the ranking commission rate functionality
func TestGetRankingCommissionRate(t *testing.T) {
	t.Run("successfully gets ranking commission rate", func(t *testing.T) {
		// Setup mocks
		commissionRepo := new(mockCommissionRepository)
		orderReader := new(mockOrderReader)
		playerRepo := new(mockPlayerRepository)
		itemRepo := new(mockServiceItemRepository)
		rankingRepo := new(mockRankingRepository)
		rankingCommissionRepo := new(mockRankingCommissionRepository)

		service := NewCommissionService(commissionRepo, orderReader, playerRepo, itemRepo, rankingRepo, rankingCommissionRepo)

		// Mock data
		playerID := uint64(1)
		lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")
		
		ranking := &model.PlayerRanking{
			PlayerID: playerID,
			Rank:     3,
			Score:    10000,
		}
		
		rulesJSON := `[
			{"rankStart": 1, "rankEnd": 1, "commissionRate": 10},
			{"rankStart": 2, "rankEnd": 3, "commissionRate": 15},
			{"rankStart": 4, "rankEnd": 10, "commissionRate": 18}
		]`
		
		config := &model.RankingCommissionConfig{
			ID:        1,
			Name:      "月度收入排名抽成",
			RulesJSON: rulesJSON,
		}

		// Setup expectations
		rankingRepo.On("GetPlayerRanking", mock.Anything, playerID, model.RankingTypeIncome, "monthly", lastMonth).
			Return(ranking, nil)
		rankingCommissionRepo.On("GetActiveConfigForMonth", mock.Anything, model.RankingTypeIncome, lastMonth).
			Return(config, nil)

		// Execute
		rate, detail := service.getRankingCommissionRate(context.Background(), playerID)

		// Assert
		assert.Equal(t, 15, rate)
		assert.Equal(t, "月度收入排名抽成第3名", detail)
		rankingRepo.AssertExpectations(t)
		rankingCommissionRepo.AssertExpectations(t)
	})

	t.Run("returns 0 when ranking not found", func(t *testing.T) {
		// Setup mocks
		commissionRepo := new(mockCommissionRepository)
		orderReader := new(mockOrderReader)
		playerRepo := new(mockPlayerRepository)
		itemRepo := new(mockServiceItemRepository)
		rankingRepo := new(mockRankingRepository)
		rankingCommissionRepo := new(mockRankingCommissionRepository)

		service := NewCommissionService(commissionRepo, orderReader, playerRepo, itemRepo, rankingRepo, rankingCommissionRepo)

		playerID := uint64(1)
		lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")

		// Setup expectations
		rankingRepo.On("GetPlayerRanking", mock.Anything, playerID, model.RankingTypeIncome, "monthly", lastMonth).
			Return(nil, errors.New("ranking not found"))

		// Execute
		rate, detail := service.getRankingCommissionRate(context.Background(), playerID)

		// Assert
		assert.Equal(t, 0, rate)
		assert.Equal(t, "", detail)
		rankingRepo.AssertExpectations(t)
	})

	t.Run("returns 0 when config not found", func(t *testing.T) {
		// Setup mocks
		commissionRepo := new(mockCommissionRepository)
		orderReader := new(mockOrderReader)
		playerRepo := new(mockPlayerRepository)
		itemRepo := new(mockServiceItemRepository)
		rankingRepo := new(mockRankingRepository)
		rankingCommissionRepo := new(mockRankingCommissionRepository)

		service := NewCommissionService(commissionRepo, orderReader, playerRepo, itemRepo, rankingRepo, rankingCommissionRepo)

		playerID := uint64(1)
		lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")
		
		ranking := &model.PlayerRanking{
			PlayerID: playerID,
			Rank:     3,
			Score:    10000,
		}

		// Setup expectations
		rankingRepo.On("GetPlayerRanking", mock.Anything, playerID, model.RankingTypeIncome, "monthly", lastMonth).
			Return(ranking, nil)
		rankingCommissionRepo.On("GetActiveConfigForMonth", mock.Anything, model.RankingTypeIncome, lastMonth).
			Return(nil, errors.New("config not found"))

		// Execute
		rate, detail := service.getRankingCommissionRate(context.Background(), playerID)

		// Assert
		assert.Equal(t, 0, rate)
		assert.Equal(t, "", detail)
		rankingRepo.AssertExpectations(t)
		rankingCommissionRepo.AssertExpectations(t)
	})

	t.Run("returns 0 when rank is not in any rule range", func(t *testing.T) {
		// Setup mocks
		commissionRepo := new(mockCommissionRepository)
		orderReader := new(mockOrderReader)
		playerRepo := new(mockPlayerRepository)
		itemRepo := new(mockServiceItemRepository)
		rankingRepo := new(mockRankingRepository)
		rankingCommissionRepo := new(mockRankingCommissionRepository)

		service := NewCommissionService(commissionRepo, orderReader, playerRepo, itemRepo, rankingRepo, rankingCommissionRepo)

		playerID := uint64(1)
		lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")
		
		ranking := &model.PlayerRanking{
			PlayerID: playerID,
			Rank:     20, // 不在规则范围内
			Score:    10000,
		}
		
		rulesJSON := `[
			{"rankStart": 1, "rankEnd": 1, "commissionRate": 10},
			{"rankStart": 2, "rankEnd": 3, "commissionRate": 15},
			{"rankStart": 4, "rankEnd": 10, "commissionRate": 18}
		]`
		
		config := &model.RankingCommissionConfig{
			ID:        1,
			Name:      "月度收入排名抽成",
			RulesJSON: rulesJSON,
		}

		// Setup expectations
		rankingRepo.On("GetPlayerRanking", mock.Anything, playerID, model.RankingTypeIncome, "monthly", lastMonth).
			Return(ranking, nil)
		rankingCommissionRepo.On("GetActiveConfigForMonth", mock.Anything, model.RankingTypeIncome, lastMonth).
			Return(config, nil)

		// Execute
		rate, detail := service.getRankingCommissionRate(context.Background(), playerID)

		// Assert
		assert.Equal(t, 0, rate)
		assert.Equal(t, "", detail)
		rankingRepo.AssertExpectations(t)
		rankingCommissionRepo.AssertExpectations(t)
	})
}

// TestCalculateOrderCommissionWithRanking tests order commission calculation with ranking
func TestCalculateOrderCommissionWithRanking(t *testing.T) {
	t.Run("applies ranking commission for non-gift orders", func(t *testing.T) {
		// Setup mocks
		commissionRepo := new(mockCommissionRepository)
		orderReader := new(mockOrderReader)
		playerRepo := new(mockPlayerRepository)
		itemRepo := new(mockServiceItemRepository)
		rankingRepo := new(mockRankingRepository)
		rankingCommissionRepo := new(mockRankingCommissionRepository)

		service := NewCommissionService(commissionRepo, orderReader, playerRepo, itemRepo, rankingRepo, rankingCommissionRepo)

		// Mock data
		order := &model.Order{
			ID:              1,
			ItemID:          1,
			TotalPriceCents: 10000,
			PlayerID:        func() *uint64 { id := uint64(1); return &id }(),
			GameID:          func() *uint64 { id := uint64(1); return &id }(),
		}

		serviceItem := &model.ServiceItem{
			ID:             1,
			Name:           "游戏陪玩",
			BasePriceCents: 10000,
			CommissionRate: 0.2, // 20%
		}

		playerRule := &model.CommissionRule{
			ID:          1,
			Name:        "陪玩师专属优惠",
			Rate:        25, // 25%
			Type:        "special",
			GameID:      func() *uint64 { id := uint64(1); return &id }(),
			PlayerID:    func() *uint64 { id := uint64(1); return &id }(),
			ServiceType: func() *string { s := "game"; return &s }(),
			IsActive:    true,
		}

		lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")
		
		ranking := &model.PlayerRanking{
			PlayerID: 1,
			Rank:     2,
			Score:    50000,
		}
		
		rulesJSON := `[
			{"rankStart": 1, "rankEnd": 1, "commissionRate": 10},
			{"rankStart": 2, "rankEnd": 3, "commissionRate": 15},
			{"rankStart": 4, "rankEnd": 10, "commissionRate": 18}
		]`
		
		config := &model.RankingCommissionConfig{
			ID:        1,
			Name:      "月度收入排名抽成",
			RulesJSON: rulesJSON,
		}

		// Setup expectations
		itemRepo.On("Get", mock.Anything, uint64(1)).Return(serviceItem, nil)
		commissionRepo.On("GetRuleForOrder", mock.Anything, order.GameID, order.PlayerID, mock.Anything).
			Return(playerRule, nil)
		rankingRepo.On("GetPlayerRanking", mock.Anything, uint64(1), model.RankingTypeIncome, "monthly", lastMonth).
			Return(ranking, nil)
		rankingCommissionRepo.On("GetActiveConfigForMonth", mock.Anything, model.RankingTypeIncome, lastMonth).
			Return(config, nil)

		// Execute
		result, err := service.CalculateOrderCommission(context.Background(), order)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 15, result.CommissionRate) // 应该应用排名优惠15%
		assert.Equal(t, "排名优惠", result.AppliedRule)
		assert.Equal(t, 3, len(result.CandidateRates))
		
		itemRepo.AssertExpectations(t)
		commissionRepo.AssertExpectations(t)
		rankingRepo.AssertExpectations(t)
		rankingCommissionRepo.AssertExpectations(t)
	})

	t.Run("does not apply ranking commission for gift orders", func(t *testing.T) {
		// Setup mocks
		commissionRepo := new(mockCommissionRepository)
		orderReader := new(mockOrderReader)
		playerRepo := new(mockPlayerRepository)
		itemRepo := new(mockServiceItemRepository)
		rankingRepo := new(mockRankingRepository)
		rankingCommissionRepo := new(mockRankingCommissionRepository)

		service := NewCommissionService(commissionRepo, orderReader, playerRepo, itemRepo, rankingRepo, rankingCommissionRepo)

		// Mock data - gift order
		order := &model.Order{
			ID:              1,
			ItemID:          1,
			TotalPriceCents: 10000,
			PlayerID:        func() *uint64 { id := uint64(1); return &id }(),
			GameID:          func() *uint64 { id := uint64(1); return &id }(),
			RecipientPlayerID: func() *uint64 { id := uint64(1); return &id }(), // This makes it a gift order
		}

		serviceItem := &model.ServiceItem{
			ID:             1,
			Name:           "礼物",
			BasePriceCents: 10000,
			CommissionRate: 0.2, // 20%
			Type:           "gift",
		}

		// Setup expectations
		itemRepo.On("Get", mock.Anything, uint64(1)).Return(serviceItem, nil)
		commissionRepo.On("GetRuleForOrder", mock.Anything, order.GameID, order.PlayerID, mock.Anything).
			Return(nil, errors.New("no rule found"))
		commissionRepo.On("GetDefaultRule", mock.Anything).
			Return(&model.CommissionRule{
				ID:   1,
				Name: "默认规则",
				Rate: 20,
				Type: "default",
			}, nil)

		// Execute
		result, err := service.CalculateOrderCommission(context.Background(), order)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 20, result.CommissionRate) // 应该应用默认20%
		assert.Equal(t, "默认规则", result.AppliedRule)
		assert.Equal(t, 2, len(result.CandidateRates)) // 只有服务项目和默认规则，没有排名优惠
		
		// Verify ranking methods were NOT called
		rankingRepo.AssertNotCalled(t, "GetPlayerRanking", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		rankingCommissionRepo.AssertNotCalled(t, "GetActiveConfigForMonth", mock.Anything, mock.Anything, mock.Anything)
		
		itemRepo.AssertExpectations(t)
		commissionRepo.AssertExpectations(t)
	})

	t.Run("does not apply ranking commission when no player", func(t *testing.T) {
		// Setup mocks
		commissionRepo := new(mockCommissionRepository)
		orderReader := new(mockOrderReader)
		playerRepo := new(mockPlayerRepository)
		itemRepo := new(mockServiceItemRepository)
		rankingRepo := new(mockRankingRepository)
		rankingCommissionRepo := new(mockRankingCommissionRepository)

		service := NewCommissionService(commissionRepo, orderReader, playerRepo, itemRepo, rankingRepo, rankingCommissionRepo)

		// Mock data - order without player
		order := &model.Order{
			ID:              1,
			ItemID:          1,
			TotalPriceCents: 10000,
			PlayerID:        nil, // No player
			GameID:          func() *uint64 { id := uint64(1); return &id }(),
		}

		serviceItem := &model.ServiceItem{
			ID:             1,
			Name:           "游戏陪玩",
			BasePriceCents: 10000,
			CommissionRate: 0.2, // 20%
		}

		// Setup expectations
		itemRepo.On("Get", mock.Anything, uint64(1)).Return(serviceItem, nil)
		commissionRepo.On("GetDefaultRule", mock.Anything).
			Return(&model.CommissionRule{
				ID:   1,
				Name: "默认规则",
				Rate: 20,
				Type: "default",
			}, nil)

		// Execute
		result, err := service.CalculateOrderCommission(context.Background(), order)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 20, result.CommissionRate) // 应该应用默认20%
		assert.Equal(t, "默认规则", result.AppliedRule)
		assert.Equal(t, 2, len(result.CandidateRates)) // 只有服务项目和默认规则，没有排名优惠
		
		// Verify ranking methods were NOT called
		rankingRepo.AssertNotCalled(t, "GetPlayerRanking", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		rankingCommissionRepo.AssertNotCalled(t, "GetActiveConfigForMonth", mock.Anything, mock.Anything, mock.Anything)
		
		itemRepo.AssertExpectations(t)
		commissionRepo.AssertExpectations(t)
	})
}

// TestParseRankingCommissionRules tests JSON rule parsing
func TestParseRankingCommissionRules(t *testing.T) {
	t.Run("successfully parses valid JSON rules", func(t *testing.T) {
		rulesJSON := `[
			{"rankStart": 1, "rankEnd": 1, "commissionRate": 10},
			{"rankStart": 2, "rankEnd": 3, "commissionRate": 15},
			{"rankStart": 4, "rankEnd": 10, "commissionRate": 18}
		]`

		rules, err := ParseRankingCommissionRules(rulesJSON)

		assert.NoError(t, err)
		assert.Equal(t, 3, len(rules))
		assert.Equal(t, 1, rules[0].RankStart)
		assert.Equal(t, 1, rules[0].RankEnd)
		assert.Equal(t, 10, rules[0].CommissionRate)
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		rulesJSON := `invalid json`

		_, err := ParseRankingCommissionRules(rulesJSON)

		assert.Error(t, err)
	})
}

// TestFindCommissionRateForRank tests finding commission rate by rank
func TestFindCommissionRateForRank(t *testing.T) {
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 1, CommissionRate: 10},
		{RankStart: 2, RankEnd: 3, CommissionRate: 15},
		{RankStart: 4, RankEnd: 10, CommissionRate: 18},
	}

	tests := []struct {
		rank     int
		expected int
	}{
		{1, 10},
		{2, 15},
		{3, 15},
		{5, 18},
		{10, 18},
		{11, 0}, // Not in any range
		{0, 0},  // Invalid rank
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("rank %d", tt.rank), func(t *testing.T) {
			rate := FindCommissionRateForRank(rules, tt.rank)
			assert.Equal(t, tt.expected, rate)
		})
	}
}

// TestValidateRankingRules tests rule validation
func TestValidateRankingRules(t *testing.T) {
	t.Run("validates correct rules", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 1, CommissionRate: 10},
			{RankStart: 2, RankEnd: 3, CommissionRate: 15},
			{RankStart: 4, RankEnd: 10, CommissionRate: 18},
		}

		err := ValidateRankingRules(rules)
		assert.NoError(t, err)
	})

	t.Run("returns error for overlapping ranges", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 5, CommissionRate: 10},
			{RankStart: 3, RankEnd: 7, CommissionRate: 15}, // Overlaps with first rule
		}

		err := ValidateRankingRules(rules)
		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})

	t.Run("returns error for invalid rank range", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 5, RankEnd: 1, CommissionRate: 10}, // Invalid range
		}

		err := ValidateRankingRules(rules)
		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})

	t.Run("returns error for invalid commission rate", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 1, CommissionRate: 150}, // Invalid rate (>100)
		}

		err := ValidateRankingRules(rules)
		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})

	t.Run("returns error for negative commission rate", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 1, CommissionRate: -5}, // Invalid rate (<0)
		}

		err := ValidateRankingRules(rules)
		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})
}

// TestSelectLowestRate tests selecting the lowest commission rate
func TestSelectLowestRate(t *testing.T) {
	t.Run("selects lowest rate from candidates", func(t *testing.T) {
		candidates := []CommissionCandidate{
			{Source: "服务项目", Rate: 20, Detail: "游戏陪玩"},
			{Source: "陪玩师专属", Rate: 25, Detail: "VIP优惠"},
			{Source: "排名优惠", Rate: 15, Detail: "第2-3名"},
		}

		lowest := selectLowestRate(candidates)

		assert.Equal(t, 15, lowest.Rate)
		assert.Equal(t, "排名优惠", lowest.Source)
	})

	t.Run("returns default when no candidates", func(t *testing.T) {
		candidates := []CommissionCandidate{}

		lowest := selectLowestRate(candidates)

		assert.Equal(t, 20, lowest.Rate)
		assert.Equal(t, "默认规则", lowest.Source)
		assert.Equal(t, "平台默认20%抽成", lowest.Detail)
	})
}