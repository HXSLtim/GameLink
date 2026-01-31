<template>
  <view class="chat-page page-container">
    <!-- 顶部导航 -->
    <ChatNavBar
      :name="chatInfo.name"
      :type="chatInfo.type"
      :is-online="chatInfo.isOnline"
      :member-count="chatInfo.memberCount"
      @back="goBack"
      @menu="showMenu = true"
    />

    <!-- WebSocket 连接状态 -->
    <ChatConnectionStatus
      :status="wsStatus"
      @reconnect="reconnectWs"
    />

    <!-- 加载状态 -->
    <view v-if="loading" class="loading-wrap">
      <uv-skeleton rows="5" title loading></uv-skeleton>
    </view>

    <!-- 消息列表 -->
    <scroll-view 
      v-else
      class="message-scroll"
      scroll-y
      :scroll-into-view="scrollToId"
      :scroll-with-animation="true"
      @scrolltoupper="loadMoreHistory"
    >
      <!-- 加载更多历史 -->
      <view v-if="hasMoreHistory" class="load-more-history" @tap="loadMoreHistory">
        <text v-if="loadingHistory">加载中...</text>
        <text v-else>加载更多历史消息</text>
      </view>
      
      <!-- 消息列表 -->
      <ChatMessageBubble
        v-for="(message, index) in messages"
        :key="message.id"
        :id="`msg-${message.id}`"
        :message="message"
        :current-user-id="currentUserId"
        :self-name="currentUserName"
        :self-avatar="currentUserAvatar"
        :show-time="shouldShowTime(message, index)"
        :show-sender-name="chatInfo.type !== 'private'"
        :is-playing="playingVoiceId === message.id"
        @avatar-click="viewProfile"
        @image-preview="previewImage"
        @voice-play="playVoice"
        @order-click="viewOrder"
        @resend="resendMessage"
      />
      
      <!-- 底部占位 -->
      <view id="msg-bottom" class="scroll-bottom"></view>
    </scroll-view>

    <!-- 输入区域 -->
    <ChatInputBar
      v-model="inputText"
      :recording="recording"
      @send="sendTextMessage"
      @image="chooseImage"
      @more="showMore = true"
    />

    <!-- 更多功能面板 -->
    <ChatMorePanel
      :show="showMore"
      :show-order="!!chatInfo.orderId"
      @close="showMore = false"
      @image="chooseImage"
      @camera="takePhoto"
      @order="viewOrder()"
      @report="reportChat"
    />

    <!-- 遮罩层 -->
    <view v-if="showMore || showMenu" class="overlay" @tap="showMore = false; showMenu = false"></view>

    <!-- 菜单弹窗 -->
    <uv-popup :show="showMenu" mode="bottom" round="24" @close="showMenu = false">
      <view class="menu-panel">
        <view class="menu-item" @tap="viewProfile(chatInfo.targetId!)">
          <text>查看资料</text>
        </view>
        <view class="menu-item" @tap="clearHistory">
          <text>清空聊天记录</text>
        </view>
        <view class="menu-item danger" @tap="blockUser">
          <text>拉黑用户</text>
        </view>
        <view class="menu-item" @tap="showMenu = false">
          <text>取消</text>
        </view>
      </view>
    </uv-popup>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onUnmounted } from 'vue'
import { onLoad, onHide } from '@dcloudio/uni-app'
// Business 组件
import ChatNavBar from '@/components/ChatNavBar/index.vue'
import ChatConnectionStatus from '@/components/ChatConnectionStatus/index.vue'
import ChatMessageBubble from '@/components/ChatMessageBubble/index.vue'
import ChatInputBar from '@/components/ChatInputBar/index.vue'
import ChatMorePanel from '@/components/ChatMorePanel/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useChatRoom } from '@/composables/useChatRoom'

const {
  loading,
  loadingHistory,
  hasMoreHistory,
  showMore,
  showMenu,
  recording,
  playingVoiceId,
  wsStatus,
  chatInfo,
  messages,
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
  goBack,
  viewProfile,
  viewOrder,
  clearHistory,
  blockUser,
  reportChat,
} = useChatRoom()

onLoad((options) => {
  if (options) {
    initChat(options)
  }
})

onHide(() => {
  disconnectWs()
})

onUnmounted(() => {
  disconnectWs()
})
</script>

<style lang="scss" scoped>
.chat-page {
  height: 100vh;
  height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  overflow: hidden;
}

.loading-wrap {
  flex: 1;
  padding: 24rpx;
}

.message-scroll {
  flex: 1;
  overflow-y: auto;
}

.load-more-history {
  text-align: center;
  padding: 24rpx;
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.scroll-bottom {
  height: 24rpx;
}

.overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 150;
}

.menu-panel {
  padding: 16rpx 0;
  padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
}

.menu-item {
  padding: 32rpx 48rpx;
  text-align: center;
  font-size: 32rpx;
  color: var(--color-text);
  
  &.danger {
    color: var(--color-error);
  }
  
  &:active {
    background: var(--color-bg-hover);
  }
}
</style>
