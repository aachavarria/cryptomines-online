import { useState } from 'react'
import { useResources } from '../../hooks/useResources.ts'
import { useGameContext } from '../../contexts/GameContext.tsx'
import { formatNumber } from '../../hooks/useCountdown.ts'

export default function ResourceHUD() {
  const { resources, collect } = useResources()
  const { state } = useGameContext()
  const [collecting, setCollecting] = useState(false)

  const civicCenter = state.buildings.find(b => b.type_name === 'civic_center')
  const ccLevel = civicCenter?.level ?? 0

  async function handleCollect() {
    setCollecting(true)
    await collect()
    setCollecting(false)
  }

  const pending = resources
    ? (resources.pending_metal + resources.pending_he3 + resources.pending_gold)
    : 0

  return (
    <div className="resource-hud">
      <div className="hud-logo">
        <span className="hud-logo-title">Cryptomines Online</span>
        {ccLevel > 0 && (
          <span className="hud-cc-badge">CC Lv {ccLevel}</span>
        )}
      </div>

      <div className="hud-divider" />

      <div className="hud-resources">
        {resources ? (
          <>
            <div className="hud-resource">
              <span className="hud-resource-icon metal">M</span>
              <span className="hud-resource-value">{formatNumber(resources.metal)}</span>
              <span className="hud-resource-rate">+{formatNumber(resources.metal_per_hour)}/hr</span>
            </div>
            <div className="hud-resource">
              <span className="hud-resource-icon he3">H</span>
              <span className="hud-resource-value">{formatNumber(resources.he3)}</span>
              <span className="hud-resource-rate">+{formatNumber(resources.he3_per_hour)}/hr</span>
            </div>
            <div className="hud-resource">
              <span className="hud-resource-icon gold">G</span>
              <span className="hud-resource-value">{formatNumber(resources.gold)}</span>
              <span className="hud-resource-rate">+{formatNumber(resources.gold_per_hour)}/hr</span>
            </div>

            <div className="hud-storage">
              Storage Cap: {formatNumber(resources.storage_capacity)}
            </div>
          </>
        ) : (
          <span style={{ color: 'var(--text-dim)' }}>Loading...</span>
        )}
      </div>

      <button
        className="hud-collect-btn"
        onClick={handleCollect}
        disabled={collecting || pending <= 0}
      >
        {collecting ? 'Collecting...' : `Collect +${formatNumber(pending)}`}
      </button>
    </div>
  )
}
