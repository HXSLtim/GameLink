import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Button, Tag, PageSkeleton, Table, Pagination } from '../../components';
import type { TableColumn } from '../../components/Table/Table';
import { userApi, UserDetail as UserDetailType } from '../../services/api/user';
import { orderApi, OrderInfo } from '../../services/api/order';
import {
  formatUserRole,
  getUserRoleColor,
  formatUserStatus,
  getUserStatusColor,
  formatVerificationStatus,
  getVerificationStatusColor,
  formatPhone,
  formatEmail,
  formatHourlyRate,
  formatRating,
  formatPrice,
} from '../../utils/userFormatters';
import {
  formatDateTime,
  formatCurrency,
  formatOrderStatus,
  getOrderStatusColor,
} from '../../utils/formatters';
import { UserRole } from '../../types/user';
import styles from './UserDetail.module.less';

export const UserDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [userDetail, setUserDetail] = useState<UserDetailType | null>(null);

  // 用户订单列表状态
  const [ordersLoading, setOrdersLoading] = useState(false);
  const [orders, setOrders] = useState<OrderInfo[]>([]);
  const [ordersTotal, setOrdersTotal] = useState(0);
  const [ordersPage, setOrdersPage] = useState(1);
  const ordersPageSize = 10;

  // 加载用户详情
  useEffect(() => {
    const loadUserDetail = async () => {
      if (!id) {
        setError('用户ID无效');
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        setError(null);
        const data = await userApi.getDetail(Number(id));
        setUserDetail(data);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '加载用户详情失败';
        setError(errorMessage);
        console.error('加载用户详情失败:', err);
      } finally {
        setLoading(false);
      }
    };

    loadUserDetail();
  }, [id]);

  // 加载用户订单列表
  useEffect(() => {
    const loadUserOrders = async () => {
      if (!id) return;

      try {
        setOrdersLoading(true);
        const result = await orderApi.getUserOrders(Number(id), {
          page: ordersPage,
          pageSize: ordersPageSize,
        });

        if (result && result.list) {
          setOrders(result.list);
          setOrdersTotal(result.total || 0);
        } else {
          setOrders([]);
          setOrdersTotal(0);
        }
      } catch (err) {
        console.error('加载用户订单失败:', err);
        setOrders([]);
        setOrdersTotal(0);
      } finally {
        setOrdersLoading(false);
      }
    };

    loadUserOrders();
  }, [id, ordersPage]);

  // 加载中状态
  if (loading) {
    return (
      <div className={styles.container}>
        <PageSkeleton />
      </div>
    );
  }

  // 错误状态
  if (error || !userDetail) {
    return (
      <div className={styles.container}>
        <Card className={styles.errorCard}>
          <h2>用户未找到</h2>
          <p>{error || `用户 ID: ${id} 不存在`}</p>
          <Button onClick={() => navigate('/users')}>返回用户列表</Button>
        </Card>
      </div>
    );
  }

  const isPlayer = userDetail.role === UserRole.PLAYER && userDetail.player;

  // 订单表格列定义
  const orderColumns: TableColumn<OrderInfo>[] = [
    {
      title: '订单ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      width: 200,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => (
        <Tag color={getOrderStatusColor(status as any)}>{formatOrderStatus(status as any)}</Tag>
      ),
    },
    {
      title: '金额',
      dataIndex: 'priceCents',
      key: 'priceCents',
      width: 100,
      render: (priceCents: number) => formatCurrency(priceCents),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 150,
      render: (createdAt: string) => formatDateTime(createdAt),
    },
    {
      title: '操作',
      key: 'actions',
      width: 100,
      render: (_: any, record: OrderInfo) => (
        <Button variant="text" onClick={() => navigate(`/orders/${record.id}`)}>
          查看详情
        </Button>
      ),
    },
  ];

  return (
    <div className={styles.container}>
      {/* 头部 */}
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <Button variant="outlined" onClick={() => navigate('/users')}>
            ← 返回列表
          </Button>
          <h1 className={styles.title}>用户详情</h1>
        </div>
        <div className={styles.headerRight}>
          <Tag color={getUserRoleColor(userDetail.role)}>{formatUserRole(userDetail.role)}</Tag>
          <Tag color={getUserStatusColor(userDetail.status)}>
            {formatUserStatus(userDetail.status)}
          </Tag>
        </div>
      </div>

      {/* 主要内容 */}
      <div className={styles.content}>
        {/* 左侧 */}
        <div className={styles.leftColumn}>
          {/* 基本信息 */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>基本信息</h2>

            <div className={styles.userHeader}>
              {userDetail.avatarUrl && (
                <img src={userDetail.avatarUrl} alt={userDetail.name} className={styles.avatar} />
              )}
              {!userDetail.avatarUrl && (
                <div className={styles.avatarPlaceholder}>{userDetail.name.charAt(0)}</div>
              )}
              <div className={styles.userBasicInfo}>
                <div className={styles.userName}>{userDetail.name}</div>
                <div className={styles.userId}>ID: {userDetail.id}</div>
              </div>
            </div>

            <div className={styles.infoGrid}>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>手机号</span>
                <span className={styles.infoValue}>{formatPhone(userDetail.phone)}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>邮箱</span>
                <span className={styles.infoValue}>{formatEmail(userDetail.email)}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>注册时间</span>
                <span className={styles.infoValue}>{formatDateTime(userDetail.createdAt)}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>最后登录</span>
                <span className={styles.infoValue}>
                  {userDetail.lastLoginAt ? formatDateTime(userDetail.lastLoginAt) : '从未登录'}
                </span>
              </div>
            </div>
          </Card>

          {/* 统计信息 */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>统计数据</h2>
            <div className={styles.statsGrid}>
              <div className={styles.statItem}>
                <div className={styles.statValue}>{userDetail.orderCount || 0}</div>
                <div className={styles.statLabel}>订单数量</div>
              </div>
              <div className={styles.statItem}>
                <div className={styles.statValue}>
                  {userDetail.totalSpent ? formatPrice(userDetail.totalSpent) : '¥0'}
                </div>
                <div className={styles.statLabel}>总消费</div>
              </div>
              <div className={styles.statItem}>
                <div className={styles.statValue}>{userDetail.reviewCount || 0}</div>
                <div className={styles.statLabel}>评价数量</div>
              </div>
            </div>
          </Card>

          {/* 陪玩师信息 */}
          {isPlayer && userDetail.player && (
            <Card className={styles.section}>
              <h2 className={styles.sectionTitle}>
                陪玩师信息
                <Tag
                  color={getVerificationStatusColor(userDetail.player.verificationStatus)}
                  className={styles.verificationTag}
                >
                  {formatVerificationStatus(userDetail.player.verificationStatus)}
                </Tag>
              </h2>

              <div className={styles.infoGrid}>
                <div className={styles.infoItem}>
                  <span className={styles.infoLabel}>昵称</span>
                  <span className={styles.infoValue}>{userDetail.player.nickname || '-'}</span>
                </div>
                <div className={styles.infoItem}>
                  <span className={styles.infoLabel}>时薪</span>
                  <span className={styles.infoValue}>
                    {formatHourlyRate(userDetail.player.hourlyRateCents)}
                  </span>
                </div>
                <div className={styles.infoItem}>
                  <span className={styles.infoLabel}>评分</span>
                  <span className={styles.infoValue}>
                    {formatRating(userDetail.player.ratingAverage)} (
                    {userDetail.player.ratingCount}条评价)
                  </span>
                </div>
                <div className={styles.infoItem}>
                  <span className={styles.infoLabel}>认证时间</span>
                  <span className={styles.infoValue}>
                    {formatDateTime(userDetail.player.createdAt)}
                  </span>
                </div>
              </div>

              {userDetail.player.bio && (
                <div className={styles.bioSection}>
                  <h3 className={styles.subTitle}>个人简介</h3>
                  <p className={styles.bioText}>{userDetail.player.bio}</p>
                </div>
              )}
            </Card>
          )}
        </div>

        {/* 右侧 */}
        <div className={styles.rightColumn}>
          {/* 操作区域 */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>账户操作</h2>
            <div className={styles.actions}>
              <Button
                variant="outlined"
                onClick={() => console.log('暂停账户')}
                className={styles.actionButton}
                disabled={userDetail.status !== 'active'}
              >
                ⏸️ 暂停账户
              </Button>
              <Button
                variant="outlined"
                onClick={() => console.log('封禁账户')}
                className={styles.actionButton}
                disabled={userDetail.status === 'banned'}
              >
                🚫 封禁账户
              </Button>
              {userDetail.status !== 'active' && (
                <Button
                  variant="primary"
                  onClick={() => console.log('解除限制')}
                  className={styles.actionButton}
                >
                  ✅ 解除限制
                </Button>
              )}
            </div>
          </Card>

          {/* 快捷入口 */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>快捷入口</h2>
            <div className={styles.quickLinks}>
              <Button
                variant="text"
                onClick={() => console.log('查看评价')}
                className={styles.linkButton}
              >
                ⭐ 查看用户评价
              </Button>
              <Button
                variant="text"
                onClick={() => console.log('查看支付')}
                className={styles.linkButton}
              >
                💳 查看支付记录
              </Button>
            </div>
          </Card>
        </div>
      </div>

      {/* 用户订单列表 */}
      <Card className={styles.ordersSection}>
        <div className={styles.ordersSectionHeader}>
          <h2 className={styles.sectionTitle}>📋 订单记录</h2>
          <Tag color={'blue' as any}>共 {ordersTotal} 条订单</Tag>
        </div>

        <Table
          columns={orderColumns}
          dataSource={orders}
          rowKey="id"
          loading={ordersLoading}
          emptyText={orders.length === 0 ? '暂无订单记录' : undefined}
        />

        {ordersTotal > ordersPageSize && (
          <div className={styles.paginationWrapper}>
            <Pagination
              current={ordersPage}
              total={ordersTotal}
              pageSize={ordersPageSize}
              onChange={setOrdersPage}
            />
          </div>
        )}
      </Card>
    </div>
  );
};
