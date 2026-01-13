Component({
  options: {
    styleIsolation: 'apply-shared',
  },
  
  properties: {
    name: {
      type: String,
      value: '',
    },
    size: {
      type: String,
      value: 'medium', // small | medium | large | xlarge
    },
    color: {
      type: String,
      value: '',
    },
  },
  
  methods: {
    handleTap() {
      this.triggerEvent('tap')
    },
  },
})
