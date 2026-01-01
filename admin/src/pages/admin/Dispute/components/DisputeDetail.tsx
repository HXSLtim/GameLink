/**
 * DisputeDetail Component
 * Displays detailed information about a dispute
 */
import React from 'react';
import { Tag, Descriptions, Card, Timeline, Divider, Image, Space, Button, Typography } from 'antd';
import {
    DISPUTE_STATUS_LABELS,
    DISPUTE_STATUS_COLORS,
    DISPUTE_TYPE_LABELS,
    DISPUTE_RESOLUTION_LABELS,
} from '@/types/dispute';
import type { Dispute } from '@/types/dispute';
import dayjs from 'dayjs';

export interface DisputeDetailProps {
    /** Dispute data */
    dispute: Dispute;
}

/**
 * DisputeDetail Component
 */
export const DisputeDetail: React.FC<DisputeDetailProps> = ({ dispute }) => {
    /**
     * Render evidence images
     */
    const renderEvidence = () => {
        if (!dispute.evidenceUrls || dispute.evidenceUrls.length === 0) {
            return null;
        }

        return (
            <Card size="small" title="证据图片" style={{ marginTop: 16 }}>
                <Image.PreviewGroup>
                    <Space wrap>
                        {dispute.evidenceUrls.map((url, index) => (
                            <Image
                                key={index}
                                src={url}
                                width={100}
                                height={100}
                                style={{ objectFit: 'cover' }}
                            />
                        ))}
                    </Space>
                </Image.PreviewGroup>
            </Card>
        );
    };

    return (
        <div>
            {/* Status Card */}
            <Card size="small" style={{ marginBottom: 16 }}>
                <div style={{ textAlign: 'center' }}>
                    <Tag
                        color={DISPUTE_STATUS_COLORS[dispute.status]}
                        style={{ fontSize: 16, padding: '4px 16px' }}
                    >
                        {DISPUTE_STATUS_LABELS[dispute.status]}
                    </Tag>
                    <div style={{ marginTop: 12 }}>
                        <Space>
                            <span>订单号:</span>
                            <Typography.Text copyable>{dispute.orderNo || `ID: ${dispute.orderId}`}</Typography.Text>
                        </Space>
                    </div>
                </div>
            </Card>

            {/* Basic Information */}
            <Descriptions title="基本信息" column={2} size="small" bordered>
                <Descriptions.Item label="纠纷ID">{dispute.id}</Descriptions.Item>
                <Descriptions.Item label="订单ID">{dispute.orderId}</Descriptions.Item>
                <Descriptions.Item label="订单号" span={2}>
                    <Typography.Text copyable>{dispute.orderNo || '-'}</Typography.Text>
                </Descriptions.Item>
                <Descriptions.Item label="发起人ID">{dispute.initiatorId}</Descriptions.Item>
                <Descriptions.Item label="发起人类型">
                    <Tag color={dispute.initiatorType === 'user' ? 'blue' : 'purple'}>
                        {dispute.initiatorType === 'user' ? '用户' : '陪玩师'}
                    </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="发起人姓名" span={2}>
                    {dispute.initiatorName || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="纠纷类型">
                    {DISPUTE_TYPE_LABELS[dispute.type]}
                </Descriptions.Item>
                <Descriptions.Item label="状态">
                    <Tag color={DISPUTE_STATUS_COLORS[dispute.status]}>
                        {DISPUTE_STATUS_LABELS[dispute.status]}
                    </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="纠纷原因" span={2}>
                    {dispute.reason}
                </Descriptions.Item>
                {dispute.evidenceText && (
                    <Descriptions.Item label="详细描述" span={2}>
                        {dispute.evidenceText}
                    </Descriptions.Item>
                )}
                {dispute.chatSnapshotId && (
                    <Descriptions.Item label="聊天快照ID" span={2}>
                        <Typography.Text copyable>{dispute.chatSnapshotId}</Typography.Text>
                    </Descriptions.Item>
                )}
                <Descriptions.Item label="创建时间">
                    {dayjs(dispute.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                </Descriptions.Item>
                <Descriptions.Item label="更新时间">
                    {dayjs(dispute.updatedAt).format('YYYY-MM-DD HH:mm:ss')}
                </Descriptions.Item>
                <Descriptions.Item label="TraceID" span={2}>
                    <Typography.Text copyable>{dispute.traceId}</Typography.Text>
                </Descriptions.Item>
            </Descriptions>

            {/* Evidence Images */}
            {renderEvidence()}

            {/* Dual-CS Information */}
            <Descriptions
                title="双客服信息"
                column={2}
                size="small"
                bordered
                style={{ marginTop: 16 }}
            >
                <Descriptions.Item label="主客服ID">
                    {dispute.assignedServiceId || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="原客服ID">
                    {dispute.originalServiceId || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="主客服姓名">
                    {dispute.assignedServiceName || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="原客服姓名">
                    {dispute.originalServiceName || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="SLA截止时间">
                    {dispute.slaDeadline
                        ? dayjs(dispute.slaDeadline).format('YYYY-MM-DD HH:mm:ss')
                        : '-'}
                </Descriptions.Item>
                <Descriptions.Item label="SLA状态">
                    {dispute.slaBreached ? (
                        <Tag color="error">已超时</Tag>
                    ) : (
                        <Tag color="success">正常</Tag>
                    )}
                    {dispute.slaBreachedAt && (
                        <span style={{ marginLeft: 8, fontSize: 12, color: '#999' }}>
                            {dayjs(dispute.slaBreachedAt).format('YYYY-MM-DD HH:mm:ss')}
                        </span>
                    )}
                </Descriptions.Item>
            </Descriptions>

            {/* Resolution Information */}
            {dispute.resolution && (
                <Descriptions
                    title="处理结果"
                    column={2}
                    size="small"
                    bordered
                    style={{ marginTop: 16 }}
                >
                    <Descriptions.Item label="解决方案">
                        {DISPUTE_RESOLUTION_LABELS[dispute.resolution]}
                    </Descriptions.Item>
                    <Descriptions.Item label="处理人ID">
                        {dispute.resolvedBy || '-'}
                    </Descriptions.Item>
                    <Descriptions.Item label="处理人姓名">
                        {dispute.resolvedByName || '-'}
                    </Descriptions.Item>
                    <Descriptions.Item label="处理时间">
                        {dispute.resolvedAt
                            ? dayjs(dispute.resolvedAt).format('YYYY-MM-DD HH:mm:ss')
                            : '-'}
                    </Descriptions.Item>
                    <Descriptions.Item label="处理备注" span={2}>
                        {dispute.resolveRemark || '-'}
                    </Descriptions.Item>
                </Descriptions>
            )}

            {/* Rollback Information */}
            {dispute.rolledBackAt && (
                <Descriptions
                    title="回滚信息"
                    column={2}
                    size="small"
                    bordered
                    style={{ marginTop: 16 }}
                >
                    <Descriptions.Item label="回滚时间">
                        {dayjs(dispute.rolledBackAt).format('YYYY-MM-DD HH:mm:ss')}
                    </Descriptions.Item>
                    <Descriptions.Item label="回滚操作人ID">
                        {dispute.rolledBackByUserId || '-'}
                    </Descriptions.Item>
                    <Descriptions.Item label="回滚原因" span={2}>
                        {dispute.rollbackReason || '-'}
                    </Descriptions.Item>
                </Descriptions>
            )}

            {/* Timeline */}
            <Card title="处理进度" size="small" style={{ marginTop: 16 }}>
                <Timeline
                    items={[
                        {
                            color: 'green',
                            children: `${dayjs(dispute.createdAt).format('MM-DD HH:mm:ss')} 纠纷创建`,
                        },
                        ...(dispute.assignedServiceId
                            ? [
                                  {
                                      color: 'blue',
                                      children: `分配给 ${dispute.assignedServiceName || '客服'}`,
                                  },
                              ]
                            : []),
                        ...(dispute.resolvedAt
                            ? [
                                  {
                                      color: 'green',
                                      children: `${dayjs(dispute.resolvedAt).format(
                                          'MM-DD HH:mm:ss'
                                      )} 处理完成`,
                                  },
                              ]
                            : []),
                        ...(dispute.rolledBackAt
                            ? [
                                  {
                                      color: 'orange',
                                      children: `${dayjs(dispute.rolledBackAt).format(
                                          'MM-DD HH:mm:ss'
                                      )} 分配回滚`,
                                  },
                              ]
                            : []),
                    ]}
                />
            </Card>

            {/* Action Buttons for Quick Actions */}
            <Divider style={{ margin: '16px 0' }} />
            <Space>
                {dispute.status === 'pending' && (
                    <Button type="primary">分配纠纷</Button>
                )}
                {['assigned', 'mediating'].includes(dispute.status) && (
                    <>
                        <Button type="primary">处理纠纷</Button>
                        <Button>回滚分配</Button>
                    </>
                )}
            </Space>
        </div>
    );
};

export default DisputeDetail;
