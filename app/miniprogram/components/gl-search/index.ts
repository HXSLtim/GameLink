/**
 * 搜索框组件
 */
Component({
  options: {
    styleIsolation: 'apply-shared',
  },

  properties: {
    placeholder: {
      type: String,
      value: 'Search...',
    },
    value: {
      type: String,
      value: '',
    },
    // 是否显示右侧图标
    showIcon: {
      type: Boolean,
      value: true,
    },
  },

  methods: {
    onInput(e: WechatMiniprogram.Input) {
      this.triggerEvent('input', { value: e.detail.value })
    },

    onFocus() {
      this.triggerEvent('focus')
    },

    onBlur() {
      this.triggerEvent('blur')
    },

    onConfirm(e: WechatMiniprogram.Input) {
      this.triggerEvent('search', { value: e.detail.value })
    },

    onClear() {
      this.triggerEvent('input', { value: '' })
      this.triggerEvent('clear')
    },
  },
})
