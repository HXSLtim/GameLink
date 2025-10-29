import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Card, Button, Tag, Modal, Input } from '../../components';
import { paymentApi } from '../../services/api/payment';
import type { PaymentDetail, RefundPaymentRequest } from '../../types/payment';
import { PAYMENT_STATUS_TEXT, PAYMENT_METHOD_TEXT, PAYMENT_METHOD_ICON } from '../../types/payment';
import { formatCurrency, formatDateTime, formatRelativeTime } from '../../utils/formatters';
import styles from './PaymentDetail.module.less';

export const PaymentDetailPage: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [loading, setLoading] = useState(false);
  const [payment, setPayment] = useState<PaymentDetail | null>(null);
  
  // 退款Modal
  const [refundModalVisible, setRefundModalVisible] = useState(false);
  const [refundReason, setRefundReason] = useState('');
  const [refunding, setRefunding] = useState(false);
  
  // 删除Modal
  const [deleteModalVisible, setDeleteModalVisible] = useState(false);

  // 加载支付详情
  const loadPaymentDetail = async () => {
    if (!id) return;

    setLoading(true);
    try {
      const data = await paymentApi.getDetail(Number(id));
      setPayment(data);
    } catch (err) {
      console.error('加载支付详情失败:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPaymentDetail();
  }, [id]);

  // 申请退款
  const handleRefund = async () => {
    if (!id || !refundReason.trim()) {
      alert('请输入退款原因');
      return;
    }

    setRefunding(true);
    try {
      const refundData: RefundPaymentRequest = {
        reason: refundReason,
      };
      await paymentApi.refund(Number(id), refundData);
      setRefundModalVisible(false);
      setRefundReason('');
      await loadPaymentDetail();
    } catch (err) {
      console.error('退款失败:', err);
      alert('退款失败，请稍后重试');
    } finally {
      setRefunding(false);
    }
  };

  // 确认收款
  const handleCapture = async () => {
    if (!id) return;
    
    if (!confirm('确定要确认收款吗？')) return;

    try {
      await paymentApi.capture(Number(id));
      await loadPaymentDetail();
    } catch (err) {
      console.error('确认收款失败:', err);
      alert('确认收款失败');
    }
  };

  // 删除支付记录
  const handleDelete = async () => {
    if (!id) return;

    try {
      await paymentApi.delete(Number(id));
      setDeleteModalVisible(false);
      navigate('/payments');
    } catch (err) {
      console.error('删除失败:', err);
      alert('删除失败');
    }
  };

  // 获取状态颜色
  const getStatusColor = (status: string): string => {
    const colorMap: Record<string, string> = {
      pending: 'orange',
      paid: 'green',
      failed: 'red',
      refunded: 'purple',
      cancelled: 'default',
    };
    return colorMap[status] || 'default';
  };

  if (loading) {
    return (
      <div className={styles.container}>
        <div className={styles.loading}>加载中...</div>
      </div>
    );
  }

  if (!payment) {
    return (
      <div className={styles.container}>
        <div className={styles.error}>支付记录不存在</div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* 头部 */}
      <div className={styles.header}>
        <Button variant="text" onClick={() => navigate('/payments')} className={styles.backButton}>
          ← 返回列表
        </Button>
        <div className={styles.headerActions}>
          {payment.status === 'pending' && (
            <Button variant="primary" onClick={handleCapture}>
              确认收款
            </Button>
          )}
          {payment.status === 'paid' && (
            <Button variant="outlined" onClick={() => setRefundModalVisible(true)}>
              申请退款
            </Button>
          )}
          <Button variant="outlined" onClick={() => setDeleteModalVisible(true)}>
            删除
          </Button>
        </div>
      </div>

      {/* 支付信息卡片 */}
      <Card className={styles.infoCard}>
        <div className={styles.paymentHeader}>
          <div className={styles.paymentIcon}>
            {React.createElement(PAYMENT_METHOD_ICON[payment.method], { size: 64 })}
          </div>
          <div className={styles.paymentInfo}>
            <div className={styles.paymentAmount}>
              {formatCurrency(payment.amountCents)}
            </div>
            <div className={styles.paymentMeta}>
              <Tag color={getStatusColor(payment.status) as any}>
                {PAYMENT_STATUS_TEXT[payment.status]}
              </Tag>
              <span className={styles.paymentMethod}>
                {PAYMENT_METHOD_TEXT[payment.method]}
              </span>
              <span className={styles.paymentId}>
                ID: {payment.id}
              </span>
            </div>
          </div>
        </div>

        <div className={styles.detailsGrid}>
          <div className={styles.detailItem}>
            <span className={styles.detailLabel}>支付方式</span>
            <span className={styles.detailValue}>
              {React.createElement(PAYMENT_METHOD_ICON[payment.method], { size: 20 })}{' '}
              {PAYMENT_METHOD_TEXT[payment.method]}
            </span>
          </div>

          <div className={styles.detailItem}>
            <span className={styles.detailLabel}>支付金额</span>
            <span className={styles.detailValue}>
              {formatCurrency(payment.amountCents)}
            </span>
          </div>

          {payment.currency && (
            <div className={styles.detailItem}>
              <span className={styles.detailLabel}>货币类型</span>
              <span className={styles.detailValue}>{payment.currency}</span>
            </div>
          )}

          <div className={styles.detailItem}>
            <span className={styles.detailLabel}>支付状态</span>
            <span className={styles.detailValue}>
              <Tag color={getStatusColor(payment.status) as any}>
                {PAYMENT_STATUS_TEXT[payment.status]}
              </Tag>
            </span>
          </div>

          {payment.providerTradeNo && (
            <div className={styles.detailItem}>
              <span className={styles.detailLabel}>交易流水号</span>
              <span className={styles.detailValue}>{payment.providerTradeNo}</span>
            </div>
          )}

          {payment.paidAt && (
            <div className={styles.detailItem}>
              <span className={styles.detailLabel}>支付时间</span>
              <span className={styles.detailValue}>
                {formatDateTime(payment.paidAt)}
              </span>
            </div>
          )}

          {payment.refundedAt && (
            <div className={styles.detailItem}>
              <span className={styles.detailLabel}>退款时间</span>
              <span className={styles.detailValue}>
                {formatDateTime(payment.refundedAt)}
              </span>
            </div>
          )}
        </div>

        <div className={styles.timestamps}>
          <div className={styles.timestamp}>
            <span className={styles.timestampLabel}>创建时间:</span>
            <span className={styles.timestampValue}>
              {formatDateTime(payment.createdAt)}
              <span className={styles.relative}>({formatRelativeTime(payment.createdAt)})</span>
            </span>
          </div>
          <div className={styles.timestamp}>
            <span className={styles.timestampLabel}>更新时间:</span>
            <span className={styles.timestampValue}>
              {formatDateTime(payment.updatedAt)}
              <span className={styles.relative}>({formatRelativeTime(payment.updatedAt)})</span>
            </span>
          </div>
        </div>
      </Card>

      {/* 关联订单信息 */}
      {payment.order && (
        <Card className={styles.orderCard}>
          <h3 className={styles.cardTitle}>关联订单</h3>
          <div className={styles.orderInfo}>
            <div className={styles.orderItem}>
              <span className={styles.orderLabel}>订单ID:</span>
              <Button 
                variant="text" 
                onClick={() => navigate(`/orders/${payment.order!.id}`)}
                className={styles.orderLink}
              >
                #{payment.order.id}
              </Button>
            </div>
            {payment.order.title && (
              <div className={styles.orderItem}>
                <span className={styles.orderLabel}>订单标题:</span>
                <span className={styles.orderValue}>{payment.order.title}</span>
              </div>
            )}
            {payment.order.status && (
              <div className={styles.orderItem}>
                <span className={styles.orderLabel}>订单状态:</span>
                <span className={styles.orderValue}>{payment.order.status}</span>
              </div>
            )}
          </div>
        </Card>
      )}

      {/* 用户信息 */}
      {payment.user && (
        <Card className={styles.userCard}>
          <h3 className={styles.cardTitle}>支付用户</h3>
          <div className={styles.userInfo}>
            <div className={styles.userItem}>
              <span className={styles.userLabel}>用户ID:</span>
              <Button 
                variant="text" 
                onClick={() => navigate(`/users/${payment.user!.id}`)}
                className={styles.userLink}
              >
                #{payment.user.id}
              </Button>
            </div>
            <div className={styles.userItem}>
              <span className={styles.userLabel}>用户名:</span>
              <span className={styles.userValue}>{payment.user.name}</span>
            </div>
            {payment.user.phone && (
              <div className={styles.userItem}>
                <span className={styles.userLabel}>手机号:</span>
                <span className={styles.userValue}>{payment.user.phone}</span>
              </div>
            )}
            {payment.user.email && (
              <div className={styles.userItem}>
                <span className={styles.userLabel}>邮箱:</span>
                <span className={styles.userValue}>{payment.user.email}</span>
              </div>
            )}
          </div>
        </Card>
      )}

      {/* 退款信息 */}
      {payment.refundInfo && (
        <Card className={styles.refundCard}>
          <h3 className={styles.cardTitle}>退款信息</h3>
          <div className={styles.refundInfo}>
            <div className={styles.refundItem}>
              <span className={styles.refundLabel}>退款金额:</span>
              <span className={styles.refundValue}>
                {formatCurrency(payment.refundInfo.refundAmount)}
              </span>
            </div>
            <div className={styles.refundItem}>
              <span className={styles.refundLabel}>退款原因:</span>
              <span className={styles.refundValue}>{payment.refundInfo.refundReason}</span>
            </div>
            <div className={styles.refundItem}>
              <span className={styles.refundLabel}>退款时间:</span>
              <span className={styles.refundValue}>
                {formatDateTime(payment.refundInfo.refundedAt)}
              </span>
            </div>
          </div>
        </Card>
      )}

      {/* 第三方数据 */}
      {payment.providerRaw && (
        <Card className={styles.providerCard}>
          <h3 className={styles.cardTitle}>第三方支付数据</h3>
          <pre className={styles.providerRaw}>
            {JSON.stringify(payment.providerRaw, null, 2)}
          </pre>
        </Card>
      )}

      {/* 退款Modal */}
      <Modal
        visible={refundModalVisible}
        title="申请退款"
        onClose={() => {
          setRefundModalVisible(false);
          setRefundReason('');
        }}
        onOk={handleRefund}
        onCancel={() => {
          setRefundModalVisible(false);
          setRefundReason('');
        }}
        okText={refunding ? '退款中...' : '确定退款'}
        cancelText="取消"
        width={500}
      >
        <div className={styles.refundForm}>
          <p className={styles.refundTip}>
            退款金额：<strong>{formatCurrency(payment.amountCents)}</strong>
          </p>
          <div className={styles.formItem}>
            <label className={styles.formLabel}>
              退款原因 <span className={styles.required}>*</span>
            </label>
            <textarea
              className={styles.formTextarea}
              value={refundReason}
              onChange={(e) => setRefundReason(e.target.value)}
              placeholder="请输入退款原因"
              rows={4}
            />
          </div>
          <p className={styles.refundWarning}>
            ⚠️ 退款操作不可撤销，请谨慎操作
          </p>
        </div>
      </Modal>

      {/* 删除确认Modal */}
      <Modal
        visible={deleteModalVisible}
        title="确认删除"
        onClose={() => setDeleteModalVisible(false)}
        onOk={handleDelete}
        onCancel={() => setDeleteModalVisible(false)}
        okText="确定删除"
        cancelText="取消"
        width={400}
      >
        <p>确定要删除这条支付记录吗？此操作不可恢复。</p>
        {payment.status === 'paid' && (
          <p className={styles.deleteWarning}>
            ⚠️ 该支付记录状态为"已支付"，删除可能影响订单和财务数据。
          </p>
        )}
      </Modal>
    </div>
  );
};

