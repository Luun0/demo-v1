/// <reference types="vitest/config" />
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'node:path'

const api = process.env.VITE_API_PROXY ?? 'http://127.0.0.1:8080'

// Same CSP as deploy/nginx.conf, applied by the preview server (prod-like).
const csp =
  "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; " +
  "connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'; object-src 'none'"

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: { alias: { '@': path.resolve(import.meta.dirname, './src') } },
  server: { host: '127.0.0.1', port: 5174, strictPort: true, proxy: { '/api': api } },
  preview: {
    host: '127.0.0.1',
    port: 4173,
    proxy: { '/api': api },
    headers: {
      'Content-Security-Policy': csp,
      'X-Content-Type-Options': 'nosniff',
      'Referrer-Policy': 'no-referrer',
      'X-Frame-Options': 'DENY',
      'Permissions-Policy': 'camera=(), microphone=(), geolocation=()',
    },
  },
  build: { target: 'es2022', cssCodeSplit: true },
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/shared/lib/test/setup.ts'],
    include: ['src/**/*.test.{ts,tsx}'],
    css: false,
    coverage: { provider: 'v8', include: ['src/**'], exclude: ['src/**/*.gen.ts', 'src/**/*.test.*', 'src/shared/lib/test/**'] },
  },
})
