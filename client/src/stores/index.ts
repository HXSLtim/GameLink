export * from './modules/auth-store';
export * from './modules/theme-store';
export * from './modules/player-store';
export * from './modules/order-store';
export * from './modules/chat-store';
export * from './modules/wallet-store';
export * from './modules/vip-store';
export * from './modules/favorite-store';
export * from './modules/notification-store';
export * from './modules/presence-store';
export * from './modules/dispute-store';

// Game Room & Voice
// Note: room-store also exports ChatGroupType and ChatGroupStatus, but we use the ones from chat-store
// which include both the const object and type. Re-export room-store types with aliases to avoid conflicts.
export {
    useRoomStore,
    type GameRoom,
    type RoomMember,
    type CreateRoomRequest,
    type UpdateRoomRequest,
    type RoomListOptions,
    type RoomState,
    selectRooms,
    selectCurrentRoom,
    selectMembers,
    selectIsLoading,
    selectRoomError,
    // Re-export with aliases for room-specific usage if needed
    type ChatGroupType as RoomGroupType,
    type ChatGroupStatus as RoomGroupStatus,
} from './modules/room-store';
export * from './modules/lfg-store';
export * from './modules/voice-store';

// Marketing & Promotion
export * from './modules/coupon-store';
export * from './modules/activity-store';
export * from './modules/referral-store';

// Additional Features
export * from './modules/block-store';
export * from './modules/team-store';
export * from './modules/gift-store';
export * from './modules/review-store';
