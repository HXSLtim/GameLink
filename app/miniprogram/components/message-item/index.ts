/**
 * 消息项组件
 */
Component({
  options: {
    styleIsolation: 'apply-shared',
  },

  properties: {
    // 消息数据
    message: {
      type: Object,
      value: {},
    },
  },

  methods: {
    onTap() {
      this.triggerEvent('tap', { message: this.data.message })
    },
  },
})
