/**
 * 创建订单专用 Hook
 */
import { ref, reactive, computed } from 'vue'
import { createOrder } from '@/api/order'
import { getPlayerDetail } from '@/api/publicPlayer'
import type { GameOption } from '@/types/game'
import type { Coupon } from '@/types/coupon'
import type { PlayerOrderInfo, PlayerServiceData } from '@/types/player'
import { formatYuan } from '@/utils/format'

export function useOrderCreate() {
  // 页面参数
  const playerId = ref(0)

  // 状态
  const loading = ref(true)
  const submitting = ref(false)
  const showDatePicker = ref(false)
  const showTimePicker = ref(false)
  const showCouponPicker = ref(false)

  // 陪玩师数据
  const player = reactive<PlayerOrderInfo & { games: GameOption[]; services: PlayerServiceData[] }>({
    id: 0,
    nickname: '',
    avatar: '',
    isOnline: false,
    rating: 5,
    games: [],
    services: [],
  })

  // 表单数据
  const selectedGameId = ref<number | undefined>()
  const selectedServiceId = ref<number | undefined>()
  const quantity = ref(1)
  const scheduledDate = ref('')
  const scheduledTime = ref('')
  const gameAccount = ref('')
  const remark = ref('')
  const selectedCoupon = ref<Coupon | null>(null)

  // 优惠券
  const availableCoupons = ref<Coupon[]>([])

  // 计算选中的服务
  const selectedService = computed(() =>
    player.services.find((s: PlayerServiceData) => s.id === selectedServiceId.value)
  )

  // 费用计算
  const serviceFee = computed(() => {
    if (!selectedService.value) return 0
    return selectedService.value.price * quantity.value
  })

  const vipDiscount = ref(0) // 根据用户 VIP 等级计算

  const totalFee = computed(() => {
    let fee = serviceFee.value
    if (selectedCoupon.value) {
      fee -= selectedCoupon.value.discount
    }
    fee -= vipDiscount.value
    return Math.max(0, fee)
  })

  // 表单校验
  const canSubmit = computed(() => {
    return selectedGameId.value &&
      selectedServiceId.value &&
      quantity.value > 0 &&
      scheduledDate.value &&
      scheduledTime.value
  })

  // 数量标题
  const quantityTitle = computed(() => {
    return selectedService.value?.unit === '小时' ? '时长（小时）' : '局数'
  })

  // 数量提示
  const quantityTip = computed(() => {
    const unit = selectedService.value?.unit || '局'
    const price = selectedService.value?.price || 0
    return `单价 ¥${price}/${unit}，共 ¥${formatYuan(serviceFee.value)}`
  })

  // 加载陪玩师详情
  const loadPlayerDetail = async (id: number) => {
    playerId.value = id
    loading.value = true

    try {
      const res = await getPlayerDetail(id)
      if (res.data) {
        const data = res.data as any
        player.id = data.id
        player.nickname = data.nickname
        player.avatar = data.avatar
        player.isOnline = data.isOnline
        player.rating = data.rating
        player.games = (data.gameRanks || []).map((g: any) => ({
          id: g.gameId,
          name: g.gameName,
          icon: g.gameIcon,
        }))
        player.services = (data.services || []).map((s: any) => ({
          id: s.id,
          name: s.serviceType,
          description: s.description,
          price: s.priceCents / 100,
          unit: s.unit || '局',
        }))

        // 默认选中第一个游戏和服务
        if (player.games.length > 0) {
          selectedGameId.value = player.games[0]?.id
        }
        if (player.services.length > 0) {
          selectedServiceId.value = player.services[0]?.id
        }
      }
    } catch (error) {
      console.error('加载陪玩师详情失败', error)
      uni.showToast({ title: '加载失败', icon: 'none' })
    } finally {
      loading.value = false
    }
  }

  // 选择优惠券
  const selectCoupon = (coupon: Coupon | null) => {
    selectedCoupon.value = coupon
    showCouponPicker.value = false
  }

  // 日期选择
  const onDateChange = (date: string) => {
    scheduledDate.value = date
    showDatePicker.value = false
  }

  // 时间选择
  const onTimeChange = (time: string) => {
    scheduledTime.value = time
    showTimePicker.value = false
  }

  // 提交订单
  const submitOrder = async () => {
    if (!canSubmit.value || submitting.value) return
    if (!selectedGameId.value || !selectedServiceId.value) return // Additional check for TS

    submitting.value = true

    try {
      const orderData = {
        playerId: playerId.value,
        gameId: selectedGameId.value,
        serviceId: selectedServiceId.value,
        quantity: quantity.value,
        scheduledStart: `${scheduledDate.value} ${scheduledTime.value}`,
        gameAccount: gameAccount.value || undefined,
        remark: remark.value || undefined,
        couponId: selectedCoupon.value?.id,
      }

      const res = await createOrder(orderData)

      uni.showToast({ title: '下单成功', icon: 'success' })

      // 跳转到支付页
      setTimeout(() => {
        uni.redirectTo({
          url: `/pages/order/pay/index?orderId=${(res.data as any).id}`
        })
      }, 500)
    } catch (error: any) {
      uni.showToast({ title: error?.message || '下单失败', icon: 'none' })
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
    showDatePicker,
    showTimePicker,
    showCouponPicker,

    // 数据
    player,
    selectedGameId,
    selectedServiceId,
    selectedService,
    quantity,
    scheduledDate,
    scheduledTime,
    gameAccount,
    remark,
    selectedCoupon,
    availableCoupons,

    // 计算
    serviceFee,
    vipDiscount,
    totalFee,
    canSubmit,
    quantityTitle,
    quantityTip,

    // 方法
    loadPlayerDetail,
    selectCoupon,
    onDateChange,
    onTimeChange,
    submitOrder,
    goBack,
  }
}
