<template>
  <view class="game-cert-item">
    <view class="game-header">
      <text class="game-name">{{ game.gameName || '选择游戏' }}</text>
      <text class="remove-btn" @tap="$emit('remove')">删除</text>
    </view>
    
    <view class="game-form">
      <FormItem
        label="游戏"
        :display-value="game.gameName"
        placeholder="请选择"
        @click="$emit('select-game')"
      />
      
      <FormItem
        label="段位"
        :display-value="game.rankName"
        placeholder="请选择"
        @click="$emit('select-rank')"
      />
      
      <view class="screenshot-row">
        <text class="row-label">段位截图</text>
        <view class="screenshot-upload" @tap="uploadScreenshot">
          <image 
            v-if="game.screenshot" 
            :src="game.screenshot" 
            mode="aspectFill"
            class="screenshot-preview"
          />
          <view v-else class="screenshot-placeholder">
            <uv-icon name="camera" size="24" color="var(--color-text-secondary)"></uv-icon>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import FormItem from '@/components/FormItem/index.vue'
import type { GameCertData } from '@/types/certification'
import { useImageTools } from '@/composables/useImageTools'

interface Props {
  game: GameCertData
}

defineProps<Props>()

const emit = defineEmits<{
  'remove': []
  'select-game': []
  'select-rank': []
  'update:screenshot': [url: string]
}>()

const { pickImages } = useImageTools()

const uploadScreenshot = async () => {
  try {
    const [tempPath] = await pickImages()
    if (!tempPath) return
    emit('update:screenshot', tempPath)
  } catch {
    // ignore cancel
  }
}
</script>

<style lang="scss" scoped>
.game-cert-item {
  padding: var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  margin-bottom: var(--spacing-sm);
}

.game-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
  padding-bottom: var(--spacing-sm);
  border-bottom: 1rpx solid var(--color-border);
}

.game-name {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  @include text-ellipsis;
}

.remove-btn {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  cursor: pointer;
  @include press-effect;
}

.game-form {
  display: flex;
  flex-direction: column;
}

.screenshot-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-sm) 0;
}

.row-label {
  font-size: var(--font-sm);
  color: var(--color-text);
}

.screenshot-upload {
  width: 100rpx;
  height: 100rpx;
  border-radius: var(--radius-sm);
  overflow: hidden;
  border: 1rpx dashed var(--color-border);
  background: var(--color-bg-card);
  cursor: pointer;
  @include press-effect;
}

.screenshot-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.screenshot-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
