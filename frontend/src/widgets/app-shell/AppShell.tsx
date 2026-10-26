import { useState } from 'react'
import { NavLink, Outlet, useLocation, useMatch } from 'react-router'
import { useTranslation } from 'react-i18next'
import { useViewer } from '@/entities/viewer'
import { LogoutButton } from '@/features/auth-logout'
import { LanguageSwitch } from '@/features/language-switch'
import { childNav, mainNav } from '@/shared/config/nav'
import { Icon, Logo, type IconName } from '@/shared/ui'
import { MobileDrawer } from './MobileDrawer'

const navIcons: Record<string, IconName> = {
  'nav:dashboard': 'home', 'nav:profile': 'child', 'nav:schedule': 'calendar',
  'nav:diary': 'diary', 'nav:team': 'people', 'nav:accessLog': 'access',
  'nav:users': 'settings', 'nav:audit': 'shield',
}

export function AppShell() {
  const { t } = useTranslation()
  const viewer = useViewer()
  const role = viewer.data?.role
  const [drawer, setDrawer] = useState(false)
  const { pathname } = useLocation()
  const childMatch = useMatch('/app/children/:childId/*')
  const childId = childMatch?.params.childId
  const items = [...mainNav(role), ...(childId ? childNav(childId, role) : [])]
  const bottom = items.slice(0, 3)
  const active = [...items].reverse().find((item) => pathname === item.to || pathname.startsWith(`${item.to}/`))
  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `relative flex min-h-12 items-center gap-3 rounded-r-xl px-4 py-3 transition-colors before:absolute before:bottom-3 before:left-0 before:top-3 before:w-1 before:origin-center before:rounded-full before:bg-accent before:transition-transform before:duration-200 motion-reduce:before:transition-none motion-reduce:transition-none ${isActive ? 'before:scale-y-100 bg-white/10 font-semibold text-white' : 'before:scale-y-0 text-white/80 hover:bg-white/10 hover:text-white'}`
  return (
    <div className="min-h-dvh bg-canvas text-ink md:flex">
      <a href="#main" className="sr-only focus:not-sr-only focus:fixed focus:left-2 focus:top-2 focus:z-50 focus:rounded focus:bg-surface focus:p-3">
        {t('skipToContent')}
      </a>
      <aside className="relative hidden w-[260px] shrink-0 flex-col overflow-hidden bg-brand px-4 py-6 text-on-brand md:flex">
        <span aria-hidden="true" className="pointer-events-none absolute -left-20 top-36 size-64 rounded-full border border-white/10" />
        <Logo inverse className="mb-10 px-2" />
        <nav aria-label={t('nav:main')} className="flex-1">
          <ul className="space-y-1">
            {items.map((item) => (
              <li key={item.to}>
                <NavLink to={item.to} end={item.end ?? false} className={linkClass}>
                  <Icon name={navIcons[item.labelKey] ?? 'home'} />
                  <span>{t(item.labelKey)}</span>
                </NavLink>
              </li>
            ))}
          </ul>
        </nav>
        <LogoutButton className="mt-6 w-full justify-start text-white hover:bg-white/10 hover:text-white" />
      </aside>
      <div className="min-w-0 flex-1">
        <header className="sticky top-0 z-30 flex min-h-16 items-center gap-3 border-b border-line bg-surface/95 px-4 py-2 backdrop-blur md:px-8">
          <button type="button" className="grid min-h-11 min-w-11 place-items-center rounded-xl border border-line md:hidden" aria-label={t('openMenu')} aria-expanded={drawer} onClick={() => setDrawer(true)}>
            <Icon name="menu" />
          </button>
          <p className="text-lg font-semibold text-ink">{t(active?.labelKey ?? 'brand')}</p>
          <div className="ml-auto flex items-center gap-3">
            <LanguageSwitch />
            <span className="hidden rounded-full border border-line bg-canvas px-3 py-1.5 text-sm sm:inline">
              {viewer.data && `${viewer.data.first_name} · ${t(`roles.${viewer.data.role}`)}`}
            </span>
          </div>
        </header>
        <main id="main" tabIndex={-1} className="page-enter mx-auto min-w-0 max-w-7xl px-4 py-7 pb-24 md:px-8 md:py-10 md:pb-10">
          <Outlet />
        </main>
      </div>
      <nav aria-label={t('menu')} className="fixed inset-x-0 bottom-0 z-30 border-t border-line bg-surface md:hidden">
        <ul className="grid grid-cols-4">
          {bottom.map((item) => (
            <li key={item.to}>
              <NavLink to={item.to} end={item.end ?? false} className="flex min-h-14 flex-col items-center justify-center gap-0.5 px-1 text-center text-xs aria-[current=page]:font-semibold aria-[current=page]:text-brand-strong">
                <Icon name={navIcons[item.labelKey] ?? 'home'} className="size-5" /><span>{t(item.labelKey)}</span>
              </NavLink>
            </li>
          ))}
          <li><button type="button" className="flex min-h-14 w-full flex-col items-center justify-center gap-0.5 text-xs" aria-expanded={drawer} onClick={() => setDrawer(true)}><Icon name="menu" className="size-5" />{t('more')}</button></li>
        </ul>
      </nav>
      <MobileDrawer open={drawer} onOpenChange={setDrawer} items={items} />
    </div>
  )
}
