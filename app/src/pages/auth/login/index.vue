<template>
  <view class="login-page" :class="{ 'theme-dark': isDark }">
    <!-- Logo 区域 -->
    <view class="logo-section">
      <image class="logo" src="/static/logo.png" mode="aspectFit" />
      <text class="app-name">GameLink</text>
      <text class="app-slogan">专业游戏陪玩平台</text>
    </view>

    <!-- 登录方式 -->
    <view class="login-section">
      <!-- 微信一键登录（小程序环境） -->
      <!-- #ifdef MP-WEIXIN -->
      <button 
        class="btn-wechat"
        :loading="wechatLoading"
        @click="handleWechatLogin"
      >
        <text class="btn-icon">微信</text>
        <text>微信一键登录</text>
      </button>
      <!-- #endif -->

      <!-- H5/App 环境显示账号登录 -->
      <!-- #ifndef MP-WEIXIN -->
      <view class="form-section">
        <view class="input-group">
          <text class="input-label">手机号/邮箱</text>
          <input
            v-model="form.username"
            class="input"
            type="text"
            placeholder="请输入手机号或邮箱"
            :placeholder-class="'input-placeholder'"
          />
        </view>

        <view class="input-group">
          <text class="input-label">密码</text>
          <input
            v-model="form.password"
            class="input"
            type="password"
            placeholder="请输入密码"
            :placeholder-class="'input-placeholder'"
          />
        </view>

        <button 
          class="btn-primary"
          :loading="loginLoading"
          :disabled="!canLogin"
          @click="handleLogin"
        >
          登录
        </button>
      </view>
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
        <view class="link" @click="showAccountLogin = true">
          <text>使用账号密码登录</text>
        </view>
        <!-- #endif -->
        
        <view class="link" @click="goToRegister">
          <text>还没有账号？立即注册</text>
        </view>
      </view>
    </view>

    <!-- 账号密码弹窗（小程序环境） -->
    <!-- #ifdef MP-WEIXIN -->
    <uv-popup v-model:show="showAccountLogin" mode="bottom" round="24">
      <view class="popup-content">
        <view class="popup-header">
          <text class="popup-title">账号登录</text>
          <view class="popup-close" @click="showAccountLogin = false">
            <text>关闭</text>
          </view>
        </view>
        
        <view class="form-section">
          <view class="input-group">
            <text class="input-label">手机号/邮箱</text>
            <input
              v-model="form.username"
              class="input"
              type="text"
              placeholder="请输入手机号或邮箱"
            />
          </view>

          <view class="input-group">
            <text class="input-label">密码</text>
            <input
              v-model="form.password"
              class="input"
              type="password"
              placeholder="请输入密码"
            />
          </view>

          <button 
            class="btn-primary"
            :loading="loginLoading"
            :disabled="!canLogin"
            @click="handleLogin"
          >
            登录
          </button>
        </view>
      </view>
    </uv-popup>
    <!-- #endif -->

    <!-- 底部协议 -->
    <view class="footer">
      <text class="footer-text">
        登录即代表同意
        <text class="link-text">《用户协议》</text>
        和
        <text class="link-text">《隐私政策》</text>
      </text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { useTheme } from '@/composables/useTheme'
import { login, doWeChatLogin } from '@/api/auth'

const userStore = useUserStore()
const { isDark } = useTheme()

// 表单数据
const form = ref({
  username: '',
  password: '',
})

// 加载状态
const loginLoading = ref(false)
const wechatLoading = ref(false)
const showAccountLogin = ref(false)

// 是否可以登录
const canLogin = computed(() => {
  return form.value.username.trim() && form.value.password.trim()
})

// 账号密码登录
async function handleLogin() {
  if (!canLogin.value) return
  
  loginLoading.value = true
  try {
    const res = await login({
      username: form.value.username.trim(),
      password: form.value.password,
    })
    
    // 保存登录状态
    userStore.login({
      accessToken: res.data.accessToken,
      refreshToken: res.data.refreshToken,
      user: res.data.user,
    })
    
    uni.showToast({ title: '登录成功', icon: 'success' })
    
    // 跳转到首页
    setTimeout(() => {
      uni.switchTab({ url: '/pages/index/index' })
    }, 500)
  } catch (error: any) {
    console.error('Login failed:', error)
    // 错误已在 request.ts 中处理
  } finally {
    loginLoading.value = false
    showAccountLogin.value = false
  }
}

