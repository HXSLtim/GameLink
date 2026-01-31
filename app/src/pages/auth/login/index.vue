<template>
  <view class="login-page page-container">
    <!-- Logo 区域 -->
    <AuthLogo />

    <!-- 登录方式 -->
    <view class="login-section">
      <!-- 微信一键登录（小程序环境） -->
      <!-- #ifdef MP-WEIXIN -->
      <view class="wechat-section">
        <GlButton 
          type="success"
          block
          round
          size="large"
          :loading="wechatLoading"
          custom-style="background: #07C160; border-color: #07C160;"
          @click="handleWechatLogin"
        >
          微信一键登录
        </GlButton>
      </view>
      <!-- #endif -->

      <!-- H5/App 环境显示账号登录 -->
      <!-- #ifndef MP-WEIXIN -->
      <LoginForm
        v-model:username="form.username"
        v-model:password="form.password"
        :loading="loginLoading"
        @submit="handleLogin"
      />
      <!-- #endif -->

      <!-- 分割线 -->
      <view class="divider">
        <view class="divider-line"></view>
        <text class="divider-text">或</text>
        <view class="divider-line"></view>
      </view>

      <!-- 其他登录方式 -->
      <view class="other-login">
        <!-- #ifdef MP-WEIXIN -->
        <text class="link" @click="showAccountLogin = true">使用账号密码登录</text>
        <!-- #endif -->
        
        <text class="link" @click="goToRegister">还没有账号？立即注册</text>
      </view>
    </view>

    <!-- 账号密码弹窗（小程序环境） -->
    <!-- #ifdef MP-WEIXIN -->
    <uv-popup :show="showAccountLogin" mode="bottom" round="24" @close="showAccountLogin = false">
      <view class="popup-content">
        <view class="popup-header">
          <text class="popup-title">账号登录</text>
          <uv-icon name="close" size="20" @click="showAccountLogin = false"></uv-icon>
        </view>
        
        <LoginForm
          v-model:username="form.username"
          v-model:password="form.password"
          :loading="loginLoading"
          @submit="handleLogin"
        />
      </view>
    </uv-popup>
    <!-- #endif -->

    <!-- 底部协议 -->
    <view class="footer">
      <text class="footer-text">
        登录即代表同意
        <text class="link-text" @tap="goToAgreement('user')">《用户协议》</text>
        和
        <text class="link-text" @tap="goToAgreement('privacy')">《隐私政策》</text>
      </text>
    </view>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
// Pattern 组件
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import AuthLogo from '@/components/AuthLogo/index.vue'
import LoginForm from '@/components/LoginForm/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useLogin } from '@/composables/useLogin'

const {
  loginLoading,
  wechatLoading,
  showAccountLogin,
  form,
  handleLogin,
  handleWechatLogin,
  goToRegister,
  goToAgreement,
} = useLogin()
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
}

.login-section {
  flex: 1;
  padding: 0 48rpx;
}

.wechat-section {
  margin-bottom: 40rpx;
}

.divider {
  display: flex;
  align-items: center;
  gap: 24rpx;
  margin: 48rpx 0;
}

.divider-line {
  flex: 1;
  height: 1rpx;
  background: var(--color-border);
}

.divider-text {
  font-size: 26rpx;
  color: var(--color-text-secondary);
}

.other-login {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24rpx;
}

.link {
  font-size: 28rpx;
  color: var(--color-primary);
}

.popup-content {
  padding: 32rpx 0;
  padding-bottom: calc(32rpx + env(safe-area-inset-bottom));
}

.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 48rpx 32rpx;
}

.popup-title {
  font-size: 34rpx;
  font-weight: 700;
  color: var(--color-text);
}

.footer {
  padding: 48rpx;
  text-align: center;
}

.footer-text {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.link-text {
  color: var(--color-primary);
}
</style>
