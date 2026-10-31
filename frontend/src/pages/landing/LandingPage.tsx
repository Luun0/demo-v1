import { useEffect, useRef, useState, type CSSProperties, type PointerEvent as ReactPointerEvent, type ReactNode } from 'react'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { LanguageSwitch } from '@/features/language-switch'
import { ThemeSwitch } from '@/features/theme-switch'
import { Icon, Logo, PrototypeNotice } from '@/shared/ui'

const features = ['f1', 'f2', 'f3'] as const
const roles = ['parent', 'specialist', 'admin'] as const
const roleIcons = { parent: 'calendar', specialist: 'diary', admin: 'shield' } as const
const roleCardLayouts = {
  parent: 'lg:col-span-4 lg:min-h-64',
  specialist: 'lg:col-start-5 lg:col-span-3 lg:mt-12 lg:min-h-72',
  admin: 'lg:col-start-8 lg:col-span-5 lg:mt-5 lg:min-h-60',
} as const
const steps = ['step1', 'step2', 'step3'] as const
type ProgressStyle = CSSProperties & { '--progress': number }

function InterfacePreview() {
  const { t } = useTranslation('landing')
  return (
    <article className="rounded-2xl border border-line bg-surface p-5 text-ink shadow-panel sm:p-7" aria-label={t('preview.eyebrow')}>
      <header className="flex flex-wrap items-start justify-between gap-4 border-b border-line pb-5">
        <div>
          <p className="text-sm font-semibold text-brand-ink">{t('preview.eyebrow')}</p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight">{t('preview.child')}</h2>
        </div>
        <span className="rounded-full bg-muted px-4 py-2 text-sm font-medium text-ink-soft">{t('preview.visibility')}</span>
      </header>
      <section className="grid gap-3 border-b border-line py-5 sm:grid-cols-[8rem_1fr]" aria-labelledby="preview-session">
        <h3 id="preview-session" className="font-semibold text-ink-soft">{t('preview.sessionHeading')}</h3>
        <div className="flex items-start gap-3">
          <span aria-hidden="true" className="mt-2 size-2.5 shrink-0 rounded-full bg-accent" />
          <div><p className="text-lg font-semibold">{t('preview.sessionTitle')}</p><p className="mt-1 text-ink-soft">{t('preview.sessionTime')}</p></div>
          <span className="ml-auto rounded-full bg-muted px-3 py-1.5 text-sm font-medium text-brand-ink">{t('preview.sessionStatus')}</span>
        </div>
      </section>
      <section className="grid gap-3 pt-5 sm:grid-cols-[8rem_1fr]" aria-labelledby="preview-observation">
        <h3 id="preview-observation" className="font-semibold text-ink-soft">{t('preview.observationHeading')}</h3>
        <p className="text-lg leading-relaxed">{t('preview.observation')}</p>
      </section>
    </article>
  )
}

