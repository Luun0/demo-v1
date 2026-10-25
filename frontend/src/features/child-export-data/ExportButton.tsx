import { useTranslation } from 'react-i18next'
import { useMutation } from '@tanstack/react-query'
import { api, unwrap } from '@/shared/api'
import { Button, useNotify } from '@/shared/ui'

export function ExportButton({ childId }: { childId: string }) {
  const { t } = useTranslation('child')
  const notify = useNotify()
  const m = useMutation({
    mutationFn: async () => unwrap(await api.POST('/children/{childId}/exports', { params: { path: { childId } } })),
    onSuccess: (data) => {
      const url = URL.createObjectURL(new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' }))
      const a = Object.assign(document.createElement('a'), { href: url, download: `child-export-${childId}.json` })
      a.click()
      URL.revokeObjectURL(url)
      notify.ok('child:exported')
    },
    onError: notify.fail,
  })
  return <Button variant="secondary" busy={m.isPending} onClick={() => m.mutate()}>{t('export')}</Button>
}
