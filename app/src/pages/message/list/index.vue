<template>
  <BasePageLayout
    class="message-page"
    :scroll="false"
    padding="0"
    title="消息"
    :show-back="false"
    :show-tab-bar="true"
    :tab-bar-current="2"
  >
    <!-- 移动端 Tab 栏（PC 端内置在分栏左侧） -->
    <template v-if="!isPC" #tabs>
      <TabsBar
        v-model="currentTab"
        :tabs="tabs"
        stretch
        @change="switchTab"
      />
    </template>

    <!-- PC 端：左右分栏布局 -->
    <view v-if="isPC" class="split-layout">
      <!-- 左侧：会话列表 -->
      <view class="split-left">
        <TabsBar
          v-model="currentTab"
          :tabs="tabs"
          stretch
          @change="switchTab"
        />
        <scroll-view class="conv-scroll" scroll-y>
          <view v-if="pageState === 'loading'" class="conv-loading">
            <Skeleton :rows="4" />
          </view>
          <view v-else-if="pageState === 'empty'" class="conv-empty">
            <text class="conv-empty-text">暂无消息</text>
          </view>
          <view v-else class="conv-list">
            <view
              v-for="item in messages"
              :key="item.id"
              class="conv-item"
              :class="{ active: selectedChatId === item.conversationId }"
              @click="handleItemClick(item)"
            >
              <MessageItem :message="item" />
            </view>
          </view>
        </scroll-view>
      </view>

      <!-- 右侧：聊天窗口 -->
      <view class="split-right">
        <!-- 未选中会话 -->
        <view v-if="!selectedChatId" class="chat-empty">
          <uv-icon name="chat" size="48" color="var(--color-text-placeholder)" />
          <text class="chat-empty-title">选择一个会话开始聊天</text>
          <text class="chat-empty-desc">从左侧列表选择聊天对象</text>
        </view>

        <!-- 已选中：嵌入聊天面板 -->
        <view v-else class="chat-panel" :key="selectedChatId">
          <ChatNavBar
            :name="chatInfo.name"
            :type="chatInfo.type"
            :is-online="chatInfo.isOnline"
            :member-count="chatInfo.memberCount"
            @back="selectedChatId = null"
            @menu="showMenu = true"
          />

          <ChatConnectionStatus
            :status="wsStatus"
            @reconnect="reconnectWs"
          />

          <ChatMessageList
            :loading="chatLoading"
            :loading-history="loadingHistory"
            :has-more-history="hasMoreHistory"
            :messages="chatMessages"
            :scroll-to-id="scrollToId"
            :current-user-id="currentUserId"
            :current-user-name="currentUserName"
            :current-user-avatar="currentUserAvatar"
            :show-time="shouldShowTime"
            :show-sender-name="chatInfo.type !== 'private'"
            :playing-voice-id="playingVoiceId"
            @load-more="loadMoreHistory"
            @avatar-click="viewProfile"
            @image-preview="previewImage"
            @voice-play="playVoice"
            @order-click="viewOrder"
            @resend="resendMessage"
          />

          <ChatInputBar
            v-model="inputText"
            :recording="recording"
            @send="sendTextMessage"
            @image="chooseImage"
            @more="showMore = true"
          />

          <!-- 更多面板 -->
          <ChatMorePanel
            :show="showMore"
            :show-order="!!chatInfo.orderId"
            @close="showMore = false"
            @image="chooseImage"
            @camera="takePhoto"
            @order="viewOrder()"
            @report="reportChat"
          />

          <!-- 遮罩 -->
          <view v-if="showMore || showMenu" class="overlay" @tap="showMore = false; showMenu = false" />

          <!-- 菜单 -->
          <ChatMenuPanel
            v-model:show="showMenu"
            :target-id="chatInfo.targetId"
            @view-profile="viewProfile"
            @clear="clearHistory"
            @block="blockUser"
          />
        </view>
      </view>
    </view>

    <!-- 移动端：纯列表，点击跳转 -->
    <InfiniteList
      v-if="!isPC"
      :state="pageState"
      :error-message="errorMessage"
      :show-load-more="false"
      empty-title="暂无消息"
      empty-desc="去找个陪玩师聊聊吧"
      padding="24rpx"
      @retry="loadMessages"
    >
      <view
        v-for="item in messages"
        :key="item.id"
        class="mobile-item"
        @click="goToChat(item)"
      >
        <MessageItem :message="item" />
      </view>
    </InfiniteList>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import TabsBar from '@/components/TabsBar/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import InfiniteList from '@/components/InfiniteList/index.vue'
