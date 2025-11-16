/**
 * 陪玩师详情页面 - Discord风格
 */

import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { DiscordLayout } from '@/components/Layout';
import styles from './PlayerDetail.module.less';

/**
 * 服务项目接口
 */
interface ServiceItem {
  id: number;
  name: string;
  description: string;
  price: number;
  duration: number; // 小时
}

/**
 * 评价接口
 */
interface Review {
  id: number;
  userId: number;
  username: string;
  userAvatar: string;
  rating: number;
  content: string;
  createdAt: string;
}

/**
 * 陪玩师详情接口
 */
interface PlayerDetail {
  id: number;
  avatar: string;
  nickname: string;
  gameName: string;
  rating: number;
  reviewCount: number;
  totalOrders: number;
  signature: string;
  bio: string;
  tags: string[];
  isOnline: boolean;
  services: ServiceItem[];
  reviews: Review[];
}

/**
 * 陪玩师详情页面组件
 */
export const PlayerDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [player, setPlayer] = useState<PlayerDetail | null>(null);
  const [selectedService, setSelectedService] = useState<number | null>(null);

  useEffect(() => {
    // Mock数据加载（后续替换为真实API调用）
    setLoading(true);
    setTimeout(() => {
      setPlayer({
        id: Number(id),
        avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=player${id}`,
        nickname: '王者大神',
        gameName: '王者荣耀',
        rating: 4.8,
        reviewCount: 256,
        totalOrders: 512,
        signature: '国服最强打野，带你上王者！',
        bio: '大家好，我是王者大神。拥有3年王者荣耀职业经验，擅长打野位置，熟练掌握多个英雄。带过上千个玩家上分，成功率高达85%。游戏风格稳健，注重团队配合，善于教学。无论你是想快速上分还是学习技术，我都能满足你的需求！',
        tags: ['打野', '上分快', '技术好', '耐心', '教学'],
        isOnline: true,
        services: [
          {
            id: 1,
            name: '上分陪玩',
            description: '带你快速上分，保证胜率',
            price: 50,
            duration: 1,
          },
          {
            id: 2,
            name: '技术教学',
            description: '一对一教学，提升游戏技术',
            price: 80,
            duration: 2,
          },
          {
            id: 3,
            name: '娱乐陪玩',
            description: '轻松娱乐，开心游戏',
            price: 40,
            duration: 1,
          },
        ],
        reviews: [
          {
            id: 1,
            userId: 101,
            username: '玩家A',
            userAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=user1',
            rating: 5,
            content: '技术很好，带我从钻石打到王者，非常耐心！强烈推荐！',
            createdAt: '2025-11-15 18:30:00',
          },
          {
            id: 2,
            userId: 102,
            username: '玩家B',
            userAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=user2',
            rating: 5,
            content: '人很好，教了很多实用技巧，游戏水平提升明显。',
            createdAt: '2025-11-14 20:15:00',
          },
          {
            id: 3,
            userId: 103,
            username: '玩家C',
            userAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=user3',
            rating: 4,
            content: '整体不错，就是有时候太认真了，压力有点大哈哈。',
            createdAt: '2025-11-13 15:45:00',
          },
        ],
      });
      setLoading(false);
    }, 500);
  }, [id]);

  const handleOrder = () => {
    if (selectedService === null) {
      alert('请先选择服务项目');
      return;
    }
    navigate(`/user/orders/create?playerId=${id}&serviceId=${selectedService}`);
  };

  const renderStars = (rating: number) => {
    const stars = [];
    for (let i = 0; i < 5; i++) {
      stars.push(
        <span key={i} className={i < rating ? styles.starFull : styles.starEmpty}>
          ★
        </span>
      );
    }
    return stars;
  };

  if (loading) {
    return (
      <DiscordLayout showMemberPanel={false}>
        <div className={styles.loading}>
          <div className={styles.spinner} />
          <p>加载中...</p>
        </div>
      </DiscordLayout>
    );
  }

  if (!player) {
    return (
      <DiscordLayout showMemberPanel={false}>
        <div className={styles.error}>
          <p>❌ 陪玩师不存在</p>
          <button onClick={() => navigate('/user/players')}>返回列表</button>
        </div>
      </DiscordLayout>
    );
  }

  return (
    <DiscordLayout
      showMemberPanel={true}
      memberPanel={
        <div className={styles.sidebar}>
          <h3 className={styles.sidebarTitle}>选择服务</h3>
          <div className={styles.serviceList}>
            {player.services.map((service) => (
              <div
                key={service.id}
                className={`${styles.serviceItem} ${
                  selectedService === service.id ? styles.selected : ''
                }`}
                onClick={() => setSelectedService(service.id)}
              >
                <div className={styles.serviceName}>{service.name}</div>
                <div className={styles.serviceDesc}>{service.description}</div>
                <div className={styles.servicePrice}>
                  <span className={styles.priceLabel}>￥</span>
                  <span className={styles.priceValue}>{service.price}</span>
                  <span className={styles.priceUnit}>/{service.duration}小时</span>
                </div>
              </div>
            ))}
          </div>
          <button className={styles.orderButton} onClick={handleOrder}>
            立即下单
          </button>
        </div>
      }
    >
      <div className={styles.playerDetail}>
        {/* 头部信息 */}
        <div className={styles.header}>
          <div className={styles.avatarSection}>
            <div className={styles.avatarWrapper}>
              <img src={player.avatar} alt={player.nickname} className={styles.avatar} />
              {player.isOnline && <span className={styles.onlineIndicator} />}
            </div>
          </div>

          <div className={styles.infoSection}>
            <div className={styles.nameRow}>
              <h1 className={styles.nickname}>{player.nickname}</h1>
              <span className={styles.gameBadge}>{player.gameName}</span>
            </div>

            <div className={styles.signature}>{player.signature}</div>

            <div className={styles.statsRow}>
              <div className={styles.statItem}>
                <div className={styles.ratingStars}>{renderStars(player.rating)}</div>
                <span className={styles.ratingText}>{player.rating.toFixed(1)}</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statLabel}>评价</span>
                <span className={styles.statValue}>{player.reviewCount}</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statLabel}>订单</span>
                <span className={styles.statValue}>{player.totalOrders}</span>
              </div>
            </div>

            <div className={styles.tagsRow}>
              {player.tags.map((tag, index) => (
                <span key={index} className={styles.tag}>
                  {tag}
                </span>
              ))}
            </div>
          </div>
        </div>

        {/* 个人简介 */}
        <div className={styles.section}>
          <h2 className={styles.sectionTitle}>个人简介</h2>
          <p className={styles.bio}>{player.bio}</p>
        </div>

        {/* 用户评价 */}
        <div className={styles.section}>
          <h2 className={styles.sectionTitle}>用户评价 ({player.reviews.length})</h2>
          <div className={styles.reviewList}>
            {player.reviews.map((review) => (
              <div key={review.id} className={styles.reviewItem}>
                <div className={styles.reviewHeader}>
                  <img
                    src={review.userAvatar}
                    alt={review.username}
                    className={styles.reviewAvatar}
                  />
                  <div className={styles.reviewInfo}>
                    <div className={styles.reviewName}>{review.username}</div>
                    <div className={styles.reviewMeta}>
                      <div className={styles.reviewStars}>{renderStars(review.rating)}</div>
                      <span className={styles.reviewDate}>{review.createdAt}</span>
                    </div>
                  </div>
                </div>
                <div className={styles.reviewContent}>{review.content}</div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </DiscordLayout>
  );
};

export default PlayerDetail;
