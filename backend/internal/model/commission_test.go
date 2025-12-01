package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"gamelink/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCommissionRuleModel(t *testing.T) {
	now := time.Now()
	gameID := uint64(10)
	serviceType := "escort"

	commissionRule := &model.CommissionRule{
		ID:          1,
		Name:        "默认抽成规则",
		Description: "这是默认的抽成规则",
		Type:        "default",
		Rate:        20, // 20%
		IsActive:    true,
		GameID:      &gameID,
		PlayerID:    nil,
		ServiceType: &serviceType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, uint64(1), commissionRule.ID)
	assert.Equal(t, "默认抽成规则", commissionRule.Name)
	assert.Equal(t, "这是默认的抽成规则", commissionRule.Description)
	assert.Equal(t, "default", commissionRule.Type)
	assert.Equal(t, 20, commissionRule.Rate)
	assert.True(t, commissionRule.IsActive)
	assert.Equal(t, &gameID, commissionRule.GameID)
	assert.Nil(t, commissionRule.PlayerID)
	assert.Equal(t, &serviceType, commissionRule.ServiceType)
	assert.Equal(t, now, commissionRule.CreatedAt)
	assert.Equal(t, now, commissionRule.UpdatedAt)
}

func TestCommissionRuleJSONSerialization(t *testing.T) {
	now := time.Now()
	gameID := uint64(10)
	serviceType := "escort"

	commissionRule := &model.CommissionRule{
		ID:          1,
		Name:        "默认抽成规则",
		Description: "这是默认的抽成规则",
		Type:        "default",
		Rate:        20,
		IsActive:    true,
		GameID:      &gameID,
		PlayerID:    nil,
		ServiceType: &serviceType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 序列化
	data, err := json.Marshal(commissionRule)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "默认抽成规则")
	assert.Contains(t, string(data), "default")

	// 反序列化
	var decoded model.CommissionRule
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, commissionRule.ID, decoded.ID)
	assert.Equal(t, commissionRule.Name, decoded.Name)
	assert.Equal(t, commissionRule.Type, decoded.Type)
	assert.Equal(t, commissionRule.Rate, decoded.Rate)
	assert.Equal(t, commissionRule.IsActive, decoded.IsActive)
}

func TestCommissionRuleTableName(t *testing.T) {
	commissionRule := model.CommissionRule{}
	tableName := commissionRule.TableName()
	assert.Equal(t, "commission_rules", tableName)
}

func TestCommissionRuleZeroValues(t *testing.T) {
	commissionRule := &model.CommissionRule{
		ID:          0,
		Name:        "",
		Description: "",
		Type:        "",
		Rate:        0,
		IsActive:    false,
		GameID:      nil,
		PlayerID:    nil,
		ServiceType: nil,
	}

	assert.Equal(t, uint64(0), commissionRule.ID)
	assert.Equal(t, "", commissionRule.Name)
	assert.Equal(t, "", commissionRule.Description)
	assert.Equal(t, "", commissionRule.Type)
	assert.Equal(t, 0, commissionRule.Rate)
	assert.False(t, commissionRule.IsActive)
	assert.Nil(t, commissionRule.GameID)
	assert.Nil(t, commissionRule.PlayerID)
	assert.Nil(t, commissionRule.ServiceType)
}

func TestCommissionRuleEdgeCases(t *testing.T) {
	// 测试高抽成率
	commissionRule1 := &model.CommissionRule{
		Rate: 100, // 100%
	}
	assert.Equal(t, 100, commissionRule1.Rate)

	// 测试0%抽成率
	commissionRule2 := &model.CommissionRule{
		Rate: 0,
	}
	assert.Equal(t, 0, commissionRule2.Rate)

	// 测试负抽成率（虽然业务上不合理）
	commissionRule3 := &model.CommissionRule{
		Rate: -10,
	}
	assert.Equal(t, -10, commissionRule3.Rate)

	// 测试大数值
	commissionRule4 := &model.CommissionRule{
		Rate: 999,
	}
	assert.Equal(t, 999, commissionRule4.Rate)
}

func TestCommissionRuleWithAllFields(t *testing.T) {
	now := time.Now()
	gameID := uint64(5)
	playerID := uint64(15)
	serviceType := "gift"

	commissionRule := &model.CommissionRule{
		ID:          1,
		Name:        "特殊抽成规则",
		Description: "针对特定游戏、陪玩师和服务类型的抽成规则",
		Type:        "special",
		Rate:        15, // 15%
		IsActive:    true,
		GameID:      &gameID,
		PlayerID:    &playerID,
		ServiceType: &serviceType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, uint64(1), commissionRule.ID)
	assert.Equal(t, "特殊抽成规则", commissionRule.Name)
	assert.Equal(t, "special", commissionRule.Type)
	assert.Equal(t, 15, commissionRule.Rate)
	assert.Equal(t, &gameID, commissionRule.GameID)
	assert.Equal(t, &playerID, commissionRule.PlayerID)
	assert.Equal(t, &serviceType, commissionRule.ServiceType)
}

