import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Card } from './Card';
import styles from './Card.module.less';

describe('Card Component', () => {
  describe('Rendering', () => {
    it('should render children content', () => {
      render(
        <Card>
          <div>Card content</div>
        </Card>,
      );
      expect(screen.getByText('Card content')).toBeInTheDocument();
    });

    it('should render with custom className', () => {
      const { container } = render(<Card className="custom-card">Content</Card>);
      expect(container.firstChild).toHaveClass('custom-card');
    });

    it('should apply custom styles', () => {
      const { container } = render(<Card style={{ padding: '20px' }}>Content</Card>);
      const card = container.firstChild as HTMLElement;
      expect(card).toHaveAttribute('style');
    });
  });

  describe('Title', () => {
    it('should render with title', () => {
      render(<Card title="Card Title">Content</Card>);
      expect(screen.getByText('Card Title')).toBeInTheDocument();
    });

    it('should not render header when title is not provided', () => {
      const { container } = render(<Card>Content</Card>);
      const header = container.querySelector(`.${styles.header}`);
      expect(header).not.toBeInTheDocument();
    });

    it('should render title with ReactNode', () => {
      render(<Card title={<strong>Bold Title</strong>}>Content</Card>);
      expect(screen.getByText('Bold Title')).toBeInTheDocument();
    });
  });

  describe('Extra', () => {
    it('should render extra content', () => {
      render(
        <Card title="Title" extra={<button>Action</button>}>
          Content
        </Card>,
      );
      expect(screen.getByText('Action')).toBeInTheDocument();
    });

    it('should render header with only extra (no title)', () => {
      render(<Card extra={<button>Action</button>}>Content</Card>);
      expect(screen.getByText('Action')).toBeInTheDocument();
    });
  });

  describe('Bordered', () => {
    it('should have bordered class by default', () => {
      const { container } = render(<Card>Content</Card>);
      expect(container.firstChild).toHaveClass(styles.bordered);
    });

    it('should not have bordered class when bordered is false', () => {
      const { container } = render(<Card bordered={false}>Content</Card>);
      expect(container.firstChild).not.toHaveClass(styles.bordered);
    });
  });

  describe('Hoverable', () => {
    it('should not have hoverable class by default', () => {
      const { container } = render(<Card>Content</Card>);
      expect(container.firstChild).not.toHaveClass(styles.hoverable);
    });

    it('should have hoverable class when hoverable is true', () => {
      const { container } = render(<Card hoverable>Content</Card>);
      expect(container.firstChild).toHaveClass(styles.hoverable);
    });
  });

  describe('Structure', () => {
    it('should have card body', () => {
      const { container } = render(<Card>Content</Card>);
      const body = container.querySelector(`.${styles.body}`);
      expect(body).toBeInTheDocument();
      expect(body).toHaveTextContent('Content');
    });

    it('should render complete structure with title and extra', () => {
      const { container } = render(
        <Card title="Title" extra="Extra">
          Body Content
        </Card>,
      );

      expect(container.querySelector(`.${styles.header}`)).toBeInTheDocument();
      expect(container.querySelector(`.${styles.title}`)).toBeInTheDocument();
      expect(container.querySelector(`.${styles.extra}`)).toBeInTheDocument();
      expect(container.querySelector(`.${styles.body}`)).toBeInTheDocument();
    });
  });
});
