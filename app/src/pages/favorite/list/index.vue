<template>
  <view class="favorite-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="我的收藏" @back="goBack">
      <template #right>
        <text v-if="isEditMode" class="nav-action" @tap="toggleEditMode">完成</text>
        <text v-else-if="favorites.length > 0" class="nav-action" @tap="toggleEditMode">管理</text>
      </template>
    </NavBar>

    <!-- 收藏列表 -->
    <InfiniteList
      :state="pageState"
      :loading="loadingMore"
      :no-more="noMore"
      empty-title="暂无收藏"
      empty-desc="去发现喜欢的陪玩师吧"
      padding="24rpx"
      @load-more="loadMore"
      @retry="refresh"
    >
      <template #empty-action>
        <GlButton type="primary" size="small" @click="goToPlayerList">去看看</GlButton>
      </template>
      
      <ListItem
        v-for="(item, index) in favorites"
        :key="item.id"
        :index="index"
        @click="goToDetail(item)"
      >
        <FavoritePlayerCard
          :player="item"
          :edit-mode="isEditMode"
          :selected="selectedIds.includes(item.id)"
          @toggle-select="toggleSelect(item.id)"
        />
      </ListItem>
    </InfiniteList>

    <!-- 底部操作栏（编辑模式） -->
    <view v-if="isEditMode" class="action-bar">
      <view class="select-all" @tap="toggleSelectAll">
        <view class="checkbox" :class="{ checked: isAllSelected }">
          <uv-icon v-if="isAllSelected" name="checkbox-mark" size="16" color="#fff"></uv-icon>
        </view>
        <text>全选</text>
      </view>
      <view class="action-info">
        <text>已选 {{ selectedIds.length }} 项</text>
      </view>
      <GlButton 
        type="error" 
        size="small" 
        :disabled="selectedIds.length === 0"
        @click="deleteSelected"
      >
        取消收藏
      </GlButton>
    </view>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import ListItem from '@/components/ListItem/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import FavoritePlayerCard from '@/components/FavoritePlayerCard/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
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
.favorite-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  
  @include desktop {
    height: 100vh;
    min-height: auto;
    overflow: hidden;
  }
}

.nav-action {
  font-size: 28rpx;
  color: var(--color-primary);
  padding: 8rpx 16rpx;
}

.action-bar {
  display: flex;
  align-items: center;
  gap: 24rpx;
  padding: 20rpx 32rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: var(--color-bg-card);
  border-top: 1rpx solid var(--color-border);
}

.select-all {
  display: flex;
  align-items: center;
  gap: 12rpx;
  
  text {
    font-size: 28rpx;
    color: var(--color-text);
  }
}

.checkbox {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 2rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  
  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
  }
}

.action-info {
  flex: 1;
  text-align: center;
  font-size: 26rpx;
  color: var(--color-text-secondary);
}
</style>
