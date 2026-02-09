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
import { useImageTools } from '@/composables/useImageTools'

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

const { pickImages } = useImageTools()

const handleUpload = async () => {
  try {
    const [tempPath] = await pickImages()
    if (!tempPath) return
    emit('update:modelValue', tempPath)
    emit('upload')
  } catch {
    // ignore cancel
  }
}
</script>

<style lang="scss" scoped>
.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--spacing-xl) var(--spacing-md);
  cursor: pointer;
  @include press-effect;
}

.avatar-wrap {
  position: relative;
}

.avatar-edit {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 40rpx;
  height: 40rpx;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid var(--color-bg-card);
}

.avatar-tip {
  margin-top: var(--spacing-sm);
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}
</style>
