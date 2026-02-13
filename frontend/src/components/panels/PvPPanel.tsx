import { useState } from 'react'
import { usePvP } from '../../hooks/usePvP.ts'
import { useFleets } from '../../hooks/useFleets.ts'
import LoadingButton from '../common/LoadingButton.tsx'
import '../../styles/pvp.css'
import '../../styles/common.css'

function formatCountdown(targetDate: string): string {
  const diff = Math.max(0, Math.floor((new Date(targetDate).getTime() - Date.now()) / 1000))
  if (diff === 0) return 'Arrived'
  const m = Math.floor(diff / 60)
  const s = diff % 60
  return `${m}m ${s}s`
}

export default function PvPPanel() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedTarget, setSelectedTarget] = useState<string | null>(null)
  const [selectedFleets, setSelectedFleets] = useState<string[]>([])
  const {
    searchResults, searching, attacking, lastDispatch, pendingAttacks,
    incoming, playerSP, cooldownSeconds,
    search, attack, cancel, clearDispatch,
  } = usePvP()
  const { fleets } = useFleets()

  const stationedFleets = fleets.filter((f) => f.status === 'stationed')
  const travelingAttacks = pendingAttacks.filter(a => a.status === 'traveling')
  const resolvedAttacks = pendingAttacks.filter(a => a.status === 'resolved').slice(0, 5)

  const handleSearch = () => search(searchQuery)

  const handleAttack = async () => {
    if (!selectedTarget || selectedFleets.length === 0) return
    await attack(selectedTarget, selectedFleets)
    setSelectedFleets([])
  }

  const toggleFleet = (fleetId: string) => {
    setSelectedFleets(prev =>
      prev.includes(fleetId) ? prev.filter(id => id !== fleetId) : [...prev, fleetId]
    )
  }

  const selectedTargetData = searchResults.find(r => r.planet_id === selectedTarget)

  return (
    <div className="panel pvp-panel">
      <div className="panel-header">
        <h2>PvP Combat</h2>
        <div className="pvp-header-info">
          {playerSP && (
            <span className="pvp-sp-badge">SP: {playerSP.space_points}/{playerSP.max_space_points}</span>
          )}
        </div>
      </div>

      <div className="panel-content">

        {/* Incoming Attacks (Radar) */}
        {incoming && incoming.incoming_attacks.length > 0 && (
          <div className="pvp-section pvp-incoming">
            <h3>Incoming Attacks (Radar Lv{incoming.radar_level})</h3>
            <div className="pvp-incoming-list">
              {incoming.incoming_attacks.map(inc => (
                <div key={inc.id} className="pvp-incoming-card">
                  <span className="pvp-incoming-eta">ETA: {formatCountdown(inc.arrival_at)}</span>
                  {inc.origin_x !== undefined && (
                    <span className="pvp-incoming-origin">From: ({inc.origin_x}, {inc.origin_y})</span>
                  )}
                  {inc.fleet_count !== undefined && (
                    <span className="pvp-incoming-strength">{inc.fleet_count} fleet(s)</span>
                  )}
                  {inc.attacker_name && (
                    <span className="pvp-incoming-attacker">{inc.attacker_name}</span>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Active Attacks in Transit */}
        {travelingAttacks.length > 0 && (
          <div className="pvp-section">
            <h3>Fleets in Transit ({travelingAttacks.length})</h3>
            <div className="pvp-pending-list">
              {travelingAttacks.map(pa => (
                <div key={pa.id} className="pvp-pending-card">
                  <div className="pvp-pending-info">
                    <span className="pvp-pending-target">{pa.planet_name} ({pa.defender_name})</span>
                    <span className="pvp-pending-eta">ETA: {formatCountdown(pa.arrival_at)}</span>
                    <span className="pvp-pending-fleets">{pa.fleet_ids.length} fleet(s)</span>
                  </div>
                  <button className="btn btn-secondary btn-sm" onClick={() => cancel(pa.id)}>
                    Recall
                  </button>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Dispatch confirmation */}
        {lastDispatch && (
          <div className="pvp-section pvp-dispatch-result">
            <h3>Fleets Dispatched!</h3>
            <p>{lastDispatch.fleets_dispatched} fleet(s) en route.</p>
            <p>Travel time: {Math.floor(lastDispatch.travel_seconds / 60)}m {lastDispatch.travel_seconds % 60}s</p>
            <p>ETA: {formatCountdown(lastDispatch.arrival_at)}</p>
            <button className="btn btn-secondary" onClick={clearDispatch}>OK</button>
          </div>
        )}

        {/* Search Section */}
        {!selectedTarget && !lastDispatch && (
          <div className="pvp-section">
            <h3>Search for Targets</h3>
            <div className="pvp-search">
              <input
                type="text"
                className="pvp-search-input"
                placeholder="Search by planet or player name..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              />
              <LoadingButton className="btn btn-primary" onClick={handleSearch} loading={searching}>
                Search
              </LoadingButton>
            </div>

            {searchResults.length > 0 && (
              <div className="pvp-results">
                <h4>Search Results ({searchResults.length})</h4>
                <div className="pvp-results-list">
                  {searchResults.map(result => (
                    <div
                      key={result.planet_id}
                      className="pvp-result-card"
                      onClick={() => setSelectedTarget(result.planet_id)}
                    >
                      <div className="pvp-result-planet">{result.planet_name}</div>
                      <div className="pvp-result-player">Owner: {result.player_name}</div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}

        {/* Fleet Selection Section */}
        {selectedTarget && !lastDispatch && (
          <div className="pvp-section">
            <div className="pvp-target-header">
              <h3>
                Target: {selectedTargetData?.planet_name} ({selectedTargetData?.player_name})
              </h3>
              <button className="btn btn-secondary" onClick={() => setSelectedTarget(null)}>
                Cancel
              </button>
            </div>

            <h4>Select Fleets to Attack</h4>
            {stationedFleets.length === 0 ? (
              <p className="pvp-empty">No stationed fleets available</p>
            ) : (
              <div className="pvp-fleet-list">
                {stationedFleets.map(fleet => (
                  <div
                    key={fleet.id}
                    className={`pvp-fleet-card ${selectedFleets.includes(fleet.id) ? 'selected' : ''}`}
                    onClick={() => toggleFleet(fleet.id)}
                  >
                    <div className="pvp-fleet-name">{fleet.name}</div>
                    <div className="pvp-fleet-ships">{fleet.stacks?.length || 0} stacks</div>
                  </div>
                ))}
              </div>
            )}

            <div className="pvp-attack-actions">
              <LoadingButton
                className="btn btn-danger"
                onClick={handleAttack}
                loading={attacking}
                disabled={selectedFleets.length === 0 || cooldownSeconds > 0 || (playerSP?.space_points ?? 0) < 1}
              >
                {cooldownSeconds > 0
                  ? `Cooldown ${Math.floor(cooldownSeconds / 60)}m ${cooldownSeconds % 60}s`
                  : (playerSP?.space_points ?? 0) < 1
                    ? 'No SP Available'
                    : `Dispatch ${selectedFleets.length} Fleet(s) (1 SP)`
                }
              </LoadingButton>
            </div>
          </div>
        )}

        {/* Recent Combat Results */}
        {resolvedAttacks.length > 0 && (
          <div className="pvp-section">
            <h3>Recent Attacks</h3>
            <div className="pvp-recent-list">
              {resolvedAttacks.map(pa => (
                <div key={pa.id} className="pvp-recent-card">
                  <span className="pvp-recent-target">{pa.planet_name}</span>
                  <span className="pvp-recent-status">
                    {pa.combat_report_id ? 'Resolved' : 'Pending'}
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* No radar warning */}
        {incoming && incoming.radar_level === 0 && (
          <div className="pvp-section pvp-no-radar">
            <p>Build a Radar to detect incoming attacks!</p>
          </div>
        )}
      </div>
    </div>
  )
}
