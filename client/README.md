# GameLink Client

用户端 + 陪玩师端前端应用

## 技术栈

- React 18 + TypeScript
- Vite
- Ant Design 5
- React Router 7
- Axios

## 开发

```bash
# 安装依赖
npm install

# 启动开发服务器 (localhost:5174)
npm run dev

# 构建生产版本
npm run build

# 代码检查
npm run lint
```

## 目录结构

```
client/
├── src/
│   ├── api/           # API 客户端
│   ├── context/       # React Context
│   ├── layouts/       # 布局组件
│   ├── pages/         # 页面组件
│   │   ├── auth/      # 登录/注册
│   │   ├── Home/      # 首页
│   │   ├── player/    # 陪玩师相关
│   │   ├── order/     # 订单相关
│   │   └── user/      # 用户中心
│   ├── styles/        # 全局样式
│   ├── types/         # TypeScript 类型
│   ├── App.tsx        # 根组件
│   └── main.tsx       # 入口文件
├── public/            # 静态资源
└── package.json
```

## 功能模块

### 用户端
- 首页展示
- 陪玩师列表/详情
- 下单/支付
- 订单管理
- 个人中心

### 陪玩师端 (TODO)
- 接单管理
- 服务设置
- 收益统计
- 团队功能
