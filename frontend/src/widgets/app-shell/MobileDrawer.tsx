import { useRef } from 'react'
import * as RD from '@radix-ui/react-dialog'
import { NavLink } from 'react-router'
import { useTranslation } from 'react-i18next'
import { LogoutButton } from '@/features/auth-logout'
import { LanguageSwitch } from '@/features/language-switch'
import { ThemeSwitch } from '@/features/theme-switch'
import type { NavItem } from '@/shared/config/nav'
import { Icon, Logo, type IconName } from '@/shared/ui'

const navIcons: Record<string, IconName> = {
  'nav:dashboard': 'home', 'nav:profile': 'child', 'nav:schedule': 'calendar',
  'nav:diary': 'diary', 'nav:team': 'people', 'nav:accessLog': 'access',
  'nav:users': 'people', 'nav:audit': 'shield',
}

/** Mobile drawer: Radix dialog gives focus trap, Esc and focus return. */
export function MobileDrawer({ open, onOpenChange, items }: { open: boolean; onOpenChange: (o: boolean) => void; items: NavItem[] }) {
  const { t } = useTranslation()
  // The drawer has two openers (hamburger, bottom-nav "More"); return focus to whichever was used.
  const opener = useRef<HTMLElement | null>(null)
  return (
    <RD.Root open={open} onOpenChange={onOpenChange}>
      <RD.Portal>
        <RD.Overlay className="fixed inset-0 z-40 bg-overlay backdrop-blur-sm" />
        <RD.Content className="fixed inset-y-0 left-0 z-50 w-72 overflow-y-auto bg-brand p-5 text-on-brand shadow-panel" aria-describedby={undefined}
          onOpenAutoFocus={() => { opener.current = document.activeElement as HTMLElement | null }}
          onCloseAutoFocus={(e) => { e.preventDefault(); opener.current?.focus() }}>
          <div className="mb-4 flex items-center justify-between">
            <RD.Title><Logo inverse /></RD.Title>
            <RD.Close className="min-h-11 min-w-11 rounded-lg transition-[background-color] hover:bg-on-brand/10 focus-visible:ring-2 focus-visible:ring-accent motion-reduce:transition-none" aria-label={t('closeMenu')}>
              <Icon name="close" />
            </RD.Close>
          </div>
          <nav aria-label={t('nav:main')}>
            <ul className="space-y-1">
              {items.map((i) => (
                <li key={i.to}>
                  <NavLink to={i.to} end={i.end ?? false} onClick={() => onOpenChange(false)} className="flex min-h-12 items-center gap-3 rounded-xl px-3 py-2.5 text-on-brand transition-[background-color] hover:bg-on-brand/10 aria-[current=page]:bg-on-brand/15 aria-[current=page]:font-semibold motion-reduce:transition-none">
                    <Icon name={navIcons[i.labelKey] ?? 'home'} />{t(i.labelKey)}
                  </NavLink>
                </li>
              ))}
            </ul>
          </nav>
          <div className="mt-6 space-y-3 border-t border-on-brand/20 pt-4">
            <ThemeSwitch dark />
            <LanguageSwitch dark />
            <LogoutButton className="w-full justify-start text-on-brand hover:bg-on-brand/10 hover:text-on-brand" />
          </div>
        </RD.Content>
      </RD.Portal>
    </RD.Root>
  )
}
