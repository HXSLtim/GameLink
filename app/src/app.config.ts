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
    color: '#000000',
    selectedColor: '#fa2c19',
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
