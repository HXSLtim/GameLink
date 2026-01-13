/**
 * 空状态组件
 */
Component({
  options: {
    styleIsolation: 'apply-shared',
  },

  properties: {
    icon: {
      type: String,
      value: '/assets/icons/message.svg',
    },
    title: {
      type: String,
      value: 'No data',
    },
    description: {
      type: String,
      value: '',
    },
    // 是否显示操作按钮
    showAction: {
      type: Boolean,
      value: false,
    },
    actionText: {
      type: String,
      value: 'Action',
    },
  },

  methods: {
    onAction() {
      this.triggerEvent('action')
    },
  },
})
