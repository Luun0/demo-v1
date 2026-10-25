import { useState } from 'react'
import { Link, NavLink, Outlet, useMatch } from 'react-router'
import { useTranslation } from 'react-i18next'
import { useViewer } from '@/entities/viewer'
import { LogoutButton } from '@/features/auth-logout'
import { LanguageSwitch } from '@/features/language-switch'
import { childNav, mainNav } from '@/shared/config/nav'
import { MobileDrawer } from './MobileDrawer'

const linkCls = 'block min-h-11 rounded-lg px-3 py-2.5 hover:bg-muted aria-[current=page]:bg-brand/15 aria-[current=page]:font-semibold'

export function AppShell() {
  const { t } = useTranslation()
  const viewer = useViewer()
  const role = viewer.data?.role
  const [drawer, setDrawer] = useState(false)
  const childMatch = useMatch('/app/children/:childId/*')
  const childId = childMatch?.params.childId
  const items = [...mainNav(role), ...(childId ? childNav(childId, role) : [])]
  const bottom = items.slice(0, 3)
  return (
    <div className="min-h-dvh bg-canvas text-ink">
      <a href="#main" className="sr-only focus:not-sr-only focus:fixed focus:left-2 focus:top-2 focus:z-50 focus:rounded focus:bg-surface focus:p-3">
        {t('skipToContent')}
      </a>
      <header className="sticky top-0 z-30 flex items-center gap-3 border-b border-line bg-surface/95 px-4 py-2 backdrop-blur">
        <button type="button" className="min-h-11 min-w-11 rounded-lg border border-line md:hidden" aria-label={t('openMenu')} aria-expanded={drawer} onClick={() => setDrawer(true)}>
          <span aria-hidden="true">☰</span>
        </button>
        <Link to="/app" className="text-lg font-bold text-brand-strong">{t('brand')}</Link>
        <span className="ml-auto hidden text-sm text-ink-soft sm:inline">
          {viewer.data && `${viewer.data.first_name} · ${t(`roles.${viewer.data.role}`)}`}
        </span>
        <div className="hidden items-center gap-3 md:flex">
          <LanguageSwitch />
          <LogoutButton />
        </div>
      </header>
      <div className="mx-auto flex max-w-6xl gap-6 px-4 py-6">
        <nav aria-label={t('nav:main')} className="hidden w-56 shrink-0 md:block">
          <ul className="sticky top-20 space-y-1">
            {items.map((i) => (
              <li key={i.to}><NavLink to={i.to} end={i.end ?? false} className={linkCls}>{t(i.labelKey)}</NavLink></li>
            ))}
          </ul>
        </nav>
        <main id="main" tabIndex={-1} className="min-w-0 flex-1 pb-24 md:pb-6">
          <Outlet />
        </main>
      </div>
      <nav aria-label={t('menu')} className="fixed inset-x-0 bottom-0 z-30 border-t border-line bg-surface md:hidden">
        <ul className="grid grid-cols-4">
          {bottom.map((i) => (
            <li key={i.to}><NavLink to={i.to} end={i.end ?? false} className="flex min-h-14 items-center justify-center px-1 text-center text-sm aria-[current=page]:font-semibold aria-[current=page]:text-brand-strong">{t(i.labelKey)}</NavLink></li>
          ))}
          <li>
            <button type="button" className="flex min-h-14 w-full items-center justify-center text-sm" aria-expanded={drawer} onClick={() => setDrawer(true)}>{t('more')}</button>
          </li>
        </ul>
      </nav>
      <MobileDrawer open={drawer} onOpenChange={setDrawer} items={items} />
    </div>
  )
}
