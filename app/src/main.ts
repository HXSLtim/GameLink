import { createSSRApp } from 'vue'
import { createPinia } from 'pinia'
import uvUI from '@climblee/uv-ui'
import App from './App.vue'

// 引入全局样式
import './styles/index.scss'

export function createApp() {
  const app = createSSRApp(App)
  
  // 使用 Pinia 状态管理
  const pinia = createPinia()
  app.use(pinia)
  
  // 使用 uv-ui 组件库
  app.use(uvUI)
  
  return {
    app,
  }
}
