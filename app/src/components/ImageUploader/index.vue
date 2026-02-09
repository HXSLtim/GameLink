<template>
  <view class="image-uploader">
    <!-- 已上传的图片 -->
    <view 
      v-for="(image, index) in modelValue" 
      :key="index"
      class="image-item"
    >
      <image :src="image" mode="aspectFill" @click="previewImages(modelValue, index)" />
      <view v-if="!readonly" class="delete-btn" @click="removeImage(index)">
        <text>×</text>
      </view>
    </view>

    <!-- 上传按钮 -->
    <view 
      v-if="!readonly && modelValue.length < maxCount"
      class="upload-btn"
      @click="chooseImage"
    >
      <view class="upload-icon">
        <text>+</text>
      </view>
      <text class="upload-text">{{ uploadText }}</text>
    </view>

    <!-- 数量提示 -->
    <view v-if="showCount && maxCount > 1" class="count-tip">
      <text>{{ modelValue.length }}/{{ maxCount }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { uploadFile } from '@/api/request'
import { useImageTools } from '@/composables/useImageTools'

const props = withDefaults(defineProps<{
  modelValue: string[]
  maxCount?: number
  maxSize?: number  // 最大文件大小（MB）
  uploadText?: string
  uploadUrl?: string
  readonly?: boolean
  showCount?: boolean
}>(), {
  maxCount: 9,
  maxSize: 5,
  uploadText: '上传图片',
  uploadUrl: '/upload/image',
  readonly: false,
  showCount: true,
})

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

const { pickImages, previewImages } = useImageTools()

// 选择图片
async function chooseImage() {
  const count = props.maxCount - props.modelValue.length
  if (count <= 0) return

  try {
    const tempFilePaths = await pickImages({ count })
    for (const filePath of tempFilePaths) {
      try {
        // 上传图片
        uni.showLoading({ title: '上传中...' })
        const result = await uploadFile(props.uploadUrl, filePath)
        uni.hideLoading()
        
        // 更新图片列表
        const newImages = [...props.modelValue, result.data.filePath]
        emit('update:modelValue', newImages)
      } catch (error) {
        uni.hideLoading()
        console.error('Upload failed:', error)
      }
    }
  } catch {
    // ignore cancel
  }
}

// 删除图片
function removeImage(index: number) {
  const newImages = [...props.modelValue]
  newImages.splice(index, 1)
  emit('update:modelValue', newImages)
}
</script>

<style lang="scss" scoped>
.image-uploader {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}

.image-item {
  position: relative;
  width: 200rpx;
  height: 200rpx;
  border-radius: var(--radius-sm);
  overflow: hidden;
  cursor: pointer;
  
  image {
    width: 100%;
    height: 100%;
  }
  
  .delete-btn {
    position: absolute;
    top: 8rpx;
    right: 8rpx;
    width: 40rpx;
    height: 40rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.5);
    border-radius: var(--radius-full);
    @include press-effect;
    
    text {
      color: #FFFFFF;
      font-size: var(--font-md);
      line-height: 1;
    }
  }
}

.upload-btn {
  width: 200rpx;
  height: 200rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-card);
  border: 1rpx dashed var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  @include press-effect;
  
  .upload-icon {
    width: 56rpx;
    height: 56rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-bg-secondary);
    border-radius: var(--radius-sm);
    border: 1rpx solid var(--color-border);
    margin-bottom: var(--spacing-xs);
    
    text {
      font-size: var(--font-xl);
      color: var(--color-text-secondary);
      line-height: 1;
    }
  }
  
  .upload-text {
    font-size: var(--font-sm);
    color: var(--color-text-placeholder);
  }
}

.count-tip {
  width: 100%;
  text-align: right;
  
  text {
    font-size: var(--font-xs);
    color: var(--color-text-placeholder);
  }
}
</style>
