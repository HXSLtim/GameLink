import { RouterProvider } from 'react-router-dom';
import { router } from './router';

/**
 * App root component
 */
export const App = () => {
  return <RouterProvider router={router} />;
};
