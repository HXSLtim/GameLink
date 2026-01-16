import path from "path"
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

import { VitePWA } from 'vite-plugin-pwa'

import { visualizer } from 'rollup-plugin-visualizer'

// https://vite.dev/config/
export default defineConfig(({ mode }) => ({
  plugins: [
    react(),
    mode === 'analyze' && visualizer({
      open: true,
      filename: 'dist/stats.html',
      gzipSize: true,
      brotliSize: true,
    }),
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['favicon.ico', 'apple-touch-icon.png', 'mask-icon.svg'],
      manifest: {
        name: 'GameLink - Find Your Pro Teammate',
        short_name: 'GameLink',
        description: 'Connect with pro gamers for the ultimate gaming experience.',
        theme_color: '#09090b',
        background_color: '#09090b',
        display: 'standalone',
        scope: '/',
        start_url: '/',
        orientation: 'portrait',
        icons: [
          {
            src: 'pwa-192x192.png',
            sizes: '192x192',
            type: 'image/png'
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png'
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any maskable'
          }
        ]
      }
    })
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  // Build optimization: Code splitting for better caching and faster loads
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          // Core vendor libraries (rarely change)
          if (id.includes('node_modules/react/') || id.includes('node_modules/react-dom/') || id.includes('node_modules/react-router-dom/')) {
            return 'vendor-react';
          }
          // State management
          if (id.includes('node_modules/zustand/')) {
            return 'vendor-state';
          }
          // UI components (Radix UI)
          if (id.includes('node_modules/@radix-ui/')) {
            return 'vendor-ui';
          }
          // Form handling
          if (id.includes('node_modules/react-hook-form/') || id.includes('node_modules/@hookform/') || id.includes('node_modules/zod/')) {
            return 'vendor-form';
          }
          // i18n
          if (id.includes('node_modules/i18next') || id.includes('node_modules/react-i18next/')) {
            return 'vendor-i18n';
          }
          // HTTP and utilities
          if (id.includes('node_modules/axios/') || id.includes('node_modules/date-fns/') || id.includes('node_modules/clsx/') || id.includes('node_modules/tailwind-merge/')) {
            return 'vendor-utils';
          }
          // Voice/Video (large, load on demand)
          if (id.includes('node_modules/trtc-js-sdk/')) {
            return 'vendor-trtc';
          }
        },
      },
    },
    // Warn if chunks exceed 500KB
    chunkSizeWarningLimit: 500,
    // Enable source maps for production debugging (optional)
    sourcemap: mode === 'development',
  },
  server: {
    port: 5000,
    host: true, // Listen on all addresses
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
    },
  }
}))
