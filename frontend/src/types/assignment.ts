import type { OrderStatus } from './order';

export interface PendingAssignment {
  orderId: number;
  userId: number;
  status: OrderStatus;
  assignmentSource: string;
  createdAt: string;
  slaDeadline: string;
  slaRemainingSeconds: number;
  isOverdue: boolean;
}

export interface AssignmentCandidate {
  playerId: number;
  nickname?: string;
  hourlyRateCents?: number;
  score?: number;
  source: string;
  reason?: string;
}

export interface AssignmentTimelineItem {
  id: number;
  status: string;
  resolution?: string;
  note?: string;
  refundAmountCents?: number;
  handledAt?: string;
  createdAt: string;
  raisedBy?: string;
}

export interface AssignmentDispute {
  id: number;
  orderId: number;
  raisedBy: string;
  raisedByUserId?: number;
  status: string;
  reason: string;
  evidenceUrls?: string[];
  resolution?: string;
  resolutionNote?: string;
  refundAmountCents?: number;
  handledAt?: string;
  handledById?: number;
  createdAt: string;
  responseDeadline?: string;
  respondedAt?: string;
  traceId?: string;
}

export interface CancelAssignmentRequest {
  reason?: string;
}

export interface MediateDisputeRequest {
  resolution: string;
  note?: string;
  refundAmountCents?: number;
  reassignPlayerId?: number;
}

export interface CreateDisputeRequest {
  reason: string;
  evidence?: string[];
}

export interface PendingAssignmentList {
  list: PendingAssignment[];
  total: number;
  page?: number;
  page_size?: number;
}
