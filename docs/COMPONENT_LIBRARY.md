# 🧩 GameLink 组件库设计规范
## Discord/KOOK 风格组件库

**创建日期**: 2025-11-16
**组件库名称**: GameLink UI Kit

---

## 📦 组件清单

### 基础组件 (Foundation)
1. **Button** - 按钮
2. **Input** - 输入框
3. **Avatar** - 头像
4. **Tag** - 标签
5. **Badge** - 徽章
6. **Icon** - 图标
7. **Divider** - 分隔线
8. **Tooltip** - 提示框

### 布局组件 (Layout)
9. **DiscordLayout** - Discord三栏布局
10. **Container** - 容器
11. **Grid** - 网格
12. **Flex** - 弹性布局

### 数据展示 (Display)
13. **Card** - 卡片
14. **PlayerCard** - 陪玩师卡片
15. **OrderCard** - 订单卡片
16. **MessageBubble** - 消息气泡
17. **StatusIndicator** - 状态指示器

### 导航组件 (Navigation)
18. **Sidebar** - 侧边栏
19. **NavBar** - 导航栏
20. **Tabs** - 标签页
21. **Breadcrumb** - 面包屑

### 表单组件 (Form)
22. **Form** - 表单容器
23. **Select** - 下拉选择
24. **Checkbox** - 复选框
25. **Radio** - 单选框
26. **Switch** - 开关
27. **Slider** - 滑块

### 反馈组件 (Feedback)
28. **Toast** - 轻提示
29. **Modal** - 模态框
30. **Drawer** - 抽屉
31. **Loading** - 加载

---

## 🎨 组件详细设计

### 1. Button 按钮

#### 变体 (Variants)

```tsx
// Primary - 主要操作
<Button type="primary" size="large">
  立即下单
</Button>

// Secondary - 次要操作
<Button type="secondary">
  查看详情
</Button>

// Ghost - 幽灵按钮
<Button type="ghost">
  取消
</Button>

// Danger - 危险操作
<Button type="danger">
  删除订单
</Button>

// Link - 文字链接
<Button type="link">
  查看更多
</Button>
```

#### Props接口

```typescript
interface ButtonProps {
  type?: 'primary' | 'secondary' | 'ghost' | 'danger' | 'link';
  size?: 'small' | 'medium' | 'large';
  disabled?: boolean;
  loading?: boolean;
  icon?: ReactNode;
  block?: boolean;  // 块级按钮
  onClick?: () => void;
  children: ReactNode;
}
```

#### 样式规范

```css
/* Primary Button */
.button-primary {
  background: var(--brand-primary);
  color: #ffffff;
  border: none;
  border-radius: var(--radius-md);
  padding: 8px 16px;
  font-size: var(--text-base);
  font-weight: var(--font-medium);
  cursor: pointer;
  transition: all 0.2s ease;
}

.button-primary:hover {
  background: #4752c4;  /* 10% darker */
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.button-primary:active {
  transform: scale(0.95);
}

/* Size variations */
.button-small {
  padding: 4px 12px;
  font-size: var(--text-sm);
  height: 32px;
}

.button-medium {
  padding: 8px 16px;
  font-size: var(--text-base);
  height: 40px;
}

.button-large {
  padding: 12px 24px;
  font-size: var(--text-lg);
  height: 48px;
}
```

---

### 2. Avatar 头像

#### 使用示例

```tsx
// 基础头像
<Avatar src="avatar.jpg" size={48} />

// 带状态指示器
<Avatar
  src="avatar.jpg"
  size={64}
  status="online"  // online | idle | dnd | offline
/>

// 带徽章
<Avatar src="avatar.jpg" size={48}>
  <Badge count={5} />
</Avatar>

// 头像组
<Avatar.Group maxCount={3}>
  <Avatar src="user1.jpg" />
  <Avatar src="user2.jpg" />
  <Avatar src="user3.jpg" />
  <Avatar src="user4.jpg" />
</Avatar.Group>
```

#### Props接口

```typescript
interface AvatarProps {
  src: string;
  alt?: string;
  size?: number | 'small' | 'medium' | 'large';
  status?: 'online' | 'idle' | 'dnd' | 'offline';
  shape?: 'circle' | 'square';
  children?: ReactNode;
}
```

#### 样式规范

```css
.avatar {
  position: relative;
  display: inline-block;
  overflow: hidden;
  border-radius: var(--radius-full);
  background: var(--secondary-bg);
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 状态指示器 */
.avatar-status {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 25%;
  height: 25%;
  border-radius: 50%;
  border: 2px solid var(--primary-bg);
}

.avatar-status.online {
  background: var(--status-online);
}

.avatar-status.idle {
  background: var(--status-idle);
}

.avatar-status.dnd {
  background: var(--status-dnd);
}

.avatar-status.offline {
  background: var(--status-offline);
}
```

---

### 3. PlayerCard 陪玩师卡片

#### 使用示例

