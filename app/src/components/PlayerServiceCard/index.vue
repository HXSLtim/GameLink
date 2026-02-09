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
      <GlSwitch 
        :model-value="service.isOnline"
        size="small"
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
        <PriceTag class="service-value price" :amount="service.price" amount-unit="yuan" :unit="service.unit" size="small" />
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
import GlSwitch from '@/components/gl/Switch/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import type { PlayerServiceCardData } from '@/types/player'

interface Props {
  service: PlayerServiceCardData
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
  width: 100%;
}

.service-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: var(--spacing-xs);
  border-bottom: 1rpx solid var(--color-border);
  margin-bottom: var(--spacing-sm);
}

.game-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.game-icon {
  width: 48rpx;
  height: 48rpx;
  border-radius: var(--radius-sm);
}

.game-name {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  @include text-ellipsis;
}

.service-content {
  margin-bottom: var(--spacing-sm);
}

.service-row {
  display: flex;
  justify-content: space-between;
  gap: var(--spacing-sm);
  padding: var(--spacing-xs) 0;
}

.service-label {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.service-value {
  font-size: var(--font-sm);
  color: var(--color-text);
  text-align: right;
  
  &.price {
    font-weight: 600;
    color: var(--color-primary);
  }
  
  &.desc {
    max-width: 400rpx;
    text-align: right;
    @include text-ellipsis;
  }
}

.service-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-sm);
  padding-top: var(--spacing-xs);
  border-top: 1rpx solid var(--color-border);
}
</style>
