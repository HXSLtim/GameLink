/**
 * 陪玩师列表专用 Hook
 * 封装陪玩师列表的数据加载、筛选、排序逻辑
 */
import { ref, computed, watch, onMounted } from 'vue'
import { useListPage } from './useListPage'
import { getPlayerList, type PlayerInfo, type PlayerListParams } from '@/api/publicPlayer'
import { getHotGames } from '@/api/game'
import { getCachedPlayers, saveCachedPlayers } from '@/utils/offlineData'
import type { GameItem } from '@/components/GameGrid/index.vue'
import type { SortOption } from '@/components/SortBar/index.vue'
import type { FilterSection, FilterValues } from '@/components/FilterPanel/index.vue'

// 陪玩师数据类型
export interface Player {
  id: number
  nickname: string
  avatar?: string
  isOnline?: boolean
  isVerified?: boolean
  status?: 'online' | 'busy' | 'offline'
  rating?: number
  orderCount?: number
  minPrice?: number
  games?: Array<{ id: number; name: string }>
}

// 默认排序选项
export const defaultSortOptions: SortOption[] = [
  { label: '推荐', value: 'recommend' },
  { label: '评分', value: 'rating' },
  { label: '价格', value: 'price' },
  { label: '销量', value: 'orders' },
]

// 默认筛选配置
export const defaultFilterSections: FilterSection[] = [
  {
    key: 'gender',
    label: '性别',
    options: [
      { label: '不限', value: '' },
      { label: '男', value: 'male' },
      { label: '女', value: 'female' },
    ],
  },
  {
    key: 'priceRange',
    label: '价格区间',
    options: [
      { label: '不限', value: '' },
      { label: '0-30元', value: '0-30' },
      { label: '30-50元', value: '30-50' },
      { label: '50元以上', value: '50+' },
    ],
  },
  {
    key: 'onlineOnly',
    label: '在线状态',
    options: [
      { label: '仅看在线', value: true },
    ],
  },
]

// 默认筛选值
export const defaultFilterValues: FilterValues = {
  gender: '',
  priceRange: '',
  onlineOnly: false,
}

// 转换 API 响应为 Player 类型
function transformPlayer(p: PlayerInfo): Player {
  return {
    id: p.id,
    nickname: p.nickname,
    avatar: p.avatar || '',
    isOnline: p.isOnline,
    isVerified: p.isVerified || false,
    status: p.isOnline ? 'online' : 'offline',
    rating: p.rating,
    orderCount: p.orderCount,
    minPrice: 20,
    games: p.tags?.slice(0, 2).map((tag, i) => ({ id: i + 1, name: tag })) || [],
  }
}

export function usePlayerList() {
  // 搜索关键词
  const searchKeyword = ref('')
  
  // 当前游戏分类
  const currentGameId = ref<number | string>('all')
  
  // 当前排序
  const currentSort = ref('recommend')
  
  // 筛选弹窗显示状态
  const showFilter = ref(false)
  
  // 筛选值
  const filterValues = ref<FilterValues>({ ...defaultFilterValues })
  
  // 游戏列表
  const games = ref<GameItem[]>([
    { id: 'all', name: '全部' },
  ])
  
  // 构建 API 参数
  const buildParams = (): Partial<PlayerListParams> => {
    const params: Partial<PlayerListParams> = {}
    
    if (currentGameId.value !== 'all') {
      params.gameId = Number(currentGameId.value)
    }
    if (filterValues.value.gender) {
      params.gender = filterValues.value.gender as 'male' | 'female'
    }
    if (filterValues.value.priceRange) {
      const priceRange = filterValues.value.priceRange as string
      const [min, max] = priceRange.split('-')
      if (min) params.minPrice = Number(min) * 100
      if (max && max !== '+') params.maxPrice = Number(max) * 100
    }
    if (filterValues.value.onlineOnly) {
      params.isOnline = true
    }
    if (searchKeyword.value.trim()) {
      params.keyword = searchKeyword.value.trim()
    }
    if (currentSort.value !== 'recommend') {
      params.sortBy = currentSort.value as PlayerListParams['sortBy']
    }
    
    return params
  }
  
  // 使用通用列表 Hook
  const listPage = useListPage<Player, PlayerListParams>({
    fetchFn: async (params) => {
      const res = await getPlayerList(params, { showError: false })
      return res
    },
    extractList: (data: any) => {
      const playerList: PlayerInfo[] = data?.players || data || []
      return playerList.map(transformPlayer)
    },
    pageSize: 10,
    getCacheFn: () => {
      const cached = getCachedPlayers()
      if (!cached) return null
      return cached.map(p => ({
        id: p.id,
        nickname: p.nickname,
        avatar: p.avatar || '',
        isOnline: p.isOnline,
        isVerified: false,
        status: (p.isOnline ? 'online' : 'offline') as 'online' | 'offline',
        rating: p.rating,
        orderCount: p.orderCount,
        minPrice: p.hourlyRate / 100,
        games: p.mainGame ? [{ id: 1, name: p.mainGame }] : [],
      }))
    },
    saveCacheFn: (items) => {
      saveCachedPlayers(items.map(p => ({
        id: p.id,
        nickname: p.nickname,
        avatar: p.avatar || '',
        rank: p.games?.[0]?.name || '',
        rating: p.rating || 5.0,
        hourlyRate: (p.minPrice || 20) * 100,
        isOnline: p.isOnline || false,
        orderCount: p.orderCount || 0,
        mainGame: p.games?.[0]?.name || '',
      })))
    },
  })
  
  // 加载游戏分类
  const loadGames = async () => {
    try {
      const res = await getHotGames(20, { showError: false })
      const resData = res.data as any
      const apiGames = resData?.items || resData?.games || resData || []
      games.value = [
        { id: 'all', name: '全部' },
        ...apiGames.map((g: any) => ({
          id: g.id,
          name: g.name,
          icon: g.iconUrl || g.icon,
        })),
      ]
    } catch (error) {
      console.error('加载游戏分类失败', error)
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
  
  // 选择游戏
  const handleGameSelect = () => {
    refreshList()
  }
  
  // 切换排序
  const handleSortChange = () => {
    refreshList()
  }
  
  // 应用筛选
  const handleFilterApply = () => {
    refreshList()
  }
  
  // 重置筛选
  const handleFilterReset = () => {
    filterValues.value = { ...defaultFilterValues }
  }
  
  // 跳转详情
  const goToDetail = (playerId: number) => {
    uni.navigateTo({ url: `/pages/player/detail/index?id=${playerId}` })
  }
  
  // 初始化
  const init = () => {
    loadGames()
    refreshList()
  }
  
  return {
    // 列表数据
    players: listPage.list,
    pageState: listPage.pageState,
    errorMessage: listPage.errorMessage,
    loadingMore: listPage.loadingMore,
    noMore: listPage.noMore,
    isOffline: listPage.isOffline,
    
    // 筛选相关
    searchKeyword,
    currentGameId,
    currentSort,
    showFilter,
    filterValues,
    games,
    
    // 配置
    sortOptions: defaultSortOptions,
    filterSections: defaultFilterSections,
    
    // 方法
    loadMore: listPage.loadMore,
    refresh: refreshList,
    handleSearch,
    handleGameSelect,
    handleSortChange,
    handleFilterApply,
    handleFilterReset,
    goToDetail,
    init,
  }
}
