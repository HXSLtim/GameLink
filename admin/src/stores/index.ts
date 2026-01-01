// Admin Panel Stores
export { useAuthStore, type AuthState } from './modules/authStore';
export { useUserStore } from './modules/userStore';
export { useMenuStore } from './modules/menuStore';
export { useOrderStore } from './modules/orderStore';
export { usePlayerStore } from './modules/playerStore';
export { useChatStore } from './modules/chatStore';

// Import for internal use in useAppStores
import { useAuthStore as _useAuthStore } from './modules/authStore';
import { useUserStore as _useUserStore } from './modules/userStore';
import { useMenuStore as _useMenuStore } from './modules/menuStore';
import { useOrderStore as _useOrderStore } from './modules/orderStore';
import { usePlayerStore as _usePlayerStore } from './modules/playerStore';
import { useChatStore as _useChatStore } from './modules/chatStore';

// Types
export * from './types';

// Combined hook - use direct imports since stores are already loaded
export const useAppStores = () => ({
  auth: _useAuthStore(),
  user: _useUserStore(),
  menu: _useMenuStore(),
  order: _useOrderStore(),
  player: _usePlayerStore(),
  chat: _useChatStore(),
});
