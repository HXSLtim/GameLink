import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import RoomListPage from '../RoomListPage';

// Mock stores
const mockFetchRooms = vi.fn();
const mockJoinRoom = vi.fn();

vi.mock('@/stores', () => ({
    useRoomStore: vi.fn(() => ({
        rooms: [],
        isLoading: false,
        pagination: { page: 1, pageSize: 20, total: 0 },
        fetchRooms: mockFetchRooms,
        joinRoom: mockJoinRoom,
    })),
    useAuthStore: vi.fn(() => ({
        user: { id: 1, nickname: 'TestUser' },
    })),
}));

// Mock react-i18next
vi.mock('react-i18next', () => ({
    useTranslation: () => ({
        t: (key: string, options?: any) => {
            const translations: Record<string, string> = {
                'room.title': 'Game Rooms',
                'room.description': 'Find or create a room',
                'room.searchPlaceholder': 'Search rooms...',
                'room.filterStatus': 'Filter by status',
                'room.status.all': 'All',
                'room.status.waiting': 'Waiting',
                'room.status.ready': 'Ready',
                'room.status.in_game': 'In Game',
                'room.status.finished': 'Finished',
                'room.create': 'Create Room',
                'room.empty': 'No rooms found',
                'room.emptyDescription': 'Be the first to create a room',
                'room.createFirst': 'Create First Room',
                'room.showing': `Showing ${options?.count || 0} of ${options?.total || 0}`,
            };
            return translations[key] || key;
        },
    }),
}));

// Mock navigate
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => mockNavigate,
    };
});

// Mock components
vi.mock('@/components/room', () => ({
    RoomCard: ({ room, onJoin }: any) => (
        <div data-testid={`room-card-${room.id}`}>
            <span>{room.name}</span>
            <button onClick={() => onJoin(room.id)}>Join</button>
        </div>
    ),
    RoomCardSkeleton: () => <div data-testid="room-skeleton">Loading...</div>,
}));

vi.mock('@/components/page-container', () => ({
    PageContainer: ({ children }: any) => <div>{children}</div>,
    PageHeader: ({ title, description }: any) => (
        <div>
            <h1>{title}</h1>
            <p>{description}</p>
        </div>
    ),
}));

import { useRoomStore } from '@/stores';

const renderWithRouter = (component: React.ReactElement) => {
    return render(<BrowserRouter>{component}</BrowserRouter>);
};

describe('RoomListPage', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should render page header', () => {
        renderWithRouter(<RoomListPage />);

        expect(screen.getByText('Game Rooms')).toBeInTheDocument();
        expect(screen.getByText('Find or create a room')).toBeInTheDocument();
    });

    it('should render search input', () => {
        renderWithRouter(<RoomListPage />);

        expect(screen.getByPlaceholderText('Search rooms...')).toBeInTheDocument();
    });

    it('should render create room button', () => {
        renderWithRouter(<RoomListPage />);

        expect(screen.getByText('Create Room')).toBeInTheDocument();
    });

    it('should fetch rooms on mount', () => {
        renderWithRouter(<RoomListPage />);

        expect(mockFetchRooms).toHaveBeenCalledWith({
            page: 1,
            pageSize: 20,
            status: undefined,
        });
    });

    it('should show loading skeletons when loading', () => {
        vi.mocked(useRoomStore).mockReturnValue({
            rooms: [],
            isLoading: true,
            pagination: { page: 1, pageSize: 20, total: 0 },
            fetchRooms: mockFetchRooms,
            joinRoom: mockJoinRoom,
        } as any);

        renderWithRouter(<RoomListPage />);

        expect(screen.getAllByTestId('room-skeleton')).toHaveLength(6);
    });

    it('should show empty state when no rooms', () => {
        vi.mocked(useRoomStore).mockReturnValue({
            rooms: [],
            isLoading: false,
            pagination: { page: 1, pageSize: 20, total: 0 },
            fetchRooms: mockFetchRooms,
            joinRoom: mockJoinRoom,
        } as any);

        renderWithRouter(<RoomListPage />);

        expect(screen.getByText('No rooms found')).toBeInTheDocument();
        expect(screen.getByText('Be the first to create a room')).toBeInTheDocument();
    });

    it('should render room cards when rooms exist', () => {
        const mockRooms = [
            { id: 1, name: 'Room 1', gameName: 'Game 1' },
            { id: 2, name: 'Room 2', gameName: 'Game 2' },
        ];

        vi.mocked(useRoomStore).mockReturnValue({
            rooms: mockRooms,
            isLoading: false,
            pagination: { page: 1, pageSize: 20, total: 2 },
            fetchRooms: mockFetchRooms,
            joinRoom: mockJoinRoom,
        } as any);

        renderWithRouter(<RoomListPage />);

        expect(screen.getByTestId('room-card-1')).toBeInTheDocument();
        expect(screen.getByTestId('room-card-2')).toBeInTheDocument();
    });

    it('should navigate to create room page when create button clicked', () => {
        renderWithRouter(<RoomListPage />);

        fireEvent.click(screen.getByText('Create Room'));

        expect(mockNavigate).toHaveBeenCalledWith('/rooms/create');
    });

    it('should filter rooms by search query', () => {
        const mockRooms = [
            { id: 1, name: 'League Room', gameName: 'League of Legends' },
            { id: 2, name: 'Valorant Room', gameName: 'Valorant' },
        ];

        vi.mocked(useRoomStore).mockReturnValue({
            rooms: mockRooms,
            isLoading: false,
            pagination: { page: 1, pageSize: 20, total: 2 },
            fetchRooms: mockFetchRooms,
            joinRoom: mockJoinRoom,
        } as any);

        renderWithRouter(<RoomListPage />);

        const searchInput = screen.getByPlaceholderText('Search rooms...');
        fireEvent.change(searchInput, { target: { value: 'League' } });

        expect(screen.getByTestId('room-card-1')).toBeInTheDocument();
        expect(screen.queryByTestId('room-card-2')).not.toBeInTheDocument();
    });

    it('should call joinRoom when join button clicked', async () => {
        const mockRooms = [{ id: 1, name: 'Test Room', gameName: 'Test Game' }];

        vi.mocked(useRoomStore).mockReturnValue({
            rooms: mockRooms,
            isLoading: false,
            pagination: { page: 1, pageSize: 20, total: 1 },
            fetchRooms: mockFetchRooms,
            joinRoom: mockJoinRoom,
        } as any);

        renderWithRouter(<RoomListPage />);

        const joinButton = screen.getByText('Join');
        fireEvent.click(joinButton);

        await waitFor(() => {
            expect(mockJoinRoom).toHaveBeenCalledWith(1, undefined);
        });
    });
});
