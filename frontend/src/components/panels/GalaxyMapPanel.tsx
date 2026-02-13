import { useState, useEffect } from 'react'
import { createPortal } from 'react-dom'
import { useGalaxy } from '../../hooks/useGalaxy'
import LoadingButton from '../common/LoadingButton'
import type { Fleet } from '../../types'
import '../../styles/common.css'

interface GalaxyMapPanelProps {
  onClose: () => void
}

export default function GalaxyMapPanel({ onClose }: GalaxyMapPanelProps) {
  const { zones, selectedZone, loading, attacking, selectZone, attack, refresh } = useGalaxy()
  const [selectedFleets, setSelectedFleets] = useState<string[]>([])
  const [fleets, setFleets] = useState<Fleet[]>([])
  const [myCorpId, setMyCorpId] = useState<string | null>(null)
  const [lastResult, setLastResult] = useState<{ result: string; rounds: number; conquered: boolean } | null>(null)

  // ESC to close
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  // Fetch fleets and corp data from bridge
  useEffect(() => {
    if (window.__gameActions?._ready) {
      const fleetState = window.__gameActions.getFleetState()
      setFleets(fleetState.fleets)

      const corpState = window.__gameActions.getCorpState()
      setMyCorpId(corpState.corp?.id || null)
    }
  }, [])

  const handleAttack = async () => {
    if (!selectedZone || selectedFleets.length === 0) return

    const response = await attack(selectedZone.rbp_planet_id, { fleet_ids: selectedFleets })
    if (response) {
      setLastResult({
        result: response.result,
        rounds: response.total_rounds,
        conquered: response.conquered,
      })
      setSelectedFleets([])
    }
  }

  const toggleFleet = (fleetId: string) => {
    setSelectedFleets((prev) =>
      prev.includes(fleetId) ? prev.filter((id) => id !== fleetId) : [...prev, fleetId],
    )
  }

  const getBorderColor = (zone: typeof zones[0]) => {
    if (zone.is_protected) return '#4a90d9'
    if (zone.controlling_corp_id === myCorpId) return '#4ade80'
    if (zone.controlling_corp_id) return '#f87171'
    return '#666'
  }

  return createPortal(
    <div
      className="chat-backdrop"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div className="panel" style={{ width: '900px', maxWidth: '95vw', maxHeight: '85vh', display: 'flex', flexDirection: 'column', background: 'var(--bg-panel)', border: '1px solid var(--border-glow)', borderRadius: '8px' }}>
        <div className="panel-header">
          <h2>Galaxy Map</h2>
          <div style={{ display: 'flex', gap: '8px' }}>
            <LoadingButton
              className="btn btn-small btn-secondary"
              onClick={refresh}
              loading={loading}
            >
              Refresh
            </LoadingButton>
            <button className="p2-modal-close" onClick={onClose}>
              X
            </button>
          </div>
        </div>

        <div className="panel-content" style={{ flex: 1, overflowY: 'auto', padding: '20px' }}>
          {loading && zones.length === 0 && <div className="inline-spinner" />}

          {!loading && zones.length === 0 && (
            <p style={{ color: '#888' }}>No galaxy zones available.</p>
          )}

          {zones.length > 0 && (
            <>
              <div
                style={{
                  display: 'grid',
                  gridTemplateColumns: 'repeat(7, 1fr)',
                  gap: '4px',
                  marginBottom: '20px',
                }}
              >
                {zones.map((zone) => (
                  <div
                    key={`${zone.zone_x}-${zone.zone_y}`}
                    onClick={() => selectZone(zone)}
                    style={{
                      padding: '8px',
                      border: `2px solid ${getBorderColor(zone)}`,
                      borderRadius: '4px',
                      cursor: 'pointer',
                      textAlign: 'center',
                      fontSize: '0.75rem',
                      background:
                        selectedZone?.zone_x === zone.zone_x &&
                        selectedZone?.zone_y === zone.zone_y
                          ? 'rgba(74, 144, 217, 0.2)'
                          : 'rgba(0,0,0,0.3)',
                      transition: 'all 0.15s',
                      position: 'relative',
                    }}
                  >
                    {zone.is_protected && (
                      <div
                        style={{
                          position: 'absolute',
                          inset: 0,
                          background: 'rgba(74, 144, 217, 0.1)',
                          animation: 'shimmer 2s infinite',
                          borderRadius: '2px',
                        }}
                      />
                    )}
                    <div style={{ fontWeight: 'bold', marginBottom: '2px' }}>
                      {zone.rbp_name}
                    </div>
                    <div style={{ color: '#aaa', fontSize: '0.7rem' }}>
                      Lvl {zone.rbp_level}
                    </div>
                    <div
                      style={{
                        marginTop: '4px',
                        fontSize: '0.65rem',
                        color: zone.controlling_corp_id ? '#4ade80' : '#999',
                      }}
                    >
                      {zone.controlling_corp_tag || 'Unclaimed'}
                    </div>
                  </div>
                ))}
              </div>

              {selectedZone && (
                <div
                  style={{
                    padding: '16px',
                    background: 'rgba(0,0,0,0.3)',
                    border: '1px solid #333',
                    borderRadius: '4px',
                  }}
                >
                  <h3 style={{ marginBottom: '12px' }}>RBP Details</h3>
                  <div style={{ marginBottom: '8px' }}>
                    <strong>Name:</strong> {selectedZone.rbp_name}
                  </div>
                  <div style={{ marginBottom: '8px' }}>
                    <strong>Level:</strong> {selectedZone.rbp_level}
                  </div>
                  <div style={{ marginBottom: '8px' }}>
                    <strong>Position:</strong> ({selectedZone.zone_x}, {selectedZone.zone_y})
                  </div>
                  <div style={{ marginBottom: '8px' }}>
                    <strong>Controlling Corp:</strong>{' '}
                    {selectedZone.controlling_corp_name
                      ? `[${selectedZone.controlling_corp_tag}] ${selectedZone.controlling_corp_name}`
                      : 'None'}
                  </div>
                  {selectedZone.is_protected && selectedZone.protection_until && (
                    <div style={{ marginBottom: '8px', color: '#4a90d9' }}>
                      <strong>Protected Until:</strong>{' '}
                      {new Date(selectedZone.protection_until).toLocaleString()}
                    </div>
                  )}

                  {lastResult && (
                    <div style={{
                      padding: '12px',
                      marginTop: '12px',
                      background: lastResult.result === 'attacker_win' ? 'rgba(74, 222, 128, 0.15)' : 'rgba(248, 113, 113, 0.15)',
                      border: `1px solid ${lastResult.result === 'attacker_win' ? '#4ade80' : '#f87171'}`,
                      borderRadius: '4px',
                    }}>
                      <div style={{ fontWeight: 'bold', marginBottom: '4px', color: lastResult.result === 'attacker_win' ? '#4ade80' : '#f87171' }}>
                        {lastResult.result === 'attacker_win' ? 'Victory!' : lastResult.result === 'defender_win' ? 'Defeat' : 'Draw'}
                      </div>
                      <div style={{ fontSize: '0.85rem', color: '#aaa' }}>
                        Combat resolved in {lastResult.rounds} rounds
                        {lastResult.conquered && <span style={{ color: '#4ade80', marginLeft: '8px' }}>RBP Conquered!</span>}
                      </div>
                      <button
                        style={{ marginTop: '8px', padding: '4px 12px', background: 'rgba(0,0,0,0.3)', border: '1px solid #555', borderRadius: '4px', color: '#aaa', cursor: 'pointer' }}
                        onClick={() => setLastResult(null)}
                      >
                        Dismiss
                      </button>
                    </div>
                  )}

                  {!myCorpId && (
                    <p style={{ color: '#f87171', fontSize: '0.85rem', marginTop: '12px' }}>
                      You must be in a corp to attack RBPs.
                    </p>
                  )}

                  {myCorpId && !selectedZone.is_protected && (
                    <div style={{ marginTop: '16px' }}>
                      <h4 style={{ marginBottom: '8px' }}>Select Fleets to Attack</h4>
                      <div style={{ maxHeight: '150px', overflowY: 'auto' }}>
                        {fleets.length === 0 && (
                          <p style={{ color: '#888', fontSize: '0.85rem' }}>
                            No fleets available. Create fleets first.
                          </p>
                        )}
                        {fleets.map((fleet) => (
                          <label
                            key={fleet.id}
                            style={{
                              display: 'block',
                              padding: '8px',
                              marginBottom: '4px',
                              background: selectedFleets.includes(fleet.id)
                                ? 'rgba(74, 144, 217, 0.2)'
                                : 'rgba(0,0,0,0.2)',
                              border: '1px solid #333',
                              borderRadius: '4px',
                              cursor: 'pointer',
                            }}
                          >
                            <input
                              type="checkbox"
                              checked={selectedFleets.includes(fleet.id)}
                              onChange={() => toggleFleet(fleet.id)}
                              style={{ marginRight: '8px' }}
                            />
                            {fleet.name} ({fleet.formation})
                          </label>
                        ))}
                      </div>
                      <LoadingButton
                        className="btn btn-primary"
                        onClick={handleAttack}
                        loading={attacking}
                        disabled={selectedFleets.length === 0}
                        style={{ marginTop: '12px' }}
                      >
                        Attack RBP
                      </LoadingButton>
                    </div>
                  )}

                  {myCorpId && selectedZone.is_protected && (
                    <p style={{ color: '#4a90d9', fontSize: '0.85rem', marginTop: '12px' }}>
                      This RBP is under protection and cannot be attacked yet.
                    </p>
                  )}
                </div>
              )}
            </>
          )}
        </div>
      </div>

      <style>{`
        @keyframes shimmer {
          0%, 100% { opacity: 0.3; }
          50% { opacity: 0.6; }
        }
      `}</style>
    </div>,
    document.body,
  )
}
