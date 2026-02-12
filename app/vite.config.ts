import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
import { resolve } from 'path'
import viteCompression from 'vite-plugin-compression'
import { visualizer } from 'rollup-plugin-visualizer'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    uni(),
    // 开启 Gzip 压缩
    viteCompression({
      verbose: true,
      disable: false,
      threshold: 10240,
      algorithm: 'gzip',
      ext: '.gz',
    }),
    // 打包分析 (仅在分析时开启)
    process.env.ANALYZE === 'true' && visualizer({
      open: true,
      gzipSize: true,
      brotliSize: true,
      filename: 'stats.html'
    })
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5174,
    strictPort: false,
  },
  css: {
    preprocessorOptions: {
      scss: {
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
    rollupOptions: {
      output: {
        // 手动分包策略
        manualChunks(id) {
          if (id.includes('node_modules')) {
            if (id.includes('vue') || id.includes('pinia')) {
              return 'vendor-core'
            }
            if (id.includes('uv-ui')) {
              return 'vendor-ui'
            }
            return 'vendor-common'
          }
        }
      }
    }
  },
  esbuild: {
    // 生产环境移除 console 和 debugger
    drop: process.env.NODE_ENV === 'production' ? ['console', 'debugger'] : [],
  },
})
