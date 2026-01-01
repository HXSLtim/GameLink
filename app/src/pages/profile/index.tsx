import { View, Text, Image } from '@tarojs/components'
import { useState } from 'react'
import { AtIcon } from 'taro-ui'
import './index.scss'

export default function Profile() {
    const [userInfo] = useState({
        name: 'Gamer123',
        id: '888888',
        avatar: 'https://via.placeholder.com/100',
        balance: '128.50',
        coins: '500'
    })

    const [menuItems] = useState([
        { id: 1, title: 'My Wallet', icon: 'credit-card' },
        { id: 2, title: 'Settings', icon: 'settings' },
        { id: 3, title: 'Help & Support', icon: 'help' },
        { id: 4, title: 'About Us', icon: 'alert-circle' },
    ])

    return (
        <View className='profile-page'>
            {/* Header */}
            <View className='header'>
                <Image className='avatar' src={userInfo.avatar} mode='aspectFill' />
                <View className='info'>
                    <Text className='name'>{userInfo.name}</Text>
                    <Text className='uid'>ID: {userInfo.id}</Text>
                </View>
                <View className='edit-btn'>Edit</View>
            </View>

            {/* Stats/Wallet */}
            <View className='stats-card'>
                <View className='stat-item'>
                    <Text className='value'>${userInfo.balance}</Text>
                    <Text className='label'>Balance</Text>
                </View>
                <View className='divider' />
                <View className='stat-item'>
                    <Text className='value'>{userInfo.coins}</Text>
                    <Text className='label'>Coins</Text>
                </View>
            </View>

            {/* Menu */}
            <View className='menu-list'>
                {menuItems.map(item => (
                    <View key={item.id} className='menu-item'>
                        <View className='left'>
                            <View className='icon'>
                                <AtIcon value={item.icon} size='20' color='#ffffff' />
                            </View>
                            <Text className='title'>{item.title}</Text>
                        </View>
                        <Text className='arrow'>{'>'}</Text>
                    </View>
                ))}
            </View>

            <View className='logout-btn'>
                <Text>Log Out</Text>
            </View>
        </View>
    )
}
