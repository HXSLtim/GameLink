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

export interface GameCertData {
  gameId?: number
  gameName: string
  rankId?: number
  rankName: string
  screenshot?: string
}

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

const uploadScreenshot = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    success: (res) => {
      emit('update:screenshot', res.tempFilePaths[0])
    }
  })
}
</script>

<style lang="scss" scoped>
.game-cert-item {
  padding: 20rpx;
  background: var(--color-bg-secondary);
  border-radius: 16rpx;
  margin-bottom: 16rpx;
}

.game-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
  padding-bottom: 16rpx;
  border-bottom: 1rpx solid var(--color-border);
}

.game-name {
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
}

.remove-btn {
  font-size: 26rpx;
  color: var(--color-error);
}

.game-form {
  display: flex;
  flex-direction: column;
}

.screenshot-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 0;
}

.row-label {
  font-size: 28rpx;
  color: var(--color-text);
}

.screenshot-upload {
  width: 100rpx;
  height: 100rpx;
  border-radius: 12rpx;
  overflow: hidden;
  border: 2rpx dashed var(--color-border);
  background: var(--color-bg);
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
