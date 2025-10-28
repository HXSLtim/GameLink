import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Button } from './Button';
import styles from './Button.module.less';

describe('Button Component', () => {
  describe('Rendering', () => {
    it('should render with text content', () => {
      render(<Button>Click me</Button>);
      expect(screen.getByRole('button', { name: /click me/i })).toBeInTheDocument();
    });

    it('should render with custom className', () => {
      const { container } = render(<Button className="custom-class">Test</Button>);
      expect(container.firstChild).toHaveClass('custom-class');
    });
  });

  describe('Variants', () => {
    it('should render primary variant by default', () => {
      const { container } = render(<Button>Primary</Button>);
      expect(container.firstChild).toHaveClass(styles.primary);
    });

    it('should render secondary variant', () => {
      const { container } = render(<Button variant="secondary">Secondary</Button>);
      expect(container.firstChild).toHaveClass(styles.secondary);
    });

    it('should render text variant', () => {
      const { container } = render(<Button variant="text">Text</Button>);
      expect(container.firstChild).toHaveClass(styles.text);
    });

    it('should render outlined variant', () => {
      const { container } = render(<Button variant="outlined">Outlined</Button>);
      expect(container.firstChild).toHaveClass(styles.outlined);
    });
  });

  describe('Sizes', () => {
    it('should render small size', () => {
      const { container } = render(<Button size="small">Small</Button>);
      expect(container.firstChild).toHaveClass(styles.small);
    });

    it('should render medium size by default', () => {
      const { container } = render(<Button>Medium</Button>);
      expect(container.firstChild).toHaveClass(styles.medium);
    });

    it('should render large size', () => {
      const { container } = render(<Button size="large">Large</Button>);
      expect(container.firstChild).toHaveClass(styles.large);
    });
  });

  describe('States', () => {
    it('should be disabled when disabled prop is true', () => {
      render(<Button disabled>Disabled</Button>);
      expect(screen.getByRole('button')).toBeDisabled();
    });

    it('should be disabled when loading prop is true', () => {
      render(<Button loading>Loading</Button>);
      expect(screen.getByRole('button')).toBeDisabled();
    });

    it('should render block button', () => {
      const { container } = render(<Button block>Block</Button>);
      expect(container.firstChild).toHaveClass(styles.block);
    });

    it('should have loading class when loading', () => {
      const { container } = render(<Button loading>Loading</Button>);
      expect(container.firstChild).toHaveClass(styles.loading);
    });
  });

  describe('Events', () => {
    it('should handle click events', () => {
      const handleClick = vi.fn();
      render(<Button onClick={handleClick}>Click</Button>);

      fireEvent.click(screen.getByRole('button'));
      expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it('should not call onClick when disabled', () => {
      const handleClick = vi.fn();
      render(
        <Button disabled onClick={handleClick}>
          Disabled
        </Button>,
      );

      const button = screen.getByRole('button');
      fireEvent.click(button);
      expect(handleClick).not.toHaveBeenCalled();
    });

    it('should not call onClick when loading', () => {
      const handleClick = vi.fn();
      render(
        <Button loading onClick={handleClick}>
          Loading
        </Button>,
      );

      const button = screen.getByRole('button');
      fireEvent.click(button);
      expect(handleClick).not.toHaveBeenCalled();
    });
  });

  describe('Icon', () => {
    it('should render with icon', () => {
      const icon = <span data-testid="icon">🔍</span>;
      render(<Button icon={icon}>Search</Button>);

      expect(screen.getByTestId('icon')).toBeInTheDocument();
      expect(screen.getByText('Search')).toBeInTheDocument();
    });

    it('should render icon only button', () => {
      const icon = <span data-testid="icon-only">🔍</span>;
      render(<Button icon={icon} aria-label="Search" />);

      expect(screen.getByTestId('icon-only')).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    it('should support type attribute', () => {
      render(<Button type="submit">Submit</Button>);
      expect(screen.getByRole('button')).toHaveAttribute('type', 'submit');
    });

    it('should support aria-label', () => {
      render(<Button aria-label="Close modal">×</Button>);
      expect(screen.getByLabelText('Close modal')).toBeInTheDocument();
    });
  });
});
