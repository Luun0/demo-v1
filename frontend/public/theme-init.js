(() => {
  const key = 'ssp_theme'
  let stored = null
  try { stored = localStorage.getItem(key) } catch { /* Keep system theme when storage is unavailable. */ }
  const mode = stored === 'light' || stored === 'dark' || stored === 'system' ? stored : 'system'
  const systemDark = window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false
  const theme = mode === 'system' ? (systemDark ? 'dark' : 'light') : mode
  const root = document.documentElement
  root.dataset.theme = theme
  root.dataset.themeMode = mode
  if (mode === 'system') root.style.removeProperty('color-scheme')
  else root.style.colorScheme = mode
})()
