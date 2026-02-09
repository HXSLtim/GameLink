# GameLink 前端重构规范（UI 布局与组件拆分）

**适用范围**
- 用户端/陪玩端：`app/`（uni-app + Vue 3 + TS）
- 管理端：`admin/`（React + TS + Vite + Ant Design）

**当前执行范围**
- 仅 App 端执行重构（Admin 规范作为预留，后续统一落地）

**目标**
- 统一 UI 布局与视觉层级
- 明确组件划分与职责边界
- 降低页面复杂度，便于重构与长期维护

---

## 1) 统一布局规范

### 1.1 页面布局结构（通用）
每个页面必须拆成 **「页面容器」→「布局区块」→「业务组件」** 三层：

- **页面容器（Page Container）**  
  负责：路由参数、数据请求、状态管理（加载/错误/空态）、页面级布局。
- **布局区块（Section / Panel）**  
  负责：页面布局分区（筛选区、列表区、统计区等）。
- **业务组件（Business Components）**  
  负责：展示逻辑与交互，禁止直接请求 API。

### 1.2 统一间距与尺寸
**禁止硬编码颜色/间距/圆角**，必须使用设计令牌：

- App：`app/src/styles/variables.scss`、`app/src/styles/mixins.scss`
- Admin：`admin/src/theme/*`（`spacing`、`semanticColors`、`borderRadius`）

### 1.3 统一页面外层布局
- App：推荐统一使用布局组件（如 `layouts/ResponsiveLayout.vue`）作为外层容器  
  - 仅当性能或交互冲突时可跳过，但需要注释说明理由
- Admin：页面统一使用 `PageContainer` 或 `layouts/*` 的结构，避免散乱的 `div` 包裹

---

## 2) 组件拆分规范（核心）

### 2.1 组件层级定义
| 类型 | 位置 | 责任 | 规则 |
|------|------|------|------|
| 基础组件 | `components/gl`（App） / `components`（Admin） | 纯展示/交互，不含业务 | 可跨页面复用 |
| 业务组件 | `pages/**/components` | 单一业务块 | 不直接请求 API |
| 页面容器 | `pages/**/index.*` | 数据请求/状态管理 | 只做组合，不做复杂 UI |

### 2.2 何时必须拆分
满足任意一项则必须拆分组件：
- 单文件超过 **300 行**（不含样式）
- 同一页面出现 **3+ 业务区块**（例如：筛选区、列表区、详情区）
- 页面内出现 **重复 UI/交互逻辑 ≥ 2 次**
- 页面内存在 **复杂状态/逻辑**（如分页+筛选+排序组合）

### 2.3 组件拆分流程（推荐）
1. **先拆布局**：把页面区块拆成组件
2. **再拆逻辑**：将数据处理逻辑放入 composable/hook
3. **最后统一样式**：用设计令牌替换硬编码

---

## 3) App（uni-app）规范

### 3.1 目录结构（目标）
```
app/src/pages/feature/name/
├── index.vue                # 页面容器
├── components/              # 页面内组件
│   ├── FilterBar/index.vue
│   ├── ResultList/index.vue
│   └── ResultCard/index.vue
├── composables/             # 页面逻辑
│   └── useFeatureList.ts
├── types.ts                 # 页面专用类型
└── constants.ts             # 页面常量/枚举
```

### 3.2 Vue 组件规范
- 使用 `<script setup lang="ts">`
- 统一 `defineProps` / `defineEmits` 使用方式  
  推荐：
  ```ts
  const props = withDefaults(defineProps<Props>(), {
    size: 'md',
  })
  ```
- 组件对外只暴露 **明确的 props/emit**，禁止依赖页面内全局变量

### 3.3 样式规范（App）
- 必须使用 `variables.scss` 中的 token
- 组件内样式 **默认 `scoped`**
- 不允许直接写 `rpx` 魔法数（除非是 token 或 mixin）
- 使用 `mixins.scss` 的响应式 mixin 替代手写 media query

---

## 4) Admin（React）规范

### 4.1 目录结构（目标）
```
admin/src/pages/Feature/
├── index.tsx                # 页面容器
├── components/              # 页面内组件
│   ├── FilterBar.tsx
│   ├── ResultTable.tsx
│   └── DetailDrawer.tsx
├── hooks/                   # 页面逻辑
│   └── useFeatureQuery.ts
├── types.ts                 # 页面专用类型
└── constants.ts
```

### 4.2 React 组件规范
- 函数组件 + Hooks
- Props 类型使用 `type` 或 `interface`，命名 `XXXProps`
- 组件文件名统一 PascalCase
- 禁止在组件内直接调用 `axios`，统一通过 `api/*` 模块

### 4.3 样式规范（Admin）
- 统一使用 **CSS Modules / Less**
- 禁止硬编码颜色与间距，使用 `theme/*` 中 tokens
- 临时内联样式必须注明 TODO 并限期替换

---

## 5) API / 状态 / 错误处理规范

### 5.1 API 调用统一入口
- App：只允许通过 `app/src/api/*` 和 `api/request.ts`
- Admin：只允许通过 `admin/src/api/*` 和 `api/client.ts`

### 5.2 状态管理规则
- **局部状态** → 组件内部
- **页面级状态** → composable/hook
- **跨页面状态** → Store（Pinia / Zustand）

### 5.3 错误与空态统一
每个页面至少包含以下状态：
- `loading`
- `error`
- `empty`
- `content`

如无现成组件，应优先补齐基础组件后复用。

---

## 6) 组件重构验收标准（DoD）

**每个被重构页面必须满足：**
- 页面文件 ≤ 200 行（不含样式）
- 数据请求逻辑已迁移到 composable/hook
- 页面使用统一布局容器
- 所有颜色/间距来自设计令牌
- 组件职责清晰（UI vs 数据 vs 状态分离）
- 抽出的业务组件具备可复用价值

---

## 7) 推荐的重构切分模板

以「列表页」为例，拆分为：

1. `FilterBar`（筛选区）
2. `SortTabs`（排序切换）
3. `ResultList`（列表容器）
4. `ResultCard`（单卡片）
5. `PageState` / `GlEmpty` / `Skeleton`
6. `useListQuery`（分页/筛选/排序逻辑）

---

## 8) 执行顺序建议

1. **建立基础组件目录与规范**
2. **统一布局与间距 token**
3. **拆分高频页面（列表/详情/表单）**
4. **清理重复样式与重复业务组件**

---

## 9) 快速检查清单

- 是否复用已有基础组件？
- 是否存在重复 UI？
- 是否存在页面直写 API 调用？
- 是否使用了设计 token？
- 页面逻辑是否已迁移到 composable/hook？
- 组件命名/目录是否符合规范？

---

## 10) 与现有文件的对应关系

- App 设计 token：`app/src/styles/variables.scss`
- App Mixins：`app/src/styles/mixins.scss`
- Admin 设计 token：`admin/src/theme/*`
- API 统一入口：`app/src/api/*`、`admin/src/api/*`

---

**本规范用于指导分组件重构执行。后续如落地问题或需要补齐基础组件，请在本文件追加“补充规范”小节。**
