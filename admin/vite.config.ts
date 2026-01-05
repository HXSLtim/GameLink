/// <reference types="vitest" />
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import viteCompression from 'vite-plugin-compression'
import { visualizer } from 'rollup-plugin-visualizer'
import os from 'os'
import pc from 'picocolors'
import { VitePWA } from 'vite-plugin-pwa'

// 获取所有可用的网络接口IP地址
function getNetworkIPs() {
  const interfaces = os.networkInterfaces()
  const ips: string[] = []

  for (const name of Object.keys(interfaces)) {
    for (const iface of interfaces[name] || []) {
      // 跳过内部IP和IPv6
      if (!iface.internal && iface.family === 'IPv4') {
        ips.push(iface.address)
      }
    }
  }

  return ips
}

// 自定义插件：显示所有可访问的IP地址
function showNetworkIPs() {
  return {
    name: 'show-network-ips',
    configureServer(server: any) {
      const { port = 5173, https } = server.config.server
      const protocol = https ? 'https' : 'http'
      const localhost = `${protocol}://localhost:${port}`
      const networkIPs = getNetworkIPs()

      // 服务器启动后显示
      server.httpServer?.once('listening', () => {
        console.log('\n' + pc.bold('🚀 GameLink Admin Server is running!'))
        console.log('\n' + pc.bold('  Local:   ') + pc.cyan(localhost))

        if (networkIPs.length > 0) {
          console.log(pc.bold('  Network: ') + pc.cyan(
            networkIPs.map(ip => `${protocol}://${ip}:${port}`).join('\n          ')
          ))
        }

        console.log(pc.bold('\n  Press ') + pc.green('Enter') + pc.bold(' to open in browser'))
        console.log('  ' + pc.cyan('-'.repeat(50)) + '\n')
      })
    }
  }
}

// https://vite.dev/config/
export default defineConfig(({ mode }) => ({
  // Vitest 测试配置
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    include: ['src/**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts,jsx,tsx}'],
    testTimeout: 10000,
    hookTimeout: 10000,
    // 缓存优化
    cache: {
      dir: 'node_modules/.vitest',
    },
    // 覆盖率配置
    coverage: {
      reporter: ['text', 'json', 'html'],
      exclude: ['node_modules/', 'src/test/', '**/*.d.ts', '**/*.test.{ts,tsx}'],
    },
  },
  plugins: [
    react(),
    showNetworkIPs(), // 显示所有网络IP地址
    // PWA 支持
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['icon.svg'],
      manifest: {
        name: 'GameLink Admin Panel',
        short_name: 'GameLink',
        description: 'GameLink 陪玩平台管理后台',
        theme_color: '#1890ff',
        background_color: '#ffffff',
        display: 'standalone',
        orientation: 'portrait-primary',
        icons: [
          {
            src: 'icon.svg',
            sizes: '192x192 512x512',
            type: 'image/svg+xml',
            purpose: 'any maskable'
          }
        ]
      },
      workbox: {
        // 缓存策略
        runtimeCaching: [
          {
            // API 请求 - NetworkFirst 策略
            urlPattern: /^https:\/\/.*\/api\/.*/i,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'api-cache',
              expiration: {
                maxEntries: 100,
                maxAgeSeconds: 60 * 60 * 24 // 24 小时
              },
              cacheableResponse: {
                statuses: [0, 200]
              }
            }
          },
          {
            // 静态资源 - CacheFirst 策略
            urlPattern: /\.(?:png|jpg|jpeg|svg|gif|webp|ico)$/i,
            handler: 'CacheFirst',
            options: {
              cacheName: 'image-cache',
              expiration: {
                maxEntries: 200,
                maxAgeSeconds: 60 * 60 * 24 * 30 // 30 天
              }
            }
          },
          {
            // CSS 和 JS - StaleWhileRevalidate 策略
            urlPattern: /\.(?:css|js)$/i,
            handler: 'StaleWhileRevalidate',
            options: {
              cacheName: 'static-resources-cache',
              expiration: {
                maxEntries: 100,
                maxAgeSeconds: 60 * 60 * 24 * 7 // 7 天
              }
            }
          }
        ],
        // 清理过时的缓存
        cleanupOutdatedCaches: true,
        // 导航预加载
        navigateFallback: null,
        // 开发环境禁用范围
        navigateFallbackDenylist: [/^\/api/]
      },
      devOptions: {
        enabled: false, // 开发环境禁用 PWA，避免重复注册
        type: 'module'
      }
    }),
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
    // 监听所有网络接口，可以通过局域网IP访问
    host: '0.0.0.0',
    // 严格端口，如果被占用则失败而不是尝试下一个端口
    strictPort: false,
    proxy: {
      '/api/v1': {
        target: 'http://localhost:8080',  // Docker后端端口
        changeOrigin: true,
        // 支持 WebSocket
        ws: true,
      },
    },
  },
  preview: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api/v1': {
        target: 'http://localhost:8081',  // 生产环境后端端口
        changeOrigin: true,
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
