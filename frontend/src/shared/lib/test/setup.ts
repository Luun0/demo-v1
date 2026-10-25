import '@testing-library/jest-dom/vitest'
import { afterAll, afterEach, beforeAll } from 'vitest'
import { cleanup } from '@testing-library/react'
import { server } from './msw'
import { initI18n } from '@/shared/i18n'

beforeAll(async () => {
  server.listen({ onUnhandledRequest: 'error' })
  await initI18n('en')
})
afterEach(() => {
  cleanup()
  server.resetHandlers()
})
afterAll(() => server.close())
