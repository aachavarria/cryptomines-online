import { useState, Suspense, useEffect } from 'react'
import { Canvas } from '@react-three/fiber'
import { useLocation } from 'react-router-dom'
import {
  Factory,
  PenTool,
  ScrollText,
  Rocket,
  Layers,
  Anchor,
  Recycle,
  Crosshair,
  ClipboardList,
  type LucideIcon,
} from 'lucide-react'
import ResourceHUD from '../components/layout/ResourceHUD.tsx'
import SideNav from '../components/layout/SideNav.tsx'
import ShipFactoryPanel from '../components/panels/ShipFactoryPanel.tsx'
import ShipDesignPanel from '../components/panels/ShipDesignPanel.tsx'
import BlueprintPanel from '../components/panels/BlueprintPanel.tsx'
import FleetPanel from '../components/panels/FleetPanel.tsx'
import InstancePanel from '../components/panels/InstancePanel.tsx'
import SpacedockPanel from '../components/panels/SpacedockPanel.tsx'
import RecyclingPlantPanel from '../components/panels/RecyclingPlantPanel.tsx'
import PvPPanel from '../components/panels/PvPPanel.tsx'
import { CombatReportsPanel } from '../components/panels/CombatReportsPanel.tsx'
import { FrigateModel, CruiserModel, BattleshipModel } from '../components/three/ships/index.ts'

type TabDef = { id: string; label: string; Icon: LucideIcon }

const TABS = [
  { id: 'factory',    label: 'Ship Factory',    Icon: Factory },
  { id: 'designs',    label: 'Designs',         Icon: PenTool },
  { id: 'blueprints', label: 'Blueprints',      Icon: ScrollText },
  { id: 'fleets',     label: 'Fleets',          Icon: Rocket },
  { id: 'instances',  label: 'Instances',       Icon: Layers },
  { id: 'spacedock',  label: 'Spacedock',       Icon: Anchor },
  { id: 'recycling',  label: 'Recycling Plant', Icon: Recycle },
  { id: 'pvp',        label: 'PvP Combat',      Icon: Crosshair },
  { id: 'reports',    label: 'Combat Reports',  Icon: ClipboardList },
] as const satisfies readonly TabDef[]

type TabId = typeof TABS[number]['id']

function ShipShowcase() {
  return (
    <>
      <ambientLight intensity={0.3} color="#1a1a3a" />
      <directionalLight position={[10, 15, 5]} intensity={1.0} color="#aabbff" />
      <directionalLight position={[-8, 5, -10]} intensity={0.2} color="#4466aa" />

      {/* Rotating showcase of ship models */}
      <group rotation={[0, 0, 0]}>
        <FrigateModel position={[-4, 0.5, 0]} scale={1.2} color="#44aaff" />
        <CruiserModel position={[0, 0.5, 0]} scale={1.0} color="#cc8844" />
        <BattleshipModel position={[5, 0.5, 0]} scale={0.8} color="#aa4444" />
      </group>

      {/* Simple starfield background */}
      {Array.from({ length: 200 }).map((_, i) => (
        <mesh key={i} position={[
          (Math.random() - 0.5) * 100,
          (Math.random() - 0.5) * 60,
          -20 - Math.random() * 40,
        ]}>
          <sphereGeometry args={[0.05 + Math.random() * 0.08, 4, 4]} />
          <meshBasicMaterial color="#ffffff" opacity={0.3 + Math.random() * 0.7} transparent />
        </mesh>
      ))}
    </>
  )
}

export default function Military() {
  const location = useLocation()
  const [activeTab, setActiveTab] = useState<TabId>(() => {
    const params = new URLSearchParams(location.search)
    const tab = params.get('tab')
    return TABS.some(t => t.id === tab) ? (tab as TabId) : 'factory'
  })

  useEffect(() => {
    const params = new URLSearchParams(location.search)
    const tab = params.get('tab')
    if (tab && TABS.some(t => t.id === tab)) {
      setActiveTab(tab as TabId)
    }
  }, [location.search])

  function renderPanel() {
    switch (activeTab) {
      case 'factory': return <ShipFactoryPanel />
      case 'designs': return <ShipDesignPanel />
      case 'blueprints': return <BlueprintPanel />
      case 'fleets': return <FleetPanel />
      case 'instances': return <InstancePanel />
      case 'spacedock': return <SpacedockPanel />
      case 'recycling': return <RecyclingPlantPanel />
      case 'pvp': return <PvPPanel />
      case 'reports': return <CombatReportsPanel />
    }
  }

  return (
    <div className="game-shell">
      {/* 3D background - ship showcase */}
      <div className="game-canvas military-canvas">
        <Canvas
          camera={{ fov: 50, near: 0.1, far: 200, position: [0, 3, 12] }}
        >
          <Suspense fallback={null}>
            <ShipShowcase />
          </Suspense>
        </Canvas>
      </div>

      {/* HTML overlay */}
      <div className="hud-overlay">
        <ResourceHUD />
        <SideNav />

        <div className="military-layout">
          {/* Tab navigation */}
          <div
            className="ds-tabs military-tabs"
            role="tablist"
            aria-label="Military sections"
          >
            {TABS.map(({ id, label, Icon }) => {
              const selected = activeTab === id
              return (
                <button
                  key={id}
                  type="button"
                  role="tab"
                  className="ds-tab"
                  aria-selected={selected}
                  aria-controls={`mil-panel-${id}`}
                  id={`mil-tab-${id}`}
                  onClick={() => setActiveTab(id)}
                >
                  <Icon
                    size={16}
                    strokeWidth={selected ? 2 : 1.75}
                    aria-hidden="true"
                  />
                  <span>{label}</span>
                </button>
              )
            })}
          </div>

          {/* Active panel */}
          <div
            className="military-content"
            role="tabpanel"
            id={`mil-panel-${activeTab}`}
            aria-labelledby={`mil-tab-${activeTab}`}
          >
            {renderPanel()}
          </div>
        </div>
      </div>

      <style>{`
        .military-canvas { opacity: 0.3; }
        .military-layout {
          position: fixed;
          top: var(--hud-height);
          left: var(--sidenav-width);
          right: 0;
          bottom: 0;
          display: flex;
          flex-direction: column;
          z-index: 50;
        }
        .military-tabs {
          padding: var(--sp-2) var(--sp-5) 0;
          background: var(--ds-surface);
          border-bottom: 1px solid var(--ds-border);
          flex-shrink: 0;
          overflow-x: auto;
          gap: var(--sp-1);
        }
        .military-tabs .ds-tab {
          display: inline-flex;
          align-items: center;
          gap: var(--sp-2);
          white-space: nowrap;
        }
        .military-content {
          flex: 1;
          overflow-y: auto;
          padding: var(--sp-5);
        }
        @media (max-width: 768px) {
          .military-tabs .ds-tab span { display: none; }
        }
      `}</style>
    </div>
  )
}
