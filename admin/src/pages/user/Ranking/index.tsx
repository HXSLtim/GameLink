/**
 * 排行榜页面
 * 陪玩师排行、收益排行、好评排行
 */
import React, { useState, useEffect, useCallback } from 'react';
import { logger } from '@/utils/logger';
import {
    Card,
    Tabs,
    List,
    Avatar,
    Tag,
    Space,
    Typography,
    Spin,
    Row,
    Col,
    Select,
    Badge,
} from 'antd';
import {
    TrophyOutlined,
    CrownOutlined,
    FireOutlined,
    StarOutlined,
    RiseOutlined,
    UserOutlined,
} from '@ant-design/icons';

const { Title, Text } = Typography;

interface RankingPlayer {
    rank: number;
    id: number;
    nickname: string;
    avatar: string;
    level: number;
    value: number;
    change: number;
    games: string[];
    tags: string[];
}

const UserRanking: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [activeTab, setActiveTab] = useState('hot');
    const [timeRange, setTimeRange] = useState('week');
    const [rankings, setRankings] = useState<RankingPlayer[]>([]);

    const loadRankings = useCallback(async () => {
        setLoading(true);
        try {
            // 使用 activeTab 和 timeRange 构建请求参数
            logger.info('Loading rankings:', activeTab, timeRange);
            await new Promise(resolve => setTimeout(resolve, 500));
            const mockData: RankingPlayer[] = [
                { rank: 1, id: 1, nickname: '小甜甜', avatar: '', level: 6, value: 2580, change: 2, games: ['王者荣耀', '英雄联盟'], tags: ['声音甜美', '技术好'] },
                { rank: 2, id: 2, nickname: '大神带飞', avatar: '', level: 7, value: 2350, change: 0, games: ['英雄联盟', 'DOTA2'], tags: ['技术大神', '上分快'] },
                { rank: 3, id: 3, nickname: '温柔学姐', avatar: '', level: 5, value: 2180, change: 1, games: ['王者荣耀', '和平精英'], tags: ['温柔体贴', '新手友好'] },
                { rank: 4, id: 4, nickname: '电竞少女', avatar: '', level: 5, value: 1950, change: -1, games: ['英雄联盟'], tags: ['颜值高', '有趣'] },
                { rank: 5, id: 5, nickname: '游戏王子', avatar: '', level: 6, value: 1820, change: 3, games: ['王者荣耀', '原神'], tags: ['幽默风趣', '耐心'] },
                { rank: 6, id: 6, nickname: '萌萌哒', avatar: '', level: 4, value: 1680, change: 0, games: ['和平精英'], tags: ['可爱', '陪聊'] },
                { rank: 7, id: 7, nickname: '职业选手', avatar: '', level: 7, value: 1550, change: -2, games: ['英雄联盟', 'DOTA2'], tags: ['职业水平', '教学'] },
                { rank: 8, id: 8, nickname: '暖心姐姐', avatar: '', level: 5, value: 1420, change: 1, games: ['王者荣耀'], tags: ['暖心', '好评多'] },
                { rank: 9, id: 9, nickname: '技术流', avatar: '', level: 6, value: 1350, change: 0, games: ['DOTA2'], tags: ['技术流', '稳定'] },
                { rank: 10, id: 10, nickname: '开心果', avatar: '', level: 4, value: 1280, change: 2, games: ['和平精英', '原神'], tags: ['开朗', '有趣'] },
            ];
            setRankings(mockData);
        } catch {
            logger.error('加载排行榜失败');
        } finally {
            setLoading(false);
        }
    }, [activeTab, timeRange]);

    useEffect(() => { loadRankings(); }, [loadRankings]);

    const getRankIcon = (rank: number) => {
        if (rank === 1) return <CrownOutlined style={{ color: '#ffd700', fontSize: 24 }} />;
        if (rank === 2) return <CrownOutlined style={{ color: '#c0c0c0', fontSize: 22 }} />;
        if (rank === 3) return <CrownOutlined style={{ color: '#cd7f32', fontSize: 20 }} />;
        return <span style={{ fontSize: 18, fontWeight: 'bold', color: '#666' }}>{rank}</span>;
    };

    const getChangeTag = (change: number) => {
        if (change > 0) return <Tag color="green"><RiseOutlined /> +{change}</Tag>;
        if (change < 0) return <Tag color="red">↓ {Math.abs(change)}</Tag>;
        return <Tag>-</Tag>;
    };

    const getValueLabel = () => {
        switch (activeTab) {
            case 'hot': return '热度值';
            case 'order': return '订单数';
            case 'rating': return '好评率';
            case 'earnings': return '收益(元)';
            default: return '数值';
        }
    };

    const formatValue = (value: number) => {
        if (activeTab === 'rating') return `${value}%`;
        return value.toLocaleString();
    };

    const tabItems = [
        { key: 'hot', label: <><FireOutlined /> 人气榜</>, icon: <FireOutlined /> },
        { key: 'order', label: <><TrophyOutlined /> 订单榜</>, icon: <TrophyOutlined /> },
        { key: 'rating', label: <><StarOutlined /> 好评榜</>, icon: <StarOutlined /> },
        { key: 'earnings', label: <><RiseOutlined /> 收益榜</>, icon: <RiseOutlined /> },
    ];

    return (
        <div style={{ padding: 24 }}>
            <Card>
                <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
                    <Col>
                        <Title level={4} style={{ margin: 0 }}><TrophyOutlined /> 陪玩师排行榜</Title>
                    </Col>
                    <Col>
                        <Select value={timeRange} onChange={setTimeRange} style={{ width: 120 }}
                            options={[
                                { value: 'day', label: '今日' },
                                { value: 'week', label: '本周' },
                                { value: 'month', label: '本月' },
                                { value: 'all', label: '总榜' },
                            ]} />
                    </Col>
                </Row>

                <Tabs activeKey={activeTab} onChange={setActiveTab} items={tabItems.map(item => ({
                    key: item.key,
                    label: item.label,
                    children: (
                        <Spin spinning={loading}>
                            <List dataSource={rankings} renderItem={(player) => (
                                <List.Item style={{ padding: '16px 0', background: player.rank <= 3 ? 'linear-gradient(90deg, rgba(255,215,0,0.1) 0%, transparent 100%)' : 'transparent' }}>
                                    <List.Item.Meta
                                        avatar={
                                            <Space size="middle">
                                                <div style={{ width: 40, textAlign: 'center' }}>{getRankIcon(player.rank)}</div>
                                                <Badge count={`Lv.${player.level}`} offset={[-5, 45]} color="#1890ff">
                                                    <Avatar size={56} src={player.avatar} icon={<UserOutlined />} />
                                                </Badge>
                                            </Space>
                                        }
                                        title={
                                            <Space>
                                                <Text strong style={{ fontSize: 16 }}>{player.nickname}</Text>
                                                {getChangeTag(player.change)}
                                            </Space>
                                        }
                                        description={
                                            <Space direction="vertical" size={4}>
                                                <Space>{player.games.map(g => <Tag key={g} color="blue">{g}</Tag>)}</Space>
                                                <Space>{player.tags.slice(0, 2).map(t => <Tag key={t}>{t}</Tag>)}</Space>
                                            </Space>
                                        }
                                    />
                                    <div style={{ textAlign: 'right' }}>
                                        <Text type="secondary">{getValueLabel()}</Text>
                                        <Title level={4} style={{ margin: 0, color: player.rank <= 3 ? '#faad14' : '#1890ff' }}>
                                            {formatValue(player.value)}
                                        </Title>
                                    </div>
                                </List.Item>
                            )} />
                        </Spin>
                    ),
                }))} />
            </Card>
        </div>
    );
};

export default UserRanking;
