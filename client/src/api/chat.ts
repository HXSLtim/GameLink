/**
 * Chat API
 * Handles messaging, chat rooms, message history
 */

import { http } from '@/lib/http';
import type {
    ChatMessage,
    ChatRoom,
    SendMessageRequest,
    PaginatedResponse
} from '@/types/api';

export const chatApi = {
    /**
     * Get chat room list
     */
    getRooms: (params: { page: number; pageSize: number }) =>
        http.get<PaginatedResponse<ChatRoom>>('/chat/rooms', { params }),

    /**
     * Get or create chat room with user
     */
    getOrCreateRoom: (userId: number) =>
        http.post<ChatRoom>('/chat/room', { userId }),

    /**
     * Get chat room detail
     */
    getRoom: (roomId: number) =>
        http.get<ChatRoom>(`/chat/room/${roomId}`),

    /**
     * Get message history
     */
    getMessages: (roomId: number, params: {
        page: number;
        pageSize: number;
        beforeId?: number;
    }) =>
        http.get<PaginatedResponse<ChatMessage>>(`/chat/room/${roomId}/messages`, { params }),

    /**
     * Send message
     */
    sendMessage: (roomId: number, data: SendMessageRequest) =>
        http.post<ChatMessage>(`/chat/room/${roomId}/message`, data),

    /**
     * Mark messages as read
     */
    markAsRead: (roomId: number, messageIds: number[]) =>
        http.post<void>(`/chat/room/${roomId}/read`, { messageIds }),

    /**
     * Delete message
     */
    deleteMessage: (roomId: number, messageId: number) =>
        http.delete<void>(`/chat/room/${roomId}/message/${messageId}`),

    /**
     * Get unread count
     */
    getUnreadCount: () =>
        http.get<{ total: number; rooms: Record<number, number> }>('/chat/unread'),

    /**
     * Upload chat image
     */
    uploadImage: (file: File) => {
        const formData = new FormData();
        formData.append('image', file);
        return http.post<{ url: string }>('/chat/upload/image', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
    },
};
