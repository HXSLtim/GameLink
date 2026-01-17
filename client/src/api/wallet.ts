/**
 * Wallet API
 * Handles balance, transactions, withdrawals
 */

import { http } from '@/lib/http';
import type {
    Wallet,
    Transaction,
    WithdrawRequest,
    TransactionListParams,
    PaginatedResponse
} from '@/types/api';

export const walletApi = {
    /**
     * Get wallet info
     */
    get: () =>
        http.get<Wallet>('/wallet'),

    /**
     * Get transaction list
     */
    getTransactions: (params: TransactionListParams) =>
        http.get<PaginatedResponse<Transaction>>('/wallet/transactions', { params }),

    /**
     * Get transaction detail
     */
    getTransaction: (id: number) =>
        http.get<Transaction>(`/wallet/transaction/${id}`),

    /**
     * Request withdrawal
     */
    withdraw: (data: WithdrawRequest) =>
        http.post<{ withdrawId: number }>('/wallet/withdraw', data),

    /**
     * Get withdrawal list
     */
    getWithdrawals: (params: { page: number; pageSize: number; status?: string }) =>
        http.get<PaginatedResponse<any>>('/wallet/withdrawals', { params }),

    /**
     * Cancel withdrawal
     */
    cancelWithdrawal: (id: number) =>
        http.post<void>(`/wallet/withdrawal/${id}/cancel`),

    /**
     * Get wallet statistics
     */
    getStats: () =>
        http.get<{
            totalIncome: number;
            totalWithdrawn: number;
            pendingWithdrawal: number;
            frozenAmount: number;
        }>('/wallet/stats'),

    /**
     * Get withdrawal methods
     */
    getWithdrawMethods: () =>
        http.get<Array<{
            id: string;
            name: string;
            minAmount: number;
            maxAmount: number;
            fee: number;
        }>>('/wallet/withdraw/methods'),
};
