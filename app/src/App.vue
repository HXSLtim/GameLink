<script setup lang="ts">
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { useUserStore, useAppStore } from '@/store'
import { initTheme } from '@/composables/useTheme'
import { getUnreadCount } from '@/api/notification'
import { setupRouteGuard } from '@/utils/routeGuard'

onLaunch(() => {
  // 初始化用户状态
  const userStore = useUserStore()
  userStore.init()
  
  // 初始化应用状态
  const appStore = useAppStore()
  appStore.init()
  
  // 初始化主题
  initTheme()
  appStore.applyTheme()

  // 路由守卫
  setupRouteGuard()
})

onShow(() => {
  // 刷新未读消息数
  syncUnreadCount()
})

onHide(() => {})

/**
 * 同步未读消息数
 */
async function syncUnreadCount() {
  const userStore = useUserStore()
  const appStore = useAppStore()
  
  if (!userStore.isLoggedIn) return
  
  try {
    const res = await getUnreadCount({ showError: false })
    if (res.data) {
      appStore.setUnreadCount(res.data.total)
    }
  } catch (error) {
    // 静默失败
    console.log('Sync unread count failed:', error)
  }
}
</script>

<style lang="scss">
/* 
 * 全局样式已在 main.ts 中引入
 * @import './styles/index.scss';
 */
</style>
