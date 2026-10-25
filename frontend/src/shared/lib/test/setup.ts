import '@testing-library/jest-dom/vitest'
import { afterAll, afterEach, beforeAll } from 'vitest'
import { cleanup } from '@testing-library/react'
import { server } from './msw'
import { initI18n } from '@/shared/i18n'

const memoryStorage = new Map<string, string>()
const testStorage = {
  getItem: (key: string) => memoryStorage.get(key) ?? null,
  setItem: (key: string, value: string) => memoryStorage.set(key, String(value)),
  removeItem: (key: string) => memoryStorage.delete(key),
  clear: () => memoryStorage.clear(),
  key: (index: number) => [...memoryStorage.keys()][index] ?? null,
  get length() { return memoryStorage.size },
}
Object.defineProperty(globalThis, 'localStorage', { configurable: true, value: testStorage })

beforeAll(async () => {
  server.listen({ onUnhandledRequest: 'error' })
  await initI18n('en')
})
afterEach(() => {
  cleanup()
  server.resetHandlers()
})
afterAll(() => server.close())
