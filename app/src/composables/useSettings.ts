/**
 * 设置页专用 Hook
 */
import { ref, reactive, computed } from 'vue'
import { useTheme } from '@/composables/useTheme'
import { useUserStore } from '@/store/user'
import { confirmDialog } from '@/composables/useConfirmDialog'
import type { SettingsItem, SettingsState } from '@/types/ui'
import { maskPhone } from '@/utils/format'

export function useSettings() {
  const { isDark, toggleTheme } = useTheme()
  const userStore = useUserStore()
  
  // 设置状态
  const settings = reactive<SettingsState>({
    pushEnabled: true,
    messageEnabled: true,
    orderEnabled: true,
    promotionEnabled: false,
    showOnlineStatus: true,
    allowStrangerMessage: true,
  })
  
  const userPhone = ref('138****8888')
  const blacklistCount = ref(0)
  const cacheSize = ref('12.5MB')
  
  // 账号安全菜单
  const securityItems = computed((): SettingsItem[] => [
    { key: 'password', label: '修改密码', icon: 'lock' },
    { key: 'phone', label: '更换手机号', icon: 'phone', value: maskPhone(userPhone.value) },
    { key: 'bind', label: '账号绑定', icon: 'link' },
  ])
  
  // 通知设置菜单
  const notificationItems = computed((): SettingsItem[] => [
    { key: 'push', label: '推送通知', icon: 'bell', type: 'switch', checked: settings.pushEnabled },
    { key: 'message', label: '消息提醒', icon: 'chat', type: 'switch', checked: settings.messageEnabled },
    { key: 'order', label: '订单通知', icon: 'file-text', type: 'switch', checked: settings.orderEnabled },
    { key: 'promotion', label: '活动推广', icon: 'gift', type: 'switch', checked: settings.promotionEnabled },
  ])
  
  // 隐私设置菜单
  const privacyItems = computed((): SettingsItem[] => [
    { key: 'blacklist', label: '黑名单', icon: 'minus-circle', value: blacklistCount.value > 0 ? `${blacklistCount.value}人` : undefined },
    { key: 'showOnline', label: '显示在线状态', icon: 'eye', type: 'switch', checked: settings.showOnlineStatus },
    { key: 'allowStranger', label: '允许陌生人私信', icon: 'account', type: 'switch', checked: settings.allowStrangerMessage },
  ])
  
  // 通用设置菜单
  const generalItems = computed((): SettingsItem[] => [
    { key: 'theme', label: '深色模式', icon: isDark.value ? 'eye-off' : 'eye', type: 'switch', checked: isDark.value },
    { key: 'language', label: '语言', icon: 'earth', value: '简体中文' },
    { key: 'cache', label: '清除缓存', icon: 'trash', value: cacheSize.value },
  ])
  
  // 关于菜单
  const aboutItems = computed((): SettingsItem[] => [
    { key: 'version', label: '版本', icon: 'info-circle', value: 'v1.0.0' },
    { key: 'agreement', label: '用户协议', icon: 'file-text' },
    { key: 'privacy', label: '隐私政策', icon: 'lock' },
    { key: 'about', label: '关于我们', icon: 'info-circle' },
  ])
  
  // 处理开关
  const handleSwitch = (key: string, value: boolean) => {
    switch (key) {
      case 'push':
        settings.pushEnabled = value
        break
      case 'message':
        settings.messageEnabled = value
        break
      case 'order':
        settings.orderEnabled = value
        break
      case 'promotion':
        settings.promotionEnabled = value
        break
      case 'showOnline':
        settings.showOnlineStatus = value
        break
      case 'allowStranger':
        settings.allowStrangerMessage = value
        break
      case 'theme':
        toggleTheme()
        break
    }
    uni.showToast({ title: '设置已保存', icon: 'none' })
  }
  
  // 处理点击
  const handleClick = (key: string) => {
    const routes: Record<string, string> = {
      password: '/pages/settings/password/index',
      phone: '/pages/settings/phone/index',
      bind: '/pages/settings/bind/index',
      blacklist: '/pages/settings/blacklist/index',
      language: '/pages/settings/language/index',
      agreement: '/pages/agreement/index?type=user',
      privacy: '/pages/agreement/index?type=privacy',
      about: '/pages/about/index',
    }
    
    const route = routes[key]
    if (route) {
      uni.navigateTo({ url: route })
    } else if (key === 'cache') {
      clearCache()
    } else if (key === 'version') {
      checkUpdate()
    } else {
      uni.showToast({ title: '功能开发中', icon: 'none' })
    }
  }
  
  // 清除缓存
  const clearCache = async () => {
    const confirmed = await confirmDialog({
      title: '清除缓存',
      content: '确定要清除所有缓存吗？',
    })
    if (!confirmed) return
    uni.showLoading({ title: '清除中...' })
    setTimeout(() => {
      cacheSize.value = '0KB'
      uni.hideLoading()
      uni.showToast({ title: '清除成功', icon: 'success' })
    }, 1000)
  }
  
  // 检查更新
  const checkUpdate = () => {
    uni.showToast({ title: '已是最新版本', icon: 'none' })
  }
  
  // 退出登录
  const handleLogout = async () => {
    const confirmed = await confirmDialog({
      title: '确认退出',
      content: '确定要退出登录吗？',
    })
    if (!confirmed) return
    userStore.logout()
    uni.reLaunch({ url: '/pages/index/index' })
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  
  return {
    // 菜单
    securityItems,
    notificationItems,
    privacyItems,
    generalItems,
    aboutItems,
    
    // 方法
    handleSwitch,
    handleClick,
    handleLogout,
    goBack,
  }
}
