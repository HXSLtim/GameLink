# 🚀 GameLink 用户侧实施计划

**执行时间**: 2025年10月30日 - 2025年12月中旬  
**开发周期**: 7周  
**开发模式**: 前后端并行开发

---

## 📂 目录结构规划

### 后端目录结构

```
backend/
├── internal/
│   ├── handler/
│   │   ├── user/              # 用户端 Handler（新增）
│   │   │   ├── game_handler.go
│   │   │   ├── player_handler.go
│   │   │   ├── order_handler.go
│   │   │   ├── payment_handler.go
│   │   │   ├── review_handler.go
│   │   │   ├── profile_handler.go
│   │   │   └── router.go
│   │   │
│   │   ├── player/            # 陪玩师端 Handler（新增）
│   │   │   ├── profile_handler.go
│   │   │   ├── order_handler.go
│   │   │   ├── earnings_handler.go
│   │   │   ├── stats_handler.go
│   │   │   └── router.go
│   │   │
│   │   └── admin/             # 管理端（已有）
│   │
│   ├── service/
│   │   ├── user/              # 用户端 Service（新增）
│   │   │   ├── player_service.go
│   │   │   ├── order_service.go
│   │   │   ├── payment_service.go
│   │   │   └── review_service.go
│   │   │
│   │   ├── player/            # 陪玩师端 Service（新增）
│   │   │   ├── profile_service.go
│   │   │   ├── order_service.go
│   │   │   ├── earnings_service.go
│   │   │   └── stats_service.go
│   │   │
│   │   └── admin/             # 管理端（已有）
│   │
│   ├── repository/            # Repository 层（复用 + 扩展）
│   │   ├── player/
│   │   │   └── player_gorm_repository.go  # 扩展查询方法
│   │   ├── order/
│   │   │   └── order_gorm_repository.go   # 扩展查询方法
│   │   ├── withdrawal/        # 新增：提现记录
│   │   │   ├── withdrawal_repository.go
│   │   │   └── withdrawal_gorm_repository.go
│   │   └── ...
│   │
│   └── model/                 # 数据模型（复用 + 扩展）
│       ├── withdrawal.go      # 新增
│       ├── player_tag.go      # 新增
│       ├── notification.go    # 新增（可选）
│       └── ...
│
└── docs/
    └── swagger/               # Swagger 文档
```

### 前端目录结构

