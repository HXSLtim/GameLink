import { Skeleton } from '@arco-design/web-react';
import styles from './PageSkeleton.module.less';

/**
 * 页面骨架屏组件
 * 用于页面懒加载时的占位显示
 */
export const PageSkeleton: React.FC = () => {
  return (
    <div className={styles.pageSkeleton} aria-busy="true" aria-label="页面加载中">
      <div className={styles.header}>
        <Skeleton
          animation
          text={{
            rows: 1,
            width: ['30%'],
          }}
          style={{ marginBottom: 16 }}
        />
      </div>
      <div className={styles.content}>
        <Skeleton
          animation
          text={{
            rows: 4,
            width: ['100%', '90%', '95%', '85%'],
          }}
          style={{ marginBottom: 24 }}
        />
        <Skeleton
          animation
          text={{
            rows: 6,
            width: ['100%', '95%', '90%', '95%', '85%', '90%'],
          }}
        />
      </div>
    </div>
  );
};

