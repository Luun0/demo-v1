import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'

export function NotFoundPage({ kind = 'notFound' }: { kind?: 'notFound' | 'forbidden' }) {
  const { t } = useTranslation('errors')
  return (
    <div className="page-enter mx-auto my-12 max-w-xl rounded-3xl border border-line bg-surface px-6 py-14 text-center shadow-soft sm:px-12">
      <span className="mx-auto mb-6 grid size-16 place-items-center rounded-2xl bg-brand/10 text-2xl font-bold text-brand-strong" aria-hidden="true">?</span>
      <h1 className="text-3xl font-bold tracking-tight">{t(`${kind}Title`)}</h1>
      <p className="mt-3 text-ink-soft">{t(`${kind}Text`)}</p>
      <Link to="/app" className="mt-6 inline-flex min-h-11 items-center rounded-lg bg-brand px-5 font-semibold text-on-brand">{t('home')}</Link>
    </div>
  )
}
