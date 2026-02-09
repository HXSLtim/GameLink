<template>
  <BasePageLayout
    class="edit-profile-page"
    padding="0"
    title="编辑资料"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #nav>
      <!-- 顶部导航 -->
      <NavBar title="编辑资料" @back="goBack">
        <template #right>
          <GlButton 
            type="primary" 
            size="mini" 
            :disabled="!hasChanges" 
            :loading="saving"
            @click="saveProfile"
          >
            保存
          </GlButton>
        </template>
      </NavBar>
    </template>

      <!-- 头像 -->
      <AvatarUploader
        v-model="form.avatar"
        :placeholder="form.nickname?.[0]"
        @upload="handleAvatarUpload"
      />

      <!-- 基本信息 -->
      <ProfileBasicSection
        :form="form"
        :gender-text="getGenderText(form.gender)"
        @pick-gender="openGenderPicker"
        @pick-birthday="showBirthdayPicker = true"
        @pick-region="showRegionPicker = true"
      />

      <!-- 个人介绍 -->
      <IntroTextCard
        v-model="form.bio"
        title="个人介绍"
        placeholder="介绍一下自己吧~"
        :max-length="200"
      />

      <!-- 联系方式 -->
      <ProfileContactSection
        :phone="formatPhone(profile.phone)"
        :wechat-bound="profile.wechatBound"
        @change-phone="changePhone"
        @bind-wechat="bindWechat"
      />

      <!-- 游戏信息 -->
      <ProfileGamesSection
        :games-count="form.games?.length"
        @edit="editGames"
      />

      <view class="bottom-placeholder"></view>

    <template #footer>
      <!-- 性别选择器 -->
      <uv-picker
        :show="showGenderPicker"
        :columns="[genderOptions.map(o => o.label)]"
        @confirm="(e: any) => { form.gender = genderOptions[e.indexs[0]]?.value || ''; showGenderPicker = false }"
        @cancel="showGenderPicker = false"
      ></uv-picker>

      <!-- 生日选择器 -->
      <uv-datetime-picker 
        :show="showBirthdayPicker"
        mode="date"
        :min-date="minBirthday"
        :max-date="maxBirthday"
        @confirm="onBirthdayConfirm"
        @cancel="showBirthdayPicker = false"
      />
    </template>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import AvatarUploader from '@/components/AvatarUploader/index.vue'
import IntroTextCard from '@/components/IntroTextCard/index.vue'
import ProfileBasicSection from '@/components/ProfileBasicSection/index.vue'
import ProfileContactSection from '@/components/ProfileContactSection/index.vue'
import ProfileGamesSection from '@/components/ProfileGamesSection/index.vue'
// Composables
import { useProfileEdit } from '@/composables/useProfileEdit'

const {
  saving,
  showGenderPicker,
  showBirthdayPicker,
  showRegionPicker,
  form,
  profile,
  genderOptions,
  minBirthday,
  maxBirthday,
  hasChanges,
  loadProfile,
  handleAvatarUpload,
  getGenderText,
  formatPhone,
  openGenderPicker,
  onBirthdayConfirm,
  saveProfile,
  changePhone,
  bindWechat,
  editGames,
  goBack,
} = useProfileEdit()

onMounted(() => {
  loadProfile()
})
</script>

<style lang="scss" scoped>


:deep(.gl-card) {
  margin: 0 24rpx 20rpx;
}

.bottom-placeholder {
  height: 100rpx;
}
</style>
