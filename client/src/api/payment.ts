/**
 * Payment API
 * Handles payment creation, status check, refunds
 */

import { http } from '@/lib/http';
import type {
    Payment,
    CreatePaymentRequest,
    PaymentListParams,
    PaginatedResponse
} from '@/types/api';

export const paymentApi = {
    /**
     * Create payment
     */
    create: (data: CreatePaymentRequest) =>
        http.post<Payment>('/payment/create', data),

    /**
     * Get payment detail
     */
    get: (id: number) =>
        http.get<Payment>(`/payment/${id}`),

    /**
     * Get payment list
     */
    list: (params: PaymentListParams) =>
        http.get<PaginatedResponse<Payment>>('/payment/list', { params }),

    /**
     * Check payment status
     */
    checkStatus: (id: number) =>
        http.get<{ status: string; paidAt?: string }>(`/payment/${id}/status`),

    /**
     * Request refund
     */
    refund: (id: number, reason: string) =>
        http.post<void>(`/payment/${id}/refund`, { reason }),

    /**
     * Get payment methods
     */
    getMethods: () =>
        http.get<Array<{ id: string; name: string; icon: string }>>('/payment/methods'),
};
