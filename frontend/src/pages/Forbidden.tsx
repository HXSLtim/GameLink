/**
 * 403 禁止访问页面
 * 当用户没有权限访问某个页面时显示
 * Requirements: 8.3
 */
import React from 'react';
import { Result, Button, Space, Typography, Card, theme } from 'antd';
import { useNavigate, useLocation } from 'react-router-dom';
import { HomeOutlined, ArrowLeftOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';

const { Text, Paragraph } = Typography;

interface LocationState {
    from?: string;
    requiredPermission?: string;
}

const Forbidden: React.FC = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const { token } = theme.useToken();
    
    // 从路由状态获取来源页面和所需权限
    const state = location.state as LocationState | null;
    const fromPath = state?.from;
    const requiredPermission = state?.requiredPermission;

    const handleGoBack = () => {
        // 尝试返回上一页，如果没有历史记录则返回首页
        if (window.history.length > 2) {
            navigate(-1);
        } else {
            navigate('/admin');
        }
    };

    const handleGoHome = () => {
        navigate('/admin');
    };

    const handleRequestPermission = () => {
        // 可选功能：跳转到权限申请页面或发送邮件
        // 这里可以根据实际需求实现
        const subject = encodeURIComponent('权限申请');
        const body = encodeURIComponent(
            `您好，\n\n我需要申请以下权限：\n\n` +
            `页面路径: ${fromPath || '未知'}\n` +
            `所需权限: ${requiredPermission || '未知'}\n\n` +
            `请协助处理，谢谢！`
        );
        window.location.href = `mailto:admin@gamelink.com?subject=${subject}&body=${body}`;
    };

    return (
        <div
            style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                minHeight: '100vh',
                background: token.colorBgLayout,
                padding: 24,
            }}
        >
            <Card
                style={{
                    maxWidth: 600,
                    width: '100%',
                    boxShadow: token.boxShadowSecondary,
                }}
                variant="borderless"
            >
                <Result
                    status="403"
                    title="403"
                    icon={<LockOutlined style={{ color: token.colorError }} />}
                    subTitle={
                        <Space orientation="vertical" size={12} style={{ width: '100%' }}>
                            <Text strong style={{ fontSize: 16 }}>
                                抱歉，您没有权限访问此页面
                            </Text>
                            <Paragraph type="secondary" style={{ marginBottom: 0 }}>
                                您尝试访问的页面需要特定权限。如需访问，请联系管理员申请相应权限。
                            </Paragraph>
                            
                            {/* 显示详细信息 */}
                            {(fromPath || requiredPermission) && (
                                <div
                                    style={{
                                        background: token.colorBgContainerDisabled,
                                        padding: '12px 16px',
                                        borderRadius: token.borderRadius,
                                        marginTop: 8,
                                    }}
                                >
                                    {fromPath && (
                                        <div style={{ marginBottom: requiredPermission ? 8 : 0 }}>
                                            <Text type="secondary">访问路径：</Text>
                                            <Text code>{fromPath}</Text>
                                        </div>
                                    )}
                                    {requiredPermission && (
                                        <div>
                                            <Text type="secondary">所需权限：</Text>
                                            <Text code>{requiredPermission}</Text>
                                        </div>
                                    )}
                                </div>
                            )}
                        </Space>
                    }
                    extra={
                        <Space wrap>
                            <Button 
                                icon={<ArrowLeftOutlined />} 
                                onClick={handleGoBack}
                            >
                                返回上页
                            </Button>
                            <Button 
                                type="primary" 
                                icon={<HomeOutlined />} 
                                onClick={handleGoHome}
                            >
                                返回首页
                            </Button>
                            <Button
                                icon={<MailOutlined />}
                                onClick={handleRequestPermission}
                            >
                                申请权限
                            </Button>
                        </Space>
                    }
                />
            </Card>
        </div>
    );
};

export default Forbidden;
