<template>
  <view class="edit-profile-page page-container">
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

    <scroll-view class="content-scroll" scroll-y>
      <!-- 头像 -->
      <AvatarUploader
        v-model="form.avatar"
        :placeholder="form.nickname?.[0]"
        @upload="handleAvatarUpload"
      />

      <!-- 基本信息 -->
      <FormSection title="基本信息">
        <FormItem
          v-model="form.nickname"
          label="昵称"
          type="input"
          placeholder="请输入昵称"
          :maxlength="20"
        />
        <FormItem
          label="性别"
          :display-value="getGenderText(form.gender)"
          @click="openGenderPicker"
        />
        <FormItem
          label="生日"
          :display-value="form.birthday"
          placeholder="请选择"
          @click="showBirthdayPicker = true"
        />
        <FormItem
          label="地区"
          :display-value="form.region"
          placeholder="请选择"
          @click="showRegionPicker = true"
        />
      </FormSection>

      <!-- 个人介绍 -->
      <GlCard title="个人介绍" :shadow="false" bordered>
        <textarea 
          v-model="form.bio"
          class="bio-textarea"
          placeholder="介绍一下自己吧~"
          :maxlength="200"
        />
        <text class="char-count">{{ form.bio?.length || 0 }}/200</text>
      </GlCard>

      <!-- 联系方式 -->
      <FormSection title="联系方式">
        <FormItem
          label="手机号"
          :display-value="formatPhone(profile.phone)"
          :clickable="false"
        >
          <template #extra>
            <text class="change-btn" @tap="changePhone">更换</text>
          </template>
        </FormItem>
        <FormItem
          label="微信"
          :display-value="profile.wechatBound ? '已绑定' : '未绑定'"
          :clickable="false"
        >
          <template #extra>
            <text v-if="!profile.wechatBound" class="bind-btn" @tap="bindWechat">绑定</text>
          </template>
        </FormItem>
      </FormSection>

      <!-- 游戏信息 -->
      <FormSection title="游戏信息">
        <FormItem
          label="常玩游戏"
          :display-value="form.games?.length ? `${form.games.length}款` : undefined"
          placeholder="请选择"
          @click="editGames"
        />
      </FormSection>

      <view class="bottom-placeholder"></view>
    </scroll-view>

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

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import GlCard from '@/components/gl/Card/index.vue'
import FormSection from '@/components/FormSection/index.vue'
import FormItem from '@/components/FormItem/index.vue'
// Business 组件
import AvatarUploader from '@/components/AvatarUploader/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
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
.edit-profile-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  
  @include desktop {
    height: 100vh;
    min-height: auto;
    overflow: hidden;
  }
}

.content-scroll {
  flex: 1;
  overflow-y: auto;
}

:deep(.gl-card) {
  margin: 0 24rpx 20rpx;
}

.bio-textarea {
  width: 100%;
  height: 160rpx;
  padding: 20rpx;
  background: var(--color-bg-secondary);
  border-radius: 12rpx;
  border: 2rpx solid var(--color-border);
  font-size: 28rpx;
  color: var(--color-text);
}

.char-count {
  display: block;
  text-align: right;
  font-size: 22rpx;
  color: var(--color-text-placeholder);
  margin-top: 8rpx;
}

.change-btn,
.bind-btn {
  font-size: 26rpx;
  color: var(--color-primary);
  font-weight: 500;
}

.bottom-placeholder {
  height: 100rpx;
}
</style>
