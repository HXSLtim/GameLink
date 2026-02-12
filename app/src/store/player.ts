/**
 * 陪玩师状态管理
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { put } from '@/api/request'
import type { PlayerService, TodayStats, PlayerSchedule } from '@/api/player'
import type { DashboardStatus } from '@/types/status'

export const usePlayerStore = defineStore('player', () => {
  // 工作状态
  const workStatus = ref<DashboardStatus>('offline')
  
  // 今日统计
  const todayStats = ref<TodayStats | null>(null)
  
  // 服务列表
  const services = ref<PlayerService[]>([])
  
  // 排班设置
  const schedule = ref<PlayerSchedule | null>(null)
  
  // 待接订单数
  const pendingOrderCount = ref(0)
  
  // 计算属性
  const isOnline = computed(() => workStatus.value === 'online')
  const isBusy = computed(() => workStatus.value === 'busy')
  const hasOnlineService = computed(() => services.value.some(s => s.isOnline))
  
  /**
   * 设置工作状态
   */
  function setWorkStatus(status: DashboardStatus) {
    workStatus.value = status
  }
  
  /**
   * 上线
   */
  async function goOnline() {
    try {
      await put('/player/online-status', { status: 'online' })
      workStatus.value = 'online'
    } catch (e) {
      console.error('Failed to go online:', e)
    }
  }
  
  /**
   * 下线
   */
  async function goOffline() {
    try {
      await put('/player/online-status', { status: 'offline' })
      workStatus.value = 'offline'
    } catch (e) {
      console.error('Failed to go offline:', e)
    }
  }
  
  /**
   * 设置忙碌
   */
  function setBusy() {
    workStatus.value = 'busy'
  }
  
  /**
   * 更新今日统计
   */
  function updateTodayStats(stats: TodayStats) {
    todayStats.value = stats
  }
  
  /**
   * 设置服务列表
   */
  function setServices(list: PlayerService[]) {
    services.value = list
  }
  
  /**
   * 更新单个服务
   */
  function updateService(id: number, data: Partial<PlayerService>) {
    const index = services.value.findIndex(s => s.id === id)
    if (index > -1) {
      const current = services.value[index]
      if (!current) return
      services.value[index] = { ...current, ...data }
    }
  }
  
  /**
   * 删除服务
   */
  function removeService(id: number) {
    services.value = services.value.filter(s => s.id !== id)
  }
  
  /**
   * 添加服务
   */
  function addService(service: PlayerService) {
    services.value.push(service)
  }
  
  /**
   * 设置排班
   */
  function setSchedule(data: PlayerSchedule) {
    schedule.value = data
  }
  
  /**
   * 设置待接订单数
   */
  function setPendingOrderCount(count: number) {
    pendingOrderCount.value = count
  }
  
  /**
   * 重置状态
   */
  function reset() {
    workStatus.value = 'offline'
    todayStats.value = null
    services.value = []
    schedule.value = null
    pendingOrderCount.value = 0
  }
  
  return {
    // 状态
    workStatus,
    todayStats,
    services,
    schedule,
    pendingOrderCount,
    
    // 计算属性
    isOnline,
    isBusy,
    hasOnlineService,
    
    // 方法
    setWorkStatus,
    goOnline,
    goOffline,
    setBusy,
    updateTodayStats,
    setServices,
    updateService,
    removeService,
    addService,
    setSchedule,
    setPendingOrderCount,
    reset,
  }
})
