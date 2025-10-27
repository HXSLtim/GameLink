import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Footer } from './Footer';

describe('Footer', () => {
  it('should render footer with copyright text', () => {
    render(<Footer />);

    const currentYear = new Date().getFullYear();
    expect(screen.getByText(`© ${currentYear} GameLink Admin`)).toBeInTheDocument();
  });

  it('should use Layout.Footer component', () => {
    const { container } = render(<Footer />);

    const footer = container.querySelector('.arco-layout-footer');
    expect(footer).toBeInTheDocument();
  });

  it('should display year dynamically', () => {
    const { rerender } = render(<Footer />);

    const currentYear = new Date().getFullYear();
    expect(screen.getByText(new RegExp(currentYear.toString()))).toBeInTheDocument();

    // Verify it's using current year, not hardcoded
    const copyrightText = screen.getByText(/©.*GameLink Admin/);
    expect(copyrightText.textContent).toContain(currentYear.toString());

    rerender(<Footer />);
    expect(screen.getByText(new RegExp(currentYear.toString()))).toBeInTheDocument();
  });
});


