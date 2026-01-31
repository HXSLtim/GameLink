/**
 * 注册专用 Hook
 */
import { ref, reactive, computed } from 'vue'
import { useUserStore, normalizeUserInfo } from '@/store/user'
import { register } from '@/api/auth'
import { consumeRedirectPath, redirectToUrl } from '@/utils/routeGuard'
import type { RoleOption } from '@/components/RoleSelector/index.vue'

export function useRegister() {
  const userStore = useUserStore()
  
  // 状态
  const loading = ref(false)
  const agreed = ref(false)
  
  // 表单
  const form = reactive({
    phone: '',
    nickname: '',
    password: '',
    confirmPassword: '',
    role: 'user' as 'user' | 'player',
  })
  
  // 角色选项
  const roleOptions: RoleOption[] = [
    { value: 'user', icon: '👤', name: '普通用户', desc: '找陪玩、享受游戏乐趣' },
    { value: 'player', icon: '🎮', name: '陪玩师', desc: '提供服务、赚取收入' },
  ]
  
  // 是否可注册
  const canRegister = computed(() => {
    return form.phone.length === 11 &&
           form.nickname.length >= 2 &&
           form.password.length >= 6 &&
           form.password === form.confirmPassword &&
           agreed.value
  })
  
  // 表单校验
  const validateForm = () => {
    if (form.phone.length !== 11) {
      uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
      return false
    }
    if (form.nickname.length < 2) {
      uni.showToast({ title: '昵称至少2个字符', icon: 'none' })
      return false
    }
    if (form.password.length < 6 || form.password.length > 20) {
      uni.showToast({ title: '密码需要6-20位', icon: 'none' })
      return false
    }
    if (form.password !== form.confirmPassword) {
      uni.showToast({ title: '两次密码输入不一致', icon: 'none' })
      return false
    }
    if (!agreed.value) {
      uni.showToast({ title: '请先同意用户协议', icon: 'none' })
      return false
    }
    return true
  }
  
  // 注册
  const handleRegister = async () => {
    if (!validateForm() || loading.value) return
    
    loading.value = true
    
    try {
      const res = await register({
        phone: form.phone,
        nickname: form.nickname,
        password: form.password,
        role: form.role,
      })
      
      const data = res.data as any
      
      // 保存用户信息和 token
      userStore.setToken(data.accessToken)
      userStore.setUserInfo(normalizeUserInfo(data.user))
      
      uni.showToast({ title: '注册成功', icon: 'success' })
      
      // 重定向
      setTimeout(() => {
        handleRedirect()
      }, 500)
    } catch (error: any) {
      uni.showToast({ title: error?.message || '注册失败', icon: 'none' })
    } finally {
      loading.value = false
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
  const goBack = () => uni.navigateBack()
  const goToLogin = () => uni.navigateTo({ url: '/pages/auth/login/index' })
  const goToAgreement = (type: string) => uni.navigateTo({ url: `/pages/agreement/index?type=${type}` })
  
  return {
    // 状态
    loading,
    agreed,
    
    // 数据
    form,
    roleOptions,
    canRegister,
    
    // 方法
    handleRegister,
    goBack,
    goToLogin,
    goToAgreement,
  }
}
