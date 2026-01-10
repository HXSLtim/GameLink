export default defineAppConfig({
  pages: [
    'pages/index/index',
    'pages/message/index',
    'pages/order-list/index',
    'pages/profile/index',
    'pages/login/index'
  ],
  tabBar: {
    custom: true,
    color: '#6B7280',           // $gray-500 - 未选中文字
    selectedColor: '#FF4755',   // $primary-color - 活力红
    backgroundColor: '#ffffff',
    list: [
      {
        pagePath: 'pages/index/index',
        text: '首页'
      },
      {
        pagePath: 'pages/message/index',
        text: '消息'
      },
      {
        pagePath: 'pages/order-list/index',
        text: '订单'
      },
      {
        pagePath: 'pages/profile/index',
        text: '我的'
      }
    ]
  },
  window: {
    backgroundTextStyle: 'light',
    navigationBarBackgroundColor: '#fff',
    navigationBarTitleText: 'GameLink',
    navigationBarTextStyle: 'black'
  }
})
