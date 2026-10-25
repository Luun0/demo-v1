import { useRef } from 'react'
import * as RD from '@radix-ui/react-dialog'
import { NavLink } from 'react-router'
import { useTranslation } from 'react-i18next'
import { LogoutButton } from '@/features/auth-logout'
import { LanguageSwitch } from '@/features/language-switch'
import type { NavItem } from '@/shared/config/nav'

/** Mobile drawer: Radix dialog gives focus trap, Esc and focus return. */
export function MobileDrawer({ open, onOpenChange, items }: { open: boolean; onOpenChange: (o: boolean) => void; items: NavItem[] }) {
  const { t } = useTranslation()
  // The drawer has two openers (hamburger, bottom-nav "More"); return focus to whichever was used.
  const opener = useRef<HTMLElement | null>(null)
  return (
    <RD.Root open={open} onOpenChange={onOpenChange}>
      <RD.Portal>
        <RD.Overlay className="fixed inset-0 z-40 bg-black/50" />
        <RD.Content className="fixed inset-y-0 left-0 z-50 w-72 overflow-y-auto bg-surface p-5 text-ink shadow-xl" aria-describedby={undefined}
          onOpenAutoFocus={() => { opener.current = document.activeElement as HTMLElement | null }}
          onCloseAutoFocus={(e) => { e.preventDefault(); opener.current?.focus() }}>
          <div className="mb-4 flex items-center justify-between">
            <RD.Title className="text-lg font-semibold">{t('menu')}</RD.Title>
            <RD.Close className="min-h-11 min-w-11 rounded-lg hover:bg-muted" aria-label={t('closeMenu')}>
              <span aria-hidden="true">×</span>
            </RD.Close>
          </div>
          <nav aria-label={t('nav:main')}>
            <ul className="space-y-1">
              {items.map((i) => (
                <li key={i.to}>
                  <NavLink to={i.to} end={i.end ?? false} onClick={() => onOpenChange(false)} className="block min-h-11 rounded-lg px-3 py-2.5 hover:bg-muted aria-[current=page]:bg-brand/15 aria-[current=page]:font-semibold">
                    {t(i.labelKey)}
                  </NavLink>
                </li>
              ))}
            </ul>
          </nav>
          <div className="mt-6 space-y-3 border-t border-line pt-4">
            <LanguageSwitch />
            <LogoutButton className="w-full justify-start" />
          </div>
        </RD.Content>
      </RD.Portal>
    </RD.Root>
  )
}
