import { createContext, useCallback, useContext, useMemo, useState, type ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import { errorKey } from '@/shared/lib'

type Tone = 'success' | 'error'
interface ToastItem {
  id: number
  text: string
  tone: Tone
}

const Ctx = createContext<(text: string, tone?: Tone) => void>(() => {})

/** Notifications announced via aria-live regions. */
export function ToastProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<ToastItem[]>([])
  const push = useCallback((text: string, tone: Tone = 'success') => {
    const id = Date.now() + Math.random()
    setItems((s) => [...s, { id, text, tone }])
    setTimeout(() => setItems((s) => s.filter((i) => i.id !== id)), 5000)
  }, [])
  const value = useMemo(() => push, [push])
  return (
    <Ctx.Provider value={value}>
      {children}
      <div className="pointer-events-none fixed bottom-20 right-4 z-50 flex flex-col gap-2 md:bottom-4">
        <div role="status" aria-live="polite">
          {items.filter((i) => i.tone === 'success').map((i) => (
            <p key={i.id} className="rounded-lg bg-ink px-4 py-3 text-surface shadow-lg">{i.text}</p>
          ))}
        </div>
        <div role="alert" aria-live="assertive">
          {items.filter((i) => i.tone === 'error').map((i) => (
            <p key={i.id} className="rounded-lg bg-danger px-4 py-3 text-white shadow-lg">{i.text}</p>
          ))}
        </div>
      </div>
    </Ctx.Provider>
  )
}

export function useToast() {
  return useContext(Ctx)
}

/** Translated success / error notifications for mutations. */
export function useNotify() {
  const push = useToast()
  const { t } = useTranslation()
  return useMemo(() => ({ ok: (key: string) => push(t(key)), fail: (e: unknown) => push(t(errorKey(e)), 'error') }), [push, t])
}
