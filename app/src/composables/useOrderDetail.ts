/**
 * 订单详情专用 Hook
 */
import { ref, computed, reactive } from 'vue'
import { getOrderDetail, cancelOrder as cancelOrderApi, completeOrder as completeOrderApi, submitReview as submitReviewApi, type OrderDetail as ApiOrderDetail } from '@/api/order'
import { payOrder } from '@/api/wallet'
import { confirmDialog } from '@/composables/useConfirmDialog'
import type { FeeItem, InfoItem, OrderActionKey } from '@/types/order'
import { formatDateTimeSafe } from '@/utils/format'
import type { OrderDetailData } from '@/types/order'
import { normalizeOrderStatus } from '@/components/OrderCard/utils'

export function useOrderDetail() {
  // 状态
  const loading = ref(true)
  const showReviewModal = ref(false)
  const reviewLoading = ref(false)
  const orderId = ref<number>(0)
  const countdown = ref(0)
  
  // 数据
  const order = ref<OrderDetailData>({} as OrderDetailData)
  
  // 评价表单
  const reviewForm = reactive({
    rating: 5,
    tags: [] as string[],
    content: '',
  })
  
  const reviewTags = ['技术过硬', '声音好听', '态度温柔', '有耐心', '准时上线', '沟通顺畅']
  
  // 服务信息
  const serviceInfo = computed((): InfoItem[] => {
    const o = order.value
    const items: InfoItem[] = [
      { label: '游戏', value: o.gameName || '-' },
      { label: '服务类型', value: o.serviceName || '-' },
      { label: '数量', value: `${o.quantity || 0}${o.unit || '局'}` },
    ]
    if (o.gameAccount) items.push({ label: '游戏账号', value: o.gameAccount })
    if (o.remark) items.push({ label: '备注', value: o.remark })
    return items
  })
  
  // 订单信息
  const orderInfo = computed((): InfoItem[] => {
    const o = order.value
    const items: InfoItem[] = [
      { label: '订单编号', value: o.orderNo || '-', copyable: true },
      { label: '创建时间', value: formatDateTimeSafe(o.createdAt) },
    ]
    if (o.scheduledStart) items.push({ label: '预约时间', value: formatDateTimeSafe(o.scheduledStart) })
    if (o.paidAt) items.push({ label: '支付时间', value: formatDateTimeSafe(o.paidAt) })
    if (o.startedAt) items.push({ label: '服务开始', value: formatDateTimeSafe(o.startedAt) })
    if (o.completedAt) items.push({ label: '完成时间', value: formatDateTimeSafe(o.completedAt) })
    return items
  })
  
  // 费用明细
  const feeItems = computed((): FeeItem[] => {
    const o = order.value
    const items: FeeItem[] = [
      { label: '服务费用', value: o.serviceFee || 0 },
    ]
    if (o.couponDiscount > 0) items.push({ label: '优惠券抵扣', value: o.couponDiscount, isDiscount: true })
    if (o.vipDiscount > 0) items.push({ label: 'VIP 折扣', value: o.vipDiscount, isDiscount: true })
    return items
  })
  
  // 加载订单详情
  const loadOrderDetail = async (id: number) => {
    orderId.value = id
    loading.value = true
    
    try {
      const res = await getOrderDetail(id)
      const data = res.data as ApiOrderDetail
      
      order.value = {
        id: data.id,
        orderNo: data.orderNo,
        status: normalizeOrderStatus(data.status, 'user'),
        player: {
          id: data.playerId,
          nickname: data.playerNickname || '陪玩师',
          avatar: data.playerAvatar,
          rating: 5.0,
          orderCount: 0,
        },
        gameName: data.gameName || '游戏',
        serviceName: data.serviceType || '陪玩服务',
        quantity: data.quantity || 1,
        unit: data.unit || '局',
        gameAccount: data.gameAccount,
        remark: data.remark,
        serviceFee: (data.totalCents || 0) / 100,
        couponDiscount: 0,
        vipDiscount: 0,
        totalAmount: (data.totalCents || 0) / 100,
        paymentMethod: data.paymentMethod,
        createdAt: data.createdAt,
        paidAt: data.paidAt,
        startedAt: data.startedAt,
        completedAt: data.completedAt,
        scheduledStart: data.scheduledStart,
        review: data.review ? {
          id: data.review.id,
          rating: data.review.rating,
          content: data.review.content,
          images: data.review.images,
          createdAt: data.review.createdAt,
        } : undefined,
      }
      
      // 启动倒计时
      if (order.value.status === 'pending') {
        startCountdown()
      }
    } catch (error: any) {
      console.error('加载订单详情失败', error)
      uni.showToast({ title: error?.message || '加载失败', icon: 'none' })
    } finally {
      loading.value = false
    }
  }
  
  // 启动倒计时
  const startCountdown = () => {
    const created = new Date(order.value.createdAt).getTime()
    const expiry = created + 30 * 60 * 1000 // 30分钟
    const now = Date.now()
    countdown.value = Math.max(0, Math.floor((expiry - now) / 1000))
    
    if (countdown.value > 0) {
      const timer = setInterval(() => {
        countdown.value--
        if (countdown.value <= 0) {
          clearInterval(timer)
          loadOrderDetail(orderId.value) // 刷新状态
        }
      }, 1000)
    }
  }
  
  // 操作处理
  const handleAction = async (action: OrderActionKey) => {
    switch (action) {
      case 'pay':
        await handlePay()
        break
      case 'cancel':
        await handleCancel()
        break
      case 'contact':
        goToChat()
        break
      case 'refund':
        handleRefund()
        break
      case 'complete':
        await handleComplete()
        break
      case 'review':
        showReviewModal.value = true
        break
      case 'reorder':
        handleReorder()
        break
    }
  }
  
  const handlePay = async () => {
    try {
      uni.showLoading({ title: '支付中...' })
      await payOrder(orderId.value, 'wallet')
      uni.hideLoading()
      uni.showToast({ title: '支付成功', icon: 'success' })
      loadOrderDetail(orderId.value)
    } catch (error: any) {
      uni.hideLoading()
      uni.showToast({ title: error?.message || '支付失败', icon: 'none' })
    }
  }
  
  const handleCancel = async () => {
    const confirmed = await confirmDialog({
      title: '确认取消',
      content: '确定要取消此订单吗？',
    })
    if (!confirmed) return
    try {
      await cancelOrderApi(orderId.value)
      uni.showToast({ title: '订单已取消', icon: 'success' })
      loadOrderDetail(orderId.value)
    } catch (error: any) {
      uni.showToast({ title: error?.message || '取消失败', icon: 'none' })
    }
  }
  
  const handleComplete = async () => {
    const confirmed = await confirmDialog({
      title: '确认完成',
      content: '确定服务已完成？',
    })
    if (!confirmed) return
    try {
      await completeOrderApi(orderId.value)
      uni.showToast({ title: '订单已完成', icon: 'success' })
      loadOrderDetail(orderId.value)
    } catch (error: any) {
      uni.showToast({ title: error?.message || '操作失败', icon: 'none' })
    }
  }
  
  const handleRefund = () => {
    uni.showToast({ title: '退款功能开发中', icon: 'none' })
  }
  
  const handleReorder = () => {
    uni.navigateTo({ url: `/pages/order/create/index?playerId=${order.value.player?.id}` })
  }
  
  // 评价标签切换
  const toggleTag = (tag: string) => {
    const index = reviewForm.tags.indexOf(tag)
    if (index > -1) {
      reviewForm.tags.splice(index, 1)
    } else {
      reviewForm.tags.push(tag)
    }
  }
  
  // 提交评价
  const submitReview = async () => {
    if (reviewForm.rating < 1) {
      uni.showToast({ title: '请选择评分', icon: 'none' })
      return
    }
    
    reviewLoading.value = true
    try {
      await submitReviewApi(orderId.value, {
        rating: reviewForm.rating,
        tags: reviewForm.tags,
        content: reviewForm.content,
      })
      uni.showToast({ title: '评价成功', icon: 'success' })
      showReviewModal.value = false
      loadOrderDetail(orderId.value)
    } catch (error: any) {
      uni.showToast({ title: error?.message || '提交失败', icon: 'none' })
    } finally {
      reviewLoading.value = false
    }
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  
  const goToPlayer = () => {
    if (order.value.player?.id) {
      uni.navigateTo({ url: `/pages/player/detail/index?id=${order.value.player.id}` })
    }
  }
  
  const goToChat = () => {
    if (order.value.player?.id) {
      uni.navigateTo({ url: `/pages/message/chat/index?playerId=${order.value.player.id}` })
    }
  }
  
  return {
    // 状态
    loading,
    showReviewModal,
    reviewLoading,
    countdown,
    order,
    
    // 计算属性
    serviceInfo,
    orderInfo,
    feeItems,
    
    // 评价
    reviewForm,
    reviewTags,
    toggleTag,
    submitReview,
    
    // 方法
    loadOrderDetail,
    handleAction,
    goBack,
    goToPlayer,
    goToChat,
  }
}
