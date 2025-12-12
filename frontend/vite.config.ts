import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import viteCompression from 'vite-plugin-compression'
import { visualizer } from 'rollup-plugin-visualizer'

// https://vite.dev/config/
export default defineConfig(({ mode }) => ({
  plugins: [
    react(),
    // Gzip 压缩
    viteCompression({
      verbose: true,
      disable: false,
      threshold: 10240, // 10KB 以上才压缩
      algorithm: 'gzip',
      ext: '.gz',
    }),
    // Brotli 压缩（更高压缩率）
    viteCompression({
      verbose: true,
      disable: false,
      threshold: 10240,
      algorithm: 'brotliCompress',
      ext: '.br',
    }),
    // Bundle 分析（仅在 analyze 模式下）
    mode === 'analyze' && visualizer({
      open: true,
      gzipSize: true,
      brotliSize: true,
      filename: 'dist/stats.html',
    }),
  ].filter(Boolean),
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      '/api/v1': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        // 支持 WebSocket
        ws: true,
      },
    },
  },
  optimizeDeps: {
    include: ['recharts'],
  },
  build: {
    // 启用 CSS 代码分割
    cssCodeSplit: true,
    // 构建后的文件目录
    outDir: 'dist',
    // 静态资源目录
    assetsDir: 'assets',
    // chunk 大小警告限制（KB）
    chunkSizeWarningLimit: 1000,
    // Rollup 配置
    rollupOptions: {
      output: {
        // 手动代码分割策略
        manualChunks: (id) => {
          // React 核心库 + React Router（合并到一起，避免加载顺序问题）
          if (id.includes('node_modules/react/') || 
              id.includes('node_modules/react-dom/') || 
              id.includes('node_modules/scheduler/') ||
              id.includes('node_modules/react-router') || 
              id.includes('node_modules/@remix-run')) {
            return 'react-vendor';
          }
          // 工具库（独立，不依赖其他库）
          if (id.includes('node_modules/axios') || 
              id.includes('node_modules/dayjs') || 
              id.includes('node_modules/lodash-es')) {
            return 'utils-vendor';
          }
          // 图表库、Ant Design 等随页面按需加载
        },
        // 自定义 chunk 文件名
        chunkFileNames: 'assets/js/[name]-[hash].js',
        // 自定义入口文件名
        entryFileNames: 'assets/js/[name]-[hash].js',
        // 自定义静态资源文件名
        assetFileNames: (assetInfo) => {
          const info = assetInfo.name?.split('.');
          let extType = info?.[info.length - 1];
          if (/\.(png|jpe?g|gif|svg|webp|ico)$/i.test(assetInfo.name || '')) {
            extType = 'images';
          } else if (/\.(woff2?|eot|ttf|otf)$/i.test(assetInfo.name || '')) {
            extType = 'fonts';
          } else if (/\.css$/i.test(assetInfo.name || '')) {
            extType = 'css';
          }
          return `assets/${extType}/[name]-[hash][extname]`;
        },
      },
    },
    // 压缩配置
    minify: 'terser',
    terserOptions: {
      compress: {
        // 生产环境移除 console
        drop_console: true,
        drop_debugger: true,
        // 移除无用代码
        pure_funcs: ['console.log', 'console.info', 'console.debug'],
      },
      format: {
        // 移除注释
        comments: false,
      },
    },
    // 启用源码映射（生产环境可选）
    sourcemap: false,
    // 报告压缩后的文件大小
    reportCompressedSize: true,
  },
}))
