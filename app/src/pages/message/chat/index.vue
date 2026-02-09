<template>
  <PageShell class="chat-page" :scroll="false" padding="0" content-class="chat-content">
    <template #header>
      <!-- 顶部导航 -->
      <ChatNavBar
        :name="chatInfo.name"
        :type="chatInfo.type"
        :is-online="chatInfo.isOnline"
        :member-count="chatInfo.memberCount"
        @back="goBack"
        @menu="showMenu = true"
      />
    </template>

    <!-- WebSocket 连接状态 -->
    <ChatConnectionStatus
      :status="wsStatus"
      @reconnect="reconnectWs"
    />

    <!-- 消息列表 -->
    <ChatMessageList
      :loading="loading"
      :loading-history="loadingHistory"
      :has-more-history="hasMoreHistory"
      :messages="messages"
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
      @order="viewOrder"
      @report="reportChat"
    />

    <!-- 遮罩层 -->
    <view v-if="showMore || showMenu" class="overlay" @tap="showMore = false; showMenu = false"></view>

    <!-- 菜单弹窗 -->
    <ChatMenuPanel
      v-model:show="showMenu"
      :target-id="chatInfo.targetId"
      @view-profile="viewProfile"
      @clear="clearHistory"
      @block="blockUser"
    />

    <template #footer>
      <!-- PC 端侧边栏 -->
      <CustomTabBar :show-mobile-tab-bar="false" />
    </template>
  </PageShell>
</template>

<script setup lang="ts">
import { onUnmounted } from 'vue'
import { onLoad, onHide } from '@dcloudio/uni-app'
// Business 组件
import ChatNavBar from '@/components/ChatNavBar/index.vue'
import PageShell from '@/components/layout/PageShell/index.vue'
import ChatConnectionStatus from '@/components/ChatConnectionStatus/index.vue'
import ChatMessageList from '@/components/ChatMessageList/index.vue'
import ChatInputBar from '@/components/ChatInputBar/index.vue'
import ChatMorePanel from '@/components/ChatMorePanel/index.vue'
import ChatMenuPanel from '@/components/ChatMenuPanel/index.vue'
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
.chat-content {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
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

</style>
