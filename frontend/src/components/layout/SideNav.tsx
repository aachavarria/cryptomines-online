import { useState, useEffect, type ComponentType } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import {
  Globe2,
  Rocket,
  FlaskConical,
  Shield,
  Box,
  Scroll,
  MessageSquare,
  UserCircle2,
  Map,
  Users,
} from 'lucide-react'
import { useQuests } from '../../hooks/useQuests.ts'
import { useGameContext } from '../../contexts/GameContext.tsx'
import QuestPanel from '../panels/QuestPanel.tsx'
import ResearchPanel from '../panels/ResearchPanel.tsx'
import ChatPanel from '../panels/ChatPanel.tsx'
import InventoryPanel from '../panels/InventoryPanel.tsx'
import CommandersListPanel from '../panels/CommandersListPanel.tsx'
import CorpsPanel from '../panels/CorpsPanel.tsx'
import GalaxyMapPanel from '../panels/GalaxyMapPanel.tsx'

type LucideIcon = ComponentType<{ size?: number | string; strokeWidth?: number | string; 'aria-hidden'?: boolean | 'true' | 'false' }>

interface NavItemConfig {
  id: string
  Icon: LucideIcon
  label: string
  route: string | null
  locked: boolean
}

const NAV_ITEMS: NavItemConfig[] = [
  { id: 'research', Icon: FlaskConical, label: 'Research', route: null, locked: false },
  { id: 'fleet', Icon: Shield, label: 'Military', route: '/military', locked: false },
  { id: 'inventory', Icon: Box, label: 'Inventory', route: null, locked: false },
  { id: 'quest', Icon: Scroll, label: 'Quests', route: null, locked: false },
  { id: 'chat', Icon: MessageSquare, label: 'Chat', route: null, locked: false },
  { id: 'commander', Icon: UserCircle2, label: 'Cmdr', route: null, locked: false },
  { id: 'galaxy', Icon: Map, label: 'Galaxy', route: null, locked: false },
  { id: 'corp', Icon: Users, label: 'Corp', route: null, locked: false },
]

