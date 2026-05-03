import { Canvas } from "@react-three/fiber";
import { Hammer } from "lucide-react";
import PlanetScene from "../three/PlanetScene.tsx";
import ResourceHUD from "./ResourceHUD.tsx";
import SideNav from "./SideNav.tsx";
import BuildingContextMenu from "../panels/BuildingContextMenu.tsx";
import BuildingDetailPanel from "../panels/BuildingDetailPanel.tsx";
import CommandCenterPanel from "../panels/CommandCenterPanel.tsx";
import ConstructionPanel from "../panels/ConstructionPanel.tsx";
import ConstructionInfoPanel from "../panels/ConstructionInfoPanel.tsx";
import { useGameContext } from "../../contexts/GameContext.tsx";

export default function GameShell() {
  const { state, deselectAll, openConstructPanel } = useGameContext();

  if (state.loading) {
    return (
      <div className="loading-screen">
        <h2 className="ds-h1">Cryptomines Online</h2>
        <div className="loading-spinner" />
        <p className="ds-text-muted">Loading your colony...</p>
      </div>
    );
  }

  if (state.error && !state.currentPlanet) {
    return (
      <div className="loading-screen">
        <h2 className="ds-h1">Cryptomines Online</h2>
        <p className="loading-error">{state.error}</p>
        <button
          className="ds-btn ds-btn-secondary"
          onClick={() => window.location.reload()}
        >
          Retry
        </button>
      </div>
    );
  }

  return (
    <div className="game-shell">
      {/* Three.js Canvas - full viewport */}
      <div className="game-canvas">
        <Canvas
          shadows
          orthographic
          camera={{
            zoom: 30,
            near: -500,
            far: 1000,
            position: [100, 100, 100],
          }}
          onPointerMissed={() => {
            if (!state.placementMode.active) {
              deselectAll();
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
        <CommandCenterPanel />
        {state.showConstructPanel && <ConstructionPanel />}

        {/* Build button - bottom right */}
        {!state.showConstructPanel && !state.placementMode.active && (
          <div className="bottom-actions">
            <button
              type="button"
              className="build-action-btn ds-btn ds-btn-primary"
              onClick={openConstructPanel}
            >
              <Hammer size={18} strokeWidth={1.75} aria-hidden="true" />
              Build
            </button>
          </div>
        )}

        {/* Placement mode indicator */}
        {state.placementMode.active && (
          <div className="placement-indicator ds-panel">
            <span className="ds-h3">
              Placing on {state.currentBase === 'space' ? 'Space Base' : 'Ground Base'}
            </span>
            <span className="placement-hint ds-text-muted">
              Click a green tile, or press ESC to cancel
            </span>
          </div>
        )}
      </div>
    </div>
  );
}
