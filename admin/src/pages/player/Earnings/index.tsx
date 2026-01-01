/**
 * 陪玩师端收益页面
 * 收益统计、提现申请、收益明细、图表展示
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
    Form,
    Select,
    Space,
    Typography,
    message,
    Tabs,
    Progress,
    DatePicker,
} from 'antd';
import {
    DollarOutlined,
    BankOutlined,
    HistoryOutlined,
    RiseOutlined,
    WalletOutlined,
    CalendarOutlined,
    LineChartOutlined,
    BarChartOutlined,
    TrophyOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
    LineChart,
    Line,
    AreaChart,
    Area,
    BarChart,
    Bar,
    PieChart,
    Pie,
    Cell,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer,
} from 'recharts';
import dayjs from 'dayjs';

const { Title, Text } = Typography;
const { RangePicker } = DatePicker;

interface EarningsInfo {
    availableBalance: number;
    frozenBalance: number;
    todayEarnings: number;
    weekEarnings: number;
    monthEarnings: number;
    totalEarnings: number;
    pendingWithdraw: number;
}

interface EarningsRecord {
    id: number;
    type: 'order' | 'gift' | 'bonus' | 'commission';
    amount: number;
    orderNo?: string;
    description: string;
    createdAt: string;
    status: 'settled' | 'pending' | 'frozen';
}

interface WithdrawRecord {
    id: number;
    amount: number;
    fee: number;
    actualAmount: number;
    bankName: string;
    bankAccount: string;
    status: 'pending' | 'processing' | 'success' | 'failed';
    createdAt: string;
    completedAt?: string;
}

const PlayerEarnings: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [earningsInfo, setEarningsInfo] = useState<EarningsInfo>({
        availableBalance: 0, frozenBalance: 0, todayEarnings: 0, weekEarnings: 0,
        monthEarnings: 0, totalEarnings: 0, pendingWithdraw: 0,
    });
    const [earningsRecords, setEarningsRecords] = useState<EarningsRecord[]>([]);
    const [withdrawRecords, setWithdrawRecords] = useState<WithdrawRecord[]>([]);
    const [withdrawVisible, setWithdrawVisible] = useState(false);
    const [withdrawLoading, setWithdrawLoading] = useState(false);
    const [form] = Form.useForm();
    const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs]>([
        dayjs().subtract(30, 'day'),
        dayjs(),
    ]);

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            await new Promise(resolve => setTimeout(resolve, 500));
            setEarningsInfo({
                availableBalance: 2580.50,
                frozenBalance: 320.00,
                todayEarnings: 180.00,
                weekEarnings: 1250.00,
                monthEarnings: 4680.00,
                totalEarnings: 28650.00,
                pendingWithdraw: 500.00,
            });
            setEarningsRecords([
                { id: 1, type: 'order', amount: 42.5, orderNo: 'ORD202412160010', description: '订单收益（扣除15%平台佣金）', createdAt: '2024-12-16 15:30:00', status: 'settled' },
                { id: 2, type: 'gift', amount: 28.0, description: '收到礼物：玫瑰花 x5', createdAt: '2024-12-16 14:20:00', status: 'settled' },
                { id: 3, type: 'order', amount: 68.0, orderNo: 'ORD202412160008', description: '订单收益（扣除15%平台佣金）', createdAt: '2024-12-16 12:00:00', status: 'pending' },
                { id: 4, type: 'bonus', amount: 50.0, description: '周活跃奖励', createdAt: '2024-12-15 00:00:00', status: 'settled' },
                { id: 5, type: 'order', amount: 85.0, orderNo: 'ORD202412150005', description: '订单收益（扣除15%平台佣金）', createdAt: '2024-12-15 20:30:00', status: 'settled' },
            ]);
            setWithdrawRecords([
                { id: 1, amount: 500, fee: 2, actualAmount: 498, bankName: '招商银行', bankAccount: '**** **** **** 6789', status: 'pending', createdAt: '2024-12-16 10:00:00' },
                { id: 2, amount: 1000, fee: 5, actualAmount: 995, bankName: '工商银行', bankAccount: '**** **** **** 1234', status: 'success', createdAt: '2024-12-10 14:30:00', completedAt: '2024-12-11 09:00:00' },
                { id: 3, amount: 800, fee: 4, actualAmount: 796, bankName: '招商银行', bankAccount: '**** **** **** 6789', status: 'success', createdAt: '2024-12-05 16:00:00', completedAt: '2024-12-06 10:30:00' },
            ]);
        } catch {
            message.error('加载数据失败');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => { loadData(); }, [loadData]);

    const handleWithdraw = async (values: { amount: number; bankId: number }) => {
        if (values.amount > earningsInfo.availableBalance) {
            message.error('提现金额不能超过可用余额');
            return;
        }
        setWithdrawLoading(true);
        try {
            await new Promise(resolve => setTimeout(resolve, 1000));
            message.success('提现申请已提交，预计1-3个工作日到账');
            setWithdrawVisible(false);
            form.resetFields();
            loadData();
        } catch {
            message.error('提现申请失败');
        } finally {
            setWithdrawLoading(false);
        }
    };

    const getTypeTag = (type: EarningsRecord['type']) => {
        const config = {
            order: { color: 'green', text: '订单收益' },
            gift: { color: 'pink', text: '礼物收益' },
            bonus: { color: 'gold', text: '奖励' },
            commission: { color: 'purple', text: '佣金' },
        };
        return <Tag color={config[type].color}>{config[type].text}</Tag>;
    };

    const getStatusTag = (status: EarningsRecord['status']) => {
        const config = {
            settled: { color: 'green', text: '已结算' },
            pending: { color: 'orange', text: '待结算' },
            frozen: { color: 'blue', text: '冻结中' },
        };
        return <Tag color={config[status].color}>{config[status].text}</Tag>;
    };

    const getWithdrawStatusTag = (status: WithdrawRecord['status']) => {
        const config = {
            pending: { color: 'orange', text: '待审核' },
            processing: { color: 'blue', text: '处理中' },
            success: { color: 'green', text: '已到账' },
            failed: { color: 'red', text: '失败' },
        };
        return <Tag color={config[status].color}>{config[status].text}</Tag>;
    };

    const earningsColumns: ColumnsType<EarningsRecord> = [
        { title: '类型', dataIndex: 'type', key: 'type', render: (type) => getTypeTag(type) },
        { title: '金额', dataIndex: 'amount', key: 'amount', render: (amount) => <Text style={{ color: '#52c41a', fontWeight: 500 }}>+¥{amount.toFixed(2)}</Text> },
        { title: '描述', dataIndex: 'description', key: 'description' },
        { title: '关联订单', dataIndex: 'orderNo', key: 'orderNo', render: (orderNo) => orderNo || '-' },
        { title: '时间', dataIndex: 'createdAt', key: 'createdAt' },
        { title: '状态', dataIndex: 'status', key: 'status', render: (status) => getStatusTag(status) },
    ];

    const withdrawColumns: ColumnsType<WithdrawRecord> = [
        { title: '提现金额', dataIndex: 'amount', key: 'amount', render: (amount) => `¥${amount.toFixed(2)}` },
        { title: '手续费', dataIndex: 'fee', key: 'fee', render: (fee) => `¥${fee.toFixed(2)}` },
        { title: '实际到账', dataIndex: 'actualAmount', key: 'actualAmount', render: (amount) => <Text strong>¥{amount.toFixed(2)}</Text> },
        { title: '收款账户', key: 'bank', render: (_, record) => <><Text>{record.bankName}</Text><br /><Text type="secondary">{record.bankAccount}</Text></> },
        { title: '申请时间', dataIndex: 'createdAt', key: 'createdAt' },
        { title: '状态', dataIndex: 'status', key: 'status', render: (status) => getWithdrawStatusTag(status) },
    ];

    const monthProgress = earningsInfo.monthEarnings > 0 ? Math.min((earningsInfo.monthEarnings / 10000) * 100, 100) : 0;

    // 图表数据
    const dailyEarningsData = [
        { date: '12-01', earnings: 120, orders: 3 },
        { date: '12-02', earnings: 180, orders: 5 },
        { date: '12-03', earnings: 150, orders: 4 },
        { date: '12-04', earnings: 220, orders: 6 },
        { date: '12-05', earnings: 190, orders: 5 },
        { date: '12-06', earnings: 280, orders: 7 },
        { date: '12-07', earnings: 310, orders: 8 },
        { date: '12-08', earnings: 250, orders: 6 },
        { date: '12-09', earnings: 200, orders: 5 },
        { date: '12-10', earnings: 240, orders: 6 },
        { date: '12-11', earnings: 170, orders: 4 },
        { date: '12-12', earnings: 290, orders: 7 },
        { date: '12-13', earnings: 320, orders: 8 },
        { date: '12-14', earnings: 350, orders: 9 },
        { date: '12-15', earnings: 280, orders: 7 },
        { date: '12-16', earnings: 380, orders: 10 },
    ];

    const weeklyEarningsData = [
        { week: '第1周', earnings: 950, orders: 23 },
        { week: '第2周', earnings: 1180, orders: 29 },
        { week: '第3周', earnings: 1340, orders: 33 },
        { week: '第4周', earnings: 1210, orders: 30 },
    ];

    const earningsTypeData = [
        { name: '订单收益', value: 4580, color: '#52c41a' },
        { name: '礼物收益', value: 680, color: '#eb2f96' },
        { name: '奖励', value: 320, color: '#faad14' },
        { name: '其他', value: 100, color: '#1890ff' },
    ];

    return (
        <div style={{ padding: 24 }}>
            <Title level={4}><DollarOutlined /> 我的收益</Title>

            {/* 余额卡片 */}
            <Card style={{ marginBottom: 16, background: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)' }}>
                <Row gutter={24} align="middle">
                    <Col flex="auto">
                        <Text style={{ color: 'rgba(255,255,255,0.8)' }}>可提现余额</Text>
                        <Title level={2} style={{ color: '#fff', margin: '8px 0' }}>¥ {earningsInfo.availableBalance.toFixed(2)}</Title>
                        <Space>
                            <Text style={{ color: 'rgba(255,255,255,0.6)' }}>冻结中: ¥{earningsInfo.frozenBalance.toFixed(2)}</Text>
                            <Text style={{ color: 'rgba(255,255,255,0.6)' }}>|</Text>
                            <Text style={{ color: 'rgba(255,255,255,0.6)' }}>提现中: ¥{earningsInfo.pendingWithdraw.toFixed(2)}</Text>
                        </Space>
                    </Col>
                    <Col>
                        <Button type="primary" size="large" icon={<BankOutlined />} onClick={() => setWithdrawVisible(true)}
                            style={{ background: '#fff', color: '#11998e', border: 'none' }}>申请提现</Button>
                    </Col>
                </Row>
            </Card>

            {/* 收益统计 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic title="今日收益" value={earningsInfo.todayEarnings} prefix={<RiseOutlined />} suffix="元" valueStyle={{ color: '#3f8600' }} />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic title="本周收益" value={earningsInfo.weekEarnings} prefix={<CalendarOutlined />} suffix="元" valueStyle={{ color: '#1890ff' }} />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic title="本月收益" value={earningsInfo.monthEarnings} prefix={<WalletOutlined />} suffix="元" />
                    </Card>
                </Col>
                <Col xs={12} sm={6}>
                    <Card loading={loading}>
                        <Statistic title="累计收益" value={earningsInfo.totalEarnings} prefix={<DollarOutlined />} suffix="元" valueStyle={{ color: '#722ed1' }} />
                    </Card>
                </Col>
            </Row>

            {/* 月度目标 */}
            <Card style={{ marginBottom: 16 }} loading={loading}>
                <Row gutter={24} align="middle">
                    <Col span={16}>
                        <Title level={5}>本月收益目标</Title>
                        <Progress percent={monthProgress} status={monthProgress >= 100 ? 'success' : 'active'}
                            strokeColor={{ '0%': '#108ee9', '100%': '#87d068' }} format={() => `¥${earningsInfo.monthEarnings} / ¥10000`} />
                    </Col>
                    <Col span={8} style={{ textAlign: 'right' }}>
                        <Text type="secondary">距离目标还差</Text>
                        <Title level={4} style={{ margin: 0, color: monthProgress >= 100 ? '#52c41a' : '#1890ff' }}>
                            {monthProgress >= 100 ? '已达成！' : `¥${(10000 - earningsInfo.monthEarnings).toFixed(2)}`}
                        </Title>
                    </Col>
                </Row>
            </Card>

            {/* 收益图表 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col xs={24} lg={16}>
                    <Card
                        title={
                            <Space>
                                <LineChartOutlined />
                                <span>每日收益趋势</span>
                            </Space>
                        }
                        extra={
                            <RangePicker
                                value={dateRange}
                                onChange={(dates) => setDateRange(dates as [dayjs.Dayjs, dayjs.Dayjs])}
                                allowClear={false}
                            />
                        }
                    >
                        <ResponsiveContainer width="100%" height={300}>
                            <AreaChart data={dailyEarningsData}>
                                <defs>
                                    <linearGradient id="colorEarnings" x1="0" y1="0" x2="0" y2="1">
                                        <stop offset="5%" stopColor="#1890ff" stopOpacity={0.8}/>
                                        <stop offset="95%" stopColor="#1890ff" stopOpacity={0}/>
                                    </linearGradient>
                                </defs>
                                <CartesianGrid strokeDasharray="3 3" />
                                <XAxis dataKey="date" />
                                <YAxis />
                                <Tooltip
                                    formatter={(value: number) => [`¥${value}`, '收益']}
                                    labelFormatter={(label) => `日期: ${label}`}
                                />
                                <Legend />
                                <Area
                                    type="monotone"
                                    dataKey="earnings"
                                    stroke="#1890ff"
                                    fillOpacity={1}
                                    fill="url(#colorEarnings)"
                                    name="收益(元)"
                                />
                            </AreaChart>
                        </ResponsiveContainer>
                    </Card>
                </Col>
                <Col xs={24} lg={8}>
                    <Card
                        title={
                            <Space>
                                <TrophyOutlined />
                                <span>收益构成</span>
                            </Space>
                        }
                    >
                        <ResponsiveContainer width="100%" height={300}>
                            <PieChart>
                                <Pie
                                    data={earningsTypeData}
                                    cx="50%"
                                    cy="50%"
                                    labelLine={false}
                                    label={({ name, percent }) => `${name} ${((percent ?? 0) * 100).toFixed(0)}%`}
                                    outerRadius={80}
                                    fill="#8884d8"
                                    dataKey="value"
                                >
                                    {earningsTypeData.map((entry, index) => (
                                        <Cell key={`cell-${index}`} fill={entry.color} />
                                    ))}
                                </Pie>
                                <Tooltip formatter={(value: number) => `¥${value}`} />
                            </PieChart>
                        </ResponsiveContainer>
                    </Card>
                </Col>
            </Row>

            {/* 周收益统计 */}
            <Row gutter={16} style={{ marginBottom: 16 }}>
                <Col span={24}>
                    <Card
                        title={
                            <Space>
                                <BarChartOutlined />
                                <span>每周收益统计</span>
                            </Space>
                        }
                    >
                        <ResponsiveContainer width="100%" height={250}>
                            <BarChart data={weeklyEarningsData}>
                                <CartesianGrid strokeDasharray="3 3" />
                                <XAxis dataKey="week" />
                                <YAxis />
                                <Tooltip
                                    formatter={(value: number, name: string) => [
                                        name === 'earnings' ? `¥${value}` : value,
                                        name === 'earnings' ? '收益' : '订单数'
                                    ]}
                                />
                                <Legend />
                                <Bar dataKey="earnings" fill="#52c41a" name="收益(元)" radius={[8, 8, 0, 0]} />
                                <Bar dataKey="orders" fill="#1890ff" name="订单数" radius={[8, 8, 0, 0]} />
                            </BarChart>
                        </ResponsiveContainer>
                    </Card>
                </Col>
            </Row>

            {/* 收益明细 */}
            <Card title={<><HistoryOutlined /> 收益明细</>}>
                <Tabs defaultActiveKey="earnings" items={[
                    { key: 'earnings', label: '收益记录', children: <Table columns={earningsColumns} dataSource={earningsRecords} rowKey="id" pagination={{ pageSize: 10 }} /> },
                    { key: 'withdraw', label: '提现记录', children: <Table columns={withdrawColumns} dataSource={withdrawRecords} rowKey="id" pagination={{ pageSize: 10 }} /> },
                ]} />
            </Card>

            {/* 提现弹窗 */}
            <Modal title="申请提现" open={withdrawVisible} onCancel={() => setWithdrawVisible(false)} footer={null} width={480}>
                <Form form={form} layout="vertical" onFinish={handleWithdraw}>
                    <Form.Item label="可提现余额">
                        <Text strong style={{ fontSize: 24, color: '#52c41a' }}>¥ {earningsInfo.availableBalance.toFixed(2)}</Text>
                    </Form.Item>
                    <Form.Item name="amount" label="提现金额" rules={[{ required: true, message: '请输入提现金额' }]}>
                        <InputNumber min={100} max={earningsInfo.availableBalance} prefix="¥" style={{ width: '100%' }} placeholder="最低提现100元" />
                    </Form.Item>
                    <Form.Item name="bankId" label="收款账户" rules={[{ required: true, message: '请选择收款账户' }]}>
                        <Select placeholder="选择收款银行卡" options={[
                            { value: 1, label: '招商银行 **** **** **** 6789' },
                            { value: 2, label: '工商银行 **** **** **** 1234' },
                        ]} />
                    </Form.Item>
                    <Form.Item>
                        <Text type="secondary">提现手续费：0.5%（最低2元），预计1-3个工作日到账</Text>
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" loading={withdrawLoading} block>确认提现</Button>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default PlayerEarnings;
