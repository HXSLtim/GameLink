import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { RequireAuth } from './RequireAuth';
import { AuthProvider } from '../contexts/AuthContext';
import * as authService from '../services/auth';

// Mock auth service
vi.mock('../services/auth', () => ({
  authService: {
    me: vi.fn(),
    login: vi.fn(),
  },
}));

const TestComponent = () => <div>Protected Content</div>;
const LoginComponent = () => <div>Login Page</div>;

const MockedAuthApp = ({ hasToken = false }: { hasToken?: boolean }) => {
  // Mock localStorage
  if (hasToken) {
    localStorage.setItem('gamelink_token', 'test-token');
    vi.mocked(authService.authService.me).mockResolvedValue({
      id: 1,
      username: 'testuser',
      role: 'admin',
    } as any);
  } else {
    localStorage.removeItem('gamelink_token');
    vi.mocked(authService.authService.me).mockResolvedValue(null as any);
  }

  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<LoginComponent />} />
          <Route
            path="/"
            element={
              <RequireAuth>
                <TestComponent />
              </RequireAuth>
            }
          />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
};

describe('RequireAuth', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it('should show loading spinner while checking auth', () => {
    render(<MockedAuthApp hasToken={true} />);

    // Initially should show loading spinner
    const spinner = document.querySelector('.arco-spin');
    expect(spinner).toBeInTheDocument();
  });

  it('should render children when authenticated', async () => {
    render(<MockedAuthApp hasToken={true} />);

    await waitFor(() => {
      expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });
  });

  it('should redirect to login when not authenticated', async () => {
    render(<MockedAuthApp hasToken={false} />);

    await waitFor(() => {
      expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    });
  });

  it('should check for token in localStorage', async () => {
    localStorage.setItem('gamelink_token', 'test-token');
    vi.mocked(authService.authService.me).mockResolvedValue({
      id: '1',
      username: 'testuser',
      role: 'admin',
    } as any);

    render(<MockedAuthApp hasToken={true} />);

    await waitFor(() => {
      expect(authService.authService.me).toHaveBeenCalled();
    });
  });
});

