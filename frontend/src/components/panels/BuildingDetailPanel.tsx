import { useGameContext, CATEGORY_COLORS, BUILDING_ABBREVIATIONS } from '../../contexts/GameContext.tsx'
import { useBuildings } from '../../hooks/useBuildings.ts'
import { useCountdown, formatNumber, formatDuration } from '../../hooks/useCountdown.ts'
import { useResources } from '../../hooks/useResources.ts'
import { useEffect, useState } from 'react'
import { cancelUpgrade } from '../../services/api.ts'

// GO2 cost formula: Cost(N) = BaseCost * CostMultiplier^(N-1)
function calcUpgradeCost(bt: { base_cost_metal: number; base_cost_he3: number; base_cost_gold: number; base_time_seconds: number; cost_multiplier: number; time_multiplier: number }, targetLevel: number) {
  const costMult = Math.pow(bt.cost_multiplier, targetLevel - 1)
  const timeMult = Math.pow(bt.time_multiplier, targetLevel - 1)
  return {
    metal: Math.round(bt.base_cost_metal * costMult),
    he3: Math.round(bt.base_cost_he3 * costMult),
    gold: Math.round(bt.base_cost_gold * costMult),
    time_seconds: Math.round(bt.base_time_seconds * timeMult),
  }
}

// GO2 production formula: Production(N) = BaseProd * ProdMult^(N-1)
function calcProduction(bt: { base_production_per_hour: number; production_multiplier: number }, level: number): number {
  if (!bt.base_production_per_hour || level < 1) return 0
  return Math.round(bt.base_production_per_hour * Math.pow(bt.production_multiplier, level - 1))
}

const RESOURCE_LABEL: Record<string, string> = {
  metal_collector: 'Metal',
  he3_extractor: 'He3',
  residential_area: 'Gold',
}

