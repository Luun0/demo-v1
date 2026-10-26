export const themeModes = ['system', 'light', 'dark'] as const
export type ThemeMode = (typeof themeModes)[number]

const key = 'ssp_theme'

export function getThemeMode(): ThemeMode {
  try {
    const value = localStorage.getItem(key)
    return themeModes.includes(value as ThemeMode) ? (value as ThemeMode) : 'system'
  } catch {
    return 'system'
  }
}

export function getThemeSnapshot() {
  const mode = getThemeMode()
  const theme = document.documentElement.dataset.theme ??
    (mode === 'system' && window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'dark' : mode === 'dark' ? 'dark' : 'light')
  return `${mode}|${theme}`
}

export function applyTheme(mode: ThemeMode) {
  const root = document.documentElement
  const systemDark = window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false
  root.dataset.themeMode = mode
  root.dataset.theme = mode === 'system' ? (systemDark ? 'dark' : 'light') : mode
  if (mode === 'system') root.style.removeProperty('color-scheme')
  else root.style.colorScheme = mode
}

const subscribers = new Set<() => void>()

export function subscribeTheme(listener: () => void) {
  const media = window.matchMedia?.('(prefers-color-scheme: dark)')
  const onSystemChange = () => {
    if (getThemeMode() === 'system') applyTheme('system')
    listener()
  }
  const onStorage = (event: StorageEvent) => {
    if (event.key === key || event.key === null) {
      const mode = getThemeMode()
      applyTheme(mode)
      listener()
    }
  }
  subscribers.add(listener)
  media?.addEventListener('change', onSystemChange)
  window.addEventListener('storage', onStorage)
  return () => {
    subscribers.delete(listener)
    media?.removeEventListener('change', onSystemChange)
    window.removeEventListener('storage', onStorage)
  }
}

export function setThemeMode(mode: ThemeMode) {
  try { localStorage.setItem(key, mode) } catch { /* Keep this tab usable when storage is blocked. */ }
  applyTheme(mode)
  subscribers.forEach((listener) => listener())
}
