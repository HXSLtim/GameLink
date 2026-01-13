Component({
    methods: {
        onMenuTap(e: WechatMiniprogram.TouchEvent) {
            const { type } = e.currentTarget.dataset
            this.triggerEvent('menuTap', { type })
        }
    }
})
