import { Component } from 'react'
import Taro from '@tarojs/taro'
import { Tabbar } from '@nutui/nutui-react-taro'
import { Home, Message, Order, User } from '@nutui/icons-react-taro'
import '@nutui/nutui-react-taro/dist/style.css' // Import styles
import './index.scss'

// 使用设计系统色卡
const COLORS = {
    primary: '#FF4755',    // 活力红 - 主品牌色
    inactive: '#6B7280',   // gray-500 - 未选中
}

export default class CustomTabBar extends Component {
    state = {
        selected: 0
    }

    componentDidMount() {
        this.updateSelected()
    }

    updateSelected = () => {
        // Logic to update selected tab based on current page
        // Needs global state or manual update from pages
    }

    switchTab = (value: number) => {
        const urls = [
            '/pages/index/index',
            '/pages/message/index',
            '/pages/order-list/index',
            '/pages/profile/index'
        ]
        Taro.switchTab({
            url: urls[value]
        })
        this.setState({ selected: value })
    }

    render() {
        const { selected } = this.state
        return (
            <Tabbar
                fixed
                value={selected}
                onSwitch={this.switchTab}
                activeColor={COLORS.primary}
                inactiveColor={COLORS.inactive}
            >
                <Tabbar.Item title="首页" icon={<Home width={20} height={20} />} />
                <Tabbar.Item title="消息" icon={<Message width={20} height={20} />} />
                <Tabbar.Item title="订单" icon={<Order width={20} height={20} />} />
                <Tabbar.Item title="我的" icon={<User width={20} height={20} />} />
            </Tabbar>
        )
    }
}
