import { View } from '@tarojs/components'
import Taro from '@tarojs/taro'
import { Cell, Grid, Button } from '@nutui/nutui-react-taro'
import { Order, Cart, Wallet, Service, Notice } from '@nutui/icons-react-taro'
import UserCard from '@/components/business/UserCard'
import { useUserStore } from '@/stores/user'
import './index.scss'

const Profile = () => {
    const { logout } = useUserStore()

    const handleLogout = () => {
        logout()
        Taro.reLaunch({ url: '/pages/login/index' })
    }

    return (
        <View className="profile-page">
            <UserCard />

            {/* Wallet Card */}
            <View className="wallet-card">
                <View className="wallet-header">
                    <View>我的钱包</View>
                    <View className="balance">¥ 0.00</View>
                </View>
                <Grid columns={3}>
                    <Grid.Item text="充值"><Wallet /></Grid.Item>
                    <Grid.Item text="提现"><Cart /></Grid.Item>
                    <Grid.Item text="账单"><Order /></Grid.Item>
                </Grid>
            </View>

            {/* Menu List */}
            <View className="menu-list">
                <Cell title="我的订单" extra={<Order />} onClick={() => Taro.switchTab({ url: '/pages/order-list/index' })} />
                <Cell title="在线客服" extra={<Service />} />
                <Cell title="消息中心" extra={<Notice />} onClick={() => Taro.switchTab({ url: '/pages/message/index' })} />
            </View>

            <View className="logout-btn">
                <Button block type="danger" onClick={handleLogout}>退出登录</Button>
            </View>
        </View>
    )
}

export default Profile
