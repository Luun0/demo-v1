import { useTranslation } from 'react-i18next'

export function PrototypeNotice({ className = '' }: { className?: string }) {
  const { t } = useTranslation()
  return (
    <aside role="note" className={`rounded-xl border border-warn/50 bg-warn/15 px-4 py-3 text-sm text-ink ${className}`}>
      <p className="font-semibold">{t('prototypeTitle')}</p>
      <p className="mt-1 leading-relaxed">{t('prototypeNotice')}</p>
    </aside>
  )
}
