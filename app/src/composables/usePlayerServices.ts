/**
 * 陪玩师服务管理专用 Hook
 */
import { ref, computed, reactive } from 'vue'
import { getPlayerServices, createPlayerService, updatePlayerService, deletePlayerService } from '@/api/player'
import { confirmDialog } from '@/composables/useConfirmDialog'
import type { PlayerServiceCardData, PlayerServiceForm } from '@/types/player'
import type { StatItem } from '@/types/ui'

export function usePlayerServices() {
  // 状态
  const loading = ref(true)
  const showEditor = ref(false)
  const saving = ref(false)
  const editingService = ref<PlayerServiceCardData | null>(null)
  
  // 数据
  const services = ref<PlayerServiceCardData[]>([])
  
  // 表单
  const form = reactive<PlayerServiceForm>({
    gameId: undefined,
    gameName: '',
    serviceType: '',
    serviceName: '',
    rankId: undefined,
    rankName: '',
    price: 0,
    unit: '局',
    description: '',
  })
  
  // 统计
  const onlineCount = computed(() => services.value.filter(s => s.isOnline).length)
  const offlineCount = computed(() => services.value.filter(s => !s.isOnline).length)
  
  const statItems = computed((): StatItem[] => [
    { value: services.value.length, label: '全部服务' },
    { value: onlineCount.value, label: '已上架', highlight: true },
    { value: offlineCount.value, label: '已下架' },
  ])
  
  // 加载服务列表
  const loadServices = async () => {
    loading.value = true
    try {
      const res = await getPlayerServices()
      services.value = ((res.data as any)?.items || res.data || []).map((s: any): PlayerServiceCardData => ({
        id: s.id,
        gameId: s.gameId,
        gameName: s.gameName || '未知游戏',
        gameIcon: s.gameIcon,
        serviceName: s.serviceName || s.serviceType || '陪玩',
        price: (s.priceCents || 0) / 100,
        unit: s.unit || '局',
        rankName: s.rankName || '不限',
        description: s.description,
        isOnline: s.isOnline ?? true,
      }))
    } catch (error) {
      console.error('加载服务失败', error)
    } finally {
      loading.value = false
    }
  }
  
  // 打开添加
  const addService = () => {
    resetForm()
    editingService.value = null
    showEditor.value = true
  }
  
  // 打开编辑
  const editService = (service: PlayerServiceCardData) => {
    editingService.value = service
    form.gameId = service.gameId
    form.gameName = service.gameName
    form.serviceName = service.serviceName
    form.rankName = service.rankName
    form.price = service.price
    form.unit = service.unit
    form.description = service.description || ''
    showEditor.value = true
  }
  
  // 关闭编辑器
  const closeEditor = () => {
    showEditor.value = false
    resetForm()
  }
  
  // 重置表单
  const resetForm = () => {
    form.gameId = undefined
    form.gameName = ''
    form.serviceType = ''
    form.serviceName = ''
    form.rankId = undefined
    form.rankName = ''
    form.price = 0
    form.unit = '局'
    form.description = ''
  }
  
  // 保存服务
  const saveService = async () => {
    if (!form.gameId || !form.serviceName || form.price <= 0) {
      uni.showToast({ title: '请填写必填项', icon: 'none' })
      return
    }
    
    saving.value = true
    try {
      const data = {
        gameId: form.gameId,
        serviceName: form.serviceName,
        rankName: form.rankName,
        priceCents: Math.round(form.price * 100),
        unit: form.unit,
        description: form.description,
      }
      
      if (editingService.value) {
        await updatePlayerService(editingService.value.id, data)
        uni.showToast({ title: '更新成功', icon: 'success' })
      } else {
        await createPlayerService(data)
        uni.showToast({ title: '添加成功', icon: 'success' })
      }
      
      closeEditor()
      loadServices()
    } catch (error: any) {
      uni.showToast({ title: error?.message || '保存失败', icon: 'none' })
    } finally {
      saving.value = false
    }
  }
  
  // 切换状态
  const toggleStatus = async (service: PlayerServiceCardData) => {
    try {
      await updatePlayerService(service.id, { isOnline: !service.isOnline })
      service.isOnline = !service.isOnline
      uni.showToast({ title: service.isOnline ? '已上架' : '已下架', icon: 'none' })
    } catch (error) {
      uni.showToast({ title: '操作失败', icon: 'none' })
    }
  }
  
  // 删除服务
  const handleDelete = async (service: PlayerServiceCardData) => {
    const confirmed = await confirmDialog({
      title: '确认删除',
      content: `确定要删除「${service.serviceName}」服务吗？`,
    })
    if (!confirmed) return
    try {
      await deletePlayerService(service.id)
      services.value = services.value.filter(s => s.id !== service.id)
      uni.showToast({ title: '删除成功', icon: 'success' })
    } catch (error) {
      uni.showToast({ title: '删除失败', icon: 'none' })
    }
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  
  return {
    // 状态
    loading,
    showEditor,
    saving,
    editingService,
    
    // 数据
    services,
    form,
    statItems,
    onlineCount,
    offlineCount,
    
    // 方法
    loadServices,
    addService,
    editService,
    closeEditor,
    saveService,
    toggleStatus,
    handleDelete,
    goBack,
  }
}