import Skeleton from '@/components/Skeleton/index.vue'
// Business 组件
import MessageItem from '@/components/MessageItem/index.vue'
import ChatNavBar from '@/components/ChatNavBar/index.vue'
import ChatConnectionStatus from '@/components/ChatConnectionStatus/index.vue'
import ChatMessageList from '@/components/ChatMessageList/index.vue'
import ChatInputBar from '@/components/ChatInputBar/index.vue'
import ChatMorePanel from '@/components/ChatMorePanel/index.vue'
import ChatMenuPanel from '@/components/ChatMenuPanel/index.vue'
// Composables
import { useMessageList } from '@/composables/useMessageList'
import { useChatRoom } from '@/composables/useChatRoom'
import { useDevice } from '@/composables/useDevice'
import type { MessageData } from '@/types/message'

const { isPC } = useDevice()

// ── 消息列表 ──
const {
  messages,
  pageState,
  errorMessage,
  tabs,
  currentTab,
  selectedChatId,
  switchTab,
  loadMessages,
  goToChat,
  selectChat,
} = useMessageList()

// ── 聊天室（PC 端嵌入用） ──
const {
  loading: chatLoading,
  loadingHistory,
  hasMoreHistory,
  showMore,
  showMenu,
  recording,
  playingVoiceId,
  wsStatus,
  chatInfo,
  messages: chatMessages,
  inputText,
  scrollToId,
  currentUserId,
  currentUserName,
  currentUserAvatar,
  initChat,
  loadMoreHistory,
  sendTextMessage,
  chooseImage,
  takePhoto,
  resendMessage,
  previewImage,
  playVoice,
  shouldShowTime,
  reconnectWs,
  disconnectWs,
  viewProfile,
  viewOrder,
  clearHistory,
  blockUser,
  reportChat,
} = useChatRoom()

// PC 端：切换选中会话时初始化聊天
watch(selectedChatId, (id) => {
  if (id) {
    disconnectWs()
    initChat({ groupId: String(id) })
  }
})

// 点击会话：PC 选中，移动端已由 goToChat 处理
const handleItemClick = (item: MessageData) => {
  selectChat(item)
}

onMounted(() => {
  loadMessages()
})

onShow(() => {
  loadMessages()
})
</script>

<style lang="scss" scoped>
.message-page {
  padding-bottom: calc(110rpx + env(safe-area-inset-bottom));

  @include desktop {
    padding-bottom: 0;
  }
}

// ============================================
// PC 端分栏布局
// ============================================
.split-layout {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.split-left {
  width: 320px;
  min-width: 280px;
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--color-border);
  background: var(--color-bg-card);
}

.conv-scroll {
  flex: 1;
  overflow-y: auto;
}

.conv-loading {
  padding: var(--spacing-md);
}

.conv-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-xl);
}

.conv-empty-text {
  font-size: var(--font-sm);
  color: var(--color-text-placeholder);
}

.conv-list {
  padding: var(--spacing-xs);
}

.conv-item {
  border-radius: var(--radius-md);
  transition: background 0.15s ease;
  cursor: pointer;

  // PC 列表中去掉 MessageItem 自带的卡片边框，统一由 conv-item 控制
  :deep(.message-item) {
    border: none;
    background: transparent;
    border-radius: 0;
  }

  &:hover {
    background: var(--color-bg-secondary);
  }

  &.active {
    background: rgba(var(--color-primary-rgb), 0.08);
  }
}

.split-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  background: var(--color-bg);
  position: relative;
}

// ── 空状态 ──
.chat-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
}

.chat-empty-title {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-top: var(--spacing-sm);
}

.chat-empty-desc {
  font-size: var(--font-sm);
  color: var(--color-text-placeholder);
}

// ── 嵌入聊天面板 ──
.chat-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 150;
}

// ============================================
// 移动端列表项
// ============================================
.mobile-item {
  margin-bottom: var(--spacing-xs);
  padding: 0 var(--spacing-sm);
}
</style>
