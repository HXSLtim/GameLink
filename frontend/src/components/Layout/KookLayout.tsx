/**
 * KOOK 风格两栏布局组件
 *
 * 布局结构:
 * ┌─────────────────────────────────────────────┐
 * │          Top Navigation (60px)              │
 * ├──────────┬──────────────────────────────────┤
 * │          │                                   │
 * │ Channel  │      Main Content Area            │
 * │  List    │                                   │
 * │          │                                   │
 * │ (240px)  │                                   │
 * │          │                                   │
 * └──────────┴──────────────────────────────────┘
 */

import { ReactNode, useState } from 'react';
import styles from './KookLayout.module.less';

export interface KookLayoutProps {
  /** 顶部导航栏 (60px) */
  topNav?: ReactNode;

  /** 左侧频道/模块列表 (240px) */
  channelList?: ReactNode;

  /** 主内容区 */
  children: ReactNode;

  /** 是否显示左侧频道列表 */
  showChannelList?: boolean;

  /** 频道列表初始折叠状态 */
  channelListCollapsed?: boolean;

  /** 自定义类名 */
  className?: string;

  /** 频道列表宽度 (默认240px) */
  channelListWidth?: number;
}

/**
 * KOOK 风格两栏布局组件
 */
export const KookLayout: React.FC<KookLayoutProps> = ({
  topNav,
  channelList,
  children,
  showChannelList = true,
  channelListCollapsed = false,
  className,
  channelListWidth = 240,
}) => {
  const [isChannelCollapsed, setIsChannelCollapsed] = useState(channelListCollapsed);

  const toggleChannelList = () => {
    setIsChannelCollapsed(!isChannelCollapsed);
  };

  return (
    <div className={`${styles.kookLayout} ${className || ''}`}>
      {/* 顶部导航栏 */}
      {topNav && (
        <div className={styles.topNav}>
          {topNav}
        </div>
      )}

      <div className={styles.contentWrapper}>
        {/* 左侧频道列表 */}
        {showChannelList && channelList && (
          <>
            <div
              className={`${styles.channelList} ${
                isChannelCollapsed ? styles.collapsed : ''
              }`}
              style={{
                width: isChannelCollapsed ? 0 : `${channelListWidth}px`,
              }}
            >
              {channelList}
            </div>

            {/* 折叠按钮 */}
            <button
              className={styles.toggleButton}
              onClick={toggleChannelList}
              aria-label={isChannelCollapsed ? '展开频道列表' : '折叠频道列表'}
              style={{
                left: isChannelCollapsed ? 0 : `${channelListWidth}px`,
              }}
            >
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                className={isChannelCollapsed ? styles.rotated : ''}
              >
                <path
                  fill="currentColor"
                  d="M15.4 7.4L14 6l-6 6l6 6l1.4-1.4L10.8 12z"
                />
              </svg>
            </button>
          </>
        )}

        {/* 主内容区 */}
        <div className={styles.mainContent}>
          {children}
        </div>
      </div>
    </div>
  );
};

export default KookLayout;
