import type { ReactNode } from 'react'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { LanguageSwitch } from '@/features/language-switch'
import { Icon, Logo } from '@/shared/ui'

export function AuthLayout({ title, children }: { title: string; children: ReactNode }) {
  const { t } = useTranslation()
  return (
    <div className="min-h-dvh bg-canvas text-ink md:grid md:grid-cols-2">
      <a href="#main" className="sr-only focus:not-sr-only focus:fixed focus:left-2 focus:top-2 focus:z-50 focus:rounded focus:bg-surface focus:p-3">{t('skipToContent')}</a>
      <section className="relative hidden overflow-hidden bg-brand px-10 py-9 text-white md:flex md:flex-col lg:px-16">
        <Link to="/" aria-label={t('brand')}><Logo inverse /></Link>
        <div className="relative z-10 my-auto max-w-xl py-14">
          <p className="text-sm font-semibold uppercase tracking-[0.16em] text-white/75">{t('auth:storyEyebrow')}</p>
          <h1 className="mt-5 text-4xl font-bold leading-tight lg:text-5xl">{t('auth:storyTitle')}</h1>
          <p className="mt-5 text-lg leading-relaxed text-white/80">{t('auth:storyText')}</p>
        </div>
        <span className="absolute -bottom-14 -right-8 grid size-72 place-items-center bg-accent/90 text-white [clip-path:polygon(50%_0,93%_25%,93%_75%,50%_100%,7%_75%,7%_25%)] lg:size-96" aria-hidden="true">
          <Icon name="flower" className="size-40 text-white dark:text-brand lg:size-52" />
        </span>
      </section>
      <div className="flex min-h-dvh flex-col px-4 py-6 sm:px-8 md:px-10 lg:px-16">
        <header className="mb-6 flex items-center justify-between md:justify-end"><Link to="/" className="md:hidden"><Logo /></Link><LanguageSwitch /></header>
        <main id="main" className="my-auto w-full self-center rounded-2xl border border-line bg-surface p-6 shadow-soft sm:p-9 lg:max-w-xl">
          <h2 className="mb-6 text-2xl font-semibold">{title}</h2>
          {children}
        </main>
      </div>
    </div>
  )
}
