import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Button, Tag } from '../../components';
import { gameApi } from '../../services/api/game';
import type { Game, GameCategory } from '../../types/game';
import { GAME_CATEGORY_TEXT } from '../../types/game';
import { formatDateTime } from '../../utils/formatters';
import styles from './GameDetail.module.less';

// 图标组件
const ArrowLeftIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <line x1="19" y1="12" x2="5" y2="12" strokeWidth="2" strokeLinecap="round" />
    <polyline
      points="12 19 5 12 12 5"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

const EditIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

const TrashIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <polyline points="3 6 5 6 21 6" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path
      d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

export const GameDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [game, setGame] = useState<Game | null>(null);

  // 加载游戏详情
  useEffect(() => {
    const loadGameDetail = async () => {
      if (!id) {
        setError('游戏ID无效');
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        setError(null);
        const data = await gameApi.getDetail(Number(id));
        setGame(data);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '加载游戏详情失败';
        setError(errorMessage);
        console.error('加载游戏详情失败:', err);
      } finally {
        setLoading(false);
      }
    };

    loadGameDetail();
  }, [id]);

  // 删除游戏
  const handleDelete = async () => {
    if (!id || !game) return;

    if (!window.confirm(`确定要删除游戏 "${game.name}" 吗？`)) {
      return;
    }

    try {
      await gameApi.delete(Number(id));
      alert('游戏删除成功');
      navigate('/games');
    } catch (err) {
      console.error('删除游戏失败:', err);
      alert('删除游戏失败: ' + (err instanceof Error ? err.message : '未知错误'));
    }
  };

  // 获取分类颜色
  const getCategoryColor = (category: string): string => {
    const colorMap: Record<string, string> = {
      moba: 'blue',
      fps: 'red',
      rpg: 'purple',
      strategy: 'green',
      sports: 'orange',
      racing: 'cyan',
      puzzle: 'pink',
      other: 'default',
    };
    return colorMap[category] || 'default';
  };

  // 加载中状态
  if (loading) {
    return (
      <div className={styles.container}>
        <p>加载中...</p>
      </div>
    );
  }

  // 错误状态
  if (error || !game) {
    return (
      <div className={styles.container}>
        <Card className={styles.errorCard}>
          <h2>游戏未找到</h2>
          <p>{error || `游戏 ID: ${id} 不存在`}</p>
          <Button onClick={() => navigate('/games')}>返回游戏列表</Button>
        </Card>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* 头部 */}
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <Button variant="outlined" onClick={() => navigate('/games')}>
            <ArrowLeftIcon />
            返回列表
          </Button>
          <h1 className={styles.title}>游戏详情</h1>
        </div>
        <div className={styles.headerRight}>
          <Tag color={getCategoryColor(game.category) as any}>
            {GAME_CATEGORY_TEXT[game.category as GameCategory] || game.category}
          </Tag>
        </div>
      </div>

      {/* 主要内容 */}
      <div className={styles.content}>
        {/* 基本信息 */}
        <Card className={styles.section}>
          <div className={styles.sectionHeader}>
            <h2 className={styles.sectionTitle}>基本信息</h2>
            <div className={styles.sectionActions}>
              <Button variant="outlined" onClick={() => console.log('编辑游戏')}>
                <EditIcon />
                编辑
              </Button>
              <Button variant="secondary" onClick={handleDelete}>
                <TrashIcon />
                删除
              </Button>
            </div>
          </div>

          <div className={styles.infoGrid}>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>游戏ID</span>
              <span className={styles.infoValue}>{game.id}</span>
            </div>

            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>游戏标识</span>
              <span className={styles.infoValue}>{game.key}</span>
            </div>

            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>游戏名称</span>
              <span className={styles.infoValue}>{game.name}</span>
            </div>

            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>分类</span>
              <span className={styles.infoValue}>
                <Tag color={getCategoryColor(game.category) as any}>
                  {GAME_CATEGORY_TEXT[game.category as GameCategory] || game.category}
                </Tag>
              </span>
            </div>

            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>图标URL</span>
              <span className={styles.infoValue}>{game.icon_url || '-'}</span>
            </div>

            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>创建时间</span>
              <span className={styles.infoValue}>{formatDateTime(game.created_at)}</span>
            </div>

            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>更新时间</span>
              <span className={styles.infoValue}>{formatDateTime(game.updated_at)}</span>
            </div>
          </div>

          {game.description && (
            <div className={styles.descriptionSection}>
              <h3 className={styles.subTitle}>游戏描述</h3>
              <p className={styles.description}>{game.description}</p>
            </div>
          )}
        </Card>

        {/* 游戏图标预览 */}
        {game.icon_url && (
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>图标预览</h2>
            <div className={styles.iconPreview}>
              <img src={game.icon_url} alt={game.name} className={styles.gameIcon} />
            </div>
          </Card>
        )}
      </div>
    </div>
  );
};
