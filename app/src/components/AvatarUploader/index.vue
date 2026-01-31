<template>
  <view class="avatar-section" @tap="handleUpload">
    <view class="avatar-wrap">
      <GlAvatar :src="modelValue" :text="placeholder" :size="160" />
      <view class="avatar-edit">
        <uv-icon name="camera" size="20" color="#fff"></uv-icon>
      </view>
    </view>
    <text class="avatar-tip">{{ tip }}</text>
  </view>
</template>

<script setup lang="ts">
import GlAvatar from '@/components/gl/Avatar/index.vue'

interface Props {
  modelValue?: string
  placeholder?: string
  tip?: string
}

withDefaults(defineProps<Props>(), {
  placeholder: 'U',
  tip: '点击更换头像',
})

const emit = defineEmits<{
  'update:modelValue': [url: string]
  upload: []
}>()

const handleUpload = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const tempPath = res.tempFilePaths[0]
      emit('update:modelValue', tempPath)
      emit('upload')
    }
  })
}
</script>

<style lang="scss" scoped>
.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48rpx 24rpx;
}

.avatar-wrap {
  position: relative;
}

.avatar-edit {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 48rpx;
  height: 48rpx;
  background: var(--color-primary);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 4rpx solid var(--color-bg-card);
}

.avatar-tip {
  margin-top: 16rpx;
  font-size: 24rpx;
  color: var(--color-text-secondary);
}
</style>
