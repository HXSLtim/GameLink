import React, { useEffect, useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button, Card, Input, Pagination, Table, Tag, message } from '../../components';
import { assignmentApi } from '../../services/api/assignment';
import { orderApi } from '../../services/api/order';
import type {
  AssignmentCandidate,
  AssignmentDispute,
  PendingAssignment,
} from '../../types/assignment';
import type { OrderDetail } from '../../types/order';
import { formatDateTime, formatOrderStatus, getOrderStatusColor } from '../../utils/formatters';
import styles from './AssignmentsWorkbench.module.less';

const PAGE_SIZE = 10;

const ASSIGNMENT_SOURCE_LABELS: Record<string, string> = {
  manual: '人工指派',
  system_recommendation: '系统推荐',
  rollback: '已回退',
  unknown: '未指定',
};

const DISPUTE_STATUS_LABELS: Record<string, string> = {
  none: '无争议',
  pending: '待处理',
  in_mediation: '调解中',
  resolved: '已解决',
};

const DISPUTE_RESOLUTION_LABELS: Record<string, string> = {
  none: '未处理',
  refund: '退款',
  reassign: '重派',
  reject: '驳回',
};

type AssignPayload = {
  orderId: number;
  playerId: number;
  source?: string;
};

type CancelPayload = {
  orderId: number;
  reason?: string;
};

type MediatePayload = {
  orderId: number;
  resolution: string;
  note?: string;
  refundAmount?: number;
  reassignPlayerId?: number;
};

