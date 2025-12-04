import React from 'react';
import { Card, Empty } from 'antd';
import { PageContainer } from '@/components';

const UserBehavior: React.FC = () => {
    return (
        <PageContainer title="用户行为分析" subTitle="分析用户登录、使用习惯和消费行为">
            <Card>
                <Empty description="用户行为分析功能开发中..." />
            </Card>
        </PageContainer>
    );
};

export default UserBehavior;
