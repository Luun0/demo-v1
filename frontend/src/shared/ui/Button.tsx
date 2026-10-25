import type { ButtonHTMLAttributes } from 'react'
import { cn } from '@/shared/lib'

type Variant = 'primary' | 'secondary' | 'danger' | 'ghost' | 'accent'

const styles: Record<Variant, string> = {
  primary: 'bg-brand text-on-brand hover:bg-brand-strong',
  accent: 'rounded-full bg-accent text-ink hover:brightness-95',
  secondary: 'border border-line bg-surface text-ink hover:bg-muted',
  danger: 'bg-danger text-white hover:opacity-90 dark:text-brand',
  ghost: 'text-ink hover:bg-muted',
}

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant
  busy?: boolean
}

export function Button({ variant = 'primary', busy = false, className, disabled, children, type = 'button', ...rest }: ButtonProps) {
  return (
    <button
      type={type}
      className={cn('inline-flex min-h-11 min-w-11 items-center justify-center gap-2 rounded-xl px-4 font-medium transition-colors disabled:cursor-not-allowed disabled:opacity-60 motion-reduce:transition-none', styles[variant], className)}
      disabled={disabled || busy}
      aria-busy={busy || undefined}
      {...rest}
    >
      {children}
    </button>
  )
}
