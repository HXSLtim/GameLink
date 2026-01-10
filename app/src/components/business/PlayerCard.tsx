import { View, Image, Text } from '@tarojs/components'
import { Tag, Price } from '@nutui/nutui-react-taro'
import { StarFill } from '@nutui/icons-react-taro'
import './PlayerCard.scss'

interface PlayerCardProps {
    player: {
        id: number
        nickname: string
        avatar: string
        gameName: string
        price: number
        rating: number
        tags: string[]
    }
    onClick?: () => void
}

const PlayerCard = ({ player, onClick }: PlayerCardProps) => {
    return (
        <View className="player-card" onClick={onClick}>
            <Image className="avatar" src={player.avatar || 'https://img12.360buyimg.com/imagetools/jfs/t1/196130/38/8105/14329/60ec7465E052ecc2f/6487e076644eb6e3.jpg'} mode="aspectFill" />
            <View className="info">
                <View className="header">
                    <Text className="nickname">{player.nickname}</Text>
                    <View className="rating">
                        <StarFill color="#fa2c19" size={12} />
                        <Text className="score">{player.rating}</Text>
                    </View>
                </View>
                <View className="game-tags">
                    <Text className="game">{player.gameName}</Text>
                    {player.tags.map(tag => (
                        <Tag key={tag} type="primary" plain style={{ marginLeft: 4 }}>{tag}</Tag>
                    ))}
                </View>
                <View className="price-row">
                    <Price price={player.price} symbol="¥" position="before" size="normal" />
                    <Text className="unit">/局</Text>
                </View>
            </View>
        </View>
    )
}

export default PlayerCard
