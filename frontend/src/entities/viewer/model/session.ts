import { useSyncExternalStore } from 'react'
import { getAccessToken, subscribeToken } from '@/shared/api'

export function useAccessToken(): string | null {
  return useSyncExternalStore(subscribeToken, getAccessToken, getAccessToken)
}
