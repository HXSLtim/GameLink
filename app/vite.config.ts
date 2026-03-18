import path from "node:path"
import react from "@vitejs/plugin-react"
import tailwindcss from "@tailwindcss/vite"
import { defineConfig } from "vite"

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    host: "0.0.0.0",
    port: 5175,
  },
  build: {
    // 启用 CSS 代码分割
    cssCodeSplit: true,
    // chunk 大小警告限制（KB）
    chunkSizeWarningLimit: 500,
    // Rollup 配置
    rollupOptions: {
      output: {
        // 自定义 chunk 文件名
        chunkFileNames: 'assets/js/[name]-[hash].js',
        // 自定义入口文件名
        entryFileNames: 'assets/js/[name]-[hash].js',
        // 自定义静态资源文件名
        assetFileNames: 'assets/[ext]/[name]-[hash][extname]',
        // 代码分割策略 - 将大型依赖分离到独立 chunk
        manualChunks: (id) => {
          // React 核心库
          if (id.includes('node_modules/react') || id.includes('node_modules/react-dom')) {
            return 'react-vendor';
          }
          // React Router
          if (id.includes('node_modules/react-router-dom')) {
            return 'router-vendor';
          }
          // Radix UI 组件
          if (id.includes('node_modules/@radix-ui')) {
            return 'radix-vendor';
          }
          // 图标库
          if (id.includes('node_modules/lucide-react')) {
            return 'icons-vendor';
          }
          // 状态管理和工具库
          if (id.includes('node_modules/zustand') ||
              id.includes('node_modules/axios') ||
              id.includes('node_modules/clsx') ||
              id.includes('node_modules/tailwind-merge')) {
            return 'utils-vendor';
          }
          // 其他 node_modules 依赖
          if (id.includes('node_modules')) {
            return 'vendor';
          }
        },
      },
    },
  },
})
