import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { http } from '@/lib/http';

export interface Transaction {
    id: string; // or number, backend sends number but we can normalize
    type: 'recharge' | 'payment' | 'refund' | 'withdrawal';
    amount: number; // Stored as units (e.g., CNY)
    status: 'pending' | 'success' | 'failed';
    createdAt: string;
    description: string;
}

export interface WalletState {
    balance: number; // Stored as units (e.g., 12.50)
    frozenBalance: number;
    currency: string;
    transactions: Transaction[];
    loading: boolean;
    error: string | null;
}

export interface WalletActions {
    fetchWallet: () => Promise<void>;
    fetchTransactions: (page?: number) => Promise<void>;
    recharge: (amount: number, method: string) => Promise<void>;
    withdraw: (amount: number, bankCardId: number) => Promise<void>;
}

export const useWalletStore = create<WalletState & WalletActions>()(
    persist(
        (set, get) => ({
            balance: 0,
            frozenBalance: 0,
            currency: 'CNY',
            transactions: [],
            loading: false,
            error: null,

            fetchWallet: async () => {
                set({ loading: true, error: null });
                try {
                    const data = await http.get<any>('/user/wallet/balance');
                    // Backend returns cents
                    const balance = (data.balanceCents || 0) / 100;
                    const frozen = (data.frozenCents || 0) / 100;

                    set({ balance: balance, frozenBalance: frozen, loading: false });
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Failed to fetch wallet' });
                }
            },

            fetchTransactions: async (page = 1) => {
                set({ loading: true, error: null });
                try {
                    const data = await http.get<any>('/user/wallet/transactions', {
                        params: { page, pageSize: 20 }
                    });

                    const items = (data.items || []).map((t: any) => ({
                        id: String(t.id),
                        type: t.type,
                        amount: t.amountCents / 100, // Convert cents to units
                        status: t.status,
                        createdAt: t.createdAt,
                        description: t.description
                    }));

                    set({ transactions: items, loading: false });
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Failed to fetch transactions' });
                }
            },

            recharge: async (amount, method) => {
                set({ loading: true, error: null });
                try {
                    // frontend amount is units, backend expects cents
                    const amountCents = Math.round(amount * 100);
                    const data = await http.post<any>('/user/wallet/recharge', {
                        amountCents,
                        method
                    });

                    // Update balance immediately if returned, or fetch
                    if (data.balanceCents !== undefined) {
                        set({ balance: data.balanceCents / 100, loading: false });
                    } else {
                        // fallback
                        await get().fetchWallet();
                        set({ loading: false });
                    }
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Recharge failed' });
                    throw err;
                }
            },

            withdraw: async (amount, bankCardId) => {
                set({ loading: true, error: null });
                try {
                    // frontend amount is units, backend expects cents
                    const amountCents = Math.round(amount * 100);
                    await http.post('/user/wallet/withdraw', {
                        amountCents,
                        bankCardId
                    });

                    // Withdrawal usually pends, balance might be frozen or deducted
                    await get().fetchWallet();
                    set({ loading: false });
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Withdrawal failed' });
                    throw err;
                }
            }
        }),
        {
            name: 'wallet-storage',
            partialize: (state) => ({
                balance: state.balance,
                currency: state.currency
            })
        }
    )
);
