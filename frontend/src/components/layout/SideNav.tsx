import { useState, useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useQuests } from '../../hooks/useQuests.ts'
import QuestPanel from '../panels/QuestPanel.tsx'
import ResearchPanel from '../panels/ResearchPanel.tsx'
import ChatPanel from '../panels/ChatPanel.tsx'
import InventoryPanel from '../panels/InventoryPanel.tsx'

interface NavItemConfig {
  id: string
  icon: string
  label: string
  route: string | null
  locked: boolean
}

const NAV_ITEMS: NavItemConfig[] = [
  { id: 'base', icon: '\u{1F30D}', label: 'Base', route: null, locked: false },
  { id: 'research', icon: '\u{1F52C}', label: 'Research', route: null, locked: false },
  { id: 'fleet', icon: '\u{1F680}', label: 'Military', route: '/military', locked: false },
  { id: 'inventory', icon: '\u{1F4E6}', label: 'Inventory', route: null, locked: false },
  { id: 'quest', icon: '\u{1F4DC}', label: 'Quests', route: null, locked: false },
  { id: 'chat', icon: '\u{1F4AC}', label: 'Chat', route: null, locked: false },
  { id: 'commander', icon: '\u{1F464}', label: 'Cmdr', route: null, locked: true },
  { id: 'galaxy', icon: '\u{1F5FA}', label: 'Galaxy', route: null, locked: true },
  { id: 'corp', icon: '\u{1F6E1}', label: 'Corp', route: null, locked: true },
]

export default function SideNav() {
  const [tooltip, setTooltip] = useState<string | null>(null)
  const [questOpen, setQuestOpen] = useState(false)
  const [researchOpen, setResearchOpen] = useState(false)
  const [chatOpen, setChatOpen] = useState(false)
  const [inventoryOpen, setInventoryOpen] = useState(false)
  const [devMode, setDevMode] = useState(() => {
    return localStorage.getItem('dev_mode') === 'true'
  })
  const navigate = useNavigate()
  const location = useLocation()
  const { claimableCount } = useQuests()

  useEffect(() => {
    localStorage.setItem('dev_mode', devMode.toString())
  }, [devMode])

  function isActive(item: NavItemConfig): boolean {
    if (item.id === 'base') return location.pathname.startsWith('/planet')
    if (item.id === 'quest') return questOpen
    if (item.id === 'research') return researchOpen
    if (item.id === 'chat') return chatOpen
    if (item.id === 'inventory') return inventoryOpen
    if (item.route) return location.pathname === item.route
    return false
  }

  function handleClick(item: NavItemConfig) {
    if (item.locked) return
    if (item.id === 'quest') {
      setQuestOpen(prev => !prev)
      setResearchOpen(false)
      setChatOpen(false)
      setInventoryOpen(false)
      return
    }
    if (item.id === 'research') {
      setResearchOpen(prev => !prev)
      setQuestOpen(false)
      setChatOpen(false)
      setInventoryOpen(false)
      return
    }
    if (item.id === 'chat') {
      setChatOpen(prev => !prev)
      setQuestOpen(false)
      setResearchOpen(false)
      setInventoryOpen(false)
      return
    }
    if (item.id === 'inventory') {
      setInventoryOpen(prev => !prev)
      setQuestOpen(false)
      setResearchOpen(false)
      setChatOpen(false)
      return
    }
    if (item.id === 'base') {
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
    <>
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
            {item.id === 'quest' && claimableCount > 0 && (
              <span className="nav-quest-badge" />
            )}
            {tooltip === item.id && item.locked && (
              <span className="nav-tooltip">Coming Soon</span>
            )}
          </button>
        ))}

        <div className="dev-mode-toggle">
          <label>
            <input
              type="checkbox"
              checked={devMode}
              onChange={(e) => setDevMode(e.target.checked)}
            />
            <span className="dev-mode-label">DEV</span>
          </label>
        </div>
      </nav>

      {questOpen && <QuestPanel onClose={() => setQuestOpen(false)} />}
      {researchOpen && <ResearchPanel onClose={() => setResearchOpen(false)} />}
      {chatOpen && <ChatPanel onClose={() => setChatOpen(false)} />}
      {inventoryOpen && <InventoryPanel />}
    </>
  )
}
