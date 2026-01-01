/**
 * 段位认证申请页面
 * 陪玩师提交游戏段位认证信息
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Form,
    Button,
    Upload,
    message,
    Space,
    Typography,
    Alert,
    Image,
    Select,
    Row,
    Col,
    Descriptions,
    Tag,
    Spin,
} from 'antd';
import {
    UploadOutlined,
    CheckCircleOutlined,
    ClockCircleOutlined,
    CloseCircleOutlined,
    TrophyOutlined,
} from '@ant-design/icons';
import type { UploadFile, UploadProps } from 'antd';
import { certificationApi, type RankCertification, type IdentityCertification, type CertificationStatus } from '@/api/certification';

const { Title, Text, Paragraph } = Typography;
const { Option } = Select;

interface CertificationStatusResponse {
    identityCertified: boolean;
    rankCertified: boolean;
    identityCertification?: IdentityCertification;
    rankCertification?: RankCertification;
}

/**
 * 游戏类型配置
 */
const GAME_TYPES = [
    { value: 'lol', label: '英雄联盟', ranks: ['黑铁', '青铜', '白银', '黄金', '铂金', '翡翠', '钻石', '大师', '宗师', '王者'] },
    { value: '王者荣耀', label: '王者荣耀', ranks: ['青铜', '白银', '黄金', '铂金', '钻石', '星耀', '王者', '荣耀王者'] },
    { value: 'dota2', label: 'DOTA 2', ranks: [' Herald', 'Crusader', 'Archon', 'Legend', 'Ancient', 'Divine', 'Immortal'] },
    { value: 'csgo', label: 'CS:GO', ranks: ['白银1', '白银2', '白银3', '白银4', '白银5', '白银6', '黄金1', '黄金2', '黄金3', '黄金4', '黄金5', '黄金6', '大师1', '大师2', '大师3', '大师4', '大师5', '大师6', '传奇', '全球精英'] },
    { value: 'valorant', label: 'Valorant', ranks: ['铁1', '铁2', '铁3', '铜1', '铜2', '铜3', '银1', '银2', '银3', '金1', '金2', '金3', '铂金1', '铂金2', '铂金3', '钻石1', '钻石2', '钻石3', '超凡1', '超凡2', '超凡3', '无敌'] },
];

/**
 * 段位认证申请页面
 */
