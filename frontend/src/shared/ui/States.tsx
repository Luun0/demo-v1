import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import { errorKey } from '@/shared/lib'
import { Button } from './Button'

export function Loading({ label }: { label?: string }) {
  const { t } = useTranslation()
  return (
    <div role="status" aria-live="polite" className="space-y-3 py-4">
      <span className="sr-only">{label ?? t('loading')}</span>
      {[0, 1, 2].map((i) => (
        <div key={i} className="skeleton-shimmer h-16 rounded-2xl border border-line bg-muted motion-reduce:animate-none" />
      ))}
    </div>
  )
}

export function EmptyState({ children, action }: { children: ReactNode; action?: ReactNode }) {
  return (
    <div className="relative overflow-hidden rounded-3xl border border-dashed border-line bg-surface px-6 py-10 text-center text-ink-soft shadow-soft sm:px-10">
      <span className="pointer-events-none absolute -right-8 -top-10 size-36 rounded-full bg-brand/5" aria-hidden="true" />
      <span className="relative mx-auto mb-4 block size-3 rounded-full bg-accent before:absolute before:-inset-2 before:rounded-full before:border before:border-accent/30" aria-hidden="true" />
      <div className="relative mx-auto max-w-xl text-ink">{children}</div>
      {action && <div className="relative mt-5">{action}</div>}
    </div>
  )
}

export function ErrorState({ error, onRetry }: { error: unknown; onRetry?: () => void }) {
  const { t } = useTranslation()
  return (
    <div role="alert" className="rounded-2xl border border-danger/40 bg-danger/5 p-6 shadow-soft">
      <p className="text-ink">{t(errorKey(error))}</p>
      {onRetry && (
        <Button variant="secondary" className="mt-3" onClick={onRetry}>
          {t('retry')}
        </Button>
      )}
    </div>
  )
}

type BadgeTone = 'neutral' | 'brand' | 'warn' | 'danger' | 'success' | 'accent'

export function Badge({ children, tone = 'neutral' }: { children: ReactNode; tone?: BadgeTone }) {
  const tones: Record<BadgeTone, string> = {
    neutral: 'bg-muted text-ink', brand: 'bg-brand/15 text-brand-strong',
    warn: 'bg-warn/20 text-ink', danger: 'bg-danger/15 text-danger',
    success: 'bg-brand text-white', accent: 'bg-accent/15 text-ink',
  }
  return <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-sm font-medium ${tones[tone]}`}>{children}</span>
}

export function PageHeader({ title, subtitle, actions }: { title: string; subtitle?: string; actions?: ReactNode }) {
  return (
    <div className="page-enter mb-8 flex flex-wrap items-end justify-between gap-4 border-b border-line pb-5">
      <div>
        <h1 className="text-3xl font-semibold tracking-tight text-ink sm:text-4xl">{title}</h1>
        {subtitle && <p className="mt-2 max-w-2xl text-base text-ink-soft">{subtitle}</p>}
      </div>
      {actions}
    </div>
  )
}
