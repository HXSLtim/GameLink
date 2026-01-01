// Admin Panel Stores
export { useAuthStore } from './modules/authStore';
export { useUserStore } from './modules/userStore';
export { useMenuStore } from './modules/menuStore';
export { useOrderStore } from './modules/orderStore';
export { usePlayerStore } from './modules/playerStore';
export { useChatStore } from './modules/chatStore';

// Types
export * from './types';

// Combined hook
export const useAppStores = () => ({
  auth: useAuthStore(),
  user: useUserStore(),
  menu: useMenuStore(),
  order: useOrderStore(),
  player: usePlayerStore(),
  chat: useChatStore(),
});