// 微信一键登录
async function handleWechatLogin() {
  wechatLoading.value = true
  try {
    const res = await doWeChatLogin()
    
    // 保存登录状态
    userStore.login({
      accessToken: res.accessToken,
      refreshToken: res.refreshToken,
      user: res.user,
    })
    
    uni.showToast({ title: '登录成功', icon: 'success' })
    
    // 跳转到首页
    setTimeout(() => {
      uni.switchTab({ url: '/pages/index/index' })
    }, 500)
  } catch (error: any) {
    console.error('WeChat login failed:', error)
    uni.showToast({ title: error.message || '微信登录失败', icon: 'none' })
  } finally {
    wechatLoading.value = false
  }
}

// 跳转到注册页
function goToRegister() {
  uni.navigateTo({ url: '/pages/auth/register/index' })
}
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  background: var(--color-bg);
  padding: 0 48rpx;
  display: flex;
  flex-direction: column;
}

// Logo 区域
.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 160rpx;
  padding-bottom: 80rpx;
  
  .logo {
    width: 160rpx;
    height: 160rpx;
    margin-bottom: 24rpx;
  }
  
  .app-name {
    font-size: 48rpx;
    font-weight: 600;
    color: var(--color-primary);
    margin-bottom: 12rpx;
  }
  
  .app-slogan {
    font-size: 28rpx;
    color: var(--color-text-secondary);
  }
}

// 登录区域
.login-section {
  flex: 1;
}

// 微信登录按钮
.btn-wechat {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 96rpx;
  background: #07C160;
  color: #FFFFFF;
  font-size: 32rpx;
  font-weight: 500;
  border-radius: 48rpx;
  border: none;
  
  .btn-icon {
    margin-right: 16rpx;
  }
  
  &::after {
    border: none;
  }
}

// 表单区域
.form-section {
  padding: 24rpx 0;
}

.input-group {
  margin-bottom: 32rpx;
  
  .input-label {
    display: block;
    font-size: 28rpx;
    color: var(--color-text-secondary);
    margin-bottom: 16rpx;
  }
  
  .input {
    width: 100%;
    height: 96rpx;
    padding: 0 32rpx;
    background: var(--color-bg-card);
    border: 2rpx solid var(--color-border);
    border-radius: 16rpx;
    font-size: 32rpx;
    color: var(--color-text);
    box-sizing: border-box;
    
    &:focus {
      border-color: var(--color-primary);
    }
  }
}

.input-placeholder {
  color: var(--color-text-placeholder);
}

// 主按钮
.btn-primary {
  width: 100%;
  height: 96rpx;
  background: var(--color-primary);
  color: #FFFFFF;
  font-size: 32rpx;
  font-weight: 500;
  border-radius: 48rpx;
  border: none;
  margin-top: 16rpx;
  
  &::after {
    border: none;
  }
  
  &[disabled] {
    background: var(--color-text-disabled);
    color: #FFFFFF;
  }
}

// 分割线
.divider {
  display: flex;
  align-items: center;
  margin: 48rpx 0;
  
  .divider-line {
    flex: 1;
    height: 2rpx;
    background: var(--color-divider);
  }
  
  .divider-text {
    padding: 0 24rpx;
    font-size: 24rpx;
    color: var(--color-text-placeholder);
  }
}

// 其他登录方式
.other-login {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24rpx;
  
  .link {
    font-size: 28rpx;
    color: var(--color-primary);
  }
}

// 弹窗内容
.popup-content {
  padding: 32rpx 48rpx 64rpx;
  padding-bottom: calc(64rpx + env(safe-area-inset-bottom));
  
  .popup-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 32rpx;
    
    .popup-title {
      font-size: 36rpx;
      font-weight: 600;
      color: var(--color-text);
    }
    
    .popup-close {
      font-size: 28rpx;
      color: var(--color-text-secondary);
    }
  }
}

// 底部协议
.footer {
  padding: 48rpx 0;
  padding-bottom: calc(48rpx + env(safe-area-inset-bottom));
  text-align: center;
  
  .footer-text {
    font-size: 24rpx;
    color: var(--color-text-placeholder);
    line-height: 1.6;
  }
  
  .link-text {
    color: var(--color-primary);
  }
}
</style>
