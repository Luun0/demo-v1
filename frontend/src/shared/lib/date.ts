import { INTL_LOCALE, isLang } from '@/shared/i18n'

function locale(lang: string): string {
  return isLang(lang) ? INTL_LOCALE[lang] : 'ru-RU'
}

export function formatDate(iso: string, lang: string): string {
  const d = /^\d{4}-\d{2}-\d{2}$/.test(iso) ? new Date(`${iso}T00:00:00Z`) : new Date(iso)
  return new Intl.DateTimeFormat(locale(lang), { dateStyle: 'long', timeZone: /^\d{4}-\d{2}-\d{2}$/.test(iso) ? 'UTC' : undefined }).format(d)
}

export function formatDateTime(iso: string, lang: string): string {
  return new Intl.DateTimeFormat(locale(lang), { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(iso))
}

export function formatRelative(iso: string, lang: string, now = Date.now()): string {
  const diffMin = Math.round((new Date(iso).getTime() - now) / 60000)
  const rtf = new Intl.RelativeTimeFormat(locale(lang), { numeric: 'auto' })
  if (Math.abs(diffMin) < 60) return rtf.format(diffMin, 'minute')
  if (Math.abs(diffMin) < 60 * 24) return rtf.format(Math.round(diffMin / 60), 'hour')
  return rtf.format(Math.round(diffMin / 1440), 'day')
}

/** Value for <input type="datetime-local"> in local time. */
export function toLocalInput(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

export function localInputToISO(v: string): string {
  return new Date(v).toISOString().replace(/\.\d{3}Z$/, 'Z')
}
