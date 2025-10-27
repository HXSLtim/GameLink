import { Card, Grid, Statistic, Space, Typography } from '@arco-design/web-react';
import { IconUser, IconThumbUp, IconClockCircle, IconApps } from '@arco-design/web-react/icon';
import styles from './Dashboard.module.less';

const { Row, Col } = Grid;

interface StatisticCardProps {
  title: string;
  value: number;
  prefix: React.ReactNode;
  suffix?: string;
  precision?: number;
}

/**
 * Statistic card component
 */
const StatisticCard: React.FC<StatisticCardProps> = ({
  title,
  value,
  prefix,
  suffix,
  precision,
}) => (
  <Card bordered>
    <Statistic title={title} value={value} prefix={prefix} suffix={suffix} precision={precision} />
  </Card>
);

/**
 * Dashboard page component
 *
 * @component
 * @description Displays an overview of key business metrics and statistics
 */
export const Dashboard = () => {
  return (
    <Space direction="vertical" size={16} className={styles.container}>
      <Typography.Title heading={4}>总览</Typography.Title>
      <Row gutter={16}>
        <Col span={6}>
          <StatisticCard title="用户数" value={1288} prefix={<IconUser />} />
        </Col>
        <Col span={6}>
          <StatisticCard title="订单数" value={342} prefix={<IconApps />} />
        </Col>
        <Col span={6}>
          <StatisticCard title="好评率" value={97.6} suffix="%" prefix={<IconThumbUp />} />
        </Col>
        <Col span={6}>
          <StatisticCard
            title="平均响应(s)"
            value={0.38}
            precision={2}
            prefix={<IconClockCircle />}
          />
        </Col>
      </Row>
      <Card title="近期动态" bordered>
        <Typography.Paragraph>这里展示最近 7 天的业务数据趋势与系统状态。</Typography.Paragraph>
      </Card>
    </Space>
  );
};
