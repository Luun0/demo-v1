import { Component, type ReactNode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ApiError } from '@/shared/api'
import { ToastProvider } from '@/shared/ui'

export function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 30_000,
        retry: (n, err) => !(err instanceof ApiError && err.status < 500) && n < 2,
        refetchOnWindowFocus: false,
      },
      mutations: { retry: false },
    },
  })
}

class ErrorBoundary extends Component<{ children: ReactNode }, { failed: boolean }> {
  override state = { failed: false }
  static getDerivedStateFromError() {
    return { failed: true }
  }
  override render() {
    if (this.state.failed) {
      return (
        <div role="alert" className="p-8 text-center">
          <p className="text-lg font-semibold">Something went wrong / Қате орын алды / Произошла ошибка</p>
          <button type="button" className="mt-4 min-h-11 rounded-lg border px-4" onClick={() => window.location.reload()}>↻</button>
        </div>
      )
    }
    return this.props.children
  }
}

export function Providers({ client, children }: { client: QueryClient; children: ReactNode }) {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={client}>
        <ToastProvider>{children}</ToastProvider>
      </QueryClientProvider>
    </ErrorBoundary>
  )
}
