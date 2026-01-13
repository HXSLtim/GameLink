Component({
    properties: {
        role: { type: String, value: 'user' },
        themeMode: { type: String, value: 'dark' },
    },
    methods: {
        onSwitchRole() {
            this.triggerEvent('switchRole')
        },
        onToggleTheme() {
            this.triggerEvent('toggleTheme')
        }
    }
})
