# E2E 测试说明

## 运行 E2E 测试

### 方法 1: 手动启动 Dev Server

```bash
# Terminal 1: 启动 dev server
cd client
npm run dev

# Terminal 2: 运行 E2E 测试
cd client
npm run test:e2e
```

### 方法 2: 使用已部署的环境

```bash
# 设置环境变量指向已部署的环境
BASE_URL=https://staging.gamelink.com npm run test:e2e
```

## 可用命令

| 命令 | 说明 |
|------|------|
| `npm run test:e2e` | 运行所有 E2E 测试（无头模式） |
| `npm run test:e2e:ui` | UI 模式运行测试 |
| `npm run test:e2e:headed` | 有头浏览器模式运行 |
| `npm run test:e2e:debug` | 调试模式 |
| `npm run test:e2e:report` | 查看测试报告 |

## 测试状态

| 测试套件 | 状态 | 通过/总计 |
|---------|------|----------|
| Client Authentication | ✅ 完全通过 | 10/10 |
| Player Browsing | ⚠️ 需要调试 | 0/13 |

### 认证测试覆盖 ✅

- ✅ 登录页面正确显示
- ✅ 表单验证（HTML5 required 属性）
- ✅ Enter 键提交表单
- ✅ 有效凭证登录成功
- ✅ 无效凭证显示错误
- ✅ 会话持久化
- ✅ 受保护路由重定向
- ✅ 登录后导航到首页
- ✅ 响应式布局
- ✅ 加载状态显示

### 陪玩师浏览测试覆盖 ⚠️

以下测试需要进一步调试元素选择器：
- 首页显示和导航
- 游戏卡片展示
- 推荐陪玩师展示
- 陪玩师列表筛选
- 陪玩师详情页
- 收藏功能

## 目录结构

```
tests/e2e/
├── fixtures/           # 测试夹具
│   └── test-data.fixture.ts
├── pages/              # Page Object Models
│   ├── HomePage.ts
│   ├── LoginPage.ts
│   ├── PlayerListPage.ts
│   └── PlayerDetailPage.ts
├── helpers/            # 辅助函数
│   └── api-helpers.ts
├── auth.spec.ts        # 认证测试 (10 passing)
└── player-browsing.spec.ts  # 陪玩师浏览测试 (needs debugging)
```
