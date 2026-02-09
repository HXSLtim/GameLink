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
        custom-style="background: var(--color-bg-secondary); border-color: var(--color-border); color: var(--color-text);"
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
        custom-style="background: var(--color-bg-secondary); border-color: var(--color-border); color: var(--color-text);"
        @click="$emit('login')"
      >
        立即登录
      </GlButton>
    </view>
    
    <!-- 统计数据 -->
    <HeaderStatsRow
      v-if="isLoggedIn"
      class="user-stats"
      :items="stats"
      :show-divider="false"
      clickable
      size="lg"
      @item-click="handleStatClick"
    >
      <template #value="{ item }">
        <PriceTag
          v-if="item.key === 'wallet'"
          :amount="balance"
          amount-unit="cents"
          size="small"
        />
        <text v-else>{{ item.value }}</text>
      </template>
    </HeaderStatsRow>
    
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
import PriceTag from '@/components/PriceTag/index.vue'
import HeaderStatsRow from '@/components/HeaderStatsRow/index.vue'
import type { HeaderStatItem } from '@/types/ui'
import type { ProfileStatKey } from '@/types/profile'

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

const displayName = computed(() => {
  if (props.userId) {
    return props.nickname || `用户${props.userId}`
  }
  return '未登录'
})

const stats = computed<Array<HeaderStatItem & { key: ProfileStatKey }>>(() => [
  { key: 'orders', label: '订单', value: props.orderCount || 0 },
  { key: 'favorites', label: '收藏', value: props.favoriteCount || 0 },
  { key: 'wallet', label: '余额', value: props.balance || 0 },
])

const emit = defineEmits<{
  edit: []
  login: []
  'stat-click': [type: ProfileStatKey]
}>()

const handleStatClick = (item: HeaderStatItem & { key?: ProfileStatKey }) => {
  if (!item.key) return
  emit('stat-click', item.key)
}

</script>

<style lang="scss" scoped>
.profile-header {
  background: linear-gradient(180deg, var(--color-bg-secondary) 0%, var(--color-bg-card) 100%);
  padding: var(--spacing-lg) var(--spacing-lg) var(--spacing-md);
  border-radius: 0 0 var(--radius-lg) var(--radius-lg);
  border-bottom: 1rpx solid var(--color-border);
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  position: relative;
  z-index: 1;
}

.user-details {
  flex: 1;
}

.user-name {
  font-size: var(--font-lg);
  font-weight: 700;
  color: var(--color-text);
}

.user-id {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  margin-top: var(--spacing-xs);
}

.player-badge {
  display: inline-flex;
  padding: 2rpx var(--spacing-sm);
  background: var(--color-bg-secondary);
  border: 1rpx solid var(--color-primary);
  border-radius: var(--radius-sm);
  margin-top: var(--spacing-sm);
  
  text {
    font-size: var(--font-xs);
    font-weight: 600;
    color: var(--color-primary);
  }
}

.login-tip {
  margin-top: var(--spacing-sm);
  
  text {
    font-size: var(--font-sm);
    color: var(--color-text-secondary);
  }
}

.user-stats {
  margin-top: var(--spacing-md);
  padding-top: var(--spacing-md);
  border-top: 1rpx solid var(--color-border);
}

.login-prompt {
  margin-top: var(--spacing-md);
  padding-top: var(--spacing-md);
  border-top: 1rpx solid var(--color-border);
  text-align: center;
}

.prompt-text {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  font-weight: 500;
}
</style>
