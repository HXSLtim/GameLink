/**
 * Discord 风格三栏布局组件
 *
 * 布局结构:
 * ┌─────┬────────────────────────────┬──────────┐
 * │     │                            │          │
 * │  S  │      Main Content          │  Member  │
 * │  e  │                            │          │
 * │  r  │                            │  Panel   │
 * │  v  │                            │          │
 * │  e  │                            │  (可选)   │
 * │  r  │                            │          │
 * │  s  │                            │          │
 * │     │                            │          │
 * └─────┴────────────────────────────┴──────────┘
 */

import { ReactNode, useState } from 'react';
import styles from './DiscordLayout.module.less';

export interface DiscordLayoutProps {
  /** 左侧服务器/模块列表 (50-60px) */
  serverList?: ReactNode;

  /** 右侧成员/详情面板 (240-280px) */
  memberPanel?: ReactNode;

  /** 主内容区 */
  children: ReactNode;

  /** 是否显示右侧面板 */
  showMemberPanel?: boolean;

  /** 右侧面板初始折叠状态 */
  memberPanelCollapsed?: boolean;

  /** 自定义类名 */
  className?: string;
}

/**
 * Discord 风格三栏布局组件
 */
export const DiscordLayout: React.FC<DiscordLayoutProps> = ({
  serverList,
  memberPanel,
  children,
  showMemberPanel = false,
  memberPanelCollapsed = false,
  className,
}) => {
  const [isPanelCollapsed, setIsPanelCollapsed] = useState(memberPanelCollapsed);

  const toggleMemberPanel = () => {
    setIsPanelCollapsed(!isPanelCollapsed);
  };

  return (
    <div className={`${styles.discordLayout} ${className || ''}`}>
      {/* 左侧服务器列表 */}
      {serverList && (
        <div className={styles.serverList}>
          {serverList}
        </div>
      )}

      {/* 主内容区 */}
      <div className={styles.mainContent}>
        {children}
      </div>

      {/* 右侧成员/详情面板 */}
      {showMemberPanel && memberPanel && (
        <>
          <div
            className={`${styles.memberPanel} ${
              isPanelCollapsed ? styles.collapsed : ''
            }`}
          >
            {memberPanel}
          </div>

          {/* 折叠按钮 */}
          <button
            className={styles.toggleButton}
            onClick={toggleMemberPanel}
            aria-label={isPanelCollapsed ? '展开侧边栏' : '折叠侧边栏'}
          >
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              className={isPanelCollapsed ? styles.rotated : ''}
            >
              <path
                fill="currentColor"
                d="M9.4 18L8 16.6l4.6-4.6L8 7.4L9.4 6l6 6z"
              />
            </svg>
          </button>
        </>
      )}
    </div>
  );
};

export default DiscordLayout;
