/**
 * 陪玩师详情专用 Hook
 */
import { ref, computed } from 'vue'
import { getPlayerDetail, getPlayerReviews, type PlayerDetail as ApiPlayerDetail } from '@/api/publicPlayer'
import { addFavorite, removeFavorite, checkFavorite as checkFavoriteApi } from '@/api/favorite'
import { useUserStore } from '@/store/user'
import { setRedirectPath } from '@/utils/routeGuard'
import type { PageStateType } from '@/types/page'
import type { PlayerDetailData, PlayerGameData, PlayerServiceData } from '@/types/player'
import type { PlayerReviewData } from '@/types/review'

export function usePlayerDetail() {
  const userStore = useUserStore()
  
  // 页面状态
  const pageState = ref<PageStateType>('loading')
  const errorMessage = ref('')
  
  // 数据
  const playerId = ref<number>(0)
  const player = ref<PlayerDetailData>({} as PlayerDetailData)
  const isFavorite = ref(false)
  const selectedService = ref<PlayerServiceData | null>(null)
  
  // 评价预览
  const displayReviews = computed(() => player.value.reviews?.slice(0, 3) || [])
  
  // 加载详情
  const loadPlayerDetail = async (id: number) => {
    playerId.value = id
    pageState.value = 'loading'
    errorMessage.value = ''
    
    try {
      const res = await getPlayerDetail(id, { showError: false })
      const data = res.data as ApiPlayerDetail
      
      if (!data) {
        pageState.value = 'empty'
        return
      }
      
      player.value = {
        id: data.id,
        nickname: data.nickname,
        avatar: data.avatar,
        coverImage: data.photos?.[0] || data.avatar,
        signature: data.bio,
        gender: data.gender === 'unknown' ? undefined : data.gender,
        isOnline: data.isOnline,
        isVerified: data.isVerified,
        rating: data.rating ?? 5,
        orderCount: data.orderCount ?? 0,
        favoriteCount: data.followerCount,
        createdAt: data.createdAt,
        games: data.gameRanks?.map(g => ({
          id: g.gameId,
          name: g.gameName,
          icon: '',
          rankName: g.rankName,
          price: 20,
        })) || [],
        services: data.services?.map(s => ({
          id: s.id,
          name: s.serviceType,
          description: s.description,
          price: s.priceCents / 100,
          unit: s.unit,
        })) || [],
        reviews: [],
      }
      
      // 加载评价
      loadReviews(id)
      
      // 检查收藏状态
      checkFavoriteStatus()
      
      // 默认选中第一个服务
      if (player.value.services?.length) {
        selectedService.value = player.value.services[0] ?? null
      }
      
      pageState.value = 'content'
    } catch (error: any) {
      console.error('加载陪玩师详情失败', error)
      pageState.value = 'error'
      errorMessage.value = error?.message || '加载失败，请稍后重试'
    }
  }
  
  // 加载评价
  const loadReviews = async (id: number) => {
    try {
      const res = await getPlayerReviews(id, { page_size: 5 }, { showError: false })
      const reviews = res.data || []
      player.value.reviews = reviews.map((r: any): PlayerReviewData => ({
        id: r.id,
        userId: r.userId,
        userName: r.userName,
        userAvatar: r.userAvatar,
        rating: r.rating,
        content: r.content || '',
        images: r.images || [],
        createdAt: r.createdAt,
      }))
    } catch (error) {
      console.error('加载评价失败', error)
    }
  }
  
  // 检查收藏状态
  const checkFavoriteStatus = async () => {
    if (!userStore.isLoggedIn) return
    try {
      const res = await checkFavoriteApi(playerId.value)
      isFavorite.value = Boolean(res.data?.isFavorite ?? res.data?.isFavorited)
    } catch (error) {
      isFavorite.value = false
    }
  }
  
  // 切换收藏
  const toggleFavorite = async () => {
    if (!userStore.isLoggedIn) {
      setRedirectPath(`/pages/player/detail/index?id=${playerId.value}`)
      uni.navigateTo({ url: '/pages/auth/login/index' })
      return
    }
    
    try {
      if (isFavorite.value) {
        await removeFavorite(playerId.value)
        isFavorite.value = false
        uni.showToast({ title: '已取消收藏', icon: 'none' })
      } else {
        await addFavorite(playerId.value)
        isFavorite.value = true
        uni.showToast({ title: '收藏成功', icon: 'success' })
      }
    } catch (error: any) {
      uni.showToast({ title: error?.message || '操作失败', icon: 'none' })
    }
  }
  
  // 选择服务
  const selectService = (service: PlayerServiceData) => {
    selectedService.value = service
  }
  
  // 重试
  const handleRetry = () => {
    if (playerId.value) {
      loadPlayerDetail(playerId.value)
    }
  }
  
  // 导航
  const goBack = () => {
    uni.navigateBack()
  }
  
  const handleShare = () => {
    uni.showToast({ title: '分享功能开发中', icon: 'none' })
  }
  
  const goToReviews = () => {
    uni.navigateTo({ url: `/pages/review/list/index?playerId=${playerId.value}` })
  }
  
  const goToChat = () => {
    if (!userStore.isLoggedIn) {
      setRedirectPath(`/pages/player/detail/index?id=${playerId.value}`)
      uni.navigateTo({ url: '/pages/auth/login/index' })
      return
    }
    uni.navigateTo({ url: `/pages/message/chat/index?playerId=${playerId.value}` })
  }
  
  const goToOrder = () => {
    if (!player.value.isOnline) {
      uni.showToast({ title: '陪玩师当前离线', icon: 'none' })
      return
    }
    if (!userStore.isLoggedIn) {
      setRedirectPath(`/pages/player/detail/index?id=${playerId.value}`)
      uni.navigateTo({ url: '/pages/auth/login/index' })
      return
    }
    
    const serviceId = selectedService.value?.id
    uni.navigateTo({ 
      url: `/pages/order/create/index?playerId=${playerId.value}${serviceId ? `&serviceId=${serviceId}` : ''}`
    })
  }
  
  return {
    // 状态
    pageState,
    errorMessage,
    player,
    isFavorite,
    selectedService,
    displayReviews,
    
    // 方法
    loadPlayerDetail,
    handleRetry,
    toggleFavorite,
    selectService,
    goBack,
    handleShare,
    goToReviews,
    goToChat,
    goToOrder,
  }
}
