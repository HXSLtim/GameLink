/**
 * 游戏列表专用 Hook
 */
import { ref, computed } from 'vue'
import { useListPage } from './useListPage'
import { getGames, getGameCategories, type Game as ApiGame, type GameCategory as ApiGameCategory } from '@/api/game'
import type { GameCardData, GameTabItem } from '@/types/game'
import type { FilterSection, FilterValues } from '@/types/filter'

export function useGameList() {
  // 搜索和筛选
  const searchKeyword = ref('')
  const showFilter = ref(false)
  const filterValues = ref<FilterValues>({ categoryId: 'all' })
  
  // 分类列表
  const categories = ref<GameTabItem[]>([{ id: 'all', name: '全部' }])
  
  // 构建 API 参数
  const buildParams = () => {
    const params: Record<string, any> = {}
    if (filterValues.value.categoryId && filterValues.value.categoryId !== 'all') {
      params.categoryId = filterValues.value.categoryId
    }
    if (searchKeyword.value.trim()) {
      params.keyword = searchKeyword.value.trim()
    }
    return params
  }
  
  // 使用通用列表 Hook
  const listPage = useListPage<GameCardData>({
    fetchFn: async (params) => {
      const res = await getGames({
        page: params.page,
        page_size: params.pageSize,
        ...buildParams(),
      }, { showError: false })
      return res
    },
    extractList: (data: any) => {
      const gameList: ApiGame[] = data?.games || data?.items || data || []
      return gameList.map((g): GameCardData => ({
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
  
  const filterSections = computed<FilterSection[]>(() => [
    {
      key: 'categoryId',
      label: '分类',
      options: categories.value.map(category => ({
        label: category.name,
        value: category.id,
      })),
    },
  ])

  // 应用筛选
  const handleFilterApply = () => {
    refreshList()
  }

  // 重置筛选
  const handleFilterReset = () => {
    filterValues.value = { categoryId: 'all' }
  }
  
  // 跳转游戏详情
  const goToGame = (game: GameCardData) => {
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
    showFilter,
    filterValues,
    categories,
    filterSections,
    
    // 方法
    loadMore: listPage.loadMore,
    refresh: refreshList,
    handleSearch,
    clearSearch,
    handleFilterApply,
    handleFilterReset,
    goToGame,
    goBack,
    init,
  }
}
