import React from 'react';
import { Card, Empty } from 'antd';
import { PageContainer } from '@/components';

const UserLevel: React.FC = () => {
    return (
        <PageContainer title="用户等级管理" subTitle="管理用户VIP等级和积分体系">
            <Card>
                <Empty description="用户等级管理功能开发中..." />
            </Card>
        </PageContainer>
    );
};

export default UserLevel;
