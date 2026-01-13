import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface Transaction {
    id: string;
    // ... (rest is same, but I can't leave ... in replace_file_content unless I target specific lines)
    // Let's target specific blocks.
    id: string;
    type: 'recharge' | 'payment' | 'refund' | 'withdrawal';
    amount: number;
    status: 'pending' | 'success' | 'failed';
    createdAt: string;
    description: string;
}

export interface WalletState {
    balance: number;
    currency: string;
    transactions: Transaction[];
    loading: boolean;
    error: string | null;
}

export interface WalletActions {
    fetchWallet: () => Promise<void>;
    fetchTransactions: () => Promise<void>;
    recharge: (amount: number, method: string) => Promise<void>;
    withdraw: (amount: number, method: string) => Promise<void>;
}

export const useWalletStore = create<WalletState & WalletActions>()(
    persist(
        (set) => ({
            balance: 0,
            currency: 'CNY',
            transactions: [],
            loading: false,
            error: null,

            fetchWallet: async () => {
                set({ loading: true, error: null });
                try {
                    // Mock API call
                    // const data = await http.get('/wallet/balance');
                    await new Promise(resolve => setTimeout(resolve, 500));

                    const mockBalance = 1250.50; // Mock data

                    set({ balance: mockBalance, loading: false });
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Failed to fetch wallet' });
                }
            },

            fetchTransactions: async () => {
                set({ loading: true, error: null });
                try {
                    // Mock API call
                    // const data = await http.get('/wallet/transactions');
                    await new Promise(resolve => setTimeout(resolve, 500));

                    const mockTransactions: Transaction[] = [
                        { id: 't1', type: 'recharge', amount: 500, status: 'success', createdAt: new Date().toISOString(), description: 'Alipay Recharge' },
                        { id: 't2', type: 'payment', amount: -100, status: 'success', createdAt: new Date(Date.now() - 86400000).toISOString(), description: 'Order Payment #1234' },
                    ];

                    set({ transactions: mockTransactions, loading: false });
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Failed to fetch transactions' });
                }
            },

            recharge: async (amount, method) => {
                set({ loading: true, error: null });
                try {
                    // await http.post('/wallet/recharge', { amount, method });
                    await new Promise(resolve => setTimeout(resolve, 1000));

                    set(state => ({
                        balance: state.balance + amount,
                        loading: false,
                        transactions: [
                            {
                                id: `t_${Date.now()}`,
                                type: 'recharge',
                                amount,
                                status: 'success',
                                createdAt: new Date().toISOString(),
                                description: `${method} Recharge`
                            },
                            ...state.transactions
                        ]
                    }));
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Recharge failed' });
                }
            },

            withdraw: async (amount, method) => {
                set({ loading: true, error: null });
                try {
                    // await http.post('/wallet/withdraw', { amount, method });
                    await new Promise(resolve => setTimeout(resolve, 1000));

                    set(state => ({
                        balance: state.balance - amount,
                        loading: false,
                        transactions: [
                            {
                                id: `t_${Date.now()}`,
                                type: 'withdrawal',
                                amount: -amount,
                                status: 'pending',
                                createdAt: new Date().toISOString(),
                                description: `${method} Withdrawal`
                            },
                            ...state.transactions
                        ]
                    }));
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Withdrawal failed' });
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
