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
        <div key={i} className="h-16 animate-pulse rounded-lg bg-muted motion-reduce:animate-none" />
      ))}
    </div>
  )
}

export function EmptyState({ children, action }: { children: ReactNode; action?: ReactNode }) {
  return (
    <div className="rounded-xl border border-dashed border-line p-8 text-center text-ink-soft">
      <p>{children}</p>
      {action && <div className="mt-4">{action}</div>}
    </div>
  )
}

export function ErrorState({ error, onRetry }: { error: unknown; onRetry?: () => void }) {
  const { t } = useTranslation()
  return (
    <div role="alert" className="rounded-xl border border-danger/40 bg-danger/5 p-6">
      <p className="text-ink">{t(errorKey(error))}</p>
      {onRetry && (
        <Button variant="secondary" className="mt-3" onClick={onRetry}>
          {t('retry')}
        </Button>
      )}
    </div>
  )
}

export function Badge({ children, tone = 'neutral' }: { children: ReactNode; tone?: 'neutral' | 'brand' | 'warn' | 'danger' }) {
  const tones = { neutral: 'bg-muted text-ink', brand: 'bg-brand/15 text-brand-strong', warn: 'bg-warn/20 text-ink', danger: 'bg-danger/15 text-danger' }
  return <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-sm font-medium ${tones[tone]}`}>{children}</span>
}

export function PageHeader({ title, subtitle, actions }: { title: string; subtitle?: string; actions?: ReactNode }) {
  return (
    <div className="mb-6 flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 className="text-2xl font-semibold text-ink">{title}</h1>
        {subtitle && <p className="mt-1 text-ink-soft">{subtitle}</p>}
      </div>
      {actions}
    </div>
  )
}
