<template>
  <GlCard title="身份证照片" :shadow="false" bordered>
    <view class="id-card-upload">
      <view class="upload-item" @tap="uploadCard('front')">
        <image 
          v-if="frontImage" 
          :src="frontImage" 
          mode="aspectFill"
          class="upload-preview"
        />
        <view v-else class="upload-placeholder">
          <uv-icon name="camera" size="32" color="var(--color-text-secondary)"></uv-icon>
          <text class="upload-text">身份证正面</text>
        </view>
      </view>
      
      <view class="upload-item" @tap="uploadCard('back')">
        <image 
          v-if="backImage" 
          :src="backImage" 
          mode="aspectFill"
          class="upload-preview"
        />
        <view v-else class="upload-placeholder">
          <uv-icon name="camera" size="32" color="var(--color-text-secondary)"></uv-icon>
          <text class="upload-text">身份证背面</text>
        </view>
      </view>
    </view>
    
    <text class="upload-tip">请上传清晰的身份证正反面照片，确保信息可见</text>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'
import { useImageTools } from '@/composables/useImageTools'

interface Props {
  frontImage?: string
  backImage?: string
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
})

const emit = defineEmits<{
  'update:frontImage': [url: string]
  'update:backImage': [url: string]
}>()

const { pickImages } = useImageTools()

const uploadCard = async (side: 'front' | 'back') => {
  if (props.disabled) return
  try {
    const [tempPath] = await pickImages()
    if (!tempPath) return
    if (side === 'front') {
      emit('update:frontImage', tempPath)
    } else {
      emit('update:backImage', tempPath)
    }
  } catch {
    // ignore cancel
  }
}
</script>

<style lang="scss" scoped>
.id-card-upload {
  display: flex;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-sm);
}

.upload-item {
  flex: 1;
  height: 180rpx;
  border-radius: var(--radius-sm);
  overflow: hidden;
  border: 1rpx dashed var(--color-border);
  background: var(--color-bg-card);
  cursor: pointer;
  @include press-effect;
}

.upload-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.upload-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-xs);
}

.upload-text {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.upload-tip {
  font-size: var(--font-sm);
  color: var(--color-text-placeholder);
}
</style>
