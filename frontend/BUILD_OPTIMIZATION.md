# 前端构建优化指南

## 已实施的优化

### 1. 代码分割策略

通过 Vite 的 `manualChunks` 配置，将代码分割为以下模块：

- **react-vendor**: React 核心库（react, react-dom, react-router-dom）
- **antd-vendor**: Ant Design UI 库及图标
- **chart-vendor**: 图表库（recharts）
- **utils-vendor**: 工具库（axios, dayjs, lodash-es）

### 2. 懒加载路由

所有页面组件使用 `React.lazy()` 动态导入：

```typescript
const Dashboard = lazy(() => import('@/pages/admin/Dashboard'));
```

配合 `<LazyLoad>` 组件提供加载状态。

### 3. 文件命名优化

- JS 文件: `assets/js/[name]-[hash].js`
- CSS 文件: `assets/css/[name]-[hash].css`
- 图片文件: `assets/images/[name]-[hash][ext]`
- 字体文件: `assets/fonts/[name]-[hash][ext]`

### 4. 压缩配置

使用 Terser 进行代码压缩：

- 移除 console.log/info/debug
- 移除 debugger 语句
- 移除注释
- 压缩变量名

### 5. 构建优化

- 启用 CSS 代码分割
- 禁用 sourcemap（生产环境）
- Chunk 大小警告限制: 1000KB

## 安装依赖

```bash
cd frontend
npm install --save-dev terser
```

## 构建命令

```bash
# 开发环境
npm run dev

# 生产构建
npm run build

# 预览生产构建
npm run preview
```

## 构建产物分析

构建后查看 `dist/` 目录结构：

```
dist/
├── assets/
│   ├── js/
│   │   ├── react-vendor-[hash].js      # React 核心 (~150KB gzip)
│   │   ├── antd-vendor-[hash].js       # Ant Design (~300KB gzip)
│   │   ├── chart-vendor-[hash].js      # 图表库 (~50KB gzip)
│   │   ├── utils-vendor-[hash].js      # 工具库 (~30KB gzip)
│   │   └── [page]-[hash].js            # 各页面模块
│   ├── css/
│   │   └── [name]-[hash].css
│   ├── images/
│   └── fonts/
└── index.html
```

## 性能指标

### 预期优化效果

- **首屏加载时间**: < 2s (3G 网络)
- **首次内容绘制 (FCP)**: < 1.5s
- **最大内容绘制 (LCP)**: < 2.5s
- **首次输入延迟 (FID)**: < 100ms
- **累积布局偏移 (CLS)**: < 0.1

### Bundle 大小目标

- 初始 JS bundle: < 200KB (gzip)
- 总 JS 大小: < 800KB (gzip)
- CSS 大小: < 100KB (gzip)

## 进一步优化建议

### 1. 图片优化

```typescript
// 使用 WebP 格式
import logo from '@/assets/logo.webp';

// 或使用 vite-plugin-imagemin
```

### 2. 预加载关键资源

```html
<link rel="preload" href="/assets/js/react-vendor.js" as="script">
```

### 3. CDN 加速

将静态资源部署到 CDN：

```typescript
// vite.config.ts
export default defineConfig({
  base: 'https://cdn.example.com/',
})
```

### 4. 启用 Brotli 压缩

```bash
npm install --save-dev vite-plugin-compression
```

```typescript
import viteCompression from 'vite-plugin-compression';

export default defineConfig({
  plugins: [
    viteCompression({
      algorithm: 'brotliCompress',
      ext: '.br',
    }),
  ],
})
```

### 5. 分析 Bundle 大小

```bash
npm install --save-dev rollup-plugin-visualizer
```

```typescript
import { visualizer } from 'rollup-plugin-visualizer';

export default defineConfig({
  plugins: [
    visualizer({
      open: true,
      gzipSize: true,
      brotliSize: true,
    }),
  ],
})
```

## 监控和测试

### Lighthouse 测试

```bash
# 安装 Lighthouse
npm install -g lighthouse

# 运行测试
lighthouse http://localhost:5173 --view
```

### Bundle 分析

```bash
npm run build
# 查看构建报告
```

## 注意事项

1. **懒加载边界**: 不要过度拆分，避免请求过多
2. **缓存策略**: 使用 hash 文件名实现长期缓存
3. **预加载**: 对关键路由使用 `<link rel="prefetch">`
4. **Tree Shaking**: 确保使用 ES6 模块导入
5. **动态导入**: 使用 `import()` 而非 `require()`

## 相关文档

- [Vite 构建优化](https://vitejs.dev/guide/build.html)
- [React 代码分割](https://react.dev/reference/react/lazy)
- [Web Vitals](https://web.dev/vitals/)
