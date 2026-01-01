export default defineAppConfig({
  pages: [
    'pages/index/index',
    'pages/orders/index',
    'pages/messages/index',
    'pages/profile/index'
  ],
  window: {
    backgroundTextStyle: 'light',
    navigationBarBackgroundColor: '#fff',
    navigationBarTitleText: 'GameLink',
    navigationBarTextStyle: 'black'
  },
  tabBar: {
    color: '#999',
    selectedColor: '#ff381a',
    backgroundColor: '#fff',
    list: [
      {
        pagePath: 'pages/index/index',
        text: 'Home'
      },
      {
        pagePath: 'pages/orders/index',
        text: 'Orders'
      },
      {
        pagePath: 'pages/messages/index',
        text: 'Messages'
      },
      {
        pagePath: 'pages/profile/index',
        text: 'Profile'
      }
    ]
  }
})
