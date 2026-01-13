Component({
  options: {
    styleIsolation: 'apply-shared',
  },
  
  properties: {
    src: {
      type: String,
      value: '',
    },
    size: {
      type: String,
      value: 'medium', // small | medium | large | xlarge
    },
    placeholder: {
      type: String,
      value: '?',
    },
    badge: {
      type: String,
      value: '',
    },
    online: {
      type: Boolean,
      value: false,
    },
  },
  
  methods: {
    handleTap() {
      this.triggerEvent('tap')
    },
    
    handleError() {
      this.setData({ src: '' })
    },
  },
})