```tsx
<PlayerCard
  player={{
    id: 1,
    avatar: 'avatar.jpg',
    name: '陪玩师A',
    status: 'online',
    rating: 4.9,
    reviewCount: 156,
    game: '王者荣耀',
    rank: '荣耀王者',
    price: 50,
    description: '专业上分，态度认真...',
    tags: ['王者荣耀', '陪玩上分']
  }}
  onOrder={() => handleOrder()}
  onViewDetail={() => handleViewDetail()}
/>
```

#### 组件结构

```tsx
const PlayerCard: React.FC<PlayerCardProps> = ({ player, onOrder, onViewDetail }) => {
  return (
    <Card hoverable className="player-card">
      <div className="player-card__header">
        <Avatar
          src={player.avatar}
          size={64}
          status={player.status}
        />
        <div className="player-card__info">
          <h3>{player.name}</h3>
          <div className="player-card__rating">
            <Star filled />
            <span>{player.rating}</span>
            <span className="muted">({player.reviewCount}评价)</span>
          </div>
          <div className="player-card__game">
            {player.game} · {player.rank}
          </div>
        </div>
        <Tag>{player.status}</Tag>
      </div>

      <div className="player-card__body">
        <p className="player-card__description">
          {player.description}
        </p>
        <div className="player-card__tags">
          {player.tags.map(tag => (
            <Tag key={tag}>{tag}</Tag>
          ))}
        </div>
      </div>

      <div className="player-card__footer">
        <div className="player-card__price">
          💰 {player.price}元/小时
        </div>
        <div className="player-card__actions">
          <Button type="primary" onClick={onOrder}>
            立即下单
          </Button>
          <Button type="secondary" onClick={onViewDetail}>
            查看详情
          </Button>
        </div>
      </div>
    </Card>
  );
};
```

#### 样式规范

```css
.player-card {
  background: var(--secondary-bg);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  transition: all 0.3s ease;
}

.player-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.player-card__header {
  display: flex;
  gap: var(--space-3);
  align-items: flex-start;
  margin-bottom: var(--space-4);
}

.player-card__info {
  flex: 1;
}

.player-card__info h3 {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--space-1);
}

.player-card__rating {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  color: var(--text-secondary);
  font-size: var(--text-sm);
}

.player-card__game {
  color: var(--text-muted);
  font-size: var(--text-xs);
  margin-top: var(--space-1);
}

.player-card__description {
  color: var(--text-secondary);
  font-size: var(--text-sm);
  line-height: var(--leading-relaxed);
  margin-bottom: var(--space-3);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.player-card__tags {
  display: flex;
  gap: var(--space-2);
  flex-wrap: wrap;
  margin-bottom: var(--space-4);
}

.player-card__footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--space-4);
  border-top: 1px solid var(--divider);
}

.player-card__price {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--brand-primary);
}

.player-card__actions {
  display: flex;
  gap: var(--space-2);
}
```

---

### 4. DiscordLayout Discord布局

#### 使用示例

```tsx
<DiscordLayout>
  <DiscordLayout.Sidebar>
    <Logo />
    <NavItem icon={<HomeIcon />} label="首页" />
    <NavItem icon={<GameIcon />} label="陪玩师" />
    <NavItem icon={<OrderIcon />} label="订单" />
    <NavItem icon={<MessageIcon />} label="消息" badge={5} />
    <NavItem icon={<UserIcon />} label="我的" />
  </DiscordLayout.Sidebar>

  <DiscordLayout.Main>
    <DiscordLayout.Header>
      <SearchBar />
      <UserMenu />
    </DiscordLayout.Header>
    <DiscordLayout.Content>
      {children}
    </DiscordLayout.Content>
  </DiscordLayout.Main>

  <DiscordLayout.Panel>
    <FilterPanel />
  </DiscordLayout.Panel>
</DiscordLayout>
```

#### 组件结构

```tsx
const DiscordLayout: React.FC<DiscordLayoutProps> = ({ children }) => {
  return (
    <div className="discord-layout">
      {children}
    </div>
  );
};

DiscordLayout.Sidebar = ({ children }) => (
  <aside className="discord-layout__sidebar">
    {children}
  </aside>
);

DiscordLayout.Main = ({ children }) => (
  <main className="discord-layout__main">
    {children}
  </main>
);

DiscordLayout.Panel = ({ children }) => (
  <aside className="discord-layout__panel">
    {children}
  </aside>
);
```

#### 样式规范

```css
.discord-layout {
  display: flex;
  height: 100vh;
  background: var(--primary-bg);
}

/* 左侧导航栏 */
.discord-layout__sidebar {
  width: 72px;
  background: var(--tertiary-bg);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--space-3) 0;
  gap: var(--space-2);
  overflow-y: auto;
}

/* 主要内容区域 */
.discord-layout__main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 头部 */
.discord-layout__header {
  height: 64px;
  background: var(--secondary-bg);
  border-bottom: 1px solid var(--divider);
  display: flex;
  align-items: center;
  padding: 0 var(--space-4);
  justify-content: space-between;
}

/* 内容区域 */
.discord-layout__content {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-6);
}

/* 右侧面板 */
.discord-layout__panel {
  width: 280px;
  background: var(--secondary-bg);
  border-left: 1px solid var(--divider);
  padding: var(--space-4);
  overflow-y: auto;
}

/* 响应式 */
@media (max-width: 768px) {
  .discord-layout__sidebar {
    display: none;
  }

  .discord-layout__panel {
    position: fixed;
    top: 0;
    right: -280px;
    height: 100vh;
    transition: right 0.3s ease;
    z-index: 1000;
  }

  .discord-layout__panel.open {
    right: 0;
  }
}
```

