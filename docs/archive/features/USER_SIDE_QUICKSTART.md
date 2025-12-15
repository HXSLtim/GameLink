# ⚡ GameLink 用户侧开发快速开始

**目标**: 5分钟内启动第一个用户侧功能的开发

---

## 🎯 开发前准备

### 确认环境

```bash
# 1. 检查 Go 版本
go version  # 应该是 1.24+

# 2. 检查 Node 版本
node --version  # 应该是 18+

# 3. 确认后端测试通过
cd backend
go test ./...  # 应该全部 PASS

# 4. 确认前端可运行
cd frontend
npm run dev    # 应该能启动
```

---

## 🚀 第一个功能：陪玩师列表

### 步骤1: 创建后端 API（15分钟）

**1.1 创建 User Handler 目录**

```bash
cd backend/internal/handler
mkdir user
cd user
```

**1.2 创建路由文件** `router.go`

```go
package user

import (
	"github.com/gin-gonic/gin"
	userservice "gamelink/internal/service/user"
)

type Handlers struct {
	Player *PlayerHandler
}

func NewHandlers(playerSvc *userservice.PlayerService) *Handlers {
	return &Handlers{
		Player: NewPlayerHandler(playerSvc),
	}
}

func RegisterUserRoutes(router gin.IRouter, handlers *Handlers) {
	user := router.Group("/user")
	{
		// 陪玩师
		user.GET("/players", handlers.Player.ListPlayers)
		user.GET("/players/:id", handlers.Player.GetPlayerDetail)
	}
}
```

**1.3 创建 Player Handler** `player_handler.go`

```go
package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	userservice "gamelink/internal/service/user"
)

type PlayerHandler struct {
	playerService *userservice.PlayerService
}

func NewPlayerHandler(svc *userservice.PlayerService) *PlayerHandler {
	return &PlayerHandler{playerService: svc}
}

// ListPlayers godoc
// @Summary 获取陪玩师列表
// @Tags User-Player
// @Param gameId query int false "游戏ID"
// @Param minPrice query int false "最低价格（分）"
// @Param maxPrice query int false "最高价格（分）"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} model.APIResponse{data=userservice.ListPlayersResponse}
// @Router /api/v1/user/players [get]
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	var req userservice.ListPlayersRequest
	
	// 解析查询参数
	if gameIDStr := c.Query("gameId"); gameIDStr != "" {
		gameID, _ := strconv.ParseUint(gameIDStr, 10, 64)
		req.GameID = &gameID
	}
	if minPriceStr := c.Query("minPrice"); minPriceStr != "" {
		minPrice, _ := strconv.ParseInt(minPriceStr, 10, 64)
		req.MinPrice = &minPrice
	}
	if maxPriceStr := c.Query("maxPrice"); maxPriceStr != "" {
		maxPrice, _ := strconv.ParseInt(maxPriceStr, 10, 64)
		req.MaxPrice = &maxPrice
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	req.Page = page
	req.PageSize = pageSize

	// 调用服务
	resp, err := h.playerService.ListPlayers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// GetPlayerDetail godoc
// @Summary 获取陪玩师详情
// @Tags User-Player
// @Param id path int true "陪玩师ID"
// @Success 200 {object} model.APIResponse{data=userservice.PlayerDetailResponse}
// @Router /api/v1/user/players/{id} [get]
func (h *PlayerHandler) GetPlayerDetail(c *gin.Context) {
	playerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的陪玩师ID",
		})
		return
	}

	resp, err := h.playerService.GetPlayerDetail(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}
```

**1.4 创建 User Service** `backend/internal/service/user/player_service.go`

