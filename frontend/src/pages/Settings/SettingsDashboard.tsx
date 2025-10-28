import React from 'react';
import { Card } from '../../components';
import styles from './SettingsDashboard.module.less';

export const SettingsDashboard: React.FC = () => {
  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>系统设置</h1>
      </div>

      <Card className={styles.content}>
        <p className={styles.placeholder}>⚙️ 系统设置模块开发中...</p>
      </Card>
    </div>
  );
};