func TestCommissionRecordModel(t *testing.T) {
	now := time.Now()
	settledAt := now.Add(1 * time.Hour)

	commissionRecord := &model.CommissionRecord{
		ID:                1,
		OrderID:           100,
		PlayerID:          200,
		TotalAmountCents:  10000, // 100元
		CommissionRate:    20,    // 20%
		CommissionCents:   2000,  // 20元
		PlayerIncomeCents: 8000,  // 80元
		SettlementStatus:  "settled",
		SettlementMonth:   "2024-01",
		SettledAt:         &settledAt,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	assert.Equal(t, uint64(1), commissionRecord.ID)
	assert.Equal(t, uint64(100), commissionRecord.OrderID)
	assert.Equal(t, uint64(200), commissionRecord.PlayerID)
	assert.Equal(t, int64(10000), commissionRecord.TotalAmountCents)
	assert.Equal(t, 20, commissionRecord.CommissionRate)
	assert.Equal(t, int64(2000), commissionRecord.CommissionCents)
	assert.Equal(t, int64(8000), commissionRecord.PlayerIncomeCents)
	assert.Equal(t, "settled", commissionRecord.SettlementStatus)
	assert.Equal(t, "2024-01", commissionRecord.SettlementMonth)
	assert.Equal(t, &settledAt, commissionRecord.SettledAt)
	assert.Equal(t, now, commissionRecord.CreatedAt)
	assert.Equal(t, now, commissionRecord.UpdatedAt)
}

func TestCommissionRecordJSONSerialization(t *testing.T) {
	now := time.Now()

	commissionRecord := &model.CommissionRecord{
		ID:                1,
		OrderID:           100,
		PlayerID:          200,
		TotalAmountCents:  10000,
		CommissionRate:    20,
		CommissionCents:   2000,
		PlayerIncomeCents: 8000,
		SettlementStatus:  "pending",
		SettlementMonth:   "2024-01",
		SettledAt:         nil,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// 序列化
	data, err := json.Marshal(commissionRecord)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "2024-01")

	// 反序列化
	var decoded model.CommissionRecord
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, commissionRecord.ID, decoded.ID)
	assert.Equal(t, commissionRecord.OrderID, decoded.OrderID)
	assert.Equal(t, commissionRecord.PlayerID, decoded.PlayerID)
	assert.Equal(t, commissionRecord.TotalAmountCents, decoded.TotalAmountCents)
	assert.Equal(t, commissionRecord.SettlementStatus, decoded.SettlementStatus)
}

func TestCommissionRecordTableName(t *testing.T) {
	commissionRecord := model.CommissionRecord{}
	tableName := commissionRecord.TableName()
	assert.Equal(t, "commission_records", tableName)
}

func TestCommissionRecordSettlementStatuses(t *testing.T) {
	// 测试待结算状态
	record1 := &model.CommissionRecord{
		SettlementStatus: "pending",
	}
	assert.Equal(t, "pending", record1.SettlementStatus)

	// 测试已结算状态
	record2 := &model.CommissionRecord{
		SettlementStatus: "settled",
	}
	assert.Equal(t, "settled", record2.SettlementStatus)
}

func TestCommissionRecordZeroValues(t *testing.T) {
	commissionRecord := &model.CommissionRecord{
		ID:                0,
		OrderID:           0,
		PlayerID:          0,
		TotalAmountCents:  0,
		CommissionRate:    0,
		CommissionCents:   0,
		PlayerIncomeCents: 0,
		SettlementStatus:  "",
		SettlementMonth:   "",
		SettledAt:         nil,
	}

	assert.Equal(t, uint64(0), commissionRecord.ID)
	assert.Equal(t, uint64(0), commissionRecord.OrderID)
	assert.Equal(t, uint64(0), commissionRecord.PlayerID)
	assert.Equal(t, int64(0), commissionRecord.TotalAmountCents)
	assert.Equal(t, 0, commissionRecord.CommissionRate)
	assert.Equal(t, int64(0), commissionRecord.CommissionCents)
	assert.Equal(t, int64(0), commissionRecord.PlayerIncomeCents)
	assert.Equal(t, "", commissionRecord.SettlementStatus)
	assert.Equal(t, "", commissionRecord.SettlementMonth)
	assert.Nil(t, commissionRecord.SettledAt)
}

