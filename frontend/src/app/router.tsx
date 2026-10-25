import { lazy, Suspense, type ComponentType } from 'react'
import { createBrowserRouter, type RouteObject } from 'react-router'
import { Loading } from '@/shared/ui'
import { AppShell } from '@/widgets/app-shell'
import { LandingPage } from '@/pages/landing'
import { LoginPage } from '@/pages/login'
import { NotFoundPage } from '@/pages/not-found'
import { GuestOnly, RequireAuth, RequireRole } from './guards'

function lazyPage<K extends string>(load: () => Promise<Record<K, ComponentType>>, name: K) {
  const C = lazy(async (): Promise<{ default: ComponentType }> => ({ default: (await load())[name] }))
  return (
    <Suspense fallback={<Loading />}>
      <C />
    </Suspense>
  )
}

export const routes: RouteObject[] = [
  { path: '/', element: <LandingPage /> },
  {
    element: <GuestOnly />,
    children: [
      { path: '/login', element: <LoginPage /> },
      { path: '/register', element: lazyPage(() => import('@/pages/register'), 'RegisterPage') },
    ],
  },
  {
    path: '/app',
    element: <RequireAuth />,
    children: [
      {
        element: <AppShell />,
        children: [
          { index: true, element: lazyPage(() => import('@/pages/dashboard'), 'DashboardPage') },
          {
            path: 'children/:childId',
            children: [
              { index: true, element: lazyPage(() => import('@/pages/child-profile'), 'ChildProfilePage') },
              { path: 'team', element: lazyPage(() => import('@/pages/team'), 'TeamPage') },
              {
                element: <RequireRole roles={['parent', 'specialist']} />,
                children: [
                  { path: 'schedule', element: lazyPage(() => import('@/pages/schedule'), 'SchedulePage') },
                  { path: 'diary', element: lazyPage(() => import('@/pages/diary'), 'DiaryPage') },
                ],
              },
              {
                element: <RequireRole roles={['parent', 'admin']} />,
                children: [{ path: 'access-log', element: lazyPage(() => import('@/pages/access-log'), 'AccessLogPage') }],
              },
            ],
          },
          {
            path: 'admin',
            element: <RequireRole roles={['admin']} />,
            children: [
              { path: 'users', element: lazyPage(() => import('@/pages/admin-users'), 'AdminUsersPage') },
              { path: 'audit', element: lazyPage(() => import('@/pages/admin-audit'), 'AdminAuditPage') },
            ],
          },
          { path: '*', element: <NotFoundPage /> },
        ],
      },
    ],
  },
  { path: '*', element: <NotFoundPage /> },
]

export const createRouter = () => createBrowserRouter(routes)
