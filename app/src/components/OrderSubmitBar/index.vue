<template>
  <view class="submit-bar-shell">
    <view class="submit-bar">
      <view class="submit-bar__content">
        <view class="price-section">
          <text class="label">合计</text>
          <view class="price-wrap">
            <text class="symbol">¥</text>
            <text class="amount">{{ total }}</text>
          </view>
        </view>

        <GlButton
          class="submit-btn"
          type="primary"
          size="large"
          round
          :disabled="disabled"
          :loading="loading"
          @click="$emit('submit')"
        >
          {{ loading ? '提交中...' : buttonText }}
        </GlButton>
      </view>
    </view>
    <!-- 占位符，防止内容被遮挡 -->
    <view class="submit-bar-placeholder"></view>
  </view>
</template>

<script setup lang="ts">
import GlButton from '@/components/gl/Button/index.vue'

interface Props {
  total: number
  disabled?: boolean
  loading?: boolean
  buttonText?: string
}

withDefaults(defineProps<Props>(), {
  disabled: false,
  loading: false,
  buttonText: '立即支付',
})

defineEmits<{
  submit: []
}>()
</script>

<style lang="scss" scoped>
.submit-bar-shell {
  position: relative;
  z-index: 100;
}

.submit-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-top: 1rpx solid rgba(0, 0, 0, 0.05);
  box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.03);
  padding-bottom: env(safe-area-inset-bottom);
  transition: all 0.3s ease;
}

.submit-bar__content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 112rpx; // 56px
  padding: 0 $gl-spacing-lg;
}

.price-section {
  display: flex;
  align-items: baseline;
  gap: $gl-spacing-xs;

  .label {
    font-size: 26rpx;
    color: $gl-text-secondary;
  }

  .price-wrap {
    color: #FF5252; // 价格强调色
    font-weight: 700;
    font-family: 'DIN Alternate', sans-serif;
    line-height: 1;

    .symbol {
      font-size: 28rpx;
      margin-right: 4rpx;
    }

    .amount {
      font-size: 44rpx;
    }
  }
}

.submit-btn {
  min-width: 240rpx;
  font-weight: 600;
  box-shadow: 0 8rpx 20rpx rgba($gl-color-primary, 0.25);
}

.submit-bar-placeholder {
  height: calc(112rpx + env(safe-area-inset-bottom));
  width: 100%;
}
</style>
