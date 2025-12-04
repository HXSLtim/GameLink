import React from 'react';
import { Card, Empty } from 'antd';
import { PageContainer } from '@/components';

const UserPortrait: React.FC = () => {
    return (
        <PageContainer title="用户画像分析" subTitle="分析用户群体特征、地域分布和偏好">
            <Card>
                <Empty description="用户画像分析功能开发中..." />
            </Card>
        </PageContainer>
    );
};

export default UserPortrait;
