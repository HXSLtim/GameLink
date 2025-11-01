import React, { useState } from 'react';
import { Tabs } from '../../components';
import { PlatformSettings } from './PlatformSettings';
import { PaymentSettings } from './PaymentSettings';
import { OrderSettings } from './OrderSettings';
import { UserSettings } from './UserSettings';
import { AnnouncementSettings } from './AnnouncementSettings';
import styles from './SettingsDashboard.module.less';

/**
 * 系统设置页面
 */
export const SettingsDashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState('platform');

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>系统设置</h1>
      </div>

      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        items={[
          {
            key: 'platform',
            label: '平台配置',
            children: <PlatformSettings />,
          },
          {
            key: 'payment',
            label: '支付设置',
            children: <PaymentSettings />,
          },
          {
            key: 'order',
            label: '订单设置',
            children: <OrderSettings />,
          },
          {
            key: 'user',
            label: '用户设置',
            children: <UserSettings />,
          },
          {
            key: 'announcement',
            label: '公告管理',
            children: <AnnouncementSettings />,
          },
        ]}
      />
    </div>
  );
};
