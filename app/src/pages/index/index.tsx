import { View, Text, Image } from '@tarojs/components'
import { useState } from 'react'
import { PLACEHOLDER_AVATAR, ICONS } from '../../constants'
import './index.scss'

export default function Index() {
  const [categories] = useState([
    { id: 1, name: 'LoL', icon: ICONS.LIGHTNING },
    { id: 2, name: 'HoK', icon: ICONS.STAR },
    { id: 3, name: 'Apex', icon: ICONS.FILTER },
    { id: 4, name: 'Valorant', icon: ICONS.CHECK },
    { id: 5, name: 'Genshin', icon: ICONS.EYE },
    { id: 6, name: 'More', icon: ICONS.MENU },
  ])

  const [companions] = useState([
    { id: 1, name: 'Alice', game: 'LoL', price: '20/hr', rating: '5.0', avatar: PLACEHOLDER_AVATAR },
    { id: 2, name: 'Bob', game: 'HoK', price: '15/hr', rating: '4.8', avatar: PLACEHOLDER_AVATAR },
    { id: 3, name: 'Charlie', game: 'Apex', price: '25/hr', rating: '4.9', avatar: PLACEHOLDER_AVATAR },
    { id: 4, name: 'Diana', game: 'Valorant', price: '30/hr', rating: '5.0', avatar: PLACEHOLDER_AVATAR },
  ])

  return (
    <View className='home-page'>
      {/* Search Bar Placeholder */}
      <View className='search-bar'>
        <Image src={ICONS.SEARCH} style={{ width: '16px', height: '16px' }} />
        <Text className='placeholder' style={{ marginLeft: '8px' }}>Search for games or companions...</Text>
      </View>

      {/* Banner */}
      <View className='banner'>
        <View className='banner-item'>
          <Text className='banner-text'>Find your perfect gaming duo!</Text>
        </View>
      </View>

      {/* Categories */}
      <View className='section-title'>Popular Games</View>
      <View className='categories-grid'>
        {categories.map(cat => (
          <View key={cat.id} className='category-item'>
            <View className='icon-wrapper'>
              <Image src={cat.icon} style={{ width: '24px', height: '24px' }} />
            </View>
            <Text className='name'>{cat.name}</Text>
          </View>
        ))}
      </View>

      {/* Featured Companions */}
      <View className='section-title'>Featured Companions</View>
      <View className='companions-list'>
        {companions.map(comp => (
          <View key={comp.id} className='companion-card'>
            <Image className='avatar' src={comp.avatar} mode='aspectFill' />
            <View className='info'>
              <View className='header'>
                <Text className='name'>{comp.name}</Text>
                <Text className='price'>{comp.price}</Text>
              </View>
              <Text className='game'>{comp.game}</Text>
              <Text className='rating'>⭐ {comp.rating}</Text>
            </View>
          </View>
        ))}
      </View>
    </View>
  )
}
