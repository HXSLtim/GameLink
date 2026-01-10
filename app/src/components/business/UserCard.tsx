import { View, Image, Text } from '@tarojs/components'
// import { Button } from '@nutui/nutui-react-taro'
import { Setting } from '@nutui/icons-react-taro'
import { useUserStore } from '@/stores/user'
import './UserCard.scss'

interface UserCardProps {
    onSettingClick?: () => void
}

const UserCard = ({ onSettingClick }: UserCardProps) => {
    const { currentUser, isPlayer, switchRole } = useUserStore()

    const handleRoleSwitch = () => {
        switchRole(isPlayer ? 'user' : 'player')
    }

    return (
        <View className="user-card-container">
            <View className="user-info">
                <Image
                    className="avatar"
                    src={currentUser?.avatarUrl || 'https://img12.360buyimg.com/imagetools/jfs/t1/143702/31/16654/116794/5fc6f541Edebf8a57/4138097748889987.png'}
                    mode="aspectFill"
                />
                <View className="info-text">
                    <Text className="nickname">{currentUser?.name || '未登录用户'}</Text>
                    <Text className="user-id">ID: {currentUser?.id || '000000'}</Text>
                </View>
                <View className="expert-tag" onClick={handleRoleSwitch}>
                    <Text>{isPlayer ? '陪玩版' : '用户版'}</Text>
                </View>
                <Setting className="setting-icon" onClick={onSettingClick} />
            </View>

            <View className="stats-row">
                <View className="stat-item">
                    <Text className="num">0</Text>
                    <Text className="label">关注</Text>
                </View>
                <View className="stat-item">
                    <Text className="num">0</Text>
                    <Text className="label">粉丝</Text>
                </View>
                <View className="stat-item">
                    <Text className="num">0</Text>
                    <Text className="label">访客</Text>
                </View>
            </View>
        </View>
    )
}

export default UserCard
