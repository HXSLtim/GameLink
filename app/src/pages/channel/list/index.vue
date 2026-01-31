<template>
  <view class="channel-list-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="公共频道" @back="goBack" />

    <!-- 搜索栏 -->
    <SearchBar
      v-model="searchKeyword"
      placeholder="搜索频道"
      @search="handleSearch"
    />

    <!-- 离线提示 -->
    <OfflineBanner
      :visible="isOffline"
      message="网络不可用，显示推荐频道"
      @action="refresh"
    />

    <!-- 游戏分类 -->
    <GameTabs
      v-model="currentGameId"
      :games="games"
      @select="handleGameSelect"
    />

    <!-- 频道列表 -->
    <InfiniteList
      :state="pageState"
      :loading="loadingMore"
      :no-more="noMore"
      :error-message="errorMessage"
      empty-title="暂无频道"
      empty-desc="换个分类试试吧"
      padding="24rpx"
      @load-more="loadMore"
      @retry="refresh"
    >
      <ListItem
        v-for="(channel, index) in channels"
        :key="channel.id"
        :index="index"
        @click="enterChannel(channel)"
      >
        <ChannelCard
          :channel="channel"
          @join="joinChannel(channel)"
          @leave="leaveChannel(channel)"
        />
      </ListItem>
    </InfiniteList>
    
    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import SearchBar from '@/components/SearchBar/index.vue'
import OfflineBanner from '@/components/OfflineBanner/index.vue'
import GameTabs from '@/components/GameTabs/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import ListItem from '@/components/ListItem/index.vue'
// Business 组件
import ChannelCard from '@/components/ChannelCard/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useChannelList } from '@/composables/useChannelList'

const {
  channels,
  pageState,
  errorMessage,
  loadingMore,
  noMore,
  isOffline,
  searchKeyword,
  currentGameId,
  games,
  loadMore,
  refresh,
  handleSearch,
  handleGameSelect,
  joinChannel,
  leaveChannel,
  enterChannel,
  goBack,
  init,
} = useChannelList()

onMounted(init)

onShow(() => {
  refresh()
})
</script>

<style lang="scss" scoped>
.channel-list-page {
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
</style>
