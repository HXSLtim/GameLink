<template>
  <PageShell class="register-page" padding="24rpx 48rpx">
    <template #header>
      <!-- 顶部导航 -->
      <NavBar title="注册账号" @back="goBack" />
    </template>

    <!-- 角色选择 -->
    <view class="section-wrap">
      <RoleSelector
        v-model="form.role"
        :roles="roleOptions"
      />
    </view>

    <!-- 表单 -->
    <RegisterForm
      :form="form"
      v-model:agreed="agreed"
      :loading="loading"
      :can-submit="canRegister"
      @submit="handleRegister"
      @go-login="goToLogin"
      @agreement="goToAgreement"
    />

    <template #footer>
      <!-- PC 端侧边栏 -->
      <CustomTabBar :show-mobile-tab-bar="false" />
    </template>
  </PageShell>
</template>

<script setup lang="ts">
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import PageShell from '@/components/layout/PageShell/index.vue'
// Business 组件
import RoleSelector from '@/components/RoleSelector/index.vue'
import RegisterForm from '@/components/RegisterForm/index.vue'
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
.section-wrap {
  margin-bottom: 24rpx;
}

</style>