const RankCertificationPage: React.FC = () => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [certification, setCertification] = useState<RankCertification | null>(null);
    const [screenshotList, setScreenshotList] = useState<UploadFile[]>([]);
    const [videoList, setVideoList] = useState<UploadFile[]>([]);
    const [selectedGame, setSelectedGame] = useState<string | null>(null);

    /**
     * 加载认证状态
     */
    const loadCertificationStatus = useCallback(async () => {
        setLoading(true);
        try {
            const response = await certificationApi.getMyCertificationStatus();
            if (response.data.success) {
                const data = response.data.data as CertificationStatusResponse;
                if (data.rankCertification) {
                    setCertification(data.rankCertification);
                }
            } else {
                message.error(response.data.message || '加载认证状态失败');
            }
        } catch (error) {
            console.error('Load certification status error:', error);
            message.error('加载认证状态失败');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadCertificationStatus();
    }, [loadCertificationStatus]);

    /**
     * 上传处理
     */
    const handleUpload: UploadProps['customRequest'] = async (options) => {
        const { file, onSuccess, onError } = options;
        try {
            // 模拟上传，实际应该调用上传 API
            await new Promise(resolve => setTimeout(resolve, 1000));
            const fileUrl = URL.createObjectURL(file as File);
            onSuccess?.({ url: fileUrl });
        } catch (error) {
            onError?.(error as Error);
        }
    };

    /**
     * 提交认证
     */
    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();

            if (screenshotList.length === 0) {
                message.error('请上传至少一张段位截图');
                return;
            }

            setSubmitting(true);

            const screenshotUrls = screenshotList.map(file => file.url || '');

            await certificationApi.createRankCertification({
                gameType: values.gameType,
                currentRank: values.currentRank,
                targetRank: values.targetRank,
                screenshotUrls,
                videoUrl: videoList.length > 0 ? videoList[0].url : undefined,
                additionalInfo: values.additionalInfo,
            });

            message.success('段位认证申请已提交，请等待审核');
            form.resetFields();
            setScreenshotList([]);
            setVideoList([]);
            loadCertificationStatus();
        } catch (error) {
            console.error('Submit certification error:', error);
            message.error('提交认证申请失败');
        } finally {
            setSubmitting(false);
        }
    };

    /**
     * 获取状态标签
     */
    const getStatusTag = (status: CertificationStatus) => {
        const config = {
            pending: { color: 'orange', icon: <ClockCircleOutlined />, text: '待审核' },
            approved: { color: 'success', icon: <CheckCircleOutlined />, text: '已通过' },
            rejected: { color: 'error', icon: <CloseCircleOutlined />, text: '已拒绝' },
        };
        const statusConfig = config[status];
        return (
            <Tag color={statusConfig.color} icon={statusConfig.icon}>
                {statusConfig.text}
            </Tag>
        );
    };

    /**
     * 认证状态显示
     */
    const renderCertificationStatus = () => {
        if (!certification) return null;

        const gameConfig = GAME_TYPES.find(g => g.value === certification.gameType);

        return (
            <Card style={{ marginBottom: 24 }}>
                <Space direction="vertical" size="large" style={{ width: '100%' }}>
                    <div>
                        <Title level={4}>认证状态</Title>
                        {getStatusTag(certification.status)}
                    </div>

                    <Descriptions column={2} bordered size="small">
                        <Descriptions.Item label="游戏类型">{gameConfig?.label || certification.gameType}</Descriptions.Item>
                        <Descriptions.Item label="当前段位">{certification.currentRank}</Descriptions.Item>
                        <Descriptions.Item label="目标段位">{certification.targetRank}</Descriptions.Item>
                        <Descriptions.Item label="段位截图" span={2}>
                            <Image.PreviewGroup>
                                <Space>
                                    {certification.screenshotUrls.map((url, index) => (
                                        <Image
                                            key={index}
                                            src={url}
                                            alt={`段位截图${index + 1}`}
                                            width={100}
                                            style={{ borderRadius: 8 }}
                                        />
                                    ))}
                                </Space>
                            </Image.PreviewGroup>
                        </Descriptions.Item>
                        {certification.videoUrl && (
                            <Descriptions.Item label="视频证明" span={2}>
                                <video
                                    src={certification.videoUrl}
                                    controls
                                    style={{ width: '100%', maxWidth: 400 }}
                                />
                            </Descriptions.Item>
                        )}
                        {certification.additionalInfo && (
                            <Descriptions.Item label="补充说明" span={2}>
                                {certification.additionalInfo}
                            </Descriptions.Item>
                        )}
                        <Descriptions.Item label="申请时间" span={2}>
                            {new Date(certification.createdAt).toLocaleString('zh-CN')}
                        </Descriptions.Item>
                        {certification.status === 'rejected' && certification.rejectReason && (
                            <Descriptions.Item label="拒绝原因" span={2}>
                                <Text type="danger">{certification.rejectReason}</Text>
                            </Descriptions.Item>
                        )}
                        {certification.reviewedAt && (
                            <Descriptions.Item label="审核时间" span={2}>
                                {new Date(certification.reviewedAt).toLocaleString('zh-CN')}
                            </Descriptions.Item>
                        )}
                    </Descriptions>

                    {certification.status === 'pending' && (
                        <Alert
                            message="认证审核中"
                            description="您的段位认证申请正在审核中，请耐心等待，通常会在 1-3 个工作日内完成审核。"
                            type="info"
                            showIcon
                        />
                    )}

                    {certification.status === 'approved' && (
                        <Alert
                            message="认证已通过"
                            description="恭喜！您的段位认证已通过审核，您的段位信息已更新，可以接对应段位的订单了。"
                            type="success"
                            showIcon
                        />
                    )}

                    {certification.status === 'rejected' && (
                        <Alert
                            message="认证未通过"
                            description={
                                <div>
                                    <p>很遗憾，您的段位认证未通过审核。</p>
                                    <p>拒绝原因：{certification.rejectReason}</p>
                                    <p>您可以修改后重新提交认证申请。</p>
                                </div>
                            }
                            type="error"
                            showIcon
                            action={
                                <Button
                                    type="primary"
                                    danger
                                    onClick={() => {
                                        setCertification(null);
                                    }}
                                >
                                    重新认证
                                </Button>
                            }
                        />
                    )}
                </Space>
            </Card>
        );
    };

    /**
     * 认证申请表单
     */
    const renderCertificationForm = () => {
        if (certification && certification.status !== 'rejected') {
            return renderCertificationStatus();
        }

        const currentRanks = selectedGame
            ? GAME_TYPES.find(g => g.value === selectedGame)?.ranks || []
            : [];

        return (
            <>
                {certification?.status === 'rejected' && renderCertificationStatus()}

                <Card>
                    <Space direction="vertical" size="large" style={{ width: '100%' }}>
                        <div>
                            <Title level={4}>
                                <TrophyOutlined /> 段位认证
                            </Title>
                            <Paragraph type="secondary">
                                请提交您的游戏段位证明，审核通过后将更新您的段位信息。
                                段位信息将影响您能够接的订单类型和价格。
                            </Paragraph>
                        </div>

                        <Alert
                            message="认证须知"
                            description={
                                <ul style={{ margin: 0, paddingLeft: 20 }}>
                                    <li>请上传清晰的段位截图，必须包含当前段位和召唤师名称</li>
                                    <li>视频证明可以帮助加速审核，建议上传对局记录视频</li>
                                    <li>审核时间为 1-3 个工作日</li>
                                    <li>虚假认证将被永久封禁</li>
                                </ul>
                            }
                            type="warning"
                            showIcon
                        />

                        <Form form={form} layout="vertical">
                            <Row gutter={16}>
                                <Col span={12}>
                                    <Form.Item
                                        name="gameType"
                                        label="游戏类型"
                                        rules={[{ required: true, message: '请选择游戏类型' }]}
                                    >
                                        <Select
                                            placeholder="请选择游戏"
                                            onChange={(value) => setSelectedGame(value)}
                                        >
                                            {GAME_TYPES.map(game => (
                                                <Option key={game.value} value={game.value}>
                                                    {game.label}
                                                </Option>
                                            ))}
                                        </Select>
                                    </Form.Item>
                                </Col>
                                <Col span={12}>
                                    <Form.Item
                                        name="currentRank"
                                        label="当前段位"
                                        rules={[{ required: true, message: '请选择当前段位' }]}
                                    >
                                        <Select placeholder="请选择当前段位">
                                            {currentRanks.map(rank => (
                                                <Option key={rank} value={rank}>
                                                    {rank}
                                                </Option>
                                            ))}
                                        </Select>
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Row gutter={16}>
                                <Col span={12}>
                                    <Form.Item
                                        name="targetRank"
                                        label="目标段位"
                                        rules={[{ required: true, message: '请选择目标段位' }]}
                                    >
                                        <Select placeholder="请选择目标段位">
                                            {currentRanks.map(rank => (
                                                <Option key={rank} value={rank}>
                                                    {rank}
                                                </Option>
                                            ))}
                                        </Select>
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Form.Item
                                label="段位截图"
                                required
                                extra="请上传至少一张清晰的段位截图，支持 JPG、PNG 格式"
                            >
                                <Upload
                                    listType="picture-card"
                                    fileList={screenshotList}
                                    onChange={({ fileList }) => setScreenshotList(fileList)}
                                    customRequest={handleUpload}
                                    multiple
                                    accept="image/*"
                                >
                                    <div>
                                        <UploadOutlined />
                                        <div style={{ marginTop: 8 }}>上传截图</div>
                                    </div>
                                </Upload>
                            </Form.Item>

                            <Form.Item
                                label="视频证明（可选）"
                                extra="可以上传游戏视频作为辅助证明，支持 MP4 格式，最大 50MB"
                            >
                                <Upload
                                    fileList={videoList}
                                    onChange={({ fileList }) => setVideoList(fileList)}
                                    customRequest={handleUpload}
                                    maxCount={1}
                                    accept="video/*"
                                >
                                    <Button icon={<UploadOutlined />}>选择视频</Button>
                                </Upload>
                            </Form.Item>

                            <Form.Item
                                name="additionalInfo"
                                label="补充说明（可选）"
                                extra="如有其他需要说明的信息，请在此填写"
                            >
                                <TextArea
                                    rows={3}
                                    placeholder="例如：账号ID、特殊说明等"
                                    maxLength={500}
                                    showCount
                                />
                            </Form.Item>

                            <Form.Item>
                                <Button
                                    type="primary"
                                    size="large"
                                    block
                                    loading={submitting}
                                    onClick={handleSubmit}
                                    icon={<TrophyOutlined />}
                                >
                                    提交认证申请
                                </Button>
                            </Form.Item>
                        </Form>
                    </Space>
                </Card>
            </>
        );
    };

    return (
        <div style={{ padding: 24 }}>
            <Spin spinning={loading}>
                {renderCertificationForm()}
            </Spin>
        </div>
    );
};

export default RankCertificationPage;
