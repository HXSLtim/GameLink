import { render, screen } from '@testing-library/react';
import { RouterProvider, createMemoryRouter } from 'react-router-dom';
import { App } from './App';

describe('App', () => {
  it('renders header title', () => {
    const router = createMemoryRouter(
      [{ path: '/', element: <App />, children: [{ index: true, element: <div /> }] }],
      { future: { v7_relativeSplatPath: true } },
    );
    render(<RouterProvider router={router} future={{ v7_startTransition: true }} />);
    expect(screen.getByText('GameLink 管理端')).toBeInTheDocument();
  });
});
