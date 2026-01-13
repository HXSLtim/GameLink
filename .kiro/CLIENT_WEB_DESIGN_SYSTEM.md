# GameLink Client Web - Design System Specification

> **Project**: GameLink Desktop PWA Client (User/Player Web Application)
> **Version**: 1.0.0
> **Status**: Design Phase
> **Last Updated**: 2025-01-11

---

## Executive Summary

This design system defines the visual language, interaction patterns, and technical specifications for the GameLink **Desktop PWA** client. The system implements a **dual-theme architecture** inspired by industry-leading gaming platforms:

| Theme | Inspiration | Mood | Target Environment |
|:------|:------------|:-----|:-------------------|
| **Day Mode** | KOOK (开黑啦) | Bright, vibrant, energetic | Daytime usage, well-lit environments |
| **Night Mode** | Discord | Dark, immersive, gaming-focused | Evening gaming, low-light environments |

### Technical Stack

| Category | Technology | Notes |
|:---------|:-----------|:------|
| **Framework** | React 18 + TypeScript | Vite build tool |
| **UI Components** | shadcn/ui | Radix UI primitives + Tailwind CSS |
| **Styling** | Tailwind CSS v4 | CSS variables for theming |
| **State** | Zustand | Lightweight state management |
| **PWA** | vite-plugin-pwa | Service worker, offline support |
| **Target** | Desktop browsers | Chrome, Edge, Firefox, Safari |

---

## Table of Contents

