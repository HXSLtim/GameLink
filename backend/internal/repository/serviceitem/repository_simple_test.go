package serviceitem_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	serviceItemRepo "gamelink/internal/repository/serviceitem"
)

type ServiceItemRepositorySimpleTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo serviceItemRepo.ServiceItemRepository
	ctx  context.Context
}

func (s *ServiceItemRepositorySimpleTestSuite) SetupSuite() {
	// 使用SQLite内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Skipf("sqlite driver unavailable: %v", err)
		return
	}

	// 迁移表结构 - 需要服务项相关的表
	err = db.AutoMigrate(&model.ServiceItem{})
	if err != nil {
		s.T().Skipf("sqlite migration unavailable: %v", err)
		return
	}

	s.db = db
	s.repo = serviceItemRepo.NewServiceItemRepository(db)
	s.ctx = context.Background()
}

func (s *ServiceItemRepositorySimpleTestSuite) TearDownSuite() {
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (s *ServiceItemRepositorySimpleTestSuite) SetupTest() {
	// 清空服务项表
	s.db.Exec("DELETE FROM service_items")
}

func TestServiceItemRepositorySimpleTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceItemRepositorySimpleTestSuite))
}

func (s *ServiceItemRepositorySimpleTestSuite) TestBasicCRUD() {
	// 测试创建服务项
	gameID := uint64(1)
	item := &model.ServiceItem{
		GameID:         &gameID,
		ItemCode:       "SERVICE-001",
		Name:           "Test Service",
		Description:    "This is a test service item",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 1000,
		IsActive:       true,
		SortOrder:      1,
		CommissionRate: 0.2,
	}

	err := s.repo.Create(s.ctx, item)
	s.NoError(err)
	s.NotZero(item.ID)

	// 测试获取服务项
	found, err := s.repo.Get(s.ctx, item.ID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(item.ID, found.ID)
	s.Equal(item.ItemCode, found.ItemCode)
	s.Equal(item.Name, found.Name)

	// 测试通过item code获取服务项
	foundByCode, err := s.repo.GetByCode(s.ctx, item.ItemCode)
	s.NoError(err)
	s.NotNil(foundByCode)
	s.Equal(item.ID, foundByCode.ID)

	// 测试更新服务项
	item.Name = "Updated Service Name"
	err = s.repo.Update(s.ctx, item)
	s.NoError(err)

	updated, err := s.repo.Get(s.ctx, item.ID)
	s.NoError(err)
	s.Equal("Updated Service Name", updated.Name)

	// 测试删除服务项
	err = s.repo.Delete(s.ctx, item.ID)
	s.NoError(err)

	deleted, err := s.repo.Get(s.ctx, item.ID)
	s.Error(err)
	s.Nil(deleted)
	s.Equal(repository.ErrNotFound, err)
}

func (s *ServiceItemRepositorySimpleTestSuite) TestListOperations() {
	// 创建多个服务项
	gameID := uint64(10)
	items := []*model.ServiceItem{
		{
			GameID:         &gameID,
			ItemCode:       "SERVICE-001",
			Name:           "Service 1",
			Description:    "First service",
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 1000,
			IsActive:       true,
			SortOrder:      1,
			CommissionRate: 0.2,
		},
		{
			GameID:         &gameID,
			ItemCode:       "SERVICE-002",
			Name:           "Service 2",
			Description:    "Second service",
			Category:       "escort",
			SubCategory:    model.SubCategoryTeam,
			BasePriceCents: 1500,
			IsActive:       true,
			SortOrder:      2,
			CommissionRate: 0.25,
		},
		{
			GameID:         &gameID,
			ItemCode:       "GIFT-001",
			Name:           "Gift 1",
			Description:    "First gift",
			Category:       "gift",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 500,
			IsActive:       false,
			SortOrder:      3,
			CommissionRate: 0.1,
		},
	}

	for _, item := range items {
		err := s.repo.Create(s.ctx, item)
		s.NoError(err)
	}

	// 测试列表所有服务项
	list, total, err := s.repo.List(s.ctx, serviceItemRepo.ServiceItemListOptions{})
	s.NoError(err)
	s.Len(list, 3)
	s.Equal(int64(3), total)

	// 测试分页
	pagedList, total, err := s.repo.List(s.ctx, serviceItemRepo.ServiceItemListOptions{
		Page:     1,
		PageSize: 2,
	})
	s.NoError(err)
	s.Len(pagedList, 2)
	s.Equal(int64(3), total)

	// 第二页
	pagedList, total, err = s.repo.List(s.ctx, serviceItemRepo.ServiceItemListOptions{
		Page:     2,
		PageSize: 2,
	})
	s.NoError(err)
	s.Len(pagedList, 1)
	s.Equal(int64(3), total)
}

func (s *ServiceItemRepositorySimpleTestSuite) TestFilters() {
	// 创建测试服务项
	gameID1 := uint64(20)
	gameID2 := uint64(21)
	playerID1 := uint64(100)
	isActive := true
	isInactive := false

	items := []*model.ServiceItem{
		{
			GameID:         &gameID1,
			PlayerID:       &playerID1,
			ItemCode:       "FILTER-001",
			Name:           "Filter Service 1",
			Description:    "Active escort solo service",
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 2000,
			IsActive:       true,
			SortOrder:      1,
			CommissionRate: 0.2,
		},
		{
			GameID:         &gameID1,
			ItemCode:       "FILTER-002",
			Name:           "Filter Service 2",
			Description:    "Active escort team service",
			Category:       "escort",
			SubCategory:    model.SubCategoryTeam,
			BasePriceCents: 2500,
			IsActive:       true,
			SortOrder:      2,
			CommissionRate: 0.25,
		},
		{
			GameID:         &gameID2,
			ItemCode:       "FILTER-003",
			Name:           "Filter Gift 3",
			Description:    "Inactive gift",
			Category:       "gift",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 800,
			IsActive:       false,
			SortOrder:      3,
			CommissionRate: 0.1,
		},
	}

	for _, item := range items {
		err := s.repo.Create(s.ctx, item)
		s.NoError(err)
	}

	// 按游戏ID过滤
	list, total, err := s.repo.List(s.ctx, serviceItemRepo.ServiceItemListOptions{
		GameID:   &gameID1,
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.Len(list, 2)
	s.Equal(int64(2), total)

	// 按陪玩师ID过滤
	list, total, err = s.repo.List(s.ctx, serviceItemRepo.ServiceItemListOptions{
		PlayerID: &playerID1,
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.Len(list, 1)
	s.Equal(int64(1), total)
	s.Equal(&playerID1, list[0].PlayerID)

	// 按活跃状态过滤
	list, total, err = s.repo.List(s.ctx, serviceItemRepo.ServiceItemListOptions{
		IsActive: &isActive,
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.Len(list, 2)
	s.Equal(int64(2), total)

	// 按非活跃状态过滤
	list, total, err = s.repo.List(s.ctx, serviceItemRepo.ServiceItemListOptions{
		IsActive: &isInactive,
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.Len(list, 1)
	s.Equal(int64(1), total)
	s.False(list[0].IsActive)

	// 按子分类过滤
	soloSubCategory := model.SubCategorySolo
	list, total, err = s.repo.List(s.ctx, serviceItemRepo.ServiceItemListOptions{
		SubCategory: &soloSubCategory,
		Page:        1,
		PageSize:    10,
	})
	s.NoError(err)
	s.Len(list, 1)
	s.Equal(int64(1), total)
	s.Equal(model.SubCategorySolo, list[0].SubCategory)
}

func (s *ServiceItemRepositorySimpleTestSuite) TestBatchOperations() {
	// 创建多个服务项
	items := []*model.ServiceItem{
		{
			GameID:         nil,
			ItemCode:       "BATCH-001",
			Name:           "Batch Service 1",
			Description:    "First batch service",
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 3000,
			IsActive:       true,
			SortOrder:      1,
			CommissionRate: 0.2,
		},
		{
			GameID:         nil,
			ItemCode:       "BATCH-002",
			Name:           "Batch Service 2",
			Description:    "Second batch service",
			Category:       "escort",
			SubCategory:    model.SubCategoryTeam,
			BasePriceCents: 3500,
			IsActive:       true,
			SortOrder:      2,
			CommissionRate: 0.25,
		},
		{
			GameID:         nil,
			ItemCode:       "BATCH-003",
			Name:           "Batch Service 3",
			Description:    "Third batch service",
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 4000,
			IsActive:       true,
			SortOrder:      3,
			CommissionRate: 0.3,
		},
	}

	var ids []uint64
	for _, item := range items {
		err := s.repo.Create(s.ctx, item)
		s.NoError(err)
		ids = append(ids, item.ID)
	}

	// 批量更新状态
	err := s.repo.BatchUpdateStatus(s.ctx, ids[:2], false)
	s.NoError(err)

	// 验证更新结果
	for i, id := range ids {
		found, err := s.repo.Get(s.ctx, id)
		s.NoError(err)
		s.NotNil(found)
		if i < 2 {
			s.False(found.IsActive, "Items 1 and 2 should be inactive")
		} else {
			s.True(found.IsActive, "Item 3 should still be active")
		}
	}

	// 批量更新价格
	newPrice := int64(9999)
	err = s.repo.BatchUpdatePrice(s.ctx, ids, newPrice)
	s.NoError(err)

	// 验证所有服务项价格都被更新
	for _, id := range ids {
		found, err := s.repo.Get(s.ctx, id)
		s.NoError(err)
		s.NotNil(found)
		s.Equal(newPrice, found.BasePriceCents, "Price should be updated")
	}
}

func (s *ServiceItemRepositorySimpleTestSuite) TestGetGifts() {
	// 创建礼物服务项和普通服务项
	items := []*model.ServiceItem{
		{
			GameID:         nil,
			ItemCode:       "GIFT-TEST-001",
			Name:           "Test Gift 1",
			Description:    "First test gift",
			Category:       "gift",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 1000,
			IsActive:       true,
			SortOrder:      1,
			CommissionRate: 0.05,
		},
		{
			GameID:         nil,
			ItemCode:       "GIFT-TEST-002",
			Name:           "Test Gift 2",
			Description:    "Second test gift",
			Category:       "gift",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 1500,
			IsActive:       true,
			SortOrder:      2,
			CommissionRate: 0.08,
		},
		{
			GameID:         nil,
			ItemCode:       "GIFT-TEST-003",
			Name:           "Test Gift 3",
			Description:    "Third test gift (inactive)",
			Category:       "gift",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 2000,
			IsActive:       false, // 非活跃状态
			SortOrder:      3,
			CommissionRate: 0.1,
		},
		{
			GameID:         nil,
			ItemCode:       "ESCORT-TEST-001",
			Name:           "Test Escort Service",
			Description:    "Test escort service (not a gift)",
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 3000,
			IsActive:       true,
			SortOrder:      4,
			CommissionRate: 0.2,
		},
	}

	for _, item := range items {
		err := s.repo.Create(s.ctx, item)
		s.NoError(err)
	}

	// 获取礼物列表
	gifts, total, err := s.repo.GetGifts(s.ctx, 1, 10)
	s.NoError(err)
	s.Len(gifts, 2) // 只有活跃的礼物
	s.Equal(int64(2), total)

	// 验证所有返回的都是礼物
	for _, gift := range gifts {
		s.Equal(model.SubCategoryGift, gift.SubCategory)
		s.True(gift.IsActive)
	}
}

func (s *ServiceItemRepositorySimpleTestSuite) TestGetGameServices() {
	// 清空表
	s.db.Exec("DELETE FROM service_items")

	gameID := uint64(100)

	// 创建指定游戏的服务项
	items := []*model.ServiceItem{
		{
			GameID:         &gameID,
			ItemCode:       "GAME-SERVICE-001",
			Name:           "Game Service 1",
			Description:    "First game service",
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 5100,
			IsActive:       true,
			SortOrder:      1,
			CommissionRate: 0.2,
		},
		{
			GameID:         &gameID,
			ItemCode:       "GAME-SERVICE-002",
			Name:           "Game Service 2",
			Description:    "Second game service",
			Category:       "escort",
			SubCategory:    model.SubCategoryTeam,
			BasePriceCents: 5200,
			IsActive:       true,
			SortOrder:      2,
			CommissionRate: 0.25,
		},
		{
			GameID:         &gameID,
			ItemCode:       "GAME-SERVICE-003",
			Name:           "Game Service 3",
			Description:    "Third game service (inactive)",
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 5300,
			IsActive:       false, // 非活跃状态
			SortOrder:      3,
			CommissionRate: 0.3,
		},
		{
			GameID:         func() *uint64 { id := uint64(101); return &id }(), // 不同游戏ID
			ItemCode:       "OTHER-GAME-SERVICE",
			Name:           "Other Game Service",
			Description:    "Service for other game",
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 6000,
			IsActive:       true,
			SortOrder:      1,
			CommissionRate: 0.2,
		},
	}

	for _, item := range items {
		err := s.repo.Create(s.ctx, item)
		s.NoError(err)
	}

	// 获取指定游戏的活跃服务
	services, err := s.repo.GetGameServices(s.ctx, gameID, nil)
	s.NoError(err)
	s.Len(services, 2) // 只有活跃的指定游戏服务

	// 验证所有返回的都是指定游戏的活跃服务
	for _, service := range services {
		s.Equal(gameID, *service.GameID)
		s.True(service.IsActive)
	}

	// 验证按sort_order排序
	for i := 0; i < len(services)-1; i++ {
		s.LessOrEqual(services[i].SortOrder, services[i+1].SortOrder)
	}
}

func (s *ServiceItemRepositorySimpleTestSuite) TestUniqueConstraints() {
	// 创建第一个服务项
	item1 := &model.ServiceItem{
		GameID:         nil,
		ItemCode:       "UNIQUE-CODE-1",
		Name:           "Service 1",
		Description:    "First service",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 1000,
		IsActive:       true,
		SortOrder:      1,
		CommissionRate: 0.2,
	}

	err := s.repo.Create(s.ctx, item1)
	s.NoError(err)

	// 测试重复item code应该失败
	item2 := &model.ServiceItem{
		GameID:         nil,
		ItemCode:       "UNIQUE-CODE-1", // 重复的item code
		Name:           "Service 2",
		Description:    "Second service",
		Category:       "escort",
		SubCategory:    model.SubCategoryTeam,
		BasePriceCents: 1500,
		IsActive:       true,
		SortOrder:      2,
		CommissionRate: 0.25,
	}

	err = s.repo.Create(s.ctx, item2)
	s.Error(err)
}
