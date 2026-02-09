/**
 * 首页专用 Hook
 */
import { ref, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { useNetwork } from '@/composables/useNetwork'
import { getHotGames, type Game } from '@/api/game'
import { getPlayerList, type PlayerInfo } from '@/api/publicPlayer'
import { getBanners as fetchBannersApi, type BannerItem } from '@/api/banner'
import { 
  getCachedHotGames, saveCachedHotGames, 
  getCachedPlayers, saveCachedPlayers 
} from '@/utils/offlineData'
import { mapCachedPlayerToCard, mapCardToCachedPlayer, mapPlayerInfoToCard } from '@/utils/playerMapper'
import type { HotGameData } from '@/types/game'
import type { RecommendPlayerData } from '@/types/player'
import type { HomeBannerItem } from '@/types/home'

/** 将后端 BannerItem 转换为前端 HomeBannerItem */
function mapBannerToHome(b: BannerItem): HomeBannerItem {
  return {
    id: b.id,
    type: b.type,
    image: b.imageUrl,
    link: b.link,
    previewImages: b.type === 'preview' ? [b.imageUrl] : undefined,
    title: b.title,
    description: b.description,
    actionText: b.actionText,
  }
}

export function useHome() {
  const userStore = useUserStore()
  const { isOnline } = useNetwork()
  
  // 离线模式：任一接口因网络原因失败即为 true
  const isOfflineMode = ref(false)
  // 记录每个接口是否加载失败（用于判断全部失败的场景）
  let loadFailCount = 0
  
  // Banner（从后端接口获取）
  const banners = ref<HomeBannerItem[]>([])
  const bannersLoading = ref(false)

  // 热门游戏
  const hotGames = ref<HotGameData[]>([])
  const gamesLoading = ref(true)
  
  // 推荐陪玩师（首页只显示8个，不分页）
  const recommendPlayers = ref<RecommendPlayerData[]>([])
  const playersLoading = ref(true)
  
  // 计算属性
  const isLoggedIn = computed(() => userStore.isLoggedIn)
  const userInfo = computed(() => userStore.userInfo)
  
  // 加载 Banner
  const loadBanners = async () => {
    bannersLoading.value = true
    try {
      const res = await fetchBannersApi()
      const resData = res.data as any
      const items: BannerItem[] = resData?.banners || resData || []
      banners.value = items.map(mapBannerToHome)
    } catch (error) {
      console.error('加载 banner 失败', error)
      loadFailCount++
      if (!isOnline.value) {
        isOfflineMode.value = true
      }
      // 接口失败时使用本地静态兜底数据
      if (banners.value.length === 0) {
        banners.value = [
          {
            id: 1,
            type: 'link',
            image: '/static/images/banner-jump.svg',
            link: '/pages/game/list/index',
            title: '探索热门游戏',
            description: '一键发现优质陪玩师，畅享游戏乐趣',
            actionText: '立即前往',
          },
          {
            id: 2,
            type: 'preview',
            image: '/static/images/banner-preview.svg',
            previewImages: ['/static/images/banner-preview.svg'],
            title: '新赛季展示',
            description: '查看最新活动海报与精彩内容',
            actionText: '查看详情',
          },
        ]
      }
    } finally {
      bannersLoading.value = false
    }
  }

  // 加载热门游戏
  const loadHotGames = async () => {
    gamesLoading.value = true
    try {
      const res = await getHotGames(10)
      const resData = res.data as any
      const games = resData?.games || resData || []
      hotGames.value = games.map((g: Game): HotGameData => ({
        id: g.id,
        name: g.name,
        icon: g.icon,
        playerCount: g.playerCount || Math.floor(Math.random() * 1000) + 100,
      }))
      
      if (hotGames.value.length > 0) {
        saveCachedHotGames(hotGames.value)
        isOfflineMode.value = false
      }
    } catch (error) {
      console.error('加载热门游戏失败', error)
      loadFailCount++
      const cachedGames = getCachedHotGames()
      if (cachedGames.length > 0) {
        hotGames.value = cachedGames
      }
      // 无论有无缓存，只要网络不可用就标记离线
      if (!isOnline.value) {
        isOfflineMode.value = true
      }
    } finally {
      gamesLoading.value = false
    }
  }
  
  // 加载推荐陪玩师（首页只加载8个）
  const loadRecommendPlayers = async () => {
    playersLoading.value = true
    try {
      const res = await getPlayerList({
        page: 1,
        page_size: 8,
        sortBy: 'rating',
      })
      
      const resData = res.data as any
      const playerList = resData?.players || resData || []
      const players = playerList.map((p: PlayerInfo): RecommendPlayerData => mapPlayerInfoToCard(p))
      
      recommendPlayers.value = players
      if (players.length > 0) {
        saveCachedPlayers(players.map(mapCardToCachedPlayer))
        isOfflineMode.value = false
      }
    } catch (error) {
      console.error('加载推荐陪玩师失败', error)
      loadFailCount++
      if (recommendPlayers.value.length === 0) {
        const cachedPlayers = getCachedPlayers()
        if (cachedPlayers.length > 0) {
          recommendPlayers.value = cachedPlayers.slice(0, 8).map(mapCachedPlayerToCard)
        }
      }
      if (!isOnline.value) {
        isOfflineMode.value = true
      }
    } finally {
      playersLoading.value = false
    }
  }
  
  // 刷新所有
  const refreshAll = async () => {
    isOfflineMode.value = false
    loadFailCount = 0
    await Promise.allSettled([
      loadBanners(),
      loadHotGames(),
      loadRecommendPlayers(),
    ])
    if (loadFailCount >= 3) {
      isOfflineMode.value = true
    }
  }
  
  // 导航
  const goToLogin = () => {
    uni.navigateTo({ url: '/pages/auth/login/index' })
  }
  
  const goToProfile = () => {
    uni.navigateTo({ url: '/pages/profile/index/index' })
  }
  
  const goToPlayerList = () => {
    uni.navigateTo({ url: '/pages/player/list/index' })
  }
  
  const goToGameList = () => {
    uni.navigateTo({ url: '/pages/game/list/index' })
  }
  
  const goToGamePlayers = (gameId: number) => {
    uni.navigateTo({ url: `/pages/player/list/index?gameId=${gameId}` })
  }
  
  const goToPlayerDetail = (playerId: number) => {
    uni.navigateTo({ url: `/pages/player/detail/index?id=${playerId}` })
  }
  
  // 初始化
  const init = async () => {
    loadFailCount = 0
    await Promise.allSettled([
      loadBanners(),
      loadHotGames(),
      loadRecommendPlayers(),
    ])
    // 所有接口都失败时（即使浏览器认为在线），也标记为离线
    if (loadFailCount >= 3) {
      isOfflineMode.value = true
    }
  }
  
  return {
    // 状态
    isOfflineMode,
    isLoggedIn,
    userInfo,
    
    // Banner
    banners,

    // 游戏
    hotGames,
    gamesLoading,
    
    // 陪玩师
    recommendPlayers,
    playersLoading,
    
    // 方法
    refreshAll,
    goToLogin,
    goToProfile,
    goToPlayerList,
    goToGameList,
    goToGamePlayers,
    goToPlayerDetail,
    init,
  }
}