1. [Design Principles](#1-design-principles)
2. [Desktop Layout Architecture](#2-desktop-layout-architecture)
3. [Color System](#3-color-system)
4. [Typography](#4-typography)
5. [Spacing & Layout](#5-spacing--layout)
6. [Component Library](#6-component-library)
7. [PWA Features](#7-pwa-features)
8. [User Experience Patterns](#8-user-experience-patterns)
9. [Accessibility Standards](#9-accessibility-standards)
10. [Responsive Breakpoints](#10-responsive-breakpoints)
11. [Animation & Micro-interactions](#11-animation--micro-interactions)

---

## 2. Desktop Layout Architecture

### Discord/Kook-Style Three-Column Layout

GameLink 桌面端采用类似 Discord 和 Kook 的三栏布局设计，这是游戏社交应用的行业标准布局模式。

```
┌──────┬─────────────┬────────────────────────────────────────────┐
│      │             │                                            │
│ 72px │   240px     │              Main Content                  │
│      │             │              (flex-1)                      │
│Server│  Channel    │                                            │
│ List │   List      │  ┌────────────────────────────────────┐   │
│      │             │  │ Header (48px)                      │   │
│ [🎮] │ ─ 发现 ─    │  ├────────────────────────────────────┤   │
│ [💬] │  热门陪玩师  │  │                                    │   │
│ [👥] │  推荐       │  │ Content Area                       │   │
│ [+]  │  排行榜     │  │ (scrollable)                       │   │
│      │             │  │                                    │   │
│      │ ─ 游戏分类 ─ │  │                                    │   │
│      │  王者荣耀   │  │                                    │   │
│      │  英雄联盟   │  │                                    │   │
│      │  和平精英   │  │                                    │   │
│      │             │  │                                    │   │
│      ├─────────────┤  └────────────────────────────────────┘   │
│ [⚙️] │ User Panel  │                                            │
│      │ (52px)      │                                            │
└──────┴─────────────┴────────────────────────────────────────────┘
```

### Layout Components

#### 1. Server Sidebar (72px)

最左侧的图标导航栏，类似 Discord 的服务器列表。

| 元素 | 说明 | 交互 |
|:-----|:-----|:-----|
| **Home Icon** | 发现陪玩师入口 | 点击跳转首页 |
| **Messages** | 消息中心 | 显示未读红点 |
| **Orders** | 我的订单 | 显示进行中订单数 |
| **Create** | 创建订单 | 快速下单入口 |
| **Settings** | 设置 | 底部固定 |

**视觉规范:**
- 图标尺寸: 48x48px
- 圆角: 默认 24px，hover/active 时 16px
- 激活指示器: 左侧 4px 白色竖条
- 通知红点: 右下角 12px 圆点

#### 2. Channel Sidebar (240px)

中间的频道/分类列表，类似 Discord 的频道列表。

| 区域 | 内容 | 高度 |
|:-----|:-----|:-----|
| **Header** | 应用名称 + 下拉菜单 | 48px |
| **Search** | 搜索框 | 44px |
| **Channel Groups** | 可折叠的分类组 | flex-1 |
| **User Panel** | 用户信息 + 操作按钮 | 52px |

**Channel Group 规范:**
- 组标题: 12px 大写字母，灰色
- 频道项: 32px 高度，左侧图标 + 文字
- 未读计数: 右侧红色徽章

#### 3. Main Content Area (flex-1)

主内容区域，根据当前路由显示不同内容。

| 页面类型 | Header 内容 | Content 布局 |
|:---------|:------------|:-------------|
| **陪玩师列表** | 标题 + 筛选按钮 | 卡片网格/列表 |
| **陪玩师详情** | 返回 + 陪玩师名 | 详情信息 |
| **聊天室** | 群组名 + 成员 | 消息列表 + 输入框 |
| **订单详情** | 订单号 + 状态 | 订单信息 |

### User Panel (底部用户栏)

```
┌─────────────────────────────────────────┐
│ [Avatar] 用户名      [🌙] [🎤] [🎧] [⚙️] │
│          在线状态                        │
└─────────────────────────────────────────┘
```

| 元素 | 功能 |
|:-----|:-----|
| **Avatar** | 用户头像 + 在线状态点 |
| **用户名** | 昵称 + 身份标识 |
| **主题切换** | 日间/夜间模式切换 |
| **麦克风** | 语音设置（预留） |
| **耳机** | 音频设置（预留） |
| **设置** | 用户设置入口 |

### Layout Implementation

```tsx
// MainLayout.tsx
<div className="flex h-screen w-screen overflow-hidden">
  {/* Server Sidebar - 72px */}
  <ServerSidebar />
  
  {/* Channel Sidebar - 240px */}
  <div className="flex w-60 flex-col bg-sidebar">
    <ChannelSidebar />
    <UserPanel />
  </div>
  
  {/* Main Content - flex-1 */}
  <main className="flex flex-1 flex-col overflow-hidden">
    <Outlet />
  </main>
</div>
```

### Responsive Behavior

| 屏幕宽度 | 布局调整 |
|:---------|:---------|
| **≥1280px** | 完整三栏布局 |
| **1024-1279px** | Channel Sidebar 可折叠 |
| **768-1023px** | 双栏布局，Server Sidebar 隐藏 |
| **<768px** | 单栏布局，底部导航 |

### Theme-Specific Styling

**Day Theme (Kook):**
- Server Sidebar: `#F3F4F6`
- Channel Sidebar: `#FFFFFF`
- Main Content: `#F5F7FA`

**Night Theme (Discord):**
- Server Sidebar: `#1E1F22`
- Channel Sidebar: `#2B2D31`
- Main Content: `#313338`

---

## 1. Design Principles

### Core Values

```
┌─────────────────────────────────────────────────────────────────┐
│                    GAMELINK DESIGN DNA                          │
├─────────────────────────────────────────────────────────────────┤
│  🎮 GAMING FIRST     │  Designed by gamers, for gamers        │
│  ⚡ INSTANT GRATIF.  │  <400ms response (Doherty Threshold)    │
│  🎯 CLEAR HIERARCHY  │  Visual flow guides user actions       │
│  ♿ ACCESSIBLE       │  WCAG 2.1 AA compliance minimum        │
│  🌙 ADAPTIVE         │  Seamless theme switching              │
└─────────────────────────────────────────────────────────────────┘
```

### Anti-Patterns (The "Don'ts")

| ❌ Avoid | ✅ Instead |
|:--------|:-----------|
| Hard-coded hex values like `#FF0000` | Semantic tokens like `color-danger-primary` |
| Designing only for desktop | Mobile-first approach |
| Ignoring empty states | Delightful zero-data experiences |
| Theme-breaking colors | Theme-aware component design |
| Infinite scroll memory leaks | Virtualized lists with cleanup |

---

## 3. Color System

### Theme Architecture

```css
/* CSS Variables Approach - enables runtime theme switching */
:root[data-theme="day"] {
  /* Kook-inspired day theme */
}

:root[data-theme="night"] {
  /* Discord-inspired night theme */
}
```

### 2.1 Day Theme (Kook-inspired)

#### Primary Palette

| Token | Hex | HSL | Usage |
|:------|:-----|:-----|:------|
| `--color-primary-50` | #E6FCF6 | hsl(153, 74%, 95%) | Subtle backgrounds |
| `--color-primary-100` | #CBF7E9 | hsl(153, 74%, 88%) | Hover states |
| `--color-primary-500` | **#00D26A** | hsl(165, 100%, 41%) | **Primary CTA, brand** |
| `--color-primary-600` | #00B55C | hsl(165, 100%, 35%) | CTA hover |
| `--color-primary-700` | #00984D | hsl(165, 100%, 30%) | Active states |

#### Semantic Colors

| Token | Hex | Usage |
|:------|:-----|:------|
| `--color-success` | #10B981 | Success messages, confirmations |
| `--color-warning` | #F59E0B | Warnings, pending states |
| `--color-error` | #EF4444 | Errors, destructive actions |
| `--color-info` | #3B82F6 | Information, notifications |

#### Neutral Palette

| Token | Hex | Usage |
|:------|:-----|:------|
| `--color-bg-primary` | #F5F7FA | Main background |
| `--color-bg-card` | #FFFFFF | Card/surface backgrounds |
| `--color-bg-secondary` | #F9FAFB | Alternative backgrounds |
| `--color-text-primary` | #1F2937 | Primary text |
| `--color-text-secondary` | #6B7280 | Secondary text |
| `--color-text-tertiary` | #9CA3AF | Disabled, placeholders |
| `--color-border` | #E5E7EB | Dividers, borders |
| `--color-border-light` | #F3F4F6 | Subtle borders |

### 2.2 Night Theme (Discord-inspired)

#### Primary Palette

| Token | Hex | HSL | Usage |
|:------|:-----|:-----|:------|
| `--color-primary-400` | #7983F5 | hsl(239, 84%, 72%) | Hover states |
| `--color-primary-500` | **#5865F2** | hsl(239, 89%, 65%) | **Primary CTA, brand (Blurple)** |
| `--color-primary-600` | #4752C4 | hsl(239, 63%, 53%) | CTA hover |

#### Semantic Colors

| Token | Hex | Usage |
|:------|:-----|:------|
| `--color-success` | #23A559 | Success messages |
| `--color-warning` | #F0B132 | Warnings |
| `--color-error` | #DA373C | Errors |
| `--color-info` | #00A8FC | Information |

#### Neutral Palette (Discord Dark)

| Token | Hex | Usage |
|:------|:-----|:------|
| `--color-bg-primary` | #313338 | Main background (base) |
| `--color-bg-card` | #2B2D31 | Card/surface backgrounds |
| `--color-bg-elevated` | #2B2D31 | Modals, popovers |
| `--color-bg-deep` | #1E1F22 | Deepest backgrounds |
| `--color-text-primary` | #F2F3F5 | Primary text |
| `--color-text-secondary` | #B5BAC1 | Secondary text |
| `--color-text-tertiary` | #949BA4 | Disabled, hints |
| `--color-border` | #1E1F22 | Dividers, borders |
| `--color-border-subtle` | #2B2D31 | Subtle borders |

### 2.3 Color Usage Guidelines

```
┌─────────────────────────────────────────────────────────────┐
│                    COLOR HIERARCHY                          │
├─────────────────────────────────────────────────────────────┤
│  Level 1 (Brand)   │  Primary-500  │  CTAs, links, brand   │
│  Level 2 (Action)  │  Primary-600  │  Hover, active states  │
│  Level 3 (Status)  │  Success/Error │  Feedback messaging   │
│  Level 4 (Neutral) │  Gray scale   │  Structure, borders   │
└─────────────────────────────────────────────────────────────┘
```

#### Contrast Requirements (WCAG AA)

| Context | Minimum Contrast | Recommended |
|:--------|:-----------------|:------------|
| Normal text | 4.5:1 | 7:1 |
| Large text (18pt+) | 3:1 | 4.5:1 |
| UI components | 3:1 | 4.5:1 |

---

## 3. Typography

### Font Stack

```css
--font-family-base: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI',
                     'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei',
                     sans-serif;

--font-family-mono: 'JetBrains Mono', 'SF Mono', 'Consolas', monospace;
```

### Type Scale (8pt Grid)

| Token | Size | Weight | Line-height | Usage |
|:------|:-----|:-------|:------------|:------|
| `--text-xs` | 12px | 400/500 | 16px | Captions, labels |
| `--text-sm` | 14px | 400/500 | 20px | Secondary text |
| `--text-base` | 16px | 400/500 | 24px | Body text |
| `--text-lg` | 18px | 400/600 | 28px | Emphasized body |
| `--text-xl` | 20px | 600 | 28px | Subheadings |
| `--text-2xl` | 24px | 600 | 32px | Section headings |
| `--text-3xl` | 30px | 700 | 38px | Page titles |
| `--text-4xl` | 36px | 700 | 44px | Hero titles |

### Font Weights

```css
--font-weight-regular: 400;  /* Body text */
--font-weight-medium: 500;   /* Emphasized body */
--font-weight-semibold: 600; /* Headings, labels */
--font-weight-bold: 700;     /* Display titles */
```

### Typography Best Practices

```
✅ DO:
- Use sentence case for UI elements
- Maintain 1.5+ line-height for body text
- Limit line length to 75 characters for readability

❌ DON'T:
- Use all caps (except for acronyms)
- Mix more than 2 weights on one screen
- Use font sizes < 12px for readable content
```

---

## 4. Spacing & Layout

### 8pt Grid System

All spacing values must align to the 8pt grid:

```css
--space-0:   0px;
--space-1:   4px;
--space-2:   8px;
--space-3:   12px;
--space-4:   16px;
--space-5:   20px;
--space-6:   24px;
--space-8:   32px;
--space-10:  40px;
--space-12:  48px;
--space-16:  64px;
--space-20:  80px;
--space-24:  96px;
```

### Common Spacing Patterns

| Context | Spacing | Token |
|:--------|:--------|:------|
| Button padding | 12px 24px | `space-3 space-6` |
| Card padding | 24px | `space-6` |
| Section gap | 48px | `space-12` |
| Form item gap | 16px | `space-4` |
| List item padding | 12px 16px | `space-3 space-4` |

### Container Widths

```css
--container-sm:  640px;   /* Mobile landscape, tablet portrait */
--container-md:  768px;   /* Tablet landscape */
--container-lg:  1024px;  /* Desktop */
--container-xl:  1280px;  /* Large desktop */
```

---

## 5. Component Library

### 5.1 Buttons

#### Primary Button

```
States: default, hover, active, disabled, loading
Padding: 12px 24px
Border-radius: 8px
Font: 600 14px / 20px
```

**Day Theme:**
- Default: `bg-primary-500 text-white`
- Hover: `bg-primary-600 text-white`
- Active: `bg-primary-700 text-white`

**Night Theme:**
- Default: `bg-primary-500 text-white`
- Hover: `bg-primary-600 text-white`
- Active: `scale-0.98`

#### Secondary Button

```
States: default, hover, active, disabled
Padding: 12px 24px
Border-radius: 8px
Font: 600 14px / 20px
```

**Day Theme:**
- Default: `bg-white border-2 border-border text-text-primary`
- Hover: `bg-gray-50 border-text-tertiary`

**Night Theme:**
- Default: `bg-transparent border border-border-subtle text-text-primary`
- Hover: `bg-bg-card border-text-tertiary`

#### Ghost Button

```
For tertiary actions and icon buttons
No background, text-only or icon-only
```

### 5.2 Cards

#### Base Card

```css
.card {
  border-radius: 12px;
  padding: var(--space-6);
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* Night theme enhanced shadow */
[data-theme="night"] .card {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}
```

#### Card Variants

| Variant | Use Case | Styling |
|:--------|:---------|:--------|
| **Elevated** | Modals, popovers | `box-shadow-lg` |
| **Interactive** | Clickable cards | `hover:shadow-md cursor-pointer` |
| **Bordered** | Emphasized sections | `border-2` |

### 5.3 Form Elements

#### Input Fields

```
Height: 44px (minimum touch target)
Padding: 12px 16px
Border-radius: 8px
Border: 1px solid
```

**States:**
- Default: `border-border`
- Focus: `border-primary-500 ring-2 ring-primary-500/20`
- Error: `border-error`
- Disabled: `bg-bg-secondary opacity-50`

#### Labels

```
Font: 600 14px
Color: var(--color-text-primary)
Margin-bottom: 6px
```

### 5.4 Navigation

#### Bottom Navigation (Mobile)

```
Height: 56px + safe-area-inset-bottom
Items: 3-5 icons
Active: Primary color with subtle background
```

#### Top Navigation (Desktop)

```
Height: 64px
Logo left, actions right
Hover states on all interactive elements
```

### 5.5 Avatar System

| Size | Diameter | Use Case |
|:-----|:---------|:---------|
| xs | 24px | Comments, mentions |
| sm | 32px | List items, chat |
| md | 40px | Cards, profiles |
| lg | 64px | Profile page |
| xl | 96px+ | Hero sections |

**Status Indicators:**
- Online: `bg-success` 8px dot
- Away: `bg-warning` 8px dot
- Offline: `bg-text-tertiary` 8px dot
- In-game: `bg-primary-500` 8px dot with ring

### 5.6 Gaming-Specific Components

#### Player Card

```
┌─────────────────────────────────────┐
│  [Avatar]  Player Name    [Rank]    │
│            Status indicator          │
│  ━━━━━━━━⭐⭐⭐⭐☆                 │
│  Games: League, Valorant, Apex      │
│  ¥50/hour  [Book Now]               │
└─────────────────────────────────────┘
```

#### Game Tag

```
Icon + Game Name
Rounded pill: 9999px background
Clickable: filter by game
```

#### Order Status Badge

| Status | Color | Icon |
|:-------|:------|:-----|
| pending | warning | clock |
| confirmed | info | check-circle |
| in_progress | primary-500 | gamepad |
| completed | success | check |
| canceled | error | x-circle |

---

## 6. PWA Features

### 6.1 Service Worker Strategy

```javascript
// Cache strategy: Stale-While-Revalidate
const CACHE_NAME = 'gamelink-v1';
const STALE_DATA = ['player-list', 'game-list'];

// Cache-first for static assets
// Network-first for API calls
// Network-only for real-time data (WebSocket)
```

### 6.2 Install Prompts

**Trigger Conditions:**
- User has visited 2+ times
- User spent 3+ minutes on site
- User engaged with 3+ features

**Custom UI:**
```
┌─────────────────────────────────────────┐
│  🎮 Install GameLink                    │
│  Get app-like experience, offline mode  │
│              [Install] [Not now]        │
└─────────────────────────────────────────┘
```

### 6.3 Offline States

| State | UI Treatment |
|:------|:-------------|
| Offline (cached) | Banner notification, cached content |
| Offline (uncached) | Friendly offline page with retry |
| Reconnecting | Shimmer UI with "Reconnecting..." |
| Online | Brief "You're back!" toast |

### 6.4 Push Notifications

**Notification Types:**
| Type | Priority | Sound |
|:------|:---------|:------|
| New message | High | Default |
| Order accepted | High | Custom chime |
| Payment received | Medium | Ka-ching |
| Promotion | Low | Silent |

**Permission Strategy:**
- Don't prompt on first visit
- Request after first meaningful interaction
- Explain value before requesting

### 6.5 App Shell

```
┌───────────────────────────────────────┐
│  Header (logo, search, profile)       │  ← Cached shell
├──────────┬────────────────────────────┤
│          │                            │
│  Nav     │  Dynamic Content           │  ← Fresh content
│  (side)  │  (cached / network)        │
│          │                            │
└──────────┴────────────────────────────┘
```

---

## 7. User Experience Patterns

### 7.1 Core User Flows

#### Guest → Registration Flow (Progressive Disclosure)

```
Guest browse → Filter/View players → Attempt to book →
  ↓
"Just one step to book your companion!"
  ↓
[Phone/Email input] → [OTP verification] → [Profile completion]
  ↓
"Now booking [Player Name]..."
```

**Rationale:** Minimize friction, capture intent first.

#### Booking Flow

```
Select Player → Choose Game/Time → Confirm Price → Payment →
  ↓
[Real-time] Order accepted → Chat room opens → Session starts
  ↓
Session ends → Rate player → Leave review
```

### 7.2 Empty States (Zero-Data Experience)

| Screen | Empty State | CTA |
|:-------|:------------|:-----|
| Orders | "No orders yet" | "Browse players" |
| Messages | "Start chatting!" | "Find players" |
| Favorites | "Save players for quick access" | "Browse players" |

### 7.3 Loading States

| Type | Pattern | Duration |
|:-----|:--------|:---------|
| Initial | Skeleton screens | <1s |
| Action | Spinner + text | 1-3s |
| Long wait | Progress bar + estimate | 3s+ |
| Optimistic | Immediate update, rollback on error | <400ms |

### 7.4 Error Handling

```
┌─────────────────────────────────────────┐
│  ⚠️ Oops! Something went wrong         │
│                                         │
│  We couldn't load the player list.     │
│  [Retry]  [Go to homepage]             │
└─────────────────────────────────────────┘
```

**Error Recovery:**
- Auto-retry for network errors (exponential backoff)
- Preserve form data on failure
- Provide clear next steps

---

## 8. Accessibility Standards

### 8.1 WCAG 2.1 AA Compliance

| Criterion | Requirement |
|:----------|:------------|
| Color contrast | 4.5:1 minimum |
| Touch targets | 44x44px minimum |
| Keyboard nav | All features accessible |
| Screen reader | Proper ARIA labels |
| Focus indicators | Visible at all times |
| Motion respect | Respect `prefers-reduced-motion` |

### 8.2 Keyboard Navigation

```
Tab: Forward focus
Shift+Tab: Backward focus
Enter/Space: Activate
Escape: Close modals/dropdowns
Arrow keys: Navigate lists, grids
```

### 8.3 Screen Reader Support

```html
<!-- Example: Player card with proper labeling -->
<div role="article" aria-label="Player: JohnDoe">
  <img src="avatar.jpg" alt="JohnDoe's avatar">
  <h2>JohnDoe <span class="sr-only">Player</span></h2>
  <p>4.5 <span aria-label="5 out of 5 stars">⭐⭐⭐⭐☆</span></p>
</div>
```

---

## 9. Responsive Breakpoints

### Breakpoint Scale

```css
/* Mobile-first approach */
--breakpoint-xs:  375px;   /* Small phones */
--breakpoint-sm:  640px;   /* Large phones */
--breakpoint-md:  768px;   /* Tablets */
--breakpoint-lg:  1024px;  /* Laptops */
--breakpoint-xl:  1280px;  /* Desktops */
--breakpoint-2xl: 1536px;  /* Large screens */
```

### Layout Patterns

| Screen Size | Layout Pattern |
|:------------|:---------------|
| < 640px | Single column, bottom nav |
| 640-1024px | Single column, side/top nav |
| > 1024px | Multi-column, persistent side nav |

### Container Queries (Future)

```css
@container (min-width: 400px) {
  .card {
    grid-template-columns: 1fr 1fr;
  }
}
```

---

## 10. Animation & Micro-interactions

### 10.1 Animation Principles

```
1. Purpose-driven: Animations must communicate meaning
2. Subtle: 200-400ms duration, ease-out timing
3. Respectful: Honor prefers-reduced-motion
4. Performant: Use transform, opacity (GPU-accelerated)
```

### 10.2 Standard Animations

| Interaction | Duration | Easing | Property |
|:------------|:---------|:--------|:---------|
| Hover | 150ms | ease-out | opacity, transform |
| Focus | 200ms | ease-out | box-shadow |
| Modal open | 250ms | cubic-bezier(0.4, 0, 0.2, 1) | opacity, transform |
| Page load | 400ms | ease-out | opacity |
| Loading | 1000ms | linear | stroke-dashoffset |

### 10.3 Micro-interactions

#### Like Animation
```
Click → Scale down (0.9) → Scale up (1.1) → Settle (1.0)
Particles burst from center
```

#### Message Sent
```
Send button → Spinner (200ms) → Checkmark (300ms) → Fade
Message slides up + fades in
```

#### Theme Toggle
```
Click → Icon rotates 180° → Colors cross-fade (300ms)
Background slides up/down for transition effect
```

### 10.4 Loading Patterns

#### Shimmer Skeleton

```css
@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

.skeleton {
  background: linear-gradient(
    90deg,
    var(--color-bg-secondary) 0%,
    var(--color-border) 50%,
    var(--color-bg-secondary) 100%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}
```

#### Spinner

```
Two rotating rings, different speeds
Outer ring: 1.5s, linear
Inner ring: 1s, linear, reverse
```

---

## 11. Icon System

### Icon Library

**Primary:** Font Awesome 6 (via CDN)
**Fallback:** SVG icons (custom)

### Icon Sizes (8pt grid)

| Token | Size | Usage |
|:------|:-----|:------|
| `--icon-xs` | 16px | Inline icons, buttons |
| `--icon-sm` | 20px | Small buttons |
| `--icon-md` | 24px | Standard icons |
| `--icon-lg` | 32px | Large icons, featured |
| `--icon-xl` | 48px | Hero icons |

### Common Icons

| Icon | Name | Usage |
|:-----|:-----|:------|
| 🔍 | search | Search players, games |
| 💬 | comment | Chat, messages |
| 🎮 | gamepad | Gaming, sessions |
| ⭐ | star | Ratings, favorites |
| 🔔 | bell | Notifications |
| ⚙️ | cog | Settings |
| 🌙 / ☀️ | moon / sun | Theme toggle |

---

## 12. Brand Assets

### Logo Usage

| Context | Size | Clear Space |
|:--------|:-----|:------------|
| Header (desktop) | 40px height | 16px around |
| Header (mobile) | 32px height | 12px around |
| Footer | 32px height | 8px around |
| Favicon | 32x32px | N/A |

### App Icons (PWA)

| Size | Use Case |
|:-----|:---------|
| 192x192 | Android icon |
| 512x512 | iOS splash screen |
| 1024x1024 | App Store / Play Store |

---

## 13. Implementation Guide

### 13.1 Tech Stack Recommendations

```json
{
  "framework": "Vue 3 / React 18",
  "build": "Vite 5",
  "styling": "Tailwind CSS + CSS Modules",
  "state": "Pinia / Zustand",
  "http": "Axios / Fetch with React Query",
  "pwa": "vite-plugin-pwa",
  "i18n": "vue-i18next / next-i18next",
  "testing": "Vitest + Testing Library"
}
```

### 13.2 CSS Architecture

```css
/* 1. CSS Variables (Theme tokens) */
:root[data-theme="day"] { /* ... */ }
:root[data-theme="night"] { /* ... */ }

/* 2. Utility classes (Tailwind) */
<button class="px-6 py-3 bg-primary-500 text-white rounded-lg">

/* 3. Component styles (CSS Modules) */
.card { /* component-specific styles */ }

/* 4. Responsive utilities */
@media (max-width: 640px) { /* mobile overrides */ }
```

### 13.3 Theme Switching Implementation

```typescript
// Theme store
export const useTheme = create<ThemeState>((set) => ({
  theme: localStorage.getItem('theme') || 'day',
  toggle: () => set((state) => {
    const newTheme = state.theme === 'day' ? 'night' : 'day';
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    return { theme: newTheme };
  }),
}));

// Auto-detect system preference
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)');
if (prefersDark.matches && !localStorage.getItem('theme')) {
  useTheme.getState().toggle();
}
```

---

## 14. Design Tokens Export

### Figma / Design Tool Variables

```json
{
  "GameLink/Day/Primary": "#00D26A",
  "GameLink/Day/Background": "#F5F7FA",
  "GameLink/Night/Primary": "#5865F2",
  "GameLink/Night/Background": "#313338",
  "GameLink/Spacing/Base": "8px",
  "GameLink/Radius/Base": "8px",
  "GameLink/Font/Size/Base": "16px"
}
```

### Tailwind Config Extension

```javascript
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: 'var(--color-primary-500)',
          600: 'var(--color-primary-600)',
        },
        // ... full color palette
      },
      spacing: {
        '18': 'var(--space-18)',
        // ... full spacing scale
      }
    }
  }
}
```

---

## 15. Quality Checklist

### Before Handoff

- [ ] All components have default, hover, active, disabled states
- [ ] All colors pass WCAG AA contrast requirements
- [ ] All text has semantic markup (not just visual styling)
- [ ] All interactive elements have 44x44px minimum touch target
- [ ] All forms have proper labels and error messages
- [ ] All images have alt text or aria-hidden="true"
- [ ] Keyboard navigation works for all features
- [ ] Theme switch preserves all functionality
- [ ] Loading states defined for all async operations
- [ ] Empty states designed for all list views
- [ ] PWA manifest includes all required fields
- [ ] Service worker implements proper caching strategy

---

## 16. Future Enhancements

### Phase 2 (Post-Launch)

| Feature | Description |
|:--------|:-------------|
| Haptic feedback | Vibration on actions (mobile) |
| 3D avatar previews | Three.js player avatar viewer |
| Voice chat UI | In-call voice chat interface |
| Streaming integration | Twitch/Kick embed for player streams |
| AR player cards | WebXR AR preview of player stats |

---

## Appendix A: Color Comparison

```
┌────────────────────────────────────────────────────────────┐
│                    THEME COMPARISON                        │
├────────────────────────────────────────────────────────────┤
│                     │   DAY (Kook)   │  NIGHT (Discord)   │
├─────────────────────┼─────────────────┼────────────────────┤
│ Primary CTA         │   #00D26A (🟢)  │  #5865F2 (🟣)      │
│ Background          │   #F5F7FA (⬜)  │  #313338 (⬛)       │
│ Card                │   #FFFFFF (⬜)  │  #2B2D31 (⬛)       │
│ Text Primary        │   #1F2937 (⬛)  │  #F2F3F5 (⬜)       │
│ Success             │   #10B981 (🟢)  │  #23A559 (🟢)       │
│ Error               │   #EF4444 (🔴)  │  #DA373C (🔴)       │
└─────────────────────┴─────────────────┴────────────────────┘
```

---

## Appendix B: References

### Design Inspiration
- **KOOK (开黑啦)**: https://www.kookapp.cn/
- **Discord**: https://discord.com/
- **LINE**: Typography and color systems
- **Telegram**: Smooth animations and micro-interactions

### Technical Resources
- **PWA Best Practices**: https://web.dev/pwa/
- **WCAG 2.1 Guidelines**: https://www.w3.org/WAI/WCAG21/quickref/
- **Material Design**: https://m3.material.io/
- **Apple HIG**: https://developer.apple.com/design/human-interface-guidelines/

---

**Document Version**: 1.0.0
**Created**: 2025-01-11
**Authors**: Super Dev Team (UI/UX Experts)
**Status**: Ready for Implementation
