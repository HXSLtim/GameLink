import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import LFGPage from '../LFGPage';

// Mock stores
const mockFetchPendingRequests = vi.fn();
const mockFetchMyRequests = vi.fn();
const mockFetchActiveRequest = vi.fn();
const mockAcceptRequest = vi.fn();
const mockCancelRequest = vi.fn();

vi.mock('@/stores', () => ({
    useLFGStore: vi.fn(() => ({
        requests: [],
        myRequests: [],
        activeRequest: null,
        isLoading: false,
        fetchPendingRequests: mockFetchPendingRequests,
        fetchMyRequests: mockFetchMyRequests,
        fetchActiveRequest: mockFetchActiveRequest,
        acceptRequest: mockAcceptRequest,
        cancelRequest: mockCancelRequest,
    })),
    useAuthStore: vi.fn(() => ({
        user: { id: 1, nickname: 'TestUser' },
    })),
}));

// Mock react-i18next
vi.mock('react-i18next', () => ({
    useTranslation: () => ({
        t: (key: string) => {
            const translations: Record<string, string> = {
                'lfg.title': 'Looking For Group',
                'lfg.description': 'Find players or teams',
                'lfg.searchPlaceholder': 'Search requests...',
                'lfg.filterType': 'Filter by type',
                'lfg.type.all': 'All',
                'lfg.type.findPlayer': 'Find Player',
                'lfg.type.findTeam': 'Find Team',
                'lfg.create': 'Create Request',
                'lfg.empty': 'No requests found',
                'lfg.emptyDescription': 'Be the first to create a request',
                'lfg.createFirst': 'Create First Request',
                'lfg.myEmpty': 'No requests yet',
                'lfg.myEmptyDescription': 'Create a request to find players',
                'lfg.tabs.browse': 'Browse',
                'lfg.tabs.my': 'My Requests',
                'lfg.activeRequest': 'Active Request',
                'lfg.untitled': 'Untitled',
                'lfg.cancel': 'Cancel',
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
vi.mock('@/components/lfg', () => ({
    LFGRequestCard: ({ request, onAccept, onCancel }: { request: { id: number; title: string }; onAccept?: (id: number) => void; onCancel?: (id: number) => void }) => (
        <div data-testid={`lfg-card-${request.id}`}>
            <span>{request.title}</span>
            {onAccept && <button onClick={() => onAccept(request.id)}>Accept</button>}
            {onCancel && <button onClick={() => onCancel(request.id)}>Cancel</button>}
        </div>
    ),
    LFGRequestCardSkeleton: () => <div data-testid="lfg-skeleton">Loading...</div>,
}));

vi.mock('@/components/page-container', () => ({
    PageContainer: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
    PageHeader: ({ title, description }: { title: string; description: string }) => (
        <div>
            <h1>{title}</h1>
            <p>{description}</p>
        </div>
    ),
}));

import { useLFGStore } from '@/stores';

const renderWithRouter = (component: React.ReactElement) => {
    return render(<BrowserRouter>{component}</BrowserRouter>);
};

describe('LFGPage', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should render page header', () => {
        renderWithRouter(<LFGPage />);

        expect(screen.getByText('Looking For Group')).toBeInTheDocument();
        expect(screen.getByText('Find players or teams')).toBeInTheDocument();
    });

    it('should render search input', () => {
        renderWithRouter(<LFGPage />);

        expect(screen.getByPlaceholderText('Search requests...')).toBeInTheDocument();
    });

    it('should render create request button', () => {
        renderWithRouter(<LFGPage />);

        expect(screen.getByText('Create Request')).toBeInTheDocument();
    });

    it('should fetch data on mount', () => {
        renderWithRouter(<LFGPage />);

        expect(mockFetchPendingRequests).toHaveBeenCalled();
        expect(mockFetchMyRequests).toHaveBeenCalled();
        expect(mockFetchActiveRequest).toHaveBeenCalled();
    });

    it('should show loading skeletons when loading', () => {
        vi.mocked(useLFGStore).mockReturnValue({
            requests: [],
            myRequests: [],
            activeRequest: null,
            isLoading: true,
            fetchPendingRequests: mockFetchPendingRequests,
            fetchMyRequests: mockFetchMyRequests,
            fetchActiveRequest: mockFetchActiveRequest,
            acceptRequest: mockAcceptRequest,
            cancelRequest: mockCancelRequest,
        } as Partial<ReturnType<typeof useLFGStore>>);

        renderWithRouter(<LFGPage />);

        expect(screen.getAllByTestId('lfg-skeleton')).toHaveLength(3);
    });

    it('should show empty state when no requests', () => {
        vi.mocked(useLFGStore).mockReturnValue({
            requests: [],
            myRequests: [],
            activeRequest: null,
            isLoading: false,
            fetchPendingRequests: mockFetchPendingRequests,
            fetchMyRequests: mockFetchMyRequests,
            fetchActiveRequest: mockFetchActiveRequest,
            acceptRequest: mockAcceptRequest,
            cancelRequest: mockCancelRequest,
        } as Partial<ReturnType<typeof useLFGStore>>);

        renderWithRouter(<LFGPage />);

        expect(screen.getByText('No requests found')).toBeInTheDocument();
    });

    it('should render request cards when requests exist', () => {
        const mockRequests = [
            { id: 1, title: 'Request 1', requestType: 'find_player' },
            { id: 2, title: 'Request 2', requestType: 'find_team' },
        ];

        vi.mocked(useLFGStore).mockReturnValue({
            requests: mockRequests,
            myRequests: [],
            activeRequest: null,
            isLoading: false,
            fetchPendingRequests: mockFetchPendingRequests,
            fetchMyRequests: mockFetchMyRequests,
            fetchActiveRequest: mockFetchActiveRequest,
            acceptRequest: mockAcceptRequest,
            cancelRequest: mockCancelRequest,
        } as Partial<ReturnType<typeof useLFGStore>>);

        renderWithRouter(<LFGPage />);

        expect(screen.getByTestId('lfg-card-1')).toBeInTheDocument();
        expect(screen.getByTestId('lfg-card-2')).toBeInTheDocument();
    });

    it('should navigate to create page when create button clicked', () => {
        renderWithRouter(<LFGPage />);

        fireEvent.click(screen.getByText('Create Request'));

        expect(mockNavigate).toHaveBeenCalledWith('/lfg/create');
    });

    it('should show active request banner when active request exists', () => {
        vi.mocked(useLFGStore).mockReturnValue({
            requests: [],
            myRequests: [],
            activeRequest: { id: 1, title: 'My Active Request' },
            isLoading: false,
            fetchPendingRequests: mockFetchPendingRequests,
            fetchMyRequests: mockFetchMyRequests,
            fetchActiveRequest: mockFetchActiveRequest,
            acceptRequest: mockAcceptRequest,
            cancelRequest: mockCancelRequest,
        } as Partial<ReturnType<typeof useLFGStore>>);

        renderWithRouter(<LFGPage />);

        expect(screen.getByText('Active Request')).toBeInTheDocument();
        expect(screen.getByText('My Active Request')).toBeInTheDocument();
    });

    it('should disable create button when active request exists', () => {
        vi.mocked(useLFGStore).mockReturnValue({
            requests: [],
            myRequests: [],
            activeRequest: { id: 1, title: 'Active' },
            isLoading: false,
            fetchPendingRequests: mockFetchPendingRequests,
            fetchMyRequests: mockFetchMyRequests,
            fetchActiveRequest: mockFetchActiveRequest,
            acceptRequest: mockAcceptRequest,
            cancelRequest: mockCancelRequest,
        } as Partial<ReturnType<typeof useLFGStore>>);

        renderWithRouter(<LFGPage />);

        const createButton = screen.getByText('Create Request').closest('button');
        expect(createButton).toBeDisabled();
    });

    it('should filter requests by search query', () => {
        const mockRequests = [
            { id: 1, title: 'League Request', requestType: 'find_player' },
            { id: 2, title: 'Valorant Request', requestType: 'find_team' },
        ];

        vi.mocked(useLFGStore).mockReturnValue({
            requests: mockRequests,
            myRequests: [],
            activeRequest: null,
            isLoading: false,
            fetchPendingRequests: mockFetchPendingRequests,
            fetchMyRequests: mockFetchMyRequests,
            fetchActiveRequest: mockFetchActiveRequest,
            acceptRequest: mockAcceptRequest,
            cancelRequest: mockCancelRequest,
        } as Partial<ReturnType<typeof useLFGStore>>);

        renderWithRouter(<LFGPage />);

        const searchInput = screen.getByPlaceholderText('Search requests...');
        fireEvent.change(searchInput, { target: { value: 'League' } });

        expect(screen.getByTestId('lfg-card-1')).toBeInTheDocument();
        expect(screen.queryByTestId('lfg-card-2')).not.toBeInTheDocument();
    });

    it('should call acceptRequest when accept button clicked', async () => {
        const mockRequests = [{ id: 1, title: 'Test Request', requestType: 'find_player' }];
        mockAcceptRequest.mockResolvedValue({ id: 10 });

        vi.mocked(useLFGStore).mockReturnValue({
            requests: mockRequests,
            myRequests: [],
            activeRequest: null,
            isLoading: false,
            fetchPendingRequests: mockFetchPendingRequests,
            fetchMyRequests: mockFetchMyRequests,
            fetchActiveRequest: mockFetchActiveRequest,
            acceptRequest: mockAcceptRequest,
            cancelRequest: mockCancelRequest,
        } as Partial<ReturnType<typeof useLFGStore>>);

        renderWithRouter(<LFGPage />);

        const acceptButton = screen.getByText('Accept');
        fireEvent.click(acceptButton);

        await waitFor(() => {
            expect(mockAcceptRequest).toHaveBeenCalledWith(1);
        });
    });

    it('should navigate to room after accepting request', async () => {
        const mockRequests = [{ id: 1, title: 'Test Request', requestType: 'find_player' }];
        mockAcceptRequest.mockResolvedValue({ id: 10 });

        vi.mocked(useLFGStore).mockReturnValue({
            requests: mockRequests,
            myRequests: [],
            activeRequest: null,
            isLoading: false,
            fetchPendingRequests: mockFetchPendingRequests,
            fetchMyRequests: mockFetchMyRequests,
            fetchActiveRequest: mockFetchActiveRequest,
            acceptRequest: mockAcceptRequest,
            cancelRequest: mockCancelRequest,
        } as Partial<ReturnType<typeof useLFGStore>>);

        renderWithRouter(<LFGPage />);

        fireEvent.click(screen.getByText('Accept'));

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith('/rooms/10');
        });
    });

    it('should call cancelRequest when cancel button clicked on active request', async () => {
        vi.mocked(useLFGStore).mockReturnValue({
            requests: [],
            myRequests: [],
            activeRequest: { id: 1, title: 'Active Request' },
            isLoading: false,
            fetchPendingRequests: mockFetchPendingRequests,
            fetchMyRequests: mockFetchMyRequests,
            fetchActiveRequest: mockFetchActiveRequest,
            acceptRequest: mockAcceptRequest,
            cancelRequest: mockCancelRequest,
        } as Partial<ReturnType<typeof useLFGStore>>);

        renderWithRouter(<LFGPage />);

        // Click cancel on the active request banner
        const cancelButton = screen.getByText('Cancel');
        fireEvent.click(cancelButton);

        await waitFor(() => {
            expect(mockCancelRequest).toHaveBeenCalledWith(1);
        });
    });
});
