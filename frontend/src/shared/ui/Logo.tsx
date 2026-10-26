export function Logo({ inverse = false, className = '' }: { inverse?: boolean; className?: string }) {
  return (
    <span className={`inline-flex items-center gap-3 font-semibold ${inverse ? 'text-white' : 'text-brand-strong'} ${className}`}>
      <svg
        role="img"
        aria-label="SEN Support mark"
        className="size-11 shrink-0 overflow-visible"
        viewBox="0 0 48 48"
        fill="none"
      >
        <path d="M8 33c1-9 5-15 12-17m8 0c7 2 11 8 12 17" stroke="currentColor" strokeWidth="4" strokeLinecap="round" />
        <circle cx="15" cy="12" r="5" fill="currentColor" />
        <circle cx="33" cy="12" r="5" fill="currentColor" />
        <path d="M8 33c4 7 9 10 16 10s12-3 16-10" stroke="currentColor" strokeWidth="4" strokeLinecap="round" />
        <circle cx="24" cy="26" r="5" fill="var(--color-accent)" />
      </svg>
      <span>SEN Support</span>
    </span>
  )
}
