import { forwardRef, useId, type InputHTMLAttributes, type ReactNode, type SelectHTMLAttributes, type TextareaHTMLAttributes } from 'react'
import { cn } from '@/shared/lib'

const control = 'field-control mt-1 block w-full min-h-11 rounded-[0.875rem] border border-line bg-surface px-3 py-2 text-ink placeholder:text-ink-soft/75 focus-visible:border-brand-strong aria-[invalid=true]:border-danger disabled:cursor-not-allowed disabled:bg-muted disabled:opacity-70 motion-reduce:transition-none'

interface Base {
  label: string
  error?: string | undefined
  hint?: string | undefined
}

function Wrap({ id, label, error, hint, children }: Base & { id: string; children: ReactNode }) {
  return (
    <div className="mb-4">
      <label htmlFor={id} className="block text-sm font-medium text-ink">
        {label}
      </label>
      {children}
      {hint && !error && (
        <p id={`${id}-hint`} className="mt-1 text-sm text-ink-soft">
          {hint}
        </p>
      )}
      {error && (
        <p id={`${id}-error`} className="mt-1 text-sm text-danger" role="alert">
          {error}
        </p>
      )}
    </div>
  )
}

function describedBy(id: string, error?: string, hint?: string) {
  return error ? `${id}-error` : hint ? `${id}-hint` : undefined
}

export const Input = forwardRef<HTMLInputElement, Base & InputHTMLAttributes<HTMLInputElement>>(function Input({ label, error, hint, className, ...rest }, ref) {
  const id = useId()
  return (
    <Wrap id={id} label={label} error={error} hint={hint}>
      <input ref={ref} id={id} className={cn(control, className)} aria-invalid={!!error} aria-describedby={describedBy(id, error, hint)} {...rest} />
    </Wrap>
  )
})

export const Select = forwardRef<HTMLSelectElement, Base & SelectHTMLAttributes<HTMLSelectElement>>(function Select({ label, error, hint, className, children, ...rest }, ref) {
  const id = useId()
  return (
    <Wrap id={id} label={label} error={error} hint={hint}>
      <select ref={ref} id={id} className={cn(control, className)} aria-invalid={!!error} aria-describedby={describedBy(id, error, hint)} {...rest}>
        {children}
      </select>
    </Wrap>
  )
})

export const Textarea = forwardRef<HTMLTextAreaElement, Base & TextareaHTMLAttributes<HTMLTextAreaElement>>(function Textarea({ label, error, hint, className, ...rest }, ref) {
  const id = useId()
  return (
    <Wrap id={id} label={label} error={error} hint={hint}>
      <textarea ref={ref} id={id} rows={5} className={cn(control, className)} aria-invalid={!!error} aria-describedby={describedBy(id, error, hint)} {...rest} />
    </Wrap>
  )
})
