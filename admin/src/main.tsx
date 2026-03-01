import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { cleanupLegacyServiceWorker } from './pwa/cleanupLegacyServiceWorker'

void cleanupLegacyServiceWorker()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
