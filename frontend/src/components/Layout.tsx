import type { ReactNode } from 'react'

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="layout">
      <header className="header">
        <h1 className="header-title">Cryptomines Online</h1>
      </header>
      <main className="main">{children}</main>
    </div>
  )
}
