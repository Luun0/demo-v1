import { useEffect, useRef, useState, type CSSProperties, type ReactNode } from 'react'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { LanguageSwitch } from '@/features/language-switch'
import { Icon, Logo } from '@/shared/ui'

const features = ['f1', 'f2', 'f3'] as const
const roles = ['parent', 'specialist', 'admin'] as const
const steps = ['step1', 'step2', 'step3'] as const
type ProgressStyle = CSSProperties & { '--progress': number }

function IntroLoader() {
  const [visible, setVisible] = useState(false)
  const [count, setCount] = useState(0)

  useEffect(() => {
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return
    setVisible(true)
    const start = performance.now()
    let frame = 0
    const tick = (now: number) => {
      setCount(Math.min(100, Math.round(((now - start) / 820) * 100)))
      if (now - start < 820) frame = requestAnimationFrame(tick)
    }
    frame = requestAnimationFrame(tick)
    const done = window.setTimeout(() => setVisible(false), 1260)
    return () => { cancelAnimationFrame(frame); clearTimeout(done) }
  }, [])

  if (!visible) return null
  return (
    <div className="intro-loader pointer-events-none fixed inset-0 z-50 grid place-items-center bg-ink-panel text-panel-ink" aria-hidden="true">
      <div className="w-[min(76vw,44rem)]">
        <div className="mb-5 flex items-end justify-between border-b border-panel-ink/25 pb-3">
          <span className="text-sm font-semibold tracking-[0.12em]">SEN SUPPORT</span>
          <span className="font-mono text-sm tabular-nums text-sun">{String(count).padStart(3, '0')}</span>
        </div>
        <div className="h-1 origin-left scale-x-[var(--progress)] bg-accent" style={{ '--progress': count / 100 } as ProgressStyle} />
      </div>
    </div>
  )
}

function Reveal({ children, className = '' }: { children: ReactNode; className?: string }) {
  const ref = useRef<HTMLDivElement>(null)
  const [shown, setShown] = useState(false)
  useEffect(() => {
    const node = ref.current
    if (!node || !('IntersectionObserver' in window)) { setShown(true); return }
    const observer = new IntersectionObserver(([entry]) => {
      if (entry?.isIntersecting) { setShown(true); observer.disconnect() }
    }, { rootMargin: '0px 0px -10% 0px' })
    observer.observe(node)
    return () => observer.disconnect()
  }, [])
  return <div ref={ref} className={`scroll-reveal ${shown ? 'is-revealed' : ''} ${className}`}>{children}</div>
}

function CountUp({ value, label }: { value: number; label: string }) {
  const ref = useRef<HTMLSpanElement>(null)
  const [count, setCount] = useState(0)
  useEffect(() => {
    const node = ref.current
    if (!node) return
    let frame = 0
    const show = () => {
      if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) { setCount(value); return }
      const start = performance.now()
      const tick = (now: number) => {
        const progress = Math.min(1, (now - start) / 620)
        setCount(Math.round(value * (1 - (1 - progress) ** 3)))
        if (progress < 1) frame = requestAnimationFrame(tick)
      }
      frame = requestAnimationFrame(tick)
    }
    const observer = new IntersectionObserver(([entry]) => {
      if (entry?.isIntersecting) { show(); observer.disconnect() }
    }, { threshold: 0.5 })
    observer.observe(node)
    return () => { observer.disconnect(); cancelAnimationFrame(frame) }
  }, [value])
  return <span className="block"><span ref={ref} aria-hidden="true" className="count-pop tabular-nums">{count}</span><span className="sr-only">{value} {label}</span></span>
}

function WordReveal({ text }: { text: string }) {
  const words = text.split(' ')
  return <p className="mt-6 max-w-2xl text-lg leading-relaxed text-ink-soft md:text-xl">
    <span className="sr-only">{text}</span><span aria-hidden="true">{words.map((word, index) => <span key={`${word}-${index}`} className="inline-block animate-[page-enter_500ms_var(--ease-out-soft)_both]" style={{ animationDelay: `${Math.min(index, 4) * 45}ms` }}>{word}{index < words.length - 1 ? '\u00a0' : ''}</span>)}</span>
  </p>
}

