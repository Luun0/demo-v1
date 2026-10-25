import { Icon } from './Icon'

export function Logo({ inverse = false, className = '' }: { inverse?: boolean; className?: string }) {
  return (
    <span className={`inline-flex items-center gap-3 font-semibold ${inverse ? 'text-white' : 'text-ink'} ${className}`}>
      <span className="grid size-11 shrink-0 place-items-center bg-accent text-white [clip-path:polygon(50%_0,93%_25%,93%_75%,50%_100%,7%_75%,7%_25%)]">
        <Icon name="flower" className="size-6 text-white dark:text-brand" />
      </span>
      <span>SEN Support</span>
    </span>
  )
}
