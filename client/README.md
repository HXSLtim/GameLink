# GameLink Client

GameLink is a modern Game Companion Platform (PWA) connecting gamers with professional playmates and coaches. Built with React, Vite, and Tailwind CSS.

## Features

- **Store & Orders**
  - Browse verified pro gamers and coaches.
  - Book sessions with real-time pricing (`/orders/create`).
  - Order management (Accept/Reject/Complete) (`/orders`).
  - Review system with ratings and tags.

- **Social & Communication**
  - Real-time Chat (WebSocket) (`/chat`).
  - 1-on-1 private messaging.
  - Online/Offline status indicators.

- **User Profile & Wallet**
  - Custom profiles with avatars and game stats.
  - Favorites system.
  - Wallet balance management (Mock Payment: WeChat/Alipay).
  - VIP Status and privileges.

- **PWA Experience**
  - Installable on Mobile/Desktop.
  - Offline support with network status detection.
  - Mobile-first responsive design.

## Tech Stack

- **Framework**: [React](https://react.dev/) + [Vite](https://vitejs.dev/)
- **Language**: [TypeScript](https://www.typescriptlang.org/)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/) + [Shadcn UI](https://ui.shadcn.com/)
- **State Management**: [Zustand](https://github.com/pmndrs/zustand)
- **Routing**: [React Router v6](https://reactrouter.com/)
- **Icons**: [Lucide React](https://lucide.dev/)
- **Testing**: [Vitest](https://vitest.dev/) + [React Testing Library](https://testing-library.com/)

## Getting Started

### Prerequisites

- Node.js (v18+)
- npm

### Installation

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd client
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Start the development server:

   ```bash
   npm run dev
   ```

4. Build for production:

   ```bash
   npm run build
   ```

### Testing

Run unit and integration tests:

```bash
npm test
```

## Project Structure

```
client/
├── src/
│   ├── components/     # Reusable UI components (Shadcn + Custom)
│   ├── layouts/        # App layouts (AppLayout, AuthLayout)
│   ├── lib/            # Utilities (HTTP, WebSocket, cn)
│   ├── pages/          # Application routes/pages
│   ├── stores/         # Zustand state stores
│   └── test/           # Test setup and utilities
├── public/             # Static assets (PWA icons)
```

## License

MIT
