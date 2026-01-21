/**
 * 团队成员卡片
 */
import React from 'react';
import {
    Card,
    Avatar,
    Tag,
    Space,
    Button,
    Popconfirm,
    Typography,
    Row,
    Col,
    Statistic,
} from 'antd';
import {
    UserOutlined,
    CrownOutlined,
    TeamOutlined,
    DeleteOutlined,
    SwapOutlined,
} from '@ant-design/icons';
import type { TeamMember } from '@/api/team';
import dayjs from 'dayjs';

const { Text } = Typography;

export interface MemberCardProps {
    /** 成员数据 */
    member: TeamMember;
    /** 是否可转让队长（仅队长可用） */
    canTransfer?: boolean;
    /** 移除成员 */
    onRemove?: (member: TeamMember) => void;
    /** 转让队长 */
    onTransfer?: (member: TeamMember) => void;
}

/**
 * 状态映射
 */
const statusMap: Record<TeamMember['status'], { color: string; text: string }> = {
    active: { color: 'success', text: '正常' },
    left: { color: 'default', text: '已退出' },
    kicked: { color: 'error', text: '已移除' },
};

/**
 * 角色映射
 */
const roleMap: Record<TeamMember['role'], { icon: React.ReactNode; text: string }> = {
    leader: { icon: <CrownOutlined />, text: '队长' },
    member: { icon: <TeamOutlined />, text: '成员' },
};

/**
 * 成员卡片组件
 * 优化: 使用 React.memo 避免不必要的重新渲染
 * 适用场景: 列表项组件，在成员列表中频繁渲染，仅在成员数据变化时重新渲染
 */
const MemberCard: React.FC<MemberCardProps> = React.memo(({
    member,
    canTransfer = false,
    onRemove,
    onTransfer,
}) => {
    const isActive = member.status === 'active';
    const isLeader = member.role === 'leader';

    return (
        <Card
            size="small"
            className={isActive ? '' : 'member-card-inactive'}
            style={{ marginBottom: 12 }}
        >
            {/* 头部：头像和基本信息 */}
            <Space style={{ width: '100%', justifyContent: 'space-between' }}>
                <Space>
                    <Avatar
                        size={48}
                        src={member.player?.avatar || undefined}
                        icon={<UserOutlined />}
                    />
                    <div>
                        <div style={{ fontWeight: 500, fontSize: 15 }}>
                            <Space size={4}>
                                {member.player?.nickname || '-'}
                                {roleMap[member.role].icon}
                            </Space>
                        </div>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            ID: {member.playerId}
                        </Text>
                        {member.player?.rank && (
                            <Tag color="blue" style={{ marginLeft: 8 }}>
                                {member.player.rank}
                            </Tag>
                        )}
                    </div>
                </Space>

                {/* 状态标签 */}
                <Space>
                    <Tag color={statusMap[member.status].color}>
                        {statusMap[member.status].text}
                    </Tag>
                    {isLeader && (
                        <Tag color="gold" icon={<CrownOutlined />}>
                            {roleMap.leader.text}
                        </Tag>
                    )}
                </Space>
            </Space>

            {/* 统计数据 */}
            {isActive && (
                <Row gutter={16} style={{ marginTop: 16 }}>
                    <Col span={12}>
                        <Statistic
                            title="订单数"
                            value={member.orderCount}
                            suffix="单"
                            />
                    </Col>
                    <Col span={12}>
                        <Statistic
                            title="收益"
                            value={member.incomeCents ? (member.incomeCents / 100).toFixed(2) : 0}
                            prefix="¥"
                            />
                    </Col>
                </Row>
            )}

            {/* 加入时间和排序 */}
            <div style={{ marginTop: 12, fontSize: 12, color: '#999' }}>
                <Space split="|">
                    <span>排序: {member.sortOrder}</span>
                    <span>加入: {member.joinedAt ? dayjs(member.joinedAt).format('YYYY-MM-DD') : '-'}</span>
                    {member.leftAt && <span>退出: {dayjs(member.leftAt).format('YYYY-MM-DD')}</span>}
                </Space>
            </div>

            {/* 操作按钮 */}
            {isActive && !isLeader && (
                <div style={{ marginTop: 12 }}>
                    <Space>
                        {canTransfer && onTransfer && (
                            <Popconfirm
                                title="确定要将队长转让给该成员吗？"
                                description="转让后您将成为普通成员"
                                onConfirm={() => onTransfer(member)}
                            >
                                <Button
                                    type="link"
                                    size="small"
                                    icon={<SwapOutlined />}
                                >
                                    转让队长
                                </Button>
                            </Popconfirm>
                        )}
                        {onRemove && (
                            <Popconfirm
                                title="确定要移除该成员吗？"
                                onConfirm={() => onRemove(member)}
                            >
                                <Button
                                    type="link"
                                    size="small"
                                    danger
                                    icon={<DeleteOutlined />}
                                >
                                    移除
                                </Button>
                            </Popconfirm>
                        )}
                    </Space>
                </div>
            )}
        </Card>
    );
});

export default MemberCard;
