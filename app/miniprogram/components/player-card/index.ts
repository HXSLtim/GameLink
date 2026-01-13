Component({
  properties: {
    player: {
      type: Object,
      value: {
        id: '',
        name: '',
        avatar: '',
        imageUrl: '',
        tags: [],
        price: 0,
        rating: 0,
        orders: 0,
        game: '',
        audioDuration: 0,
        isOnline: false,
      },
    },
  },

  methods: {
    handleTap() {
      this.triggerEvent('tap', { player: this.data.player })
    },
  },
})
