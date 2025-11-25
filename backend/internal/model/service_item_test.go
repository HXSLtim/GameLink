package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gamelink/internal/model"
)

func TestServiceItemModel(t *testing.T) {
	now := time.Now()
	serviceItem := &model.ServiceItem{
		ID:             1,
		ItemCode:       "SERVICE001",
		Name:           "护航服务",
		Description:    "专业护航服务",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		GameID:         nil,
		PlayerID:       nil,
		RankLevel:      "Diamond",
		BasePriceCents: 10000, // 100元
		ServiceHours:   2,
		CommissionRate: 0.20, // 20%
		MinUsers:       1,
		MaxPlayers:     1,
		Tags:           `["护航", "专业"]`,
		IconURL:        "https://example.com/icon.png",
		IsActive:       true,
		SortOrder:      10,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	assert.Equal(t, uint64(1), serviceItem.ID)
	assert.Equal(t, "SERVICE001", serviceItem.ItemCode)
	assert.Equal(t, "护航服务", serviceItem.Name)
	assert.Equal(t, "专业护航服务", serviceItem.Description)
	assert.Equal(t, "escort", serviceItem.Category)
	assert.Equal(t, model.SubCategorySolo, serviceItem.SubCategory)
	assert.Nil(t, serviceItem.GameID)
	assert.Nil(t, serviceItem.PlayerID)
	assert.Equal(t, "Diamond", serviceItem.RankLevel)
	assert.Equal(t, int64(10000), serviceItem.BasePriceCents)
	assert.Equal(t, 2, serviceItem.ServiceHours)
	assert.Equal(t, 0.20, serviceItem.CommissionRate)
	assert.Equal(t, 1, serviceItem.MinUsers)
	assert.Equal(t, 1, serviceItem.MaxPlayers)
	assert.Equal(t, `["护航", "专业"]`, serviceItem.Tags)
	assert.Equal(t, "https://example.com/icon.png", serviceItem.IconURL)
	assert.True(t, serviceItem.IsActive)
	assert.Equal(t, 10, serviceItem.SortOrder)
	assert.Equal(t, now, serviceItem.CreatedAt)
	assert.Equal(t, now, serviceItem.UpdatedAt)
}

func TestServiceItemJSONSerialization(t *testing.T) {
	now := time.Now()
	serviceItem := &model.ServiceItem{
		ID:             1,
		ItemCode:       "SERVICE001",
		Name:           "护航服务",
		Description:    "专业护航服务",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		GameID:         nil,
		PlayerID:       nil,
		RankLevel:      "Diamond",
		BasePriceCents: 10000,
		ServiceHours:   2,
		CommissionRate: 0.20,
		MinUsers:       1,
		MaxPlayers:     1,
		Tags:           `["护航", "专业"]`,
		IconURL:        "https://example.com/icon.png",
		IsActive:       true,
		SortOrder:      10,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 序列化
	data, err := json.Marshal(serviceItem)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "SERVICE001")
	assert.Contains(t, string(data), "护航服务")

	// 反序列化
	var decoded model.ServiceItem
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, serviceItem.ID, decoded.ID)
	assert.Equal(t, serviceItem.ItemCode, decoded.ItemCode)
	assert.Equal(t, serviceItem.Name, decoded.Name)
	assert.Equal(t, serviceItem.Category, decoded.Category)
	assert.Equal(t, serviceItem.SubCategory, decoded.SubCategory)
	assert.Equal(t, serviceItem.BasePriceCents, decoded.BasePriceCents)
	assert.Equal(t, serviceItem.CommissionRate, decoded.CommissionRate)
}

func TestServiceItemConstants(t *testing.T) {
	// 测试服务子类别常量
	assert.Equal(t, model.ServiceItemSubCategory("solo"), model.SubCategorySolo)
	assert.Equal(t, model.ServiceItemSubCategory("team"), model.SubCategoryTeam)
	assert.Equal(t, model.ServiceItemSubCategory("gift"), model.SubCategoryGift)
}

func TestServiceItemTableName(t *testing.T) {
	serviceItem := model.ServiceItem{}
	tableName := serviceItem.TableName()
	assert.Equal(t, "service_items", tableName)
}

func TestServiceItemIsGift(t *testing.T) {
	// 单人护航服务
	serviceItem1 := &model.ServiceItem{
		SubCategory: model.SubCategorySolo,
	}
	assert.False(t, serviceItem1.IsGift())

	// 团队护航服务
	serviceItem2 := &model.ServiceItem{
		SubCategory: model.SubCategoryTeam,
	}
	assert.False(t, serviceItem2.IsGift())

	// 礼物服务
	serviceItem3 := &model.ServiceItem{
		SubCategory: model.SubCategoryGift,
	}
	assert.True(t, serviceItem3.IsGift())
}

