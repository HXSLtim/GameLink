<template>
  <view class="register-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="注册账号" @back="goBack" />

    <!-- 注册表单 -->
    <scroll-view class="register-content" scroll-y>
      <!-- 角色选择 -->
      <view class="section-wrap">
        <RoleSelector
          v-model="form.role"
          :roles="roleOptions"
        />
      </view>

      <!-- 表单 -->
      <view class="form-section">
        <view class="input-group">
          <text class="input-label">手机号</text>
          <input
            v-model="form.phone"
            class="input"
            type="number"
            maxlength="11"
            placeholder="请输入手机号"
          />
        </view>

        <view class="input-group">
          <text class="input-label">昵称</text>
          <input
            v-model="form.nickname"
            class="input"
            type="text"
            maxlength="20"
            placeholder="请输入昵称"
          />
        </view>

        <view class="input-group">
          <text class="input-label">密码</text>
          <input
            v-model="form.password"
            class="input"
            type="password"
            placeholder="请设置密码（6-20位）"
          />
        </view>

        <view class="input-group">
          <text class="input-label">确认密码</text>
          <input
            v-model="form.confirmPassword"
            class="input"
            type="password"
            placeholder="请再次输入密码"
          />
        </view>

        <!-- 协议勾选 -->
        <view class="agreement" @click="agreed = !agreed">
          <view class="checkbox" :class="{ checked: agreed }">
            <uv-icon v-if="agreed" name="checkbox-mark" size="14" color="#fff"></uv-icon>
          </view>
          <text class="agreement-text">
            我已阅读并同意
            <text class="link-text" @tap.stop="goToAgreement('user')">《用户协议》</text>
            和
            <text class="link-text" @tap.stop="goToAgreement('privacy')">《隐私政策》</text>
          </text>
        </view>

        <GlButton 
          type="primary"
          block
          round
          size="large"
          :loading="loading"
          :disabled="!canRegister"
          @click="handleRegister"
        >
          注册
        </GlButton>

        <view class="login-link" @click="goToLogin">
          <text>已有账号？去登录</text>
        </view>
      </view>
    </scroll-view>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import RoleSelector from '@/components/RoleSelector/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useRegister } from '@/composables/useRegister'

const {
  loading,
  agreed,
  form,
  roleOptions,
  canRegister,
  handleRegister,
  goBack,
  goToLogin,
  goToAgreement,
} = useRegister()
</script>

<style lang="scss" scoped>
.register-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
}

.register-content {
  flex: 1;
  padding: 24rpx 48rpx;
}

.section-wrap {
  margin-bottom: 24rpx;
}

.form-section {
  padding-top: 16rpx;
}

.input-group {
  margin-bottom: 32rpx;
}

.input-label {
  display: block;
  font-size: 28rpx;
  font-weight: 500;
  color: var(--color-text);
  margin-bottom: 16rpx;
}

.input {
  width: 100%;
  height: 88rpx;
  padding: 0 24rpx;
  background: var(--color-bg-secondary);
  border: 2rpx solid var(--color-border);
  border-radius: 16rpx;
  font-size: 30rpx;
  color: var(--color-text);
}

.agreement {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 32rpx;
}

.checkbox {
  width: 36rpx;
  height: 36rpx;
  border-radius: 8rpx;
  border: 2rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  
  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
  }
}

.agreement-text {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.link-text {
  color: var(--color-primary);
}

.login-link {
  margin-top: 32rpx;
  text-align: center;
  font-size: 28rpx;
  color: var(--color-primary);
}
</style>
