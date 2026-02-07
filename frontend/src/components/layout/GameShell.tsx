import { Canvas } from '@react-three/fiber'
import PlanetScene from '../three/PlanetScene.tsx'
import ResourceHUD from './ResourceHUD.tsx'
import SideNav from './SideNav.tsx'
import BuildingContextMenu from '../panels/BuildingContextMenu.tsx'
import BuildingDetailPanel from '../panels/BuildingDetailPanel.tsx'
import ConstructionPanel from '../panels/ConstructionPanel.tsx'
import ConstructionInfoPanel from '../panels/ConstructionInfoPanel.tsx'
import { useGameContext } from '../../contexts/GameContext.tsx'

export default function GameShell() {
  const { state, deselectAll, openConstructPanel } = useGameContext()

  if (state.loading) {
    return (
      <div className="loading-screen">
        <h2>Cryptomines Online</h2>
        <div className="loading-spinner" />
        <p>Loading your colony...</p>
      </div>
    )
  }

  if (state.error && !state.currentPlanet) {
    return (
      <div className="loading-screen">
        <h2>Cryptomines Online</h2>
        <p style={{ color: 'var(--accent-danger)' }}>{state.error}</p>
        <button
          className="hud-collect-btn"
          onClick={() => window.location.reload()}
        >
          Retry
        </button>
      </div>
    )
  }

  return (
    <div className="game-shell">
      {/* Three.js Canvas - full viewport */}
      <div className="game-canvas">
        <Canvas
          shadows
          orthographic
          camera={{
            zoom: 10,
            near: -500,
            far: 1000,
            position: [100, 100, 100],
          }}
          onPointerMissed={() => {
            if (!state.placementMode.active) {
              deselectAll()
            }
          }}
        >
          <PlanetScene />
        </Canvas>
      </div>

      {/* HTML Overlay Layer */}
      <div className="hud-overlay">
        <ResourceHUD />
        <SideNav />
        <BuildingContextMenu />
        <ConstructionInfoPanel />
        <BuildingDetailPanel />
        {state.showConstructPanel && <ConstructionPanel />}

        {/* Build button - bottom right */}
        {!state.showConstructPanel && !state.placementMode.active && (
          <div className="bottom-actions">
            <button
              className="build-action-btn"
              onClick={openConstructPanel}
            >
              Build
            </button>
          </div>
        )}

        {/* Placement mode indicator */}
        {state.placementMode.active && (
          <div className="placement-indicator">
            <span>Click a green tile to place building</span>
            <span className="placement-hint">Press ESC to cancel</span>
          </div>
        )}
      </div>
    </div>
  )
}