function trackGlow(event: ReactPointerEvent<HTMLElement>) {
  if (event.pointerType !== 'mouse') return
  const bounds = event.currentTarget.getBoundingClientRect()
  event.currentTarget.style.setProperty('--mx', `${event.clientX - bounds.left}px`)
  event.currentTarget.style.setProperty('--my', `${event.clientY - bounds.top}px`)
}

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
    <div className="landing-page min-h-dvh overflow-hidden bg-canvas text-ink">
      <IntroLoader />
      <a href="#main" className="sr-only focus:not-sr-only focus:fixed focus:left-2 focus:top-2 focus:z-50 focus:rounded focus:bg-surface focus:p-3">{t('common:skipToContent')}</a>
      <header className="landing-header sticky top-0 z-30 border-b border-line/60">
        <div className="mx-auto flex max-w-[90rem] items-center justify-between gap-4 px-5 py-4 md:px-10">
          <Logo />
          <nav aria-label={t('common:nav:main')} className="hidden items-center gap-5 md:flex">
            <a href="#for-people" className="font-medium text-ink-soft hover:text-brand-ink">{t('navPeople')}</a>
            <a href="#how-it-works" className="font-medium text-ink-soft hover:text-brand-ink">{t('navWorkflow')}</a>
            <a href="#privacy" className="font-medium text-ink-soft hover:text-brand-ink">{t('navPrivacy')}</a>
            <ThemeSwitch />
            <LanguageSwitch />
            <Link to="/login" className="rounded-lg px-3 py-2 font-medium text-brand-ink transition-[background-color] hover:bg-brand/10 focus-visible:ring-2 focus-visible:ring-accent motion-reduce:transition-none">{t('ctaLogin')}</Link>
            <Link to="/register" className="inline-flex min-h-11 items-center rounded-lg bg-accent px-4 font-semibold text-on-accent shadow-soft transition-[transform,box-shadow] duration-200 ease-out-soft hover:-translate-y-0.5 hover:shadow-lift focus-visible:ring-2 focus-visible:ring-accent motion-reduce:transition-none motion-reduce:transform-none">{t('ctaRegister')}</Link>
          </nav>
          <details className="group relative md:hidden">
            <summary aria-label={t('common:menu')} className="grid min-h-11 min-w-11 cursor-pointer list-none place-items-center rounded-xl border border-line bg-surface text-brand-ink focus-visible:ring-2 focus-visible:ring-accent"><Icon name="menu" /></summary>
            <div className="absolute right-0 top-14 z-40 grid min-w-64 gap-4 rounded-2xl border border-line bg-surface p-4 shadow-panel">
              <a href="#for-people" className="min-h-11 py-2 font-medium text-ink">{t('navPeople')}</a>
              <a href="#how-it-works" className="min-h-11 py-2 font-medium text-ink">{t('navWorkflow')}</a>
              <a href="#privacy" className="min-h-11 py-2 font-medium text-ink">{t('navPrivacy')}</a>
              <div className="flex items-center justify-between gap-3 border-t border-line pt-3"><ThemeSwitch /><LanguageSwitch /></div>
              <Link to="/login" className="min-h-11 py-2 font-medium text-brand-ink">{t('ctaLogin')}</Link>
              <Link to="/register" className="inline-flex min-h-11 items-center justify-center rounded-lg bg-accent px-4 font-semibold text-on-accent">{t('ctaRegister')}</Link>
            </div>
          </details>
        </div>
      </header>

      <div className="mx-auto max-w-[90rem] px-5 md:px-10">
        <PrototypeNotice />
      </div>

      <main id="main" className="mx-auto max-w-[90rem] px-5 pb-16 md:px-10">
        <section className="relative grid min-h-[min(76vh,54rem)] items-center gap-9 border-b chapter-rule py-10 md:grid-cols-[1.05fr_.95fr] md:py-14">
          <span aria-hidden="true" className="aurora"><span className="aurora-noise block size-full" /></span>
          <span aria-hidden="true" className="landing-watermark pointer-events-none absolute bottom-0 right-0 translate-x-[7%] translate-y-[24%] select-none" />
          <div className="relative z-10">
            <p className="inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.14em] text-brand-ink"><span className="size-2 rounded-full bg-accent" />{t('eyebrow')}</p>
            <h1 aria-label={title} className="relative mt-7 max-w-3xl font-display text-[clamp(3.2rem,6.3vw,5.6rem)] font-normal leading-[.98] tracking-[-0.04em] text-ink"><span className="sr-only">{title}</span><span aria-hidden="true">{titleLines.map((line, index) => <span key={line} className="block animate-[page-enter_600ms_var(--ease-out-soft)_both]" style={{ animationDelay: `${index * 90}ms` }}>{line}</span>)}</span></h1>
            <WordReveal text={t('text')} />
            <div className="mt-8 flex flex-wrap items-center gap-4">
              <Link to="/register" className="brand-action group inline-flex min-h-12 items-center gap-3 rounded-lg px-5 font-semibold shadow-soft transition-[transform,box-shadow] duration-200 ease-out-soft hover:-translate-y-0.5 hover:scale-[1.015] hover:shadow-lift active:translate-y-px active:scale-[0.985] focus-visible:ring-2 focus-visible:ring-accent focus-visible:ring-offset-2 motion-reduce:transform-none motion-reduce:transition-none">{t('ctaRegister')}<Icon name="arrow" className="size-5 transition-transform duration-200 group-hover:translate-x-1 motion-reduce:transition-none" /></Link>
              <a href="#example" className="inline-flex min-h-12 items-center rounded-lg border border-line bg-surface px-4 font-semibold text-brand-ink transition-[transform,box-shadow] duration-200 ease-out-soft hover:-translate-y-0.5 hover:shadow-soft focus-visible:ring-2 focus-visible:ring-accent motion-reduce:transform-none motion-reduce:transition-none">{t('ctaExample')}</a>
              <span className="text-sm text-ink-soft">{t('common:demoNotice')}</span>
            </div>
          </div>
          <figure className="relative z-10 aspect-[1.55] overflow-hidden rounded-[1.75rem] border border-line bg-muted shadow-soft md:ml-5 md:self-center">
            <img src="/landing-care-collage.jpg" alt="" aria-hidden="true" fetchPriority="high" className="absolute inset-0 size-full object-cover object-[72%_center]" />
            <span aria-hidden="true" className="absolute inset-0 bg-gradient-to-t from-ink-panel/35 via-transparent to-transparent dark:bg-ink-panel/45 dark:from-ink-panel/65" />
            <figcaption className="sr-only">{t('previewLabel')}</figcaption>
          </figure>
        </section>

        <Reveal className="py-14 md:py-20">
          <section id="for-people" aria-labelledby="for-people-title">
            <div className="mb-8 max-w-2xl"><p className="text-sm font-semibold uppercase tracking-[0.14em] text-accent-ink">{t('rolesEyebrow')}</p><h2 id="for-people-title" className="mt-3 text-3xl font-semibold tracking-[-0.04em] md:text-5xl">{t('rolesTitle')}</h2></div>
            <div className="relative grid gap-4 md:grid-cols-2 lg:grid-cols-12 lg:items-start">
              <svg aria-hidden="true" viewBox="0 0 1000 240" preserveAspectRatio="none" className="pointer-events-none absolute inset-x-0 top-0 hidden h-60 w-full lg:block">
                <path d="M36 45 C188 8 252 182 405 143 S654 36 786 106 S921 171 974 115" fill="none" stroke="var(--color-accent)" strokeDasharray="3 9" strokeLinecap="round" strokeWidth="2" opacity=".58" />
              </svg>
              {roles.map((role, index) => {
                return (
                  <article key={role} className={`surface-lift landing-role-cutout landing-role-cutout--${role} relative z-10 min-h-52 border border-line bg-surface p-6 text-ink md:p-7 ${roleCardLayouts[role]}`}>
                    <div className="flex items-start justify-between">
                      <p aria-hidden="true" className="font-display text-5xl leading-none text-accent-ink/65">0{index + 1}</p>
                      {role === 'parent' ? <svg aria-hidden="true" viewBox="0 0 56 44" className="size-14"><path d="M5 27 Q28 3 51 27" fill="none" stroke="var(--color-accent)" strokeWidth="2.5" /><g fill="var(--color-brand-strong)" stroke="var(--color-brand-strong)" strokeLinecap="round" strokeWidth="2.5"><circle cx="14" cy="15" r="4" /><path d="M14 21v14m-5-8 5-4 5 4" fill="none" /><circle cx="28" cy="18" r="3" /><path d="M28 23v12" fill="none" /><circle cx="42" cy="15" r="4" /><path d="M42 21v14m-5-8 5-4 5 4" fill="none" /></g></svg> : <span className="grid size-10 place-items-center rounded-full bg-muted text-brand-ink"><Icon name={roleIcons[role]} /></span>}
                    </div>
                    <h3 className="mt-6 text-xl font-semibold">{t(`${role}Title`)}</h3>
                    <p className="mt-2 max-w-sm text-ink-soft">{t(`${role}Text`)}</p>
                  </article>
                )
              })}
            </div>
          </section>
        </Reveal>

        <Reveal className="py-14 md:py-20">
          <div className="mb-8 flex flex-wrap items-end justify-between gap-4">
            <h2 className="max-w-xl text-3xl font-semibold tracking-[-0.04em] md:text-5xl">{t('workspaceTitle')}</h2>
            <p className="max-w-sm text-ink-soft">{t('whyText')}</p>
          </div>
          <figure className="family-panorama mb-10 overflow-hidden">
            <img src="/family-mountain-picnic.jpg" alt={t('picnicAlt')} className="block aspect-[3.8/1] max-h-[30rem] min-h-44 w-full object-cover object-center" loading="lazy" />
          </figure>
          <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
            {features.map((key, index) => <article key={key} onPointerMove={trackGlow} className={`glow-card surface-lift landing-feature-cutout landing-feature-cutout--${index + 1} group relative min-h-64 overflow-hidden border border-line p-6 md:p-7 ${index === 0 ? 'bg-ink-panel text-panel-ink' : 'bg-surface text-ink'}`}>
              <div className="flex items-start justify-between"><span className={`text-sm font-semibold ${index === 0 ? 'text-lagoon' : 'text-brand-ink'}`}>0{index + 1}</span><span className={`grid size-9 place-items-center rounded-full transition-transform duration-200 group-hover:-translate-y-0.5 group-hover:translate-x-0.5 motion-reduce:transition-none ${index === 0 ? 'bg-panel-ink/10 text-panel-ink' : 'bg-muted text-brand-ink'}`}><Icon name="arrow" className="size-4 -rotate-45" /></span></div>
              <h3 className="mt-14 text-2xl font-semibold tracking-tight">{t(`${key}t`)}</h3>
              <p className={`mt-3 max-w-sm ${index === 0 ? 'text-panel-ink/75' : 'text-ink-soft'}`}>{t(`${key}d`)}</p>
            </article>)}
          </div>
        </Reveal>

        <Reveal>
          <section id="example" aria-labelledby="example-title" className="border-y border-line py-14 md:py-20">
            <div className="mb-8 grid items-center gap-7 md:grid-cols-[.8fr_1.2fr]">
              <div className="max-w-2xl">
                <p className="text-sm font-semibold uppercase tracking-[0.14em] text-accent-ink">{t('exampleEyebrow')}</p>
                <h2 id="example-title" className="mt-3 text-3xl font-semibold tracking-[-0.04em] md:text-5xl">{t('exampleTitle')}</h2>
                <p className="mt-4 text-ink-soft">{t('exampleIntro')}</p>
              </div>
              <InterfacePreview />
            </div>
            <ol className="relative grid gap-8 pt-8 md:grid-cols-3 md:gap-0 md:pt-48">
              <svg aria-hidden="true" viewBox="0 0 1200 190" preserveAspectRatio="none" className="pointer-events-none absolute inset-x-0 top-0 hidden h-48 w-full md:block">
                <g fill="var(--color-canvas)" stroke="var(--color-brand-strong)" strokeLinejoin="round" strokeWidth="2.5">
                  <path d="M34 72 73 40l39 32v35H34Z" /><path d="M28 73 73 34l45 39" fill="none" stroke="var(--color-accent)" strokeWidth="4" />
                  <path d="M66 82h14v25H66z" /><path d="M43 77h12v12H43zM91 77h12v12H91z" />
                  <path d="M363 145h22v22h-22zM385 153h20v14h-20z" fill="var(--color-accent)" />
                  <circle cx="433" cy="157" r="11" fill="var(--color-muted)" />
                  <path d="M1048 29v61m-29-40h58m-48 0-9 32m19-32v32m19-32 9 32" fill="none" />
                  <path d="M1019 49q29-34 58 0Z" fill="var(--color-accent)" />
                  <path d="M1139 70c-12-13 5-31 17-18 12-13 29 5 17 18l-17 17z" fill="var(--color-accent)" />
                  <path d="m1156 86-8 20" fill="none" />
                </g>
                <path d="M12 112 C150 112 150 24 292 24 S430 112 572 112 S710 24 852 24 S990 112 1188 112" fill="none" stroke="var(--color-brand-strong)" strokeWidth="3" />
                <path d="M12 120 C150 120 150 32 292 32 S430 120 572 120 S710 32 852 32 S990 120 1188 120" fill="none" stroke="var(--color-accent)" strokeWidth="3" />
                <g fill="var(--color-accent)">
                  <circle cx="95" cy="111" r="3" /><circle cx="170" cy="87" r="3" /><circle cx="215" cy="53" r="3" />
                  <circle cx="369" cy="51" r="3" /><circle cx="425" cy="87" r="3" /><circle cx="500" cy="111" r="3" />
                  <circle cx="644" cy="87" r="3" /><circle cx="700" cy="51" r="3" /><circle cx="778" cy="27" r="3" />
                  <circle cx="922" cy="52" r="3" /><circle cx="980" cy="88" r="3" /><circle cx="1060" cy="112" r="3" />
                </g>
                <g transform="translate(530 73) rotate(8 42 22)" stroke="var(--color-brand)" strokeWidth="3">
                  <g fill="var(--color-canvas)"><circle cx="29" cy="5" r="6" /><circle cx="53" cy="5" r="6" /></g>
                  <path d="M19 19 Q19 10 29 10 Q39 10 39 19 M43 19 Q43 10 53 10 Q63 10 63 19" fill="var(--color-brand-strong)" />
                  <path d="M5 19 H78 L74 30 H12 Z" fill="var(--color-accent)" />
                  <circle cx="23" cy="34" r="7" fill="var(--color-canvas)" />
                  <circle cx="67" cy="34" r="7" fill="var(--color-canvas)" />
                </g>
              </svg>
              {['exampleParent', 'exampleSpecialist', 'exampleControl'].map((key, index) => (
                <li key={key} className={`relative border-l border-accent/55 pl-5 md:border-l-0 md:px-6 md:pl-6 ${index === 1 ? 'md:translate-y-10' : index === 2 ? 'md:translate-y-2' : ''}`}>
                  <p className="font-display text-4xl leading-none text-accent-ink/80">0{index + 1}</p>
                  <p className="mt-5 max-w-xs text-lg font-semibold leading-snug">{t(`${key}Title`)}</p>
                  <p className="mt-3 max-w-sm leading-relaxed text-ink-soft">{t(`${key}Text`)}</p>
                </li>
              ))}
            </ol>
          </section>
        </Reveal>

        <section id="how-it-works" className="grid gap-10 py-16 md:grid-cols-[.8fr_1.2fr] md:py-24">
          <div className="max-w-md"><p className="text-sm font-semibold uppercase tracking-[0.14em] text-accent-ink">{t('howTitle')}</p><h2 className="mt-3 text-3xl font-semibold tracking-[-0.04em] md:text-5xl">{t('whyTitle')}</h2><p className="mt-5 text-ink-soft">{t('whyText')}</p></div>
          <ol className="border-t border-line">
            {steps.map((step, index) => <li key={step} className="group flex items-start gap-5 border-b border-line px-3 py-6 transition-[background-color] duration-200 hover:bg-surface focus-within:bg-surface motion-reduce:transition-none"><span className="pt-1 text-sm font-semibold text-accent-ink">0{index + 1}</span><p className="flex-1 text-lg font-medium leading-snug md:text-xl">{t(step)}</p><Icon name="arrow" className="mt-1 size-5 -rotate-45 text-brand-ink transition-transform duration-200 group-hover:-translate-y-0.5 group-hover:translate-x-0.5 motion-reduce:transition-none" /></li>)}
          </ol>
        </section>

        <section id="privacy" className="grid gap-10 border-y border-line py-14 md:grid-cols-2 md:py-20">
          <div><p className="text-sm font-semibold uppercase tracking-[0.14em] text-brand-ink">{t('privacyTitle')}</p><h2 className="mt-3 text-3xl font-semibold tracking-[-0.04em] md:text-4xl">{t('privacyTitle')}</h2></div>
          <div><p className="text-lg leading-relaxed text-ink-soft">{t('privacyText')}</p><div className="mt-6 flex items-center gap-3 border-t border-line pt-4"><span className="grid size-10 place-items-center rounded-full bg-brand text-on-brand"><Icon name="lock" /></span><p className="font-semibold">{t('adminText')}</p></div></div>
        </section>

        <section className="my-14 grid gap-6 bg-accent px-6 py-10 text-on-accent md:my-20 md:grid-cols-[1fr_auto] md:items-end md:px-10 md:py-12">
          <div><h2 className="max-w-2xl text-3xl font-semibold tracking-[-0.04em] md:text-5xl">{t('finalTitle')}</h2><p className="mt-3 max-w-xl text-on-accent">{t('finalText')}</p></div>
          <Link to="/register" className="group inline-flex min-h-12 items-center justify-center gap-3 rounded-lg bg-on-accent px-5 font-semibold text-accent transition-[transform,box-shadow] duration-200 ease-out-soft hover:-translate-y-0.5 hover:shadow-lift active:translate-y-px focus-visible:ring-2 focus-visible:ring-on-accent focus-visible:ring-offset-2 motion-reduce:transform-none motion-reduce:transition-none">{t('ctaRegister')}<Icon name="arrow" className="size-5 transition-transform duration-200 group-hover:translate-x-1 motion-reduce:transition-none" /></Link>
        </section>
      </main>
      <footer className="relative overflow-hidden border-t border-line px-5 py-10 md:px-10 md:py-14">
        <span aria-hidden="true" className="footer-watermark pointer-events-none absolute bottom-[-.25em] right-0 select-none" />
        <div className="relative mx-auto grid max-w-[90rem] gap-10 md:grid-cols-[1.4fr_.7fr_.7fr]">
          <div><Logo /><p className="mt-5 max-w-sm text-sm leading-relaxed text-ink-soft">{t('footerNote')}</p></div>
          <div><h2 className="font-semibold">{t('footerExplore')}</h2><ul className="mt-3 space-y-2 text-sm"><li><a className="text-ink-soft hover:text-brand-ink" href="#for-people">{t('navPeople')}</a></li><li><a className="text-ink-soft hover:text-brand-ink" href="#how-it-works">{t('navWorkflow')}</a></li><li><a className="text-ink-soft hover:text-brand-ink" href="#privacy">{t('navPrivacy')}</a></li></ul></div>
          <div><h2 className="font-semibold">{t('footerAccount')}</h2><ul className="mt-3 space-y-2 text-sm"><li><Link className="text-ink-soft hover:text-brand-ink" to="/login">{t('ctaLogin')}</Link></li><li><Link className="text-ink-soft hover:text-brand-ink" to="/register">{t('ctaRegister')}</Link></li></ul></div>
        </div>
      </footer>
    </div>
  )
}
