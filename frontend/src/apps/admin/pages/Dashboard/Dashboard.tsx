/**
 * 管理后台 - Dashboard页面
 */
import { Card, Grid, Statistic, Space, Typography } from '@arco-design/web-react';
import { IconUser, IconShoppingCart, IconDollarCircle, IconStar } from '@arco-design/web-react/icon';
import './Dashboard.less';

const { Row, Col } = Grid;
const { Title } = Typography;

/**
 * 统计卡片数据
 */
const statisticsData = [
  {
    title: '总用户数',
    value: 1234,
    prefix: <IconUser style={{ color: '#165DFF' }} />,
    suffix: '人',
  },
  {
    title: '今日订单',
    value: 89,
    prefix: <IconShoppingCart style={{ color: '#00B42A' }} />,
    suffix: '单',
  },
  {
    title: '总收入',
    value: 45678,
    prefix: <IconDollarCircle style={{ color: '#F7BA1E' }} />,
    suffix: '元',
    precision: 2,
  },
  {
    title: '平均评分',
    value: 4.8,
    prefix: <IconStar style={{ color: '#FF7D00' }} />,
    suffix: '分',
    precision: 1,
  },
];

/**
 * Admin Dashboard组件
 */
export const Dashboard = () => {
  return (
    <div className="admin-dashboard-page">
      <div className="page-header">
        <Title heading={4}>管理工作台</Title>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
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

      {/* 快捷操作 */}
      <Row gutter={[16, 16]}>
        <Col span={24}>
          <Card title="快速管理" bordered={false}>
            <Space size="large">
              <Card hoverable style={{ width: 200, textAlign: 'center', cursor: 'pointer' }}>
                <IconUser style={{ fontSize: 48, color: '#165DFF' }} />
                <div style={{ marginTop: 12 }}>用户管理</div>
              </Card>
              <Card hoverable style={{ width: 200, textAlign: 'center', cursor: 'pointer' }}>
                <IconShoppingCart style={{ fontSize: 48, color: '#00B42A' }} />
                <div style={{ marginTop: 12 }}>订单管理</div>
              </Card>
            </Space>
          </Card>
        </Col>
      </Row>
    </div>
  );
};
