import type { ButtonHTMLAttributes } from 'react'
import { cn } from '@/shared/lib'

type Variant = 'primary' | 'secondary' | 'danger' | 'ghost' | 'accent'

const styles: Record<Variant, string> = {
  primary: 'bg-brand text-on-brand hover:bg-brand-strong hover:shadow-lift',
  accent: 'bg-accent text-on-accent hover:brightness-95 hover:shadow-lift',
  secondary: 'border border-line bg-surface text-ink hover:border-brand/40 hover:bg-muted hover:shadow-soft',
  danger: 'bg-danger-solid text-on-brand hover:brightness-95 hover:shadow-soft',
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
      className={cn('inline-flex min-h-11 min-w-11 items-center justify-center gap-2 rounded-[0.875rem] px-4 font-medium transition-[transform,background-color,box-shadow,border-color] duration-200 ease-out-soft hover:-translate-y-0.5 hover:scale-[1.015] active:translate-y-px active:scale-[0.985] focus-visible:ring-2 focus-visible:ring-accent focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-55 disabled:shadow-none motion-reduce:transform-none motion-reduce:transition-none', styles[variant], className)}
      disabled={disabled || busy}
      aria-busy={busy || undefined}
      {...rest}
    >
      {children}
    </button>
  )
}
