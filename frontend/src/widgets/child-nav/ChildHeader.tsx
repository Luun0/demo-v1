import { NavLink } from 'react-router'
import { useTranslation } from 'react-i18next'
import type { Role } from '@/shared/api'
import { useChild } from '@/entities/child'
import { childNav } from '@/shared/config/nav'

/** Child name + section tabs (also reachable on mobile). */
export function ChildHeader({ childId, role }: { childId: string; role: Role | undefined }) {
  const { t } = useTranslation()
  const child = useChild(childId)
  return (
    <div className="mb-6">
      <p className="text-sm font-semibold tracking-wide text-brand-ink">{child.data ? `${child.data.first_name} ${child.data.last_name}` : '\u00a0'}</p>
      <nav aria-label={t('nav:childSections')} className="mt-2 overflow-x-auto">
        <ul className="flex gap-1 border-b border-line">
          {childNav(childId, role).map((i) => (
            <li key={i.to}>
              <NavLink to={i.to} end={i.end ?? false} className="inline-flex min-h-11 items-center whitespace-nowrap border-b-2 border-transparent px-3 transition-[border-color] hover:text-brand-ink aria-[current=page]:border-accent aria-[current=page]:font-semibold aria-[current=page]:text-brand-ink motion-reduce:transition-none">
                {t(i.labelKey)}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </div>
  )
}
