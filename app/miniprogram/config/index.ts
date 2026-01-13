/**
 * 应用配置
 */

// 环境类型
type EnvType = 'development' | 'production'

// 当前环境（可通过微信开发者工具的编译模式切换）
const ENV: EnvType = 'development'

// API 配置
const API_CONFIG = {
  development: {
    baseUrl: 'http://localhost:8080/api/v1',
    wsUrl: 'ws://localhost:8080/ws',
  },
  production: {
    baseUrl: 'https://api.gamelink.com/api/v1',
    wsUrl: 'wss://api.gamelink.com/ws',
  },
}

// 导出配置
export const config = {
  env: ENV,
  api: API_CONFIG[ENV],
  
  // 应用信息
  appName: 'GameLink',
  version: '1.0.0',
  
  // 分页默认值
  pageSize: 10,
  
  // 请求超时（毫秒）
  timeout: 10000,
  
  // 是否开启调试日志
  debug: ENV === 'development',
}

export default config
