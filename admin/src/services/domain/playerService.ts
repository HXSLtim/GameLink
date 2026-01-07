/**
 * Player Domain Service
 * Encapsulates all player-related business logic
 *
 * @module services/domain/playerService
 */

import {
  BaseService,
  type ServiceDependencies,
} from './base';
import {
  ServiceErrorCodes,
  type ServiceResult,
  type BatchResult,
  ServiceResultHelper,
} from '../utils';
import type {
  Player,
  PlayerQueryParams,
  CreatePlayerDto,
  UpdatePlayerDto,
  Order,
  CommissionRule,
} from '@/api/admin';

/**
 * Verification check result
 */
export interface VerificationCheckResult {
  allowed: boolean;
  reason?: string;
}

/**
 * Earnings calculation result
 */
export interface EarningsCalculation {
  /** Gross amount in cents (before commission) */
  grossAmount: number;
  /** Commission rate as decimal (0-1) */
  commissionRate: number;
  /** Commission amount in cents (platform takes) */
  commissionAmount: number;
  /** Net amount in cents (player receives) */
  netAmount: number;
}

/**
 * Player statistics
 */
export interface PlayerStatistics {
  /** Total earnings in cents (all time) */
  totalEarnings: number;
  /** Monthly earnings in cents (current month) */
  monthlyEarnings: number;
  /** Number of completed orders */
  completedOrders: number;
  /** Average rating (0-5) */
  averageRating: number;
  /** Total number of ratings */
  ratingCount: number;
}

/**
 * Skill tag validation result
 */
export interface SkillTagValidationResult {
  valid: boolean;
  invalidTags: string[];
}

/**
 * Player Service Interface
 */
export interface IPlayerService {
  // CRUD Operations
  getPlayers(params?: PlayerQueryParams): Promise<ServiceResult<Player[]>>;
  getPlayerById(id: number): Promise<ServiceResult<Player>>;
  createPlayer(data: CreatePlayerDto): Promise<ServiceResult<Player>>;
  updatePlayer(id: number, data: UpdatePlayerDto): Promise<ServiceResult<Player>>;
  deletePlayer(id: number): Promise<ServiceResult<void>>;

  // Verification
  verifyPlayer(id: number, status: string, remark?: string): Promise<ServiceResult<Player>>;
  canVerify(player: Player, newStatus: string): VerificationCheckResult;

  // Skill Tags
  updateSkillTags(id: number, tags: string[]): Promise<ServiceResult<void>>;
  validateSkillTags(tags: string[]): SkillTagValidationResult;
  parseSkillTags(tagsString: string): string[];

  // Batch Operations
  batchUpdateStatus(playerIds: number[], status: string): Promise<BatchResult<void>>;
  batchDelete(playerIds: number[]): Promise<BatchResult<void>>;

  // Earnings
  calculateEarnings(order: Order, commissionRule?: CommissionRule): EarningsCalculation;
  computeStatistics(player: Player, orders: Order[]): PlayerStatistics;
}

/**
 * Valid verification statuses
 */
export const VERIFICATION_STATUSES = ['pending', 'verified', 'rejected'] as const;
export type VerificationStatus = typeof VERIFICATION_STATUSES[number];

/**
 * Valid verification state transitions
 * pending -> verified | rejected
 * verified -> rejected (can be revoked)
 * rejected -> pending (can reapply)
 */
const VALID_VERIFICATION_TRANSITIONS: Record<VerificationStatus, VerificationStatus[]> = {
  pending: ['verified', 'rejected'],
  verified: ['rejected'],
  rejected: ['pending'],
};

/**
 * Default commission rate (20%)
 */
const DEFAULT_COMMISSION_RATE = 0.20;

/**
 * Allowed skill tags (can be extended)
 */
const ALLOWED_SKILL_TAGS = [
  '上分',
  '陪玩',
  '教学',
  '娱乐',
  '语音',
  '高端局',
  '新手友好',
  '耐心',
  '技术流',
  '幽默',
  '女生',
  '男生',
  '深夜档',
  '白天档',
  '周末档',
];

/**
 * Player Service Implementation
 *
 * Provides all player-related business logic including:
 * - CRUD operations
 * - Verification workflow
 * - Skill tag management
 * - Earnings calculation
 * - Statistics computation
 * - Batch operations
 */
export class PlayerService extends BaseService implements IPlayerService {
  constructor(deps: ServiceDependencies = {}) {
    super(deps);
  }

  // ==================== CRUD Operations ====================

