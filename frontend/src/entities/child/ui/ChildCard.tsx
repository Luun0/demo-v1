import type { ReactNode } from 'react'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import type { Child } from '@/shared/api'
import { formatDate } from '@/shared/lib'
import type { Role } from '@/shared/api'
import { childNav } from '@/shared/config/nav'

export function ChildCard({ child, extra, role }: { child: Child; extra?: ReactNode; role: Role }) {
  const { t, i18n } = useTranslation()
  const { t: dashboardT } = useTranslation('dashboard')
  const sections = childNav(child.id, role)
  const profile = sections.find((item) => item.end)
  return (
    <article className="surface-lift rounded-2xl border border-line bg-surface p-5 shadow-soft">
      <h3 className="text-lg font-semibold text-ink">
        {child.first_name} {child.last_name}
      </h3>
      <p className="text-ink-soft">{dashboardT('born', { date: formatDate(child.birth_date, i18n.language) })}</p>
      {extra}
      {profile && <Link to={profile.to} className="mt-2 inline-flex min-h-10 items-center font-medium text-brand-ink underline underline-offset-4">{dashboardT('open')} <span className="sr-only">: {child.first_name}</span></Link>}
      <nav aria-label={`${t('nav:childSections')}: ${child.first_name}`} className="mt-3 flex flex-wrap gap-2 border-t border-line pt-3">
        {sections.filter((item) => !item.end).map((item) => <Link key={item.to} to={item.to} className="inline-flex min-h-10 items-center rounded-lg border border-line px-3 text-sm font-medium text-brand-ink hover:bg-muted">{t(item.labelKey)}</Link>)}
      </nav>
    </article>
  )
}
