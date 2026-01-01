/**
 * 实名认证申请页面
 * 陪玩师提交实名认证信息
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Form,
    Input,
    Button,
    Upload,
    message,
    Space,
    Typography,
    Alert,
    Image,
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
    SafetyOutlined,
} from '@ant-design/icons';
import type { UploadFile, UploadProps } from 'antd';
import { certificationApi, type IdentityCertification, type RankCertification, type CertificationStatus } from '@/api/certification';

const { Title, Text, Paragraph } = Typography;

interface CertificationStatusResponse {
    identityCertified: boolean;
    rankCertified: boolean;
    identityCertification?: IdentityCertification;
    rankCertification?: RankCertification;
}

/**
 * 实名认证申请页面
 */
const IdentityCertificationPage: React.FC = () => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [certification, setCertification] = useState<IdentityCertification | null>(null);
    const [idCardFrontList, setIdCardFrontList] = useState<UploadFile[]>([]);
    const [idCardBackList, setIdCardBackList] = useState<UploadFile[]>([]);

    /**
     * 加载认证状态
     */
    const loadCertificationStatus = useCallback(async () => {
        setLoading(true);
        try {
            const response = await certificationApi.getMyCertificationStatus();
            if (response.data.success) {
                const data = response.data.data as CertificationStatusResponse;
                if (data.identityCertification) {
                    setCertification(data.identityCertification);
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

            if (idCardFrontList.length === 0) {
                message.error('请上传身份证正面照片');
                return;
            }
            if (idCardBackList.length === 0) {
                message.error('请上传身份证反面照片');
                return;
            }

            setSubmitting(true);

            await certificationApi.createIdentityCertification({
                realName: values.realName,
                idCardNumber: values.idCardNumber,
                idCardFrontUrl: idCardFrontList[0].url || '',
                idCardBackUrl: idCardBackList[0].url || '',
            });

            message.success('认证申请已提交，请等待审核');
            form.resetFields();
            setIdCardFrontList([]);
            setIdCardBackList([]);
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

        return (
            <Card style={{ marginBottom: 24 }}>
                <Space direction="vertical" size="large" style={{ width: '100%' }}>
                    <div>
                        <Title level={4}>认证状态</Title>
                        {getStatusTag(certification.status)}
                    </div>

                    <Descriptions column={2} bordered size="small">
                        <Descriptions.Item label="真实姓名">{certification.realName}</Descriptions.Item>
                        <Descriptions.Item label="身份证号">
                            {certification.idCardNumber.replace(/^(.{6})(.*)(.{4})$/, '$1********$3')}
                        </Descriptions.Item>
                        <Descriptions.Item label="身份证正面" span={2}>
                            <Image
                                src={certification.idCardFrontUrl}
                                alt="身份证正面"
                                width={200}
                                style={{ borderRadius: 8 }}
                            />
                        </Descriptions.Item>
                        <Descriptions.Item label="身份证反面" span={2}>
                            <Image
                                src={certification.idCardBackUrl}
                                alt="身份证反面"
                                width={200}
                                style={{ borderRadius: 8 }}
                            />
                        </Descriptions.Item>
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
                            description="您的实名认证申请正在审核中，请耐心等待，通常会在 1-3 个工作日内完成审核。"
                            type="info"
                            showIcon
                        />
                    )}

                    {certification.status === 'approved' && (
                        <Alert
                            message="认证已通过"
                            description="恭喜！您的实名认证已通过审核，现在可以接单了。"
                            type="success"
                            showIcon
                        />
                    )}

                    {certification.status === 'rejected' && (
                        <Alert
                            message="认证未通过"
                            description={
                                <div>
                                    <p>很遗憾，您的实名认证未通过审核。</p>
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

        return (
            <>
                {certification?.status === 'rejected' && renderCertificationStatus()}

                <Card>
                    <Space direction="vertical" size="large" style={{ width: '100%' }}>
                        <div>
                            <Title level={4}>
                                <SafetyOutlined /> 实名认证
                            </Title>
                            <Paragraph type="secondary">
                                为了保障平台安全和用户权益，所有陪玩师都需要完成实名认证才能开始接单。
                                您的个人信息将被严格保密，仅用于身份验证。
                            </Paragraph>
                        </div>

                        <Alert
                            message="认证须知"
                            description={
                                <ul style={{ margin: 0, paddingLeft: 20 }}>
                                    <li>请确保上传的身份证照片清晰完整，信息可见</li>
                                    <li>姓名和身份证号必须与身份证照片一致</li>
                                    <li>审核时间为 1-3 个工作日</li>
                                    <li>认证通过后，姓名将在订单中展示给用户</li>
                                </ul>
                            }
                            type="info"
                            showIcon
                        />

                        <Form form={form} layout="vertical">
                            <Row gutter={16}>
                                <Col span={12}>
                                    <Form.Item
                                        name="realName"
                                        label="真实姓名"
                                        rules={[
                                            { required: true, message: '请输入真实姓名' },
                                            { min: 2, max: 20, message: '姓名长度为 2-20 个字符' },
                                        ]}
                                    >
                                        <Input placeholder="请输入您的真实姓名" />
                                    </Form.Item>
                                </Col>
                                <Col span={12}>
                                    <Form.Item
                                        name="idCardNumber"
                                        label="身份证号"
                                        rules={[
                                            { required: true, message: '请输入身份证号' },
                                            {
                                                pattern: /^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$/,
                                                message: '请输入正确的身份证号',
                                            },
                                        ]}
                                    >
                                        <Input placeholder="请输入18位身份证号" maxLength={18} />
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Row gutter={16}>
                                <Col span={12}>
                                    <Form.Item
                                        label="身份证正面照片"
                                        required
                                        extra="请上传身份证正面（人像面）照片，支持 JPG、PNG 格式"
                                    >
                                        <Upload
                                            listType="picture-card"
                                            fileList={idCardFrontList}
                                            onChange={({ fileList }) => setIdCardFrontList(fileList)}
                                            customRequest={handleUpload}
                                            maxCount={1}
                                            accept="image/*"
                                        >
                                            {idCardFrontList.length === 0 && (
                                                <div>
                                                    <UploadOutlined />
                                                    <div style={{ marginTop: 8 }}>上传照片</div>
                                                </div>
                                            )}
                                        </Upload>
                                    </Form.Item>
                                </Col>
                                <Col span={12}>
                                    <Form.Item
                                        label="身份证反面照片"
                                        required
                                        extra="请上传身份证反面（国徽面）照片，支持 JPG、PNG 格式"
                                    >
                                        <Upload
                                            listType="picture-card"
                                            fileList={idCardBackList}
                                            onChange={({ fileList }) => setIdCardBackList(fileList)}
                                            customRequest={handleUpload}
                                            maxCount={1}
                                            accept="image/*"
                                        >
                                            {idCardBackList.length === 0 && (
                                                <div>
                                                    <UploadOutlined />
                                                    <div style={{ marginTop: 8 }}>上传照片</div>
                                                </div>
                                            )}
                                        </Upload>
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Form.Item>
                                <Button
                                    type="primary"
                                    size="large"
                                    block
                                    loading={submitting}
                                    onClick={handleSubmit}
                                    icon={<SafetyOutlined />}
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

export default IdentityCertificationPage;
