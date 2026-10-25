import type { ReactNode } from 'react'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import type { Child } from '@/shared/api'
import { formatDate } from '@/shared/lib'

export function ChildCard({ child, extra }: { child: Child; extra?: ReactNode }) {
  const { t, i18n } = useTranslation('dashboard')
  return (
    <article className="rounded-xl border border-line bg-surface p-5 shadow-sm">
      <h3 className="text-lg font-semibold text-ink">
        {child.first_name} {child.last_name}
      </h3>
      <p className="text-ink-soft">{t('born', { date: formatDate(child.birth_date, i18n.language) })}</p>
      {extra}
      <Link to={`/app/children/${child.id}`} className="mt-3 inline-flex min-h-11 items-center font-medium text-brand-strong underline">
        {t('open')} <span className="sr-only">: {child.first_name}</span>
      </Link>
    </article>
  )
}