func TestServiceItemCalculateCommission(t *testing.T) {
	// 测试20%抽成率
	serviceItem1 := &model.ServiceItem{
		BasePriceCents: 10000, // 100元
		CommissionRate: 0.20,  // 20%
	}

	platformCommission, playerIncome := serviceItem1.CalculateCommission(1)
	assert.Equal(t, int64(2000), platformCommission) // 20元
	assert.Equal(t, int64(8000), playerIncome)       // 80元

	// 测试数量大于1的情况
	platformCommission2, playerIncome2 := serviceItem1.CalculateCommission(3)
	assert.Equal(t, int64(6000), platformCommission2) // 60元
	assert.Equal(t, int64(24000), playerIncome2)      // 240元

	// 测试0%抽成率
	serviceItem2 := &model.ServiceItem{
		BasePriceCents: 5000, // 50元
		CommissionRate: 0.0,  // 0%
	}

	platformCommission3, playerIncome3 := serviceItem2.CalculateCommission(1)
	assert.Equal(t, int64(0), platformCommission3)
	assert.Equal(t, int64(5000), playerIncome3)

	// 测试100%抽成率
	serviceItem3 := &model.ServiceItem{
		BasePriceCents: 8000, // 80元
		CommissionRate: 1.0,  // 100%
	}

	platformCommission4, playerIncome4 := serviceItem3.CalculateCommission(1)
	assert.Equal(t, int64(8000), platformCommission4)
	assert.Equal(t, int64(0), playerIncome4)

	// 测试50%抽成率
	serviceItem4 := &model.ServiceItem{
		BasePriceCents: 20000, // 200元
		CommissionRate: 0.50,  // 50%
	}

	platformCommission5, playerIncome5 := serviceItem4.CalculateCommission(1)
	assert.Equal(t, int64(10000), platformCommission5)
	assert.Equal(t, int64(10000), playerIncome5)
}

func TestServiceItemZeroValues(t *testing.T) {
	serviceItem := &model.ServiceItem{
		ID:             0,
		ItemCode:       "",
		Name:           "",
		Description:    "",
		Category:       "",
		SubCategory:    "",
		GameID:         nil,
		PlayerID:       nil,
		RankLevel:      "",
		BasePriceCents: 0,
		ServiceHours:   0,
		CommissionRate: 0.0,
		MinUsers:       0,
		MaxPlayers:     0,
		Tags:           "",
		IconURL:        "",
		IsActive:       false,
		SortOrder:      0,
	}

	assert.Equal(t, uint64(0), serviceItem.ID)
	assert.Equal(t, "", serviceItem.ItemCode)
	assert.Equal(t, "", serviceItem.Name)
	assert.Equal(t, "", serviceItem.Description)
	assert.Equal(t, "", serviceItem.Category)
	assert.Equal(t, model.ServiceItemSubCategory(""), serviceItem.SubCategory)
	assert.Nil(t, serviceItem.GameID)
	assert.Nil(t, serviceItem.PlayerID)
	assert.Equal(t, "", serviceItem.RankLevel)
	assert.Equal(t, int64(0), serviceItem.BasePriceCents)
	assert.Equal(t, 0, serviceItem.ServiceHours)
	assert.Equal(t, 0.0, serviceItem.CommissionRate)
	assert.Equal(t, 0, serviceItem.MinUsers)
	assert.Equal(t, 0, serviceItem.MaxPlayers)
	assert.Equal(t, "", serviceItem.Tags)
	assert.Equal(t, "", serviceItem.IconURL)
	assert.False(t, serviceItem.IsActive)
	assert.Equal(t, 0, serviceItem.SortOrder)
}

