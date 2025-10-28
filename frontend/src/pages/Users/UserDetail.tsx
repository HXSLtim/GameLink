import React, { useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Button, Tag } from '../../components';
import { getMockUserDetail } from '../../services/userMockData';
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
import { formatDateTime } from '../../utils/formatters';
import { UserRole } from '../../types/user.types';
import styles from './UserDetail.module.less';

export const UserDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const userDetail = useMemo(() => {
    if (!id) return null;
    return getMockUserDetail(Number(id));
  }, [id]);

  if (!userDetail) {
    return (
      <div className={styles.container}>
        <Card className={styles.errorCard}>
          <h2>用户未找到</h2>
          <p>用户 ID: {id} 不存在</p>
          <Button onClick={() => navigate('/users')}>返回用户列表</Button>
        </Card>
      </div>
    );
  }

  const isPlayer = userDetail.role === UserRole.PLAYER && userDetail.player;

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
              {userDetail.avatar_url && (
                <img src={userDetail.avatar_url} alt={userDetail.name} className={styles.avatar} />
              )}
              {!userDetail.avatar_url && (
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
                <span className={styles.infoValue}>{formatDateTime(userDetail.created_at)}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>最后登录</span>
                <span className={styles.infoValue}>
                  {userDetail.last_login_at ? formatDateTime(userDetail.last_login_at) : '从未登录'}
                </span>
              </div>
            </div>
          </Card>

          {/* 统计信息 */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>统计数据</h2>
            <div className={styles.statsGrid}>
              <div className={styles.statItem}>
                <div className={styles.statValue}>{userDetail.order_count || 0}</div>
                <div className={styles.statLabel}>订单数量</div>
              </div>
              <div className={styles.statItem}>
                <div className={styles.statValue}>
                  {userDetail.total_spent ? formatPrice(userDetail.total_spent) : '¥0'}
                </div>
                <div className={styles.statLabel}>总消费</div>
              </div>
              <div className={styles.statItem}>
                <div className={styles.statValue}>{userDetail.review_count || 0}</div>
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
                  color={getVerificationStatusColor(userDetail.player.verification_status)}
                  className={styles.verificationTag}
                >
                  {formatVerificationStatus(userDetail.player.verification_status)}
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
                    {formatHourlyRate(userDetail.player.hourly_rate_cents)}
                  </span>
                </div>
                <div className={styles.infoItem}>
                  <span className={styles.infoLabel}>评分</span>
                  <span className={styles.infoValue}>
                    {formatRating(userDetail.player.rating_average)} (
                    {userDetail.player.rating_count}条评价)
                  </span>
                </div>
                <div className={styles.infoItem}>
                  <span className={styles.infoLabel}>认证时间</span>
                  <span className={styles.infoValue}>
                    {formatDateTime(userDetail.player.created_at)}
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
                onClick={() => console.log('查看订单')}
                className={styles.linkButton}
              >
                📋 查看订单记录
              </Button>
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
    </div>
  );
};
