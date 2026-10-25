import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { LanguageSwitch } from '@/features/language-switch'

export function LandingPage() {
  const { t } = useTranslation('landing')
  return (
    <div className="min-h-dvh bg-canvas text-ink">
      <a href="#main" className="sr-only focus:not-sr-only focus:fixed focus:left-2 focus:top-2 focus:rounded focus:bg-surface focus:p-3">{t('common:skipToContent')}</a>
      <header className="mx-auto flex max-w-5xl items-center justify-between gap-3 px-4 py-4">
        <span className="text-xl font-bold text-brand-strong">{t('common:brand')}</span>
        <div className="flex items-center gap-3">
          <LanguageSwitch />
          <Link to="/login" className="inline-flex min-h-11 items-center rounded-lg px-3 font-medium hover:bg-muted">{t('ctaLogin')}</Link>
        </div>
      </header>
      <main id="main" className="mx-auto max-w-5xl px-4 pb-16">
        <section className="py-12">
          <p className="font-medium text-brand-strong">{t('eyebrow')}</p>
          <h1 className="mt-2 max-w-3xl text-4xl font-bold leading-tight">{t('title')}</h1>
          <p className="mt-4 max-w-2xl text-lg text-ink-soft">{t('text')}</p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link to="/register" className="inline-flex min-h-11 items-center rounded-lg bg-brand px-5 font-semibold text-on-brand hover:bg-brand-strong">{t('ctaRegister')}</Link>
            <Link to="/login" className="inline-flex min-h-11 items-center rounded-lg border border-line px-5 font-semibold hover:bg-muted">{t('ctaLogin')}</Link>
          </div>
        </section>
        <section aria-label={t('common:brand')} className="grid gap-4 md:grid-cols-3">
          {(['f1', 'f2', 'f3'] as const).map((f) => (
            <article key={f} className="rounded-xl border border-line bg-surface p-5">
              <h2 className="text-lg font-semibold">{t(`${f}t`)}</h2>
              <p className="mt-2 text-ink-soft">{t(`${f}d`)}</p>
            </article>
          ))}
        </section>
      </main>
    </div>
  )
}
