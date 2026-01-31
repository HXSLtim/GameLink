<template>
  <view class="message-page page-container">
    <!-- 顶部标签 -->
    <TabsBar
      v-model="currentTab"
      :tabs="tabs"
      stretch
      @change="switchTab"
    />
    
    <!-- 消息列表 -->
    <InfiniteList
      :state="pageState"
      :error-message="errorMessage"
      :show-load-more="false"
      empty-title="暂无消息"
      empty-desc="去找个陪玩师聊聊吧"
      padding="0"
      @retry="loadMessages"
    >
      <ListItem
        v-for="(item, index) in messages"
        :key="item.id"
        :index="index"
        :animated="true"
        @click="goToChat(item)"
      >
        <MessageItem :message="item" />
      </ListItem>
    </InfiniteList>
    
    <!-- 自定义 TabBar -->
    <CustomTabBar :current="2" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import TabsBar from '@/components/TabsBar/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import ListItem from '@/components/ListItem/index.vue'
// Business 组件
import MessageItem from '@/components/MessageItem/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useMessageList } from '@/composables/useMessageList'

// 使用消息列表 Hook
const {
  messages,
  pageState,
  errorMessage,
  tabs,
  currentTab,
  switchTab,
  loadMessages,
  goToChat,
} = useMessageList()

onMounted(() => {
  loadMessages()
})

onShow(() => {
  loadMessages()
})
</script>

<style lang="scss" scoped>
.message-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  padding-bottom: calc(110rpx + env(safe-area-inset-bottom));
  
  @include desktop {
    height: 100vh;
    min-height: auto;
    padding-bottom: 0;
    overflow: hidden;
  }
}
</style>
