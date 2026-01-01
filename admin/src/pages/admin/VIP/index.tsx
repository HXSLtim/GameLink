/**
 * VIP等级管理页面
 * 卡片式展示VIP等级，支持增删改查操作
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Row,
    Col,
    Button,
    Space,
    Tag,
    Typography,
    Image,
    App,
    Popconfirm,
    Switch,
    Statistic,
    Segmented,
    Input,
    theme,
} from 'antd';
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    CrownOutlined,
    StarOutlined,
    GiftOutlined,
    ReloadOutlined,
    CheckCircleOutlined,
} from '@ant-design/icons';
import { vipApi } from '@/api/vip';
import type { VIPLevel } from '@/api/vip';
import { VIP_PERMISSIONS } from '@/constants/permissions';
import { PermissionGuard } from '@/components/PermissionGuard';
import { StateContainer } from '@/components/common/StateContainer';
import LevelForm from './components/LevelForm';
import BenefitsEditor from './components/BenefitsEditor';
import dayjs from 'dayjs';

const { Title, Text, Paragraph } = Typography;
const { Search } = Input;

const VIPPage: React.FC = () => {
    const { message } = App.useApp();
    const { token } = theme.useToken();
    const [loading, setLoading] = useState(false);
    const [levels, setLevels] = useState<VIPLevel[]>([]);
    const [filteredLevels, setFilteredLevels] = useState<VIPLevel[]>([]);
    const [loadError, setLoadError] = useState<string | null>(null);
    const [formVisible, setFormVisible] = useState(false);
    const [currentLevel, setCurrentLevel] = useState<VIPLevel | null>(null);
    const [benefitsEditorVisible, setBenefitsEditorVisible] = useState(false);
    const [benefitsEditingLevel, setBenefitsEditingLevel] = useState<VIPLevel | null>(null);
    const [viewMode, setViewMode] = useState<'all' | 'active' | 'inactive'>('all');
    const [keyword, setKeyword] = useState('');

    const loadData = useCallback(async () => {
        setLoading(true);
        setLoadError(null);
        try {
            const response = await vipApi.getVIPLevels({ page_size: 100 });
            if (response.data.success) {
                const data = response.data.data || [];
                setLevels(data);
                setFilteredLevels(data);
            } else {
                setLoadError(response.data.message || '加载失败');
                message.error(response.data.message || '加载失败');
            }
        } catch (error) {
            console.error('Load VIP levels error:', error);
            setLoadError('加载VIP等级失败');
            message.error('加载VIP等级失败');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadData();
    }, [loadData]);

    useEffect(() => {
        let filtered = [...levels];

        // 状态筛选
        if (viewMode === 'active') {
            filtered = filtered.filter(l => l.isActive);
        } else if (viewMode === 'inactive') {
            filtered = filtered.filter(l => !l.isActive);
        }

        // 关键词搜索
        if (keyword) {
            const kw = keyword.toLowerCase();
            filtered = filtered.filter(l =>
                l.title.toLowerCase().includes(kw) ||
                l.slug.toLowerCase().includes(kw)
            );
        }

        // 按排序顺序排序
        filtered.sort((a, b) => a.sortOrder - b.sortOrder);

        setFilteredLevels(filtered);
    }, [levels, viewMode, keyword]);

    const handleAdd = () => {
        setCurrentLevel(null);
        setFormVisible(true);
    };

    const handleEdit = (level: VIPLevel) => {
        setCurrentLevel(level);
        setFormVisible(true);
    };

    const handleDelete = async (level: VIPLevel) => {
        try {
            await vipApi.deleteVIPLevel(level.id);
            message.success('删除成功');
            loadData();
        } catch (error) {
            console.error('Delete VIP level error:', error);
            message.error('删除失败');
        }
    };

    const handleSetDefault = async (level: VIPLevel) => {
        try {
            await vipApi.setDefaultVIPLevel(level.id);
            message.success('设置成功');
            loadData();
        } catch (error) {
            console.error('Set default VIP level error:', error);
            message.error('设置失败');
        }
    };

    const handleToggleActive = async (level: VIPLevel, isActive: boolean) => {
        try {
            await vipApi.updateVIPLevel(level.id, { ...level, isActive });
            message.success(isActive ? '已启用' : '已禁用');
            loadData();
        } catch (error) {
            console.error('Toggle VIP level status error:', error);
            message.error('操作失败');
        }
    };

    const handleEditBenefits = (level: VIPLevel) => {
        setBenefitsEditingLevel(level);
        setBenefitsEditorVisible(true);
    };

    const handleSaveBenefits = async (benefits: string) => {
        if (benefitsEditingLevel) {
            try {
                await vipApi.updateVIPLevel(benefitsEditingLevel.id, {
                    ...benefitsEditingLevel,
                    benefits,
                });
                message.success('更新权益成功');
                loadData();
            } catch (error) {
                console.error('Update benefits error:', error);
                message.error('更新权益失败');
            }
        }
    };

    const parseBenefits = (benefitsStr: string): string[] => {
        if (!benefitsStr) return [];
        try {
            const parsed = JSON.parse(benefitsStr);
            return Array.isArray(parsed) ? parsed : [];
        } catch {
            return benefitsStr.split('\n').filter(s => s.trim());
        }
    };

    const stats = {
        total: levels.length,
        active: levels.filter(l => l.isActive).length,
        defaultCount: levels.filter(l => l.isDefault).length,
    };

    const renderVIPCard = (level: VIPLevel) => {
        const benefits = parseBenefits(level.benefits);
        const levelColor = level.color || '#1890ff';

        return (
            <Col key={level.id} xs={24} sm={12} lg={8} xl={6}>
                <Card
                    hoverable
                    style={{
                        height: '100%',
                        borderTop: `4px solid ${levelColor}`,
                        position: 'relative',
                    }}
                    bodyStyle={{ padding: 16 }}
                >
                    {/* 状态标签 */}
                    <div style={{ position: 'absolute', top: 12, right: 12 }}>
                        <Space size={4}>
                            {level.isDefault && (
                                <Tag color="gold" icon={<StarOutlined />}>
                                    默认
                                </Tag>
                            )}
                            {!level.isActive && (
                                <Tag color="default">已禁用</Tag>
                            )}
                        </Space>
                    </div>

                    {/* 等级头部 */}
                    <div style={{ textAlign: 'center', marginBottom: 16 }}>
                        {level.iconUrl && (
                            <div style={{ marginBottom: 8 }}>
                                <Image
                                    src={level.iconUrl}
                                    alt={level.title}
                                    width={48}
                                    height={48}
                                    preview={false}
                                    style={{ borderRadius: '50%' }}
                                    fallback="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mN8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="
                                />
                            </div>
                        )}
                        <Title level={5} style={{ margin: 0, color: levelColor }}>
                            <CrownOutlined style={{ marginRight: 4 }} />
                            {level.title}
                        </Title>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            {level.slug}
                        </Text>
                    </div>

                    {/* 等级属性 */}
                    <div style={{ marginBottom: 16 }}>
                        <Row gutter={[8, 8]}>
                            <Col span={12}>
                                <Statistic
                                    title="升级经验"
                                    value={level.expRequired}
                                    valueStyle={{ fontSize: 16, color: token.colorPrimary }}
                                />
                            </Col>
                            <Col span={12}>
                                <Statistic
                                    title="订单折扣"
                                    value={level.orderDiscount * 100}
                                    suffix="%"
                                    valueStyle={{ fontSize: 16, color: token.colorSuccess }}
                                />
                            </Col>
                        </Row>
                        {level.monthlyCouponCount > 0 && (
                            <div style={{ marginTop: 8 }}>
                                <Tag color="purple" icon={<GiftOutlined />}>
                                    每月{level.monthlyCouponCount}张优惠券
                                </Tag>
                            </div>
                        )}
                    </div>

                    {/* 会员权益 */}
                    {benefits.length > 0 && (
                        <div style={{ marginBottom: 16 }}>
                            <Text strong style={{ fontSize: 12 }}>会员权益：</Text>
                            <div style={{ marginTop: 4, maxHeight: 80, overflow: 'auto' }}>
                                {benefits.slice(0, 3).map((benefit, index) => (
                                    <div key={index} style={{ fontSize: 12, color: '#666', marginBottom: 2 }}>
                                        <CheckCircleOutlined style={{ color: levelColor, marginRight: 4 }} />
                                        {benefit}
                                    </div>
                                ))}
                                {benefits.length > 3 && (
                                    <Text type="secondary" style={{ fontSize: 12 }}>
                                        +{benefits.length - 3} 更多权益
                                    </Text>
                                )}
                            </div>
                        </div>
                    )}

                    {/* 操作按钮 */}
                    <div style={{ borderTop: '1px solid #f0f0f0', paddingTop: 12 }}>
                        <Space size={4} style={{ width: '100%' }} wrap>
                            <PermissionGuard permission={VIP_PERMISSIONS.UPDATE}>
                                <Button
                                    type="link"
                                    size="small"
                                    icon={<EditOutlined />}
                                    onClick={() => handleEdit(level)}
                                >
                                    编辑
                                </Button>
                            </PermissionGuard>
                            <Button
                                type="link"
                                size="small"
                                onClick={() => handleEditBenefits(level)}
                            >
                                权益
                            </Button>
                            {!level.isDefault && (
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={() => handleSetDefault(level)}
                                >
                                    设为默认
                                </Button>
                            )}
                            <PermissionGuard permission={VIP_PERMISSIONS.DELETE}>
                                <Popconfirm
                                    title="确定要删除此等级吗？"
                                    onConfirm={() => handleDelete(level)}
                                    okText="确定"
                                    cancelText="取消"
                                >
                                    <Button
                                        type="link"
                                        size="small"
                                        danger
                                        icon={<DeleteOutlined />}
                                    >
                                        删除
                                    </Button>
                                </Popconfirm>
                            </PermissionGuard>
                        </Space>
                        <div style={{ marginTop: 8 }}>
                            <Switch
                                size="small"
                                checked={level.isActive}
                                onChange={(checked) => handleToggleActive(level, checked)}
                                checkedChildren="启用"
                                unCheckedChildren="禁用"
                            />
                        </div>
                    </div>

                    {/* 时间信息 */}
                    <div style={{ marginTop: 8, borderTop: '1px solid #f0f0f0', paddingTop: 8 }}>
                        <Text type="secondary" style={{ fontSize: 11 }}>
                            创建于 {dayjs(level.createdAt).format('YYYY-MM-DD')}
                        </Text>
                    </div>
                </Card>
            </Col>
        );
    };

    return (
        <div style={{ padding: 24 }}>
            {/* 页面头部 */}
            <div style={{ marginBottom: 24 }}>
                <Row justify="space-between" align="middle">
                    <Col>
                        <Title level={4} style={{ margin: 0 }}>
                            <CrownOutlined style={{ marginRight: 8 }} />
                            VIP等级管理
                        </Title>
                        <Text type="secondary">管理平台VIP会员等级和权益配置</Text>
                    </Col>
                    <Col>
                        <Space>
                            <Button icon={<ReloadOutlined />} onClick={loadData} loading={loading}>
                                刷新
                            </Button>
                            <PermissionGuard permission={VIP_PERMISSIONS.CREATE}>
                                <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
                                    新增等级
                                </Button>
                            </PermissionGuard>
                        </Space>
                    </Col>
                </Row>
            </div>

            {/* 统计卡片 */}
            <Row gutter={16} style={{ marginBottom: 24 }}>
                <Col xs={24} sm={8}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic title="等级总数" value={stats.total} prefix={<CrownOutlined />} />
                    </Card>
                </Col>
                <Col xs={24} sm={8}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="已启用"
                            value={stats.active}
                            valueStyle={{ color: token.colorSuccess }}
                            prefix={<CheckCircleOutlined />}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={8}>
                    <Card style={{ minHeight: 120 }}>
                        <Statistic
                            title="默认等级"
                            value={stats.defaultCount}
                            valueStyle={{ color: token.colorWarning }}
                            prefix={<StarOutlined />}
                        />
                    </Card>
                </Col>
            </Row>

            {/* 筛选栏 */}
            <Card style={{ marginBottom: 16 }}>
                <Row justify="space-between" align="middle" gutter={16}>
                    <Col xs={24} sm={12}>
                        <Segmented
                            options={[
                                { label: '全部', value: 'all' },
                                { label: '已启用', value: 'active' },
                                { label: '已禁用', value: 'inactive' },
                            ]}
                            value={viewMode}
                            onChange={(v) => setViewMode(v as typeof viewMode)}
                        />
                    </Col>
                    <Col xs={24} sm={12} style={{ textAlign: 'right' }}>
                        <Search
                            placeholder="搜索等级名称或标识"
                            allowClear
                            style={{ maxWidth: 300 }}
                            onChange={(e) => setKeyword(e.target.value)}
                        />
                    </Col>
                </Row>
            </Card>

            {/* VIP等级卡片列表 */}
            <StateContainer
                loading={loading}
                data={filteredLevels}
                error={loadError}
                emptyType={keyword ? 'no-search' : 'no-data'}
                emptyTitle={keyword ? '未找到匹配的VIP等级' : '暂无VIP等级'}
                emptyDescription={keyword ? '请尝试调整搜索条件' : '创建第一个VIP等级开始使用'}
                emptyActionText={!keyword ? '创建第一个等级' : undefined}
                onEmptyAction={!keyword ? handleAdd : undefined}
                loadingConfig={{ card: false, rows: 4 }}
            >
                <Row gutter={[16, 16]}>
                    {filteredLevels.map(renderVIPCard)}
                </Row>
            </StateContainer>

            {/* 编辑表单弹窗 */}
            <LevelForm
                visible={formVisible}
                level={currentLevel}
                onCancel={() => setFormVisible(false)}
                onSuccess={loadData}
            />

            {/* 权益编辑器 */}
            <BenefitsEditor
                visible={benefitsEditorVisible}
                benefits={benefitsEditingLevel?.benefits || '[]'}
                onSave={handleSaveBenefits}
                onCancel={() => setBenefitsEditorVisible(false)}
            />
        </div>
    );
};

export default VIPPage;
