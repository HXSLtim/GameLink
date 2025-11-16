/**
 * 下单页面 - Discord风格
 */

import { useState, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { DiscordLayout } from '@/components/Layout';
import styles from './CreateOrder.module.less';

/**
 * 服务项目接口
 */
interface ServiceItem {
  id: number;
  name: string;
  description: string;
  pricePerHour: number;
}

/**
 * 陪玩师信息接口
 */
interface PlayerInfo {
  id: number;
  avatar: string;
  nickname: string;
  gameName: string;
  rating: number;
  reviewCount: number;
}

/**
 * 下单页面组件
 */
export const CreateOrder: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const playerId = searchParams.get('playerId');
  const serviceId = searchParams.get('serviceId');

  const [loading, setLoading] = useState(true);
  const [player, setPlayer] = useState<PlayerInfo | null>(null);
  const [services, setServices] = useState<ServiceItem[]>([]);
  const [selectedServiceId, setSelectedServiceId] = useState<number | null>(
    serviceId ? Number(serviceId) : null
  );
  const [duration, setDuration] = useState(1);
  const [note, setNote] = useState('');

  useEffect(() => {
    // Mock数据加载（后续替换为真实API调用）
    setLoading(true);
    setTimeout(() => {
      setPlayer({
        id: Number(playerId),
        avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=player${playerId}`,
        nickname: '王者大神',
        gameName: '王者荣耀',
        rating: 4.8,
        reviewCount: 256,
      });
      setServices([
        {
          id: 1,
          name: '上分陪玩',
          description: '带你快速上分，保证胜率',
          pricePerHour: 50,
        },
        {
          id: 2,
          name: '技术教学',
          description: '一对一教学，提升游戏技术',
          pricePerHour: 80,
        },
        {
          id: 3,
          name: '娱乐陪玩',
          description: '轻松娱乐，开心游戏',
          pricePerHour: 40,
        },
      ]);
      setLoading(false);
    }, 500);
  }, [playerId]);

  const selectedService = services.find((s) => s.id === selectedServiceId);
  const totalPrice = selectedService ? selectedService.pricePerHour * duration : 0;

  const handleSubmit = () => {
    if (!selectedServiceId) {
      alert('请选择服务项目');
      return;
    }

    // 跳转到支付页面
    navigate(
      `/user/payment?playerId=${playerId}&serviceId=${selectedServiceId}&duration=${duration}&totalPrice=${totalPrice}`
    );
  };

  const renderStars = (rating: number) => {
    const stars = [];
    for (let i = 0; i < 5; i++) {
      stars.push(
        <span key={i} className={i < rating ? styles.starFull : styles.starEmpty}>
          ★
        </span>
      );
    }
    return stars;
  };

  if (loading) {
    return (
      <DiscordLayout showMemberPanel={false}>
        <div className={styles.loading}>
          <div className={styles.spinner} />
          <p>加载中...</p>
        </div>
      </DiscordLayout>
    );
  }

  if (!player) {
    return (
      <DiscordLayout showMemberPanel={false}>
        <div className={styles.error}>
          <p>❌ 陪玩师不存在</p>
          <button onClick={() => navigate('/user/players')}>返回列表</button>
        </div>
      </DiscordLayout>
    );
  }

  return (
    <DiscordLayout
      showMemberPanel={true}
      memberPanel={
        <div className={styles.sidebar}>
          <h3 className={styles.sidebarTitle}>订单摘要</h3>

          {/* 陪玩师信息 */}
          <div className={styles.playerSummary}>
            <img src={player.avatar} alt={player.nickname} className={styles.playerAvatar} />
            <div className={styles.playerInfo}>
              <div className={styles.playerName}>{player.nickname}</div>
              <div className={styles.playerGame}>{player.gameName}</div>
              <div className={styles.playerRating}>
                {renderStars(player.rating)}
                <span className={styles.ratingText}>{player.rating}</span>
              </div>
            </div>
          </div>

          {/* 订单详情 */}
          {selectedService && (
            <>
              <div className={styles.divider} />
              <div className={styles.orderDetails}>
                <div className={styles.detailRow}>
                  <span className={styles.detailLabel}>服务项目</span>
                  <span className={styles.detailValue}>{selectedService.name}</span>
                </div>
                <div className={styles.detailRow}>
                  <span className={styles.detailLabel}>单价</span>
                  <span className={styles.detailValue}>¥{selectedService.pricePerHour}/小时</span>
                </div>
                <div className={styles.detailRow}>
                  <span className={styles.detailLabel}>时长</span>
                  <span className={styles.detailValue}>{duration}小时</span>
                </div>
              </div>

              <div className={styles.divider} />

              {/* 总价 */}
              <div className={styles.totalPrice}>
                <span className={styles.totalLabel}>总计</span>
                <div className={styles.totalValue}>
                  <span className={styles.priceSymbol}>¥</span>
                  <span className={styles.priceAmount}>{totalPrice}</span>
                </div>
              </div>

              <button className={styles.submitButton} onClick={handleSubmit}>
                确认下单
              </button>
            </>
          )}
        </div>
      }
    >
      <div className={styles.createOrder}>
        {/* 页面标题 */}
        <div className={styles.header}>
          <button className={styles.backButton} onClick={() => navigate(-1)}>
            ← 返回
          </button>
          <h1 className={styles.title}>创建订单</h1>
        </div>

        {/* 服务项目选择 */}
        <div className={styles.section}>
          <h2 className={styles.sectionTitle}>
            选择服务项目 <span className={styles.required}>*</span>
          </h2>
          <div className={styles.serviceList}>
            {services.map((service) => (
              <div
                key={service.id}
                className={`${styles.serviceCard} ${
                  selectedServiceId === service.id ? styles.selected : ''
                }`}
                onClick={() => setSelectedServiceId(service.id)}
              >
                <div className={styles.serviceHeader}>
                  <div className={styles.serviceName}>{service.name}</div>
                  <div className={styles.servicePrice}>
                    ¥{service.pricePerHour}
                    <span className={styles.priceUnit}>/小时</span>
                  </div>
                </div>
                <div className={styles.serviceDesc}>{service.description}</div>
                {selectedServiceId === service.id && (
                  <div className={styles.selectedBadge}>✓ 已选择</div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* 时长选择 */}
        <div className={styles.section}>
          <h2 className={styles.sectionTitle}>
            选择时长 <span className={styles.required}>*</span>
          </h2>
          <div className={styles.durationOptions}>
            {[1, 2, 3, 4, 5].map((hours) => (
              <button
                key={hours}
                className={`${styles.durationButton} ${
                  duration === hours ? styles.active : ''
                }`}
                onClick={() => setDuration(hours)}
              >
                {hours}小时
              </button>
            ))}
          </div>
          <div className={styles.customDuration}>
            <label>自定义时长：</label>
            <input
              type="number"
              min="1"
              max="24"
              value={duration}
              onChange={(e) => setDuration(Math.max(1, Math.min(24, Number(e.target.value))))}
              className={styles.durationInput}
            />
            <span>小时</span>
          </div>
        </div>

        {/* 备注 */}
        <div className={styles.section}>
          <h2 className={styles.sectionTitle}>备注（可选）</h2>
          <textarea
            className={styles.noteTextarea}
            placeholder="有什么特殊要求或想对陪玩师说的话..."
            value={note}
            onChange={(e) => setNote(e.target.value)}
            maxLength={200}
            rows={4}
          />
          <div className={styles.charCount}>{note.length}/200</div>
        </div>

        {/* 移动端提交按钮 */}
        <div className={styles.mobileSubmit}>
          <div className={styles.mobilePrice}>
            <span className={styles.mobilePriceLabel}>总计</span>
            <span className={styles.mobilePriceValue}>¥{totalPrice}</span>
          </div>
          <button
            className={styles.mobileSubmitButton}
            onClick={handleSubmit}
            disabled={!selectedServiceId}
          >
            确认下单
          </button>
        </div>
      </div>
    </DiscordLayout>
  );
};

export default CreateOrder;
