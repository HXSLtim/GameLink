import { View, Image, ScrollView, Text } from '@tarojs/components'
import Taro, { useLoad } from '@tarojs/taro'
import { Swiper, Grid } from '@nutui/nutui-react-taro'
import PlayerCard from '@/components/business/PlayerCard'
import './index.scss'

// Mock Data
const BANNERS = [
  'https://storage.360buyimg.com/jdc-article/NutUItaro34.jpg',
  'https://storage.360buyimg.com/jdc-article/NutUItaro2.jpg',
  'https://storage.360buyimg.com/jdc-article/welcomenutui.jpg',
]

const GAMES = [
  { name: '王者荣耀', icon: 'https://img12.360buyimg.com/imagetools/jfs/t1/147573/29/16034/8504/5fa46dfdE684d0b0b/314e30528574d5c9.png' },
  { name: '和平精英', icon: 'https://img12.360buyimg.com/imagetools/jfs/t1/128318/20/12985/6852/5f617497E3a8fcd22/a933224b4249a0e6.png' },
  { name: '英雄联盟', icon: 'https://img14.360buyimg.com/imagetools/jfs/t1/139366/30/13217/6334/5fa46dfdE4613c757/845184b2e153f581.png' },
  { name: '绝地求生', icon: 'https://img14.360buyimg.com/imagetools/jfs/t1/123136/1/14131/6179/5f617497E69b18349/30018d9cc93740e2.png' },
]

const PLAYERS = [
  { id: 1, nickname: '带妹上分', avatar: '', gameName: '王者荣耀', price: 20, rating: 5.0, tags: ['人皮话多', '技术流'] },
  { id: 2, nickname: '绝地刚枪王', avatar: '', gameName: '绝地求生', price: 35, rating: 4.8, tags: ['钢枪', '幽默'] },
  { id: 3, nickname: 'LOL陪玩', avatar: '', gameName: '英雄联盟', price: 25, rating: 4.9, tags: ['御姐音', '全能'] },
]

export default function Index() {
  useLoad(() => {
    console.log('Home Page loaded.')
  })

  return (
    <ScrollView className="home-page" scrollY>
      {/* Banner */}
      <View className="banner-section">
        <Swiper defaultValue={0} autoPlay indicator>
          {BANNERS.map((item, index) => (
            <Swiper.Item key={index}>
              <Image src={item} className="banner-img" mode="aspectFill" />
            </Swiper.Item>
          ))}
        </Swiper>
      </View>

      {/* Game Categories */}
      <View className="category-section">
        <Grid columns={4}>
          {GAMES.map((item, index) => (
            <Grid.Item key={index} text={item.name}>
              <Image src={item.icon} style={{ width: 40, height: 40, marginBottom: 8 }} />
            </Grid.Item>
          ))}
        </Grid>
      </View>

      {/* Recommended Players */}
      <View className="recommend-section">
        <Text className="section-title">推荐大神</Text>
        <View className="player-list">
          {PLAYERS.map(player => (
            <PlayerCard key={player.id} player={player} onClick={() => Taro.navigateTo({ url: `/pages/profile/index?id=${player.id}` })} />
          ))}
        </View>
      </View>
    </ScrollView>
  )
}

