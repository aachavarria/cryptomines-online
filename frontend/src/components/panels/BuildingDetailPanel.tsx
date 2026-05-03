import { useGameContext } from '../../contexts/GameContext.tsx'
import { useBuildings } from '../../hooks/useBuildings.ts'
import { useCountdown, formatNumber, formatDuration } from '../../hooks/useCountdown.ts'
import { useResources } from '../../hooks/useResources.ts'
import { useEffect, useState } from 'react'
import { cancelUpgrade } from '../../services/api.ts'
import {
  Building2,
  Factory,
  Hammer,
  Home,
  Rocket,
  ShieldHalf,
  Target,
  X,
  CircleCheck,
  CircleX,
  type LucideIcon,
} from 'lucide-react'

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

// Map a building category to its design-system badge variant.
const CATEGORY_BADGE: Record<string, string> = {
  resource: 'ds-badge--success',
  core: 'ds-badge--teal',
  military: 'ds-badge--danger',
  defense: 'ds-badge--orange',
  space: 'ds-badge--orange',
}

// Map a category (or building type) to a lucide icon.
function pickIcon(category: string, typeName: string): LucideIcon {
  if (typeName === 'metal_collector' || typeName === 'he3_extractor') return Factory
  if (typeName === 'residential_area') return Home
  switch (category) {
    case 'resource':
      return Factory
    case 'core':
      return Hammer
    case 'military':
      return Target
    case 'defense':
      return ShieldHalf
    case 'space':
      return Rocket
    default:
      return Building2
  }
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

  const buttonText = isUpgrading
    ? `UPGRADING... ${countdown}`
    : `UPGRADE TO LV ${level + 1}`

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

  const Icon = building ? pickIcon(building.category, building.type_name) : Building2
  const badgeClass = building ? CATEGORY_BADGE[building.category] || 'ds-badge--neutral' : 'ds-badge--neutral'

  return (
    <div className={`building-detail-panel ds-panel--flush ${isOpen ? 'open' : ''}`}>
      {building && (
        <>
          <div className="bdp-header">
            <button
              className="bdp-close ds-btn-icon"
              onClick={closeDetailPanel}
              aria-label="Close"
            >
              <X size={16} />
            </button>

            <div className="bdp-icon">
              <Icon strokeWidth={1.75} />
            </div>

            <span className={`ds-badge ${badgeClass}`}>{building.category}</span>

            <h2 className="bdp-name ds-h2">{building.display_name}</h2>
            <div className="bdp-level ds-text-muted">
              Level <span className="ds-mono">{building.level}</span> · <span className="ds-mono">{building.level}/{building.max_level}</span>
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
                Current: <span className="ds-mono">{formatNumber(calcProduction(bt, building.level))}</span> {resourceLabel}/hr
              </div>
              {!isMaxLevel && (
                <div className="bdp-production-next">
                  Next Lv: <span className="ds-mono">{formatNumber(calcProduction(bt, building.level + 1))}</span> /hr
                  {calcProduction(bt, building.level) > 0 && (
                    <> (+<span className="ds-mono">{((calcProduction(bt, building.level + 1) / calcProduction(bt, building.level) - 1) * 100).toFixed(1)}</span>%)</>
                  )}
                </div>
              )}
            </div>
          )}

          {/* Upgrade cost section */}
          {!isMaxLevel && (
            <div className="bdp-section">
              <div className="bdp-section-title">Upgrade to Level {level + 1}</div>

              <CostRow
                label="Metal"
                amount={upgradeCost.metal}
                ok={hasMetal}
                dotClass="ds-resource-dot--metal"
              />
              <CostRow
                label="He3"
                amount={upgradeCost.he3}
                ok={hasHe3}
                dotClass="ds-resource-dot--he3"
              />
              <CostRow
                label="Gold"
                amount={upgradeCost.gold}
                ok={hasGold}
                dotClass="ds-resource-dot--gold"
              />

              <div className="bdp-time-row">
                <span>Time</span>
                <span className="bdp-time-value ds-mono">{formatDuration(upgradeCost.time_seconds)}</span>
              </div>
            </div>
          )}

          <div className="bdp-section">
            {isMaxLevel ? (
              <span className="ds-badge ds-badge--orange ds-btn--block" style={{ justifyContent: 'center', padding: '10px 16px' }}>
                MAX LEVEL
              </span>
            ) : isUpgrading ? (
              <button
                className="bdp-upgrade-btn ds-btn ds-btn--block upgrading-cancel"
                onClick={handleCancel}
                disabled={cancelling}
              >
                {cancelling ? 'CANCELLING...' : `CANCEL UPGRADE (${countdown})`}
              </button>
            ) : (
              <button
                className="bdp-upgrade-btn ds-btn-secondary ds-btn--block"
                onClick={handleUpgrade}
                disabled={!canAfford}
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

function CostRow({
  label,
  amount,
  ok,
  dotClass,
}: {
  label: string
  amount: number
  ok: boolean
  dotClass: string
}) {
  return (
    <div className="bdp-cost-row">
      <div className="bdp-cost-left">
        <span className={`ds-resource-dot ${dotClass}`} />
        <span className="bdp-cost-name">{label}</span>
      </div>
      <span className="bdp-cost-amount ds-mono">{formatNumber(amount)}</span>
      <span className={`bdp-cost-status ${ok ? 'ok' : 'insufficient'}`} aria-label={ok ? 'sufficient' : 'insufficient'}>
        {ok ? <CircleCheck size={16} /> : <CircleX size={16} />}
      </span>
    </div>
  )
}