func TestServiceItemEdgeCases(t *testing.T) {
	// 测试长字符串
	longCode := "VERY_LONG_SERVICE_ITEM_CODE_123456789"
	longName := "这是一个非常非常长的服务名称，用于测试字符串长度的边界情况"
	longDescription := "这是一个非常长的描述，可以包含很多关于服务的详细信息，比如服务的特点、优势、适用场景、注意事项等等。"
	longRank := "VeryLongRankName123456789"

	serviceItem1 := &model.ServiceItem{
		ItemCode:    longCode,
		Name:        longName,
		Description: longDescription,
		RankLevel:   longRank,
	}

	assert.Equal(t, longCode, serviceItem1.ItemCode)
	assert.Equal(t, longName, serviceItem1.Name)
	assert.Equal(t, longDescription, serviceItem1.Description)
	assert.Equal(t, longRank, serviceItem1.RankLevel)

	// 测试特殊字符
	serviceItem2 := &model.ServiceItem{
		ItemCode: "CODE_123#测试",
		Name:     "服务@#$%^&*()",
		Tags:     `["标签1", "tag2", "测试#标签"]`,
	}
	assert.Equal(t, "CODE_123#测试", serviceItem2.ItemCode)
	assert.Equal(t, "服务@#$%^&*()", serviceItem2.Name)
	assert.Equal(t, `["标签1", "tag2", "测试#标签"]`, serviceItem2.Tags)

	// 测试大数值
	serviceItem3 := &model.ServiceItem{
		BasePriceCents: ^int64(0), // 最大int64值
		ServiceHours:   ^int(0) >> 1, // 最大int值
		MinUsers:       ^int(0) >> 1,
		MaxPlayers:     ^int(0) >> 1,
		SortOrder:      ^int(0) >> 1,
	}

	assert.Equal(t, ^int64(0), serviceItem3.BasePriceCents)
	assert.Equal(t, ^int(0)>>1, serviceItem3.ServiceHours)
	assert.Equal(t, ^int(0)>>1, serviceItem3.MinUsers)
	assert.Equal(t, ^int(0)>>1, serviceItem3.MaxPlayers)
	assert.Equal(t, ^int(0)>>1, serviceItem3.SortOrder)

	// 测试浮点数精度
	serviceItem4 := &model.ServiceItem{
		CommissionRate: 0.333333,
	}
	assert.Equal(t, 0.333333, serviceItem4.CommissionRate)

	serviceItem5 := &model.ServiceItem{
		CommissionRate: 0.15,
	}
	assert.Equal(t, 0.15, serviceItem5.CommissionRate)
}

func TestServiceItemWithGameAndPlayerIDs(t *testing.T) {
	gameID := uint64(10)
	playerID := uint64(20)

	serviceItem := &model.ServiceItem{
		GameID:   &gameID,
		PlayerID: &playerID,
	}

	assert.NotNil(t, serviceItem.GameID)
	assert.Equal(t, uint64(10), *serviceItem.GameID)
	assert.NotNil(t, serviceItem.PlayerID)
	assert.Equal(t, uint64(20), *serviceItem.PlayerID)
}

func TestServiceItemJSONFields(t *testing.T) {
	serviceItem := &model.ServiceItem{
		ID:             1,
		ItemCode:       "SERVICE001",
		Name:           "护航服务",
		Description:    "专业护航服务",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		RankLevel:      "Diamond",
		BasePriceCents: 10000,
		ServiceHours:   2,
		CommissionRate: 0.20,
		MinUsers:       1,
		MaxPlayers:     1,
		Tags:           `["护航", "专业"]`,
		IconURL:        "https://example.com/icon.png",
		IsActive:       true,
		SortOrder:      10,
	}

	data, err := json.Marshal(serviceItem)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// 检查必需的字段
	assert.Contains(t, result, "id")
	assert.Contains(t, result, "itemCode")
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "description")
	assert.Contains(t, result, "category")
	assert.Contains(t, result, "subCategory")
	assert.Contains(t, result, "basePriceCents")
	assert.Contains(t, result, "commissionRate")
	assert.Contains(t, result, "isActive")

	// 验证值
	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, "SERVICE001", result["itemCode"])
	assert.Equal(t, "护航服务", result["name"])
	assert.Equal(t, "escort", result["category"])
	assert.Equal(t, "solo", result["subCategory"])
	assert.Equal(t, float64(10000), result["basePriceCents"])
	assert.Equal(t, 0.2, result["commissionRate"])
	assert.Equal(t, true, result["isActive"])
}

func TestServiceItemCalculateCommissionWithZeroQuantity(t *testing.T) {
	serviceItem := &model.ServiceItem{
		BasePriceCents: 10000,
		CommissionRate: 0.20,
	}

	platformCommission, playerIncome := serviceItem.CalculateCommission(0)
	assert.Equal(t, int64(0), platformCommission)
	assert.Equal(t, int64(0), playerIncome)
}

func TestServiceItemCalculateCommissionWithNegativeCommissionRate(t *testing.T) {
	// 测试负抽成率（虽然业务上不合理，但测试边界情况）
	serviceItem := &model.ServiceItem{
		BasePriceCents: 10000,
		CommissionRate: -0.1, // -10%
	}

	platformCommission, playerIncome := serviceItem.CalculateCommission(1)
	assert.Equal(t, int64(-1000), platformCommission) // -10元
	assert.Equal(t, int64(11000), playerIncome)       // 110元
}