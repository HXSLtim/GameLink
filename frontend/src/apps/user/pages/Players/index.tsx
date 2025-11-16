/**
 * 陪玩师浏览页面
 */
import React, { useState } from 'react';
import { DiscordLayout } from '@/shared/components/DiscordLayout';
import { Button } from '@/shared/components/Button';
import { PlayerCard, type PlayerCardProps } from '../../components/PlayerCard';
import { usePlayerList } from '../../hooks/usePlayerList';
import type { Player } from '@/api/player';
import './Players.less';

/**
 * 筛选条件接口
 */
interface FilterState {
  game: string;
  status: string;
  priceRange: string;
  sortBy: string;
}

const PlayersPage: React.FC = () => {
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState(false);
  const [filters, setFilters] = useState<FilterState>({
    game: 'all',
    status: 'all',
    priceRange: 'all',
    sortBy: 'rating',
  });
  const [searchQuery, setSearchQuery] = useState('');

  // TODO: 后续将替换为真实API调用
  const mockPlayers: Player[] = [
    {
      id: 1,
      avatar: 'https://via.placeholder.com/64',
      name: '王者大神',
      game: '王者荣耀',
      rating: 4.8,
      reviewCount: 256,
      pricePerHour: 50,
      status: 'online',
      bio: '国服最强打野，带你上王者！擅长各种英雄，游戏经验丰富。',
      tags: ['打野', '上分快', '技术好'],
      isBusy: false,
    },
    {
      id: 2,
      avatar: 'https://via.placeholder.com/64',
      name: '温柔小姐姐',
      game: '王者荣耀',
      rating: 4.9,
      reviewCount: 512,
      pricePerHour: 80,
      status: 'online',
      bio: '声音甜美，游戏技术也不错哦~陪你一起开黑，快乐上分！',
      tags: ['声音甜', '辅助位', '陪聊'],
      isBusy: false,
    },
    {
      id: 3,
      avatar: 'https://via.placeholder.com/64',
      name: 'LOL王者',
      game: '英雄联盟',
      rating: 4.7,
      reviewCount: 189,
      pricePerHour: 60,
      status: 'idle',
      bio: '钻石选手，擅长中单和ADC，带你轻松上分。',
      tags: ['中单', 'ADC', '钻石'],
      isBusy: false,
    },
    {
      id: 4,
      avatar: 'https://via.placeholder.com/64',
      name: '吃鸡高手',
      game: '和平精英',
      rating: 4.6,
      reviewCount: 145,
      pricePerHour: 45,
      status: 'online',
      bio: '枪法精准，意识好，带你吃鸡！',
      tags: ['枪法准', '意识好', '吃鸡率高'],
      isBusy: false,
    },
    {
      id: 5,
      avatar: 'https://via.placeholder.com/64',
      name: '专业陪练',
      game: '原神',
      rating: 4.5,
      reviewCount: 98,
      pricePerHour: 40,
      status: 'busy',
      bio: '原神深渊满星，帮你通关各种副本。',
      tags: ['深渊', '副本', '原神'],
      isBusy: true,
    },
    {
      id: 6,
      avatar: 'https://via.placeholder.com/64',
      name: 'CS高手',
      game: 'CS:GO',
      rating: 4.9,
      reviewCount: 276,
      pricePerHour: 70,
      status: 'offline',
      bio: '职业选手退役，擅长狙击和突破，带你冲分！',
      tags: ['职业', '狙击', 'FPS'],
      isBusy: false,
    },
  ];

  useEffect(() => {
    // 模拟数据加载
    setLoading(true);
    setTimeout(() => {
      setPlayers(mockPlayers);
      setLoading(false);
    }, 500);
  }, []);

  const handleFilterChange = (key: keyof FilterState, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
  };

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value);
  };

  const filteredPlayers = players.filter((player) => {
    // 搜索过滤
    if (
      searchQuery &&
      !player.name.toLowerCase().includes(searchQuery.toLowerCase()) &&
      !player.game.toLowerCase().includes(searchQuery.toLowerCase())
    ) {
      return false;
    }

    // 游戏过滤
    if (filters.game !== 'all' && player.game !== filters.game) {
      return false;
    }

    // 状态过滤
    if (filters.status !== 'all' && player.status !== filters.status) {
      return false;
    }

    // 价格过滤
    if (filters.priceRange !== 'all') {
      const price = player.pricePerHour;
      switch (filters.priceRange) {
        case 'low':
          if (price >= 50) return false;
          break;
        case 'medium':
          if (price < 50 || price >= 80) return false;
          break;
        case 'high':
          if (price < 80) return false;
          break;
      }
    }

    return true;
  });

  // 排序
  const sortedPlayers = [...filteredPlayers].sort((a, b) => {
    switch (filters.sortBy) {
      case 'rating':
        return b.rating - a.rating;
      case 'price-low':
        return a.pricePerHour - b.pricePerHour;
      case 'price-high':
        return b.pricePerHour - a.pricePerHour;
      case 'reviews':
        return b.reviewCount - a.reviewCount;
      default:
        return 0;
    }
  });

  return (
    <DiscordLayout>
      <DiscordLayout.Main>
        <DiscordLayout.Header>
          <div className="players-header">
            <h1 className="players-header__title">发现陪玩师</h1>
            <div className="players-header__actions">
              <input
                type="text"
                className="players-header__search"
                placeholder="搜索陪玩师或游戏..."
                value={searchQuery}
                onChange={handleSearch}
              />
            </div>
          </div>
        </DiscordLayout.Header>

        <DiscordLayout.Content>
          <div className="players-page">
            {/* 筛选栏 */}
            <div className="players-filters">
              <div className="players-filters__group">
                <label className="players-filters__label">游戏:</label>
                <select
                  className="players-filters__select"
                  value={filters.game}
                  onChange={(e) => handleFilterChange('game', e.target.value)}
                >
                  <option value="all">全部游戏</option>
                  <option value="王者荣耀">王者荣耀</option>
                  <option value="英雄联盟">英雄联盟</option>
                  <option value="和平精英">和平精英</option>
                  <option value="原神">原神</option>
                  <option value="CS:GO">CS:GO</option>
                </select>
              </div>

              <div className="players-filters__group">
                <label className="players-filters__label">状态:</label>
                <select
                  className="players-filters__select"
                  value={filters.status}
                  onChange={(e) => handleFilterChange('status', e.target.value)}
                >
                  <option value="all">全部状态</option>
                  <option value="online">在线</option>
                  <option value="idle">离开</option>
                  <option value="busy">忙碌</option>
                  <option value="offline">离线</option>
                </select>
              </div>

              <div className="players-filters__group">
                <label className="players-filters__label">价格:</label>
                <select
                  className="players-filters__select"
                  value={filters.priceRange}
                  onChange={(e) => handleFilterChange('priceRange', e.target.value)}
                >
                  <option value="all">全部价格</option>
                  <option value="low">&lt; ¥50/小时</option>
                  <option value="medium">¥50-80/小时</option>
                  <option value="high">&gt; ¥80/小时</option>
                </select>
              </div>

              <div className="players-filters__group">
                <label className="players-filters__label">排序:</label>
                <select
                  className="players-filters__select"
                  value={filters.sortBy}
                  onChange={(e) => handleFilterChange('sortBy', e.target.value)}
                >
                  <option value="rating">评分最高</option>
                  <option value="reviews">评价最多</option>
                  <option value="price-low">价格最低</option>
                  <option value="price-high">价格最高</option>
                </select>
              </div>
            </div>

            {/* 陪玩师列表 */}
            <div className="players-list">
              {loading ? (
                <div className="players-loading">加载中...</div>
              ) : sortedPlayers.length > 0 ? (
                sortedPlayers.map((player) => (
                  <PlayerCard key={player.id} {...player} />
                ))
              ) : (
                <div className="players-empty">
                  <p>没有找到符合条件的陪玩师</p>
                  <Button type="secondary" onClick={() => setFilters({
                    game: 'all',
                    status: 'all',
                    priceRange: 'all',
                    sortBy: 'rating',
                  })}>
                    清除筛选
                  </Button>
                </div>
              )}
            </div>
          </div>
        </DiscordLayout.Content>
      </DiscordLayout.Main>
    </DiscordLayout>
  );
};

export default PlayersPage;