```
frontend/
├── src/
│   ├── pages/
│   │   ├── Home/              # 首页（新增）
│   │   │   ├── Home.tsx
│   │   │   ├── Home.module.less
│   │   │   ├── components/
│   │   │   │   ├── HeroSection.tsx
│   │   │   │   ├── HotGames.tsx
│   │   │   │   └── FeaturedPlayers.tsx
│   │   │   └── index.ts
│   │   │
│   │   ├── Players/           # 陪玩师（新增）
│   │   │   ├── PlayerList.tsx
│   │   │   ├── PlayerList.module.less
│   │   │   ├── PlayerDetail.tsx
│   │   │   ├── PlayerDetail.module.less
│   │   │   ├── BookingPage.tsx
│   │   │   ├── components/
│   │   │   │   ├── PlayerCard.tsx
│   │   │   │   ├── FilterSidebar.tsx
│   │   │   │   ├── BookingCard.tsx
│   │   │   │   └── ReviewList.tsx
│   │   │   └── index.ts
│   │   │
│   │   ├── UserCenter/        # 用户中心（新增）
│   │   │   ├── MyOrders/
│   │   │   │   ├── OrderList.tsx
│   │   │   │   ├── OrderDetail.tsx
│   │   │   │   └── components/
│   │   │   ├── MyReviews/
│   │   │   │   └── ReviewList.tsx
│   │   │   ├── Profile/
│   │   │   │   └── ProfilePage.tsx
│   │   │   ├── Wallet/
│   │   │   │   └── WalletPage.tsx
│   │   │   └── index.ts
│   │   │
│   │   ├── Payment/           # 支付（新增）
│   │   │   ├── PaymentPage.tsx
│   │   │   ├── PaymentResult.tsx
│   │   │   └── index.ts
│   │   │
│   │   ├── PlayerCenter/      # 陪玩师中心（新增）
│   │   │   ├── Apply/
│   │   │   │   └── ApplyPage.tsx
│   │   │   ├── Dashboard/
│   │   │   │   └── PlayerDashboard.tsx
│   │   │   ├── Orders/
│   │   │   │   ├── OrderHall.tsx
│   │   │   │   └── MyOrders.tsx
│   │   │   ├── Earnings/
│   │   │   │   ├── EarningsPage.tsx
│   │   │   │   └── WithdrawModal.tsx
│   │   │   ├── Stats/
│   │   │   │   └── StatsPage.tsx
│   │   │   └── index.ts
│   │   │
│   │   └── Login/             # 登录（已有）
│   │       └── ...
│   │
│   ├── components/            # 通用组件
│   │   ├── PlayerCard/        # 陪玩师卡片（新增）
│   │   ├── OrderCard/         # 订单卡片（新增）
│   │   ├── StarRating/        # 星级评分（新增）
│   │   ├── PriceDisplay/      # 价格展示（新增）
│   │   ├── StatusTag/         # 状态标签（新增）
│   │   └── ...
│   │
│   ├── services/              # API 服务
│   │   ├── userApi.ts         # 用户端 API（新增）
│   │   ├── playerApi.ts       # 陪玩师端 API（新增）
│   │   ├── orderApi.ts        # 订单 API（新增）
│   │   ├── paymentApi.ts      # 支付 API（新增）
│   │   └── ...
│   │
│   ├── types/                 # 类型定义
│   │   ├── player.ts          # 陪玩师类型（扩展）
│   │   ├── order.ts           # 订单类型（扩展）
│   │   ├── payment.ts         # 支付类型
│   │   └── ...
│   │
│   ├── hooks/                 # 自定义 Hooks
│   │   ├── usePlayer.ts       # 陪玩师相关（新增）
│   │   ├── useOrder.ts        # 订单相关（新增）
│   │   ├── usePayment.ts      # 支付相关（新增）
│   │   └── ...
│   │
│   ├── layouts/
│   │   ├── UserLayout.tsx     # 用户端布局（新增）
│   │   ├── PlayerLayout.tsx   # 陪玩师端布局（新增）
│   │   └── AdminLayout.tsx    # 管理端布局（已有）
│   │
│   └── router/
│       └── routes.tsx         # 路由配置（扩展）
```

---

## 🎯 第一阶段：用户端基础功能（第1-2周）

### Week 1: 陪玩师模块 + 订单创建

#### 后端任务

**1. 扩展 Player Repository**

文件：`backend/internal/repository/player/player_gorm_repository.go`

```go
// 新增方法
func (r *playerRepository) ListWithFilters(ctx context.Context, opts PlayerListOptions) ([]model.Player, int64, error)
func (r *playerRepository) GetDetailByID(ctx context.Context, id uint64) (*PlayerDetail, error)
func (r *playerRepository) UpdateOnlineStatus(ctx context.Context, playerID uint64, online bool) error
```

**2. 创建 User Service**

文件：`backend/internal/service/user/player_service.go`

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
    cache      cache.Cache
}

func NewPlayerService(...) *PlayerService { }

// 获取陪玩师列表（带筛选）
func (s *PlayerService) ListPlayers(ctx context.Context, req ListPlayersRequest) (*ListPlayersResponse, error) { }

// 获取陪玩师详情
func (s *PlayerService) GetPlayerDetail(ctx context.Context, playerID uint64) (*PlayerDetailResponse, error) { }
```

**3. 创建 User Handler**

文件：`backend/internal/handler/user/player_handler.go`

```go
package user

import (
    "github.com/gin-gonic/gin"
    userservice "gamelink/internal/service/user"
)

type PlayerHandler struct {
    playerService *userservice.PlayerService
}

func NewPlayerHandler(svc *userservice.PlayerService) *PlayerHandler { }

// GET /api/v1/user/players
func (h *PlayerHandler) ListPlayers(c *gin.Context) { }

// GET /api/v1/user/players/:id
func (h *PlayerHandler) GetPlayerDetail(c *gin.Context) { }
```

**4. 注册路由**

文件：`backend/internal/handler/user/router.go`

```go
package user

