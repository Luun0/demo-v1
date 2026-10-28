import type { ReactNode } from 'react'

export type IconName =
  | 'home' | 'child' | 'calendar' | 'diary' | 'people' | 'access' | 'shield' | 'exit' | 'flower'
  | 'arrow' | 'check' | 'chevron' | 'clock' | 'close' | 'eye' | 'lock' | 'menu' | 'minus' | 'moon' | 'notebook' | 'plus' | 'search' | 'settings' | 'sun' | 'chat'

const paths: Record<IconName, ReactNode> = {
  home: <><path d="m3 11 9-8 9 8"/><path d="M5 10v10h14V10M9 20v-6h6v6"/></>,
  child: <><circle cx="12" cy="7" r="3"/><path d="M7 21v-3a5 5 0 0 1 10 0v3M8 11l-3 4M16 11l3 4"/></>,
  calendar: <><rect x="3.5" y="5.5" width="17" height="15" rx="2"/><path d="M8 3.5v4M16 3.5v4M3.5 10h17M8 14h2M14 14h2M8 17h2"/></>,
  diary: <><path d="M5 4h12a2 2 0 0 1 2 2v14H7a3 3 0 0 1-3-3V5a1 1 0 0 1 1-1Z"/><path d="M8 4v16M11 9h5M11 13h4"/></>,
  people: <><circle cx="9" cy="8" r="3"/><circle cx="17" cy="10" r="2.5"/><path d="M3.5 19v-1.5A4.5 4.5 0 0 1 8 13h2a4.5 4.5 0 0 1 4.5 4.5V19M14.5 14.5a4 4 0 0 1 6 3.5v1"/></>,
  access: <><path d="M4 19v-2a4 4 0 0 1 4-4h4"/><circle cx="10" cy="7" r="4"/><path d="M15 17h6M18 14v6"/></>,
  shield: <><path d="M12 3 20 6v5c0 5-3.2 8.2-8 10-4.8-1.8-8-5-8-10V6l8-3Z"/><path d="m8.5 12 2.2 2.2 4.8-5"/></>,
  exit: <><path d="M10 5H5v14h5M14 8l4 4-4 4M8 12h10"/></>,
  flower: <><circle cx="12" cy="12" r="2.2"/><circle cx="12" cy="6.5" r="3"/><circle cx="17" cy="10" r="3"/><circle cx="15" cy="16" r="3"/><circle cx="9" cy="16" r="3"/><circle cx="7" cy="10" r="3"/></>,
  arrow: <path d="M5 12h14m-6-6 6 6-6 6"/>,
  check: <path d="m5 12 4 4L19 6"/>,
  chevron: <path d="m9 18 6-6-6-6"/>,
  clock: <><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></>,
  close: <path d="m6 6 12 12M18 6 6 18"/>,
  eye: <><path d="M2.5 12s3.5-6 9.5-6 9.5 6 9.5 6-3.5 6-9.5 6-9.5-6-9.5-6Z"/><circle cx="12" cy="12" r="2.5"/></>,
  lock: <><rect x="4" y="10" width="16" height="11" rx="2"/><path d="M8 10V7a4 4 0 0 1 8 0v3"/></>,
  menu: <path d="M4 6h16M4 12h16M4 18h16"/>,
  moon: <path d="M20.5 14.2A8.5 8.5 0 0 1 9.8 3.5 8.5 8.5 0 1 0 20.5 14.2Z"/>,
  minus: <path d="M5 12h14"/>,
  notebook: <><path d="M6 3h12v18H6zM9 7h6M9 11h6M9 15h4"/><path d="M4 6h2M4 10h2M4 14h2M4 18h2"/></>,
  plus: <path d="M12 5v14M5 12h14"/>,
  search: <><circle cx="10.8" cy="10.8" r="6.8"/><path d="m16 16 5 5"/></>,
  settings: <><circle cx="12" cy="12" r="3"/><path d="m19.4 15 .1.1 1.4 1.1-1.4 2.4-1.7-.7a8 8 0 0 1-1.8 1l-.3 1.8h-2.8l-.3-1.8a8 8 0 0 1-1.8-1l-1.7.7-1.4-2.4 1.4-1.1a7 7 0 0 1 0-2l-1.4-1.1 1.4-2.4 1.7.7a8 8 0 0 1 1.8-1l.3-1.8h2.8l.3 1.8a8 8 0 0 1 1.8 1l1.7-.7 1.4 2.4-1.4 1.1a7 7 0 0 1 0 2Z"/></>,
  sun: <><circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M4.93 4.93l1.42 1.42m11.3 11.3 1.42 1.42M2 12h2m16 0h2M4.93 19.07l1.42-1.42m11.3-11.3 1.42-1.42"/></>,
  chat: <><path d="M20 11.5a7.5 7.5 0 0 1-7.5 7.5H5l-1-4a7.5 7.5 0 1 1 16-3.5Z"/><path d="M8 11h8M8 14h5"/></>,
}

export function Icon({ name, className = 'size-5' }: { name: IconName; className?: string }) {
  return <svg aria-hidden="true" className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={1.8} strokeLinecap="round" strokeLinejoin="round">{paths[name]}</svg>
}
