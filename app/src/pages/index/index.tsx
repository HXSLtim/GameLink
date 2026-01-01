import { View, Text, Image } from '@tarojs/components'
import { useState } from 'react'
import { AtIcon } from 'taro-ui'
import './index.scss'

export default function Index() {
  const [categories] = useState([
    { id: 1, name: 'LoL', icon: 'lightning-bolt' },
    { id: 2, name: 'HoK', icon: 'star' },
    { id: 3, name: 'Apex', icon: 'filter' },
    { id: 4, name: 'Valorant', icon: 'check-circle' },
    { id: 5, name: 'Genshin', icon: 'eye' },
    { id: 6, name: 'More', icon: 'menu' },
  ])

  const [companions] = useState([
    { id: 1, name: 'Alice', game: 'LoL', price: '20/hr', rating: '5.0', avatar: 'https://via.placeholder.com/100' },
    { id: 2, name: 'Bob', game: 'HoK', price: '15/hr', rating: '4.8', avatar: 'https://via.placeholder.com/100' },
    { id: 3, name: 'Charlie', game: 'Apex', price: '25/hr', rating: '4.9', avatar: 'https://via.placeholder.com/100' },
    { id: 4, name: 'Diana', game: 'Valorant', price: '30/hr', rating: '5.0', avatar: 'https://via.placeholder.com/100' },
  ])

  return (
    <View className='home-page'>
      {/* Search Bar Placeholder */}
      <View className='search-bar'>
        <AtIcon value='search' size='16' color='#999' />
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
              <AtIcon value={cat.icon} size='24' color='#ffffff' />
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
