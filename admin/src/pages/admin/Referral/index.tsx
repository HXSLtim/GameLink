/**
 * Referral Management Page - Main Entry Point
 * 推荐管理页面 - 主入口
 *
 * Provides tabs for managing referrals, codes, and rewards.
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Tabs,
    Row,
    Col,
    Statistic,
    theme,
} from 'antd';
import {
    TeamOutlined,
    GiftOutlined,
    DollarOutlined,
    UserAddOutlined,
} from '@ant-design/icons';
import { referralApi, type ReferralStats } from '@/api/referral';
import ReferralList from './ReferralList';
import ReferralCodes from './Codes';
import ReferralRewards from './Rewards';

/**
 * Referral Management Page
 */
const ReferralPage: React.FC = () => {
    const { token } = theme.useToken();
    const [stats, setStats] = useState<ReferralStats | null>(null);
    const [statsLoading, setStatsLoading] = useState(false);
    const [activeTab, setActiveTab] = useState('1');

    /**
     * Load referral statistics
     */
    const loadStats = useCallback(async () => {
        try {
            setStatsLoading(true);
            const response = await referralApi.getReferralStats();
            if (response.data.success) {
                setStats(response.data.data);
            }
        } catch (err) {
            console.error('Failed to load referral stats:', err);
            // Silent fail, don't show error message
        } finally {
            setStatsLoading(false);
        }
    }, []);

    useEffect(() => {
        loadStats();
    }, [loadStats]);

    /**
     * Handle refresh stats when data changes
     */
    const handleDataChange = useCallback(() => {
        loadStats();
    }, [loadStats]);

    /**
     * Handle tab change
     */
    const handleTabChange = (key: string) => {
        setActiveTab(key);
    };

    return (
        <div style={{ padding: 24 }}>
            {/* Statistics Cards */}
            <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
                <Col xs={24} sm={12} lg={6}>
                    <Card variant="borderless" loading={statsLoading} style={{ minHeight: 120 }}>
                        <Statistic
                            title="总推荐数"
                            value={stats?.totalReferrals || 0}
                            prefix={<TeamOutlined />}
                            valueStyle={{ color: token.colorPrimary }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card variant="borderless" loading={statsLoading} style={{ minHeight: 120 }}>
                        <Statistic
                            title="已完成推荐"
                            value={stats?.completedReferrals || 0}
                            prefix={<UserAddOutlined />}
                            valueStyle={{ color: token.colorSuccess }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card variant="borderless" loading={statsLoading} style={{ minHeight: 120 }}>
                        <Statistic
                            title="已发放奖励"
                            value={(stats?.issuedRewardsCents || 0) / 100}
                            prefix={<DollarOutlined />}
                            precision={2}
                            suffix="元"
                            valueStyle={{ color: token.colorWarning }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card variant="borderless" loading={statsLoading} style={{ minHeight: 120 }}>
                        <Statistic
                            title="活跃邀请码"
                            value={stats?.activeCodes || 0}
                            prefix={<GiftOutlined />}
                            valueStyle={{ color: token.colorInfo }}
                        />
                        <div style={{ fontSize: 12, color: '#8c8c8c', marginTop: 4 }}>
                            总计: {stats?.totalCodes || 0} 个
                        </div>
                    </Card>
                </Col>
            </Row>

            {/* Main Content Tabs */}
            <Card>
                <Tabs
                    activeKey={activeTab}
                    onChange={handleTabChange}
                    items={[
                        {
                            key: '1',
                            label: '推荐关系',
                            children: (
                                <ReferralList onDataChange={handleDataChange} />
                            ),
                        },
                        {
                            key: '2',
                            label: '邀请码管理',
                            children: (
                                <ReferralCodes onDataChange={handleDataChange} />
                            ),
                        },
                        {
                            key: '3',
                            label: '奖励管理',
                            children: (
                                <ReferralRewards onDataChange={handleDataChange} />
                            ),
                        },
                    ]}
                />
            </Card>
        </div>
    );
};

export default ReferralPage;
