import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import LoginPage from '../LoginPage';
import { BrowserRouter } from 'react-router-dom';
import { useAuthStore } from '@/stores';

// Mock Modules
vi.mock('react-i18next', () => ({
    useTranslation: () => ({ t: (key: string) => key }),
}));

vi.mock('@/stores', () => ({
    useAuthStore: vi.fn(),
}));

vi.mock('@/components/language-switcher', () => ({
    LanguageSwitcher: () => <div data-testid="lang-switcher" />,
}));

vi.mock('@/components/mode-toggle', () => ({
    ModeToggle: () => <div data-testid="mode-toggle" />,
}));

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => mockNavigate,
    };
});

describe('LoginPage', () => {
    const mockLogin = vi.fn();
    const mockRegister = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useAuthStore as any).mockReturnValue({
            login: mockLogin,
            register: mockRegister,
            loading: false,
            error: null,
        });
    });

    it('renders login form by default', () => {
        render(
            <BrowserRouter>
                <LoginPage />
            </BrowserRouter>
        );
        expect(screen.getByPlaceholderText('auth.username')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('••••••••')).toBeInTheDocument();
        // The button text might be composed, checking key presence
        expect(screen.getByRole('button', { name: 'auth.sign_in' })).toBeInTheDocument();
    });

    it('validates input and calls login on submit', async () => {
        render(
            <BrowserRouter>
                <LoginPage />
            </BrowserRouter>
        );

        fireEvent.change(screen.getByPlaceholderText('auth.username'), { target: { value: 'testuser' } });
        fireEvent.change(screen.getByPlaceholderText('••••••••'), { target: { value: 'password123' } });

        fireEvent.click(screen.getByRole('button', { name: 'auth.sign_in' }));

        await waitFor(() => {
            expect(mockLogin).toHaveBeenCalledWith({ username: 'testuser', password: 'password123' });
        });
    });

    it('switches to register mode and validates password length', async () => {
        render(
            <BrowserRouter>
                <LoginPage />
            </BrowserRouter>
        );

        // Switch to register (Look for button that toggles mode)
        // Text is "auth.create_account" (from t('auth.create_account'))
        // Note: There are two "auth.create_account" texts: one header, one button.
        // The toggle button is at the bottom.
        const toggleButtons = screen.getAllByText('auth.create_account');
        fireEvent.click(toggleButtons[toggleButtons.length - 1]); // Last one is likely the toggle button

        // Fill form with short password
        fireEvent.change(screen.getByPlaceholderText('auth.nickname'), { target: { value: 'New User' } });
        fireEvent.change(screen.getByPlaceholderText('name@example.com'), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByPlaceholderText('••••••••'), { target: { value: '123' } }); // Short password

        // Click submit (Text changes to auth.create_account)
        // We need to target the submit button specifically to avoid confusion with toggle
        const submitButton = screen.getByRole('button', { name: 'auth.create_account' });
        fireEvent.click(submitButton);

        // Multiple error messages may appear, check that at least one exists
        const errorMessages = screen.getAllByText('auth.password_min_length');
        expect(errorMessages.length).toBeGreaterThan(0);
        expect(mockRegister).not.toHaveBeenCalled();
    });
});
