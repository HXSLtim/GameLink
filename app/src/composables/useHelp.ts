/**
 * 帮助中心专用 Hook
 */
import { ref, computed } from 'vue'
import type { HelpCategory, HelpFaq } from '@/types/help'

export function useHelp() {
  // 状态
  const searchKeyword = ref('')
  const selectedCategory = ref<string | null>(null)
  const expandedId = ref<number | null>(null)
  
  // 分类（icon 为 uv-icon 名称）
  const categories = ref<HelpCategory[]>([
    { id: 'order', name: '订单相关', icon: 'list' },
    { id: 'payment', name: '支付问题', icon: 'wallet' },
    { id: 'account', name: '账号安全', icon: 'lock' },
    { id: 'service', name: '服务规则', icon: 'file-text' },
    { id: 'player', name: '陪玩入驻', icon: 'grid-fill' },
    { id: 'other', name: '其他问题', icon: 'question-circle' },
  ])
  
  // FAQ 数据
  const faqs = ref<HelpFaq[]>([
    { id: 1, categoryId: 'order', question: '如何下单预约陪玩师？', answer: '1. 在陪玩师列表中选择心仪的陪玩师\n2. 进入详情页查看服务项目\n3. 点击"立即下单"按钮\n4. 选择游戏、服务类型和时长\n5. 确认订单信息并完成支付' },
    { id: 2, categoryId: 'order', question: '订单如何取消？', answer: '待支付订单可直接取消。\n已支付订单：\n- 服务开始前可申请全额退款\n- 服务进行中需联系客服协商处理\n- 订单完成后无法取消' },
    { id: 3, categoryId: 'order', question: '如何评价陪玩师？', answer: '订单完成后，您可以在订单详情页对陪玩师进行评价。评价包括星级评分和文字描述，帮助其他用户更好地选择陪玩师。' },
    { id: 4, categoryId: 'payment', question: '支持哪些支付方式？', answer: '目前支持以下支付方式：\n- 微信支付\n- 支付宝\n- 平台余额支付\n- 组合支付（余额+第三方）' },
    { id: 5, categoryId: 'payment', question: '充值后可以退款吗？', answer: '充值金额原则上不支持退款，仅限在平台内消费使用。如遇特殊情况，可联系客服申请处理。' },
    { id: 6, categoryId: 'account', question: '如何修改密码？', answer: '进入「我的」-「设置」-「账号安全」-「修改密码」，验证手机号后即可设置新密码。' },
    { id: 7, categoryId: 'player', question: '如何成为陪玩师？', answer: '1. 在个人中心点击「成为陪玩师」\n2. 填写基本信息和游戏段位\n3. 上传身份证照片进行实名认证\n4. 等待平台审核（1-3个工作日）\n5. 审核通过后即可开始接单' },
  ])
  
  // 过滤后的 FAQ
  const displayFaqs = computed(() => {
    let result = faqs.value
    
    if (selectedCategory.value) {
      result = result.filter(f => f.categoryId === selectedCategory.value)
    }
    
    if (searchKeyword.value.trim()) {
      const keyword = searchKeyword.value.toLowerCase()
      result = result.filter(f => 
        f.question.toLowerCase().includes(keyword) ||
        f.answer.toLowerCase().includes(keyword)
      )
    }
    
    return result
  })
  
  // 选择分类
  const selectCategory = (id: string) => {
    selectedCategory.value = selectedCategory.value === id ? null : id
    expandedId.value = null
  }
  
  // 展开/收起 FAQ
  const toggleFaq = (id: number) => {
    expandedId.value = expandedId.value === id ? null : id
  }
  
  // 搜索
  const handleSearch = () => {
    selectedCategory.value = null
    expandedId.value = null
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  const goToService = () => uni.navigateTo({ url: '/pages/service/index' })
  
  return {
    // 数据
    searchKeyword,
    selectedCategory,
    expandedId,
    categories,
    displayFaqs,
    
    // 方法
    selectCategory,
    toggleFaq,
    handleSearch,
    goBack,
    goToService,
  }
}
