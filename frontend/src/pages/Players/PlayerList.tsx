import React from 'react';
import { Card } from '../../components';
import styles from './PlayerList.module.less';

export const PlayerList: React.FC = () => {
  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>陪玩师管理</h1>
      </div>

      <Card className={styles.content}>
        <p className={styles.placeholder}>🎯 陪玩师管理模块开发中...</p>
      </Card>
    </div>
  );
};