---

### 5. MessageBubble 消息气泡

#### 使用示例

```tsx
// 对方的消息
<MessageBubble
  type="received"
  avatar="avatar.jpg"
  username="陪玩师A"
  timestamp="10:30"
  content="你好，什么时候可以开始？"
/>

// 自己的消息
<MessageBubble
  type="sent"
  timestamp="10:35"
  content="现在可以"
/>

// 系统消息
<MessageBubble
  type="system"
  content="订单已确认"
/>
```

#### 组件结构

```tsx
const MessageBubble: React.FC<MessageBubbleProps> = ({
  type,
  avatar,
  username,
  timestamp,
  content,
}) => {
  return (
    <div className={`message-bubble message-bubble--${type}`}>
      {type === 'received' && (
        <Avatar src={avatar} size={40} />
      )}
      <div className="message-bubble__content">
        {type === 'received' && (
          <div className="message-bubble__header">
            <span className="message-bubble__username">{username}</span>
            <span className="message-bubble__timestamp">{timestamp}</span>
          </div>
        )}
        <div className="message-bubble__text">
          {content}
        </div>
        {type === 'sent' && (
          <span className="message-bubble__timestamp">{timestamp}</span>
        )}
      </div>
    </div>
  );
};
```

#### 样式规范

```css
.message-bubble {
  display: flex;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

/* 接收的消息 */
.message-bubble--received {
  justify-content: flex-start;
}

.message-bubble--received .message-bubble__text {
  background: var(--secondary-bg);
  color: var(--text-primary);
  border-radius: 0 var(--radius-lg) var(--radius-lg) var(--radius-lg);
}

/* 发送的消息 */
.message-bubble--sent {
  justify-content: flex-end;
}

.message-bubble--sent .message-bubble__text {
  background: var(--brand-primary);
  color: #ffffff;
  border-radius: var(--radius-lg) 0 var(--radius-lg) var(--radius-lg);
}

/* 系统消息 */
.message-bubble--system {
  justify-content: center;
}

.message-bubble--system .message-bubble__text {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-muted);
  font-size: var(--text-xs);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-full);
}

/* 消息内容 */
.message-bubble__content {
  max-width: 70%;
}

.message-bubble__header {
  display: flex;
  gap: var(--space-2);
  margin-bottom: var(--space-1);
}

.message-bubble__username {
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  font-size: var(--text-sm);
}

.message-bubble__timestamp {
  color: var(--text-muted);
  font-size: var(--text-xs);
}

.message-bubble__text {
  padding: var(--space-3) var(--space-4);
  line-height: var(--leading-relaxed);
  word-wrap: break-word;
}
```

---

### 6. StatusIndicator 状态指示器

#### 使用示例

```tsx
<StatusIndicator status="online" />
<StatusIndicator status="idle" />
<StatusIndicator status="dnd" label="游戏中" />
<StatusIndicator status="offline" />
```

#### 样式规范

```css
.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
}

.status-indicator__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.status-indicator__dot.online {
  background: var(--status-online);
}

.status-indicator__dot.idle {
  background: var(--status-idle);
}

.status-indicator__dot.dnd {
  background: var(--status-dnd);
}

.status-indicator__dot.offline {
  background: var(--status-offline);
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
```

---

## 🎨 主题配置

### 主题Provider

```tsx
import { ThemeProvider } from './theme';

function App() {
  return (
    <ThemeProvider theme="dark">
      <YourApp />
    </ThemeProvider>
  );
}
```

### 主题切换

```tsx
const { theme, setTheme } = useTheme();

<Switch
  checked={theme === 'dark'}
  onChange={(checked) => setTheme(checked ? 'dark' : 'light')}
/>
```

---

## 📚 使用指南

### 安装

```bash
npm install @gamelink/ui-kit
```

### 导入

```tsx
import {
  Button,
  Avatar,
  PlayerCard,
  DiscordLayout
} from '@gamelink/ui-kit';
import '@gamelink/ui-kit/dist/styles.css';
```

### 全局样式

```tsx
import { GlobalStyles } from '@gamelink/ui-kit';

function App() {
  return (
    <>
      <GlobalStyles />
      <YourApp />
    </>
  );
}
```

---

## 🧪 组件测试

### 测试覆盖率要求

- 单元测试覆盖率 ≥ 80%
- 组件快照测试
- 交互测试
- 可访问性测试

---

**组件库设计完成，包含30+个核心组件！** 🎨
