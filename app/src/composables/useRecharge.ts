/**
 * 充值专用 Hook
 */
import { ref, computed } from 'vue'
import { getWalletInfo, recharge } from '@/api/wallet'
import type { AmountOption } from '@/components/AmountSelector/index.vue'
import type { PaymentMethod } from '@/components/PaymentMethodSelector/index.vue'

export function useRecharge() {
  // 状态
  const submitting = ref(false)
  const agreeTerms = ref(false)
  const currentBalance = ref(0)
  
  // 金额
  const selectedAmount = ref(0)
  const customAmount = ref('')
  
  const amountOptions: AmountOption[] = [
    { value: 10 },
    { value: 30 },
    { value: 50 },
    { value: 100, bonus: 5 },
    { value: 200, bonus: 15 },
    { value: 500, bonus: 50 },
  ]
  
  // 支付方式
  const selectedMethod = ref('wechat')
  
  const paymentMethods: PaymentMethod[] = [
    { value: 'wechat', name: '微信支付', icon: '💚', enabled: true },
    { value: 'alipay', name: '支付宝', icon: '💙', enabled: true },
    { value: 'apple', name: 'Apple Pay', icon: '🍎', enabled: false, tip: '暂不支持' },
  ]
  
  // 计算最终金额
  const finalAmount = computed(() => {
    if (selectedAmount.value > 0) {
      return selectedAmount.value
    }
    return parseFloat(customAmount.value) || 0
  })
  
  // 赠送金额
  const bonusAmount = computed(() => {
    const option = amountOptions.find(o => o.value === selectedAmount.value)
    return option?.bonus || 0
  })
  
  // 是否可提交
  const canSubmit = computed(() => {
    return finalAmount.value >= 10 && 
           finalAmount.value <= 10000 && 
           selectedMethod.value && 
           agreeTerms.value
  })
  
  // 格式化余额
  const formatBalance = (cents: number) => (cents / 100).toFixed(2)
  
  // 加载余额
  const loadBalance = async () => {
    try {
      const res = await getWalletInfo()
      currentBalance.value = (res.data as any)?.balanceCents || 0
    } catch (error) {
      console.error('加载余额失败', error)
    }
  }
  
  // 选择金额
  const selectAmount = (amount: number) => {
    selectedAmount.value = amount
    customAmount.value = ''
  }
  
  // 提交充值
  const submitRecharge = async () => {
    if (!canSubmit.value || submitting.value) return
    
    submitting.value = true
    
    try {
      const res = await recharge({
        amountCents: Math.round(finalAmount.value * 100),
        paymentMethod: selectedMethod.value,
      })
      
      // 调用支付
      const paymentData = res.data as any
      
      if (selectedMethod.value === 'wechat') {
        // #ifdef MP-WEIXIN
        uni.requestPayment({
          provider: 'wxpay',
          ...paymentData,
          success: () => {
            uni.showToast({ title: '充值成功', icon: 'success' })
            setTimeout(() => {
              uni.navigateBack()
            }, 1000)
          },
          fail: () => {
            uni.showToast({ title: '支付取消', icon: 'none' })
          }
        })
        // #endif
        // #ifdef H5
        uni.showToast({ title: '请在微信中完成支付', icon: 'none' })
        // #endif
      } else {
        uni.showToast({ title: '支付发起成功', icon: 'success' })
      }
    } catch (error: any) {
      uni.showToast({ title: error?.message || '充值失败', icon: 'none' })
    } finally {
      submitting.value = false
    }
  }
  
  // 查看协议
  const viewAgreement = () => {
    uni.navigateTo({ url: '/pages/agreement/index?type=recharge' })
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  
  // 初始化
  const init = () => {
    loadBalance()
    selectedAmount.value = amountOptions[0].value
  }
  
  return {
    // 状态
    submitting,
    agreeTerms,
    currentBalance,
    selectedAmount,
    customAmount,
    selectedMethod,
    
    // 数据
    amountOptions,
    paymentMethods,
    finalAmount,
    bonusAmount,
    canSubmit,
    
    // 方法
    formatBalance,
    selectAmount,
    submitRecharge,
    viewAgreement,
    goBack,
    init,
  }
}
