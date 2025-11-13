import { useEffect, useMemo, useState } from 'react';

import { feedApi, notificationApi } from '../../services/api';
import type { Feed, CreateFeedPayload } from '../../types/feed';
import type { NotificationEvent } from '../../types/notification';

const MAX_IMAGES = 9;

const parseImageInput = (raw?: string): CreateFeedPayload['images'] => {
  if (!raw) return [];
  return raw
    .split('\n')
    .map((item) => item.trim())
    .filter(Boolean)
    .map((url) => ({ url }));
};

const formatDateTime = (value?: string) => {
  if (!value) return '';
  try {
    return new Date(value).toLocaleString('zh-CN', { hour12: false });
  } catch (error) {
    return value;
  }
};

export const CommunityPage = () => {
  const [feeds, setFeeds] = useState<Feed[]>([]);
  const [loading, setLoading] = useState(false);
  const [composerOpen, setComposerOpen] = useState(false);
  const [nextCursor, setNextCursor] = useState<string | undefined>();
  const [publishing, setPublishing] = useState(false);
  const [likeSet, setLikeSet] = useState<Set<number>>(new Set());
  const [notifications, setNotifications] = useState<NotificationEvent[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [notificationOpen, setNotificationOpen] = useState(false);
  const [notificationLoading, setNotificationLoading] = useState(false);
  const [formContent, setFormContent] = useState('');
  const [formImages, setFormImages] = useState('');
  const [formVisibility, setFormVisibility] = useState<'public' | 'followers' | 'private'>('public');
  const [feedback, setFeedback] = useState<string | null>(null);

  const hasMore = useMemo(() => Boolean(nextCursor), [nextCursor]);

  const loadFeeds = async (cursor?: string, append = false) => {
    setLoading(true);
    try {
      const result = await feedApi.list({ cursor, limit: 10 });
      setFeeds((prev) => (append ? [...prev, ...result.items] : result.items));
      setNextCursor(result.nextCursor);
    } catch (error) {
      setFeedback((error as Error).message || '加载动态失败');
    } finally {
      setLoading(false);
    }
  };

  const refreshUnread = async () => {
    try {
      const data = await notificationApi.unreadCount();
      setUnreadCount(data.unread ?? 0);
    } catch (error) {
      // 静默处理，避免干扰用户
    }
  };

  const openNotifications = async () => {
    setNotificationOpen(true);
    setNotificationLoading(true);
    try {
      const result = await notificationApi.list({ page: 1, pageSize: 10, unread: false });
      setNotifications(result.items);
    } catch (error) {
      setFeedback((error as Error).message || '加载通知失败');
    } finally {
      setNotificationLoading(false);
    }
  };

  const handleMarkAllRead = async () => {
    if (!notifications.length) return;
    try {
      await notificationApi.markRead({ ids: notifications.map((n) => n.id) });
      setNotifications((prev) => prev.map((n) => ({ ...n, readAt: n.readAt ?? new Date().toISOString() })));
      refreshUnread();
    } catch (error) {
      setFeedback((error as Error).message || '标记已读失败');
    }
  };

  const handlePublish = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    try {
      if (!formContent.trim()) {
        setFeedback('动态内容不能为空');
        return;
      }
      const images = parseImageInput(formImages);
      if (images.length > MAX_IMAGES) {
        setFeedback(`单次最多上传 ${MAX_IMAGES} 张图片`);
        return;
      }
      const payload: CreateFeedPayload = {
        content: formContent,
        visibility: formVisibility,
        images,
      };
      setPublishing(true);
      await feedApi.create(payload);
      setFeedback('动态已提交，等待审核');
      setComposerOpen(false);
      setFormContent('');
      setFormImages('');
      loadFeeds();
    } catch (error) {
      if ((error as Error).message) {
        setFeedback((error as Error).message);
      }
    } finally {
      setPublishing(false);
    }
  };

  const handleReport = async (feed: Feed) => {
    const confirmed = window.confirm('确认举报该动态？我们会在 10 分钟内处理。');
    if (!confirmed) return;
    try {
      await feedApi.report(feed.id, { reason: '用户主动举报' });
      setFeedback('已提交举报');
    } catch (error) {
      setFeedback((error as Error).message || '举报失败');
    }
  };

  const toggleLike = (feedId: number) => {
    setLikeSet((prev) => {
      const next = new Set(prev);
      if (next.has(feedId)) {
        next.delete(feedId);
      } else {
        next.add(feedId);
      }
      return next;
    });
  };

  useEffect(() => {
    loadFeeds();
    refreshUnread();
    const timer = window.setInterval(refreshUnread, 15000);
    return () => window.clearInterval(timer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="community-page" style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem', width: '100%' }}>
      <div className="community-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2 style={{ margin: 0 }}>社区动态</h2>
        <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
          <button type="button" onClick={openNotifications}>
            通知中心 {unreadCount > 0 ? `(${unreadCount})` : ''}
          </button>
          <button type="button" onClick={() => setComposerOpen((prev) => !prev)} disabled={publishing}>
            {composerOpen ? '收起发布' : '发布动态'}
          </button>
        </div>
      </div>

      {feedback ? <div role="status">{feedback}</div> : null}

      {composerOpen ? (
        <form onSubmit={handlePublish} style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem', border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
          <label>
            <span>动态内容</span>
            <textarea value={formContent} onChange={(e) => setFormContent(e.target.value)} rows={4} placeholder='分享你的日常，最多500字' required />
          </label>
          <label>
            <span>图片地址（每行一条，最多9张）</span>
            <textarea value={formImages} onChange={(e) => setFormImages(e.target.value)} rows={3} placeholder='https://example.com/image1.png' />
          </label>
          <label>
            <span>可见范围</span>
            <select value={formVisibility} onChange={(e) => setFormVisibility(e.target.value as typeof formVisibility)}>
              <option value='public'>public</option>
              <option value='followers'>followers</option>
              <option value='private'>private</option>
            </select>
          </label>
          <button type="submit" disabled={publishing}>{publishing ? '发布中...' : '确认发布'}</button>
        </form>
      ) : null}

      <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
        {loading ? <div>加载中...</div> : null}
        {!loading && feeds.length === 0 ? <div>暂无动态</div> : null}
        {feeds.map((item) => (
          <div key={item.id} style={{ border: '1px solid #e5e7eb', borderRadius: '8px', padding: '1rem' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem', color: '#555' }}>
              <span>作者 #{item.authorId}</span>
              <span>{formatDateTime(item.createdAt)}</span>
            </div>
            <div style={{ marginBottom: '0.5rem' }}>{item.content}</div>
            {item.images?.length ? (
              <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap', marginBottom: '0.5rem' }}>
                {item.images.map((img) => (
                  <img key={img.url} src={img.url} alt='feed' style={{ width: '96px', height: '96px', objectFit: 'cover', borderRadius: '6px' }} />
                ))}
              </div>
            ) : null}
            <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
              <button type="button" onClick={() => toggleLike(item.id)}>
                {likeSet.has(item.id) ? '已赞' : '点赞'}
              </button>
              <button type="button" onClick={() => handleReport(item)}>举报</button>
              <span style={{ fontSize: '0.85rem', color: '#666' }}>审核状态：{item.moderationStatus}</span>
            </div>
          </div>
        ))}
      </div>

      {hasMore ? (
        <button type="button" onClick={() => loadFeeds(nextCursor, true)} disabled={loading}>
          加载更多
        </button>
      ) : null}

      {notificationOpen ? (
        <div style={{ border: '1px solid #d1d5db', borderRadius: '8px', padding: '1rem' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.75rem' }}>
            <h3 style={{ margin: 0 }}>通知中心</h3>
            <div style={{ display: 'flex', gap: '0.5rem' }}>
              <button type="button" onClick={() => setNotificationOpen(false)}>关闭</button>
              <button type="button" onClick={handleMarkAllRead} disabled={!notifications.length}>标记全部已读</button>
            </div>
          </div>
          {notificationLoading ? (
            <div>加载中...</div>
          ) : notifications.length === 0 ? (
            <div>暂无通知</div>
          ) : (
            <ul style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem', padding: 0 }}>
              {notifications.map((item) => (
                <li key={item.id} style={{ listStyle: 'none', borderBottom: '1px solid #eee', paddingBottom: '0.5rem' }}>
                  <div style={{ fontWeight: 600 }}>{item.title}</div>
                  <div style={{ margin: '0.25rem 0' }}>{item.message}</div>
                  <div style={{ fontSize: '0.85rem', color: '#666' }}>{formatDateTime(item.createdAt)}</div>
                </li>
              ))}
            </ul>
          )}
        </div>
      ) : null}
    </div>
  );
};

export default CommunityPage;
