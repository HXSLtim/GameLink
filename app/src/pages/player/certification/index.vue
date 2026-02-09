<template>
  <BasePageLayout
    class="certification-page"
    padding="0"
    title="陪玩认证"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #nav>
      <!-- 顶部导航 -->
      <NavBar title="陪玩认证" @back="goBack" />
    </template>

    <!-- 认证状态卡片 -->
    <CertStatusCard :status="certStatus" />

    <!-- 基本信息 -->
    <CertificationBasicSection
      :form="form"
      :gender-text="getGenderText(form.gender)"
      @pick-gender="showGenderPicker = true"
    />

    <!-- 身份证照片 -->
    <IdCardUploader
      v-model:front-image="form.idCardFront"
      v-model:back-image="form.idCardBack"
      :disabled="isApproved"
    />

    <!-- 游戏认证 -->
    <GameCertSection
      :games="form.games"
      @add="addGameCert"
      @remove="removeGameCert"
      @update:screenshot="updateScreenshot"
    />

    <!-- 个人介绍 -->
    <IntroTextCard v-model="form.introduction" />

    <!-- 语音样本 -->
    <VoiceSampleCard
      :sample="form.voiceSample"
      :duration="form.voiceDuration"
      :recording="recording"
      :is-playing="isPlaying"
      @play="playVoice"
      @delete="deleteVoice"
      @record-start="startRecord"
      @record-end="stopRecord"
    />

    <view class="bottom-placeholder"></view>

    <template #footer>
      <!-- 底部操作栏 -->
      <view class="action-bar" v-if="!isApproved">
        <GlButton 
          type="primary" 
          block 
          round 
          size="large"
          :disabled="!canSubmit"
          :loading="submitting"
          @click="submitForm"
        >
          提交认证
        </GlButton>
      </view>

      <!-- 性别选择器 -->
      <uv-picker
        :show="showGenderPicker"
        :columns="[genderOptions.map(o => o.label)]"
        @confirm="(e: any) => { form.gender = genderOptions[e.indexs[0]]?.value || ''; showGenderPicker = false }"
        @cancel="showGenderPicker = false"
      ></uv-picker>

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
import CertStatusCard from '@/components/CertStatusCard/index.vue'
import CertificationBasicSection from '@/components/CertificationBasicSection/index.vue'
import IdCardUploader from '@/components/IdCardUploader/index.vue'
import GameCertSection from '@/components/GameCertSection/index.vue'
import IntroTextCard from '@/components/IntroTextCard/index.vue'
import VoiceSampleCard from '@/components/VoiceSampleCard/index.vue'
// Composables
import { usePlayerCertification } from '@/composables/usePlayerCertification'

const {
  submitting,
  certStatus,
  showGenderPicker,
  recording,
  isPlaying,
  form,
  genderOptions,
  isApproved,
  canSubmit,
  getGenderText,
  loadCertStatus,
  addGameCert,
  removeGameCert,
  updateScreenshot,
  startRecord,
  stopRecord,
  playVoice,
  deleteVoice,
  submitForm,
  goBack,
} = usePlayerCertification()

onMounted(loadCertStatus)
</script>

<style lang="scss" scoped>


:deep(.gl-card) {
  margin: 0 24rpx 20rpx;
}


.action-bar {
  padding: 20rpx 24rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: var(--color-bg-card);
  border-top: 1rpx solid var(--color-border);
}

.bottom-placeholder {
  height: 180rpx;
}
</style>
