import { View, Text, Image } from '@tarojs/components'
import { useState } from 'react'
import { PLACEHOLDER_AVATAR } from '../../constants'
import './index.scss'

export default function Messages() {
    const [messages] = useState([
        { id: 1, name: 'Alice', message: 'Hi, are you clear for tonight?', time: '14:30', avatar: PLACEHOLDER_AVATAR, unread: 2 },
        { id: 2, name: 'System', message: 'Your order has been completed.', time: 'Yesterday', avatar: PLACEHOLDER_AVATAR, unread: 0 },
        { id: 3, name: 'Bob', message: 'GG!', time: 'Yesterday', avatar: PLACEHOLDER_AVATAR, unread: 0 },
    ])

    return (
        <View className='messages-page'>
            {messages.map(msg => (
                <View key={msg.id} className='message-item'>
                    <View className='avatar-container'>
                        <Image className='avatar' src={msg.avatar} mode='aspectFill' />
                        {msg.unread > 0 && <View className='badge'>{msg.unread}</View>}
                    </View>
                    <View className='content'>
                        <View className='header'>
                            <Text className='name'>{msg.name}</Text>
                            <Text className='time'>{msg.time}</Text>
                        </View>
                        <Text className='message-preview'>{msg.message}</Text>
                    </View>
                </View>
            ))}
        </View>
    )
}
