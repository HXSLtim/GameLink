import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';
import LoginPage from './login-page';

// Storage key constant (same as in login-page.tsx)
const REMEMBER_CREDENTIALS_KEY = 'gamelink_remembered_credentials';

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

    describe('Remember Me - Load Saved Credentials', () => {
        it('should load saved credentials from localStorage on mount', () => {
            const savedCredentials = { username: 'saveduser', password: 'savedpass' };
            localStorage.setItem(REMEMBER_CREDENTIALS_KEY, JSON.stringify(savedCredentials));

            renderLoginPage();

            const usernameInput = screen.getByPlaceholderText('auth.username') as HTMLInputElement;
            const passwordInput = screen.getByPlaceholderText('••••••••') as HTMLInputElement;

            expect(usernameInput.value).toBe('saveduser');
            expect(passwordInput.value).toBe('savedpass');
        });

        it('should not crash if localStorage has invalid JSON', () => {
            localStorage.setItem(REMEMBER_CREDENTIALS_KEY, 'invalid-json');

            expect(() => renderLoginPage()).not.toThrow();
        });

        it('should show empty fields if no saved credentials', () => {
            renderLoginPage();

            const usernameInput = screen.getByPlaceholderText('auth.username') as HTMLInputElement;
            const passwordInput = screen.getByPlaceholderText('••••••••') as HTMLInputElement;

            expect(usernameInput.value).toBe('');
            expect(passwordInput.value).toBe('');
        });
    });

    describe('Remember Me - Save Credentials', () => {
        it('should save credentials to localStorage when remember me is checked', async () => {
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

            // Ensure remember me is checked (default is true)
            expect(rememberCheckbox).toBeChecked();

            fireEvent.click(submitButton);

            await waitFor(() => {
                expect(mockLogin).toHaveBeenCalledWith({
                    username: 'testuser',
                    password: 'testpass',
                });
            });

            // Check localStorage was updated
            await waitFor(() => {
                const saved = localStorage.getItem(REMEMBER_CREDENTIALS_KEY);
                expect(saved).not.toBeNull();
                const parsed = JSON.parse(saved!);
                expect(parsed.username).toBe('testuser');
                expect(parsed.password).toBe('testpass');
            });
        });

        it('should remove credentials from localStorage when remember me is unchecked', async () => {
            // Pre-save some credentials
            localStorage.setItem(REMEMBER_CREDENTIALS_KEY, JSON.stringify({ username: 'old', password: 'old' }));

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

            // Uncheck remember me
            fireEvent.click(rememberCheckbox);
            expect(rememberCheckbox).not.toBeChecked();

            fireEvent.click(submitButton);

            await waitFor(() => {
                expect(mockLogin).toHaveBeenCalled();
            });

            // Check localStorage was cleared
            await waitFor(() => {
                expect(localStorage.getItem(REMEMBER_CREDENTIALS_KEY)).toBeNull();
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
