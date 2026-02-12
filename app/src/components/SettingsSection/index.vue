<template>
  <SectionCard :title="title">
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
        
        <GlSwitch 
          v-if="item.type === 'switch'"
          :model-value="!!item.checked"
          size="small"
          @update:modelValue="(value) => $emit('switch', item.key, value)"
        />
        
        <uv-icon 
          v-else
          name="arrow-right" 
          size="14" 
          color="var(--color-text-secondary)"
        ></uv-icon>
      </view>
    </view>
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard/index.vue'
import GlSwitch from '@/components/gl/Switch/index.vue'
import type { SettingsItem } from '@/types/ui'

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
.settings-list {
  display: flex;
  flex-direction: column;
}

.settings-item {
  display: flex;
  align-items: center;
  padding: var(--spacing-sm) 0;
  border-bottom: 1rpx solid var(--color-border);
  transition: background 0.2s ease;
  @include press-effect;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:active {
    background: var(--color-bg-secondary);
  }
}

.item-icon {
  width: 40rpx;
  height: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: var(--spacing-sm);
}

.item-label {
  flex: 1;
  font-size: var(--font-md);
  color: var(--color-text);
}

.item-value {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  margin-right: var(--spacing-sm);
}
</style>
