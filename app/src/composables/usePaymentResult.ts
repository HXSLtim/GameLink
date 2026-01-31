/**
 * 支付结果专用 Hook
 */
import { ref, computed } from 'vue'
import { getOrderPaymentStatus } from '@/api/order'

type ResultType = 'success' | 'pending' | 'failed'

interface OrderInfo {
  orderId?: number
  orderNo?: string
  amount?: number
  method?: string
  paidAt?: string
}

export function usePaymentResult() {
  // 状态
  const resultType = ref<ResultType>('success')
  const orderInfo = ref<OrderInfo>({})
  
  // 结果文案
  const resultTitle = computed(() => {
    const titles: Record<ResultType, string> = {
      success: '支付成功',
      pending: '支付处理中',
      failed: '支付失败',
    }
    return titles[resultType.value]
  })
  
  const resultDesc = computed(() => {
    const descs: Record<ResultType, string> = {
      success: '感谢您的支付，订单已生效',
      pending: '正在确认支付结果，请稍候...',
      failed: '支付未完成，请重新尝试',
    }
    return descs[resultType.value]
  })
  
  const paymentMethodText = computed(() => {
    const methods: Record<string, string> = {
      wechat: '微信支付',
      alipay: '支付宝',
      wallet: '余额支付',
      combined: '组合支付',
    }
    return methods[orderInfo.value.method || ''] || '在线支付'
  })
  
  // 格式化时间
  const formatTime = (dateStr?: string) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
  }
  
  // 初始化
  const init = (options: any) => {
    resultType.value = (options?.status as ResultType) || 'success'
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
        const status = (res.data as any)?.status
        
        if (status === 'paid') {
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
    uni.setClipboardData({
      data: orderInfo.value.orderNo,
      success: () => uni.showToast({ title: '已复制', icon: 'success' }),
    })
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
