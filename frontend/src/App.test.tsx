import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { App } from './App';

describe('App', () => {
  it('renders app title', () => {
    render(<App />);
    expect(screen.getByText(/GameLink 管理系统/i)).toBeDefined();
  });
});