export default function SideNav() {
  const [tooltip, setTooltip] = useState<string | null>(null)
  const [questOpen, setQuestOpen] = useState(false)
  const [researchOpen, setResearchOpen] = useState(false)
  const [chatOpen, setChatOpen] = useState(false)
  const [inventoryOpen, setInventoryOpen] = useState(false)
  const [commanderOpen, setCommanderOpen] = useState(false)
  const [corpOpen, setCorpOpen] = useState(false)
  const [galaxyOpen, setGalaxyOpen] = useState(false)
  const [devMode, setDevMode] = useState(() => {
    return localStorage.getItem('dev_mode') === 'true'
  })
  const navigate = useNavigate()
  const location = useLocation()
  const { claimableCount } = useQuests()
  const { state: gameState, setBaseView } = useGameContext()

  useEffect(() => {
    localStorage.setItem('dev_mode', devMode.toString())
  }, [devMode])

  // Listen for bridge panel control events (game:panel custom events)
  useEffect(() => {
    const handler = (e: Event) => {
      const { panel, action } = (e as CustomEvent<{ panel: string; action: 'open' | 'close' | 'toggle' }>).detail
      const setters: Record<string, (v: boolean | ((p: boolean) => boolean)) => void> = {
        quest: setQuestOpen,
        research: setResearchOpen,
        chat: setChatOpen,
        inventory: setInventoryOpen,
        commander: setCommanderOpen,
        corp: setCorpOpen,
        galaxy: setGalaxyOpen,
      }
      const setter = setters[panel]
      if (!setter) return

      if (action === 'open') {
        // Close all others first, then open the target
        Object.entries(setters).forEach(([key, s]) => {
          s(key === panel)
        })
      } else if (action === 'close') {
        setter(false)
      } else if (action === 'toggle') {
        // Close all others, toggle this one
        Object.entries(setters).forEach(([key, s]) => {
          if (key !== panel) s(false)
        })
        setter(prev => !prev)
      }
    }
    window.addEventListener('game:panel', handler)
    return () => window.removeEventListener('game:panel', handler)
  }, [])

  // Listen for bridge navigation events (game:navigate custom events)
  useEffect(() => {
    const handler = (e: Event) => {
      const { path } = (e as CustomEvent<{ path: string }>).detail
      navigate(path)
    }
    window.addEventListener('game:navigate', handler)
    return () => window.removeEventListener('game:navigate', handler)
  }, [navigate])

  function isActive(item: NavItemConfig): boolean {
    if (item.id === 'quest') return questOpen
    if (item.id === 'research') return researchOpen
    if (item.id === 'chat') return chatOpen
    if (item.id === 'inventory') return inventoryOpen
    if (item.id === 'commander') return commanderOpen
    if (item.id === 'corp') return corpOpen
    if (item.id === 'galaxy') return galaxyOpen
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
      setCommanderOpen(false)
      setCorpOpen(false)
      setGalaxyOpen(false)
      return
    }
    if (item.id === 'research') {
      setResearchOpen(prev => !prev)
      setQuestOpen(false)
      setChatOpen(false)
      setInventoryOpen(false)
      setCommanderOpen(false)
      setCorpOpen(false)
      setGalaxyOpen(false)
      return
    }
    if (item.id === 'chat') {
      setChatOpen(prev => !prev)
      setQuestOpen(false)
      setResearchOpen(false)
      setInventoryOpen(false)
      setCommanderOpen(false)
      setCorpOpen(false)
      setGalaxyOpen(false)
      return
    }
    if (item.id === 'inventory') {
      setInventoryOpen(prev => !prev)
      setQuestOpen(false)
      setResearchOpen(false)
      setChatOpen(false)
      setCommanderOpen(false)
      setCorpOpen(false)
      setGalaxyOpen(false)
      return
    }
    if (item.id === 'commander') {
      setCommanderOpen(prev => !prev)
      setQuestOpen(false)
      setResearchOpen(false)
      setChatOpen(false)
      setInventoryOpen(false)
      setCorpOpen(false)
      setGalaxyOpen(false)
      return
    }
    if (item.id === 'corp') {
      setCorpOpen(prev => !prev)
      setQuestOpen(false)
      setResearchOpen(false)
      setChatOpen(false)
      setInventoryOpen(false)
      setCommanderOpen(false)
      setGalaxyOpen(false)
      return
    }
    if (item.id === 'galaxy') {
      setGalaxyOpen(prev => !prev)
      setQuestOpen(false)
      setResearchOpen(false)
      setChatOpen(false)
      setInventoryOpen(false)
      setCommanderOpen(false)
      setCorpOpen(false)
      return
    }
    if (item.route) {
      navigate(item.route)
    }
  }

  function handleBaseSwitch(view: 'ground' | 'space') {
    setBaseView(view)
    const storedPlanetId = localStorage.getItem('current_planet_id')
    if (!location.pathname.startsWith('/planet')) {
      if (storedPlanetId) navigate(`/planet/${storedPlanetId}`)
      else navigate('/')
    }
  }

  const onPlanet = location.pathname.startsWith('/planet')
  const groundActive = onPlanet && gameState.currentBase === 'ground'
  const spaceActive = onPlanet && gameState.currentBase === 'space'

  return (
    <>
      <nav className="side-nav">
        {/* Ground / Space pill switcher */}
        <div className="nav-base-switcher" role="tablist" aria-label="Base view">
          <button
            type="button"
            role="tab"
            aria-selected={groundActive}
            className={`nav-base-pill ${groundActive ? 'active' : ''}`}
            onClick={() => handleBaseSwitch('ground')}
            onMouseEnter={() => setTooltip('base')}
            onMouseLeave={() => setTooltip(null)}
            title="Ground Base"
          >
            <Globe2 size={20} strokeWidth={1.75} aria-hidden="true" />
            <span className="nav-base-pill-label">Ground</span>
            {tooltip === 'base' && <span className="ds-tooltip nav-tooltip">Ground Base</span>}
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={spaceActive}
            className={`nav-base-pill ${spaceActive ? 'active' : ''}`}
            onClick={() => handleBaseSwitch('space')}
            onMouseEnter={() => setTooltip('space')}
            onMouseLeave={() => setTooltip(null)}
            title="Space Base"
          >
            <Rocket size={20} strokeWidth={1.75} aria-hidden="true" />
            <span className="nav-base-pill-label">Space</span>
            {tooltip === 'space' && <span className="ds-tooltip nav-tooltip">Space Base</span>}
          </button>
        </div>

        <div className="nav-divider" />

        {NAV_ITEMS.map(item => {
          const Icon = item.Icon
          const active = isActive(item)
          return (
            <button
              key={item.id}
              type="button"
              className={`nav-item ${active ? 'active' : ''} ${item.locked ? 'locked' : ''}`}
              aria-pressed={active}
              onMouseEnter={() => setTooltip(item.id)}
              onMouseLeave={() => setTooltip(null)}
              onClick={() => handleClick(item)}
            >
              <Icon size={22} strokeWidth={1.75} aria-hidden="true" />
              <span className="nav-item-label">{item.label}</span>
              {item.id === 'quest' && claimableCount > 0 && (
                <span className="nav-quest-badge" aria-label={`${claimableCount} claimable quests`} />
              )}
              {tooltip === item.id && (
                <span className="ds-tooltip nav-tooltip">
                  {item.locked ? 'Coming Soon' : item.label}
                </span>
              )}
            </button>
          )
        })}

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
      {inventoryOpen && <InventoryPanel onClose={() => setInventoryOpen(false)} />}
      {commanderOpen && <CommandersListPanel onClose={() => setCommanderOpen(false)} />}
      {corpOpen && <CorpsPanel onClose={() => setCorpOpen(false)} />}
      {galaxyOpen && <GalaxyMapPanel onClose={() => setGalaxyOpen(false)} />}
    </>
  )
}