```go
package user

import (
	"context"
	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type PlayerService struct {
	playerRepo repository.PlayerRepository
	gameRepo   repository.GameRepository
	reviewRepo repository.ReviewRepository
}

func NewPlayerService(
	playerRepo repository.PlayerRepository,
	gameRepo repository.GameRepository,
	reviewRepo repository.ReviewRepository,
) *PlayerService {
	return &PlayerService{
		playerRepo: playerRepo,
		gameRepo:   gameRepo,
		reviewRepo: reviewRepo,
	}
}

type ListPlayersRequest struct {
	GameID   *uint64
	MinPrice *int64
	MaxPrice *int64
	Page     int
	PageSize int
}

type PlayerCardDTO struct {
	ID              uint64  `json:"id"`
	UserID          uint64  `json:"userId"`
	Nickname        string  `json:"nickname"`
	Bio             string  `json:"bio"`
	Rank            string  `json:"rank"`
	RatingAverage   float32 `json:"ratingAverage"`
	RatingCount     uint32  `json:"ratingCount"`
	HourlyRateCents int64   `json:"hourlyRateCents"`
	MainGame        string  `json:"mainGame"`
	IsOnline        bool    `json:"isOnline"`
}

type ListPlayersResponse struct {
	Players []PlayerCardDTO `json:"players"`
	Total   int64           `json:"total"`
}

func (s *PlayerService) ListPlayers(ctx context.Context, req ListPlayersRequest) (*ListPlayersResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 获取陪玩师列表
	players, total, err := s.playerRepo.ListPaged(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	dtos := make([]PlayerCardDTO, len(players))
	for i, p := range players {
		// TODO: 获取游戏名称、在线状态等
		dtos[i] = PlayerCardDTO{
			ID:              p.ID,
			UserID:          p.UserID,
			Nickname:        p.Nickname,
			Bio:             p.Bio,
			Rank:            p.Rank,
			RatingAverage:   p.RatingAverage,
			RatingCount:     p.RatingCount,
			HourlyRateCents: p.HourlyRateCents,
			MainGame:        "游戏名称", // TODO
			IsOnline:        false,     // TODO
		}
	}

	return &ListPlayersResponse{
		Players: dtos,
		Total:   total,
	}, nil
}

type PlayerDetailResponse struct {
	Player PlayerCardDTO `json:"player"`
	// TODO: 添加更多详情信息
}

func (s *PlayerService) GetPlayerDetail(ctx context.Context, playerID uint64) (*PlayerDetailResponse, error) {
	player, err := s.playerRepo.Get(ctx, playerID)
	if err != nil {
		return nil, err
	}

	dto := PlayerCardDTO{
		ID:              player.ID,
		UserID:          player.UserID,
		Nickname:        player.Nickname,
		Bio:             player.Bio,
		Rank:            player.Rank,
		RatingAverage:   player.RatingAverage,
		RatingCount:     player.RatingCount,
		HourlyRateCents: player.HourlyRateCents,
		MainGame:        "游戏名称",
		IsOnline:        false,
	}

	return &PlayerDetailResponse{
		Player: dto,
	}, nil
}
```

**1.5 在 main.go 中注册路由**

编辑 `backend/cmd/user-service/main.go`：

```go
// 在现有导入中添加
import (
    userhandler "gamelink/internal/handler/user"
    userservice "gamelink/internal/service/user"
)

// 在 setupRoutes 函数中添加
func setupRoutes(r *gin.Engine, deps *dependencies) {
    // ... 现有代码 ...
    
    // 用户端路由
    playerService := userservice.NewPlayerService(
        deps.playerRepo,
        deps.gameRepo,
        deps.reviewRepo,
    )
    userHandlers := userhandler.NewHandlers(playerService)
    userhandler.RegisterUserRoutes(api, userHandlers)
}
```

**1.6 测试 API**

```bash
# 启动后端
cd backend
make run CMD=user-service

# 在另一个终端测试
curl http://localhost:8080/api/v1/user/players?page=1&pageSize=10
```

---

### 步骤2: 创建前端页面（15分钟）

**2.1 创建类型定义** `frontend/src/types/player.ts`

```typescript
export interface PlayerCardDTO {
  id: number;
  userId: number;
  nickname: string;
  bio: string;
  rank: string;
  ratingAverage: number;
  ratingCount: number;
  hourlyRateCents: number;
  mainGame: string;
  isOnline: boolean;
}

export interface ListPlayersResponse {
  players: PlayerCardDTO[];
  total: number;
}
```

**2.2 创建 API 服务** `frontend/src/services/playerApi.ts`

```typescript
import { request } from './request';
import { PlayerCardDTO, ListPlayersResponse } from '@/types/player';

export interface PlayerListParams {
  gameId?: number;
  minPrice?: number;
  maxPrice?: number;
  page: number;
  pageSize: number;
}

export const playerApi = {
  async getPlayers(params: PlayerListParams) {
    const res = await request.get<ListPlayersResponse>(
      '/api/v1/user/players',
      { params }
    );
    return res.data;
  },

  async getPlayerDetail(id: number) {
    const res = await request.get<{ player: PlayerCardDTO }>(
      `/api/v1/user/players/${id}`
    );
    return res.data;
  },
};
```

**2.3 创建陪玩师卡片组件** `frontend/src/components/PlayerCard/PlayerCard.tsx`

