import { useId } from 'react'
import { useTranslation } from 'react-i18next'
import { changeLanguage, isLang, LANGS } from '@/shared/i18n'

const labels = { kz: 'Қазақша', ru: 'Русский', en: 'English' } as const

export function LanguageSwitch() {
  const { t, i18n } = useTranslation()
  const id = useId()
  return (
    <div className="flex items-center gap-2">
      <label htmlFor={id} className="text-sm text-ink-soft">{t('language')}</label>
      <select id={id} value={i18n.language} onChange={(e) => isLang(e.target.value) && void changeLanguage(e.target.value)} className="min-h-11 rounded-lg border border-line bg-surface px-2 text-ink">
        {LANGS.map((l) => (
          <option key={l} value={l} lang={l === 'kz' ? 'kk' : l}>{labels[l]}</option>
        ))}
      </select>
    </div>
  )
}
