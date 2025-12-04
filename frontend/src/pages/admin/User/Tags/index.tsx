import React from 'react';
import { Card, Empty } from 'antd';
import { PageContainer } from '@/components';

const UserTags: React.FC = () => {
    return (
        <PageContainer title="用户标签管理" subTitle="管理用户标签体系">
            <Card>
                <Empty description="用户标签管理功能开发中..." />
            </Card>
        </PageContainer>
    );
};

export default UserTags;
