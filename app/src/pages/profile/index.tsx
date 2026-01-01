import { View, Text, Image } from '@tarojs/components'
import { useState } from 'react'
import { PLACEHOLDER_AVATAR, ICONS } from '../../constants'
import './index.scss'

export default function Profile() {
    const [userInfo] = useState({
        name: 'Gamer123',
        id: '888888',
        avatar: PLACEHOLDER_AVATAR,
        balance: '128.50',
        coins: '500'
    })

    const [menuItems] = useState([
        { id: 1, title: 'My Wallet', icon: ICONS.CARD },
        { id: 2, title: 'Settings', icon: ICONS.SETTINGS },
        { id: 3, title: 'Help & Support', icon: ICONS.HELP },
        { id: 4, title: 'About Us', icon: ICONS.INFO },
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
                                <Image src={item.icon} style={{ width: '20px', height: '20px' }} />
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
