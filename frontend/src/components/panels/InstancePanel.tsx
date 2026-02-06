import { useState } from 'react'
import { useInstances } from '../../hooks/useInstances.ts'
import { useFleets } from '../../hooks/useFleets.ts'
import { formatNumber } from '../../hooks/useCountdown.ts'
import type { Instance, InstanceDetail, InstanceAttemptResponse, Fleet } from '../../types'

interface InstanceDetailModalProps {
  instance: Instance
  detail: InstanceDetail | null
  loadingDetail: boolean
  fleets: Fleet[]
  onAttempt: (instanceId: number, fleetIds: string[]) => Promise<InstanceAttemptResponse>
  onClose: () => void
}

function InstanceDetailModal({
  instance, detail, loadingDetail, fleets, onAttempt, onClose,
}: InstanceDetailModalProps) {
  const [selectedFleets, setSelectedFleets] = useState<string[]>([])
  const [attacking, setAttacking] = useState(false)
  const [result, setResult] = useState<InstanceAttemptResponse | null>(null)

  function toggleFleet(id: string) {
    if (selectedFleets.includes(id)) {
      setSelectedFleets(selectedFleets.filter(f => f !== id))
    } else if (selectedFleets.length < instance.max_fleets) {
      setSelectedFleets([...selectedFleets, id])
    }
  }

  async function handleAttack() {
    if (selectedFleets.length === 0) return
    setAttacking(true)
    try {
      const res = await onAttempt(instance.id, selectedFleets)
      setResult(res)
    } catch {
      // Error handled by hook
    } finally {
      setAttacking(false)
    }
  }

  const stationedFleets = fleets.filter(f => f.status === 'stationed' && (f.stacks?.length || 0) > 0)

  return (
    <div className="p2-modal-backdrop" onClick={e => { if (e.target === e.currentTarget) onClose() }}>
      <div className="p2-modal p2-modal-md">
        <div className="p2-modal-header">
          <span className="p2-modal-title">{instance.name}</span>
          <button className="p2-modal-close" onClick={onClose}>X</button>
        </div>
        <div className="p2-modal-body">
          {result ? (
            /* Battle Result */
            <div className="inst-result">
              <div className={`inst-result-banner ${result.result === 'attacker_win' ? 'victory' : 'defeat'}`}>
                {result.result === 'attacker_win' ? 'VICTORY' : 'DEFEAT'}
              </div>
              <div className="inst-result-stats">
                <div className="inst-result-row">
                  <span>Rounds:</span><span>{result.total_rounds}</span>
                </div>
                <div className="inst-result-row">
                  <span>EXP Gained:</span><span className="inst-exp">+{formatNumber(result.exp_gained)}</span>
                </div>
                {result.treasure_box && (
                  <>
                    <div className="p2-section-title">Rewards</div>
                    <div className="inst-result-row">
                      <span className="res-dot metal" /> Metal: +{formatNumber(result.treasure_box.resources.metal)}
                    </div>
                    <div className="inst-result-row">
                      <span className="res-dot he3" /> He3: +{formatNumber(result.treasure_box.resources.he3)}
                    </div>
                    <div className="inst-result-row">
                      <span className="res-dot gold" /> Gold: +{formatNumber(result.treasure_box.resources.gold)}
                    </div>
                    {result.treasure_box.blueprint_id && (
                      <div className="inst-result-row inst-bp-drop">
                        Blueprint Drop! ID: {result.treasure_box.blueprint_id}
                      </div>
                    )}
                  </>
                )}
                {Object.keys(result.losses.ships_destroyed).length > 0 && (
                  <>
                    <div className="p2-section-title">Losses</div>
                    {Object.entries(result.losses.ships_destroyed).map(([designId, count]) => (
                      <div key={designId} className="inst-result-row inst-loss">
                        <span>Design {designId.slice(0, 8)}...</span>
                        <span>-{count} ships</span>
                      </div>
                    ))}
                    <div className="inst-result-row inst-loss">
                      <span>He3 Consumed:</span>
                      <span>{formatNumber(result.losses.he3_consumed)}</span>
                    </div>
                  </>
                )}
              </div>
              <button className="p2-btn p2-btn-primary p2-btn-full" onClick={onClose}>Close</button>
            </div>
          ) : (
            /* Pre-battle screen */
            <>
              <div className="inst-info">
                <div className="inst-info-row">
                  <span>Level Requirement:</span>
                  <span>{instance.required_level}</span>
                </div>
                <div className="inst-info-row">
                  <span>Max Fleets:</span>
                  <span>{instance.max_fleets}</span>
                </div>
                <div className="inst-info-row">
                  <span>EXP Reward:</span>
                  <span className="inst-exp">{formatNumber(instance.exp_reward)}</span>
                </div>
              </div>

              {loadingDetail ? (
                <div className="p2-panel-loading"><div className="loading-spinner" /></div>
              ) : detail && (
                <div className="inst-enemies">
                  <div className="p2-section-title">Enemy Forces</div>
                  {detail.enemy_fleets.map((ef, i) => (
                    <div key={i} className="inst-enemy-row">
                      <span className={`inst-enemy-class ${ef.hull_class}`}>{ef.hull_class}</span>
                      <span>x{ef.ship_count}</span>
                      <span className="inst-enemy-power">~{formatNumber(ef.power_estimate)} power</span>
                    </div>
                  ))}
                </div>
              )}

              <div className="inst-fleet-select">
                <div className="p2-section-title">
                  Select Fleets ({selectedFleets.length}/{instance.max_fleets})
                </div>
                {stationedFleets.length === 0 ? (
                  <div className="p2-empty-state">No stationed fleets available. Create a fleet first.</div>
                ) : (
                  <div className="inst-fleet-list">
                    {stationedFleets.map(f => {
                      const isSelected = selectedFleets.includes(f.id)
                      const totalShips = f.stacks?.reduce((sum, s) => sum + s.ship_count, 0) || 0
                      return (
                        <button
                          key={f.id}
                          className={`inst-fleet-btn ${isSelected ? 'selected' : ''}`}
                          onClick={() => toggleFleet(f.id)}
                        >
                          <span className="inst-fleet-name">{f.name}</span>
                          <span className="inst-fleet-ships">{formatNumber(totalShips)} ships</span>
                        </button>
                      )
                    })}
                  </div>
                )}
              </div>

              <button
                className="p2-btn p2-btn-danger p2-btn-full"
                onClick={handleAttack}
                disabled={selectedFleets.length === 0 || attacking}
              >
                {attacking ? 'Attacking...' : 'ATTACK'}
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

export default function InstancePanel() {
  const { instances, loading, error, getDetail, attempt, isCompleted } = useInstances()
  const { fleets } = useFleets()
  const [selectedInstance, setSelectedInstance] = useState<Instance | null>(null)
  const [instanceDetail, setInstanceDetail] = useState<InstanceDetail | null>(null)
  const [loadingDetail, setLoadingDetail] = useState(false)

  if (loading) {
    return <div className="p2-panel-loading"><div className="loading-spinner" /><span>Loading Instances...</span></div>
  }

  async function handleSelect(inst: Instance) {
    setSelectedInstance(inst)
    setLoadingDetail(true)
    setInstanceDetail(null)
    try {
      const detail = await getDetail(inst.id)
      setInstanceDetail(detail)
    } catch {
      // Error handled
    } finally {
      setLoadingDetail(false)
    }
  }

  return (
    <div className="p2-panel">
      <div className="p2-panel-header">
        <div className="p2-panel-icon inst-icon">IN</div>
        <div>
          <div className="p2-panel-title">Instances</div>
          <div className="p2-panel-subtitle">
            Normal Instances - {instances.filter(i => isCompleted(i.id)).length}/{instances.length} completed
          </div>
        </div>
      </div>

      {error && <div className="p2-error-msg">{error}</div>}

      <div className="inst-list">
        {instances.map(inst => {
          const completed = isCompleted(inst.id)
          return (
            <div
              key={inst.id}
              className={`inst-card ${completed ? 'completed' : ''}`}
              onClick={() => handleSelect(inst)}
            >
              <div className="inst-card-number">{inst.difficulty}</div>
              <div className="inst-card-info">
                <div className="inst-card-name">{inst.name}</div>
                <div className="inst-card-req">
                  Lv {inst.required_level} | {inst.max_fleets} fleet{inst.max_fleets > 1 ? 's' : ''}
                </div>
              </div>
              <div className="inst-card-reward">
                <span className="inst-exp">+{formatNumber(inst.exp_reward)} EXP</span>
              </div>
              {completed && <div className="inst-card-check">Done</div>}
            </div>
          )
        })}
      </div>

      {selectedInstance && (
        <InstanceDetailModal
          instance={selectedInstance}
          detail={instanceDetail}
          loadingDetail={loadingDetail}
          fleets={fleets}
          onAttempt={attempt}
          onClose={() => { setSelectedInstance(null); setInstanceDetail(null) }}
        />
      )}
    </div>
  )
}
