<script setup lang="ts">
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { useUserStore } from '@/store/user'
import { initTheme } from '@/composables/useTheme'

onLaunch(() => {
  console.log('App Launch')
  
  // 初始化用户状态
  const userStore = useUserStore()
  userStore.init()
  
  // 初始化主题
  initTheme()
})

onShow(() => {
  console.log('App Show')
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style lang="scss">
/* 
 * 全局样式已在 main.ts 中引入
 * @import './styles/index.scss';
 */
</style>
