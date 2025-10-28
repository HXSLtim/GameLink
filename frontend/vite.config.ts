import { defineConfig, Plugin } from 'vite';
import react from '@vitejs/plugin-react';
import FullReload from 'vite-plugin-full-reload';
import path from 'path';

/**
 * 开发环境认证 Mock 插件
 * 从环境变量读取凭据，避免硬编码
 */
function devAuthMock(): Plugin {
  // 从环境变量读取 Mock 凭据
  const MOCK_USERNAME = process.env.VITE_DEV_MOCK_USERNAME || 'admin';
  const MOCK_PASSWORD = process.env.VITE_DEV_MOCK_PASSWORD || 'admin123';
  const MOCK_TOKEN = process.env.VITE_DEV_MOCK_TOKEN || 'dev-token';

  return {
    name: 'dev-auth-mock',
    apply: 'serve',
    configureServer(server) {
      server.middlewares.use(async (req, res, next) => {
        const url = req.url || '';
        const method = (req.method || 'GET').toUpperCase();

        // Mock 登录接口
        if (url.startsWith('/api/v1/auth/login') && method === 'POST') {
          let body = '';
          req.on('data', (c) => (body += c));
          req.on('end', () => {
            try {
              const { username, password } = JSON.parse(body || '{}');
              if (username === MOCK_USERNAME && password === MOCK_PASSWORD) {
                const payload = {
                  success: true,
                  code: 0,
                  message: 'ok',
                  data: { token: MOCK_TOKEN, expires_in: 3600 },
                };
                res.setHeader('Content-Type', 'application/json');
                res.end(JSON.stringify(payload));
              } else {
                const payload = {
                  success: false,
                  code: 401,
                  message: 'Invalid credentials',
                  data: null,
                };
                res.statusCode = 401;
                res.setHeader('Content-Type', 'application/json');
                res.end(JSON.stringify(payload));
              }
            } catch (e) {
              res.statusCode = 400;
              res.end('bad request');
            }
          });
          return;
        }

        // Mock 获取当前用户接口
        if (url.startsWith('/api/v1/auth/me') && method === 'GET') {
          const auth = req.headers['authorization'] || '';
          if (typeof auth === 'string' && auth === `Bearer ${MOCK_TOKEN}`) {
            const payload = {
              success: true,
              code: 0,
              message: 'ok',
              data: { id: 1, username: MOCK_USERNAME, role: 'admin' },
            };
            res.setHeader('Content-Type', 'application/json');
            res.end(JSON.stringify(payload));
          } else {
            const payload = { success: false, code: 401, message: 'unauthorized', data: null };
            res.statusCode = 401;
            res.setHeader('Content-Type', 'application/json');
            res.end(JSON.stringify(payload));
          }
          return;
        }
        next();
      });
    },
  };
}

export default defineConfig({
  plugins: [
    react(),
    devAuthMock(),
    // 当后端或配置文件变更时，触发浏览器整页刷新
    FullReload(['backend/**', 'configs/**', 'frontend/index.html']),
  ],
  resolve: {
    alias: {
      components: path.resolve(__dirname, './src/components'),
      pages: path.resolve(__dirname, './src/pages'),
      utils: path.resolve(__dirname, './src/utils'),
      hooks: path.resolve(__dirname, './src/hooks'),
      services: path.resolve(__dirname, './src/services'),
      types: path.resolve(__dirname, './src/types'),
      contexts: path.resolve(__dirname, './src/contexts'),
      styles: path.resolve(__dirname, './src/styles'),
    },
  },
  css: {
    preprocessorOptions: {
      less: {
        javascriptEnabled: true,
        modifyVars: {
          // 品牌主色（示例，可按需修改）
          'arcoblue-6': '#1772f6',
        },
      },
    },
  },
  server: {
    host: true,
    port: 5173,
    open: false,
    hmr: {
      protocol: 'ws',
      host: 'localhost',
      overlay: true,
    },
    watch: {
      // 在 WSL 下启用 polling 以确保文件变化被检测到
      usePolling: true,
      interval: 100,
    },
    proxy: {
      // 代理后端 API，方便本地联调
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path,
      },
    },
  },
  build: {
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          // 将React相关库打包
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
        },
        chunkFileNames: 'assets/js/[name]-[hash].js',
        entryFileNames: 'assets/js/[name]-[hash].js',
        assetFileNames: 'assets/[ext]/[name]-[hash].[ext]',
      },
    },
    chunkSizeWarningLimit: 1000,
  },
});
