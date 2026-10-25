import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { localInputToISO, toLocalInput } from '@/shared/lib'
import { Button, Select, Textarea, Input } from '@/shared/ui'

export const entrySchema = z.object({
  type: z.enum(['observation', 'progress', 'recommendation']),
  visibility: z.enum(['shared', 'specialists_only']),
  content: z.string().trim().min(1).max(10000),
  observed_at: z.string().min(1),
})
export type EntryValues = z.infer<typeof entrySchema>

export function toEntryBody(v: EntryValues) {
  return { ...v, observed_at: localInputToISO(v.observed_at) }
}

interface Props {
  defaults?: Partial<EntryValues>
  canChooseVisibility: boolean
  busy: boolean
  onSubmit: (v: EntryValues) => void
}

/** Parents never get a visibility choice: their entries are always shared. */
export function EntryForm({ defaults, canChooseVisibility, busy, onSubmit }: Props) {
  const { t } = useTranslation('diary')
  const { register, handleSubmit, formState } = useForm<EntryValues>({
    resolver: zodResolver(entrySchema),
    defaultValues: { type: 'observation', visibility: 'shared', observed_at: toLocalInput(new Date()), content: '', ...defaults },
  })
  return (
    <form onSubmit={handleSubmit(onSubmit)} noValidate>
      <Select label={t('type')} {...register('type')}>
        {(['observation', 'progress', 'recommendation'] as const).map((k) => <option key={k} value={k}>{t(`types.${k}`)}</option>)}
      </Select>
      {canChooseVisibility ? (
        <Select label={t('visibility')} {...register('visibility')}>
          <option value="shared">{t('visibilities.shared')}</option>
          <option value="specialists_only">{t('visibilities.specialists_only')}</option>
        </Select>
      ) : (
        <p className="mb-4 text-sm text-ink-soft">{t('parentShared')}</p>
      )}
      <Input label={t('observedAt')} type="datetime-local" error={formState.errors.observed_at && t('common:required')} {...register('observed_at')} />
      <Textarea label={t('content')} maxLength={10000} error={formState.errors.content && t('common:required')} {...register('content')} />
      <Button type="submit" busy={busy}>{t('common:save')}</Button>
    </form>
  )
}
