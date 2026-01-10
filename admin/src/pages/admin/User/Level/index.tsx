import React from 'react';
import { Card, Table, Tag, Button, Space } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageContainer } from '@/components';
import { EditOutlined } from '@ant-design/icons';

interface LevelData {
    key: string;
    level: number;
    name: string;
    expRequired: number;
    benefits: string[];
}

const columns: ColumnsType<LevelData> = [
    {
        title: '等级',
        dataIndex: 'level',
        key: 'level',
        render: (text: number) => <Tag color="gold">Lv.{text}</Tag>,
    },
    {
        title: '名称',
        dataIndex: 'name',
        key: 'name',
    },
    {
        title: '所需经验值',
        dataIndex: 'expRequired',
        key: 'expRequired',
    },
    {
        title: '权益',
        dataIndex: 'benefits',
        key: 'benefits',
        render: (benefits: string[]) => (
            <Space wrap>
                {benefits.map((benefit, index) => (
                    <Tag key={index}>{benefit}</Tag>
                ))}
            </Space>
        ),
    },
    {
        title: '操作',
        key: 'action',
        fixed: 'right',
        width: 80,
        render: () => (
            <Space size={4}>
            <Button type="link" icon={<EditOutlined />}>
                编辑
            </Button>
        ),
    },
];

const data = [
    {
        key: '1',
        level: 1,
        name: '青铜会员',
        expRequired: 0,
        benefits: ['专属徽章'],
    },
    {
        key: '2',
        level: 2,
        name: '白银会员',
        expRequired: 1000,
        benefits: ['专属徽章', '生日礼包'],
    },
    {
        key: '3',
        level: 3,
        name: '黄金会员',
        expRequired: 5000,
        benefits: ['专属徽章', '生日礼包', '专属客服'],
    },
    {
        key: '4',
        level: 4,
        name: '铂金会员',
        expRequired: 20000,
        benefits: ['专属徽章', '生日礼包', '专属客服', '线下活动优先权'],
    },
    {
        key: '5',
        level: 5,
        name: '钻石会员',
        expRequired: 100000,
        benefits: ['专属徽章', '生日礼包', '专属客服', '线下活动优先权', '定制礼物'],
    },
];

const UserLevel: React.FC = () => {
    return (
        <PageContainer title="用户等级管理" subTitle="管理用户VIP等级和积分体系">
            <Card>
                <Table<LevelData>
                    columns={columns}
                    dataSource={data}
                    pagination={false}
                    scroll={{ x: 1000 }}
                />
            </Card>
        </PageContainer>
    );
};

export default UserLevel;
