import React from 'react';
import { Card } from '../../components';
import styles from './GameList.module.less';

export const GameList: React.FC = () => {
  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>游戏管理</h1>
      </div>

      <Card className={styles.content}>
        <p className={styles.placeholder}>🎮 游戏管理模块开发中...</p>
      </Card>
    </div>
  );
};
