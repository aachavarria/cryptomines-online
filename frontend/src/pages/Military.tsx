import { useState, Suspense } from 'react'
import { Canvas } from '@react-three/fiber'
import ResourceHUD from '../components/layout/ResourceHUD.tsx'
import SideNav from '../components/layout/SideNav.tsx'
import ShipFactoryPanel from '../components/panels/ShipFactoryPanel.tsx'
import ShipDesignPanel from '../components/panels/ShipDesignPanel.tsx'
import BlueprintPanel from '../components/panels/BlueprintPanel.tsx'
import FleetPanel from '../components/panels/FleetPanel.tsx'
import InstancePanel from '../components/panels/InstancePanel.tsx'
import SpacedockPanel from '../components/panels/SpacedockPanel.tsx'
import { FrigateModel, CruiserModel, BattleshipModel } from '../components/three/ships/index.ts'

const TABS = [
  { id: 'factory', label: 'Ship Factory', icon: 'SF' },
  { id: 'designs', label: 'Designs', icon: 'DS' },
  { id: 'blueprints', label: 'Blueprints', icon: 'BP' },
  { id: 'fleets', label: 'Fleets', icon: 'FL' },
  { id: 'instances', label: 'Instances', icon: 'IN' },
  { id: 'spacedock', label: 'Spacedock', icon: 'SD' },
] as const

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
  const [activeTab, setActiveTab] = useState<TabId>('factory')

  function renderPanel() {
    switch (activeTab) {
      case 'factory': return <ShipFactoryPanel />
      case 'designs': return <ShipDesignPanel />
      case 'blueprints': return <BlueprintPanel />
      case 'fleets': return <FleetPanel />
      case 'instances': return <InstancePanel />
      case 'spacedock': return <SpacedockPanel />
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
          <div className="mil-tabs">
            {TABS.map(tab => (
              <button
                key={tab.id}
                className={`mil-tab ${activeTab === tab.id ? 'active' : ''}`}
                onClick={() => setActiveTab(tab.id)}
              >
                <span className="mil-tab-icon">{tab.icon}</span>
                <span className="mil-tab-label">{tab.label}</span>
              </button>
            ))}
          </div>

          {/* Active panel */}
          <div className="mil-content">
            {renderPanel()}
          </div>
        </div>
      </div>
    </div>
  )
}
