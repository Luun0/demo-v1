import { useSyncExternalStore } from 'react'
import { useTranslation } from 'react-i18next'
import { Icon } from '@/shared/ui'
import { getThemeSnapshot, setThemeMode, subscribeTheme, type ThemeMode } from './model/theme'

const nextMode: Record<ThemeMode, ThemeMode> = { system: 'light', light: 'dark', dark: 'system' }

export function ThemeSwitch({ dark = false }: { dark?: boolean }) {
  const { t } = useTranslation()
  const snapshot = useSyncExternalStore(subscribeTheme, getThemeSnapshot, () => 'system|light')
  const mode = snapshot.split('|')[0] as ThemeMode
  const selected = nextMode[mode]
  const isDark = snapshot.endsWith('|dark')
  const label = t('common:theme.switchLabel', {
    current: t(`common:theme.${mode}`),
    next: t(`common:theme.${selected}`),
  })
  return (
    <button
      type="button"
      aria-label={label}
      aria-pressed={isDark}
      title={label}
      onClick={() => setThemeMode(selected)}
      className={`grid min-h-11 min-w-11 place-items-center rounded-xl border transition-[transform,background-color,border-color] duration-200 ease-out-soft hover:-translate-y-0.5 focus-visible:ring-2 focus-visible:ring-accent active:translate-y-px motion-reduce:transition-none motion-reduce:transform-none ${dark ? 'border-on-brand/30 text-on-brand hover:bg-on-brand/10' : 'border-line bg-surface text-brand-ink hover:bg-muted'}`}
    >
      <Icon name={isDark ? 'sun' : 'moon'} />
    </button>
  )
}