import (
    "github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router gin.IRouter, handlers *Handlers) {
    user := router.Group("/user")
    
    // 陪玩师相关
    user.GET("/players", handlers.Player.ListPlayers)
    user.GET("/players/:id", handlers.Player.GetPlayerDetail)
    
    // 订单相关
    user.POST("/orders", handlers.Order.CreateOrder)
    user.GET("/orders", handlers.Order.ListMyOrders)
    user.GET("/orders/:id", handlers.Order.GetOrderDetail)
    user.PUT("/orders/:id/cancel", handlers.Order.CancelOrder)
}
```

**5. 创建订单服务**

文件：`backend/internal/service/user/order_service.go`

```go
package user

type OrderService struct {
    orderRepo   repository.OrderRepository
    playerRepo  repository.PlayerRepository
    userRepo    repository.UserRepository
}

// 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, userID uint64, req CreateOrderRequest) (*CreateOrderResponse, error) {
    // 1. 验证陪玩师存在且可接单
    // 2. 验证游戏存在
    // 3. 计算价格
    // 4. 创建订单
    // 5. 返回订单信息
}

// 获取我的订单列表
func (s *OrderService) ListMyOrders(ctx context.Context, userID uint64, req MyOrderListRequest) (*MyOrderListResponse, error) { }

// 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, userID uint64, orderID uint64, reason string) error { }
```

#### 前端任务

**1. 创建 API 服务**

文件：`frontend/src/services/playerApi.ts`

```typescript
import { request } from './request';

export interface PlayerListParams {
  gameId?: number;
  minPrice?: number;
  maxPrice?: number;
  minRating?: number;
  onlineOnly?: boolean;
  sortBy?: 'price' | 'rating' | 'orders';
  page: number;
  pageSize: number;
}

export interface PlayerCardDTO {
  id: number;
  userId: number;
  nickname: string;
  avatarUrl: string;
  bio: string;
  rank: string;
  ratingAverage: number;
  ratingCount: number;
  hourlyRateCents: number;
  mainGame: string;
  isOnline: boolean;
  orderCount: number;
}

export const playerApi = {
  // 获取陪玩师列表
  async getPlayers(params: PlayerListParams) {
    return request<{
      players: PlayerCardDTO[];
      total: number;
    }>('/api/v1/user/players', { params });
  },

  // 获取陪玩师详情
  async getPlayerDetail(id: number) {
    return request<PlayerDetailResponse>(`/api/v1/user/players/${id}`);
  },
};
```

**2. 创建陪玩师卡片组件**

文件：`frontend/src/components/PlayerCard/PlayerCard.tsx`

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
    <Card
      className={styles.playerCard}
      hoverable
      onClick={onClick}
      cover={
        <div className={styles.avatarWrapper}>
          <Avatar size={100} src={player.avatarUrl}>
            {player.nickname[0]}
          </Avatar>
          {player.isOnline && (
            <Tag color="green" className={styles.onlineTag}>
              在线
            </Tag>
          )}
        </div>
      }
    >
      <Card.Meta
        title={player.nickname}
        description={
          <div className={styles.info}>
            <div className={styles.game}>
              {player.mainGame} · {player.rank}
            </div>
            <div className={styles.rating}>
              <Rate readonly value={player.ratingAverage} />
              <span className={styles.count}>
                ({player.ratingCount})
              </span>
            </div>
            <div className={styles.stats}>
              已接单 {player.orderCount} 次
            </div>
            <div className={styles.price}>
              ¥{priceYuan}/小时
            </div>
          </div>
        }
      />
    </Card>
  );
};
```

**3. 创建陪玩师列表页**

文件：`frontend/src/pages/Players/PlayerList.tsx`

```tsx
import React, { useState, useEffect } from 'react';
import { Grid, Pagination, Spin } from '@arco-design/web-react';
import { PlayerCard } from '@/components/PlayerCard';
import { FilterSidebar } from './components/FilterSidebar';
import { playerApi } from '@/services/playerApi';
import { useNavigate } from 'react-router-dom';
import styles from './PlayerList.module.less';

const { Row, Col } = Grid;

export const PlayerList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [players, setPlayers] = useState([]);
  const [total, setTotal] = useState(0);
  const [filters, setFilters] = useState({
    page: 1,
    pageSize: 20,
  });

  const fetchPlayers = async () => {
    setLoading(true);
    try {
      const res = await playerApi.getPlayers(filters);
      setPlayers(res.data.players);
      setTotal(res.data.total);
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPlayers();
  }, [filters]);

  return (
    <div className={styles.playerList}>
      <Row gutter={24}>
        <Col span={6}>
          <FilterSidebar
            filters={filters}
            onChange={(newFilters) => setFilters({ ...filters, ...newFilters })}
          />
        </Col>
        <Col span={18}>
          <Spin loading={loading}>
            <Row gutter={[16, 16]}>
              {players.map((player) => (
                <Col key={player.id} span={8}>
                  <PlayerCard
                    player={player}
                    onClick={() => navigate(`/players/${player.id}`)}
                  />
                </Col>
              ))}
            </Row>
            <div className={styles.pagination}>
              <Pagination
                total={total}
                current={filters.page}
                pageSize={filters.pageSize}
                onChange={(page) => setFilters({ ...filters, page })}
              />
            </div>
          </Spin>
        </Col>
      </Row>
    </div>
  );
};
```

