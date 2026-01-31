/**
 * 首页专用 Hook
 */
import { ref, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { useNetwork } from '@/composables/useNetwork'
import { getHotGames, type Game } from '@/api/game'
import { getPlayerList, type PlayerInfo } from '@/api/publicPlayer'
import { 
  getCachedHotGames, saveCachedHotGames, 
  getCachedPlayers, saveCachedPlayers 
} from '@/utils/offlineData'
import type { HotGameData } from '@/components/HotGamesScroll/index.vue'
import type { RecommendPlayerData } from '@/components/RecommendPlayersSection/index.vue'

export function useHome() {
  const userStore = useUserStore()
  const { isOnline } = useNetwork()
  
  // 离线模式状态
  const isOfflineMode = ref(false)
  
  // 热门游戏
  const hotGames = ref<HotGameData[]>([])
  const gamesLoading = ref(true)
  
  // 推荐陪玩师
  const recommendPlayers = ref<RecommendPlayerData[]>([])
  const playersLoading = ref(true)
  const noMorePlayers = ref(false)
  const playerPage = ref(1)
  
  // 计算属性
  const isLoggedIn = computed(() => userStore.isLoggedIn)
  const userInfo = computed(() => userStore.userInfo)
  
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
      const cachedGames = getCachedHotGames()
      if (cachedGames.length > 0) {
        hotGames.value = cachedGames
        isOfflineMode.value = true
      }
    } finally {
      gamesLoading.value = false
    }
  }
  
  // 加载推荐陪玩师
  const loadRecommendPlayers = async (refresh = true) => {
    if (refresh) {
      playerPage.value = 1
      noMorePlayers.value = false
      recommendPlayers.value = []
    }
    
    playersLoading.value = true
    try {
      const res = await getPlayerList({
        page: playerPage.value,
        pageSize: 10,
        sortBy: 'rating',
      })
      
      const resData = res.data as any
      const playerList = resData?.players || resData || []
      const players = playerList.map((p: PlayerInfo): RecommendPlayerData => ({
        id: p.id,
        nickname: p.nickname,
        avatar: p.avatar,
        rank: p.tags?.[0] || '',
        rating: p.rating || 5.0,
        hourlyRate: 2000,
        isOnline: p.isOnline,
        orderCount: p.orderCount || 0,
        mainGame: '',
      }))
      
      if (refresh) {
        recommendPlayers.value = players
        if (players.length > 0) {
          saveCachedPlayers(players)
          isOfflineMode.value = false
        }
      } else {
        recommendPlayers.value.push(...players)
      }
      
      if (players.length < 10) {
        noMorePlayers.value = true
      }
      
      playerPage.value++
    } catch (error) {
      console.error('加载推荐陪玩师失败', error)
      if (refresh && recommendPlayers.value.length === 0) {
        const cachedPlayers = getCachedPlayers()
        if (cachedPlayers.length > 0) {
          recommendPlayers.value = cachedPlayers
          isOfflineMode.value = true
          noMorePlayers.value = true
        }
      }
    } finally {
      playersLoading.value = false
    }
  }
  
  // 加载更多
  const loadMorePlayers = () => {
    if (playersLoading.value || noMorePlayers.value) return
    loadRecommendPlayers(false)
  }
  
  // 刷新所有
  const refreshAll = () => {
    isOfflineMode.value = false
    loadHotGames()
    loadRecommendPlayers()
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
  const init = () => {
    loadHotGames()
    loadRecommendPlayers()
  }
  
  return {
    // 状态
    isOfflineMode,
    isLoggedIn,
    userInfo,
    
    // 游戏
    hotGames,
    gamesLoading,
    
    // 陪玩师
    recommendPlayers,
    playersLoading,
    noMorePlayers,
    
    // 方法
    loadMorePlayers,
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
