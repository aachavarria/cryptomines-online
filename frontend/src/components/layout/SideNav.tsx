import { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'

interface NavItemConfig {
  id: string
  icon: string
  label: string
  route: string | null
  locked: boolean
}

const NAV_ITEMS: NavItemConfig[] = [
  { id: 'base', icon: '\u{1F30D}', label: 'Base', route: null, locked: false },
  { id: 'research', icon: '\u{1F52C}', label: 'Research', route: null, locked: true },
  { id: 'fleet', icon: '\u{1F680}', label: 'Military', route: '/military', locked: false },
  { id: 'commander', icon: '\u{1F464}', label: 'Cmdr', route: null, locked: true },
  { id: 'galaxy', icon: '\u{1F5FA}', label: 'Galaxy', route: null, locked: true },
  { id: 'corp', icon: '\u{1F6E1}', label: 'Corp', route: null, locked: true },
]

export default function SideNav() {
  const [tooltip, setTooltip] = useState<string | null>(null)
  const navigate = useNavigate()
  const location = useLocation()

  function isActive(item: NavItemConfig): boolean {
    if (item.id === 'base') return location.pathname.startsWith('/planet')
    if (item.route) return location.pathname === item.route
    return false
  }

  function handleClick(item: NavItemConfig) {
    if (item.locked) return
    if (item.id === 'base') {
      // Navigate to planet - use stored planet or go home
      const storedPlanetId = localStorage.getItem('current_planet_id')
      if (storedPlanetId) {
        navigate(`/planet/${storedPlanetId}`)
      } else {
        navigate('/')
      }
      return
    }
    if (item.route) {
      navigate(item.route)
    }
  }

  return (
    <nav className="side-nav">
      {NAV_ITEMS.map(item => (
        <button
          key={item.id}
          className={`nav-item ${isActive(item) ? 'active' : ''} ${item.locked ? 'locked' : ''}`}
          onMouseEnter={() => setTooltip(item.id)}
          onMouseLeave={() => setTooltip(null)}
          onClick={() => handleClick(item)}
        >
          <span>{item.icon}</span>
          <span className="nav-item-label">{item.label}</span>
          {tooltip === item.id && item.locked && (
            <span className="nav-tooltip">Coming Soon</span>
          )}
        </button>
      ))}
    </nav>
  )
}
