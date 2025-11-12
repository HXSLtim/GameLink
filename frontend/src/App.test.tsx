import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react';
import { App } from './App';

describe('App', () => {
  it('renders without crashing', () => {
    // App已经包含了所有必要的Provider (router已内置AuthProvider和ThemeProvider)
    const { container } = render(<App />);
    
    // 验证App组件能够正常渲染
    expect(container).toBeDefined();
    expect(container.firstChild).toBeDefined();
  });
});
