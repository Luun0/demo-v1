import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from 'react-router'
import { refreshSession } from '@/shared/api'
import { initI18n } from '@/shared/i18n'
import { makeQueryClient, Providers } from './providers'
import { createRouter } from './router'
import './styles.css'

async function start() {
  // Restore the session from the HttpOnly refresh cookie; tokens live in memory only.
  await Promise.all([initI18n(), refreshSession().catch(() => false)])
  const root = document.getElementById('root')
  if (!root) throw new Error('#root missing')
  createRoot(root).render(
    <StrictMode>
      <Providers client={makeQueryClient()}>
        <RouterProvider router={createRouter()} />
      </Providers>
    </StrictMode>,
  )
}

void start()
