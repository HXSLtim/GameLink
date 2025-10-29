import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Card, Button, Tag, Modal } from '../../components';
import { gameApi } from '../../services/api/game';
import type { GameDetail, UpdateGameRequest } from '../../types/game';
import { GAME_CATEGORY_TEXT, GAME_STATUS_TEXT, GAME_STATUS_COLOR } from '../../types/game';
import { formatCurrency, formatDateTime, formatRelativeTime } from '../../utils/formatters';
import { GameFormModal } from './GameFormModal';
import styles from './GameDetail.module.less';

const GameDetail: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [loading, setLoading] = useState(false);
  const [game, setGame] = useState<GameDetail | null>(null);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [deleteModalVisible, setDeleteModalVisible] = useState(false);

  // 加载游戏详情
  const loadGameDetail = async () => {
    if (!id) return;

    setLoading(true);
    try {
      const data = await gameApi.getDetail(Number(id));
      setGame(data);
    } catch (err) {
      console.error('加载游戏详情失败:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadGameDetail();
  }, [id]);

  // 编辑游戏
  const handleEdit = async (data: UpdateGameRequest) => {
    if (!id) return;

    try {
      await gameApi.update(Number(id), data);
      setEditModalVisible(false);
      await loadGameDetail();
    } catch (err) {
      console.error('更新游戏失败:', err);
      throw err;
    }
  };

  // 删除游戏
  const handleDelete = async () => {
    if (!id) return;

    try {
      await gameApi.delete(Number(id));
      setDeleteModalVisible(false);
      navigate('/games');
    } catch (err) {
      console.error('删除游戏失败:', err);
    }
  };

  if (loading) {
    return (
      <div className={styles.container}>
        <div className={styles.loading}>加载中...</div>
      </div>
    );
  }

  if (!game) {
    return (
      <div className={styles.container}>
        <div className={styles.error}>游戏不存在</div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* 头部 */}
      <div className={styles.header}>
        <Button variant="text" onClick={() => navigate('/games')} className={styles.backButton}>
          ← 返回列表
        </Button>
        <div className={styles.headerActions}>
          <Button variant="outlined" onClick={() => setEditModalVisible(true)}>
            编辑
          </Button>
          <Button variant="outlined" onClick={() => setDeleteModalVisible(true)}>
            删除
          </Button>
        </div>
      </div>

      {/* 基本信息卡片 */}
      <Card className={styles.infoCard}>
        <div className={styles.gameHeader}>
          {game.iconUrl && (
            <img src={game.iconUrl} alt={game.name} className={styles.gameIcon} />
          )}
          <div className={styles.gameInfo}>
            <div className={styles.gameName}>
              {game.name}
              {game.status && (
                <Tag color={GAME_STATUS_COLOR[game.status] as any}>
                  {GAME_STATUS_TEXT[game.status]}
                </Tag>
              )}
            </div>
            <div className={styles.gameMeta}>
              <span className={styles.metaItem}>
                <strong>游戏标识:</strong> {game.key}
              </span>
              <span className={styles.metaItem}>
                <strong>分类:</strong> {GAME_CATEGORY_TEXT[game.category as keyof typeof GAME_CATEGORY_TEXT] || game.category}
              </span>
              <span className={styles.metaItem}>
                <strong>ID:</strong> {game.id}
              </span>
            </div>
            {game.tags && game.tags.length > 0 && (
              <div className={styles.tags}>
                {game.tags.map((tag) => (
                  <Tag key={tag} className={styles.tag}>
                    {tag}
                  </Tag>
                ))}
              </div>
            )}
          </div>
        </div>

        {game.description && (
          <div className={styles.description}>
            <h3>游戏简介</h3>
            <p>{game.description}</p>
          </div>
        )}

        <div className={styles.timestamps}>
          <div className={styles.timestamp}>
            <span className={styles.timestampLabel}>创建时间:</span>
            <span className={styles.timestampValue}>
              {formatDateTime(game.createdAt)}
              <span className={styles.relative}>({formatRelativeTime(game.createdAt)})</span>
            </span>
          </div>
          <div className={styles.timestamp}>
            <span className={styles.timestampLabel}>更新时间:</span>
            <span className={styles.timestampValue}>
              {formatDateTime(game.updatedAt)}
              <span className={styles.relative}>({formatRelativeTime(game.updatedAt)})</span>
            </span>
          </div>
        </div>
      </Card>

      {/* 统计信息卡片 */}
      {game.stats && (
        <Card className={styles.statsCard}>
          <h3 className={styles.cardTitle}>统计数据</h3>
          <div className={styles.statsGrid}>
            <div className={styles.statItem}>
              <div className={styles.statLabel}>陪玩师数量</div>
              <div className={styles.statValue}>{game.stats.totalPlayers}</div>
            </div>
            <div className={styles.statItem}>
              <div className={styles.statLabel}>订单数量</div>
              <div className={styles.statValue}>{game.stats.totalOrders}</div>
            </div>
            <div className={styles.statItem}>
              <div className={styles.statLabel}>总收入</div>
              <div className={styles.statValue}>
                {formatCurrency(game.stats.totalRevenue)}
              </div>
            </div>
            <div className={styles.statItem}>
              <div className={styles.statLabel}>平均评分</div>
              <div className={styles.statValue}>
                {game.stats.avgRating.toFixed(1)} ⭐
              </div>
            </div>
          </div>
        </Card>
      )}

      {/* 热门陪玩师卡片 */}
      {game.topPlayers && game.topPlayers.length > 0 && (
        <Card className={styles.playersCard}>
          <h3 className={styles.cardTitle}>热门陪玩师</h3>
          <div className={styles.playersList}>
            {game.topPlayers.map((player) => (
              <div
                key={player.id}
                className={styles.playerItem}
                onClick={() => navigate(`/players/${player.id}`)}
              >
                {player.avatarUrl && (
                  <img src={player.avatarUrl} alt={player.nickname} className={styles.playerAvatar} />
                )}
                <div className={styles.playerInfo}>
                  <div className={styles.playerName}>{player.nickname}</div>
                  <div className={styles.playerRating}>⭐ {player.rating.toFixed(1)}</div>
                </div>
              </div>
            ))}
          </div>
        </Card>
      )}

      {/* 编辑Modal */}
      <GameFormModal
        visible={editModalVisible}
        game={game}
        onClose={() => setEditModalVisible(false)}
        onSubmit={handleEdit}
      />

      {/* 删除确认Modal */}
      <Modal
        visible={deleteModalVisible}
        title="确认删除"
        onClose={() => setDeleteModalVisible(false)}
        onOk={handleDelete}
        onCancel={() => setDeleteModalVisible(false)}
        okText="确定删除"
        cancelText="取消"
        width={400}
      >
        <p>确定要删除游戏 "{game.name}" 吗？此操作不可恢复。</p>
        {game.stats && game.stats.totalOrders > 0 && (
          <p className={styles.warning}>
            ⚠️ 该游戏已有 {game.stats.totalOrders} 个订单，删除可能影响历史数据。
          </p>
        )}
      </Modal>
    </div>
  );
};

export default GameDetail;
