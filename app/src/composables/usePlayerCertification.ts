/**
 * 陪玩认证专用 Hook
 */
import { ref, reactive, computed } from 'vue'
import { getCertificationStatus, submitCertification } from '@/api/player'
import type { GameCertData } from '@/components/GameCertItem/index.vue'

type CertStatus = 'none' | 'pending' | 'approved' | 'rejected'

interface CertForm {
  realName: string
  idNumber: string
  gender: string
  idCardFront: string
  idCardBack: string
  games: GameCertData[]
  introduction: string
  voiceSample?: string
  voiceDuration?: number
}

export function usePlayerCertification() {
  // 状态
  const loading = ref(true)
  const submitting = ref(false)
  const certStatus = ref<CertStatus>('none')
  const showGenderPicker = ref(false)
  const recording = ref(false)
  const isPlaying = ref(false)
  
  // 表单
  const form = reactive<CertForm>({
    realName: '',
    idNumber: '',
    gender: '',
    idCardFront: '',
    idCardBack: '',
    games: [],
    introduction: '',
    voiceSample: undefined,
    voiceDuration: undefined,
  })
  
  // 性别选项
  const genderOptions = [
    { label: '男', value: 'male' },
    { label: '女', value: 'female' },
  ]
  
  // 是否已认证
  const isApproved = computed(() => certStatus.value === 'approved')
  
  // 是否可提交
  const canSubmit = computed(() => {
    return form.realName && 
           form.idNumber && 
           form.gender &&
           form.idCardFront && 
           form.idCardBack && 
           form.games.length > 0 &&
           form.games.every(g => g.gameId && g.rankId && g.screenshot)
  })
  
  // 获取性别文案
  const getGenderText = (gender: string) => {
    const option = genderOptions.find(o => o.value === gender)
    return option?.label || '请选择'
  }
  
  // 加载认证状态
  const loadCertStatus = async () => {
    loading.value = true
    try {
      const res = await getCertificationStatus()
      const data = res.data as any
      certStatus.value = data?.status || 'none'
      
      if (data?.form) {
        form.realName = data.form.realName || ''
        form.idNumber = data.form.idNumber || ''
        form.gender = data.form.gender || ''
        form.idCardFront = data.form.idCardFront || ''
        form.idCardBack = data.form.idCardBack || ''
        form.games = data.form.games || []
        form.introduction = data.form.introduction || ''
        form.voiceSample = data.form.voiceSample
        form.voiceDuration = data.form.voiceDuration
      }
    } catch (error) {
      console.error('加载认证状态失败', error)
    } finally {
      loading.value = false
    }
  }
  
  // 添加游戏认证
  const addGameCert = () => {
    form.games.push({
      gameId: undefined,
      gameName: '',
      rankId: undefined,
      rankName: '',
      screenshot: undefined,
    })
  }
  
  // 移除游戏认证
  const removeGameCert = (index: number) => {
    form.games.splice(index, 1)
  }
  
  // 选择游戏（需要外部实现 Picker）
  const selectGame = (index: number, gameId: number, gameName: string) => {
    form.games[index].gameId = gameId
    form.games[index].gameName = gameName
  }
  
  // 选择段位（需要外部实现 Picker）
  const selectRank = (index: number, rankId: number, rankName: string) => {
    form.games[index].rankId = rankId
    form.games[index].rankName = rankName
  }
  
  // 更新截图
  const updateScreenshot = (index: number, url: string) => {
    form.games[index].screenshot = url
  }
  
  // 录音相关
  const startRecord = () => {
    recording.value = true
    // 实际实现需要调用 uni.getRecorderManager()
  }
  
  const stopRecord = () => {
    recording.value = false
    // 停止录音并保存
  }
  
  const playVoice = () => {
    isPlaying.value = !isPlaying.value
    // 播放/暂停语音
  }
  
  const deleteVoice = () => {
    form.voiceSample = undefined
    form.voiceDuration = undefined
  }
  
  // 提交认证
  const submitForm = async () => {
    if (!canSubmit.value || submitting.value) return
    
    submitting.value = true
    
    try {
      await submitCertification(form)
      certStatus.value = 'pending'
      uni.showToast({ title: '提交成功，请等待审核', icon: 'success' })
    } catch (error: any) {
      uni.showToast({ title: error?.message || '提交失败', icon: 'none' })
    } finally {
      submitting.value = false
    }
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  
  return {
    // 状态
    loading,
    submitting,
    certStatus,
    showGenderPicker,
    recording,
    isPlaying,
    
    // 数据
    form,
    genderOptions,
    isApproved,
    canSubmit,
    
    // 方法
    getGenderText,
    loadCertStatus,
    addGameCert,
    removeGameCert,
    selectGame,
    selectRank,
    updateScreenshot,
    startRecord,
    stopRecord,
    playVoice,
    deleteVoice,
    submitForm,
    goBack,
  }
}
