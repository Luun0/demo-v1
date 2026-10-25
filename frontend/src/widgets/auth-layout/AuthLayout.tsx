import type { ReactNode } from 'react'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { LanguageSwitch } from '@/features/language-switch'

export function AuthLayout({ title, children }: { title: string; children: ReactNode }) {
  const { t } = useTranslation()
  return (
    <div className="min-h-dvh bg-canvas px-4 py-8 text-ink">
      <header className="mx-auto mb-8 flex max-w-md items-center justify-between">
        <Link to="/" className="text-lg font-bold text-brand-strong">{t('brand')}</Link>
        <LanguageSwitch />
      </header>
      <main id="main" className="mx-auto max-w-md rounded-xl border border-line bg-surface p-6 shadow-sm">
        <h1 className="mb-6 text-2xl font-semibold">{title}</h1>
        {children}
      </main>
    </div>
  )
}
