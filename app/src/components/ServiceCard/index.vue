<template>
  <GlCard :shadow="false" bordered class="service-card">
    <!-- 服务头部 -->
    <view class="service-header">
      <view class="game-info">
        <image 
          v-if="service.gameIcon" 
          :src="service.gameIcon" 
          mode="aspectFit"
          class="game-icon"
        />
        <text class="game-name">{{ service.gameName }}</text>
      </view>
      <switch 
        :checked="service.isOnline" 
        color="#00D26A"
        @change="$emit('toggle-status')"
      />
    </view>
    
    <!-- 服务内容 -->
    <view class="service-content">
      <view class="service-row">
        <text class="service-label">服务类型</text>
        <text class="service-value">{{ service.serviceName }}</text>
      </view>
      <view class="service-row">
        <text class="service-label">价格</text>
        <text class="service-value price">¥{{ service.price }}/{{ service.unit }}</text>
      </view>
      <view class="service-row">
        <text class="service-label">段位</text>
        <text class="service-value">{{ service.rankName }}</text>
      </view>
      <view v-if="service.description" class="service-row">
        <text class="service-label">服务介绍</text>
        <text class="service-value desc">{{ service.description }}</text>
      </view>
    </view>
    
    <!-- 服务操作 -->
    <view class="service-actions">
      <GlButton size="small" plain @click="$emit('edit')">编辑</GlButton>
      <GlButton size="small" plain type="error" @click="$emit('delete')">删除</GlButton>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'
import GlButton from '@/components/gl/Button/index.vue'

export interface PlayerService {
  id: number
  gameId: number
  gameName: string
  gameIcon?: string
  serviceName: string
  price: number
  unit: string
  rankName: string
  description?: string
  isOnline: boolean
}

interface Props {
  service: PlayerService
}

defineProps<Props>()

defineEmits<{
  'toggle-status': []
  edit: []
  delete: []
}>()
</script>

<style lang="scss" scoped>
.service-card {
  margin-bottom: 20rpx;
}

.service-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 20rpx;
  border-bottom: 1rpx solid var(--color-border);
  margin-bottom: 16rpx;
}

.game-info {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.game-icon {
  width: 48rpx;
  height: 48rpx;
  border-radius: 8rpx;
}

.game-name {
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
}

.service-content {
  margin-bottom: 16rpx;
}

.service-row {
  display: flex;
  justify-content: space-between;
  padding: 12rpx 0;
}

.service-label {
  font-size: 26rpx;
  color: var(--color-text-secondary);
}

.service-value {
  font-size: 26rpx;
  color: var(--color-text);
  
  &.price {
    font-weight: 600;
    color: var(--color-primary);
  }
  
  &.desc {
    max-width: 400rpx;
    text-align: right;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.service-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid var(--color-border);
}
</style>
