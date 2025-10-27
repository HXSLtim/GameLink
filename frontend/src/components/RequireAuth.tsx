import { Navigate, useLocation } from 'react-router-dom';
import { Spin } from '@arco-design/web-react';
import { useAuth } from '../contexts/AuthContext';
import styles from './RequireAuth.module.less';

export interface RequireAuthProps {
  /** Child components to render if authenticated */
  children: React.ReactNode;
}

/**
 * Authentication guard component
 *
 * @component
 * @description Protects routes by requiring user authentication.
 * Redirects to login page if user is not authenticated.
 * Shows loading spinner while checking authentication status.
 *
 * @example
 * ```tsx
 * <RequireAuth>
 *   <ProtectedPage />
 * </RequireAuth>
 * ```
 */
export const RequireAuth: React.FC<RequireAuthProps> = ({ children }) => {
  const { token, loading } = useAuth();
  const location = useLocation();

  // Show loading spinner while authentication status is being checked
  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <Spin dot />
      </div>
    );
  }

  // Redirect to login if not authenticated
  if (!token) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }

  // Render protected content
  return <>{children}</>;
};
