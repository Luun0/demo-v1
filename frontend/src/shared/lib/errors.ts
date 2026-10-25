import { isApiError } from '@/shared/api'

/** Maps an error to a common.* i18n key; never shows server internals. */
export function errorKey(e: unknown): string {
  if (isApiError(e)) {
    switch (e.status) {
      case 403:
        return 'common:errorForbidden'
      case 404:
        return 'common:errorNotFound'
      case 409:
        return 'common:errorConflict'
      case 429:
        return 'common:errorRateLimited'
    }
    return 'common:errorGeneric'
  }
  if (e instanceof TypeError) return 'common:errorNetwork'
  return 'common:errorGeneric'
}
