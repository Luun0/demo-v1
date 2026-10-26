import type { ReactNode } from 'react'
import * as RD from '@radix-ui/react-dialog'
import { useTranslation } from 'react-i18next'

interface DialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  description?: string
  children: ReactNode
}

/** Accessible modal: focus trap, Esc, focus return, scroll lock (Radix). */
export function Dialog({ open, onOpenChange, title, description, children }: DialogProps) {
  const { t } = useTranslation()
  return (
    <RD.Root open={open} onOpenChange={onOpenChange}>
      <RD.Portal>
        <RD.Overlay className="fixed inset-0 z-40 bg-overlay backdrop-blur-sm" />
        <RD.Content
          className="dialog-surface fixed left-1/2 top-1/2 z-50 max-h-[90vh] w-[min(92vw,36rem)] -translate-x-1/2 -translate-y-1/2 overflow-y-auto rounded-[1.25rem] border border-line bg-surface p-6 text-ink shadow-panel"
          {...(description ? {} : { 'aria-describedby': undefined })}
        >
          <div className="mb-4 flex items-start justify-between gap-4">
            <RD.Title className="text-xl font-semibold">{title}</RD.Title>
            <RD.Close className="min-h-11 min-w-11 rounded-xl transition-[background-color] hover:bg-muted focus-visible:ring-2 focus-visible:ring-accent motion-reduce:transition-none" aria-label={t('close')}>
              <span aria-hidden="true">×</span>
            </RD.Close>
          </div>
          {description && <RD.Description className="mb-4 text-ink-soft">{description}</RD.Description>}
          {children}
        </RD.Content>
      </RD.Portal>
    </RD.Root>
  )
}
