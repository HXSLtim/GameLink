/**
 * 支付结果专用 Hook
 */
import { ref, computed } from 'vue'
import { getOrderPaymentStatus } from '@/api/order'
import { formatDateTimeSafe } from '@/utils/format'
import { copyToClipboard } from '@/utils'
import type { ResultStatusType } from '@/types/status'
import type { OrderPaymentInfo, OrderPaymentMethod, OrderPaymentStatus } from '@/types/order'

export function usePaymentResult() {
  // 状态
  const resultType = ref<ResultStatusType>('success')
  const orderInfo = ref<OrderPaymentInfo>({})
  
  // 结果文案
  const resultTitle = computed(() => {
    const titles: Record<ResultStatusType, string> = {
      success: '支付成功',
      pending: '支付处理中',
      failed: '支付失败',
      warning: '支付异常',
    }
    return titles[resultType.value]
  })
  
  const resultDesc = computed(() => {
    const descs: Record<ResultStatusType, string> = {
      success: '感谢您的支付，订单已生效',
      pending: '正在确认支付结果，请稍候...',
      failed: '支付未完成，请重新尝试',
      warning: '支付状态异常，请稍后重试',
    }
    return descs[resultType.value]
  })
  
  const paymentMethodText = computed(() => {
    const methods: Record<OrderPaymentMethod, string> = {
      wechat: '微信支付',
      alipay: '支付宝',
      wallet: '余额支付',
      combined: '组合支付',
    }
    const method = orderInfo.value.method
    return (method && methods[method]) || '在线支付'
  })
  
  // 格式化时间
  const formatTime = (dateStr?: string) => formatDateTimeSafe(dateStr)
  
  // 初始化
  const init = (options: any) => {
    resultType.value = (options?.status as ResultStatusType) || 'success'
    orderInfo.value = {
      orderId: options?.orderId ? parseInt(options.orderId) : undefined,
      orderNo: options?.orderNo,
      amount: options?.amount ? parseInt(options.amount) : undefined,
      method: options?.method,
      paidAt: options?.paidAt || new Date().toISOString(),
    }
    
    // 如果是 pending，可以轮询查询状态
    if (resultType.value === 'pending' && orderInfo.value.orderId) {
      pollPaymentStatus()
    }
  }
  
  // 轮询支付状态
  const pollPaymentStatus = async () => {
    let attempts = 0
    const maxAttempts = 10
    
    const check = async () => {
      attempts++
      try {
        const res = await getOrderPaymentStatus(orderInfo.value.orderId!)
        const status = (res.data as { status?: OrderPaymentStatus })?.status
        
        if (status === 'paid' || status === 'success') {
          resultType.value = 'success'
          return
        } else if (status === 'failed') {
          resultType.value = 'failed'
          return
        }
      } catch (error) {
        console.error('查询支付状态失败', error)
      }
      
      if (attempts < maxAttempts) {
        setTimeout(check, 2000)
      } else {
        resultType.value = 'failed'
      }
    }
    
    setTimeout(check, 2000)
  }
  
  // 复制订单号
  const copyOrderNo = () => {
    if (!orderInfo.value.orderNo) return
    copyToClipboard(orderInfo.value.orderNo).catch(() => undefined)
  }
  
  // 导航
  const goToOrderDetail = () => {
    if (orderInfo.value.orderId) {
      uni.redirectTo({ url: `/pages/order/detail/index?id=${orderInfo.value.orderId}` })
    } else {
      uni.redirectTo({ url: '/pages/order/list/index' })
    }
  }
  
  const goToHome = () => uni.switchTab({ url: '/pages/index/index' })
  
  return {
    // 状态
    resultType,
    orderInfo,
    resultTitle,
    resultDesc,
    paymentMethodText,
    
    // 方法
    formatTime,
    init,
    copyOrderNo,
    goToOrderDetail,
    goToHome,
  }
}
