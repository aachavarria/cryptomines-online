import { useState } from 'react'
import { useCommanders } from '../../hooks/useCommanders'
import { dismissCommander, type Commander } from '../../services/api'
import './CommandersListPanel.css'

type FilterRarity = 'all' | 'common' | 'skill' | 'super'
type SortBy = 'star_rank' | 'accuracy' | 'dodge' | 'speed' | 'electron' | 'name'

export default function CommandersListPanel() {
  const { commanders, loading, refresh } = useCommanders()
  const [filterRarity, setFilterRarity] = useState<FilterRarity>('all')
  const [sortBy, setSortBy] = useState<SortBy>('star_rank')
  const [selectedCommander, setSelectedCommander] = useState<Commander | null>(null)

  // Filter and sort commanders
  const filteredCommanders = commanders
    .filter(c => filterRarity === 'all' || c.rarity === filterRarity)
    .sort((a, b) => {
      if (sortBy === 'name') return a.name.localeCompare(b.name)
      return (b[sortBy] as number) - (a[sortBy] as number)
    })

  // Group by rarity for display
  const grouped = {
    super: filteredCommanders.filter(c => c.rarity === 'super'),
    skill: filteredCommanders.filter(c => c.rarity === 'skill'),
    common: filteredCommanders.filter(c => c.rarity === 'common'),
  }

  const handleDismiss = async (commander: Commander) => {
    if (commander.is_deployed) {
      alert('Cannot dismiss deployed commander. Unassign from fleet first.')
      return
    }

    if (!confirm(`Dismiss ${commander.name}? This action cannot be undone.`)) return

    try {
      await dismissCommander(commander.id)
      alert(`✓ ${commander.name} dismissed`)
      await refresh()
      setSelectedCommander(null)
    } catch (err: any) {
      alert(`✗ Failed to dismiss commander: ${err.message}`)
    }
  }

  if (loading) {
    return <div className="commanders-list-panel loading">Loading commanders...</div>
  }

  return (
    <div className="commanders-list-panel">
      <div className="panel-header">
        <h2>My Commanders</h2>
        <div className="commander-count">{commanders.length} / 60</div>
      </div>

      <div className="controls">
        <div className="filters">
          <label>Filter:</label>
          <button
            className={filterRarity === 'all' ? 'active' : ''}
            onClick={() => setFilterRarity('all')}
          >
            All ({commanders.length})
          </button>
          <button
            className={filterRarity === 'super' ? 'active super' : 'super'}
            onClick={() => setFilterRarity('super')}
          >
            Super ({commanders.filter(c => c.rarity === 'super').length})
          </button>
          <button
            className={filterRarity === 'skill' ? 'active skill' : 'skill'}
            onClick={() => setFilterRarity('skill')}
          >
            Skill ({commanders.filter(c => c.rarity === 'skill').length})
          </button>
          <button
            className={filterRarity === 'common' ? 'active common' : 'common'}
            onClick={() => setFilterRarity('common')}
          >
            Common ({commanders.filter(c => c.rarity === 'common').length})
          </button>
        </div>

        <div className="sort">
          <label>Sort by:</label>
          <select value={sortBy} onChange={e => setSortBy(e.target.value as SortBy)}>
            <option value="star_rank">Star Rank</option>
            <option value="accuracy">Accuracy</option>
            <option value="dodge">Dodge</option>
            <option value="speed">Speed</option>
            <option value="electron">Electron</option>
            <option value="name">Name</option>
          </select>
        </div>
      </div>

      {filteredCommanders.length === 0 ? (
        <div className="empty-state">
          <p>No commanders found with current filters.</p>
          {commanders.length === 0 && (
            <p className="hint">Visit the Command Center to recruit your first commander!</p>
          )}
        </div>
      ) : (
        <div className="commanders-content">
          {filterRarity === 'all' ? (
            // Group view when showing all
            <>
              {(['super', 'skill', 'common'] as const).map(rarity => (
                grouped[rarity].length > 0 && (
                  <div key={rarity} className="rarity-section">
                    <h3 className={`rarity-header ${rarity}`}>
                      {rarity.toUpperCase()} ({grouped[rarity].length})
                    </h3>
                    <div className="commanders-grid">
                      {grouped[rarity].map(commander => (
                        <CommanderCard
                          key={commander.id}
                          commander={commander}
                          isSelected={selectedCommander?.id === commander.id}
                          onClick={() => setSelectedCommander(commander)}
                        />
                      ))}
                    </div>
                  </div>
                )
              ))}
            </>
          ) : (
            // Single grid when filtered
            <div className="commanders-grid">
              {filteredCommanders.map(commander => (
                <CommanderCard
                  key={commander.id}
                  commander={commander}
                  isSelected={selectedCommander?.id === commander.id}
                  onClick={() => setSelectedCommander(commander)}
                />
              ))}
            </div>
          )}
        </div>
      )}

      {/* Commander Detail Modal */}
      {selectedCommander && (
        <div className="commander-detail-modal" onClick={() => setSelectedCommander(null)}>
          <div className="modal-content" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{selectedCommander.name}</h3>
              <button className="close-btn" onClick={() => setSelectedCommander(null)}>×</button>
            </div>

            <div className="modal-body">
              <div className={`rarity-badge ${selectedCommander.rarity}`}>
                {selectedCommander.rarity.toUpperCase()}
              </div>

              <div className="star-rank">
                {'★'.repeat(selectedCommander.star_rank)}
                {'☆'.repeat(15 - selectedCommander.star_rank)}
                <span className="rank-text">Rank {selectedCommander.star_rank}/15</span>
              </div>

              <div className="stats-display">
                <div className="stat-row">
                  <span className="label">Accuracy:</span>
                  <span className="value">{selectedCommander.accuracy}</span>
                </div>
                <div className="stat-row">
                  <span className="label">Dodge:</span>
                  <span className="value">{selectedCommander.dodge}</span>
                </div>
                <div className="stat-row">
                  <span className="label">Speed:</span>
                  <span className="value">{selectedCommander.speed}</span>
                </div>
                <div className="stat-row">
                  <span className="label">Electron:</span>
                  <span className="value">{selectedCommander.electron}</span>
                </div>
              </div>

              <div className="deployment-status">
                {selectedCommander.is_deployed ? (
                  <span className="deployed">⚔️ Deployed to Fleet</span>
                ) : (
                  <span className="available">✓ Available</span>
                )}
              </div>
            </div>

            <div className="modal-actions">
              <button
                className="dismiss-btn"
                onClick={() => handleDismiss(selectedCommander)}
                disabled={selectedCommander.is_deployed}
              >
                Dismiss Commander
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

// Commander Card Component
function CommanderCard({
  commander,
  isSelected,
  onClick,
}: {
  commander: Commander
  isSelected: boolean
  onClick: () => void
}) {
  return (
    <div
      className={`commander-card ${commander.rarity} ${isSelected ? 'selected' : ''}`}
      onClick={onClick}
    >
      <div className="card-header">
        <div className={`rarity-indicator ${commander.rarity}`} />
        <div className="name">{commander.name}</div>
      </div>

      <div className="star-rank">
        {'★'.repeat(commander.star_rank)}
        {'☆'.repeat(Math.min(5, 15 - commander.star_rank))}
        {commander.star_rank > 5 && `+${commander.star_rank - 5}`}
      </div>

      <div className="stats">
        <div className="stat">
          <span className="label">ACC</span>
          <span className="value">{commander.accuracy}</span>
        </div>
        <div className="stat">
          <span className="label">DOD</span>
          <span className="value">{commander.dodge}</span>
        </div>
        <div className="stat">
          <span className="label">SPD</span>
          <span className="value">{commander.speed}</span>
        </div>
        <div className="stat">
          <span className="label">ELEC</span>
          <span className="value">{commander.electron}</span>
        </div>
      </div>

      {commander.is_deployed && (
        <div className="deployed-badge">DEPLOYED</div>
      )}
    </div>
  )
}
