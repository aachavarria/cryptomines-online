import { useState, type CSSProperties } from 'react'
import { Sword, Crosshair, Radar, Rocket, Search, X } from 'lucide-react'
import { usePvP } from '../../hooks/usePvP.ts'
import { useFleets } from '../../hooks/useFleets.ts'
import LoadingButton from '../common/LoadingButton.tsx'

function formatCountdown(targetDate: string): string {
  const diff = Math.max(0, Math.floor((new Date(targetDate).getTime() - Date.now()) / 1000))
  if (diff === 0) return 'Arrived'
  const m = Math.floor(diff / 60)
  const s = diff % 60
  return `${m}m ${s}s`
}

const styles: Record<string, CSSProperties> = {
  panel: {
    maxWidth: 900,
    margin: '0 auto',
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-4)',
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 'var(--sp-3)',
  },
  title: {
    display: 'inline-flex',
    alignItems: 'center',
    gap: 'var(--sp-2)',
    margin: 0,
  },
  section: {
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-3)',
  },
  sectionHeader: {
    display: 'flex',
    alignItems: 'center',
    gap: 'var(--sp-2)',
    fontSize: 'var(--fs-h3)',
    fontWeight: 600,
    color: 'var(--ds-text)',
    margin: 0,
  },
  searchRow: {
    display: 'flex',
    gap: 'var(--sp-2)',
    alignItems: 'center',
  },
  resultsGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))',
    gap: 'var(--sp-2)',
  },
  fleetGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))',
    gap: 'var(--sp-2)',
  },
  pendingList: {
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-2)',
  },
  pendingRow: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 'var(--sp-3)',
  },
  incomingMeta: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: 'var(--sp-3)',
    fontSize: 'var(--fs-sm)',
    color: 'var(--ds-text-muted)',
  },
  attackActions: {
    display: 'flex',
    justifyContent: 'flex-end',
    gap: 'var(--sp-2)',
    paddingTop: 'var(--sp-3)',
    borderTop: '1px solid var(--ds-border)',
  },
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
    <div className="ds-panel" style={styles.panel}>
      <div style={styles.header}>
        <h2 className="ds-h2" style={styles.title}>
          <Sword size={20} strokeWidth={1.75} aria-hidden="true" />
          PvP Combat
        </h2>
        {playerSP && (
          <span className="ds-badge ds-badge--info">
            SP: <span className="ds-mono" style={{ marginLeft: 4 }}>
              {playerSP.space_points}/{playerSP.max_space_points}
            </span>
          </span>
        )}
      </div>

      {/* Incoming Attacks (Radar) */}
      {incoming && incoming.incoming_attacks.length > 0 && (
        <section className="ds-card" style={{ ...styles.section, borderColor: 'var(--ds-danger)' }}>
          <h3 style={styles.sectionHeader}>
            <Radar size={18} strokeWidth={1.75} aria-hidden="true" style={{ color: 'var(--ds-danger)' }} />
            Incoming Attacks (Radar Lv{incoming.radar_level})
          </h3>
          <div style={styles.pendingList}>
            {incoming.incoming_attacks.map(inc => (
              <div key={inc.id} className="ds-card" style={{ background: 'var(--ds-danger-tint)' }}>
                <div style={styles.incomingMeta}>
                  <span style={{ color: 'var(--ds-danger)', fontWeight: 600 }}>
                    ETA: <span className="ds-mono">{formatCountdown(inc.arrival_at)}</span>
                  </span>
                  {inc.origin_x !== undefined && (
                    <span>From: <span className="ds-mono">({inc.origin_x}, {inc.origin_y})</span></span>
                  )}
                  {inc.fleet_count !== undefined && (
                    <span>{inc.fleet_count} fleet(s)</span>
                  )}
                  {inc.attacker_name && (
                    <span style={{ color: 'var(--ds-owner-enemy)' }}>{inc.attacker_name}</span>
                  )}
                </div>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* Active Attacks in Transit */}
      {travelingAttacks.length > 0 && (
        <section style={styles.section}>
          <h3 style={styles.sectionHeader}>
            <Rocket size={18} strokeWidth={1.75} aria-hidden="true" />
            Fleets in Transit ({travelingAttacks.length})
          </h3>
          <div style={styles.pendingList}>
            {travelingAttacks.map(pa => (
              <div key={pa.id} className="ds-card" style={styles.pendingRow}>
                <div style={styles.incomingMeta}>
                  <span style={{ fontWeight: 600 }}>
                    {pa.planet_name} <span className="ds-text-muted">({pa.defender_name})</span>
                  </span>
                  <span>ETA: <span className="ds-mono">{formatCountdown(pa.arrival_at)}</span></span>
                  <span>{pa.fleet_ids.length} fleet(s)</span>
                </div>
                <button className="ds-btn ds-btn-ghost ds-btn--sm" onClick={() => cancel(pa.id)}>
                  Recall
                </button>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* Dispatch confirmation */}
      {lastDispatch && (
        <section className="ds-card" style={{ ...styles.section, borderColor: 'var(--ds-teal)', background: 'var(--ds-teal-tint)' }}>
          <h3 style={styles.sectionHeader}>
            <Rocket size={18} strokeWidth={1.75} aria-hidden="true" style={{ color: 'var(--ds-teal)' }} />
            Fleets Dispatched!
          </h3>
          <p style={{ margin: 0 }}>{lastDispatch.fleets_dispatched} fleet(s) en route.</p>
          <p style={{ margin: 0 }}>
            Travel time:{' '}
            <span className="ds-mono">
              {Math.floor(lastDispatch.travel_seconds / 60)}m {lastDispatch.travel_seconds % 60}s
            </span>
          </p>
          <p style={{ margin: 0 }}>
            ETA: <span className="ds-mono">{formatCountdown(lastDispatch.arrival_at)}</span>
          </p>
          <button className="ds-btn ds-btn-secondary" onClick={clearDispatch} style={{ alignSelf: 'flex-start' }}>
            OK
          </button>
        </section>
      )}

      {/* Search Section */}
      {!selectedTarget && !lastDispatch && (
        <section style={styles.section}>
          <h3 style={styles.sectionHeader}>
            <Search size={18} strokeWidth={1.75} aria-hidden="true" />
            Search for Targets
          </h3>
          <div style={styles.searchRow}>
            <input
              type="text"
              className="ds-input"
              placeholder="Search by planet or player name..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
            />
            <LoadingButton className="ds-btn ds-btn-primary" onClick={handleSearch} loading={searching}>
              Search
            </LoadingButton>
          </div>

          {searchResults.length > 0 && (
            <>
              <div className="ds-caption">Search Results ({searchResults.length})</div>
              <div style={styles.resultsGrid}>
                {searchResults.map(result => (
                  <div
                    key={result.planet_id}
                    className="ds-list-item"
                    onClick={() => setSelectedTarget(result.planet_id)}
                  >
                    <div style={{ fontWeight: 600 }}>{result.planet_name}</div>
                    <div className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)' }}>
                      Owner: <span style={{ color: 'var(--ds-owner-enemy)' }}>{result.player_name}</span>
                    </div>
                  </div>
                ))}
              </div>
            </>
          )}
        </section>
      )}

      {/* Fleet Selection Section */}
      {selectedTarget && !lastDispatch && (
        <section style={styles.section}>
          <div style={styles.header}>
            <h3 style={styles.sectionHeader}>
              <Crosshair size={18} strokeWidth={1.75} aria-hidden="true" style={{ color: 'var(--ds-orange)' }} />
              Target: {selectedTargetData?.planet_name}{' '}
              <span className="ds-text-muted">({selectedTargetData?.player_name})</span>
            </h3>
            <button className="ds-btn-icon" onClick={() => setSelectedTarget(null)} aria-label="Clear target">
              <X size={16} strokeWidth={2} aria-hidden="true" />
            </button>
          </div>

          <div className="ds-caption">Select Fleets to Attack</div>
          {stationedFleets.length === 0 ? (
            <p className="ds-text-muted" style={{ textAlign: 'center', padding: 'var(--sp-4)' }}>
              No stationed fleets available
            </p>
          ) : (
            <div style={styles.fleetGrid}>
              {stationedFleets.map(fleet => {
                const isSelected = selectedFleets.includes(fleet.id)
                return (
                  <div
                    key={fleet.id}
                    className="ds-list-item"
                    aria-selected={isSelected}
                    onClick={() => toggleFleet(fleet.id)}
                  >
                    <div style={{ fontWeight: 600, display: 'inline-flex', alignItems: 'center', gap: 'var(--sp-2)' }}>
                      <Rocket size={14} strokeWidth={1.75} aria-hidden="true" />
                      {fleet.name}
                    </div>
                    <div className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)', marginTop: 4 }}>
                      {fleet.stacks?.length || 0} stacks
                    </div>
                  </div>
                )
              })}
            </div>
          )}

          <div style={styles.attackActions}>
            <LoadingButton
              className="ds-btn ds-btn-primary"
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
        </section>
      )}

      {/* Recent Combat Results */}
      {resolvedAttacks.length > 0 && (
        <section style={styles.section}>
          <h3 style={styles.sectionHeader}>Recent Attacks</h3>
          <div style={styles.pendingList}>
            {resolvedAttacks.map(pa => (
              <div key={pa.id} className="ds-card" style={styles.pendingRow}>
                <span style={{ fontWeight: 600 }}>{pa.planet_name}</span>
                <span className={`ds-badge ${pa.combat_report_id ? 'ds-badge--success' : 'ds-badge--neutral'}`}>
                  {pa.combat_report_id ? 'Resolved' : 'Pending'}
                </span>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* No radar warning */}
      {incoming && incoming.radar_level === 0 && (
        <section className="ds-card" style={{ background: 'var(--ds-warning-tint)', borderColor: 'var(--ds-warning)' }}>
          <p style={{ margin: 0, color: '#B45309', fontWeight: 600 }}>
            Build a Radar to detect incoming attacks!
          </p>
        </section>
      )}
    </div>
  )
}
