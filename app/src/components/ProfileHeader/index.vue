<template>
  <view class="profile-header">
    <view class="user-info">
      <GlAvatar
        :src="avatar"
        :text="displayName"
        :size="120"
        shape="circle"
        bordered
      />
      <view class="user-details">
        <text class="user-name">{{ displayName }}</text>
        <view v-if="userId" class="user-id">
          <text>ID: {{ userId }}</text>
        </view>
        <view v-if="!isLoggedIn" class="login-tip">
          <text>点击登录，享受更多服务</text>
        </view>
        <view v-else-if="isPlayer" class="player-badge">
          <text>陪玩师</text>
        </view>
      </view>
      <GlButton
        v-if="isLoggedIn"
        size="small"
        type="default"
        round
        plain
        custom-style="background: rgba(255,255,255,0.2); border-color: rgba(255,255,255,0.3); color: #fff;"
        @click="$emit('edit')"
      >
        编辑资料
      </GlButton>
      <GlButton
        v-else
        size="small"
        type="default"
        round
        plain
        custom-style="background: rgba(255,255,255,0.2); border-color: rgba(255,255,255,0.3); color: #fff;"
        @click="$emit('login')"
      >
        立即登录
      </GlButton>
    </view>
    
    <!-- 统计数据 -->
    <view v-if="isLoggedIn" class="user-stats">
      <view class="stat-item" @tap="$emit('stat-click', 'orders')">
        <text class="stat-value">{{ orderCount }}</text>
        <text class="stat-label">订单</text>
      </view>
      <view class="stat-item" @tap="$emit('stat-click', 'favorites')">
        <text class="stat-value">{{ favoriteCount }}</text>
        <text class="stat-label">收藏</text>
      </view>
      <view class="stat-item" @tap="$emit('stat-click', 'wallet')">
        <text class="stat-value">¥{{ balanceYuan }}</text>
        <text class="stat-label">余额</text>
      </view>
    </view>
    
    <!-- 未登录提示 -->
    <view v-else class="login-prompt">
      <text class="prompt-text">登录后查看您的订单、收藏和钱包</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlAvatar from '@/components/gl/Avatar/index.vue'
import GlButton from '@/components/gl/Button/index.vue'

interface Props {
  avatar?: string
  nickname?: string
  userId?: number
  isLoggedIn: boolean
  isPlayer?: boolean
  orderCount?: number
  favoriteCount?: number
  balance?: number // 分
}

const props = withDefaults(defineProps<Props>(), {
  orderCount: 0,
  favoriteCount: 0,
  balance: 0,
})

defineEmits<{
  edit: []
  login: []
  'stat-click': [type: 'orders' | 'favorites' | 'wallet']
}>()

const displayName = computed(() => {
  if (props.userId) {
    return props.nickname || `用户${props.userId}`
  }
  return '未登录'
})

const balanceYuan = computed(() => (props.balance / 100).toFixed(2))
</script>

<style lang="scss" scoped>
.profile-header {
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-light, #4ADE80) 100%);
  padding: 80rpx 32rpx 48rpx;
  border-radius: 0 0 48rpx 48rpx;
  position: relative;
  overflow: hidden;
  
  &::before {
    content: '';
    position: absolute;
    top: -50%;
    right: -30%;
    width: 400rpx;
    height: 400rpx;
    background: radial-gradient(circle, rgba(255, 255, 255, 0.15) 0%, transparent 70%);
    border-radius: 50%;
  }
}

.user-info {
  display: flex;
  align-items: center;
  gap: 28rpx;
  position: relative;
  z-index: 1;
}

.user-details {
  flex: 1;
}

.user-name {
  font-size: 40rpx;
  font-weight: 800;
  color: #FFFFFF;
  text-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.2);
}

.user-id {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.85);
  margin-top: 10rpx;
}

.player-badge {
  display: inline-flex;
  padding: 8rpx 20rpx;
  background: rgba(255, 255, 255, 0.2);
  border: 1rpx solid rgba(255, 255, 255, 0.3);
  border-radius: 20rpx;
  margin-top: 16rpx;
  
  text {
    font-size: 22rpx;
    font-weight: 600;
    color: #FFFFFF;
  }
}

.login-tip {
  margin-top: 16rpx;
  
  text {
    font-size: 26rpx;
    color: rgba(255, 255, 255, 0.8);
  }
}

.user-stats {
  display: flex;
  margin-top: 36rpx;
  padding-top: 36rpx;
  border-top: 1rpx solid rgba(255, 255, 255, 0.25);
  position: relative;
  z-index: 1;
}

.stat-item {
  flex: 1;
  text-align: center;
  padding: 8rpx 0;
  transition: all 0.2s;
  
  &:active {
    transform: scale(0.95);
  }
}

.stat-value {
  display: block;
  font-size: 44rpx;
  font-weight: 800;
  color: #FFFFFF;
  text-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.2);
}

.stat-label {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.85);
  margin-top: 10rpx;
  font-weight: 500;
}

.login-prompt {
  margin-top: 28rpx;
  padding-top: 28rpx;
  border-top: 1rpx solid rgba(255, 255, 255, 0.25);
  text-align: center;
  position: relative;
  z-index: 1;
}

.prompt-text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.85);
  font-weight: 500;
}
</style>
