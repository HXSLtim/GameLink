<template>
  <GlCard :title="title" :shadow="false" bordered class="section-card">
    <view class="settings-list">
      <view 
        v-for="item in items" 
        :key="item.key" 
        class="settings-item"
        @tap="handleClick(item)"
      >
        <view class="item-icon" v-if="item.icon">
          <uv-icon :name="item.icon" size="18" :color="item.iconColor || 'var(--color-text-secondary)'"></uv-icon>
        </view>
        <text class="item-label">{{ item.label }}</text>
        
        <!-- 右侧内容 -->
        <text v-if="item.value" class="item-value">{{ item.value }}</text>
        
        <switch 
          v-if="item.type === 'switch'"
          :checked="item.checked"
          color="#00D26A"
          @change="(e: any) => $emit('switch', item.key, e.detail.value)"
        />
        
        <uv-icon 
          v-else-if="item.type !== 'switch'" 
          name="arrow-right" 
          size="14" 
          color="var(--color-text-secondary)"
        ></uv-icon>
      </view>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'

export interface SettingsItem {
  key: string
  label: string
  icon?: string
  iconColor?: string
  value?: string
  type?: 'link' | 'switch'
  checked?: boolean
}

interface Props {
  title: string
  items: SettingsItem[]
}

defineProps<Props>()

const emit = defineEmits<{
  click: [key: string]
  switch: [key: string, value: boolean]
}>()

const handleClick = (item: SettingsItem) => {
  if (item.type !== 'switch') {
    emit('click', item.key)
  }
}
</script>

<style lang="scss" scoped>
.section-card {
  margin: 0 24rpx 20rpx;
}

.settings-list {
  display: flex;
  flex-direction: column;
}

.settings-item {
  display: flex;
  align-items: center;
  padding: 28rpx 0;
  border-bottom: 1rpx solid var(--color-border);
  
  &:last-child {
    border-bottom: none;
  }
}

.item-icon {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
}

.item-label {
  flex: 1;
  font-size: 30rpx;
  color: var(--color-text);
}

.item-value {
  font-size: 28rpx;
  color: var(--color-text-secondary);
  margin-right: 16rpx;
}
</style>
