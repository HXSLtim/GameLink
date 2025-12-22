/**
 * 用户端钱包页面
 * 余额查看、充值、交易记录
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Row,
    Col,
    Statistic,
    Button,
    Table,
    Tag,
    Modal,
    InputNumber,
    Radio,
    Space,
    Typography,
    message,
    Tabs,
} from 'antd';
import {
    WalletOutlined,
    PlusOutlined,
    HistoryOutlined,
    CreditCardOutlined,
    AlipayOutlined,
    WechatOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

const { Title, Text } = Typography;

interface WalletInfo {
    balance: number;
    frozenAmount: number;
    totalRecharge: number;
    totalSpent: number;
}

interface Transaction {
    id: number;
    type: 'recharge' | 'consume' | 'refund' | 'withdraw';
    amount: number;
    balance: number;
    description: string;
    createdAt: string;
    status: 'success' | 'pending' | 'failed';
}

const UserWallet: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [walletInfo, setWalletInfo] = useState<WalletInfo>({
        balance: 0,
        frozenAmount: 0,
        totalRecharge: 0,
        totalSpent: 0,
    });
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [rechargeVisible, setRechargeVisible] = useState(false);
    const [rechargeAmount, setRechargeAmount] = useState<number>(100);
    const [paymentMethod, setPaymentMethod] = useState<string>('alipay');
    const [rechargeLoading, setRechargeLoading] = useState(false);

    const loadWalletInfo = useCallback(async () => {
        setLoading(true);
        try {
            await new Promise(resolve => setTimeout(resolve, 500));
            setWalletInfo({
                balance: 256.50,
                frozenAmount: 50.00,
                totalRecharge: 1000.00,
                totalSpent: 693.50,
            });
        } catch {
            message.error('加载钱包信息失败');
        } finally {
            setLoading(false);
        }
    }, []);

    const loadTransactions = useCallback(async () => {
        try {
            await new Promise(resolve => setTimeout(resolve, 300));
            setTransactions([
                { id: 1, type: 'recharge', amount: 100, balance: 256.50, description: '支付宝充值', createdAt: '2024-12-16 10:30:00', status: 'success' },
                { id: 2, type: 'consume', amount: -50, balance: 156.50, description: '订单支付 ORD202412160001', createdAt: '2024-12-15 14:20:00', status: 'success' },
                { id: 3, type: 'refund', amount: 30, balance: 206.50, description: '订单退款 ORD202412140005', createdAt: '2024-12-14 16:00:00', status: 'success' },
                { id: 4, type: 'consume', amount: -80, balance: 176.50, description: '订单支付 ORD202412130002', createdAt: '2024-12-13 20:15:00', status: 'success' },
                { id: 5, type: 'recharge', amount: 200, balance: 256.50, description: '微信充值', createdAt: '2024-12-12 09:00:00', status: 'success' },
            ]);
        } catch {
            message.error('加载交易记录失败');
        }
    }, []);

    useEffect(() => {
        loadWalletInfo();
        loadTransactions();
    }, [loadWalletInfo, loadTransactions]);

    const handleRecharge = async () => {
        if (!rechargeAmount || rechargeAmount < 1) {
            message.warning('请输入有效的充值金额');
            return;
        }
        setRechargeLoading(true);
        try {
            await new Promise(resolve => setTimeout(resolve, 1000));
            message.success(`充值 ¥${rechargeAmount} 成功`);
            setRechargeVisible(false);
            loadWalletInfo();
            loadTransactions();
        } catch {
            message.error('充值失败');
        } finally {
            setRechargeLoading(false);
        }
    };

    const getTypeTag = (type: Transaction['type']) => {
        const config = {
            recharge: { color: 'green', text: '充值' },
            consume: { color: 'red', text: '消费' },
            refund: { color: 'blue', text: '退款' },
            withdraw: { color: 'orange', text: '提现' },
        };
        return <Tag color={config[type].color}>{config[type].text}</Tag>;
    };

    const columns: ColumnsType<Transaction> = [
        { title: '交易类型', dataIndex: 'type', key: 'type', render: (type) => getTypeTag(type) },
        { title: '金额', dataIndex: 'amount', key: 'amount', render: (amount) => (
            <Text style={{ color: amount > 0 ? '#52c41a' : '#ff4d4f', fontWeight: 500 }}>
                {amount > 0 ? '+' : ''}{amount.toFixed(2)}
            </Text>
        )},
        { title: '余额', dataIndex: 'balance', key: 'balance', render: (balance) => `¥${balance.toFixed(2)}` },
        { title: '描述', dataIndex: 'description', key: 'description' },
        { title: '时间', dataIndex: 'createdAt', key: 'createdAt' },
        { title: '状态', dataIndex: 'status', key: 'status', render: (status) => (
            <Tag color={status === 'success' ? 'green' : status === 'pending' ? 'orange' : 'red'}>
                {status === 'success' ? '成功' : status === 'pending' ? '处理中' : '失败'}
            </Tag>
        )},
    ];

    const quickAmounts = [50, 100, 200, 500, 1000];

    return (
        <div style={{ padding: 24 }}>
            <Title level={4}><WalletOutlined /> 我的钱包</Title>

            {/* 余额卡片 */}
            <Card style={{ marginBottom: 16, background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}>
                <Row gutter={24} align="middle">
                    <Col flex="auto">
                        <Text style={{ color: 'rgba(255,255,255,0.8)' }}>账户余额</Text>
                        <Title level={2} style={{ color: '#fff', margin: '8px 0' }}>
                            ¥ {walletInfo.balance.toFixed(2)}
                        </Title>
                        <Text style={{ color: 'rgba(255,255,255,0.6)' }}>
                            冻结金额: ¥{walletInfo.frozenAmount.toFixed(2)}
                        </Text>
                    </Col>
                    <Col>
                        <Button type="primary" size="large" icon={<PlusOutlined />} onClick={() => setRechargeVisible(true)}
                            style={{ background: '#fff', color: '#667eea', border: 'none' }}>
                            充值
                        </Button>
                    </Col>
                </Row>
            </Card>

            {/* 统计数据 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic title="累计充值" value={walletInfo.totalRecharge} prefix="¥" precision={2} />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic title="累计消费" value={walletInfo.totalSpent} prefix="¥" precision={2} />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic title="可用余额" value={walletInfo.balance} prefix="¥" precision={2} valueStyle={{ color: '#3f8600' }} />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic title="冻结金额" value={walletInfo.frozenAmount} prefix="¥" precision={2} valueStyle={{ color: '#faad14' }} />
                    </Card>
                </Col>
            </Row>

            {/* 交易记录 */}
            <Card title={<><HistoryOutlined /> 交易记录</>}>
                <Tabs defaultActiveKey="all" items={[
                    { key: 'all', label: '全部', children: <Table columns={columns} dataSource={transactions} rowKey="id" pagination={{ pageSize: 10 }} /> },
                    { key: 'recharge', label: '充值', children: <Table columns={columns} dataSource={transactions.filter(t => t.type === 'recharge')} rowKey="id" /> },
                    { key: 'consume', label: '消费', children: <Table columns={columns} dataSource={transactions.filter(t => t.type === 'consume')} rowKey="id" /> },
                    { key: 'refund', label: '退款', children: <Table columns={columns} dataSource={transactions.filter(t => t.type === 'refund')} rowKey="id" /> },
                ]} />
            </Card>

            {/* 充值弹窗 */}
            <Modal title="账户充值" open={rechargeVisible} onCancel={() => setRechargeVisible(false)} onOk={handleRecharge}
                confirmLoading={rechargeLoading} okText="确认充值" width={480}>
                <div style={{ padding: '16px 0' }}>
                    <div style={{ marginBottom: 24 }}>
                        <Text strong>选择充值金额</Text>
                        <div style={{ marginTop: 12 }}>
                            <Space wrap>
                                {quickAmounts.map(amount => (
                                    <Button key={amount} type={rechargeAmount === amount ? 'primary' : 'default'}
                                        onClick={() => setRechargeAmount(amount)}>¥{amount}</Button>
                                ))}
                            </Space>
                        </div>
                        <div style={{ marginTop: 12 }}>
                            <Text type="secondary">自定义金额：</Text>
                            <InputNumber min={1} max={10000} value={rechargeAmount} onChange={(v) => setRechargeAmount(v || 0)}
                                prefix="¥" style={{ width: 150, marginLeft: 8 }} />
                        </div>
                    </div>
                    <div>
                        <Text strong>选择支付方式</Text>
                        <Radio.Group value={paymentMethod} onChange={(e) => setPaymentMethod(e.target.value)} style={{ marginTop: 12, display: 'block' }}>
                            <Space direction="vertical" style={{ width: '100%' }}>
                                <Radio.Button value="alipay" style={{ width: '100%', height: 48, lineHeight: '48px' }}>
                                    <AlipayOutlined style={{ color: '#1677ff', marginRight: 8 }} /> 支付宝
                                </Radio.Button>
                                <Radio.Button value="wechat" style={{ width: '100%', height: 48, lineHeight: '48px' }}>
                                    <WechatOutlined style={{ color: '#52c41a', marginRight: 8 }} /> 微信支付
                                </Radio.Button>
                                <Radio.Button value="card" style={{ width: '100%', height: 48, lineHeight: '48px' }}>
                                    <CreditCardOutlined style={{ marginRight: 8 }} /> 银行卡
                                </Radio.Button>
                            </Space>
                        </Radio.Group>
                    </div>
                </div>
            </Modal>
        </div>
    );
};

export default UserWallet;
