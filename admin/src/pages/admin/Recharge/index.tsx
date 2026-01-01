/**
 * 充值管理主页面
 * 包含充值选项管理和充值记录管理两个子页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import { Tabs, Card, Statistic, Row, Col, App, Spin, theme } from 'antd';
import {
    WalletOutlined,
    TransactionOutlined,
    DollarOutlined,
    CheckCircleOutlined,
    CloseCircleOutlined,
    SyncOutlined,
} from '@ant-design/icons';
import { rechargeApi, type RechargeStats } from '@/api/recharge';
import { PageContainer } from '@/components';
import RechargeOptions from './Options';
import RechargeRecords from './Records';

const { TabPane } = Tabs;

/**
 * 充值管理主页面
 */
const RechargePage: React.FC = () => {
    const { message } = App.useApp();
    const { token } = theme.useToken();
    const [activeTab, setActiveTab] = useState('options');
    const [loading, setLoading] = useState(false);
    const [stats, setStats] = useState<RechargeStats | null>(null);

    /**
     * 加载统计数据
     */
    const loadStats = useCallback(async () => {
        setLoading(true);
        try {
            const response = await rechargeApi.getRechargeStats();
            if (response.data.success) {
                setStats(response.data.data);
            } else {
                message.error(response.data.message || '加载统计数据失败');
            }
        } catch (error) {
            console.error('Load stats error:', error);
            message.error('加载统计数据失败');
        } finally {
            setLoading(false);
        }
    }, [message]);

    useEffect(() => {
        loadStats();
    }, [loadStats]);

    return (
        <PageContainer title="充值管理" subTitle="管理充值档位和充值记录">
            {/* 统计卡片 */}
            <Spin spinning={loading && !stats}>
                <Row gutter={16} style={{ marginBottom: 24 }}>
                    <Col xs={24} sm={12} md={6}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="总充值订单"
                                value={stats?.totalOrders || 0}
                                prefix={<TransactionOutlined />}
                                valueStyle={{ color: token.colorPrimary }}
                            />
                        </Card>
                    </Col>
                    <Col xs={24} sm={12} md={6}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="总充值金额"
                                value={(stats?.totalAmountCents || 0) / 100}
                                precision={2}
                                prefix="¥"
                                valueStyle={{ color: token.colorSuccess }}
                            />
                        </Card>
                    </Col>
                    <Col xs={24} sm={12} md={6}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="成功订单"
                                value={stats?.paidOrders || 0}
                                prefix={<CheckCircleOutlined />}
                                valueStyle={{ color: token.colorSuccess }}
                            />
                        </Card>
                    </Col>
                    <Col xs={24} sm={12} md={6}>
                        <Card style={{ minHeight: 120 }}>
                            <Statistic
                                title="失败订单"
                                value={stats?.failedOrders || 0}
                                prefix={<CloseCircleOutlined />}
                                valueStyle={{ color: token.colorError }}
                            />
                        </Card>
                    </Col>
                </Row>
            </Spin>

            {/* 今日/本月统计 */}
            <Row gutter={16} style={{ marginBottom: 24 }}>
                <Col xs={24} sm={12} md={8}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="今日充值订单"
                            value={stats?.todayOrders || 0}
                            prefix={<SyncOutlined />}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} md={8}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="今日充值金额"
                            value={(stats?.todayAmountCents || 0) / 100}
                            precision={2}
                            prefix="¥"
                            valueStyle={{ color: token.colorWarning }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} md={8}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="本月充值金额"
                            value={(stats?.monthAmountCents || 0) / 100}
                            precision={2}
                            prefix="¥"
                            valueStyle={{ color: token.colorInfo }}
                        />
                    </Card>
                </Col>
            </Row>

            {/* Tabs */}
            <Card>
                <Tabs
                    activeKey={activeTab}
                    onChange={(key) => {
                        setActiveTab(key);
                        // 切换tab时刷新统计数据
                        loadStats();
                    }}
                    items={[
                        {
                            key: 'options',
                            label: (
                                <span>
                                    <WalletOutlined />
                                    充值档位
                                </span>
                            ),
                            children: <RechargeOptions onStatsUpdate={loadStats} />,
                        },
                        {
                            key: 'records',
                            label: (
                                <span>
                                    <DollarOutlined />
                                    充值记录
                                </span>
                            ),
                            children: <RechargeRecords onStatsUpdate={loadStats} />,
                        },
                    ]}
                />
            </Card>
        </PageContainer>
    );
};

export default RechargePage;