**4. 创建订单 API**

文件：`frontend/src/services/orderApi.ts`

```typescript
export const orderApi = {
  // 创建订单
  async createOrder(data: CreateOrderRequest) {
    return request<CreateOrderResponse>('/api/v1/user/orders', {
      method: 'POST',
      data,
    });
  },

  // 获取我的订单
  async getMyOrders(params: MyOrderListParams) {
    return request<MyOrderListResponse>('/api/v1/user/orders', { params });
  },

  // 获取订单详情
  async getOrderDetail(id: number) {
    return request<OrderDetailResponse>(`/api/v1/user/orders/${id}`);
  },

  // 取消订单
  async cancelOrder(id: number, reason: string) {
    return request(`/api/v1/user/orders/${id}/cancel`, {
      method: 'PUT',
      data: { reason },
    });
  },
};
```

**5. 更新路由**

文件：`frontend/src/router/routes.tsx`

```tsx
import { RouteObject } from 'react-router-dom';
import { PlayerList } from '@/pages/Players/PlayerList';
import { PlayerDetail } from '@/pages/Players/PlayerDetail';
import { UserLayout } from '@/layouts/UserLayout';

export const userRoutes: RouteObject[] = [
  {
    path: '/',
    element: <UserLayout />,
    children: [
      {
        path: 'players',
        element: <PlayerList />,
      },
      {
        path: 'players/:id',
        element: <PlayerDetail />,
      },
      // ... 更多路由
    ],
  },
];
```

#### 测试验收

- [ ] 陪玩师列表展示正常
- [ ] 筛选功能可用
- [ ] 点击卡片跳转详情
- [ ] 陪玩师详情页展示完整
- [ ] 预约下单流程顺畅
- [ ] 订单创建成功

---

### Week 2: 支付 + 评价

#### 后端任务

**1. 支付服务（Mock 版）**

文件：`backend/internal/service/user/payment_service.go`

```go
package user

type PaymentService struct {
    paymentRepo repository.PaymentRepository
    orderRepo   repository.OrderRepository
}

// 创建支付（Mock）
func (s *PaymentService) CreatePayment(ctx context.Context, userID uint64, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
    // 1. 验证订单存在且待支付
    // 2. 创建支付记录
    // 3. 返回支付参数（Mock，实际对接支付SDK）
    // 4. 在测试环境直接标记为已支付
}

// 支付回调（Mock）
func (s *PaymentService) HandlePaymentCallback(ctx context.Context, data map[string]interface{}) error {
    // 1. 验证签名
    // 2. 更新支付状态
    // 3. 更新订单状态
}
```

**2. 评价服务**

文件：`backend/internal/service/user/review_service.go`

```go
package user

type ReviewService struct {
    reviewRepo repository.ReviewRepository
    orderRepo  repository.OrderRepository
    playerRepo repository.PlayerRepository
}

// 创建评价
func (s *ReviewService) CreateReview(ctx context.Context, userID uint64, req CreateReviewRequest) (*CreateReviewResponse, error) {
    // 1. 验证订单已完成且属于该用户
    // 2. 验证订单未被评价
    // 3. 创建评价记录
    // 4. 更新陪玩师评分统计
}

// 获取我的评价
func (s *ReviewService) GetMyReviews(ctx context.Context, userID uint64, page, pageSize int) (*MyReviewListResponse, error) { }
```

#### 前端任务

**1. 支付页面**

文件：`frontend/src/pages/Payment/PaymentPage.tsx`

