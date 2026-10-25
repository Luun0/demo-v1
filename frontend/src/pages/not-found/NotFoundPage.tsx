import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'

export function NotFoundPage({ kind = 'notFound' }: { kind?: 'notFound' | 'forbidden' }) {
  const { t } = useTranslation('errors')
  return (
    <div className="mx-auto max-w-lg py-16 text-center">
      <h1 className="text-3xl font-bold">{t(`${kind}Title`)}</h1>
      <p className="mt-3 text-ink-soft">{t(`${kind}Text`)}</p>
      <Link to="/app" className="mt-6 inline-flex min-h-11 items-center rounded-lg bg-brand px-5 font-semibold text-on-brand">{t('home')}</Link>
    </div>
  )
}
