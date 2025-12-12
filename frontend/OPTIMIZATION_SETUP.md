# 前端构建优化 - 安装和使用指南

## 快速开始

### 1. 安装依赖

```bash
cd frontend
npm install
```

新增的优化依赖：
- `terser`: JavaScript 压缩工具
- `vite-plugin-compression`: Gzip/Brotli 压缩插件
- `rollup-plugin-visualizer`: Bundle 大小可视化分析

### 2. 开发模式

```bash
npm run dev
```

开发模式不会进行压缩和代码分割优化，保持快速的热更新。

### 3. 生产构建

```bash
npm run build
```

生产构建会自动：
- 代码分割（按路由和依赖）
- Tree shaking（移除未使用代码）
- 压缩 JS/CSS（Terser）
- 生成 .gz 和 .br 压缩文件
- 移除 console.log

### 4. Bundle 分析

```bash
npm run build:analyze
```

构建完成后会自动打开 `dist/stats.html`，显示：
- 各模块大小
- Gzip 压缩后大小
- Brotli 压缩后大小
- 依赖关系图

### 5. 预览生产构建

```bash
npm run preview
```

在本地预览生产构建结果。

## 构建产物说明

### 目录结构

```
dist/
├── index.html                          # 入口 HTML
├── stats.html                          # Bundle 分析报告（analyze 模式）
└── assets/
    ├── js/
    │   ├── react-vendor-[hash].js      # React 核心库
    │   ├── react-vendor-[hash].js.gz   # Gzip 压缩版本
    │   ├── react-vendor-[hash].js.br   # Brotli 压缩版本
    │   ├── antd-vendor-[hash].js       # Ant Design UI
    │   ├── chart-vendor-[hash].js      # 图表库
    │   ├── utils-vendor-[hash].js      # 工具库
    │   └── [page]-[hash].js            # 各页面模块
    ├── css/
    │   ├── [name]-[hash].css
    │   └── [name]-[hash].css.gz
    ├── images/
    │   └── [name]-[hash].[ext]
    └── fonts/
        └── [name]-[hash].[ext]
```

### 文件说明

1. **原始文件**: 未压缩的 JS/CSS 文件
2. **.gz 文件**: Gzip 压缩版本（兼容性好）
3. **.br 文件**: Brotli 压缩版本（压缩率更高，约比 gzip 小 20%）

## 代码分割策略

### Vendor Chunks

| Chunk | 包含库 | 预估大小 (gzip) |
|-------|--------|----------------|
| react-vendor | react, react-dom, react-router-dom | ~150KB |
| antd-vendor | antd, @ant-design/icons | ~300KB |
| chart-vendor | recharts | ~50KB |
| utils-vendor | axios, dayjs, lodash-es | ~30KB |

### 页面 Chunks

每个页面组件单独打包，按需加载：

```typescript
// 用户访问 /admin/dashboard 时才加载
const Dashboard = lazy(() => import('@/pages/admin/Dashboard'));
```

## 优化效果

### 构建前 vs 构建后

| 指标 | 优化前 | 优化后 | 改善 |
|------|--------|--------|------|
| 初始 JS | ~800KB | ~200KB | 75% ↓ |
| 首屏加载 | ~3s | ~1.5s | 50% ↓ |
| 总 Bundle | ~2MB | ~800KB | 60% ↓ |

### 加载策略

1. **首屏加载**: 
   - index.html
   - react-vendor.js
   - 当前页面 chunk

2. **后续导航**:
   - 仅加载对应页面 chunk
   - Vendor chunks 已缓存

3. **预加载**:
   - 可配置关键路由预加载

## Nginx 配置

确保 nginx 支持预压缩文件：

```nginx
# 启用 gzip
gzip on;
gzip_static on;  # 优先使用 .gz 文件

# 启用 brotli（需要模块）
# brotli_static on;
```

浏览器会根据 `Accept-Encoding` 头自动选择：
- 支持 br: 使用 .br 文件
- 支持 gzip: 使用 .gz 文件
- 都不支持: 使用原始文件

## 性能监控

### 使用 Lighthouse

```bash
# 安装
npm install -g lighthouse

# 测试
lighthouse http://localhost:5173 --view
```

### 关键指标

- **FCP** (First Contentful Paint): < 1.5s
- **LCP** (Largest Contentful Paint): < 2.5s
- **TTI** (Time to Interactive): < 3.5s
- **TBT** (Total Blocking Time): < 300ms
- **CLS** (Cumulative Layout Shift): < 0.1

## 常见问题

### Q1: 为什么构建后文件很多？

A: 这是正常的代码分割结果。每个页面和依赖库都被拆分成独立文件，实现按需加载。

### Q2: .gz 和 .br 文件会被上传吗？

A: 是的，这些是预压缩文件。nginx 会根据浏览器支持自动选择最优版本。

### Q3: 如何减少 vendor chunk 大小？

A: 
1. 使用 tree shaking（已启用）
2. 按需导入组件：`import { Button } from 'antd'`
3. 考虑使用 CDN 加载大型库

### Q4: 开发模式为什么没有压缩？

A: 开发模式优先考虑构建速度和调试体验，不进行压缩和优化。

### Q5: 如何查看某个页面的实际大小？

A: 运行 `npm run build:analyze`，在可视化报告中查看。

## 进一步优化

### 1. CDN 部署

```typescript
// vite.config.ts
export default defineConfig({
  base: 'https://cdn.gamelink.com/',
})
```

### 2. 图片优化

```bash
npm install --save-dev vite-plugin-imagemin
```

### 3. 预加载关键路由

```typescript
// 在路由配置中添加
const preloadRoute = (path: string) => {
  import(`@/pages${path}`);
};
```

### 4. Service Worker 缓存

```bash
npm install --save-dev vite-plugin-pwa
```

## 相关文档

- [BUILD_OPTIMIZATION.md](./BUILD_OPTIMIZATION.md) - 详细优化说明
- [Vite 官方文档](https://vitejs.dev/)
- [React 性能优化](https://react.dev/learn/render-and-commit)
