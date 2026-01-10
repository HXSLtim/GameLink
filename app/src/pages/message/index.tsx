import { View, Text } from '@tarojs/components'
import './index.scss'

function Message() {
    return (
        <View className="message-page">
            <View className="empty-state">
                <Text className="empty-icon">💬</Text>
                <Text className="empty-text">暂无消息</Text>
            </View>
        </View>
    )
}

export default Message
