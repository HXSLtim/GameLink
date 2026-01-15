import { render, screen, waitFor } from '@testing-library/react';
import { vi, describe, beforeEach, test, expect } from 'vitest';
import FavoritesPage from './FavoritesPage';
import { BrowserRouter } from 'react-router-dom';
import { http } from '@/lib/http';

// Mock dependencies
vi.mock('@/lib/http');
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => vi.fn(),
    };
});
vi.mock('react-i18next', () => ({
    useTranslation: () => ({
        t: (key: string, options?: any) => options?.defaultValue || key,
    }),
}));

describe('FavoritesPage', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    test('renders loading state initially', () => {
        (http.get as any).mockReturnValue(new Promise(() => { })); // pending promise
        render(
            <BrowserRouter>
                <FavoritesPage />
            </BrowserRouter>
        );
        expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
    });

    test('renders favorite players', async () => {
        const mockPlayers = [
            { id: 1, nickname: 'ProPlayer1', gameName: 'LOL', price: 50, avatar: 'avatar1.png', tags: ['pro'], rating: 5 },
            { id: 2, nickname: 'ProPlayer2', gameName: 'DOTA2', price: 60, avatar: 'avatar2.png', tags: ['fun'], rating: 4.5 }
        ];

        (http.get as any).mockResolvedValue(mockPlayers);

        render(
            <BrowserRouter>
                <FavoritesPage />
            </BrowserRouter>
        );

        await waitFor(() => {
            expect(screen.getByText('ProPlayer1')).toBeInTheDocument();
            expect(screen.getByText('ProPlayer2')).toBeInTheDocument();
        });
    });

    test('shows empty state when no favorites', async () => {
        (http.get as any).mockResolvedValue([]);

        render(
            <BrowserRouter>
                <FavoritesPage />
            </BrowserRouter>
        );

        await waitFor(() => {
            expect(screen.getByTestId('empty-state')).toBeInTheDocument();
            expect(screen.getByText('No favorites yet')).toBeInTheDocument();
        });
    });
});
