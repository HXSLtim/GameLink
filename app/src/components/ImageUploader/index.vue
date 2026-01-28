<template>
  <view class="image-uploader">
    <!-- 已上传的图片 -->
    <view 
      v-for="(image, index) in modelValue" 
      :key="index"
      class="image-item"
    >
      <image :src="image" mode="aspectFill" @click="previewImage(index)" />
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

// 选择图片
function chooseImage() {
  const count = props.maxCount - props.modelValue.length
  if (count <= 0) return
  
  uni.chooseImage({
    count,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: async (res) => {
      const tempFilePaths = res.tempFilePaths || []
      
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
    }
  })
}

// 预览图片
function previewImage(index: number) {
  uni.previewImage({
    urls: props.modelValue,
    current: index
  })
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
  gap: 16rpx;
}

.image-item {
  position: relative;
  width: 200rpx;
  height: 200rpx;
  border-radius: 16rpx;
  overflow: hidden;
  
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
    border-radius: 50%;
    
    text {
      color: #FFFFFF;
      font-size: 28rpx;
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
  background: var(--color-bg-secondary);
  border: 2rpx dashed var(--color-border);
  border-radius: 16rpx;
  
  .upload-icon {
    width: 64rpx;
    height: 64rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-bg-card);
    border-radius: 50%;
    margin-bottom: 12rpx;
    
    text {
      font-size: 40rpx;
      color: var(--color-text-placeholder);
      line-height: 1;
    }
  }
  
  .upload-text {
    font-size: 24rpx;
    color: var(--color-text-placeholder);
  }
}

.count-tip {
  width: 100%;
  text-align: right;
  
  text {
    font-size: 22rpx;
    color: var(--color-text-placeholder);
  }
}
</style>
