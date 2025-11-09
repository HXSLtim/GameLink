import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react';
import { App } from './App';
import { AuthProvider } from './contexts/AuthContext';
import { ThemeProvider } from './contexts/ThemeContext';

describe('App', () => {
  it('renders without crashing', () => {
    const { container } = render(
      <ThemeProvider>
        <AuthProvider>
          <App />
        </AuthProvider>
      </ThemeProvider>
    );
    // 验证App组件能够正常渲染
    expect(container).toBeDefined();
    expect(container.firstChild).toBeDefined();
  });
});
