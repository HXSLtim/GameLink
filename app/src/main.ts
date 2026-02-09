import { createSSRApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'

// 引入全局样式
import './styles/index.scss'

// uv-ui 通过 easycom 按需自动导入（见 pages.json），无需全量注册

export function createApp() {
  const app = createSSRApp(App)
  
  // 使用 Pinia 状态管理
  const pinia = createPinia()
  app.use(pinia)
  
  return {
    app,
  }
}
