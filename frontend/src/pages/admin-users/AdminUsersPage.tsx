import { useTranslation } from 'react-i18next'
import { Badge, ErrorState, Loading, PageHeader } from '@/shared/ui'
import { useUsers } from '@/entities/user'
import { useMe } from '@/entities/viewer'
import { CreateUserDialog, UserActions } from '@/features/admin-manage-user'

export function AdminUsersPage() {
  const { t } = useTranslation('admin')
  const me = useMe()
  const q = useUsers()
  return (
    <>
      <PageHeader title={t('usersTitle')} actions={<CreateUserDialog />} />
      {q.isPending ? <Loading /> : q.isError ? <ErrorState error={q.error} onRetry={() => void q.refetch()} /> : (
        <div className="overflow-x-auto rounded-2xl border border-line bg-surface shadow-soft">
          <table className="w-full text-left text-sm">
            <caption className="sr-only">{t('usersTitle')}</caption>
            <thead className="border-b border-line bg-muted">
              <tr><th scope="col" className="p-3">{t('name')}</th><th scope="col" className="p-3">{t('email')}</th><th scope="col" className="p-3">{t('role')}</th><th scope="col" className="p-3">{t('status')}</th><th scope="col" className="p-3"><span className="sr-only">{t('actions')}</span></th></tr>
            </thead>
            <tbody>
              {q.data.map((u) => (
                <tr key={u.id} className="border-b border-line last:border-0">
                  <td className="p-3 font-medium">{u.first_name} {u.last_name}</td>
                  <td className="p-3">{u.email}</td>
                  <td className="p-3">{t(`common:roles.${u.role}`)}</td>
                  <td className="p-3">
                    <Badge tone={u.status === 'active' ? 'brand' : 'neutral'}>{t(`statuses.${u.status}`)}</Badge>
                    {u.role === 'specialist' && !u.specialist_approved && <Badge tone="warn">{t('notApproved')}</Badge>}
                  </td>
                  <td className="p-3"><UserActions user={u} self={u.id === me.id} /></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </>
  )
}
