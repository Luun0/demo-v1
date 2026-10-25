import type { ReactElement, ReactNode } from 'react'
import { render } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router'
import { ToastProvider } from '@/shared/ui'

export function renderWithProviders(ui: ReactElement, { route = '/' } = {}) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  const Wrapper = ({ children }: { children: ReactNode }) => (
    <QueryClientProvider client={qc}>
      <ToastProvider>
        <MemoryRouter initialEntries={[route]}>{children}</MemoryRouter>
      </ToastProvider>
    </QueryClientProvider>
  )
  return render(ui, { wrapper: Wrapper })
}
