/**
 * 登录专用 Hook
 */
import { ref, reactive, computed } from 'vue'
import { useUserStore, normalizeUserInfo } from '@/store/user'
import { login, doWeChatLogin } from '@/api/auth'
import { consumeRedirectPath, redirectToUrl } from '@/utils/routeGuard'

export function useLogin() {
  const userStore = useUserStore()
  
  // 状态
  const loginLoading = ref(false)
  const wechatLoading = ref(false)
  const showAccountLogin = ref(false)
  
  // 表单
  const form = reactive({
    username: '',
    password: '',
  })
  
  // 是否可登录
  const canLogin = computed(() => form.username && form.password)
  
  // 账号密码登录
  const handleLogin = async () => {
    if (!canLogin.value || loginLoading.value) return
    
    loginLoading.value = true
    
    try {
      const res = await login({
        identifier: form.username,
        password: form.password,
      })
      
      const data = res.data as any
      
      // 保存用户信息和 token
      userStore.setToken(data.accessToken)
      userStore.setUserInfo(normalizeUserInfo(data.user))
      
      uni.showToast({ title: '登录成功', icon: 'success' })
      
      // 重定向
      setTimeout(() => {
        handleRedirect()
      }, 500)
    } catch (error: any) {
      uni.showToast({ title: error?.message || '登录失败', icon: 'none' })
    } finally {
      loginLoading.value = false
    }
  }
  
  // 微信登录
  const handleWechatLogin = async () => {
    if (wechatLoading.value) return
    
    wechatLoading.value = true
    
    try {
      // #ifdef MP-WEIXIN
      const loginRes = await new Promise<UniApp.LoginRes>((resolve, reject) => {
        uni.login({
          provider: 'weixin',
          success: resolve,
          fail: reject,
        })
      })
      
      const res = await doWeChatLogin({ code: loginRes.code })
      const data = res.data as any
      
      userStore.setToken(data.accessToken)
      userStore.setUserInfo(normalizeUserInfo(data.user))
      
      uni.showToast({ title: '登录成功', icon: 'success' })
      
      setTimeout(() => {
        handleRedirect()
      }, 500)
      // #endif
    } catch (error: any) {
      uni.showToast({ title: error?.message || '微信登录失败', icon: 'none' })
    } finally {
      wechatLoading.value = false
    }
  }
  
  // 处理重定向
  const handleRedirect = () => {
    const redirectPath = consumeRedirectPath()
    if (redirectPath) {
      redirectToUrl(redirectPath)
    } else {
      uni.switchTab({ url: '/pages/index/index' })
    }
  }
  
  // 导航
  const goToRegister = () => uni.navigateTo({ url: '/pages/auth/register/index' })
  const goToAgreement = (type: string) => uni.navigateTo({ url: `/pages/agreement/index?type=${type}` })
  
  return {
    // 状态
    loginLoading,
    wechatLoading,
    showAccountLogin,
    
    // 数据
    form,
    canLogin,
    
    // 方法
    handleLogin,
    handleWechatLogin,
    goToRegister,
    goToAgreement,
  }
}