func TestCommissionRecordEdgeCases(t *testing.T) {
	// 测试高金额
	record1 := &model.CommissionRecord{
		TotalAmountCents:  ^int64(0), // 最大int64值
		CommissionCents:   ^int64(0),
		PlayerIncomeCents: ^int64(0),
	}
	assert.Equal(t, ^int64(0), record1.TotalAmountCents)
	assert.Equal(t, ^int64(0), record1.CommissionCents)
	assert.Equal(t, ^int64(0), record1.PlayerIncomeCents)

	// 测试100%抽成
	record2 := &model.CommissionRecord{
		TotalAmountCents:  10000,
		CommissionRate:    100,
		CommissionCents:   10000,
		PlayerIncomeCents: 0,
	}
	assert.Equal(t, int64(10000), record2.TotalAmountCents)
	assert.Equal(t, 100, record2.CommissionRate)
	assert.Equal(t, int64(10000), record2.CommissionCents)
	assert.Equal(t, int64(0), record2.PlayerIncomeCents)

	// 测试0%抽成
	record3 := &model.CommissionRecord{
		TotalAmountCents:  5000,
		CommissionRate:    0,
		CommissionCents:   0,
		PlayerIncomeCents: 5000,
	}
	assert.Equal(t, int64(5000), record3.TotalAmountCents)
	assert.Equal(t, 0, record3.CommissionRate)
	assert.Equal(t, int64(0), record3.CommissionCents)
	assert.Equal(t, int64(5000), record3.PlayerIncomeCents)

	// 测试特殊月份格式
	record4 := &model.CommissionRecord{
		SettlementMonth: "2024-12",
	}
	assert.Equal(t, "2024-12", record4.SettlementMonth)
}

func TestMonthlySettlementModel(t *testing.T) {
	now := time.Now()
	settledAt := now.Add(2 * time.Hour)
	rank1 := 1
	rank2 := 5
	rank3 := 10

	monthlySettlement := &model.MonthlySettlement{
		ID:                   1,
		PlayerID:             300,
		SettlementMonth:      "2024-01",
		TotalOrderCount:      50,
		TotalAmountCents:     50000, // 500元
		TotalCommissionCents: 10000, // 100元
		TotalIncomeCents:     40000, // 400元
		BonusCents:           5000,  // 50元奖金
		FinalIncomeCents:     45000, // 450元最终收入
		Status:               "paid",
		IncomeRank:           &rank1,
		OrderRank:            &rank2,
		QualityRank:          &rank3,
		CreatedAt:            now,
		UpdatedAt:            now,
		SettledAt:            &settledAt,
	}

	assert.Equal(t, uint64(1), monthlySettlement.ID)
	assert.Equal(t, uint64(300), monthlySettlement.PlayerID)
	assert.Equal(t, "2024-01", monthlySettlement.SettlementMonth)
	assert.Equal(t, int64(50), monthlySettlement.TotalOrderCount)
	assert.Equal(t, int64(50000), monthlySettlement.TotalAmountCents)
	assert.Equal(t, int64(10000), monthlySettlement.TotalCommissionCents)
	assert.Equal(t, int64(40000), monthlySettlement.TotalIncomeCents)
	assert.Equal(t, int64(5000), monthlySettlement.BonusCents)
	assert.Equal(t, int64(45000), monthlySettlement.FinalIncomeCents)
	assert.Equal(t, "paid", monthlySettlement.Status)
	assert.Equal(t, &rank1, monthlySettlement.IncomeRank)
	assert.Equal(t, &rank2, monthlySettlement.OrderRank)
	assert.Equal(t, &rank3, monthlySettlement.QualityRank)
	assert.Equal(t, now, monthlySettlement.CreatedAt)
	assert.Equal(t, now, monthlySettlement.UpdatedAt)
	assert.Equal(t, &settledAt, monthlySettlement.SettledAt)
}

func TestMonthlySettlementJSONSerialization(t *testing.T) {
	now := time.Now()
	rank1 := 1

	monthlySettlement := &model.MonthlySettlement{
		ID:                   1,
		PlayerID:             300,
		SettlementMonth:      "2024-01",
		TotalOrderCount:      50,
		TotalAmountCents:     50000,
		TotalCommissionCents: 10000,
		TotalIncomeCents:     40000,
		BonusCents:           5000,
		FinalIncomeCents:     45000,
		Status:               "pending",
		IncomeRank:           &rank1,
		OrderRank:            nil,
		QualityRank:          nil,
		CreatedAt:            now,
		UpdatedAt:            now,
		SettledAt:            nil,
	}

	// 序列化
	data, err := json.Marshal(monthlySettlement)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "2024-01")

	// 反序列化
	var decoded model.MonthlySettlement
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, monthlySettlement.ID, decoded.ID)
	assert.Equal(t, monthlySettlement.PlayerID, decoded.PlayerID)
	assert.Equal(t, monthlySettlement.SettlementMonth, decoded.SettlementMonth)
	assert.Equal(t, monthlySettlement.TotalOrderCount, decoded.TotalOrderCount)
	assert.Equal(t, monthlySettlement.TotalAmountCents, decoded.TotalAmountCents)
	assert.Equal(t, monthlySettlement.Status, decoded.Status)
}

