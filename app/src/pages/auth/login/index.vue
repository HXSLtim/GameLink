<template>
  <PageShell class="login-page" padding="0">
    <view class="login-wrap">
      <view class="login-panel">
        <view class="login-body">
          <view class="login-brand">
            <!-- Logo 区域 -->
            <AuthLogo variant="icon" />
          </view>

          <view class="login-form">
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
                  custom-style="background: var(--color-wechat, #07C160); border-color: var(--color-wechat, #07C160);"
                  @click="handleWechatLogin"
                >
                  微信一键登录
                </GlButton>
              </view>
              <!-- #endif -->

              <!-- H5/App 环境显示账号登录 -->
              <!-- #ifndef MP-WEIXIN -->
              <LoginForm
                v-model:account="form.account"
                v-model:password="form.password"
                :loading="loginLoading"
                @submit="handleLogin"
              />
              <!-- #endif -->

              <!-- 分割线 -->
              <AuthDivider />

              <!-- 其他登录方式 -->
              <LoginOtherActions
                @open-account="showAccountLogin = true"
                @register="goToRegister"
              />
            </view>
          </view>
        </view>

        <!-- 底部协议 -->
        <view class="login-footer">
          <AuthAgreementFooter @agreement="goToAgreement" />
        </view>
      </view>

      <!-- 账号密码弹窗（小程序环境） -->
      <!-- #ifdef MP-WEIXIN -->
      <LoginAccountPopup
        v-model:show="showAccountLogin"
        :loading="loginLoading"
        :form="form"
        @submit="handleLogin"
      />
      <!-- #endif -->
    </view>

  </PageShell>
</template>

<script setup lang="ts">
// Pattern 组件
import GlButton from '@/components/gl/Button/index.vue'
import PageShell from '@/components/layout/PageShell/index.vue'
// Business 组件
import AuthLogo from '@/components/AuthLogo/index.vue'
import LoginForm from '@/components/LoginForm/index.vue'
import AuthDivider from '@/components/AuthDivider/index.vue'
import LoginOtherActions from '@/components/LoginOtherActions/index.vue'
import LoginAccountPopup from '@/components/LoginAccountPopup/index.vue'
import AuthAgreementFooter from '@/components/AuthAgreementFooter/index.vue'
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
.login-wrap {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 100%;
  padding: 0 48rpx;
  gap: var(--spacing-lg);
}

.login-panel {
  width: 100%;
  background: var(--color-bg-card);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--spacing-lg);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.login-body {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.login-brand {
  display: flex;
  justify-content: center;
}

.login-section {
  display: flex;
  flex-direction: column;
}

.wechat-section {
  margin-bottom: 40rpx;
}

.login-footer {
  margin-top: var(--spacing-md);
  padding-bottom: var(--spacing-md);
  text-align: center;
}

.login-page {
  @include desktop {
    margin-left: 0 !important;
    width: 100% !important;
  }
}

@include desktop {
  .login-wrap {
    justify-content: center;
    padding: 72rpx 0;
  }

  .login-panel {
    max-width: 520px;
    padding: var(--spacing-xl);
  }

  .login-body {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .login-brand {
    justify-content: center;
    padding: 0;

    :deep(.logo-section) {
      align-items: center;
      padding: 0;
    }

    :deep(.logo) {
      width: 96rpx;
      height: 96rpx;
    }
  }

  .login-form {
    padding: 0;
  }

  .login-footer {
    padding: 0;
    margin-top: var(--spacing-md);
    text-align: center;
  }
}

</style>
