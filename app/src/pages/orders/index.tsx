import { View, Text, Image } from '@tarojs/components'
import { useState } from 'react'
import './index.scss'

export default function Orders() {
    const [activeTab, setActiveTab] = useState(0)
    const tabs = ['Pending', 'Ongoing', 'Completed']

    const [orders] = useState([
        { id: 1, companion: 'Alice', game: 'LoL', status: 0, price: '40', time: '2023-10-01 14:00' },
        { id: 2, companion: 'Bob', game: 'HoK', status: 1, price: '30', time: '2023-10-02 16:00' },
        { id: 3, companion: 'Charlie', game: 'Apex', status: 2, price: '50', time: '2023-09-30 20:00' },
    ])

    // Filter orders based on active tab (this is a simple mock, so I'll just show all or filter if status matched tab index)
    // For this prototype, let's just show relevant ones or all for simplicity, but let's do a simple filter.
    // Mapping status: 0 -> Pending, 1 -> Ongoing, 2 -> Completed
    const filteredOrders = orders.filter(o => o.status === activeTab)

    return (
        <View className='orders-page'>
            <View className='tabs'>
                {tabs.map((tab, index) => (
                    <View
                        key={index}
                        className={`tab-item ${activeTab === index ? 'active' : ''}`}
                        onClick={() => setActiveTab(index)}
                    >
                        <Text>{tab}</Text>
                        {activeTab === index && <View className='active-line' />}
                    </View>
                ))}
            </View>

            <View className='order-list'>
                {filteredOrders.length === 0 ? (
                    <View className='empty-state'>
                        <Text>No orders yet</Text>
                    </View>
                ) : (
                    filteredOrders.map(order => (
                        <View key={order.id} className='order-card'>
                            <View className='header'>
                                <Text className='time'>{order.time}</Text>
                                <Text className='status-text'>{tabs[order.status]}</Text>
                            </View>
                            <View className='content'>
                                <View className='info'>
                                    <Text className='companion-name'>Companion: {order.companion}</Text>
                                    <Text className='game'>Game: {order.game}</Text>
                                </View>
                                <Text className='price'>${order.price}</Text>
                            </View>
                            <View className='actions'>
                                <View className='btn'>Contact</View>
                                {order.status === 0 && <View className='btn primary'>Pay</View>}
                                {order.status === 1 && <View className='btn primary'>Complete</View>}
                                {order.status === 2 && <View className='btn'>Rate</View>}
                            </View>
                        </View>
                    ))
                )}
            </View>
        </View>
    )
}
