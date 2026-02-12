/**
 * 登录专用 Hook
 */
import { ref, reactive, computed } from 'vue'
import { useUserStore, normalizeUserInfo } from '@/store/user'
import { login, doWeChatLogin } from '@/api/auth'
import { ApiError } from '@/api/request'
import { consumeRedirectPath, redirectToUrl } from '@/utils/routeGuard'
import type { AgreementType } from '@/types/agreement'

export function useLogin() {
  const userStore = useUserStore()

  // 状态
  const loginLoading = ref(false)
  const wechatLoading = ref(false)
  const showAccountLogin = ref(false)

  // 表单
  const form = reactive({
    account: '',  // 手机号或邮箱
    password: '',
  })

  // 是否可登录
  const canLogin = computed(() => form.account.trim() && form.password)

  const resolveLoginErrorMessage = (error: unknown) => {
    if (!form.account.trim() || !form.password) {
      return '请输入账号和密码'
    }
    if (error instanceof ApiError) {
      if (error.code === 400 || error.code === 401) {
        return '账号或密码错误'
      }
      if (error.code === 429) {
        return '请求过于频繁，请稍后再试'
      }
    }
    return '登录失败，请稍后重试'
  }

  // 账号密码登录
  const handleLogin = async () => {
    if (loginLoading.value) return
    if (!canLogin.value) {
      uni.showToast({ title: '请输入账号和密码', icon: 'none' })
      return
    }

    loginLoading.value = true

    try {
      // 判断是手机号还是邮箱
      const account = form.account.trim()
      const isPhone = /^1[3-9]\d{9}$/.test(account)
      const isEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(account)

      if (!isPhone && !isEmail) {
        uni.showToast({ title: '请输入有效的手机号或邮箱', icon: 'none' })
        loginLoading.value = false
        return
      }

      const res = await login({
        username: account,
        password: form.password,
      }, { showError: false })

      const data = res.data as any

      const accessToken = data.accessToken || data.token || ''
      const refreshToken = data.refreshToken || ''
      if (!accessToken) {
        throw new ApiError('登录凭证异常', -1)
      }

      // 保存用户信息和 token
      userStore.setToken(accessToken, refreshToken)
      userStore.setUserInfo(normalizeUserInfo(data.user))

      uni.showToast({ title: '登录成功', icon: 'success' })

      // 重定向
      setTimeout(() => {
        handleRedirect()
      }, 500)
    } catch (error: any) {
      uni.showToast({ title: resolveLoginErrorMessage(error), icon: 'none' })
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
      const data = await doWeChatLogin()

      const accessToken = data.accessToken || data.token || ''
      const refreshToken = data.refreshToken || ''
      if (!accessToken) {
        throw new ApiError('登录凭证异常', -1)
      }

      userStore.setToken(accessToken, refreshToken)
      userStore.setUserInfo(normalizeUserInfo(data.user))

      uni.showToast({ title: '登录成功', icon: 'success' })

      setTimeout(() => {
        handleRedirect()
      }, 500)
      // #endif
    } catch (error: any) {
      uni.showToast({ title: '微信登录失败，请稍后重试', icon: 'none' })
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
  const goToAgreement = (type: AgreementType) =>
    uni.navigateTo({ url: `/pages/agreement/index?type=${type}` })

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
