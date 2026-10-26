import type { ReactNode } from 'react'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { LanguageSwitch } from '@/features/language-switch'
import { ThemeSwitch } from '@/features/theme-switch'
import { Logo } from '@/shared/ui'

export function AuthLayout({ title, children }: { title: string; children: ReactNode }) {
  const { t } = useTranslation()
  return (
    <div className="min-h-dvh bg-canvas text-ink md:grid md:grid-cols-2">
      <a href="#main" className="sr-only focus:not-sr-only focus:fixed focus:left-2 focus:top-2 focus:z-50 focus:rounded focus:bg-surface focus:p-3">{t('skipToContent')}</a>
      <section className="relative hidden overflow-hidden bg-ink-panel px-10 py-9 text-panel-ink md:flex md:flex-col lg:px-16">
        <Link to="/" aria-label={t('brand')}><Logo inverse /></Link>
        <div className="relative z-10 my-auto max-w-xl py-14">
          <p className="text-sm font-semibold uppercase tracking-[0.16em] text-lagoon">{t('auth:storyEyebrow')}</p>
          <h1 className="mt-5 text-4xl font-semibold leading-[1.08] tracking-[-0.05em] lg:text-6xl">{t('auth:storyTitle')}</h1>
          <p className="mt-5 max-w-lg text-lg leading-relaxed text-panel-ink/80">{t('auth:storyText')}</p>
        </div>
        <span aria-hidden="true" className="pointer-events-none absolute -bottom-12 -right-10 size-72 rounded-full border-[3rem] border-panel-ink/[0.04]" />
        <span aria-hidden="true" className="absolute bottom-0 left-10 h-1 w-32 bg-accent lg:left-16" />
      </section>
      <div className="flex min-h-dvh flex-col px-4 py-6 sm:px-8 md:px-10 lg:px-16">
        <header className="mb-6 flex items-center justify-between md:justify-end"><Link to="/" className="md:hidden"><Logo /></Link><div className="flex items-center gap-3"><ThemeSwitch /><LanguageSwitch /></div></header>
        <main id="main" className="page-enter my-auto w-full self-center rounded-[1.25rem] border border-line bg-surface p-6 shadow-soft sm:p-9 lg:max-w-xl">
          <h2 className="mb-6 text-2xl font-semibold">{title}</h2>
          {children}
        </main>
      </div>
    </div>
  )
}
