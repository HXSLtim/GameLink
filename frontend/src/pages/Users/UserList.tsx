import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Button, Input, Select, Table, Tag, Pagination } from '../../components';
import { User, UserRole, UserStatus } from '../../types/user.types';
import { getMockUserList } from '../../services/userMockData';
import {
  formatUserRole,
  getUserRoleColor,
  formatUserStatus,
  getUserStatusColor,
  formatPhone,
  formatEmail,
} from '../../utils/userFormatters';
import { formatDateTime, formatRelativeTime } from '../../utils/formatters';
import styles from './UserList.module.less';

const SearchIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <circle cx="11" cy="11" r="8" strokeWidth="2" />
    <path d="M21 21L16.65 16.65" strokeWidth="2" strokeLinecap="round" />
  </svg>
);

export const UserList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [users, setUsers] = useState<User[]>([]);
  const [total, setTotal] = useState(0);

  // 筛选条件
  const [keyword, setKeyword] = useState('');
  const [role, setRole] = useState<UserRole | ''>('');
  const [status, setStatus] = useState<UserStatus | ''>('');

  // 分页
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  // 加载数据
  const loadData = () => {
    setLoading(true);

    // 模拟 API 调用延迟
    setTimeout(() => {
      const result = getMockUserList({
        page,
        page_size: pageSize,
        keyword: keyword || undefined,
        role: role || undefined,
        status: status || undefined,
      });

      setUsers(result.list);
      setTotal(result.total);
      setLoading(false);
    }, 300);
  };

  // 搜索
  const handleSearch = () => {
    setPage(1);
    loadData();
  };

  // 重置筛选
  const handleReset = () => {
    setKeyword('');
    setRole('');
    setStatus('');
    setPage(1);
  };

  // 页码变化
  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };

  // 初始加载和依赖变化时重新加载
  useEffect(() => {
    loadData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  // 表格列定义
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: '80px',
    },
    {
      title: '用户信息',
      key: 'userInfo',
      render: (_: unknown, record: User) => (
        <div className={styles.userInfo}>
          {record.avatar_url && (
            <img src={record.avatar_url} alt={record.name} className={styles.avatar} />
          )}
          {!record.avatar_url && (
            <div className={styles.avatarPlaceholder}>{record.name.charAt(0)}</div>
          )}
          <div className={styles.userDetails}>
            <div className={styles.userName}>{record.name}</div>
            <div className={styles.userContact}>
              {formatPhone(record.phone)} · {formatEmail(record.email)}
            </div>
          </div>
        </div>
      ),
    },
    {
      title: '角色',
      key: 'role',
      width: '120px',
      render: (_: unknown, record: User) => (
        <Tag color={getUserRoleColor(record.role)}>{formatUserRole(record.role)}</Tag>
      ),
    },
    {
      title: '状态',
      key: 'status',
      width: '100px',
      render: (_: unknown, record: User) => (
        <Tag color={getUserStatusColor(record.status)}>{formatUserStatus(record.status)}</Tag>
      ),
    },
    {
      title: '最后登录',
      key: 'lastLogin',
      width: '160px',
      render: (_: unknown, record: User) => (
        <div className={styles.timeInfo}>
          {record.last_login_at ? (
            <>
              <div>{formatRelativeTime(record.last_login_at)}</div>
              <div className={styles.timeDetail}>{formatDateTime(record.last_login_at)}</div>
            </>
          ) : (
            '从未登录'
          )}
        </div>
      ),
    },
    {
      title: '注册时间',
      key: 'createdAt',
      width: '160px',
      render: (_: unknown, record: User) => formatDateTime(record.created_at),
    },
    {
      title: '操作',
      key: 'actions',
      width: '120px',
      render: (_: unknown, record: User) => (
        <div className={styles.actions}>
          <Button
            variant="text"
            onClick={() => navigate(`/users/${record.id}`)}
            className={styles.actionButton}
          >
            查看详情
          </Button>
        </div>
      ),
    },
  ];

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>用户管理</h1>
      </div>

      {/* 筛选区域 */}
      <Card className={styles.filterCard}>
        <div className={styles.filters}>
          <div className={styles.filterRow}>
            <div className={styles.filterItem}>
              <label className={styles.filterLabel}>搜索</label>
              <Input
                value={keyword}
                onChange={(e) => setKeyword(e.target.value)}
                placeholder="用户名/手机号/邮箱"
                className={styles.filterInput}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
              />
            </div>

            <div className={styles.filterItem}>
              <label className={styles.filterLabel}>角色</label>
              <Select
                value={role}
                onChange={(value) => setRole(value as UserRole | '')}
                options={[
                  { label: '全部角色', value: '' },
                  { label: '普通用户', value: UserRole.USER },
                  { label: '陪玩师', value: UserRole.PLAYER },
                  { label: '管理员', value: UserRole.ADMIN },
                ]}
                className={styles.filterSelect}
              />
            </div>

            <div className={styles.filterItem}>
              <label className={styles.filterLabel}>状态</label>
              <Select
                value={status}
                onChange={(value) => setStatus(value as UserStatus | '')}
                options={[
                  { label: '全部状态', value: '' },
                  { label: '正常', value: UserStatus.ACTIVE },
                  { label: '暂停', value: UserStatus.SUSPENDED },
                  { label: '封禁', value: UserStatus.BANNED },
                ]}
                className={styles.filterSelect}
              />
            </div>
          </div>

          <div className={styles.filterActions}>
            <Button variant="primary" onClick={handleSearch} className={styles.filterButton}>
              <SearchIcon />
              搜索
            </Button>
            <Button variant="outlined" onClick={handleReset} className={styles.filterButton}>
              重置
            </Button>
          </div>
        </div>
      </Card>

      {/* 数据表格 */}
      <Card className={styles.tableCard}>
        <Table columns={columns} dataSource={users} loading={loading} rowKey="id" />

        {/* 分页 */}
        {total > 0 && (
          <div className={styles.pagination}>
            <Pagination
              current={page}
              pageSize={pageSize}
              total={total}
              onChange={handlePageChange}
            />
          </div>
        )}
      </Card>
    </div>
  );
};
