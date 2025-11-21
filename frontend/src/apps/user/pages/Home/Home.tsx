/**
 * 用户端 - 首页
 */
import { Card, Grid, Typography, Button, Space } from '@arco-design/web-react';
import { IconSearch, IconShoppingCart } from '@arco-design/web-react/icon';
import './Home.less';

const { Row, Col } = Grid;
const { Title, Text } = Typography;

/**
 * User Home组件
 */
export const Home = () => {
  return (
    <div className="user-home-page">
      {/* Banner */}
      <Card className="banner-card" bordered={false}>
        <div className="banner-content">
          <Title heading={2}>欢迎来到 GameLink</Title>
          <Text style={{ fontSize: 16, marginTop: 12, display: 'block' }}>
            找到最适合你的游戏陪玩师
          </Text>
          <Space size="large" style={{ marginTop: 24 }}>
            <Button type="primary" size="large" icon={<IconSearch />}>
              浏览陪玩师
            </Button>
            <Button size="large" icon={<IconShoppingCart />}>
              我的订单
            </Button>
          </Space>
        </div>
      </Card>

      {/* 热门游戏 */}
      <div style={{ marginTop: 24 }}>
        <Title heading={4} style={{ marginBottom: 16 }}>
          热门游戏
        </Title>
        <Row gutter={[16, 16]}>
          {[1, 2, 3, 4].map((item) => (
            <Col key={item} xs={24} sm={12} md={6}>
              <Card
                hoverable
                cover={
                  <div
                    style={{
                      height: 200,
                      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                    }}
                  />
                }
              >
                <Card.Meta title={`游戏 ${item}`} description="热门陪玩游戏" />
              </Card>
            </Col>
          ))}
        </Row>
      </div>
    </div>
  );
};
