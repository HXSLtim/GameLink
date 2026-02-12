# GameLink 系统架构图

> **版本**: v2.0
> **更新日期**: 2026-02-11
> **项目**: GameLink 游戏陪玩社交平台

---

## 目录

1. [系统整体架构](#1-系统整体架构)
2. [后端分层架构](#2-后端分层架构)
3. [前端组件架构](#3-前端组件架构)
4. [数据库关系图](#4-数据库关系图)
5. [API路由结构](#5-api路由结构)
6. [WebSocket消息流](#6-websocket消息流)
7. [部署架构](#7-部署架构)

---

## 1. 系统整体架构

```mermaid
graph TB
    subgraph "客户端层"
        A1[小程序/H5<br/>uni-app + Vue 3]
        A2[管理后台<br/>React 19 + AntD]
        A3[未来: APP<br/>uni-app编译]
    end

    subgraph "接入层"
        B1[Nginx<br/>反向代理 + 负载均衡]
    end

    subgraph "API层 - Go 1.24"
        C1[Handler层<br/>路由处理]
        C2[Service层<br/>业务逻辑]
        C3[Repository层<br/>数据访问]
        C4[WebSocket<br/>实时通讯]
        C5[Scheduler<br/>定时任务]
        C6[Middleware<br/>中间件]
    end

    subgraph "数据层"
        D1[(PostgreSQL 16+<br/>主存储)]
        D2[(Redis 7+<br/>缓存/会话)]
    end

    subgraph "外部服务"
        E1[微信支付/登录]
        E2[支付宝]
        E3[腾讯云TRTC]
        E4[OSS对象存储]
        E5[短信服务]
    end

    A1 -->|HTTPS/WSS| B1
    A2 -->|HTTPS/WSS| B1
    A3 -->|HTTPS/WSS| B1

    B1 --> C1
    C1 --> C2
    C2 --> C3
    C2 --> D2
    C3 --> D1

    C1 --> C4
    C1 --> C5
    C1 -.->|认证| C6

    C2 --> E1
    C2 --> E2
    C2 --> E3
    C2 --> E4
    C2 --> E5

    style A1 fill:#e1f5ff
    style A2 fill:#e1f5ff
    style C1 fill:#fff4e6
    style C2 fill:#fff4e6
    style C3 fill:#fff4e6
    style D1 fill:#e8f5e9
    style D2 fill:#e8f5e9
```

---

## 2. 后端分层架构

```mermaid
graph TB
    subgraph "Handler层 - 路由处理"
        A1[admin/ - 管理端路由]
        A2[user/ - 用户端路由]
        A3[player/ - 陪玩师端路由]
        A4[public/ - 公开接口]
    end

    subgraph "Service层 - 业务逻辑 (57个模块)"
        B1[订单服务<br/>orderService]
        B2[支付服务<br/>paymentService]
        B3[钱包服务<br/>walletService]
        B4[聊天服务<br/>chatService]
        B5[认证服务<br/>authService]
        B6[其他52个服务模块...]
    end

    subgraph "Repository层 - 数据访问 (56个模块)"
        C1[订单仓储<br/>orderRepo]
        C2[支付仓储<br/>paymentRepo]
        C3[用户仓储<br/>userRepo]
        C4[其他53个仓储...]
    end

    subgraph "数据模型 (67个)"
        D1[User, Player, Order]
        D2[Payment, Wallet, Chat]
        D3[其他61个模型...]
    end

    A1 --> B1
    A2 --> B1
    A3 --> B1
    A4 --> B1

    B1 --> C1
    B2 --> C2
    B3 --> C3
    B4 --> C4

    C1 --> D1
    C2 --> D2
    C3 --> D1
    C4 --> D3

    style A1 fill:#ffebee
    style A2 fill:#e8f5e9
    style A3 fill:#e3f2fd
    style A4 fill:#fff3e0
    style B1 fill:#f3e5f5
    style B2 fill:#f3e5f5
    style C1 fill:#e0f2f1
    style C2 fill:#e0f2f1
```

---

## 3. 前端组件架构

### 3.1 管理后台 (React 19)

```mermaid
graph TB
    subgraph "页面层 - 40+页面"
        A1[Dashboard]
        A2[User Management]
        A3[Order Management]
        A4[其他37个页面...]
    end

    subgraph "业务组件"
        B1[UserTable]
        B2[OrderCard]
        B3[PlayerCard]
        B4[其他业务组件...]
    end

    subgraph "通用组件"
        C1[DataTable]
        C2[SearchBar]
        C3[Modal]
        C4[Form]
        C5[其他通用组件...]
    end

    subgraph "基础组件 (Ant Design 6)"
        D1[Button, Input, Select]
        D2[Table, Modal, Form]
        D3[其他AntD组件...]
    end

    subgraph "状态管理"
        E1[Zustand Store]
        E2[React Context]
    end

    subgraph "API层"
        F1[apiClient]
        F2[api/admin.ts]
        F3[api/auth.ts]
        F4[其他API模块...]
    end

    A1 --> B1
    A2 --> B1
    A3 --> B2

    B1 --> C1
    B2 --> C2
    B3 --> C1

    C1 --> D1
    C2 --> D2
    C3 --> D2

    A1 -.->|读取| E1
    A1 -.->|请求| F1

    style A1 fill:#e3f2fd
    style B1 fill:#fff3e0
    style C1 fill:#f1f8e9
    style E1 fill:#fce4ec
```

### 3.2 移动端 (uni-app + Vue 3)

```mermaid
graph TB
    subgraph "页面层 - 28个页面"
        A1[index - 首页]
        A2[player/list - 陪玩师列表]
        A3[order/create - 创建订单]
        A4[其他25个页面...]
    end

    subgraph "业务组件 (133个)"
        B1[PlayerCard]
        B2[OrderCard]
        C1[MessageItem]
        C2[ReviewCard]
        D1[其他业务组件...]
    end

    subgraph "模式组件"
        E1[SearchBar]
        E2[FilterPanel]
        E3[InfiniteList]
        E4[PageState]
    end

    subgraph "基础组件 (gl/)"
        F1[GlButton]
        F2[GlInput]
        F3[GlCard]
        F4[GlTag]
        F5[其他基础组件...]
    end

    subgraph "Composables (38个)"
        G1[usePagination]
        G2[useWebSocket]
        G3[useAuth]
        G4[其他Hook...]
    end

    subgraph "Store (Pinia)"
        H1[userStore]
        H2[chatStore]
        H3[configStore]
    end

    subgraph "API层 (14个)"
        I1[request.ts]
        I2[auth.ts]
        I3[order.ts]
        I4[其他API模块...]
    end

    A1 --> B1
    A2 --> B1
    A3 --> B2

    B1 --> E1
    B2 --> E2
    C1 --> E3

    E1 --> F1
    E2 --> F2

    A1 -.->|调用| G1
    A1 -.->|读取| H1
    A1 -.->|请求| I1

    style A1 fill:#e1f5ff
    style B1 fill:#fff9c4
    style E1 fill:#f3e5f5
    style F1 fill:#e8f5e9
    style G1 fill:#ffe0b2
```

---

## 4. 数据库关系图

### 4.1 核心业务表关系

```mermaid
erDiagram
    users ||--o{ players : "是陪玩师"
    users ||--|| wallets : "拥有"
    users ||--o{ orders : "创建"
    users ||--o{ chat_group_members : "加入"

    players ||--o{ player_services : "提供"
    players ||--o{ player_certifications : "认证"
    players ||--o{ orders : "服务"
    players ||--o{ player_presence : "状态"

    orders ||--|| order_groups : "属于"
    orders ||--o{ payments : "支付"
    orders ||--o{ reviews : "评价"
    orders ||--o{ chat_groups : "聊天"
    orders ||--o{ order_disputes : "争议"

    payments ||--o{ refund_records : "退款"

    users ||--o{ withdrawals : "提现"
    players ||--o{ commission_records : "佣金"

    users {
        number id PK
        string email
        string phone
        string role
        number vip_level_id
    }

    players {
        number id PK
        number user_id FK
        string nickname
        number rating_average
        string verification_status
    }

    orders {
        number id PK
        string order_no
        number user_id FK
        number player_id FK
        number total_price_cents
        string status
    }

    payments {
        number id PK
        number order_id FK
        string payment_method
        number amount_cents
        string status
    }

    wallets {
        number id PK
        number user_id FK
        number balance_cents
        number frozen_cents
    }
```

### 4.2 表分类统计

| 分类 | 表数量 | 代表表 |
|------|--------|--------|
| 用户与陪玩师 | 7 | users, players, wallets, player_certifications |
| 订单与支付 | 6 | orders, payments, refund_records, order_disputes |
| 服务与商品 | 3 | service_items, games, game_ranks |
| 聊天与通讯 | 4 | chat_groups, chat_messages, chat_snapshots |
| 营销与会员 | 6 | vip_levels, coupons, activities, referrals |
| 财务与结算 | 5 | commission_records, settlements, withdraws |
| 内容与审核 | 4 | reviews, feeds, sensitive_words |
| 权限与管理 | 5 | roles, permissions, menus, user_roles |
| 团队与社交 | 4 | teams, team_members, lfg_requests, favorites |

---

## 5. API路由结构

```
/api/v1
├── /auth              # 认证相关
│   ├── POST /register
│   ├── POST /login
│   ├── POST /logout
│   ├── POST /refresh
│   └── GET  /me
│
├── /public            # 公开接口
│   ├── GET  /players
│   ├── GET  /games
│   └── GET  /game-categories
│
├── /user              # 用户端接口
│   ├── GET  /orders
│   ├── POST /orders
│   ├── GET  /wallet
│   ├── POST /recharge
│   └── GET  /chats
│
├── /player            # 陪玩师端接口
│   ├── GET  /orders (可接订单)
│   ├── POST /orders/:id/accept
│   ├── PUT  /orders/:id/complete
│   ├── GET  /services
│   └── POST /services
│
└── /admin             # 管理端接口
    ├── /users         # 用户管理
    ├── /players       # 陪玩师管理
    ├── /orders        # 订单管理
    ├── /games         # 游戏管理
    ├── /roles         # 角色管理
    ├── /permissions   # 权限管理
    ├── /menus         # 菜单管理
    ├── /dashboard     # 仪表盘
    ├── /settlements   # 结算管理
    ├── /commissions   # 佣金管理
    ├── /withdraws     # 提现管理
    └── ...           # 其他管理接口

/ws                   # WebSocket连接
```

---

## 6. WebSocket消息流

```mermaid
sequenceDiagram
    participant C as 客户端
    participant W as WebSocket Server
    participant S as 业务服务
    participant R as Redis Pub/Sub

    C->>W: 1. 连接请求 (带Token)
    W->>S: 2. 验证Token
    S-->>W: 3. 验证成功
    W-->>C: 4. 连接建立

    loop 心跳机制
        C->>W: ping (每30秒)
        W-->>C: pong
    end

    par 系统监控推送
        W-->>C: system_status
    and 在线用户推送
        W-->>C: online_users
    and 订单队列推送
        W-->>C: order_queue
    and 告警推送
        W-->>C: alert
    end

    C->>W: 订阅频道 (如: admin)
    W->>R: 订阅Redis频道

    R->>W: 发布消息
    W-->>C: 推送消息

    C->>W: 业务消息 (如: 聊天)
    W->>S: 处理业务
    S->>R: 发布给其他用户
```

### WebSocket消息类型

| 类型 | 方向 | 说明 |
|------|------|------|
| `ping` | 客户端→服务端 | 心跳 |
| `pong` | 服务端→客户端 | 心跳响应 |
| `system_status` | 服务端→客户端 | 系统状态 |
| `online_users` | 服务端→客户端 | 在线用户数 |
| `order_queue` | 服务端→客户端 | 订单队列 |
| `alert` | 服务端→客户端 | 告警消息 |
| `presence_update` | 双向 | 在线状态更新 |
| `room_created` | 服务端→客户端 | 房间创建 |
| `room_member_joined` | 服务端→客户端 | 成员加入 |
| `chat_message` | 双向 | 聊天消息 |

---

## 7. 部署架构

```mermaid
graph TB
    subgraph "用户端"
        U1[用户浏览器<br/>小程序]
    end

    subgraph "CDN"
        CDN[静态资源CDN]
    end

    subgraph "负载均衡"
        LB[Nginx<br/>负载均衡]
    end

    subgraph "应用服务器集群"
        APP1[API服务器 1<br/>Go 1.24]
        APP2[API服务器 2<br/>Go 1.24]
        APP3[API服务器 N<br/>Go 1.24]
    end

    subgraph "数据层"
        DB[(PostgreSQL<br/>主从复制)]
        RDS[(Redis<br/>集群模式)]
    end

    subgraph "监控"
        M1[Prometheus]
        M2[Grafana]
        M3[Alertmanager]
    end

    U1 --> CDN
    U1 --> LB

    LB --> APP1
    LB --> APP2
    LB --> APP3

    APP1 --> DB
    APP2 --> DB
    APP3 --> DB

    APP1 --> RDS
    APP2 --> RDS
    APP3 --> RDS

    APP1 -.->|指标| M1
    APP2 -.->|指标| M1
    APP3 -.->|指标| M1

    M1 --> M2
    M1 --> M3

    style U1 fill:#e1f5ff
    style APP1 fill:#fff4e6
    style APP2 fill:#fff4e6
    style APP3 fill:#fff4e6
    style DB fill:#e8f5e9
    style RDS fill:#e8f5e9
```

---

## 附录：技术栈版本

| 层级 | 技术 | 版本 |
|------|------|------|
| 后端语言 | Go | 1.24.5 |
| Web框架 | Gin | 最新 |
| ORM | GORM | 最新 |
| 数据库 | PostgreSQL | 16+ |
| 缓存 | Redis | 7+ |
| 前端框架 | React | 19 |
| UI组件 | Ant Design | 6 |
| 状态管理 | Zustand | 最新 |
| 移动端框架 | uni-app | 最新 |
| 移动端框架 | Vue | 3.4+ |
| 容器化 | Docker | 最新 |
| CI/CD | GitHub Actions | - |

---

**文档维护**: 产品经理团队
**最后更新**: 2026-02-11
