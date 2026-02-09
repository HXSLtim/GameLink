<template>
  <BasePageLayout
    class="favorite-page"
    :scroll="false"
    padding="0"
    title="我的收藏"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #nav>
      <!-- 顶部导航 -->
      <NavBar title="我的收藏" @back="goBack">
        <template #right>
          <text v-if="isEditMode" class="nav-action" @tap="toggleEditMode">完成</text>
          <text v-else-if="favorites.length > 0" class="nav-action" @tap="toggleEditMode">管理</text>
        </template>
      </NavBar>
    </template>

    <!-- 收藏列表 -->
    <FavoriteListPanel
      :favorites="favorites"
      :page-state="pageState"
      :loading-more="loadingMore"
      :no-more="noMore"
      :is-edit-mode="isEditMode"
      :selected-ids="selectedIds"
      @load-more="loadMore"
      @retry="refresh"
      @empty-action="goToPlayerList"
      @item-click="goToDetail"
      @toggle-select="toggleSelect"
    />

    <template #footer>
      <!-- 底部操作栏（编辑模式） -->
      <FavoriteEditBar
        v-if="isEditMode"
        :all-selected="isAllSelected"
        :selected-count="selectedIds.length"
        @toggle-all="toggleSelectAll"
        @delete="deleteSelected"
      />
    </template>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
// Business 组件
import FavoriteListPanel from '@/components/FavoriteListPanel/index.vue'
import FavoriteEditBar from '@/components/FavoriteEditBar/index.vue'
// Composables
import { useFavoriteList } from '@/composables/useFavoriteList'

const {
  favorites,
  pageState,
  loadingMore,
  noMore,
  isEditMode,
  selectedIds,
  isAllSelected,
  loadMore,
  refresh,
  toggleEditMode,
  toggleSelect,
  toggleSelectAll,
  deleteSelected,
  goBack,
  goToDetail,
  goToPlayerList,
} = useFavoriteList()

onShow(() => {
  refresh()
})
</script>

<style lang="scss" scoped>
.nav-action {
  font-size: var(--font-sm);
  color: var(--color-primary);
  padding: 8rpx 16rpx;
}

</style>
