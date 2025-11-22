/**
 * 陪玩师列表页面 - Discord风格
 * 使用新的DiscordLayout和PlayerCard组件
 */

import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { DiscordLayout } from '@/shared/components/DiscordLayout';
import { PlayerCard, PlayerCardProps } from '@/apps/user/components/PlayerCard';
import styles from './PlayersList.module.less';

/**
 * 筛选条件接口
 */
interface FilterState {
  gameId: number | null;
  minRating: number;
  maxPrice: number;
  onlineOnly: boolean;
  sortBy: 'rating' | 'price-low' | 'price-high' | 'reviews';
}

/**
 * 游戏选项
 */
const GAME_OPTIONS = [
  { id: 0, name: '全部游戏' },
  { id: 1, name: '王者荣耀' },
  { id: 2, name: '英雄联盟' },
  { id: 3, name: '和平精英' },
  { id: 4, name: '原神' },
  { id: 5, name: 'CS:GO' },
];

/**
 * 陪玩师列表页面组件
 */
export const PlayersList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [filters, setFilters] = useState<FilterState>({
    gameId: null,
    minRating: 0,
    maxPrice: 999,
    onlineOnly: false,
    sortBy: 'rating',
  });

  // Mock数据（后续替换为真实API调用）
  const [players] = useState<PlayerCardProps[]>([
    {
      id: 1,
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player1',
      nickname: '王者大神',
      gameName: '王者荣耀',
      rating: 4.8,
      reviewCount: 256,
      price: 50,
      isOnline: true,
      tags: ['打野', '上分快', '技术好'],
      signature: '国服最强打野，带你上王者！擅长各种英雄，游戏经验丰富。',
    },
    {
      id: 2,
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player2',
      nickname: '温柔小姐姐',
      gameName: '王者荣耀',
      rating: 4.9,
      reviewCount: 512,
      price: 80,
      isOnline: true,
      tags: ['声音甜', '辅助位', '陪聊'],
      signature: '声音甜美，游戏技术也不错哦~陪你一起开黑，快乐上分！',
    },
    {
      id: 3,
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player3',
      nickname: 'LOL王者',
      gameName: '英雄联盟',
      rating: 4.7,
      reviewCount: 189,
      price: 60,
      isOnline: false,
      tags: ['中单', 'ADC', '钻石'],
      signature: '钻石选手，擅长中单和ADC，带你轻松上分。',
    },
    {
      id: 4,
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player4',
      nickname: '吃鸡高手',
      gameName: '和平精英',
      rating: 4.6,
      reviewCount: 145,
      price: 45,
      isOnline: true,
      tags: ['枪法准', '意识好', '吃鸡率高'],
      signature: '枪法精准，意识好，带你吃鸡！',
    },
    {
      id: 5,
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player5',
      nickname: '专业陪练',
      gameName: '原神',
      rating: 4.5,
      reviewCount: 98,
      price: 40,
      isOnline: false,
      tags: ['深渊', '副本', '原神'],
      signature: '原神深渊满星，帮你通关各种副本。',
    },
    {
      id: 6,
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player6',
      nickname: 'CS高手',
      gameName: 'CS:GO',
      rating: 4.9,
      reviewCount: 276,
      price: 70,
      isOnline: false,
      tags: ['职业', '狙击', 'FPS'],
      signature: '职业选手退役，擅长狙击和突破，带你冲分！',
    },
    {
      id: 7,
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player7',
      nickname: '元气少女',
      gameName: '王者荣耀',
      rating: 4.8,
      reviewCount: 342,
      price: 65,
      isOnline: true,
      tags: ['甜美', '射手', '陪聊'],
      signature: '元气满满的小姐姐，陪你开心打游戏~',
    },
    {
      id: 8,
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player8',
      nickname: '钻石打野',
      gameName: '英雄联盟',
      rating: 4.7,
      reviewCount: 198,
      price: 55,
      isOnline: true,
      tags: ['打野', '节奏', '带飞'],
      signature: '专业打野，节奏大师，带你轻松上分！',
    },
  ]);

  const handleFilterChange = (key: keyof FilterState, value: any) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
  };

  const handlePlayerClick = (id: number) => {
    navigate(`/user/players/${id}`);
  };

  // 筛选逻辑
  const filteredPlayers = players.filter((player) => {
    // 搜索过滤
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      if (
        !player.nickname.toLowerCase().includes(query) &&
        !player.gameName.toLowerCase().includes(query)
      ) {
        return false;
      }
    }

    // 游戏过滤
    if (filters.gameId !== null && filters.gameId !== 0) {
      const selectedGame = GAME_OPTIONS.find((g) => g.id === filters.gameId);
      if (selectedGame && player.gameName !== selectedGame.name) {
        return false;
      }
    }

    // 评分过滤
    if (player.rating < filters.minRating) {
      return false;
    }

    // 价格过滤
    if (player.price > filters.maxPrice) {
      return false;
    }

    // 在线状态过滤
    if (filters.onlineOnly && !player.isOnline) {
      return false;
    }

    return true;
  });

  // 排序逻辑
  const sortedPlayers = [...filteredPlayers].sort((a, b) => {
    switch (filters.sortBy) {
      case 'rating':
        return b.rating - a.rating;
      case 'price-low':
        return a.price - b.price;
      case 'price-high':
        return b.price - a.price;
      case 'reviews':
        return b.reviewCount - a.reviewCount;
      default:
        return 0;
    }
  });

  const resetFilters = () => {
    setFilters({
      gameId: null,
      minRating: 0,
      maxPrice: 999,
      onlineOnly: false,
      sortBy: 'rating',
    });
    setSearchQuery('');
  };

  return (
    <DiscordLayout
      serverList={
        <div className={styles.serverList}>
          <div className={styles.serverItem}>
            <span className={styles.serverIcon}>🎮</span>
          </div>
          <div className={styles.serverItem}>
            <span className={styles.serverIcon}>🏆</span>
          </div>
          <div className={styles.serverItem}>
            <span className={styles.serverIcon}>⭐</span>
          </div>
        </div>
      }
      showMemberPanel={false}
    >
      <div className={styles.playersPage}>
        {/* 顶部搜索栏 */}
        <div className={styles.header}>
          <h1 className={styles.title}>发现陪玩师</h1>
          <div className={styles.searchBar}>
            <input
              type="text"
              className={styles.searchInput}
              placeholder="搜索陪玩师或游戏..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
        </div>

        {/* 筛选栏 */}
        <div className={styles.filters}>
          <div className={styles.filterGroup}>
            <label className={styles.filterLabel}>游戏</label>
            <select
              className={styles.filterSelect}
              value={filters.gameId ?? 0}
              onChange={(e) => handleFilterChange('gameId', Number(e.target.value))}
            >
              {GAME_OPTIONS.map((game) => (
                <option key={game.id} value={game.id}>
                  {game.name}
                </option>
              ))}
            </select>
          </div>

          <div className={styles.filterGroup}>
            <label className={styles.filterLabel}>最低评分</label>
            <select
              className={styles.filterSelect}
              value={filters.minRating}
              onChange={(e) => handleFilterChange('minRating', Number(e.target.value))}
            >
              <option value={0}>全部</option>
              <option value={3.0}>3.0+</option>
              <option value={4.0}>4.0+</option>
              <option value={4.5}>4.5+</option>
              <option value={4.8}>4.8+</option>
            </select>
          </div>

          <div className={styles.filterGroup}>
            <label className={styles.filterLabel}>最高价格</label>
            <select
              className={styles.filterSelect}
              value={filters.maxPrice}
              onChange={(e) => handleFilterChange('maxPrice', Number(e.target.value))}
            >
              <option value={999}>不限</option>
              <option value={50}>≤ ¥50/小时</option>
              <option value={80}>≤ ¥80/小时</option>
              <option value={100}>≤ ¥100/小时</option>
            </select>
          </div>

          <div className={styles.filterGroup}>
            <label className={styles.filterLabel}>排序</label>
            <select
              className={styles.filterSelect}
              value={filters.sortBy}
              onChange={(e) =>
                handleFilterChange('sortBy', e.target.value as FilterState['sortBy'])
              }
            >
              <option value="rating">评分最高</option>
              <option value="reviews">评价最多</option>
              <option value="price-low">价格最低</option>
              <option value="price-high">价格最高</option>
            </select>
          </div>

          <div className={styles.filterGroup}>
            <label className={styles.filterCheckbox}>
              <input
                type="checkbox"
                checked={filters.onlineOnly}
                onChange={(e) => handleFilterChange('onlineOnly', e.target.checked)}
              />
              <span>仅在线</span>
            </label>
          </div>

          <button className={styles.resetButton} onClick={resetFilters}>
            重置
          </button>
        </div>

        {/* 陪玩师列表 */}
        <div className={styles.content}>
          {loading ? (
            <div className={styles.loading}>
              <div className={styles.spinner} />
              <p>加载中...</p>
            </div>
          ) : sortedPlayers.length > 0 ? (
            <div className={styles.playerGrid}>
              {sortedPlayers.map((player) => (
                <PlayerCard key={player.id} {...player} onClick={handlePlayerClick} />
              ))}
            </div>
          ) : (
            <div className={styles.empty}>
              <p>😔 没有找到符合条件的陪玩师</p>
              <button className={styles.emptyButton} onClick={resetFilters}>
                清除筛选条件
              </button>
            </div>
          )}
        </div>
      </div>
    </DiscordLayout>
  );
};

export default PlayersList;