```tsx
import React from 'react';
import { Card, Avatar, Tag, Rate } from '@arco-design/web-react';
import { PlayerCardDTO } from '@/types/player';
import styles from './PlayerCard.module.less';

interface PlayerCardProps {
  player: PlayerCardDTO;
  onClick?: () => void;
}

export const PlayerCard: React.FC<PlayerCardProps> = ({ player, onClick }) => {
  const priceYuan = (player.hourlyRateCents / 100).toFixed(2);

  return (
    <Card className={styles.card} hoverable onClick={onClick}>
      <div className={styles.header}>
        <Avatar size={64} src={`https://via.placeholder.com/64`}>
          {player.nickname[0]}
        </Avatar>
        {player.isOnline && (
          <Tag color="green" className={styles.onlineTag}>
            在线
          </Tag>
        )}
      </div>
      <div className={styles.content}>
        <h3 className={styles.nickname}>{player.nickname}</h3>
        <div className={styles.game}>
          {player.mainGame} · {player.rank}
        </div>
        <div className={styles.rating}>
          <Rate readonly value={player.ratingAverage} />
          <span>({player.ratingCount})</span>
        </div>
        <div className={styles.price}>¥{priceYuan}/小时</div>
      </div>
    </Card>
  );
};
```

**2.4 创建样式文件** `frontend/src/components/PlayerCard/PlayerCard.module.less`

```less
.card {
  height: 100%;
  
  .header {
    text-align: center;
    position: relative;
    margin-bottom: 16px;
    
    .onlineTag {
      position: absolute;
      top: 0;
      right: 0;
    }
  }
  
  .content {
    text-align: center;
    
    .nickname {
      margin: 8px 0;
      font-size: 16px;
      font-weight: 500;
    }
    
    .game {
      color: #86909c;
      font-size: 14px;
      margin-bottom: 8px;
    }
    
    .rating {
      margin-bottom: 12px;
      
      span {
        margin-left: 8px;
        color: #86909c;
        font-size: 12px;
      }
    }
    
    .price {
      color: #f53f3f;
      font-size: 18px;
      font-weight: 600;
    }
  }
}
```

**2.5 创建列表页** `frontend/src/pages/Players/PlayerList.tsx`

```tsx
import React, { useState, useEffect } from 'react';
import { Grid, Pagination, Spin, Empty } from '@arco-design/web-react';
import { PlayerCard } from '@/components/PlayerCard';
import { playerApi } from '@/services/playerApi';
import { useNavigate } from 'react-router-dom';
import { PlayerCardDTO } from '@/types/player';
import styles from './PlayerList.module.less';

const { Row, Col } = Grid;

export const PlayerList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [players, setPlayers] = useState<PlayerCardDTO[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);

  const fetchPlayers = async () => {
    setLoading(true);
    try {
      const data = await playerApi.getPlayers({ page, pageSize });
      setPlayers(data.players || []);
      setTotal(data.total);
    } catch (error) {
      console.error('获取陪玩师列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPlayers();
  }, [page]);

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1>陪玩师列表</h1>
      </div>

      <Spin loading={loading}>
        {players.length === 0 && !loading ? (
          <Empty description="暂无陪玩师" />
        ) : (
          <>
            <Row gutter={[16, 16]}>
              {players.map((player) => (
                <Col key={player.id} xs={24} sm={12} md={8} lg={6}>
                  <PlayerCard
                    player={player}
                    onClick={() => navigate(`/players/${player.id}`)}
                  />
                </Col>
              ))}
            </Row>

            {total > pageSize && (
              <div className={styles.pagination}>
                <Pagination
                  total={total}
                  current={page}
                  pageSize={pageSize}
                  onChange={setPage}
                  showTotal
                />
              </div>
            )}
          </>
        )}
      </Spin>
    </div>
  );
};
```

**2.6 创建样式** `frontend/src/pages/Players/PlayerList.module.less`

```less
.container {
  padding: 24px;
  
  .header {
    margin-bottom: 24px;
    
    h1 {
      margin: 0;
      font-size: 24px;
      font-weight: 600;
    }
  }
  
  .pagination {
    margin-top: 24px;
    display: flex;
    justify-content: center;
  }
}
```

**2.7 添加路由** `frontend/src/router/routes.tsx`

```tsx
import { PlayerList } from '@/pages/Players/PlayerList';

// 在路由配置中添加
{
  path: '/players',
  element: <PlayerList />,
},
```

**2.8 测试前端**

```bash
cd frontend
npm run dev

# 浏览器访问
http://localhost:5173/players
```

---

## ✅ 验收标准

完成以上步骤后，你应该能够：

- ✅ 后端 API `/api/v1/user/players` 可以返回数据
- ✅ 前端页面 `/players` 能够展示陪玩师列表
- ✅ 点击陪玩师卡片能够跳转到详情页（虽然详情页还未实现）
- ✅ 分页功能正常工作

---

## 🎉 恭喜！

你已经完成了第一个用户侧功能的开发！

### 下一步

1. **完善陪玩师详情页** - 参考 `USER_SIDE_IMPLEMENTATION.md`
2. **添加筛选功能** - 游戏、价格、评分筛选
3. **实现订单创建** - 允许用户预约陪玩师

---

## 📚 相关文档

- [完整功能规划](./USER_SIDE_PLANNING.md)
- [详细实施计划](./USER_SIDE_IMPLEMENTATION.md)
- [Go 编码规范](./go-coding-standards.md)
- [API 设计规范](./api-design-standards.md)

---

## 💡 开发小贴士

### 后端开发

1. **遵循分层架构**
   - Handler: 处理 HTTP 请求
   - Service: 业务逻辑
   - Repository: 数据访问

2. **错误处理**
   ```go
   if err != nil {
       return nil, fmt.Errorf("failed to get player: %w", err)
   }
   ```

3. **使用统一响应格式**
   ```go
   c.JSON(http.StatusOK, gin.H{
       "success": true,
       "data": resp,
   })
   ```

### 前端开发

1. **类型安全**
   - 所有 API 响应都定义 TypeScript 类型
   - 使用 interface 而不是 type（除非需要联合类型）

2. **组件复用**
   - 通用组件放在 `components/`
   - 页面特定组件放在页面目录下的 `components/`

3. **样式隔离**
   - 使用 CSS Modules（`.module.less`）
   - 避免全局样式污染

---

**Happy Coding! 🚀**


