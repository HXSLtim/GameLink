/**
 * 登录专用 Hook
 */
import { ref, reactive, computed, onUnmounted } from 'vue'
import { useUserStore, normalizeUserInfo } from '@/store/user'
import { login, doWeChatLogin, sendSmsCode, loginWithPhone } from '@/api/auth'
import { ApiError } from '@/api/request'
import { consumeRedirectPath, redirectToUrl } from '@/utils/routeGuard'
import type { AgreementType } from '@/types/agreement'

export function useLogin() {
  const userStore = useUserStore()

  // 状态
  const loginLoading = ref(false)
  const wechatLoading = ref(false)
  const smsSending = ref(false)
  const showAccountLogin = ref(false)
  const loginMode = ref<'password' | 'sms'>('password')
  const smsCooldown = ref(0)
  let smsCooldownTimer: ReturnType<typeof setInterval> | null = null

  // 表单
  const form = reactive({
    account: '',  // 手机号或邮箱
    password: '',
  })
  const phoneForm = reactive({
    phone: '',
    code: '',
  })

  // 是否可登录
  const canLogin = computed(() => form.account.trim() && form.password)
  const canPhoneLogin = computed(() => /^1[3-9]\d{9}$/.test(phoneForm.phone) && /^\d{6}$/.test(phoneForm.code))
  const canSendSmsCode = computed(() => /^1[3-9]\d{9}$/.test(phoneForm.phone) && smsCooldown.value === 0 && !smsSending.value)

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

  const startSmsCooldown = (seconds: number) => {
    if (smsCooldownTimer) {
      clearInterval(smsCooldownTimer)
      smsCooldownTimer = null
    }
    smsCooldown.value = seconds
    smsCooldownTimer = setInterval(() => {
      smsCooldown.value -= 1
      if (smsCooldown.value <= 0) {
        smsCooldown.value = 0
        if (smsCooldownTimer) {
          clearInterval(smsCooldownTimer)
          smsCooldownTimer = null
        }
      }
    }, 1000)
  }

  const handleSendSmsCode = async () => {
    if (!/^1[3-9]\d{9}$/.test(phoneForm.phone)) {
      uni.showToast({ title: '请输入有效手机号', icon: 'none' })
      return
    }
    if (smsSending.value || smsCooldown.value > 0) return

    smsSending.value = true
    try {
      const res = await sendSmsCode({ phone: phoneForm.phone, scene: 'login' })
      const message = (res.data as any)?.message || '验证码已发送'
      const masterCode = (res.data as any)?.masterCode
      if (masterCode) {
        uni.showToast({ title: `${message}（${masterCode}）`, icon: 'none' })
      } else {
        uni.showToast({ title: message, icon: 'success' })
      }
      startSmsCooldown(60)
    } catch (error: any) {
      const msg = error?.message || '发送验证码失败'
      uni.showToast({ title: msg, icon: 'none' })
    } finally {
      smsSending.value = false
    }
  }

  const handlePhoneLogin = async () => {
    if (loginLoading.value) return
    if (!canPhoneLogin.value) {
      uni.showToast({ title: '请输入手机号和6位验证码', icon: 'none' })
      return
    }

    loginLoading.value = true
    try {
      const res = await loginWithPhone({
        phone: phoneForm.phone.trim(),
        code: phoneForm.code.trim(),
      })
      const data = res.data as any
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
    } catch (error: any) {
      uni.showToast({ title: error?.message || '登录失败，请稍后重试', icon: 'none' })
    } finally {
      loginLoading.value = false
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

  onUnmounted(() => {
    if (smsCooldownTimer) {
      clearInterval(smsCooldownTimer)
      smsCooldownTimer = null
    }
  })

  return {
    // 状态
    loginLoading,
    wechatLoading,
    smsSending,
    showAccountLogin,
    loginMode,
    smsCooldown,

    // 数据
    form,
    phoneForm,
    canLogin,
    canPhoneLogin,
    canSendSmsCode,

    // 方法
    handleLogin,
    handlePhoneLogin,
    handleSendSmsCode,
    handleWechatLogin,
    goToRegister,
    goToAgreement,
  }
}
