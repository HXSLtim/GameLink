Component({
  options: {
    multipleSlots: true,
    styleIsolation: 'apply-shared',
  },
  
  properties: {
    title: {
      type: String,
      value: '',
    },
    shadow: {
      type: Boolean,
      value: true,
    },
    border: {
      type: Boolean,
      value: false,
    },
    hoverable: {
      type: Boolean,
      value: false,
    },
  },
  
  methods: {
    handleTap() {
      this.triggerEvent('tap')
    },
  },
})
