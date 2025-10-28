import React from 'react';
import { Card } from '../../components';
import styles from './PaymentList.module.less';

export const PaymentList: React.FC = () => {
  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>支付管理</h1>
      </div>

      <Card className={styles.content}>
        <p className={styles.placeholder}>💰 支付管理模块开发中...</p>
      </Card>
    </div>
  );
};
