Component({
  options: {
    styleIsolation: 'apply-shared',
  },
  
  properties: {
    type: {
      type: String,
      value: 'default', // default | primary | success | warning | error | info
    },
    size: {
      type: String,
      value: 'medium', // small | medium
    },
  },
})
