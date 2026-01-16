import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';
import LoginPage from './LoginPage';

// Storage key constant (same as in login-page.tsx) - only username is saved for security
const REMEMBER_USERNAME_KEY = 'gamelink_remembered_username';

// Mock stores
vi.mock('@/stores', () => ({
    useAuthStore: vi.fn(() => ({
        login: vi.fn(),
        register: vi.fn(),
        loading: false,
        error: null,
    })),
}));

// Mock i18next
vi.mock('react-i18next', () => ({
    useTranslation: () => ({
        t: (key: string, options?: { defaultValue?: string }) => options?.defaultValue || key,
        i18n: { language: 'en' },
    }),
}));

// Mock components
vi.mock('@/components/language-switcher', () => ({
    LanguageSwitcher: () => <div data-testid="language-switcher" />,
}));

vi.mock('@/components/mode-toggle', () => ({
    ModeToggle: () => <div data-testid="mode-toggle" />,
}));

vi.mock('sonner', () => ({
    toast: { error: vi.fn(), success: vi.fn() },
}));

import { useAuthStore } from '@/stores';

const renderLoginPage = () => {
    return render(
        <BrowserRouter>
            <LoginPage />
        </BrowserRouter>
    );
};

describe('LoginPage', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        localStorage.clear();
    });

    describe('Remember Me - Load Saved Username', () => {
        it('should load saved username from localStorage on mount', () => {
            localStorage.setItem(REMEMBER_USERNAME_KEY, 'saveduser');

            renderLoginPage();

            const usernameInput = screen.getByPlaceholderText('auth.username') as HTMLInputElement;
            const passwordInput = screen.getByPlaceholderText('••••••••') as HTMLInputElement;

            expect(usernameInput.value).toBe('saveduser');
            // Password should be empty (not saved for security)
            expect(passwordInput.value).toBe('');
        });

        it('should not crash if localStorage has invalid data', () => {
            localStorage.setItem(REMEMBER_USERNAME_KEY, '');

            expect(() => renderLoginPage()).not.toThrow();
        });

        it('should show empty fields if no saved username', () => {
            renderLoginPage();

            const usernameInput = screen.getByPlaceholderText('auth.username') as HTMLInputElement;
            const passwordInput = screen.getByPlaceholderText('••••••••') as HTMLInputElement;

            expect(usernameInput.value).toBe('');
            expect(passwordInput.value).toBe('');
        });
    });

    describe('Remember Me - Save Username', () => {
        it('should save username to localStorage when remember me is checked', async () => {
            const mockLogin = vi.fn().mockResolvedValue(undefined);
            vi.mocked(useAuthStore).mockReturnValue({
                login: mockLogin,
                register: vi.fn(),
                loading: false,
                error: null,
            });

            renderLoginPage();

            const usernameInput = screen.getByPlaceholderText('auth.username');
            const passwordInput = screen.getByPlaceholderText('••••••••');
            const rememberCheckbox = screen.getByRole('checkbox');
            const submitButton = screen.getByRole('button', { name: /auth.sign_in/i });

            await userEvent.type(usernameInput, 'testuser');
            await userEvent.type(passwordInput, 'testpass');

            // Check remember me (default is false when no saved username)
            await userEvent.click(rememberCheckbox);
            expect(rememberCheckbox).toBeChecked();

            fireEvent.click(submitButton);

            await waitFor(() => {
                expect(mockLogin).toHaveBeenCalledWith({
                    username: 'testuser',
                    password: 'testpass',
                });
            });

            // Check localStorage was updated with username only (not password)
            await waitFor(() => {
                const saved = localStorage.getItem(REMEMBER_USERNAME_KEY);
                expect(saved).toBe('testuser');
            });
        });

        it('should remove username from localStorage when remember me is unchecked', async () => {
            // Pre-save a username
            localStorage.setItem(REMEMBER_USERNAME_KEY, 'olduser');

            const mockLogin = vi.fn().mockResolvedValue(undefined);
            vi.mocked(useAuthStore).mockReturnValue({
                login: mockLogin,
                register: vi.fn(),
                loading: false,
                error: null,
            });

            renderLoginPage();

            const usernameInput = screen.getByPlaceholderText('auth.username');
            const passwordInput = screen.getByPlaceholderText('••••••••');
            const rememberCheckbox = screen.getByRole('checkbox');
            const submitButton = screen.getByRole('button', { name: /auth.sign_in/i });

            // Clear and type new values
            await userEvent.clear(usernameInput);
            await userEvent.clear(passwordInput);
            await userEvent.type(usernameInput, 'newuser');
            await userEvent.type(passwordInput, 'newpass');

            // Checkbox should be checked by default since we had saved username
            // Uncheck remember me
            fireEvent.click(rememberCheckbox);
            expect(rememberCheckbox).not.toBeChecked();

            fireEvent.click(submitButton);

            await waitFor(() => {
                expect(mockLogin).toHaveBeenCalled();
            });

            // Check localStorage was cleared
            await waitFor(() => {
                expect(localStorage.getItem(REMEMBER_USERNAME_KEY)).toBeNull();
            });
        });
    });

    describe('Form Validation', () => {
        it('should show error for short password during registration', async () => {
            const mockRegister = vi.fn();
            vi.mocked(useAuthStore).mockReturnValue({
                login: vi.fn(),
                register: mockRegister,
                loading: false,
                error: null,
            });

            renderLoginPage();

            // Switch to register mode
            const toggleButton = screen.getByRole('button', { name: /auth.create_account/i });
            fireEvent.click(toggleButton);

            // Fill form with short password
            const nicknameInput = screen.getByPlaceholderText('auth.nickname');
            const emailInput = screen.getByPlaceholderText('name@example.com');
            const passwordInput = screen.getByPlaceholderText('••••••••');

            await userEvent.type(nicknameInput, 'Test');
            await userEvent.type(emailInput, 'test@example.com');
            await userEvent.type(passwordInput, '12345'); // Only 5 chars

            const submitButton = screen.getByRole('button', { name: /auth.create_account/i });
            fireEvent.click(submitButton);

            // Should show validation error
            await waitFor(() => {
                expect(screen.getByText(/Password must be at least 6 characters/i)).toBeInTheDocument();
            });

            // Register should not be called
            expect(mockRegister).not.toHaveBeenCalled();
        });
    });
});
