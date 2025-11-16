/**
 * 陪玩师卡片组件
 * 用于展示陪玩师的基本信息、评分、价格等
 */

import { FC } from 'react';
import styles from './PlayerCard.module.less';

export interface PlayerCardProps {
  /** 陪玩师ID */
  id: number;

  /** 头像URL */
  avatar: string;

  /** 昵称 */
  nickname: string;

  /** 游戏名称 */
  gameName: string;

  /** 评分 (0-5) */
  rating: number;

  /** 评价数量 */
  reviewCount: number;

  /** 价格（元/小时） */
  price: number;

  /** 在线状态 */
  isOnline: boolean;

  /** 服务标签 */
  tags?: string[];

  /** 个性签名 */
  signature?: string;

  /** 点击事件 */
  onClick?: (id: number) => void;

  /** 自定义类名 */
  className?: string;
}

/**
 * 陪玩师卡片组件
 */
export const PlayerCard: FC<PlayerCardProps> = ({
  id,
  avatar,
  nickname,
  gameName,
  rating,
  reviewCount,
  price,
  isOnline,
  tags = [],
  signature,
  onClick,
  className,
}) => {
  const handleClick = () => {
    onClick?.(id);
  };

  const renderStars = (rating: number) => {
    const stars = [];
    const fullStars = Math.floor(rating);
    const hasHalfStar = rating % 1 >= 0.5;

    for (let i = 0; i < 5; i++) {
      if (i < fullStars) {
        stars.push(
          <span key={i} className={styles.starFull}>
            ★
          </span>
        );
      } else if (i === fullStars && hasHalfStar) {
        stars.push(
          <span key={i} className={styles.starHalf}>
            ★
          </span>
        );
      } else {
        stars.push(
          <span key={i} className={styles.starEmpty}>
            ☆
          </span>
        );
      }
    }
    return stars;
  };

  return (
    <div
      className={`${styles.playerCard} ${className || ''}`}
      onClick={handleClick}
      role="button"
      tabIndex={0}
      onKeyPress={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          handleClick();
        }
      }}
    >
      {/* 头像区域 */}
      <div className={styles.avatarSection}>
        <div className={styles.avatarWrapper}>
          <img
            src={avatar}
            alt={nickname}
            className={styles.avatar}
            loading="lazy"
          />
          {isOnline && <span className={styles.onlineIndicator} />}
        </div>
      </div>

      {/* 信息区域 */}
      <div className={styles.infoSection}>
        {/* 昵称和游戏 */}
        <div className={styles.header}>
          <h3 className={styles.nickname}>{nickname}</h3>
          <span className={styles.gameName}>{gameName}</span>
        </div>

        {/* 评分 */}
        <div className={styles.ratingSection}>
          <div className={styles.stars}>{renderStars(rating)}</div>
          <span className={styles.ratingText}>
            {rating.toFixed(1)}
          </span>
          <span className={styles.reviewCount}>({reviewCount})</span>
        </div>

        {/* 个性签名 */}
        {signature && (
          <p className={styles.signature}>{signature}</p>
        )}

        {/* 服务标签 */}
        {tags.length > 0 && (
          <div className={styles.tagsSection}>
            {tags.slice(0, 3).map((tag, index) => (
              <span key={index} className={styles.tag}>
                {tag}
              </span>
            ))}
          </div>
        )}

        {/* 价格 */}
        <div className={styles.priceSection}>
          <span className={styles.priceLabel}>￥</span>
          <span className={styles.priceValue}>{price}</span>
          <span className={styles.priceUnit}>/小时</span>
        </div>
      </div>
    </div>
  );
};

export default PlayerCard;
