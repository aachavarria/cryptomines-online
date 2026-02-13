import { useState, useEffect } from 'react'
import { useResources } from '../../hooks/useResources.ts'
import { useGameContext } from '../../contexts/GameContext.tsx'
import { formatNumber } from '../../hooks/useCountdown.ts'
import { devGiveResources, getPlayerSP, getActiveBuffs } from '../../services/api.ts'
import type { ResourcesResponse } from '../../types'

// Warehouse fill color: green < 80%, yellow 80-95%, red > 95%
function getWarehouseColorClass(resources: ResourcesResponse): string {
  const totalStored = resources.warehouse_metal + resources.warehouse_he3 + resources.warehouse_gold
  const cap = resources.warehouse_capacity
  if (cap <= 0) return ''
  const pct = totalStored / cap
  if (pct > 0.95) return 'hud-storage-critical'
  if (pct > 0.80) return 'hud-storage-warning'
  return 'hud-storage-ok'
}

interface ActiveBuffData {
  buffs: { buff_type: string; buff_value: number; expires_at: string }[]
  protection_until: string | null
}

export default function ResourceHUD() {
  const { resources, collect, refreshResources } = useResources()
  const { state } = useGameContext()
  const [collecting, setCollecting] = useState(false)
  const [givingResources, setGivingResources] = useState(false)
  const [sp, setSP] = useState<{ space_points: number; max_space_points: number } | null>(null)
  const [buffData, setBuffData] = useState<ActiveBuffData | null>(null)

  useEffect(() => {
    getPlayerSP().then(setSP).catch(() => {})
    getActiveBuffs().then(setBuffData).catch(() => {})
    const interval = setInterval(() => {
      getPlayerSP().then(setSP).catch(() => {})
      getActiveBuffs().then(setBuffData).catch(() => {})
    }, 30000)
    return () => clearInterval(interval)
  }, [])

  const civicCenter = state.buildings.find(b => b.type_name === 'civic_center')
  const ccLevel = civicCenter?.level ?? 0
  const isDevMode = localStorage.getItem('dev_mode') === 'true'

  async function handleCollect() {
    setCollecting(true)
    await collect()
    setCollecting(false)
  }

  async function handleDevGiveResources() {
    setGivingResources(true)
    try {
      await devGiveResources(500000, 500000, 500000)
      await refreshResources()
    } catch {
      // dev endpoint failed
    } finally {
      setGivingResources(false)
    }
  }

  // Total pending = what would be collected (warehouse + production since last update)
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

            <div className={`hud-storage ${getWarehouseColorClass(resources)}`}>
              Storage Cap: {formatNumber(resources.storage_capacity)}
            </div>

            {sp && (
              <div className="hud-sp">
                <span className="hud-sp-label">SP</span>
                <span className="hud-sp-value">{sp.space_points}/{sp.max_space_points}</span>
              </div>
            )}

            {buffData?.protection_until && (
              <div className="hud-truce" title={`Protected until ${new Date(buffData.protection_until).toLocaleString()}`}>
                Truce Active
              </div>
            )}

            {buffData && buffData.buffs.length > 0 && (
              <div className="hud-buffs">
                {buffData.buffs.map(b => (
                  <span key={b.buff_type} className="hud-buff-badge" title={`Expires: ${new Date(b.expires_at).toLocaleString()}`}>
                    {b.buff_type.replace(/_/g, ' ')}
                  </span>
                ))}
              </div>
            )}
          </>
        ) : (
          <span style={{ color: 'var(--text-dim)' }}>Loading...</span>
        )}
      </div>

      <button
        className="hud-collect-btn"
        onClick={handleCollect}
        disabled={collecting || pending <= 0}
        title={`Collect ${formatNumber(pending)} resources from warehouse`}
      >
        {collecting
          ? 'Collecting...'
          : `Collect +${formatNumber(pending)}`
        }
      </button>

      {isDevMode && (
        <button
          className="hud-dev-btn"
          onClick={handleDevGiveResources}
          disabled={givingResources}
          title="DEV: +500k Metal, He3, Gold"
        >
          {givingResources ? '...' : '+500k'}
        </button>
      )}
    </div>
  )
}
