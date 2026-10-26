import { useId } from 'react'

export function Logo({ inverse = false, className = '' }: { inverse?: boolean; className?: string }) {
  const id = useId()
  return (
    <span className={`logo inline-flex items-center gap-3 font-semibold ${inverse ? 'text-on-brand' : 'text-brand-ink'} ${className}`}>
      <svg
        role="img"
        aria-labelledby={`${id}-title ${id}-desc`}
        className="logo-mark size-11 shrink-0 overflow-visible"
        viewBox="0 0 48 48"
        fill="none"
      >
        <title id={`${id}-title`}>SEN Support mark</title>
        <desc id={`${id}-desc`}>Two caregivers gently cradle a child.</desc>
        <path d="M8 33c0-9 3-15 10-19m12 0c7 4 10 10 10 19" stroke="currentColor" strokeWidth="4" strokeLinecap="round" />
        <circle cx="15" cy="12" r="5" fill="currentColor" />
        <circle cx="33" cy="12" r="5" fill="currentColor" />
        <path d="M8 32c4 8 9 12 16 12s12-4 16-12" stroke="currentColor" strokeWidth="4" strokeLinecap="round" />
        <circle className="logo-child" cx="24" cy="26" r="5.5" fill="var(--color-accent)" />
      </svg>
      <span>SEN Support</span>
    </span>
  )
}