export function LandingPage() {
  const { t } = useTranslation('landing')
  const title = t('title')
  const titleLines = title.split(/\s+—\s+/)
  return (
    <div className="min-h-dvh overflow-hidden bg-canvas text-ink">
      <IntroLoader />
      <a href="#main" className="sr-only focus:not-sr-only focus:fixed focus:left-2 focus:top-2 focus:z-50 focus:rounded focus:bg-surface focus:p-3">{t('common:skipToContent')}</a>
      <header className="mx-auto flex max-w-[90rem] items-center justify-between gap-4 px-5 py-5 md:px-10">
        <Logo />
        <nav aria-label={t('common:nav:main')} className="flex items-center gap-4">
          <LanguageSwitch />
          <Link to="/login" className="rounded-lg px-3 py-2 font-medium text-brand-ink transition-[background-color] hover:bg-brand/10 focus-visible:ring-2 focus-visible:ring-accent motion-reduce:transition-none">{t('ctaLogin')}</Link>
        </nav>
      </header>

      <main id="main" className="mx-auto max-w-[90rem] px-5 pb-16 md:px-10">
        <section className="relative grid min-h-[min(76vh,54rem)] items-center gap-9 border-b chapter-rule py-10 md:grid-cols-[1.05fr_.95fr] md:py-14">
          <span aria-hidden="true" className="landing-watermark pointer-events-none absolute bottom-0 right-0 translate-x-[7%] translate-y-[24%] select-none" />
          <div className="relative z-10">
            <p className="inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.14em] text-brand-ink"><span className="size-2 rounded-full bg-accent" />{t('eyebrow')}</p>
            <h1 aria-label={title} className="mt-7 max-w-3xl text-[clamp(2.65rem,5vw,4.8rem)] font-semibold leading-[1.02] tracking-[-0.06em] text-ink"><span className="sr-only">{title}</span><span aria-hidden="true">{titleLines.map((line, index) => <span key={line} className="block animate-[page-enter_600ms_var(--ease-out-soft)_both]" style={{ animationDelay: `${index * 90}ms` }}>{line}</span>)}</span></h1>
            <WordReveal text={t('text')} />
            <div className="mt-8 flex flex-wrap items-center gap-4">
              <Link to="/register" className="group inline-flex min-h-12 items-center gap-3 rounded-lg bg-accent px-5 font-semibold text-on-accent shadow-soft transition-[transform,box-shadow] duration-200 ease-out-soft hover:-translate-y-0.5 hover:scale-[1.015] hover:shadow-lift active:translate-y-px active:scale-[0.985] focus-visible:ring-2 focus-visible:ring-accent focus-visible:ring-offset-2 motion-reduce:transform-none motion-reduce:transition-none">{t('ctaRegister')}<Icon name="arrow" className="size-5 transition-transform duration-200 group-hover:translate-x-1 motion-reduce:transition-none" /></Link>
              <span className="text-sm text-ink-soft">{t('common:demoNotice')}</span>
            </div>
          </div>
          <div className="relative z-10 md:pl-5">
            <div className="relative overflow-hidden rounded-[1.5rem] bg-ink-panel p-6 text-panel-ink shadow-panel sm:p-9">
              <div className="flex items-start justify-between gap-4 border-b border-panel-ink/20 pb-5">
                <div><p className="text-xs font-semibold uppercase tracking-[0.14em] text-lagoon">{t('previewLabel')}</p><h2 className="mt-2 text-2xl font-semibold tracking-tight">{t('workspaceTitle')}</h2></div>
                <span className="grid size-11 place-items-center rounded-full bg-sun text-ink-panel"><Icon name="people" className="size-6" /></span>
              </div>
              <div className="mt-7 grid grid-cols-[2.5rem_1fr] gap-x-4">
                <div className="flex flex-col items-center"><span className="grid size-9 place-items-center rounded-full border border-lagoon/60 text-sm font-semibold text-lagoon">01</span><span className="my-2 min-h-12 w-px flex-1 bg-panel-ink/30" /><span className="grid size-9 place-items-center rounded-full bg-accent text-sm font-semibold text-on-accent">02</span><span className="my-2 min-h-12 w-px flex-1 bg-panel-ink/30" /><span className="grid size-9 place-items-center rounded-full border border-sun/60 text-sm font-semibold text-sun">03</span></div>
                <div className="space-y-4">
                  {roles.map((role, index) => <div key={role} className={`min-h-[5.3rem] border-b border-panel-ink/20 pb-4 ${index === 2 ? 'border-0' : ''}`}><p className="text-xs text-panel-ink/65">0{index + 1}</p><h3 className="mt-1 font-semibold">{t(`${role}Title`)}</h3><p className="text-sm text-panel-ink/75">{index === 0 ? t('parentText') : index === 1 ? t('specialistText') : t('adminText')}</p></div>)}
                </div>
              </div>
              <div className="mt-3 flex items-center justify-between border-t border-panel-ink/20 pt-4 text-xs text-panel-ink/70"><span>{t('footerNote')}</span><span className="size-2 rounded-full bg-sun" aria-hidden="true" /></div>
            </div>
            <span aria-hidden="true" className="absolute -bottom-5 -left-5 -z-10 size-24 border border-brand/30" />
          </div>
        </section>

        <Reveal className="py-14 md:py-20">
          <div className="mb-8 flex flex-wrap items-end justify-between gap-4">
            <h2 className="max-w-xl text-3xl font-semibold tracking-[-0.04em] md:text-5xl">{t('workspaceTitle')}</h2>
            <p className="max-w-sm text-ink-soft">{t('whyText')}</p>
          </div>
          <div className="grid gap-3 md:grid-cols-3">
            {features.map((key, index) => <article key={key} className={`surface-lift group relative min-h-64 overflow-hidden border border-line p-6 md:p-7 ${index === 0 ? 'bg-ink-panel text-panel-ink' : 'bg-surface text-ink'}`}>
              <div className="flex items-start justify-between"><span className={`text-sm font-semibold ${index === 0 ? 'text-lagoon' : 'text-brand-ink'}`}>0{index + 1}</span><span className={`grid size-9 place-items-center rounded-full transition-transform duration-200 group-hover:-translate-y-0.5 group-hover:translate-x-0.5 motion-reduce:transition-none ${index === 0 ? 'bg-panel-ink/10 text-panel-ink' : 'bg-muted text-brand-ink'}`}><Icon name="arrow" className="size-4 -rotate-45" /></span></div>
              <h3 className="mt-14 text-2xl font-semibold tracking-tight">{t(`${key}t`)}</h3>
              <p className={`mt-3 max-w-sm ${index === 0 ? 'text-panel-ink/75' : 'text-ink-soft'}`}>{t(`${key}d`)}</p>
            </article>)}
          </div>
        </Reveal>

        <Reveal>
          <section aria-label={t('statsTitle')} className="grid gap-8 rounded-[1.5rem] bg-ink-panel px-6 py-9 text-panel-ink shadow-panel md:grid-cols-[.8fr_2fr] md:items-center md:px-10 md:py-11">
            <div><p className="text-sm font-semibold uppercase tracking-[0.14em] text-sun">{t('statsEyebrow')}</p><h2 className="mt-3 text-2xl font-semibold tracking-tight">{t('statsTitle')}</h2></div>
            <div className="grid grid-cols-3 gap-4">
              {[{ value: 2, key: 'statOne' }, { value: 1, key: 'statTwo' }, { value: 3, key: 'statThree' }].map(({ value, key }) => <div key={key} className="border-l border-panel-ink/25 pl-4"><p className="text-4xl font-semibold tracking-[-0.06em] text-lagoon sm:text-5xl"><CountUp value={value} label={t(`${key}Label`)} /></p><p aria-hidden="true" className="mt-2 text-sm leading-snug text-panel-ink/75">{t(`${key}Label`)}</p></div>)}
            </div>
          </section>
        </Reveal>

        <section className="grid gap-10 py-16 md:grid-cols-[.8fr_1.2fr] md:py-24">
          <div className="max-w-md"><p className="text-sm font-semibold uppercase tracking-[0.14em] text-accent-ink">{t('howTitle')}</p><h2 className="mt-3 text-3xl font-semibold tracking-[-0.04em] md:text-5xl">{t('whyTitle')}</h2><p className="mt-5 text-ink-soft">{t('whyText')}</p></div>
          <ol className="border-t border-line">
            {steps.map((step, index) => <li key={step} className="group flex items-start gap-5 border-b border-line px-3 py-6 transition-[background-color] duration-200 hover:bg-surface focus-within:bg-surface motion-reduce:transition-none"><span className="pt-1 text-sm font-semibold text-accent-ink">0{index + 1}</span><p className="flex-1 text-lg font-medium leading-snug md:text-xl">{t(step)}</p><Icon name="arrow" className="mt-1 size-5 -rotate-45 text-brand-ink transition-transform duration-200 group-hover:-translate-y-0.5 group-hover:translate-x-0.5 motion-reduce:transition-none" /></li>)}
          </ol>
        </section>

        <section className="grid gap-10 border-y border-line py-14 md:grid-cols-2 md:py-20">
          <div><p className="text-sm font-semibold uppercase tracking-[0.14em] text-brand-ink">{t('privacyTitle')}</p><h2 className="mt-3 text-3xl font-semibold tracking-[-0.04em] md:text-4xl">{t('privacyTitle')}</h2></div>
          <div><p className="text-lg leading-relaxed text-ink-soft">{t('privacyText')}</p><div className="mt-6 flex items-center gap-3 border-t border-line pt-4"><span className="grid size-10 place-items-center rounded-full bg-brand text-on-brand"><Icon name="lock" /></span><p className="font-semibold">{t('adminText')}</p></div></div>
        </section>

        <section className="my-14 grid gap-6 bg-accent px-6 py-10 text-on-accent md:my-20 md:grid-cols-[1fr_auto] md:items-end md:px-10 md:py-12">
          <div><h2 className="max-w-2xl text-3xl font-semibold tracking-[-0.04em] md:text-5xl">{t('finalTitle')}</h2><p className="mt-3 max-w-xl text-on-accent">{t('finalText')}</p></div>
          <Link to="/register" className="group inline-flex min-h-12 items-center justify-center gap-3 rounded-lg bg-on-accent px-5 font-semibold text-accent transition-[transform,box-shadow] duration-200 ease-out-soft hover:-translate-y-0.5 hover:shadow-lift active:translate-y-px focus-visible:ring-2 focus-visible:ring-on-accent focus-visible:ring-offset-2 motion-reduce:transform-none motion-reduce:transition-none">{t('ctaRegister')}<Icon name="arrow" className="size-5 transition-transform duration-200 group-hover:translate-x-1 motion-reduce:transition-none" /></Link>
        </section>
      </main>
      <footer className="border-t border-line px-5 py-8 md:px-10"><div className="mx-auto flex max-w-[90rem] flex-wrap items-end justify-between gap-5"><Logo /><span className="text-sm text-ink-soft">{t('footerNote')}</span></div></footer>
    </div>
  )
}
