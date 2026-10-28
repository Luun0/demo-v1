import { useEffect, useRef, useState, type FormEvent } from 'react'
import { useTranslation } from 'react-i18next'
import { Button, Icon } from '@/shared/ui'
import { faqItems, findFaq, type FaqId } from './faq'
import { FaqMark } from './FaqMark'

type Message = { id: number; author: 'user' | 'guide'; text: string }
const suggestions: FaqId[] = ['registration', 'approval', 'diary', 'privacy']

export function SupportAssistant() {
  const { t } = useTranslation('supportChat')
  const [open, setOpen] = useState(false)
  const [draft, setDraft] = useState('')
  const [messages, setMessages] = useState<Message[]>([])
  const sequence = useRef(0)
  const input = useRef<HTMLInputElement>(null)
  const list = useRef<HTMLDivElement>(null)
  const launcher = useRef<HTMLButtonElement>(null)

  function closeChat() {
    setOpen(false)
    launcher.current?.focus()
  }

  useEffect(() => {
    if (open) input.current?.focus()
  }, [open])

  useEffect(() => {
    if (!open) return
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') closeChat()
    }
    window.addEventListener('keydown', closeOnEscape)
    return () => window.removeEventListener('keydown', closeOnEscape)
  }, [open])

  useEffect(() => {
    if (list.current) list.current.scrollTop = list.current.scrollHeight
  }, [messages, open])

  function ask(question: string) {
    const text = question.trim()
    if (!text) return
    const id = ++sequence.current
    const answerId = findFaq(text)
    const answer = answerId ? t(`faq.${answerId}.answer`) : t('fallback')
    setMessages((current) => [...current, { id, author: 'user', text }, { id: ++sequence.current, author: 'guide', text: answer }])
    setDraft('')
  }

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    ask(draft)
  }

  return (
    <div className="fixed bottom-5 right-5 z-[60] flex flex-col items-end gap-3 sm:bottom-6 sm:right-6">
      {open && (
        <section id="support-chat-panel" aria-labelledby="support-chat-title" className="flex h-[min(38rem,calc(100dvh-7rem))] w-[min(25rem,calc(100vw-2rem))] flex-col overflow-hidden rounded-2xl border border-line bg-surface shadow-panel page-enter">
          <header className="flex items-center gap-3 border-b border-line bg-brand px-4 py-3 text-on-brand">
            <span className="grid size-10 shrink-0 place-items-center rounded-xl bg-on-brand/15"><FaqMark /></span>
            <div className="min-w-0 flex-1">
              <h2 id="support-chat-title" className="font-semibold leading-tight">{t('title')}</h2>
              <p className="text-sm text-on-brand/80">{t('subtitle')}</p>
            </div>
            <button type="button" onClick={closeChat} aria-label={t('close')} className="grid size-10 place-items-center rounded-xl transition-colors hover:bg-on-brand/15 focus-visible:ring-2 focus-visible:ring-accent">
              <Icon name="close" />
            </button>
          </header>

          <div ref={list} role="log" aria-live="polite" aria-relevant="additions" className="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto bg-canvas p-4">
            <p className="max-w-[90%] self-start rounded-2xl rounded-tl-sm border border-line bg-surface px-3.5 py-3 text-sm text-ink shadow-soft">{t('welcome')}</p>
            {messages.map((message) => (
              <p key={message.id} className={`max-w-[90%] whitespace-pre-wrap rounded-2xl px-3.5 py-3 text-sm ${message.author === 'user' ? 'self-end rounded-tr-sm bg-brand text-on-brand' : 'self-start rounded-tl-sm border border-line bg-surface text-ink shadow-soft'}`}>
                {message.text}
              </p>
            ))}
          </div>

          <div className="border-t border-line bg-surface p-3">
            {messages.length === 0 && <p className="mb-2 px-1 text-xs font-semibold uppercase tracking-wide text-ink-soft">{t('suggestions')}</p>}
            {messages.length === 0 && <div className="mb-3 flex flex-wrap gap-2">{suggestions.map((id) => <button key={id} type="button" onClick={() => ask(t(`faq.${id}.question`))} className="rounded-full border border-line bg-muted px-3 py-1.5 text-left text-xs text-ink transition-colors hover:border-brand hover:bg-surface focus-visible:ring-2 focus-visible:ring-accent">{t(`faq.${id}.question`)}</button>)}</div>}
            <form onSubmit={submit} className="flex items-end gap-2">
              <label className="sr-only" htmlFor="support-chat-input">{t('inputLabel')}</label>
              <input ref={input} id="support-chat-input" value={draft} onChange={(event) => setDraft(event.target.value)} maxLength={300} placeholder={t('placeholder')} className="field-control min-h-11 min-w-0 flex-1 rounded-xl border border-line bg-canvas px-3 text-sm text-ink placeholder:text-ink-soft focus-visible:border-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/30" />
              <Button type="submit" aria-label={t('send')} disabled={!draft.trim()} className="size-11 min-w-11 shrink-0 px-0"><Icon name="arrow" /></Button>
            </form>
            <p className="mt-2 px-1 text-xs leading-snug text-ink-soft">{t('privacy')}</p>
          </div>
        </section>
      )}

      <button ref={launcher} type="button" aria-label={open ? t('close') : t('launcher')} aria-expanded={open} aria-controls={open ? 'support-chat-panel' : undefined} onClick={() => setOpen((value) => !value)} className="inline-flex h-14 min-w-14 items-center justify-center gap-2 rounded-full bg-brand px-4 font-medium text-on-brand shadow-lift transition-[transform,background-color,box-shadow] duration-200 ease-out-soft hover:-translate-y-0.5 hover:scale-[1.015] hover:bg-brand-strong hover:shadow-lift active:translate-y-px active:scale-[0.985] focus-visible:ring-2 focus-visible:ring-accent focus-visible:ring-offset-2 motion-reduce:transform-none motion-reduce:transition-none">
        <Icon name={open ? 'close' : 'chat'} />
        <span className="hidden sm:inline">{open ? t('close') : t('launcher')}</span>
      </button>
    </div>
  )
}

export { faqItems }
