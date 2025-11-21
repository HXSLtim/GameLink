/**
 * 陪玩师端 - Dashboard页面
 */
import { Card, Grid, Statistic, Typography } from '@arco-design/web-react';
import { IconShoppingCart, IconDollarCircle, IconStar, IconClockCircle } from '@arco-design/web-react/icon';
import './Dashboard.less';

const { Row, Col } = Grid;
const { Title } = Typography;

/**
 * 陪玩师统计数据
 */
const statisticsData = [
  {
    title: '今日订单',
    value: 5,
    prefix: <IconShoppingCart style={{ color: '#00B42A' }} />,
    suffix: '单',
  },
  {
    title: '本月收入',
    value: 3580,
    prefix: <IconDollarCircle style={{ color: '#F7BA1E' }} />,
    suffix: '元',
    precision: 2,
  },
  {
    title: '服务评分',
    value: 4.9,
    prefix: <IconStar style={{ color: '#FF7D00' }} />,
    suffix: '分',
    precision: 1,
  },
  {
    title: '在线时长',
    value: 156,
    prefix: <IconClockCircle style={{ color: '#165DFF' }} />,
    suffix: '小时',
  },
];

/**
 * Player Dashboard组件
 */
export const Dashboard = () => {
  return (
    <div className="player-dashboard-page">
      <div className="page-header">
        <Title heading={4}>陪玩工作台</Title>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]}>
        {statisticsData.map((stat, index) => (
          <Col key={index} xs={24} sm={12} md={12} lg={6} xl={6}>
            <Card hoverable>
              <Statistic
                title={stat.title}
                value={stat.value}
                prefix={stat.prefix}
                suffix={stat.suffix}
                precision={stat.precision}
                styleValue={{ fontSize: 28 }}
              />
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  );
};
