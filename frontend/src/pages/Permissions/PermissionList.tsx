import React from 'react';
import { Card } from '../../components';
import styles from './PermissionList.module.less';

export const PermissionList: React.FC = () => {
  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>权限管理</h1>
      </div>

      <Card className={styles.content}>
        <p className={styles.placeholder}>🔐 权限管理模块开发中...</p>
      </Card>
    </div>
  );
};
