import i18next from 'i18next'
import { initReactI18next } from 'react-i18next'
import ru from './locales/ru.json'

export const LANGS = ['kz', 'ru', 'en'] as const
export type Lang = (typeof LANGS)[number]
const STORAGE_KEY = 'ssp_lang' // UI preference only, never credentials

/** BCP-47 tags for Intl (kz is the product code, kk is the ISO code). */
export const INTL_LOCALE: Record<Lang, string> = { kz: 'kk-KZ', ru: 'ru-RU', en: 'en-GB' }

const loaders: Record<Lang, () => Promise<{ default: Record<string, unknown> }>> = {
  ru: async () => ({ default: ru }),
  kz: () => import('./locales/kz.json'),
  en: () => import('./locales/en.json'),
}

export function isLang(v: unknown): v is Lang {
  return typeof v === 'string' && (LANGS as readonly string[]).includes(v)
}

export function initialLang(): Lang {
  const saved = typeof localStorage !== 'undefined' ? localStorage.getItem(STORAGE_KEY) : null
  return isLang(saved) ? saved : 'ru'
}

async function load(lang: Lang): Promise<void> {
  if (i18next.hasResourceBundle(lang, 'common')) return
  const mod = await loaders[lang]()
  for (const [ns, res] of Object.entries(mod.default)) i18next.addResourceBundle(lang, ns, res, true, true)
}

/** Switches the UI language, lazy-loading the locale bundle. */
export async function changeLanguage(lang: Lang): Promise<void> {
  await load(lang)
  await i18next.changeLanguage(lang)
  document.documentElement.lang = INTL_LOCALE[lang].slice(0, 2)
  localStorage.setItem(STORAGE_KEY, lang)
}

export async function initI18n(lang: Lang = initialLang()): Promise<typeof i18next> {
  await i18next.use(initReactI18next).init({
    lng: 'ru',
    fallbackLng: 'ru',
    ns: Object.keys(ru),
    defaultNS: 'common',
    resources: { ru },
    interpolation: { escapeValue: false }, // React escapes output
    returnNull: false,
  })
  if (lang !== 'ru') await changeLanguage(lang)
  return i18next
}

export { i18next }
