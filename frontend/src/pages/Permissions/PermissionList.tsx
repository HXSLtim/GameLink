import React, { useState } from 'react';
import { Tabs } from '../../components';
import { RoleManagement } from './RoleManagement';
import { PermissionManagement } from './PermissionManagement';
import styles from './PermissionList.module.less';

/**
 * 权限管理页面
 * 包含角色管理和权限管理两个Tab
 */
export const PermissionList: React.FC = () => {
  const [activeTab, setActiveTab] = useState('roles');

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>权限管理</h1>
      </div>

      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        items={[
          {
            key: 'roles',
            label: '角色管理',
            children: <RoleManagement />,
          },
          {
            key: 'permissions',
            label: '权限管理',
            children: <PermissionManagement />,
          },
        ]}
      />
    </div>
  );
};
