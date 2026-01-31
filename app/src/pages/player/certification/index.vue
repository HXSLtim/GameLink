<template>
  <view class="certification-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="陪玩认证" @back="goBack" />

    <!-- 认证状态卡片 -->
    <CertStatusCard :status="certStatus" />

    <scroll-view class="content-scroll" scroll-y>
      <!-- 基本信息 -->
      <FormSection title="基本信息">
        <template #extra>
          <text class="required-tip">* 为必填项</text>
        </template>
        
        <FormItem
          v-model="form.realName"
          label="真实姓名"
          type="input"
          placeholder="请输入真实姓名"
          required
        />
        
        <FormItem
          v-model="form.idNumber"
          label="身份证号"
          type="input"
          placeholder="请输入身份证号"
          required
        />
        
        <FormItem
          label="性别"
          :display-value="getGenderText(form.gender)"
          required
          @click="showGenderPicker = true"
        />
      </FormSection>

      <!-- 身份证照片 -->
      <IdCardUploader
        v-model:front-image="form.idCardFront"
        v-model:back-image="form.idCardBack"
        :disabled="isApproved"
      />

      <!-- 游戏认证 -->
      <GlCard title="游戏认证" :shadow="false" bordered>
        <template #extra>
          <GlButton type="primary" size="mini" plain @click="addGameCert">+ 添加</GlButton>
        </template>
        
        <GameCertItem
          v-for="(game, index) in form.games"
          :key="index"
          :game="game"
          @remove="removeGameCert(index)"
          @select-game="() => {}"
          @select-rank="() => {}"
          @update:screenshot="(url) => updateScreenshot(index, url)"
        />
        
        <GlEmpty v-if="form.games.length === 0" title="请添加至少一个游戏认证" compact />
      </GlCard>

      <!-- 个人介绍 -->
      <GlCard title="个人介绍" :shadow="false" bordered>
        <textarea 
          v-model="form.introduction"
          class="intro-textarea"
          placeholder="介绍一下自己的陪玩特色和优势吧~"
          :maxlength="500"
        />
        <text class="char-count">{{ form.introduction?.length || 0 }}/500</text>
      </GlCard>

      <!-- 语音样本 -->
      <GlCard title="语音样本（选填）" :shadow="false" bordered>
        <view class="voice-upload">
          <view v-if="form.voiceSample" class="voice-item">
            <view class="voice-play" @tap="playVoice">
              <uv-icon :name="isPlaying ? 'pause-circle' : 'play-circle'" size="32" color="var(--color-primary)"></uv-icon>
            </view>
            <text class="voice-duration">{{ form.voiceDuration }}s</text>
            <view class="voice-delete" @tap="deleteVoice">
              <uv-icon name="close" size="16" color="var(--color-text-secondary)"></uv-icon>
            </view>
          </view>
          
          <view 
            v-else 
            class="voice-record"
            @touchstart="startRecord"
            @touchend="stopRecord"
          >
            <uv-icon name="mic" size="32" :color="recording ? 'var(--color-primary)' : 'var(--color-text-secondary)'"></uv-icon>
            <text class="record-text">{{ recording ? '松开结束' : '按住录制语音' }}</text>
          </view>
        </view>
      </GlCard>

      <view class="bottom-placeholder"></view>
    </scroll-view>

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

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import FormSection from '@/components/FormSection/index.vue'
import FormItem from '@/components/FormItem/index.vue'
import GlCard from '@/components/gl/Card/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
// Business 组件
import CertStatusCard from '@/components/CertStatusCard/index.vue'
import IdCardUploader from '@/components/IdCardUploader/index.vue'
import GameCertItem from '@/components/GameCertItem/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
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
.certification-page {
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
  padding-bottom: 160rpx;
}

.required-tip {
  font-size: 24rpx;
  color: var(--color-error);
}

:deep(.gl-card) {
  margin: 0 24rpx 20rpx;
}

.intro-textarea {
  width: 100%;
  height: 200rpx;
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

.voice-upload {
  padding: 20rpx 0;
}

.voice-item {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 20rpx;
  background: var(--color-bg-secondary);
  border-radius: 12rpx;
}

.voice-duration {
  flex: 1;
  font-size: 28rpx;
  color: var(--color-text);
}

.voice-record {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
  padding: 40rpx;
  background: var(--color-bg-secondary);
  border-radius: 16rpx;
  border: 2rpx dashed var(--color-border);
}

.record-text {
  font-size: 26rpx;
  color: var(--color-text-secondary);
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
