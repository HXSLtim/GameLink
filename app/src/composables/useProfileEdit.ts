/**
 * 编辑资料专用 Hook
 */
import { ref, reactive, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { getUserProfile, updateUserProfile, uploadAvatar } from '@/api/user'
import { formatDate, maskPhone } from '@/utils/format'
import type { ProfileForm, ProfileContactInfo } from '@/types/profile'
import type { ProfileGender, ProfileGenderValue } from '@/types/common'

export function useProfileEdit() {
  const userStore = useUserStore()
  
  // 状态
  const saving = ref(false)
  const showGenderPicker = ref(false)
  const showBirthdayPicker = ref(false)
  const showRegionPicker = ref(false)
  const tempGender = ref<ProfileGenderValue>('')
  
  // 原始数据
  const originalForm = ref<ProfileForm | null>(null)
  
  // 表单数据
  const form = reactive<ProfileForm>({
    avatar: '',
    nickname: '',
    gender: '',
    birthday: '',
    region: '',
    bio: '',
    games: [],
  })
  
  // 联系方式（只读）
  const profile = reactive<ProfileContactInfo>({
    phone: '',
    wechatBound: false,
  })
  
  // 性别选项
  const genderOptions: Array<{ label: string; value: ProfileGender }> = [
    { label: '男', value: 'male' },
    { label: '女', value: 'female' },
    { label: '保密', value: 'unknown' },
  ]
  
  // 日期范围
  const minBirthday = new Date('1950-01-01').getTime()
  const maxBirthday = new Date().getTime()
  
  // 是否有修改
  const hasChanges = computed(() => {
    if (!originalForm.value) return false
    return JSON.stringify(form) !== JSON.stringify(originalForm.value)
  })
  
  // 获取性别文案
  const getGenderText = (gender: ProfileGenderValue) => {
    const option = genderOptions.find(o => o.value === gender)
    return option?.label || '请选择'
  }
  
  // 格式化手机号
  const formatPhone = (phone: string) => maskPhone(phone)
  
  // 加载资料
  const loadProfile = async () => {
    try {
      const res = await getUserProfile()
      if (res.data) {
        const data = res.data as any
        form.avatar = data.avatar || ''
        form.nickname = data.nickname || ''
        form.gender = data.gender || ''
        form.birthday = data.birthday || ''
        form.region = data.region || ''
        form.bio = data.bio || ''
        form.games = data.games || []
        
        profile.phone = data.phone || ''
        profile.wechatBound = !!data.wechatOpenId
        
        // 保存原始数据
        originalForm.value = { ...form, games: [...form.games] }
      }
    } catch (error) {
      console.error('加载资料失败', error)
    }
  }
  
  // 更换头像
  const handleAvatarUpload = async () => {
    if (!form.avatar.startsWith('http')) {
      // 上传新头像
      try {
        uni.showLoading({ title: '上传中...' })
        const res = await uploadAvatar(form.avatar)
        form.avatar = (res.data as any).url
        uni.hideLoading()
      } catch (error) {
        uni.hideLoading()
        uni.showToast({ title: '上传失败', icon: 'none' })
      }
    }
  }
  
  // 确认性别
  const confirmGender = () => {
    form.gender = tempGender.value
    showGenderPicker.value = false
  }
  
  // 打开性别选择
  const openGenderPicker = () => {
    tempGender.value = form.gender
    showGenderPicker.value = true
  }
  
  // 生日确认
  const onBirthdayConfirm = (e: any) => {
    const date = new Date(e.value)
  form.birthday = formatDate(date)
    showBirthdayPicker.value = false
  }
  
  // 保存资料
  const saveProfile = async () => {
    if (!hasChanges.value || saving.value) return
    
    saving.value = true
    
    try {
      const normalizedGender = form.gender === '' ? undefined : form.gender
      await updateUserProfile({
        avatar: form.avatar,
        nickname: form.nickname,
        gender: normalizedGender,
        birthday: form.birthday,
        region: form.region,
        bio: form.bio,
        games: form.games,
      })
      
      // 更新 store
      userStore.setUserInfo({
        ...userStore.userInfo,
        nickname: form.nickname,
        avatar: form.avatar,
      } as any)
      
      uni.showToast({ title: '保存成功', icon: 'success' })
      
      setTimeout(() => {
        uni.navigateBack()
      }, 500)
    } catch (error: any) {
      uni.showToast({ title: error?.message || '保存失败', icon: 'none' })
    } finally {
      saving.value = false
    }
  }
  
  // 更换手机
  const changePhone = () => {
    uni.navigateTo({ url: '/pages/settings/phone/index' })
  }
  
  // 绑定微信
  const bindWechat = () => {
    uni.showToast({ title: '微信绑定功能开发中', icon: 'none' })
  }
  
  // 编辑游戏
  const editGames = () => {
    uni.navigateTo({ url: '/pages/profile/games/index' })
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  
  return {
    // 状态
    saving,
    showGenderPicker,
    showBirthdayPicker,
    showRegionPicker,
    tempGender,
    
    // 数据
    form,
    profile,
    genderOptions,
    minBirthday,
    maxBirthday,
    hasChanges,
    
    // 方法
    loadProfile,
    handleAvatarUpload,
    getGenderText,
    formatPhone,
    openGenderPicker,
    confirmGender,
    onBirthdayConfirm,
    saveProfile,
    changePhone,
    bindWechat,
    editGames,
    goBack,
  }
}
