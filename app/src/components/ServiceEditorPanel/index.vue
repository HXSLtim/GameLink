<template>
  <uv-popup :show="show" mode="bottom" round="24" @close="$emit('close')">
    <view class="editor-panel">
      <view class="editor-header">
        <text class="editor-title">{{ title }}</text>
        <uv-icon name="close" size="20" @click="$emit('close')"></uv-icon>
      </view>

      <scroll-view class="editor-content" scroll-y>
        <FormItem label="游戏" :display-value="form.gameName" required @click="$emit('pick-game')" />
        <FormItem label="服务类型" :display-value="form.serviceName" required @click="$emit('pick-service')" />
        <FormItem label="段位" :display-value="form.rankName" required @click="$emit('pick-rank')" />

        <view class="form-row">
          <text class="form-label">价格 <text class="required">*</text></text>
          <view class="price-input-wrap">
            <text class="currency">¥</text>
            <GlInput
              v-model.number="form.price"
              class="price-input"
              type="digit"
              size="small"
              variant="plain"
              placeholder="0.00"
            />
            <text class="unit">/{{ form.unit }}</text>
          </view>
        </view>

        <view class="form-row">
          <text class="form-label">服务介绍</text>
          <GlInput
            v-model="form.description"
            class="desc-input"
            type="textarea"
            size="small"
            placeholder="介绍您的服务特点..."
            :maxlength="200"
          />
        </view>
      </scroll-view>

      <view class="editor-footer">
        <GlButton type="primary" block round :loading="saving" @click="$emit('save')">
          {{ editing ? '保存修改' : '添加服务' }}
        </GlButton>
      </view>
    </view>
  </uv-popup>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlButton from '@/components/gl/Button/index.vue'
import GlInput from '@/components/gl/Input/index.vue'
import FormItem from '@/components/FormItem/index.vue'
import type { PlayerServiceForm } from '@/types/player'

interface Props {
  show: boolean
  editing: boolean
  saving: boolean
  form: PlayerServiceForm
}

const props = defineProps<Props>()

defineEmits<{
  close: []
  save: []
  'pick-game': []
  'pick-service': []
  'pick-rank': []
}>()

const title = computed(() => (props.editing ? '编辑服务' : '添加服务'))
</script>

<style lang="scss" scoped>
.editor-panel {
  padding: var(--spacing-md);
  padding-bottom: calc(var(--spacing-md) + env(safe-area-inset-bottom));
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-md);
}

.editor-title {
  font-size: var(--font-lg);
  font-weight: 700;
  color: var(--color-text);
}

.editor-content {
  flex: 1;
  overflow-y: auto;
}

.form-row {
  margin-bottom: var(--spacing-md);
}

.form-label {
  display: block;
  font-size: var(--font-md);
  color: var(--color-text);
  margin-bottom: var(--spacing-xs);

  .required {
    color: var(--color-error);
  }
}

.price-input-wrap {
  display: flex;
  align-items: center;
  padding: var(--spacing-sm);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-md);
  border: 1rpx solid var(--color-border);
}

.currency {
  font-size: var(--font-base);
  font-weight: 600;
  color: var(--color-text);
  margin-right: var(--spacing-xs);
}

.price-input {
  flex: 1;
  min-width: 0;
  
  :deep(.gl-input__field) {
    font-size: var(--font-base);
    font-weight: 600;
    color: var(--color-text);
  }
}

.unit {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.desc-input {
  :deep(.gl-input__textarea) {
    min-height: 160rpx;
  }
}

.editor-footer {
  padding-top: var(--spacing-md);
}
</style>
