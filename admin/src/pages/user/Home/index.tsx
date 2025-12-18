/**
 * 用户端首页
 * 展示陪玩师列表、推荐、搜索功能
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Row,
    Col,
    Input,
    Tag,
    Avatar,
    Rate,
    Button,
    Spin,
    Empty,
    Badge,
    Space,
    Typography,
    message,
} from 'antd';
import {
    SearchOutlined,
    FireOutlined,
    StarOutlined,
    UserOutlined,
    PlayCircleOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import styles from './index.module.css';

const { Title, Text, Paragraph } = Typography;
const { Search } = Input;

interface Player {
    id: number;
    nickname: string;
    avatar: string;
    level: number;
    rating: number;
    orderCount: number;
    price: number;
    games: string[];
    tags: string[];
    status: 'online' | 'busy' | 'offline';
    introduction: string;
}

interface Game {
    id: number;
    name: string;
    icon: string;
}

const UserHome: React.FC = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [players, setPlayers] = useState<Player[]>([]);
    const [games, setGames] = useState<Game[]>([]);
    const [selectedGame, setSelectedGame] = useState<number | undefined>();
    const [keyword, setKeyword] = useState('');


    // 模拟数据加载
    const loadPlayers = useCallback(async () => {
        setLoading(true);
        try {
            // TODO: 替换为真实 API 调用
            // const res = await userApi.getPlayers({ game_id: selectedGame, keyword });
            await new Promise(resolve => setTimeout(resolve, 500));
            
            // 模拟数据
            const mockPlayers: Player[] = [
                {
                    id: 1,
                    nickname: '小甜甜',
                    avatar: '',
                    level: 5,
                    rating: 4.9,
                    orderCount: 1280,
                    price: 50,
                    games: ['王者荣耀', '英雄联盟'],
                    tags: ['声音甜美', '技术好', '有耐心'],
                    status: 'online',
                    introduction: '王者荣耀国服最强辅助，带你上分不是梦~',
                },
                {
                    id: 2,
                    nickname: '大神带飞',
                    avatar: '',
                    level: 6,
                    rating: 4.8,
                    orderCount: 2150,
                    price: 80,
                    games: ['英雄联盟', 'DOTA2'],
                    tags: ['技术大神', '幽默风趣', '上分快'],
                    status: 'busy',
                    introduction: '职业选手退役，带你体验什么叫真正的Carry',
                },
                {
                    id: 3,
                    nickname: '温柔学姐',
                    avatar: '',
                    level: 4,
                    rating: 4.95,
                    orderCount: 890,
                    price: 45,
                    games: ['王者荣耀', '和平精英'],
                    tags: ['温柔体贴', '新手友好', '陪聊'],
                    status: 'online',
                    introduction: '不只是游戏，更是陪伴。新手玩家的最佳选择~',
                },
            ];
            setPlayers(mockPlayers);
        } catch (error) {
            console.error(error);
            message.error('加载失败');
        } finally {
            setLoading(false);
        }
    }, []);

    const loadGames = useCallback(async () => {
        // 模拟游戏列表
        setGames([
            { id: 1, name: '王者荣耀', icon: '👑' },
            { id: 2, name: '英雄联盟', icon: '⚔️' },
            { id: 3, name: '和平精英', icon: '🔫' },
            { id: 4, name: 'DOTA2', icon: '🎮' },
            { id: 5, name: '原神', icon: '🌟' },
        ]);
    }, []);

    useEffect(() => {
        loadGames();
        loadPlayers();
    }, [loadGames, loadPlayers]);

    const getStatusBadge = (status: Player['status']) => {
        const config = {
            online: { status: 'success' as const, text: '在线' },
            busy: { status: 'warning' as const, text: '忙碌' },
            offline: { status: 'default' as const, text: '离线' },
        };
        return config[status];
    };

    const handlePlayerClick = (player: Player) => {
        navigate(`/user/player/${player.id}`);
    };

    const handleOrder = (player: Player, e: React.MouseEvent) => {
        e.stopPropagation();
        navigate(`/user/order/create?playerId=${player.id}`);
    };

    return (
        <div className={styles.container}>
            {/* 搜索区域 */}
            <div className={styles.searchSection}>
                <Title level={2} style={{ color: '#fff', marginBottom: 24 }}>
                    找到你的专属陪玩
                </Title>
                <Row gutter={16} justify="center">
                    <Col xs={24} sm={16} md={12} lg={8}>
                        <Search
                            placeholder="搜索陪玩师昵称"
                            allowClear
                            enterButton={<><SearchOutlined /> 搜索</>}
                            size="large"
                            value={keyword}
                            onChange={e => setKeyword(e.target.value)}
                            onSearch={loadPlayers}
                        />
                    </Col>
                </Row>
                <div className={styles.gameFilter}>
                    <Space wrap>
                        <Tag
                            color={!selectedGame ? 'blue' : 'default'}
                            onClick={() => setSelectedGame(undefined)}
                            style={{ cursor: 'pointer', padding: '4px 12px' }}
                        >
                            全部游戏
                        </Tag>
                        {games.map(game => (
                            <Tag
                                key={game.id}
                                color={selectedGame === game.id ? 'blue' : 'default'}
                                onClick={() => setSelectedGame(game.id)}
                                style={{ cursor: 'pointer', padding: '4px 12px' }}
                            >
                                {game.icon} {game.name}
                            </Tag>
                        ))}
                    </Space>
                </div>
            </div>

            {/* 推荐陪玩师 */}
            <div className={styles.section}>
                <div className={styles.sectionHeader}>
                    <Title level={4}>
                        <FireOutlined style={{ color: '#ff4d4f' }} /> 热门推荐
                    </Title>
                </div>
                
                <Spin spinning={loading}>
                    {players.length === 0 ? (
                        <Empty description="暂无陪玩师" />
                    ) : (
                        <Row gutter={[16, 16]}>
                            {players.map((player, index) => (
                                <Col xs={24} sm={12} md={8} lg={6} key={player.id}>
                                    <motion.div
                                        initial={{ opacity: 0, y: 20 }}
                                        animate={{ opacity: 1, y: 0 }}
                                        transition={{ delay: index * 0.1 }}
                                    >
                                        <Card
                                            hoverable
                                            className={styles.playerCard}
                                            onClick={() => handlePlayerClick(player)}
                                        >
                                            <div className={styles.cardHeader}>
                                                <Badge {...getStatusBadge(player.status)} dot offset={[-5, 5]}>
                                                    <Avatar
                                                        size={64}
                                                        src={player.avatar}
                                                        icon={<UserOutlined />}
                                                    />
                                                </Badge>
                                                <div className={styles.levelBadge}>
                                                    Lv.{player.level}
                                                </div>
                                            </div>
                                            
                                            <div className={styles.cardBody}>
                                                <Title level={5} className={styles.nickname}>
                                                    {player.nickname}
                                                </Title>
                                                
                                                <div className={styles.rating}>
                                                    <Rate disabled defaultValue={player.rating} allowHalf />
                                                    <Text type="secondary">{player.rating}</Text>
                                                </div>
                                                
                                                <div className={styles.stats}>
                                                    <Space>
                                                        <StarOutlined />
                                                        <Text type="secondary">{player.orderCount} 单</Text>
                                                    </Space>
                                                </div>
                                                
                                                <Paragraph
                                                    ellipsis={{ rows: 2 }}
                                                    className={styles.intro}
                                                >
                                                    {player.introduction}
                                                </Paragraph>
                                                
                                                <div className={styles.tags}>
                                                    {player.tags.slice(0, 3).map(tag => (
                                                        <Tag key={tag} color="blue">{tag}</Tag>
                                                    ))}
                                                </div>
                                                
                                                <div className={styles.games}>
                                                    {player.games.map(game => (
                                                        <Tag key={game}>{game}</Tag>
                                                    ))}
                                                </div>
                                            </div>
                                            
                                            <div className={styles.cardFooter}>
                                                <div className={styles.price}>
                                                    <Text type="secondary">¥</Text>
                                                    <span className={styles.priceValue}>{player.price}</span>
                                                    <Text type="secondary">/小时</Text>
                                                </div>
                                                <Button
                                                    type="primary"
                                                    icon={<PlayCircleOutlined />}
                                                    onClick={(e) => handleOrder(player, e)}
                                                >
                                                    立即下单
                                                </Button>
                                            </div>
                                        </Card>
                                    </motion.div>
                                </Col>
                            ))}
                        </Row>
                    )}
                </Spin>
            </div>
        </div>
    );
};

export default UserHome;
