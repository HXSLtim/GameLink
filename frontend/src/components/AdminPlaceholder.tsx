import React from 'react';
import { Card } from 'antd';

const PageComponent: React.FC<{ title: string }> = ({ title }) => {
    return (
        <Card title={title} bordered={false}>
            <p>Module: {title}</p>
            <p>Status: Under Development</p>
        </Card>
    );
};

export default PageComponent;
