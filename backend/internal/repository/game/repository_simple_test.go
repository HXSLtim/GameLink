package game_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	gameRepo "gamelink/internal/repository/game"
)

type GameRepositorySimpleTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo repository.GameRepository
	ctx  context.Context
}

func (s *GameRepositorySimpleTestSuite) SetupSuite() {
	// 使用SQLite内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Skipf("sqlite driver unavailable: %v", err)
		return
	}

	// 迁移表结构
	err = db.AutoMigrate(&model.Game{})
	if err != nil {
		s.T().Skipf("sqlite migration unavailable: %v", err)
		return
	}

	s.db = db
	s.repo = gameRepo.NewGameRepository(db)
	s.ctx = context.Background()
}

func (s *GameRepositorySimpleTestSuite) TearDownSuite() {
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (s *GameRepositorySimpleTestSuite) SetupTest() {
	if s.db == nil {
		s.T().Skip("sqlite unavailable")
		return
	}
	// 清空游戏表
	s.db.Exec("DELETE FROM games")
}

func TestGameRepositorySimpleTestSuite(t *testing.T) {
	suite.Run(t, new(GameRepositorySimpleTestSuite))
}

func (s *GameRepositorySimpleTestSuite) TestBasicCRUD() {
	// 测试创建游戏
	game := &model.Game{
		Key:         "test-game-1",
		Name:        "Test Game 1",
		Category:    "action",
		IconURL:     "https://example.com/game1.png",
		Description: "This is a test game",
	}

	err := s.repo.Create(s.ctx, game)
	s.NoError(err)
	s.NotZero(game.ID)

	// 测试获取游戏
	found, err := s.repo.Get(s.ctx, game.ID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(game.ID, found.ID)
	s.Equal(game.Key, found.Key)
	s.Equal(game.Name, found.Name)

	// 测试更新游戏
	game.Name = "Updated Game Name"
	err = s.repo.Update(s.ctx, game)
	s.NoError(err)

	updated, err := s.repo.Get(s.ctx, game.ID)
	s.NoError(err)
	s.Equal("Updated Game Name", updated.Name)

	// 测试删除游戏
	err = s.repo.Delete(s.ctx, game.ID)
	s.NoError(err)

	deleted, err := s.repo.Get(s.ctx, game.ID)
	s.Error(err)
	s.Nil(deleted)
	s.Equal(repository.ErrNotFound, err)
}

func (s *GameRepositorySimpleTestSuite) TestListOperations() {
	// 创建多个游戏
	games := []*model.Game{
		{
			Key:         "game-1",
			Name:        "Game 1",
			Category:    "action",
			Description: "First game",
		},
		{
			Key:         "game-2",
			Name:        "Game 2",
			Category:    "strategy",
			Description: "Second game",
		},
		{
			Key:         "game-3",
			Name:        "Game 3",
			Category:    "rpg",
			Description: "Third game",
		},
	}

	for _, game := range games {
		err := s.repo.Create(s.ctx, game)
		s.NoError(err)
	}

	// 测试列表所有游戏
	list, err := s.repo.List(s.ctx)
	s.NoError(err)
	s.Len(list, 3)

	// 测试分页
	pagedList, total, err := s.repo.ListPaged(s.ctx, 1, 2)
	s.NoError(err)
	s.Len(pagedList, 2)
	s.Equal(int64(3), total)

	// 第二页
	pagedList, total, err = s.repo.ListPaged(s.ctx, 2, 2)
	s.NoError(err)
	s.Len(pagedList, 1)
	s.Equal(int64(3), total)
}

func (s *GameRepositorySimpleTestSuite) TestUniqueConstraints() {
	// 创建第一个游戏
	game1 := &model.Game{
		Key:         "unique-key-1",
		Name:        "Game 1",
		Category:    "action",
		Description: "First game",
	}

	err := s.repo.Create(s.ctx, game1)
	s.NoError(err)

	// 测试重复key应该失败
	game2 := &model.Game{
		Key:         "unique-key-1", // 重复的key
		Name:        "Game 2",
		Category:    "strategy",
		Description: "Second game",
	}

	err = s.repo.Create(s.ctx, game2)
	s.Error(err)
}

func (s *GameRepositorySimpleTestSuite) TestGameCategories() {
	// 创建不同类别的游戏
	categories := []string{"action", "strategy", "rpg", "puzzle", "sports"}
	for i, category := range categories {
		game := &model.Game{
			Key:         "category-game-" + string(rune(i)),
			Name:        "Category Game " + string(rune(i)),
			Category:    category,
			Description: "Game in " + category + " category",
		}
		err := s.repo.Create(s.ctx, game)
		s.NoError(err)
	}

	// 验证所有游戏都被创建
	list, err := s.repo.List(s.ctx)
	s.NoError(err)
	s.Len(list, len(categories))

	// 验证每个游戏都有不同的类别
	categoryMap := make(map[string]bool)
	for _, game := range list {
		categoryMap[game.Category] = true
	}
	s.Equal(len(categories), len(categoryMap))
}
