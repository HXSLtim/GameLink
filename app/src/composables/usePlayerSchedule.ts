import { ref } from 'vue'
import { getPlayerSchedule, updatePlayerSchedule, type ScheduleSlot } from '@/api/player'

const DAY_LABELS = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

function createDefaultSlots(): ScheduleSlot[] {
  return Array.from({ length: 7 }, (_, dayOfWeek) => ({
    dayOfWeek,
    startTime: '19:00',
    endTime: '23:00',
    isAvailable: dayOfWeek >= 1 && dayOfWeek <= 5,
  }))
}

function normalizeSlots(slots?: ScheduleSlot[]): ScheduleSlot[] {
  const source = Array.isArray(slots) ? slots : []
  const byDay = new Map<number, ScheduleSlot>()

  for (const slot of source) {
    if (typeof slot.dayOfWeek !== 'number' || slot.dayOfWeek < 0 || slot.dayOfWeek > 6) {
      continue
    }
    byDay.set(slot.dayOfWeek, {
      dayOfWeek: slot.dayOfWeek,
      startTime: slot.startTime || '19:00',
      endTime: slot.endTime || '23:00',
      isAvailable: Boolean(slot.isAvailable),
    })
  }

  return Array.from({ length: 7 }, (_, dayOfWeek) => {
    const existed = byDay.get(dayOfWeek)
    if (existed) return existed
    return {
      dayOfWeek,
      startTime: '19:00',
      endTime: '23:00',
      isAvailable: false,
    }
  })
}

export function usePlayerSchedule() {
  const loading = ref(true)
  const saving = ref(false)
  const autoOffline = ref(true)
  const timezone = ref('Asia/Shanghai')
  const slots = ref<ScheduleSlot[]>(createDefaultSlots())

  const getDayLabel = (dayOfWeek: number) => DAY_LABELS[dayOfWeek] || `周${dayOfWeek}`

  const loadSchedule = async () => {
    loading.value = true
    try {
      const res = await getPlayerSchedule()
      const data = res.data
      if (data) {
        autoOffline.value = Boolean(data.autoOffline)
        timezone.value = data.timezone || 'Asia/Shanghai'
        slots.value = normalizeSlots(data.slots)
      }
    } catch (error) {
      console.error('加载排班失败', error)
      slots.value = createDefaultSlots()
    } finally {
      loading.value = false
    }
  }

  const updateSlotAvailability = (index: number, value: boolean) => {
    const slot = slots.value[index]
    if (!slot) return
    slot.isAvailable = value
  }

  const updateSlotTime = (index: number, key: 'startTime' | 'endTime', value: string) => {
    const slot = slots.value[index]
    if (!slot) return
    slot[key] = value
  }

  const saveSchedule = async () => {
    if (saving.value) return
    saving.value = true
    try {
      await updatePlayerSchedule({
        autoOffline: autoOffline.value,
        timezone: timezone.value,
        slots: slots.value,
      })
      uni.showToast({ title: '保存成功', icon: 'success' })
    } catch (error: any) {
      uni.showToast({ title: error?.message || '保存失败', icon: 'none' })
    } finally {
      saving.value = false
    }
  }

  const goBack = () => uni.navigateBack()

  return {
    loading,
    saving,
    autoOffline,
    timezone,
    slots,
    getDayLabel,
    loadSchedule,
    updateSlotAvailability,
    updateSlotTime,
    saveSchedule,
    goBack,
  }
}
