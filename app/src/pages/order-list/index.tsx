import { View, Text } from '@tarojs/components'
import './index.scss'

function OrderList() {
    return (
        <View className="order-list-page">
            <View className="order-tabs">
                <View className="tab-item tab-item--active">全部</View>
                <View className="tab-item">待支付</View>
                <View className="tab-item">进行中</View>
                <View className="tab-item">已完成</View>
            </View>
            <View className="empty-state">
                <Text className="empty-icon">📋</Text>
                <Text className="empty-text">暂无订单</Text>
                <View className="empty-btn">去下单</View>
            </View>
        </View>
    )
}

export default OrderList
