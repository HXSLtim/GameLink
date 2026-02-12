import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [uni()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 3000,
    strictPort: true,  // 端口被占用时报错而不是尝试下一个端口
  },
  css: {
    preprocessorOptions: {
      scss: {
        // 全局引入变量和 mixins 文件
        additionalData: `@use "@/styles/variables.scss" as *;
@use "@/styles/mixins.scss" as *;`,
        api: 'modern-compiler',
        silenceDeprecations: ['legacy-js-api', 'import'],
        quietDeps: true,
      },
    },
  },
  build: {
    // 使用 esbuild 压缩（Vite 内置，无需额外依赖）
    minify: 'esbuild',
    // 启用 CSS 代码分割
    cssCodeSplit: true,
    // chunk 大小警告阈值 (KB)
    chunkSizeWarningLimit: 500,
    // 资源内联阈值 - 小于 4KB 的资源内联为 base64
    assetsInlineLimit: 4096,
  },
  esbuild: {
    // 生产环境移除 console 和 debugger
    drop: process.env.NODE_ENV === 'production' ? ['console', 'debugger'] : [],
  },
})
