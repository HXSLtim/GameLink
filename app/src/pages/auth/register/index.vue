<template>
  <view class="register-page" :class="{ 'theme-dark': isDark }">
    <!-- 返回按钮 -->
    <view class="nav-bar safe-area-top">
      <view class="nav-back" @click="goBack">
        <text class="nav-back-icon">返回</text>
      </view>
      <text class="nav-title">注册账号</text>
      <view class="nav-placeholder"></view>
    </view>

    <!-- 注册表单 -->
    <view class="register-content">
      <!-- 角色选择 -->
      <view class="role-section">
        <text class="section-title">选择注册身份</text>
        <view class="role-cards">
          <view 
            class="role-card"
            :class="{ active: form.role === 'user' }"
            @click="form.role = 'user'"
          >
            <view class="role-icon">用户</view>
            <text class="role-name">普通用户</text>
            <text class="role-desc">找陪玩、享受游戏乐趣</text>
          </view>
          
          <view 
            class="role-card"
            :class="{ active: form.role === 'player' }"
            @click="form.role = 'player'"
          >
            <view class="role-icon">陪</view>
            <text class="role-name">陪玩师</text>
            <text class="role-desc">提供服务、赚取收入</text>
          </view>
        </view>
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
            :placeholder-class="'input-placeholder'"
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
            :placeholder-class="'input-placeholder'"
          />
        </view>

        <view class="input-group">
          <text class="input-label">密码</text>
          <input
            v-model="form.password"
            class="input"
            type="password"
            placeholder="请设置密码（6-20位）"
            :placeholder-class="'input-placeholder'"
          />
        </view>

        <view class="input-group">
          <text class="input-label">确认密码</text>
          <input
            v-model="form.confirmPassword"
            class="input"
            type="password"
            placeholder="请再次输入密码"
            :placeholder-class="'input-placeholder'"
          />
        </view>

        <!-- 协议勾选 -->
        <view class="agreement" @click="agreed = !agreed">
          <view class="checkbox" :class="{ checked: agreed }">
            <text v-if="agreed">✓</text>
          </view>
          <text class="agreement-text">
            我已阅读并同意
            <text class="link-text">《用户协议》</text>
            和
            <text class="link-text">《隐私政策》</text>
          </text>
        </view>

        <button 
          class="btn-primary"
          :loading="loading"
          :disabled="!canRegister"
          @click="handleRegister"
        >
          注册
        </button>

        <view class="login-link" @click="goToLogin">
          <text>已有账号？去登录</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { useTheme } from '@/composables/useTheme'
import { register } from '@/api/auth'

const userStore = useUserStore()
const { isDark } = useTheme()

// 表单数据
const form = ref({
  phone: '',
  nickname: '',
  password: '',
  confirmPassword: '',
  role: 'user' as 'user' | 'player',
})

// 状态
const loading = ref(false)
const agreed = ref(false)

// 是否可以注册
const canRegister = computed(() => {
  const { phone, nickname, password, confirmPassword } = form.value
  return (
    phone.length === 11 &&
    nickname.trim().length >= 2 &&
    password.length >= 6 &&
    password === confirmPassword &&
    agreed.value
  )
})

// 注册
async function handleRegister() {
  if (!canRegister.value) {
    // 具体错误提示
    if (form.value.phone.length !== 11) {
      uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
      return
    }
    if (form.value.nickname.trim().length < 2) {
      uni.showToast({ title: '昵称至少2个字符', icon: 'none' })
      return
    }
    if (form.value.password.length < 6) {
      uni.showToast({ title: '密码至少6位', icon: 'none' })
      return
    }
    if (form.value.password !== form.value.confirmPassword) {
      uni.showToast({ title: '两次密码不一致', icon: 'none' })
      return
    }
    if (!agreed.value) {
      uni.showToast({ title: '请同意用户协议', icon: 'none' })
      return
    }
    return
  }
  
  loading.value = true
  try {
    const res = await register({
      phone: form.value.phone,
      nickname: form.value.nickname.trim(),
      password: form.value.password,
      role: form.value.role,
    })
    
    // 保存登录状态
    userStore.login({
      accessToken: res.data.accessToken,
      refreshToken: res.data.refreshToken,
      user: res.data.user,
    })
    
    uni.showToast({ title: '注册成功', icon: 'success' })
    
    // 跳转到首页
    setTimeout(() => {
      uni.switchTab({ url: '/pages/index/index' })
    }, 500)
  } catch (error: any) {
    console.error('Register failed:', error)
  } finally {
    loading.value = false
  }
}

// 返回
function goBack() {
  uni.navigateBack()
}

// 去登录
function goToLogin() {
  uni.navigateBack()
}
</script>

<style lang="scss" scoped>
.register-page {
  min-height: 100vh;
  background: var(--color-bg);
}

// 导航栏
.nav-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
  padding: 0 32rpx;
  background: var(--color-bg-card);
  border-bottom: 2rpx solid var(--color-divider);
  
  .nav-back {
    width: 100rpx;
    
    .nav-back-icon {
      font-size: 28rpx;
      color: var(--color-primary);
    }
  }
  
  .nav-title {
    font-size: 32rpx;
    font-weight: 600;
    color: var(--color-text);
  }
  
  .nav-placeholder {
    width: 100rpx;
  }
}

// 内容区
.register-content {
  padding: 32rpx 48rpx;
}

// 角色选择
.role-section {
  margin-bottom: 48rpx;
  
  .section-title {
    display: block;
    font-size: 28rpx;
    color: var(--color-text-secondary);
    margin-bottom: 24rpx;
  }
}

.role-cards {
  display: flex;
  gap: 24rpx;
}

.role-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32rpx 24rpx;
  background: var(--color-bg-card);
  border: 2rpx solid var(--color-border);
  border-radius: 24rpx;
  transition: all 0.2s;
  
  .role-icon {
    width: 80rpx;
    height: 80rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-bg-secondary);
    border-radius: 50%;
    font-size: 28rpx;
    color: var(--color-text-secondary);
    margin-bottom: 16rpx;
  }
  
  .role-name {
    font-size: 30rpx;
    font-weight: 500;
    color: var(--color-text);
    margin-bottom: 8rpx;
  }
  
  .role-desc {
    font-size: 22rpx;
    color: var(--color-text-placeholder);
    text-align: center;
  }
  
  &.active {
    border-color: var(--color-primary);
    background: rgba(0, 210, 106, 0.05);
    
    .role-icon {
      background: var(--color-primary);
      color: #FFFFFF;
    }
    
    .role-name {
      color: var(--color-primary);
    }
  }
}

// 夜间模式
.theme-dark .role-card.active {
  background: rgba(124, 58, 237, 0.1);
}

// 表单
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

// 协议勾选
.agreement {
  display: flex;
  align-items: flex-start;
  margin: 32rpx 0;
  
  .checkbox {
    width: 40rpx;
    height: 40rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2rpx solid var(--color-border);
    border-radius: 8rpx;
    margin-right: 16rpx;
    margin-top: 4rpx;
    font-size: 24rpx;
    color: #FFFFFF;
    
    &.checked {
      background: var(--color-primary);
      border-color: var(--color-primary);
    }
  }
  
  .agreement-text {
    flex: 1;
    font-size: 24rpx;
    color: var(--color-text-secondary);
    line-height: 1.6;
  }
  
  .link-text {
    color: var(--color-primary);
  }
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

// 登录链接
.login-link {
  text-align: center;
  margin-top: 32rpx;
  font-size: 28rpx;
  color: var(--color-primary);
}
</style>
