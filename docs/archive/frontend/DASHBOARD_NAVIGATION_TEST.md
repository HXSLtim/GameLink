# 仪表盘导航功能测试

## 测试目标

验证从仪表盘点击快捷入口和订单状态卡片后，能够正确导航到目标页面并自动应用筛选条件。

## 测试环境

- 前端：React + TypeScript + React Router
- 后端：Go + Gin
- 状态值：前后端完全一致（pending, in_progress, completed, canceled等）

## 测试用例

### 用例1：订单状态卡片 - 待处理

**操作步骤：**
1. 访问仪表盘 `/dashboard`
2. 点击"待处理"订单状态卡片

**预期结果：**
- 导航到 `/orders?status=pending`
- 订单列表页面的状态筛选自动选中"待处理"
- 表格数据自动加载并显示待处理订单
- 分页重置到第1页

### 用例2：订单状态卡片 - 进行中

**操作步骤：**
1. 访问仪表盘 `/dashboard`
2. 点击"进行中"订单状态卡片

**预期结果：**
- 导航到 `/orders?status=in_progress`
- 订单列表页面的状态筛选自动选中"进行中"
- 表格数据自动加载并显示进行中订单

### 用例3：订单状态卡片 - 已完成

**操作步骤：**
1. 访问仪表盘 `/dashboard`
2. 点击"已完成"订单状态卡片

**预期结果：**
- 导航到 `/orders?status=completed`
- 订单列表页面的状态筛选自动选中"已完成"
- 表格数据自动加载并显示已完成订单

### 用例4：订单状态卡片 - 已取消

**操作步骤：**
1. 访问仪表盘 `/dashboard`
2. 点击"已取消"订单状态卡片

**预期结果：**
- 导航到 `/orders?status=canceled`
- 订单列表页面的状态筛选自动选中"已取消"
- 表格数据自动加载并显示已取消订单

### 用例5：快捷入口 - 所有订单

**操作步骤：**
1. 访问仪表盘 `/dashboard`
2. 点击"所有订单"快捷入口

**预期结果：**
- 导航到 `/orders`
- 订单列表页面显示所有状态的订单
- 状态筛选保持默认值（全部状态）

### 用例6：快捷入口 - 待处理订单

**操作步骤：**
1. 访问仪表盘 `/dashboard`
2. 点击"待处理订单"快捷入口

**预期结果：**
- 导航到 `/orders?status=pending`
- 订单列表页面的状态筛选自动选中"待处理"
- 表格数据自动加载并显示待处理订单

### 用例7：快捷入口 - 进行中订单

**操作步骤：**
1. 访问仪表盘 `/dashboard`
2. 点击"进行中订单"快捷入口

**预期结果：**
- 导航到 `/orders?status=in_progress`
- 订单列表页面的状态筛选自动选中"进行中"
- 表格数据自动加载并显示进行中订单

### 用例8：快捷入口 - 用户管理

**操作步骤：**
1. 访问仪表盘 `/dashboard`
2. 点击"用户管理"快捷入口

**预期结果：**
- 导航到 `/users`
- 用户列表页面正常加载
- 显示所有用户

### 用例9：直接URL访问

**操作步骤：**
1. 直接在浏览器地址栏输入 `/orders?status=completed`

**预期结果：**
- 订单列表页面加载
- 状态筛选自动选中"已完成"
- 表格数据显示已完成订单

### 用例10：URL参数变化

**操作步骤：**
1. 访问 `/orders?status=pending`
2. 在地址栏将URL改为 `/orders?status=completed`
3. 按回车

**预期结果：**
- 页面不刷新，但数据自动更新
- 状态筛选变为"已完成"
- 表格数据显示已完成订单

### 用例11：从订单列表返回后再次进入

**操作步骤：**
1. 从仪表盘进入待处理订单 (`/orders?status=pending`)
2. 点击浏览器后退按钮回到仪表盘
3. 再次点击"进行中"订单状态卡片

**预期结果：**
- 导航到 `/orders?status=in_progress`
- 状态筛选正确更新为"进行中"
- 数据正确加载

## 技术验证点

### 前端验证

1. **URL参数读取**
   - 使用 `useSearchParams` 正确读取URL参数
   - `getInitialParams` 函数能从URL提取status
   - 初始状态正确设置

2. **URL参数监听**
   - `useEffect` 监听 `searchParams` 变化
   - URL变化时自动更新 `queryParams.status`
   - 页面重置到第1页

3. **数据加载触发**
   - `loadOrders` 函数在status变化时触发
   - API请求携带正确的status参数
   - 数据正确返回和显示

### 后端验证

1. **参数解析**
   - `buildOrderListOptions` 正确解析status参数
   - `normalizeOrderStatus` 处理字符串格式
   - 支持单个和多个status值

2. **数据查询**
   - Repository层按status筛选数据
   - 返回正确的订单列表
   - 分页信息正确

3. **前后端值匹配**
   - 前端 `OrderStatus.PENDING` = `'pending'`
   - 后端 `OrderStatusPending` = `"pending"`
   - 所有状态值完全一致

## 修复文件清单

1. `frontend/src/utils/urlParams.ts` - 新建URL参数工具
2. `frontend/src/hooks/useListPage.ts` - 增强Hook支持URL参数
3. `frontend/src/pages/Orders/OrderList.tsx` - 添加URL参数读取和监听
4. `frontend/src/pages/Users/UserList.tsx` - 使用新的urlParamKeys配置
5. `frontend/URL_PARAMS_NAVIGATION_FIX.md` - 修复文档
6. `frontend/DASHBOARD_NAVIGATION_TEST.md` - 本测试文档

## 执行测试

### 启动服务

```bash
# 后端
cd backend && make run CMD=user-service

# 前端
cd frontend && npm run dev
```

### 手动测试

1. 访问 http://localhost:5173/dashboard
2. 依次执行上述测试用例
3. 验证每个用例的预期结果

### 自动化测试（未来）

建议添加E2E测试：
- 使用Playwright或Cypress
- 覆盖所有导航场景
- 验证URL和数据状态

## 已知问题

无

## 备注

- 所有状态值使用小写下划线格式（如 `in_progress`）
- 后端支持多个status值，但前端暂时只使用单个
- 未来可以扩展支持更多URL参数（如日期范围、用户ID等）