func TestMonthlySettlementTableName(t *testing.T) {
	monthlySettlement := model.MonthlySettlement{}
	tableName := monthlySettlement.TableName()
	assert.Equal(t, "monthly_settlements", tableName)
}

func TestMonthlySettlementStatuses(t *testing.T) {
	// 测试所有状态
	statuses := []string{"pending", "confirmed", "paid"}

	for _, status := range statuses {
		settlement := &model.MonthlySettlement{
			Status: status,
		}
		assert.Equal(t, status, settlement.Status)
	}
}

func TestMonthlySettlementZeroValues(t *testing.T) {
	monthlySettlement := &model.MonthlySettlement{
		ID:                   0,
		PlayerID:             0,
		SettlementMonth:      "",
		TotalOrderCount:      0,
		TotalAmountCents:     0,
		TotalCommissionCents: 0,
		TotalIncomeCents:     0,
		BonusCents:           0,
		FinalIncomeCents:     0,
		Status:               "",
		IncomeRank:           nil,
		OrderRank:            nil,
		QualityRank:          nil,
		SettledAt:            nil,
	}

	assert.Equal(t, uint64(0), monthlySettlement.ID)
	assert.Equal(t, uint64(0), monthlySettlement.PlayerID)
	assert.Equal(t, "", monthlySettlement.SettlementMonth)
	assert.Equal(t, int64(0), monthlySettlement.TotalOrderCount)
	assert.Equal(t, int64(0), monthlySettlement.TotalAmountCents)
	assert.Equal(t, int64(0), monthlySettlement.TotalCommissionCents)
	assert.Equal(t, int64(0), monthlySettlement.TotalIncomeCents)
	assert.Equal(t, int64(0), monthlySettlement.BonusCents)
	assert.Equal(t, int64(0), monthlySettlement.FinalIncomeCents)
	assert.Equal(t, "", monthlySettlement.Status)
	assert.Nil(t, monthlySettlement.IncomeRank)
	assert.Nil(t, monthlySettlement.OrderRank)
	assert.Nil(t, monthlySettlement.QualityRank)
	assert.Nil(t, monthlySettlement.SettledAt)
}

func TestMonthlySettlementEdgeCases(t *testing.T) {
	// 测试高数值
	settlement1 := &model.MonthlySettlement{
		TotalOrderCount:      ^int64(0),
		TotalAmountCents:     ^int64(0),
		TotalCommissionCents: ^int64(0),
		TotalIncomeCents:     ^int64(0),
		BonusCents:           ^int64(0),
		FinalIncomeCents:     ^int64(0),
	}
	assert.Equal(t, ^int64(0), settlement1.TotalOrderCount)
	assert.Equal(t, ^int64(0), settlement1.TotalAmountCents)
	assert.Equal(t, ^int64(0), settlement1.TotalCommissionCents)
	assert.Equal(t, ^int64(0), settlement1.TotalIncomeCents)
	assert.Equal(t, ^int64(0), settlement1.BonusCents)
	assert.Equal(t, ^int64(0), settlement1.FinalIncomeCents)

	// 测试高排名
	rank1 := ^int(0) >> 1 // 最大int值
	rank2 := 1
	rank3 := 100

	settlement2 := &model.MonthlySettlement{
		IncomeRank:  &rank1,
		OrderRank:   &rank2,
		QualityRank: &rank3,
	}
	assert.Equal(t, &rank1, settlement2.IncomeRank)
	assert.Equal(t, &rank2, settlement2.OrderRank)
	assert.Equal(t, &rank3, settlement2.QualityRank)

	// 测试负奖金（虽然业务上不合理）
	settlement3 := &model.MonthlySettlement{
		TotalIncomeCents: 10000,
		BonusCents:       -1000,
		FinalIncomeCents: 9000,
	}
	assert.Equal(t, int64(10000), settlement3.TotalIncomeCents)
	assert.Equal(t, int64(-1000), settlement3.BonusCents)
	assert.Equal(t, int64(9000), settlement3.FinalIncomeCents)
}

func TestCommissionRuleTypes(t *testing.T) {
	// 测试不同类型的抽成规则
	rule1 := &model.CommissionRule{
		Type: "default",
		Rate: 20,
	}
	assert.Equal(t, "default", rule1.Type)

	rule2 := &model.CommissionRule{
		Type: "special",
		Rate: 15,
	}
	assert.Equal(t, "special", rule2.Type)

	rule3 := &model.CommissionRule{
		Type: "gift",
		Rate: 30,
	}
	assert.Equal(t, "gift", rule3.Type)
}
