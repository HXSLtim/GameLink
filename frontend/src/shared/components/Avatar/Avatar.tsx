/**
 * Discord风格头像组件
 */
import React from 'react';
import './Avatar.less';

export interface AvatarProps {
  /**
   * 头像URL
   */
  src?: string;

  /**
   * 备用文本（当图片加载失败时显示）
   */
  alt?: string;

  /**
   * 头像大小
   */
  size?: 'small' | 'medium' | 'large' | 'xlarge';

  /**
   * 在线状态
   */
  status?: 'online' | 'idle' | 'busy' | 'offline';

  /**
   * 是否显示状态指示器
   */
  showStatus?: boolean;

  /**
   * 自定义类名
   */
  className?: string;

  /**
   * 点击事件
   */
  onClick?: () => void;
}

const Avatar: React.FC<AvatarProps> = ({
  src,
  alt = 'Avatar',
  size = 'medium',
  status = 'offline',
  showStatus = false,
  className = '',
  onClick,
}) => {
  const [imageError, setImageError] = React.useState(false);

  const handleImageError = () => {
    setImageError(true);
  };

  const classNames = [
    'discord-avatar',
    `discord-avatar--${size}`,
    onClick && 'discord-avatar--clickable',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  const statusClassName = [
    'discord-avatar__status',
    `discord-avatar__status--${status}`,
  ].join(' ');

  // 从alt获取首字母作为备用显示
  const fallbackText = alt
    .split(' ')
    .map((word) => word[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);

  return (
    <div className={classNames} onClick={onClick}>
      <div className="discord-avatar__image">
        {!imageError && src ? (
          <img src={src} alt={alt} onError={handleImageError} />
        ) : (
          <div className="discord-avatar__fallback">{fallbackText}</div>
        )}
      </div>
      {showStatus && <span className={statusClassName} />}
    </div>
  );
};

export default Avatar;
