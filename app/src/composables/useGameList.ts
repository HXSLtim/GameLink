/**
 * 游戏列表专用 Hook
 */
import { ref, computed } from 'vue'
import { useListPage } from './useListPage'
import { getGames, getGameCategories, type Game as ApiGame, type GameCategory as ApiGameCategory } from '@/api/game'
import type { GameData } from '@/components/GameCardLarge/index.vue'

interface GameCategory {
  id: string
  name: string
}

export function useGameList() {
  // 搜索和筛选
  const searchKeyword = ref('')
  const currentCategory = ref('all')
  
  // 分类列表
  const categories = ref<GameCategory[]>([{ id: 'all', name: '全部' }])
  
  // 构建 API 参数
  const buildParams = () => {
    const params: Record<string, any> = {}
    if (currentCategory.value !== 'all') {
      params.categoryId = currentCategory.value
    }
    if (searchKeyword.value.trim()) {
      params.keyword = searchKeyword.value.trim()
    }
    return params
  }
  
  // 使用通用列表 Hook
  const listPage = useListPage<GameData>({
    fetchFn: async (params) => {
      const res = await getGames({
        page: params.page,
        pageSize: params.pageSize,
        ...buildParams(),
      }, { showError: false })
      return res
    },
    extractList: (data: any) => {
      const gameList: ApiGame[] = data?.games || data?.items || data || []
      return gameList.map((g): GameData => ({
        id: g.id,
        name: g.name,
        coverImage: g.coverImage || g.icon,
        categoryId: g.categoryId?.toString(),
        isHot: g.isHot || false,
        playerCount: g.playerCount || 0,
        minPrice: g.minPrice ? g.minPrice / 100 : 20,
        maxPrice: g.maxPrice ? g.maxPrice / 100 : 100,
      }))
    },
    pageSize: 20,
  })
  
  // 过滤后的游戏列表
  const filteredGames = computed(() => {
    if (!searchKeyword.value.trim()) return listPage.list.value
    const keyword = searchKeyword.value.toLowerCase()
    return listPage.list.value.filter(game => 
      game.name.toLowerCase().includes(keyword)
    )
  })
  
  // 加载分类
  const loadCategories = async () => {
    try {
      const res = await getGameCategories({ showError: false })
      const resData = res.data as any
      const cats = resData?.categories || resData || []
      categories.value = [
        { id: 'all', name: '全部' },
        ...cats.map((c: ApiGameCategory) => ({
          id: c.id.toString(),
          name: c.name,
        })),
      ]
    } catch (error) {
      console.error('加载分类失败', error)
    }
  }
  
  // 刷新列表
  const refreshList = () => {
    listPage.refresh(buildParams())
  }
  
  // 搜索
  const handleSearch = () => {
    refreshList()
  }
  
  // 清除搜索
  const clearSearch = () => {
    searchKeyword.value = ''
    refreshList()
  }
  
  // 选择分类
  const selectCategory = (id: string) => {
    currentCategory.value = id
    refreshList()
  }
  
  // 跳转游戏详情
  const goToGame = (game: GameData) => {
    uni.navigateTo({ url: `/pages/player/list/index?gameId=${game.id}` })
  }
  
  // 返回
  const goBack = () => {
    uni.navigateBack()
  }
  
  // 初始化
  const init = () => {
    loadCategories()
    refreshList()
  }
  
  return {
    // 列表数据
    games: listPage.list,
    filteredGames,
    pageState: listPage.pageState,
    errorMessage: listPage.errorMessage,
    loading: listPage.loading,
    loadingMore: listPage.loadingMore,
    noMore: listPage.noMore,
    
    // 筛选
    searchKeyword,
    currentCategory,
    categories,
    
    // 方法
    loadMore: listPage.loadMore,
    refresh: refreshList,
    handleSearch,
    clearSearch,
    selectCategory,
    goToGame,
    goBack,
    init,
  }
}