  /**
   * Get players with optional filtering
   */
  async getPlayers(params?: PlayerQueryParams): Promise<ServiceResult<Player[]>> {
    return this.withLogging('getPlayers', { params: params ?? {} }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.getPlayers(params);
        return response.data.data;
      }, 'Failed to fetch players');
    });
  }

  /**
   * Get a single player by ID
   */
  async getPlayerById(id: number): Promise<ServiceResult<Player>> {
    return this.withLogging('getPlayerById', { id }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.getPlayer(id);
        return response.data.data;
      }, `Failed to fetch player ${id}`);
    });
  }

  /**
   * Create a new player
   */
  async createPlayer(data: CreatePlayerDto): Promise<ServiceResult<Player>> {
    return this.withLogging('createPlayer', { userId: data.userId }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.createPlayer(data);
        return response.data.data;
      }, 'Failed to create player');
    });
  }

  /**
   * Update an existing player
   */
  async updatePlayer(id: number, data: UpdatePlayerDto): Promise<ServiceResult<Player>> {
    return this.withLogging('updatePlayer', { id }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.updatePlayer(id, data);
        return response.data.data;
      }, `Failed to update player ${id}`);
    });
  }

  /**
   * Delete a player
   */
  async deletePlayer(id: number): Promise<ServiceResult<void>> {
    return this.withLogging('deletePlayer', { id }, async () => {
      return this.wrapAsync(async () => {
        await this.api.deletePlayer(id);
      }, `Failed to delete player ${id}`);
    });
  }

  // ==================== Verification ====================

  /**
   * Verify a player (change verification status)
   */
  async verifyPlayer(
    id: number,
    status: string,
    remark?: string
  ): Promise<ServiceResult<Player>> {
    return this.withLogging('verifyPlayer', { id, status }, async () => {
      // First fetch the player to validate state transition
      const playerResult = await this.getPlayerById(id);
      if (!playerResult.success || !playerResult.data) {
        return ServiceResultHelper.failure({
          code: ServiceErrorCodes.PLAYER_NOT_FOUND,
          message: `Player ${id} not found`,
        });
      }

      // Check if verification transition is allowed
      const verificationCheck = this.canVerify(playerResult.data, status);
      if (!verificationCheck.allowed) {
        return ServiceResultHelper.failure({
          code: ServiceErrorCodes.PLAYER_INVALID_VERIFICATION,
          message: verificationCheck.reason || 'Invalid verification status transition',
          details: {
            playerId: id,
            currentStatus: playerResult.data.verificationStatus,
            requestedStatus: status,
          },
        });
      }

      return this.wrapAsync(async () => {
        const response = await this.api.updatePlayerVerification(id, status, remark);
        return response.data.data;
      }, `Failed to verify player ${id}`);
    });
  }

  /**
   * Check if a verification status transition is allowed
   * Enforces valid state transitions:
   * - pending -> verified | rejected
   * - verified -> rejected (can be revoked)
   * - rejected -> pending (can reapply)
   */
  canVerify(player: Player, newStatus: string): VerificationCheckResult {
    const currentStatus = player.verificationStatus as VerificationStatus;
    
    // Validate new status is a valid verification status
    if (!VERIFICATION_STATUSES.includes(newStatus as VerificationStatus)) {
      return {
        allowed: false,
        reason: `Invalid verification status: ${newStatus}. Valid statuses are: ${VERIFICATION_STATUSES.join(', ')}`,
      };
    }

    // Same status - no change needed
    if (currentStatus === newStatus) {
      return {
        allowed: false,
        reason: `Player is already in '${currentStatus}' status`,
      };
    }

    // Check if transition is valid
    const allowedTransitions = VALID_VERIFICATION_TRANSITIONS[currentStatus];
    if (!allowedTransitions.includes(newStatus as VerificationStatus)) {
      const transitionMessages: Record<string, string> = {
        'verified->pending': 'Verified players cannot be set back to pending. Use rejected status first.',
        'rejected->verified': 'Rejected players must reapply (set to pending) before being verified.',
      };

      const transitionKey = `${currentStatus}->${newStatus}`;
      return {
        allowed: false,
        reason: transitionMessages[transitionKey] || 
          `Cannot transition from '${currentStatus}' to '${newStatus}'`,
      };
    }

    return { allowed: true };
  }

  // ==================== Skill Tags ====================

  /**
   * Update player skill tags
   */
  async updateSkillTags(id: number, tags: string[]): Promise<ServiceResult<void>> {
    // Validate tags first
    const validation = this.validateSkillTags(tags);
    if (!validation.valid) {
      return ServiceResultHelper.failure({
        code: ServiceErrorCodes.VALIDATION_ERROR,
        message: `Invalid skill tags: ${validation.invalidTags.join(', ')}`,
        details: { invalidTags: validation.invalidTags },
      });
    }

    return this.withLogging('updateSkillTags', { id, tags }, async () => {
      return this.wrapAsync(async () => {
        await this.api.updatePlayerSkillTags(id, tags);
      }, `Failed to update skill tags for player ${id}`);
    });
  }

  /**
   * Validate skill tags against allowed values
   */
  validateSkillTags(tags: string[]): SkillTagValidationResult {
    const invalidTags = tags.filter(tag => !ALLOWED_SKILL_TAGS.includes(tag.trim()));
    return {
      valid: invalidTags.length === 0,
      invalidTags,
    };
  }

  /**
   * Parse comma-separated skill tags string
   * Trims whitespace from each tag
   */
  parseSkillTags(tagsString: string): string[] {
    if (!tagsString || typeof tagsString !== 'string') {
      return [];
    }
    
    return tagsString
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag.length > 0);
  }

  // ==================== Batch Operations ====================

  /**
   * Batch update player verification status
   */
  async batchUpdateStatus(
    playerIds: number[],
    status: string
  ): Promise<BatchResult<void>> {
    // Validate status
    if (!VERIFICATION_STATUSES.includes(status as VerificationStatus)) {
      return {
        success: false,
        total: playerIds.length,
        succeeded: 0,
        failed: playerIds.length,
        results: playerIds.map((_, index) => ({
          index,
          success: false,
          error: {
            code: ServiceErrorCodes.PLAYER_INVALID_VERIFICATION,
            message: `Invalid verification status: ${status}`,
          },
        })),
      };
    }

    if (playerIds.length === 0) {
      return ServiceResultHelper.emptyBatch(0);
    }

    return this.withLogging(
      'batchUpdateStatus',
      { playerIds, status },
      async () => {
        return this.executeBatch(
          playerIds,
          async (playerId) => {
            // Fetch player first to validate state transition
            const playerResult = await this.getPlayerById(playerId);
            if (!playerResult.success || !playerResult.data) {
              throw new Error(`Player ${playerId} not found`);
            }

            const verificationCheck = this.canVerify(playerResult.data, status);
            if (!verificationCheck.allowed) {
              throw new Error(verificationCheck.reason || 'Invalid verification transition');
            }

            await this.api.updatePlayerVerification(playerId, status);
          },
          'batchUpdateStatus'
        );
      }
    );
  }

  /**
   * Batch delete players
   */
  async batchDelete(playerIds: number[]): Promise<BatchResult<void>> {
    if (playerIds.length === 0) {
      return ServiceResultHelper.emptyBatch(0);
    }

    return this.withLogging('batchDelete', { playerIds }, async () => {
      return this.executeBatch(
        playerIds,
        async (playerId) => {
          await this.api.deletePlayer(playerId);
        },
        'batchDelete'
      );
    });
  }

  // ==================== Earnings ====================

  /**
   * Calculate player earnings from an order
   * Applies commission rules to determine net earnings
   */
  calculateEarnings(order: Order, commissionRule?: CommissionRule): EarningsCalculation {
    const grossAmount = order.totalPriceCents;
    
    // Determine commission rate
    let commissionRate = DEFAULT_COMMISSION_RATE;
    if (commissionRule && commissionRule.status === 'active') {
      commissionRate = commissionRule.ratePercent / 100;
    }

    // Ensure commission rate is within valid range (0-1)
    commissionRate = Math.max(0, Math.min(1, commissionRate));

    // Calculate commission and net amounts
    const commissionAmount = Math.round(grossAmount * commissionRate);
    const netAmount = grossAmount - commissionAmount;

    return {
      grossAmount,
      commissionRate,
      commissionAmount,
      netAmount,
    };
  }

  /**
   * Compute player statistics from orders
   */
  computeStatistics(player: Player, orders: Order[]): PlayerStatistics {
    // Filter completed orders for this player
    const completedOrders = orders.filter(
      order => order.playerId === player.id && order.status === 'completed'
    );

    // Calculate total earnings (sum of all completed order amounts minus commission)
    let totalEarnings = 0;
    let monthlyEarnings = 0;
    const currentMonth = new Date().getMonth();
    const currentYear = new Date().getFullYear();

    for (const order of completedOrders) {
      const earnings = this.calculateEarnings(order);
      totalEarnings += earnings.netAmount;

      // Check if order is from current month
      const orderDate = new Date(order.completedAt || order.createdAt);
      if (orderDate.getMonth() === currentMonth && orderDate.getFullYear() === currentYear) {
        monthlyEarnings += earnings.netAmount;
      }
    }

    return {
      totalEarnings,
      monthlyEarnings,
      completedOrders: completedOrders.length,
      averageRating: player.ratingAverage,
      ratingCount: player.ratingCount,
    };
  }
}

/**
 * Default PlayerService instance
 */
export const playerService = new PlayerService();
