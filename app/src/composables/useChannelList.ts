/**
 * 频道列表专用 Hook
 */
import { ref, onMounted, computed } from 'vue'
import { useListPage } from './useListPage'
import { getPublicChannels, joinChatGroup, leaveChatGroup, type PublicChannel } from '@/api/chat'
import { getHotGames, type Game } from '@/api/game'
import { useUserStore } from '@/store/user'
import { setRedirectPath } from '@/utils/routeGuard'
import { getCachedChannels, saveCachedChannels } from '@/utils/offlineData'
import { confirmDialog } from '@/composables/useConfirmDialog'
import type { ChannelData } from '@/types/community'
import type { GameTabItem } from '@/types/game'
import type { FilterSection, FilterValues } from '@/types/filter'

export function useChannelList() {
  const userStore = useUserStore()
  
  // 搜索和筛选
  const searchKeyword = ref('')
  const showFilter = ref(false)
  const filterValues = ref<FilterValues>({ gameId: 'all' })
  const isOfflineMode = ref(false)
  
  // 游戏分类
  const games = ref<GameTabItem[]>([{ id: 'all', name: '全部' }])
  
  // 构建 API 参数
  const buildParams = () => {
    const params: Record<string, any> = {}
    if (filterValues.value.gameId && filterValues.value.gameId !== 'all') {
      params.gameId = Number(filterValues.value.gameId)
    }
    if (searchKeyword.value.trim()) {
      params.keyword = searchKeyword.value.trim()
    }
    return params
  }
  
  // 使用通用列表 Hook
  const listPage = useListPage<ChannelData>({
    fetchFn: async (params) => {
      const res = await getPublicChannels({
        page: params.page,
        page_size: params.pageSize,
        ...buildParams(),
      }, { showError: false })
      return res
    },
    extractList: (data: any) => {
      const channelList: PublicChannel[] = data?.channels || data || []
      return channelList.map((c): ChannelData => ({
        id: c.id,
        name: c.name,
        description: c.description,
        avatar: c.avatar,
        memberCount: c.memberCount,
        maxMembers: c.maxMembers,
        isActive: c.memberCount > 5,
        isJoined: c.isJoined || false,
        gameId: c.gameId,
        gameName: c.gameName,
      }))
    },
    pageSize: 20,
    getCacheFn: () => {
      const cached = getCachedChannels()
      if (!cached) return null
      return cached.map(c => ({
        id: c.id,
        name: c.name,
        description: c.description,
        avatar: c.avatar || c.avatarUrl,
        memberCount: c.memberCount ?? c.currentMembers ?? 0,
        maxMembers: 100,
        isActive: (c.memberCount ?? c.currentMembers ?? 0) > 5,
        isJoined: false,
        gameId: c.gameId,
        gameName: c.gameName,
      }))
    },
    saveCacheFn: (items) => {
      saveCachedChannels(items.map(c => ({
        id: c.id,
        name: c.name,
        description: c.description || '',
        avatar: c.avatar || '',
        memberCount: c.memberCount,
        gameId: c.gameId,
        gameName: c.gameName,
      })))
    },
  })
  
  // 加载游戏分类
  const loadGames = async () => {
    try {
      const res = await getHotGames(10, { showError: false })
      const resData = res.data as any
      const apiGames = resData?.games || resData || []
      games.value = [
        { id: 'all', name: '全部' },
        ...apiGames.map((g: Game) => ({
          id: g.id,
          name: g.name,
          icon: g.icon,
        })),
      ]
    } catch (error) {
      console.error('加载游戏分类失败', error)
    }
  }
  
  // 刷新列表
  const refreshList = () => {
    isOfflineMode.value = false
    listPage.refresh(buildParams())
  }
  
  // 搜索
  const handleSearch = () => {
    refreshList()
  }
  
  const filterSections = computed<FilterSection[]>(() => [
    {
      key: 'gameId',
      label: '分类',
      options: games.value.map(game => ({
        label: game.name,
        value: game.id,
      })),
    },
  ])

  const handleFilterApply = () => {
    refreshList()
  }

  const handleFilterReset = () => {
    filterValues.value = { gameId: 'all' }
  }
  
  // 加入频道
  const joinChannel = async (channel: ChannelData) => {
    if (!userStore.isLoggedIn) {
      setRedirectPath('/pages/channel/list/index')
      uni.navigateTo({ url: '/pages/auth/login/index' })
      return
    }
    
    try {
      uni.showLoading({ title: '加入中...' })
      await joinChatGroup(channel.id)
      channel.isJoined = true
      channel.memberCount++
      uni.hideLoading()
      uni.showToast({ title: '加入成功', icon: 'success' })
    } catch (error: any) {
      uni.hideLoading()
      uni.showToast({ title: error?.message || '加入失败', icon: 'none' })
    }
  }
  
  // 离开频道
  const leaveChannel = async (channel: ChannelData) => {
    const confirmed = await confirmDialog({
      title: '确认离开',
      content: `确定要离开"${channel.name}"频道吗？`,
    })
    if (!confirmed) return
    try {
      uni.showLoading({ title: '处理中...' })
      await leaveChatGroup(channel.id)
      channel.isJoined = false
      if (channel.memberCount > 0) channel.memberCount--
      uni.hideLoading()
      uni.showToast({ title: '已离开', icon: 'success' })
    } catch (error: any) {
      uni.hideLoading()
      uni.showToast({ title: error?.message || '操作失败', icon: 'none' })
    }
  }
  
  // 进入频道
  const enterChannel = (channel: ChannelData) => {
    if (!channel.isJoined) {
      joinChannel(channel)
      return
    }
    uni.navigateTo({ url: `/pages/message/chat/index?groupId=${channel.id}` })
  }
  
  // 返回
  const goBack = () => {
    uni.navigateBack()
  }
  
  // 初始化
  const init = () => {
    loadGames()
    refreshList()
  }
  
  return {
    // 列表数据
    channels: listPage.list,
    pageState: listPage.pageState,
    errorMessage: listPage.errorMessage,
    loadingMore: listPage.loadingMore,
    noMore: listPage.noMore,
    isOffline: listPage.isOffline,
    
    // 筛选
    searchKeyword,
    showFilter,
    filterValues,
    games,
    filterSections,
    isOfflineMode,
    
    // 方法
    loadMore: listPage.loadMore,
    refresh: refreshList,
    handleSearch,
    handleFilterApply,
    handleFilterReset,
    joinChannel,
    leaveChannel,
    enterChannel,
    goBack,
    init,
  }
}