const formatSlaTimer = (seconds: number): string => {
  const safe = Math.max(0, Math.floor(seconds));
  const hours = Math.floor(safe / 3600);
  const minutes = Math.floor((safe % 3600) / 60);
  const secs = safe % 60;
  if (hours > 0) {
    return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
  }
  return `${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
};

const computeRemainingSeconds = (assignment: PendingAssignment, now: number): number => {
  if (assignment.slaDeadline) {
    const deadline = new Date(assignment.slaDeadline).getTime();
    if (!Number.isNaN(deadline)) {
      return Math.max(0, Math.floor((deadline - now) / 1000));
    }
  }
  return Math.max(0, assignment.slaRemainingSeconds ?? 0);
};

const actionDescriptions: Record<string, string> = {
  assign_player: '指派陪玩师',
  assign_rollback: '回退指派',
  create_dispute: '用户发起争议',
  mediate_dispute: '客服裁决争议',
};

const formatDisputeLabel = (dispute: AssignmentDispute): string => {
  const status = DISPUTE_STATUS_LABELS[dispute.status] || dispute.status;
  const resolution = dispute.resolution ? DISPUTE_RESOLUTION_LABELS[dispute.resolution] || dispute.resolution : '';
  return resolution ? `${status} · ${resolution}` : status;
};

const parseRefundInput = (input: string): number | undefined => {
  if (!input) return undefined;
  const value = Number.parseFloat(input);
  if (Number.isNaN(value)) return undefined;
  return Math.round(value * 100);
};

export const AssignmentsWorkbench: React.FC = () => {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [selectedOrderId, setSelectedOrderId] = useState<number | null>(null);
  const [now, setNow] = useState(() => Date.now());
  const [cancelReason, setCancelReason] = useState('');
  const [resolution, setResolution] = useState('refund');
  const [refundAmount, setRefundAmount] = useState('');
  const [reassignPlayerId, setReassignPlayerId] = useState('');
  const [mediateNote, setMediateNote] = useState('');

  useEffect(() => {
    const timer = window.setInterval(() => setNow(Date.now()), 1000);
    return () => window.clearInterval(timer);
  }, []);

  const pendingQuery = useQuery({
    queryKey: ['pendingAssignments', page, PAGE_SIZE],
    queryFn: () => assignmentApi.getPending({ page, pageSize: PAGE_SIZE }),
    refetchInterval: 60_000,
  });

  const pendingList = useMemo(
    () => pendingQuery.data?.list ?? [],
    [pendingQuery.data?.list]
  );

  useEffect(() => {
    if (pendingList.length === 0) {
      setSelectedOrderId(null);
      return;
    }
    const exists = selectedOrderId && pendingList.some((item) => item.orderId === selectedOrderId);
    if (!selectedOrderId || !exists) {
      setSelectedOrderId(pendingList[0].orderId);
    }
  }, [pendingList, selectedOrderId]);

  const selectedAssignment = useMemo(() => {
    if (!selectedOrderId) return undefined;
    return pendingList.find((item) => item.orderId === selectedOrderId);
  }, [pendingList, selectedOrderId]);

  const orderDetailQuery = useQuery<OrderDetail | undefined>({
    queryKey: ['orderDetail', selectedOrderId],
    queryFn: () => (selectedOrderId ? orderApi.getDetail(selectedOrderId) : Promise.resolve(undefined)),
    enabled: Boolean(selectedOrderId),
  });

  const orderLogsQuery = useQuery({
    queryKey: ['orderLogs', selectedOrderId],
    queryFn: () => (selectedOrderId ? orderApi.getLogs(selectedOrderId) : Promise.resolve([])),
    enabled: Boolean(selectedOrderId),
  });

  const candidatesQuery = useQuery<AssignmentCandidate[]>({
    queryKey: ['assignmentCandidates', selectedOrderId],
    queryFn: () => (selectedOrderId ? assignmentApi.getCandidates(selectedOrderId, { limit: 8 }) : Promise.resolve([])),
    enabled: Boolean(selectedOrderId),
    staleTime: 60_000,
  });

  const disputesQuery = useQuery<AssignmentDispute[]>({
    queryKey: ['assignmentDisputes', selectedOrderId],
    queryFn: () => (selectedOrderId ? assignmentApi.getDisputes(selectedOrderId) : Promise.resolve([])),
    enabled: Boolean(selectedOrderId),
    refetchInterval: 60_000,
  });

  const assignMutation = useMutation({
    mutationFn: ({ orderId, playerId, source }: AssignPayload) =>
      assignmentApi.assign(orderId, { playerId, source }),
    onSuccess: (_, variables) => {
      message.success('指派操作成功');
      queryClient.invalidateQueries({ queryKey: ['pendingAssignments'] });
      queryClient.invalidateQueries({ queryKey: ['assignmentCandidates', variables.orderId] });
      queryClient.invalidateQueries({ queryKey: ['orderDetail', variables.orderId] });
      queryClient.invalidateQueries({ queryKey: ['orderLogs', variables.orderId] });
      queryClient.invalidateQueries({ queryKey: ['assignmentDisputes', variables.orderId] });
    },
    onError: (err) => {
      message.error(err instanceof Error ? err.message : '指派失败');
    },
  });

  const cancelMutation = useMutation({
    mutationFn: ({ orderId, reason }: CancelPayload) =>
      assignmentApi.cancelAssignment(orderId, { reason }),
    onSuccess: (_, variables) => {
      message.success('指派已回退');
      setCancelReason('');
      queryClient.invalidateQueries({ queryKey: ['pendingAssignments'] });
      queryClient.invalidateQueries({ queryKey: ['orderDetail', variables.orderId] });
      queryClient.invalidateQueries({ queryKey: ['orderLogs', variables.orderId] });
    },
    onError: (err) => {
      message.error(err instanceof Error ? err.message : '回退失败');
    },
  });

  const mediateMutation = useMutation({
    mutationFn: ({ orderId, resolution: decision, note, refundAmount, reassignPlayerId }: MediatePayload) =>
      assignmentApi.mediate(orderId, {
        resolution: decision,
        note,
        refundAmountCents: refundAmount,
        reassignPlayerId,
      }),
    onSuccess: (_, variables) => {
      message.success('争议已处理');
      setMediateNote('');
      setRefundAmount('');
      setReassignPlayerId('');
      queryClient.invalidateQueries({ queryKey: ['orderDetail', variables.orderId] });
      queryClient.invalidateQueries({ queryKey: ['assignmentDisputes', variables.orderId] });
      queryClient.invalidateQueries({ queryKey: ['orderLogs', variables.orderId] });
      queryClient.invalidateQueries({ queryKey: ['pendingAssignments'] });
    },
    onError: (err) => {
      message.error(err instanceof Error ? err.message : '争议处理失败');
    },
  });

  const tableColumns = [
    {
      key: 'orderId',
      title: '订单ID',
      render: (_: unknown, record: PendingAssignment) => record.orderId,
      width: '90px',
    },
    {
      key: 'status',
      title: '订单状态',
      render: (_: unknown, record: PendingAssignment) => (
        <Tag color={getOrderStatusColor(record.status)}>{formatOrderStatus(record.status)}</Tag>
      ),
      width: '120px',
    },
    {
      key: 'source',
      title: '来源',
      render: (_: unknown, record: PendingAssignment) =>
        ASSIGNMENT_SOURCE_LABELS[record.assignmentSource || 'unknown'] || record.assignmentSource,
    },
    {
      key: 'createdAt',
      title: '创建时间',
      render: (_: unknown, record: PendingAssignment) => formatDateTime(record.createdAt),
    },
    {
      key: 'sla',
      title: 'SLA 倒计时',
      render: (_: unknown, record: PendingAssignment) => {
        const remaining = computeRemainingSeconds(record, now);
        const overdue = remaining <= 0;
        return (
          <Tag color={overdue ? 'error' : remaining < 300 ? 'warning' : 'info'}>
            {overdue ? '已超时' : formatSlaTimer(remaining)}
          </Tag>
        );
      },
      width: '130px',
    },
  ];

  const timelineItems = useMemo(() => {
    const logs = (orderLogsQuery.data ?? []).map((log: any) => ({
      id: `log-${log.id}`,
      createdAt: log.createdAt,
      title: actionDescriptions[log.action] || log.action,
      note: log.reason,
      traceId: log.traceId,
      actorId: log.actorUserId,
    }));
    const disputes = (disputesQuery.data ?? []).map((dispute) => ({
      id: `dispute-${dispute.id}`,
      createdAt: dispute.createdAt,
      title: formatDisputeLabel(dispute),
      note: dispute.reason,
      traceId: dispute.traceId,
    }));
    return [...logs, ...disputes].sort((a, b) => {
      const aTime = a.createdAt ? new Date(a.createdAt).getTime() : 0;
      const bTime = b.createdAt ? new Date(b.createdAt).getTime() : 0;
      return bTime - aTime;
    });
  }, [orderLogsQuery.data, disputesQuery.data]);

  const handleAssign = (candidate: AssignmentCandidate) => {
    if (!selectedOrderId) return;
    assignMutation.mutate({
      orderId: selectedOrderId,
      playerId: candidate.playerId,
      source: candidate.source,
    });
  };

  const handleCancelAssign = () => {
    if (!selectedOrderId) return;
    cancelMutation.mutate({ orderId: selectedOrderId, reason: cancelReason.trim() || undefined });
  };

  const handleMediate = () => {
    if (!selectedOrderId) return;
    const payload: MediatePayload = {
      orderId: selectedOrderId,
      resolution,
      note: mediateNote.trim() || undefined,
    };
    if (resolution === 'refund') {
      const cents = parseRefundInput(refundAmount);
      if (cents === undefined) {
        message.error('请输入有效的退款金额');
        return;
      }
      payload.refundAmount = cents;
    }
    if (resolution === 'reassign') {
      const targetId = Number.parseInt(reassignPlayerId, 10);
      if (Number.isNaN(targetId)) {
        message.error('请输入有效的陪玩师ID');
        return;
      }
      payload.reassignPlayerId = targetId;
    }
    mediateMutation.mutate(payload);
  };

  return (
    <div className={styles.container}>
      <div className={styles.listCard}>
        <Card
          title="待指派订单"
          extra={
            <div className={styles.filters}>
              <Button
                variant="secondary"
                size="small"
                onClick={() => pendingQuery.refetch()}
                loading={pendingQuery.isFetching}
              >
                刷新列表
              </Button>
            </div>
          }
        >
          {pendingList.length === 0 && !pendingQuery.isLoading ? (
            <div className={styles.emptyState}>暂无待处理的指派订单</div>
          ) : (
            <div className={styles.tableWrapper}>
              <Table<PendingAssignment>
                columns={tableColumns}
                dataSource={pendingList}
                loading={pendingQuery.isLoading}
                rowKey={(record) => String(record.orderId)}
                onRow={(record) => ({
                  onClick: () => setSelectedOrderId(record.orderId),
                  style: { cursor: 'pointer' },
                })}
                rowClassName={(record) =>
                  record.orderId === selectedOrderId ? styles.tableRowActive : ''
                }
              />
            </div>
          )}
          <Pagination
            current={page}
            total={pendingQuery.data?.total ?? 0}
            pageSize={PAGE_SIZE}
            onChange={setPage}
          />
        </Card>
      </div>

      <div className={styles.detailPane}>
        <Card
          title="指派详情"
          extra={
            selectedAssignment ? (
              <Tag color={selectedAssignment?.isOverdue ? 'error' : 'info'}>
                {selectedAssignment?.isOverdue
                  ? 'SLA 已超时'
                  : `SLA 剩余 ${formatSlaTimer(computeRemainingSeconds(selectedAssignment, now))}`}
              </Tag>
            ) : null
          }
        >
          {!selectedOrderId || !orderDetailQuery.data ? (
            <div className={styles.emptyState}>请选择左侧订单查看详情</div>
          ) : (
            <>
              {selectedAssignment?.isOverdue && (
                <div className={styles.alertCard}>
                  <span>⚠</span>
                  <span>该订单已超过 30 分钟 SLA，已触发升级告警</span>
                </div>
              )}
              <div className={styles.detailHeader}>
                <div>
                  <h3>订单 #{orderDetailQuery.data.id}</h3>
                  <div>{orderDetailQuery.data.title || '未命名订单'}</div>
                </div>
                <Tag color={getOrderStatusColor(orderDetailQuery.data.status)}>
                  {formatOrderStatus(orderDetailQuery.data.status)}
                </Tag>
              </div>
              <div className={styles.detailMeta}>
                <div className={styles.metaItem}>
                  <span className={styles.metaLabel}>客户</span>
                  <span>{orderDetailQuery.data.user?.name ?? `用户 #${orderDetailQuery.data.userId}`}</span>
                </div>
                <div className={styles.metaItem}>
                  <span className={styles.metaLabel}>当前陪玩师</span>
                  <span>
                    {orderDetailQuery.data.player
                      ? `${orderDetailQuery.data.player.nickname ?? '陪玩师'} (#${orderDetailQuery.data.player.id})`
                      : '未指派'}
                  </span>
                </div>
                <div className={styles.metaItem}>
                  <span className={styles.metaLabel}>指派来源</span>
                  <span>
                    {ASSIGNMENT_SOURCE_LABELS[orderDetailQuery.data.assignmentSource || 'unknown'] ||
                      orderDetailQuery.data.assignmentSource ||
                      '未指定'}
                  </span>
                </div>
                <div className={styles.metaItem}>
                  <span className={styles.metaLabel}>争议状态</span>
                  <span>
                    {DISPUTE_STATUS_LABELS[orderDetailQuery.data.disputeStatus || 'none'] ||
                      orderDetailQuery.data.disputeStatus ||
                      '无'}
                  </span>
                </div>
                <div className={styles.metaItem}>
                  <span className={styles.metaLabel}>下单时间</span>
                  <span>{formatDateTime(orderDetailQuery.data.createdAt)}</span>
                </div>
                <div className={styles.metaItem}>
                  <span className={styles.metaLabel}>更新时间</span>
                  <span>{formatDateTime(orderDetailQuery.data.updatedAt)}</span>
                </div>
              </div>

              <div>
                <h4>争议记录</h4>
                <div className={styles.timeline}>
                  {disputesQuery.isLoading && <div>加载中...</div>}
                  {!disputesQuery.isLoading && disputesQuery.data?.length === 0 && (
                    <div className={styles.emptyState}>暂无争议</div>
                  )}
                  {(disputesQuery.data ?? []).map((dispute) => {
                    const deadlineText = dispute.responseDeadline ? formatDateTime(dispute.responseDeadline) : undefined;
                    return (
                      <div key={dispute.id} className={styles.timelineItem}>
                        <div className={styles.timelineTitle}>
                          <span>{formatDisputeLabel(dispute)}</span>
                          <span>{formatDateTime(dispute.createdAt)}</span>
                        </div>
                        <div className={styles.timelineMeta}>
                          <span>发起人: {dispute.raisedBy}</span>
                          {deadlineText && <span>响应截止: {deadlineText}</span>}
                          {dispute.traceId && <span>traceId: {dispute.traceId}</span>}
                        </div>
                        <div>{dispute.reason}</div>
                        {dispute.resolutionNote && <div>备注：{dispute.resolutionNote}</div>}
                        {typeof dispute.refundAmountCents === 'number' && dispute.refundAmountCents > 0 && (
                          <div>退款金额：¥{(dispute.refundAmountCents / 100).toFixed(2)}</div>
                        )}
                      </div>
                    );
                  })}
                </div>
              </div>

              <div>
                <h4>时间线</h4>
                <div className={styles.timeline}>
                  {timelineItems.length === 0 && <div className={styles.emptyState}>暂无日志</div>}
                  {timelineItems.map((item) => (
                    <div key={item.id} className={styles.timelineItem}>
                      <div className={styles.timelineTitle}>
                        <span>{item.title}</span>
                        <span>{item.createdAt ? formatDateTime(item.createdAt) : '-'}</span>
                      </div>
                      <div className={styles.timelineMeta}>
                        {item.traceId && <span>traceId: {item.traceId}</span>}
                        {item.actorId && <span>操作人: #{item.actorId}</span>}
                      </div>
                      {item.note && <div>{item.note}</div>}
                    </div>
                  ))}
                </div>
              </div>
            </>
          )}
        </Card>

        <Card title="推荐候选">
          {candidatesQuery.isLoading ? (
            <div className={styles.emptyState}>加载候选中...</div>
          ) : (candidatesQuery.data ?? []).length === 0 ? (
            <div className={styles.emptyState}>暂无候选，请稍后刷新</div>
          ) : (
            <div className={styles.candidates}>
              {(candidatesQuery.data ?? []).map((candidate) => (
                <div key={candidate.playerId} className={styles.candidateCard}>
                  <div className={styles.candidateInfo}>
                    <strong>{candidate.nickname || `陪玩师 #${candidate.playerId}`}</strong>
                    <div className={styles.candidateMeta}>
                      来源：{candidate.source === 'recommendation' ? '系统推荐' : '陪玩师列表'}
                      {typeof candidate.score === 'number' && candidate.score > 0 ? ` · 评分 ${candidate.score.toFixed(2)}` : ''}
                    </div>
                    {candidate.reason && <div className={styles.candidateMeta}>推荐理由：{candidate.reason}</div>}
                    {typeof candidate.hourlyRateCents === 'number' && candidate.hourlyRateCents > 0 && (
                      <div className={styles.candidateMeta}>
                        标准价：¥{(candidate.hourlyRateCents / 100).toFixed(2)} / 小时
                      </div>
                    )}
                  </div>
                  <Button
                    size="small"
                    onClick={() => handleAssign(candidate)}
                    loading={assignMutation.isPending}
                  >
                    指派
                  </Button>
                </div>
              ))}
            </div>
          )}
        </Card>

        <Card title="调解操作">
          {!selectedOrderId ? (
            <div className={styles.emptyState}>请选择订单后操作</div>
          ) : (
            <>
              <div className={styles.actionBar}>
                <Input
                  placeholder="填写指派回退原因（选填）"
                  value={cancelReason}
                  onChange={(e) => setCancelReason(e.target.value)}
                  size="small"
                />
                <Button
                  variant="outlined"
                  size="small"
                  onClick={handleCancelAssign}
                  disabled={cancelMutation.isPending}
                >
                  回退指派
                </Button>
              </div>

              <div className={styles.formField}>
                <label>裁决结果</label>
                <select value={resolution} onChange={(e) => setResolution(e.target.value)}>
                  <option value="refund">退款</option>
                  <option value="reassign">重派</option>
                  <option value="reject">驳回</option>
                </select>
              </div>
              {resolution === 'refund' && (
                <div className={styles.formField}>
                  <label>退款金额（元）</label>
                  <Input
                    placeholder="例如 50.00"
                    value={refundAmount}
                    onChange={(e) => setRefundAmount(e.target.value)}
                    size="small"
                  />
                </div>
              )}
              {resolution === 'reassign' && (
                <div className={styles.formField}>
                  <label>重新指派陪玩师ID</label>
                  <Input
                    placeholder="输入陪玩师ID"
                    value={reassignPlayerId}
                    onChange={(e) => setReassignPlayerId(e.target.value)}
                    size="small"
                  />
                </div>
              )}
              <div className={styles.formField}>
                <label>备注</label>
                <textarea
                  className={styles.textarea}
                  placeholder="填写调解说明（可选）"
                  value={mediateNote}
                  onChange={(e) => setMediateNote(e.target.value)}
                />
              </div>
              <Button
                onClick={handleMediate}
                loading={mediateMutation.isPending}
              >
                提交裁决
              </Button>
            </>
          )}
        </Card>
      </div>
    </div>
  );
};

export default AssignmentsWorkbench;