```tsx
import React, { useState } from 'react';
import { Card, Radio, Button, Message } from '@arco-design/web-react';
import { paymentApi } from '@/services/paymentApi';
import { useParams, useNavigate } from 'react-router-dom';
import styles from './PaymentPage.module.less';

export const PaymentPage: React.FC = () => {
  const { orderId } = useParams();
  const navigate = useNavigate();
  const [method, setMethod] = useState<'wechat' | 'alipay'>('wechat');
  const [loading, setLoading] = useState(false);

  const handlePay = async () => {
    setLoading(true);
    try {
      const res = await paymentApi.createPayment({
        orderId: Number(orderId),
        method,
      });
      
      // Mock: 直接跳转成功页
      navigate(`/payment/result?success=true&orderId=${orderId}`);
    } catch (error) {
      Message.error('支付失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.paymentPage}>
      <Card title="选择支付方式">
        <Radio.Group value={method} onChange={setMethod}>
          <Radio value="wechat">微信支付</Radio>
          <Radio value="alipay">支付宝</Radio>
        </Radio.Group>
        
        <div className={styles.actions}>
          <Button
            type="primary"
            size="large"
            loading={loading}
            onClick={handlePay}
          >
            立即支付
          </Button>
        </div>
      </Card>
    </div>
  );
};
```

**2. 评价组件**

文件：`frontend/src/pages/UserCenter/MyOrders/components/ReviewModal.tsx`

```tsx
import React, { useState } from 'react';
import { Modal, Rate, Input, Message } from '@arco-design/web-react';
import { reviewApi } from '@/services/reviewApi';

interface ReviewModalProps {
  visible: boolean;
  orderId: number;
  onClose: () => void;
  onSuccess: () => void;
}

export const ReviewModal: React.FC<ReviewModalProps> = ({
  visible,
  orderId,
  onClose,
  onSuccess,
}) => {
  const [rating, setRating] = useState(5);
  const [comment, setComment] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    if (rating === 0) {
      Message.warning('请选择评分');
      return;
    }

    setLoading(true);
    try {
      await reviewApi.createReview({
        orderId,
        rating,
        comment,
      });
      Message.success('评价成功');
      onSuccess();
      onClose();
    } catch (error) {
      Message.error('评价失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      title="订单评价"
      visible={visible}
      onOk={handleSubmit}
      onCancel={onClose}
      confirmLoading={loading}
    >
      <div style={{ marginBottom: 16 }}>
        <div style={{ marginBottom: 8 }}>评分</div>
        <Rate value={rating} onChange={setRating} />
      </div>
      <div>
        <div style={{ marginBottom: 8 }}>评价内容（选填）</div>
        <Input.TextArea
          value={comment}
          onChange={setComment}
          placeholder="分享您的体验"
          maxLength={500}
          showWordLimit
        />
      </div>
    </Modal>
  );
};
```

#### 测试验收

- [ ] 支付流程完整
- [ ] Mock 支付成功
- [ ] 订单状态更新
- [ ] 评价功能正常
- [ ] 评分同步到陪玩师

---

## 📦 第一阶段交付物

### 后端
- [x] 陪玩师列表和详情 API
- [x] 订单创建和管理 API
- [x] 支付 API（Mock）
- [x] 评价 API
- [x] Swagger 文档更新

### 前端
- [x] 陪玩师列表页
- [x] 陪玩师详情页
- [x] 预约下单页
- [x] 订单列表页
- [x] 订单详情页
- [x] 支付页面
- [x] 评价功能

### 测试
- [x] 单元测试（核心业务逻辑）
- [x] 集成测试（API 端到端）
- [x] 前端功能测试

---

## 🎯 后续阶段概要

### 第二阶段（第3-4周）
- 陪玩师申请和认证
- 陪玩师接单功能
- 订单状态管理
- 收益统计

### 第三阶段（第5-6周）
- 首页和搜索
- 个人中心完善
- 提现功能
- 数据统计

### 第四阶段（第7周）
- 性能优化
- 移动端适配
- 文档完善
- 上线准备

---

## 📋 开发检查清单

### 每个功能完成后检查

- [ ] 后端代码符合 Go 编码规范
- [ ] 前端代码符合 TypeScript 规范
- [ ] API 文档已更新
- [ ] 单元测试已编写
- [ ] 功能测试已通过
- [ ] 代码已提交并注释清晰
- [ ] 相关文档已更新

---

**下一步**: 开始第一阶段 Week 1 开发


