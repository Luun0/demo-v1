import { Navigate, Outlet, useLocation } from 'react-router'
import type { Role } from '@/shared/api'
import { ErrorState, Loading } from '@/shared/ui'
import { useAccessToken, useViewer } from '@/entities/viewer'
import { NotFoundPage } from '@/pages/not-found'

export function RequireAuth() {
  const token = useAccessToken()
  const viewer = useViewer()
  const loc = useLocation()
  if (!token) return <Navigate to="/login" replace state={{ from: loc.pathname }} />
  if (viewer.isPending) return <Loading />
  if (viewer.isError) return <ErrorState error={viewer.error} onRetry={() => void viewer.refetch()} />
  return <Outlet />
}

export function RequireRole({ roles }: { roles: Role[] }) {
  const viewer = useViewer()
  if (!viewer.data || !roles.includes(viewer.data.role)) return <NotFoundPage kind="forbidden" />
  return <Outlet />
}

export function GuestOnly() {
  return useAccessToken() ? <Navigate to="/app" replace /> : <Outlet />
}