export default function BuildingDetailPanel() {
  const { state, closeDetailPanel } = useGameContext()
  const { upgrade, refreshBuildings } = useBuildings()
  const { resources } = useResources()
  const building = state.selectedBuilding
  const countdown = useCountdown(building?.upgrade_finish_at ?? null)
  const [cancelling, setCancelling] = useState(false)

  const isOpen = state.showDetailPanel && building !== null

  // Close on Escape
  useEffect(() => {
    if (!isOpen) return
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') closeDetailPanel()
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [isOpen, closeDetailPanel])

  // Find the building type data from backend
  const bt = building ? state.buildingTypes.find(t => t.name === building.type_name) : null

  const level = building?.level ?? 0
  const targetLevel = level + 1

  const upgradeCost = bt
    ? calcUpgradeCost(bt, targetLevel)
    : { metal: 0, he3: 0, gold: 0, time_seconds: 0 }

  const hasMetal = resources ? resources.metal >= upgradeCost.metal : false
  const hasHe3 = resources ? resources.he3 >= upgradeCost.he3 : false
  const hasGold = resources ? resources.gold >= upgradeCost.gold : false
  const canAfford = hasMetal && hasHe3 && hasGold
  const isMaxLevel = building ? building.level >= building.max_level : false
  const isUpgrading = building?.is_upgrading ?? false

  const color = building ? CATEGORY_COLORS[building.category] || '#888' : '#888'
  const abbr = building ? BUILDING_ABBREVIATIONS[building.type_name] || '??' : ''

  let buttonClass = 'bdp-upgrade-btn '
  let buttonText = ''
  if (isMaxLevel) {
    buttonClass += 'max-level'
    buttonText = 'MAX LEVEL'
  } else if (isUpgrading) {
    buttonClass += 'upgrading'
    buttonText = `UPGRADING... ${countdown}`
  } else if (canAfford) {
    buttonClass += 'can-upgrade'
    buttonText = `UPGRADE TO LV ${level + 1}`
  } else {
    buttonClass += 'disabled'
    buttonText = `UPGRADE TO LV ${level + 1}`
  }

  async function handleUpgrade() {
    if (!building || isMaxLevel || isUpgrading || !canAfford) return
    await upgrade(building.id)
  }

  async function handleCancel() {
    if (!building || !state.currentPlanet) return
    setCancelling(true)
    try {
      await cancelUpgrade(state.currentPlanet.id, building.id)
      await refreshBuildings()
    } catch {
      // Error displayed via context
    } finally {
      setCancelling(false)
    }
  }

  const hasProduction = bt && bt.base_production_per_hour > 0
  const resourceLabel = building ? RESOURCE_LABEL[building.type_name] || '' : ''

  return (
    <div className={`building-detail-panel ${isOpen ? 'open' : ''}`}>
      {building && (
        <>
          <div className="bdp-header">
            <button className="bdp-close" onClick={closeDetailPanel}>
              X
            </button>
            <span className={`bdp-category-badge ${building.category}`}>
              {building.category}
            </span>

            <div
              className="bdp-icon"
              style={{ background: color, color: '#000' }}
            >
              {abbr}
            </div>

            <div className="bdp-name">{building.display_name}</div>
            <div className="bdp-level">
              Level {building.level} / {building.max_level}
            </div>
            <div className="bdp-level-bar">
              <div
                className="bdp-level-bar-fill"
                style={{ width: `${(building.level / building.max_level) * 100}%` }}
              />
            </div>
          </div>

          {/* Production section for resource buildings */}
          {hasProduction && bt && (
            <div className="bdp-section">
              <div className="bdp-section-title">Production</div>
              <div className="bdp-production-current">
                Current: {formatNumber(calcProduction(bt, building.level))} {resourceLabel}/hr
              </div>
              {!isMaxLevel && (
                <div className="bdp-production-next">
                  Next Lv: {formatNumber(calcProduction(bt, building.level + 1))} /hr
                  {calcProduction(bt, building.level) > 0 && (
                    <> (+{((calcProduction(bt, building.level + 1) / calcProduction(bt, building.level) - 1) * 100).toFixed(1)}%)</>
                  )}
                </div>
              )}
            </div>
          )}

          {/* Upgrade cost section */}
          {!isMaxLevel && (
            <div className="bdp-section">
              <div className="bdp-section-title">Upgrade to Level {level + 1}</div>

              <div className="bdp-cost-row">
                <div className="bdp-cost-left">
                  <span className="bdp-cost-icon metal" />
                  <span className="bdp-cost-name">Metal</span>
                </div>
                <span className="bdp-cost-amount">{formatNumber(upgradeCost.metal)}</span>
                <span className={`bdp-cost-status ${hasMetal ? 'ok' : 'insufficient'}`}>
                  {hasMetal ? '\u2713' : '\u2717'}
                </span>
              </div>

              <div className="bdp-cost-row">
                <div className="bdp-cost-left">
                  <span className="bdp-cost-icon he3" />
                  <span className="bdp-cost-name">He3</span>
                </div>
                <span className="bdp-cost-amount">{formatNumber(upgradeCost.he3)}</span>
                <span className={`bdp-cost-status ${hasHe3 ? 'ok' : 'insufficient'}`}>
                  {hasHe3 ? '\u2713' : '\u2717'}
                </span>
              </div>

              <div className="bdp-cost-row">
                <div className="bdp-cost-left">
                  <span className="bdp-cost-icon gold" />
                  <span className="bdp-cost-name">Gold</span>
                </div>
                <span className="bdp-cost-amount">{formatNumber(upgradeCost.gold)}</span>
                <span className={`bdp-cost-status ${hasGold ? 'ok' : 'insufficient'}`}>
                  {hasGold ? '\u2713' : '\u2717'}
                </span>
              </div>

              <div className="bdp-time-row">
                <span>Time:</span>
                <span className="bdp-time-value">{formatDuration(upgradeCost.time_seconds)}</span>
              </div>
            </div>
          )}

          <div className="bdp-section">
            {isUpgrading ? (
              <button
                className="bdp-upgrade-btn upgrading-cancel"
                onClick={handleCancel}
                disabled={cancelling}
              >
                {cancelling ? 'CANCELLING...' : `CANCEL UPGRADE (${countdown})`}
              </button>
            ) : (
              <button
                className={buttonClass}
                onClick={handleUpgrade}
                disabled={isMaxLevel || !canAfford}
              >
                {buttonText}
              </button>
            )}
          </div>
        </>
      )}
    </div>
  )
}
