import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, unwrap, type User } from '@/shared/api'
import { Button, useNotify } from '@/shared/ui'
import { userKeys } from '@/entities/user'

type Patch = { status?: 'active' | 'inactive'; specialist_approved?: boolean }

export function UserActions({ user, self }: { user: User; self: boolean }) {
  const { t } = useTranslation('admin')
  const qc = useQueryClient()
  const notify = useNotify()
  const m = useMutation({
    mutationFn: async (body: Patch) => unwrap(await api.PATCH('/users/{userId}', { params: { path: { userId: user.id } }, body })),
    onSuccess: () => {
      notify.ok('admin:updated')
      void qc.invalidateQueries({ queryKey: userKeys.all })
    },
    onError: notify.fail,
  })
  const name = `${user.first_name} ${user.last_name}`
  return (
    <div className="flex flex-wrap gap-2">
      {user.role === 'specialist' && !user.specialist_approved && (
        <Button variant="secondary" busy={m.isPending} onClick={() => m.mutate({ specialist_approved: true })}>{t('approve')}<span className="sr-only">: {name}</span></Button>
      )}
      {!self && (user.status === 'active' ? (
        <Button variant="secondary" busy={m.isPending} onClick={() => m.mutate({ status: 'inactive' })}>{t('deactivate')}<span className="sr-only">: {name}</span></Button>
      ) : (
        <Button variant="secondary" busy={m.isPending} onClick={() => m.mutate({ status: 'active' })}>{t('activate')}<span className="sr-only">: {name}</span></Button>
      ))}
    </div>
  )
}
