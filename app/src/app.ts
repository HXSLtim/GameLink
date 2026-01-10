import { PropsWithChildren } from 'react'
import { useLaunch } from '@tarojs/taro'
import { initializeAuthState } from './stores/user'

import '@nutui/nutui-react-taro/dist/style.css'
import './app.scss'

function App({ children }: PropsWithChildren<any>) {
  useLaunch(() => {
    console.log('App launched.')

    // 初始化认证状态（从本地存储恢复）
    initializeAuthState()
  })

  // children 是将要会渲染的页面
  return children
}
  


export default App
