import { useState } from 'react'
import { usePvP } from '../../hooks/usePvP.ts'
import { useFleets } from '../../hooks/useFleets.ts'
import { formatNumber } from '../../hooks/useCountdown.ts'
import LoadingButton from '../common/LoadingButton.tsx'
import '../../styles/pvp.css'
import '../../styles/common.css'

export default function PvPPanel() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedTarget, setSelectedTarget] = useState<string | null>(null)
  const [selectedFleets, setSelectedFleets] = useState<string[]>([])
  const { searchResults, searching, attacking, lastBattle, search, attack, clearBattle } = usePvP()
  const { fleets } = useFleets()

  const stationedFleets = fleets.filter((f) => f.status === 'stationed')

  const handleSearch = () => {
    search(searchQuery)
  }

  const handleAttack = async () => {
    if (!selectedTarget || selectedFleets.length === 0) return
    await attack(selectedTarget, selectedFleets)
    setSelectedFleets([])
  }

  const toggleFleet = (fleetId: string) => {
    setSelectedFleets((prev) =>
      prev.includes(fleetId) ? prev.filter((id) => id !== fleetId) : [...prev, fleetId]
    )
  }

  const selectedTargetData = searchResults.find((r) => r.planet_id === selectedTarget)

  return (
    <div className="panel pvp-panel">
      <div className="panel-header">
        <h2>PvP Combat</h2>
        <p className="panel-subtitle">Attack other players for resources</p>
      </div>

      <div className="panel-content">
        {/* Search Section */}
        {!selectedTarget && !lastBattle && (
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
                  {searchResults.map((result) => (
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
        {selectedTarget && !lastBattle && (
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
                {stationedFleets.map((fleet) => (
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
                disabled={selectedFleets.length === 0}
              >
                Attack with {selectedFleets.length} Fleet(s)
              </LoadingButton>
            </div>
          </div>
        )}

        {/* Battle Results Section */}
        {lastBattle && (
          <div className="pvp-section">
            <div className="pvp-result-header">
              <h3>Battle Report</h3>
              <button className="btn btn-secondary" onClick={clearBattle}>
                New Attack
              </button>
            </div>

            <div className={`pvp-result-outcome ${lastBattle.result}`}>
              {lastBattle.result === 'attacker_win' && '🎉 Victory!'}
              {lastBattle.result === 'defender_win' && '💀 Defeat'}
              {lastBattle.result === 'draw' && '🤝 Draw'}
            </div>

            <div className="pvp-result-stats">
              <div className="pvp-stat">
                <span className="pvp-stat-label">Rounds:</span>
                <span className="pvp-stat-value">{lastBattle.total_rounds}</span>
              </div>
            </div>

            {lastBattle.loot_gained && (
              <div className="pvp-loot">
                <h4>Loot Gained</h4>
                <div className="pvp-loot-items">
                  <span>🪙 {formatNumber(lastBattle.loot_gained.metal)} M</span>
                  <span>⚡ {formatNumber(lastBattle.loot_gained.he3)} H3</span>
                  <span>💰 {formatNumber(lastBattle.loot_gained.gold)} G</span>
                </div>
              </div>
            )}

            <div className="pvp-losses">
              <div className="pvp-losses-col">
                <h5>Your Losses</h5>
                <p>Ships: {Object.keys(lastBattle.attacker_losses.ships_destroyed).length}</p>
                <p>He3: {formatNumber(lastBattle.attacker_losses.he3_consumed)}</p>
              </div>
              <div className="pvp-losses-col">
                <h5>Enemy Losses</h5>
                <p>Ships: {Object.keys(lastBattle.defender_losses.ships_destroyed).length}</p>
                <p>He3: {formatNumber(lastBattle.defender_losses.he3_consumed)}</p>
              </div>
            </div>

            <div className="pvp-cooldown-notice">
              ⏱️ 5-minute attack cooldown active for this target
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
