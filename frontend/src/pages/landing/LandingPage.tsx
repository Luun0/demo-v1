import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { LanguageSwitch } from '@/features/language-switch'
import { Icon, Logo } from '@/shared/ui'

const features = [
  { key: 'f1', icon: 'calendar' },
  { key: 'f2', icon: 'diary' },
  { key: 'f3', icon: 'access' },
] as const

const roles = ['parent', 'specialist', 'admin'] as const

export function LandingPage() {
  const { t } = useTranslation('landing')
  return (
    <div className="min-h-dvh overflow-hidden bg-canvas text-ink">
      <a href="#main" className="sr-only focus:not-sr-only focus:fixed focus:left-2 focus:top-2 focus:z-50 focus:rounded focus:bg-surface focus:p-3">{t('common:skipToContent')}</a>
      <header className="mx-auto flex max-w-7xl items-center justify-between gap-4 px-5 py-5">
        <Logo />
        <div className="flex items-center gap-3">
          <LanguageSwitch />
          <Link to="/login" className="rounded-full px-4 py-2 font-medium text-brand hover:bg-brand/5 dark:text-brand-strong">{t('ctaLogin')}</Link>
        </div>
      </header>
      <main id="main" className="mx-auto max-w-7xl px-5 pb-14">
        <section className="grid items-center gap-10 py-12 md:grid-cols-[1.2fr_.8fr] md:py-20">
          <div>
            <p className="inline-flex items-center gap-2 rounded-full bg-accent/10 px-4 py-2 text-sm font-semibold text-brand-strong dark:text-white"><span className="size-2 rounded-full bg-accent" />{t('eyebrow')}</p>
            <h1 className="mt-6 max-w-3xl text-4xl font-bold leading-tight tracking-tight sm:text-5xl lg:text-6xl">{t('title')}</h1>
            <p className="mt-5 max-w-2xl text-lg leading-relaxed text-ink-soft">{t('text')}</p>
            <div className="mt-8 flex flex-wrap gap-3">
              <Link to="/register" className="inline-flex min-h-11 items-center gap-2 rounded-full bg-accent px-5 font-semibold text-ink hover:brightness-95 dark:text-brand">{t('ctaRegister')}<Icon name="arrow" /></Link>
              <Link to="/login" className="inline-flex min-h-11 items-center rounded-full border border-line bg-surface px-5 font-semibold text-ink hover:bg-muted">{t('ctaLogin')}</Link>
            </div>
          </div>
          <div className="relative mx-auto grid aspect-square w-full max-w-sm place-items-center rounded-[2.5rem] bg-gradient-to-br from-brand/10 via-accent/10 to-brand/5">
            <span className="absolute right-6 top-8 size-20 rounded-full bg-accent/15" aria-hidden="true" />
            <span className="absolute bottom-8 left-5 size-12 rounded-full bg-brand/10" aria-hidden="true" />
            <span className="grid size-56 place-items-center bg-accent text-white shadow-soft dark:text-brand [clip-path:polygon(50%_0,93%_25%,93%_75%,50%_100%,7%_75%,7%_25%)] sm:size-64">
              <Icon name="flower" className="size-32" />
            </span>
            <span className="absolute bottom-5 right-5 rounded-2xl border border-line bg-surface px-4 py-3 text-sm font-medium shadow-soft">{t('common:demoNotice')}</span>
          </div>
        </section>

        <section aria-label={t('workspaceTitle')} className="grid gap-4 md:grid-cols-3">
          {features.map(({ key, icon }) => (
            <article key={key} className="rounded-2xl border border-line bg-surface p-6 shadow-soft">
              <span className="mb-5 grid size-12 place-items-center rounded-2xl bg-brand/10 text-brand"><Icon name={icon} /></span>
              <h2 className="text-lg font-semibold">{t(`${key}t`)}</h2>
              <p className="mt-2 text-ink-soft">{t(`${key}d`)}</p>
            </article>
          ))}
        </section>

        <section className="py-16">
          <h2 className="mb-6 text-2xl font-semibold">{t('workspaceTitle')}</h2>
          <div className="grid gap-4 md:grid-cols-3">
            {roles.map((role, index) => (
              <article key={role} className="rounded-2xl border border-line bg-surface p-6 shadow-soft">
                <span className="text-sm font-semibold text-brand-strong">0{index + 1}</span>
                <h3 className="mt-3 text-lg font-semibold">{t(`${role}Title`)}</h3>
                <p className="mt-2 text-ink-soft">{t(`${role}Text`)}</p>
              </article>
            ))}
          </div>
        </section>

        <section className="relative overflow-hidden rounded-[2rem] bg-gradient-to-r from-brand to-brand-strong px-7 py-10 text-white shadow-soft md:px-12 md:py-14">
          <div className="absolute -right-8 -top-16 size-64 rotate-12 border-[24px] border-white/10 [clip-path:polygon(50%_0,93%_25%,93%_75%,50%_100%,7%_75%,7%_25%)]" aria-hidden="true" />
          <p className="relative max-w-3xl text-2xl font-semibold leading-relaxed">{t('text')}</p>
        </section>

        <section className="grid gap-10 py-16 md:grid-cols-2">
          <div>
            <h2 className="text-2xl font-semibold">{t('whyTitle')}</h2>
            <details className="mt-5 rounded-2xl border border-line bg-surface p-5 shadow-soft">
              <summary className="cursor-pointer font-semibold marker:text-accent">{t('whyTitle')}</summary>
              <p className="mt-3 leading-relaxed text-ink-soft">{t('whyText')}</p>
            </details>
          </div>
          <div>
            <h2 className="text-2xl font-semibold">{t('howTitle')}</h2>
            <ol className="mt-5 space-y-3 rounded-2xl bg-gradient-to-br from-brand/10 via-accent/10 to-brand/5 p-5">
              {(['step1', 'step2', 'step3'] as const).map((step, index) => (
                <li key={step} className="flex gap-4 rounded-xl bg-surface/80 p-4"><span className="font-bold text-accent">0{index + 1}</span><span>{t(step)}</span></li>
              ))}
            </ol>
          </div>
        </section>

        <section className="grid items-center gap-5 rounded-2xl border border-line bg-surface p-7 shadow-soft md:grid-cols-[auto_1fr] md:p-10">
          <span className="grid size-14 place-items-center rounded-2xl bg-brand text-white"><Icon name="shield" className="size-7" /></span>
          <div><h2 className="text-xl font-semibold">{t('privacyTitle')}</h2><p className="mt-2 text-ink-soft">{t('privacyText')}</p></div>
        </section>

        <section className="my-16 rounded-[2rem] bg-accent px-6 py-12 text-center text-white shadow-soft md:px-12">
          <h2 className="text-3xl font-bold text-ink dark:text-brand">{t('finalTitle')}</h2>
          <p className="mx-auto mt-3 max-w-xl text-ink dark:text-brand">{t('finalText')}</p>
          <Link to="/register" className="mt-7 inline-flex min-h-12 items-center rounded-full bg-white px-6 font-semibold text-brand hover:bg-white/90">{t('ctaRegister')}</Link>
        </section>
      </main>
      <footer className="border-t border-line px-5 py-7 text-center text-sm text-ink-soft">{t('footerNote')}</footer>
    </div>
  )
}
