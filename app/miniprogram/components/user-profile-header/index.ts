Component({
    properties: {
        role: {
            type: String,
            value: 'user',
        },
        statusBarHeight: {
            type: Number,
            value: 44,
        },
    },

    methods: {
        onEditProfile() {
            wx.navigateTo({ url: '/pages/profile/edit' })
        },
    },
})
