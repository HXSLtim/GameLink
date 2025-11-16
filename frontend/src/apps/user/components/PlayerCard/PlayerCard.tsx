/**
 * 陪玩师卡片组件
 */
import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '@/shared/components/Card';
import { Avatar } from '@/shared/components/Avatar';
import { Button } from '@/shared/components/Button';
import './PlayerCard.less';

export interface PlayerCardProps {
  /**
   * 陪玩师ID
   */
  id: number;

  /**
   * 头像URL
   */
  avatar: string;

  /**
   * 陪玩师名称
   */
  name: string;

  /**
   * 游戏名称
   */
  game: string;

  /**
   * 评分 (1-5)
   */
  rating: number;

  /**
   * 评价数量
   */
  reviewCount: number;

  /**
   * 价格（每小时）
   */
  pricePerHour: number;

  /**
   * 在线状态
   */
  status: 'online' | 'idle' | 'busy' | 'offline';

  /**
   * 简介
   */
  bio?: string;

  /**
   * 标签
   */
  tags?: string[];

  /**
   * 是否已接单（busy状态显示）
   */
  isBusy?: boolean;
}

const PlayerCard: React.FC<PlayerCardProps> = ({
  id,
  avatar,
  name,
  game,
  rating,
  reviewCount,
  pricePerHour,
  status,
  bio,
  tags = [],
  isBusy = false,
}) => {
  const navigate = useNavigate();

  const handleCardClick = () => {
    navigate(`/user/players/${id}`);
  };

  const handleOrderClick = (e: React.MouseEvent) => {
    e.stopPropagation(); // 阻止卡片点击事件
    navigate(`/user/players/${id}?action=order`);
  };

  const renderStars = () => {
    const stars = [];
    const fullStars = Math.floor(rating);
    const hasHalfStar = rating % 1 >= 0.5;

    for (let i = 0; i < 5; i++) {
      if (i < fullStars) {
        stars.push(
          <span key={i} className="player-card__star player-card__star--full">
            ★
          </span>
        );
      } else if (i === fullStars && hasHalfStar) {
        stars.push(
          <span key={i} className="player-card__star player-card__star--half">
            ★
          </span>
        );
      } else {
        stars.push(
          <span key={i} className="player-card__star player-card__star--empty">
            ☆
          </span>
        );
      }
    }

    return stars;
  };

  const getStatusText = () => {
    if (isBusy) return '接单中';
    switch (status) {
      case 'online':
        return '在线';
      case 'idle':
        return '离开';
      case 'busy':
        return '忙碌';
      case 'offline':
        return '离线';
      default:
        return '离线';
    }
  };

  return (
    <Card
      className="player-card"
      hoverable
      clickable
      onClick={handleCardClick}
      padding="medium"
      bordered
    >
      <div className="player-card__content">
        {/* 头像和在线状态 */}
        <div className="player-card__avatar-section">
          <Avatar
            src={avatar}
            alt={name}
            size="large"
            status={status}
            showStatus
          />
        </div>

        {/* 主要信息 */}
        <div className="player-card__info">
          <div className="player-card__header">
            <h3 className="player-card__name">{name}</h3>
            <span className={`player-card__status player-card__status--${status}`}>
              {getStatusText()}
            </span>
          </div>

          <div className="player-card__game">{game}</div>

          {/* 评分和评价 */}
          <div className="player-card__rating">
            <div className="player-card__stars">{renderStars()}</div>
            <span className="player-card__review-count">
              {rating.toFixed(1)} ({reviewCount}评价)
            </span>
          </div>

          {/* 简介 */}
          {bio && <p className="player-card__bio">{bio}</p>}

          {/* 标签 */}
          {tags.length > 0 && (
            <div className="player-card__tags">
              {tags.map((tag, index) => (
                <span key={index} className="player-card__tag">
                  {tag}
                </span>
              ))}
            </div>
          )}

          {/* 价格和下单按钮 */}
          <div className="player-card__footer">
            <div className="player-card__price">
              <span className="player-card__price-amount">¥{pricePerHour}</span>
              <span className="player-card__price-unit">/小时</span>
            </div>
            <Button
              type="primary"
              size="small"
              disabled={status === 'offline' || isBusy}
              onClick={handleOrderClick}
            >
              {isBusy ? '接单中' : '立即下单'}
            </Button>
          </div>
        </div>
      </div>
    </Card>
  );
};

export default PlayerCard;
