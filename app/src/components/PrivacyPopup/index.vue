<template>
  <view v-if="visible" class="privacy-popup-mask" @tap.stop="handleMaskClick">
    <view class="privacy-popup" @tap.stop>
      <view class="popup-header">
        <text class="popup-title">用户隐私保护提示</text>
      </view>
      
      <view class="popup-content">
        <text class="content-text">
          在您使用 GameLink 服务前，请仔细阅读
          <text class="link" @tap="openPrivacyContract">{{ privacyContractName }}</text>
          。如您同意，请点击"同意"开始使用。
        </text>
      </view>
      
      <view class="popup-actions">
        <view class="action-btn cancel" @tap="handleDisagree">
          <text>不同意</text>
        </view>
        <button 
          class="action-btn confirm" 
          id="agree-btn"
          open-type="agreePrivacyAuthorization"
          @agreeprivacyauthorization="handleAgree"
        >
          <text>同意</text>
        </button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

// 声明 wx 全局变量
declare const wx: any

interface Props {
  visible?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
})

const emit = defineEmits<{
  agree: []
  disagree: []
  'update:visible': [value: boolean]
}>()

const privacyContractName = ref('《用户隐私保护指引》')

// 获取隐私协议名称
onMounted(() => {
  // #ifdef MP-WEIXIN
  if (wx.getPrivacySetting) {
    wx.getPrivacySetting({
      success: (res: any) => {
        if (res.privacyContractName) {
          privacyContractName.value = res.privacyContractName
        }
      }
    })
  }
  // #endif
})

// 打开隐私协议
const openPrivacyContract = () => {
  // #ifdef MP-WEIXIN
  if (wx.openPrivacyContract) {
    wx.openPrivacyContract({
      fail: () => {
        // 降级跳转到自定义协议页
        uni.navigateTo({ url: '/pages/agreement/index?type=privacy' })
      }
    })
  } else {
    uni.navigateTo({ url: '/pages/agreement/index?type=privacy' })
  }
  // #endif
  
  // #ifndef MP-WEIXIN
  uni.navigateTo({ url: '/pages/agreement/index?type=privacy' })
  // #endif
}

// 同意
const handleAgree = () => {
  emit('update:visible', false)
  emit('agree')
}

// 不同意
const handleDisagree = () => {
  emit('update:visible', false)
  emit('disagree')
}

// 点击遮罩
const handleMaskClick = () => {
  // 不允许点击遮罩关闭
}
</script>

<style lang="scss" scoped>
.privacy-popup-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-xl);
}

.privacy-popup {
  width: 100%;
  max-width: 600rpx;
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
  overflow: hidden;
  
  .popup-header {
    padding: var(--spacing-md) var(--spacing-md) var(--spacing-sm);
    text-align: center;
    
    .popup-title {
      font-size: var(--font-md);
      font-weight: 600;
      color: var(--color-text);
    }
  }
  
  .popup-content {
    padding: 0 var(--spacing-md) var(--spacing-md);
    
    .content-text {
      font-size: var(--font-sm);
      color: var(--color-text-secondary);
      line-height: 1.6;
      
      .link {
        color: var(--color-primary);
      }
    }
  }
  
  .popup-actions {
    display: flex;
    border-top: 1rpx solid var(--color-border);
    
    .action-btn {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      height: 80rpx;
      font-size: var(--font-sm);
      background: transparent;
      border: none;
      border-radius: 0;
      cursor: pointer;
      @include press-effect;
      
      &::after {
        display: none;
      }
      
      &.cancel {
        color: var(--color-text-secondary);
        border-right: 1rpx solid var(--color-border);
      }
      
      &.confirm {
        color: var(--color-primary);
        font-weight: 600;
      }
    }
  }
}
</style>
