/**
 * 收藏列表专用 Hook
 */
import { ref, computed } from 'vue'
import { useListPage } from './useListPage'
import { getFavorites, batchRemoveFavorites, type FavoritePlayer as ApiFavoritePlayer } from '@/api/favorite'
import type { FavoritePlayerData } from '@/types/player'
import { confirmDialog } from '@/composables/useConfirmDialog'

export function useFavoriteList() {
  // 编辑模式
  const isEditMode = ref(false)
  const selectedIds = ref<number[]>([])
  
  // 使用通用列表 Hook
  const listPage = useListPage<FavoritePlayerData>({
    fetchFn: async (params) => {
      const res = await getFavorites({
        page: params.page,
        page_size: params.pageSize,
      }, { showError: false })
      return res
    },
    extractList: (data: any) => {
      const list: ApiFavoritePlayer[] = data?.items || data || []
      return list.map((item): FavoritePlayerData => ({
        id: item.id,
        playerId: item.playerId,
        nickname: item.nickname || item.playerName || '陪玩师',
        avatar: item.avatar || item.playerAvatar,
        isOnline: item.isOnline,
        rating: item.rating || item.playerRating || 5.0,
        orderCount: item.orderCount || 0,
        minPrice: item.minPrice || 20,
        games: item.games || item.gameNames || [],
      }))
    },
    pageSize: 20,
  })
  
  // 全选状态
  const isAllSelected = computed(() => {
    return listPage.list.value.length > 0 && selectedIds.value.length === listPage.list.value.length
  })
  
  // 切换编辑模式
  const toggleEditMode = () => {
    isEditMode.value = !isEditMode.value
    if (!isEditMode.value) {
      selectedIds.value = []
    }
  }
  
  // 切换选择
  const toggleSelect = (id: number) => {
    const index = selectedIds.value.indexOf(id)
    if (index > -1) {
      selectedIds.value.splice(index, 1)
    } else {
      selectedIds.value.push(id)
    }
  }
  
  // 全选/取消全选
  const toggleSelectAll = () => {
    if (isAllSelected.value) {
      selectedIds.value = []
    } else {
      selectedIds.value = listPage.list.value.map(item => item.id)
    }
  }
  
  // 删除选中
  const deleteSelected = async () => {
    if (selectedIds.value.length === 0) return

    const confirmed = await confirmDialog({
      title: '确认取消收藏',
      content: `确定要取消收藏选中的 ${selectedIds.value.length} 位陪玩师吗？`,
    })
    if (!confirmed) return

    try {
      uni.showLoading({ title: '处理中...' })
      const selectedPlayerIds = listPage.list.value
        .filter(item => selectedIds.value.includes(item.id))
        .map(item => item.playerId)
      await batchRemoveFavorites(selectedPlayerIds)
      uni.hideLoading()
      uni.showToast({ title: '已取消收藏', icon: 'success' })
      
      // 移除已删除的项
      listPage.list.value = listPage.list.value.filter(
        item => !selectedIds.value.includes(item.id)
      )
      selectedIds.value = []
      
      // 如果列表为空，退出编辑模式
      if (listPage.list.value.length === 0) {
        isEditMode.value = false
        listPage.pageState.value = 'empty'
      }
    } catch (error: any) {
      uni.hideLoading()
      uni.showToast({ title: error?.message || '操作失败', icon: 'none' })
    }
  }
  
  // 导航
  const goBack = () => uni.navigateBack()
  
  const goToDetail = (player: FavoritePlayerData) => {
    if (isEditMode.value) {
      toggleSelect(player.id)
      return
    }
    uni.navigateTo({ url: `/pages/player/detail/index?id=${player.playerId}` })
  }
  
  const goToPlayerList = () => {
    uni.switchTab({ url: '/pages/player/list/index' })
  }
  
  return {
    // 列表数据
    favorites: listPage.list,
    pageState: listPage.pageState,
    loading: listPage.loading,
    loadingMore: listPage.loadingMore,
    noMore: listPage.noMore,
    
    // 编辑模式
    isEditMode,
    selectedIds,
    isAllSelected,
    
    // 方法
    loadMore: listPage.loadMore,
    refresh: listPage.refresh,
    toggleEditMode,
    toggleSelect,
    toggleSelectAll,
    deleteSelected,
    goBack,
    goToDetail,
    goToPlayerList,
  }
}
