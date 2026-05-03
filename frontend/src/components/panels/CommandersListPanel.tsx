import { useState, useEffect } from 'react'
import { createPortal } from 'react-dom'
import { X, Plus } from 'lucide-react'
import { useCommanders } from '../../hooks/useCommanders'
import { useGameContext } from '../../contexts/GameContext.tsx'
import { dismissCommander, type Commander } from '../../services/api'
import './CommandersListPanel.css'

// Expertise grade descriptions per GO2
function getGradeDesc(type: 'weapon' | 'ship', grade: string): string {
  if (type === 'weapon') {
    switch (grade) {
      case 'S': return '+30% damage'
      case 'A': return '+10% damage'
      case 'B': return 'No bonus'
      case 'C': return '-10% damage'
      case 'D': return '-30% damage'
      default: return 'Unknown'
    }
  }
  switch (grade) {
    case 'S': return '+10% dealt, -10% received'
    case 'A': return '+5% dealt, -10% received'
    case 'B': return 'No bonus'
    case 'C': return '-5% dealt, +5% received'
    case 'D': return '-10% dealt, +10% received'
    default: return 'Unknown'
  }
}

function getExpertiseTooltip(type: 'weapon' | 'ship'): string {
  if (type === 'weapon') {
    return `Weapon Expertise: Affects damage dealt by each weapon class. S=+30%, A=+10%, B=0%, C=-10%, D=-30%`
  }
  return `Ship Expertise: Affects damage dealt/received per ship class. S=+10%/-10%, A=+5%/-10%, B=0%, C=-5%/+5%, D=-10%/+10%`
}

function formatCategory(key: string): string {
  return key.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

function hasAnyGrades(map: Record<string, string> | undefined): boolean {
  return !!map && Object.keys(map).length > 0
}

type FilterRarity = 'all' | 'common' | 'skill' | 'super'
type SortBy = 'star_rank' | 'accuracy' | 'dodge' | 'speed' | 'electron' | 'name'

interface CommandersListPanelProps {
  onClose: () => void
}

// Map game rarity → DS rarity badge
function rarityBadgeClass(rarity: string): string {
  switch (rarity) {
    case 'super': return 'ds-badge ds-badge--rarity-epic'
    case 'skill': return 'ds-badge ds-badge--rarity-rare'
    case 'common':
    default:      return 'ds-badge ds-badge--rarity-common'
  }
}

export default function CommandersListPanel({ onClose }: CommandersListPanelProps) {
  const { commanders, loading, refresh } = useCommanders()
  const { state, openCommandCenterPanel } = useGameContext()
  const [filterRarity, setFilterRarity] = useState<FilterRarity>('all')
  const [sortBy, setSortBy] = useState<SortBy>('star_rank')
  const [selectedCommander, setSelectedCommander] = useState<Commander | null>(null)

  const hasCommandCenter = state.buildings.some(b => b.type_name === 'command_center')

  const handleOpenRecruit = () => {
    openCommandCenterPanel()
    onClose()
  }

  // ESC to close
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

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
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error'
      alert(`Failed to dismiss commander: ${message}`)
    }
  }

  if (loading) {
    return createPortal(
      <div className="ds-modal-backdrop">
        <div
          className="ds-modal"
          style={{ width: 'min(1200px, 90vw)', maxHeight: '85vh', padding: 'var(--sp-10)', textAlign: 'center' }}
        >
          Loading commanders...
        </div>
      </div>,
      document.body
    )
  }

  return createPortal(
    <div
      className="ds-modal-backdrop"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div
        className="ds-modal"
        style={{
          width: 'min(1200px, 90vw)',
          maxHeight: '85vh',
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <div className="ds-modal-header">
          <div className="ds-row" style={{ flex: 1 }}>
            <h2 className="ds-modal-title">My Commanders</h2>
            <span className="ds-badge ds-badge--neutral ds-mono">{commanders.length} / 60</span>
          </div>
          {hasCommandCenter && (
            <button
              className="ds-btn-secondary ds-btn--sm"
              onClick={handleOpenRecruit}
              title="Open Command Center recruitment"
            >
              <Plus size={14} aria-hidden /> Recruit
            </button>
          )}
          <button className="ds-btn-icon" onClick={onClose} aria-label="Close">
            <X size={16} />
          </button>
        </div>

        <div className="ds-modal-body" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--sp-4)' }}>
          {/* Controls */}
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              gap: 'var(--sp-4)',
              flexWrap: 'wrap',
            }}
          >
            <div className="ds-row" style={{ flexWrap: 'wrap' }}>
              <span className="ds-caption">Filter:</span>
              <button
                className={filterRarity === 'all' ? 'ds-btn-secondary ds-btn--sm' : 'ds-btn-ghost ds-btn--sm'}
                onClick={() => setFilterRarity('all')}
              >
                All ({commanders.length})
              </button>
              <button
                className={filterRarity === 'super' ? 'ds-btn-secondary ds-btn--sm' : 'ds-btn-ghost ds-btn--sm'}
                onClick={() => setFilterRarity('super')}
              >
                Super ({commanders.filter(c => c.rarity === 'super').length})
              </button>
              <button
                className={filterRarity === 'skill' ? 'ds-btn-secondary ds-btn--sm' : 'ds-btn-ghost ds-btn--sm'}
                onClick={() => setFilterRarity('skill')}
              >
                Skill ({commanders.filter(c => c.rarity === 'skill').length})
              </button>
              <button
                className={filterRarity === 'common' ? 'ds-btn-secondary ds-btn--sm' : 'ds-btn-ghost ds-btn--sm'}
                onClick={() => setFilterRarity('common')}
              >
                Common ({commanders.filter(c => c.rarity === 'common').length})
              </button>
            </div>

            <div className="ds-row">
              <label className="ds-caption" htmlFor="commanders-sort">Sort by:</label>
              <select
                id="commanders-sort"
                className="ds-select"
                value={sortBy}
                onChange={e => setSortBy(e.target.value as SortBy)}
                style={{ width: 'auto' }}
              >
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
            <div style={{ textAlign: 'center', padding: 'var(--sp-10) var(--sp-5)' }}>
              <p className="ds-text-muted" style={{ marginBottom: 'var(--sp-2)' }}>
                No commanders found with current filters.
              </p>
              {commanders.length === 0 && (
                <>
                  {hasCommandCenter ? (
                    <>
                      <p className="ds-text-muted" style={{ fontStyle: 'italic', marginBottom: 'var(--sp-3)' }}>
                        Recruit your first commander to begin!
                      </p>
                      <button className="ds-btn-secondary" onClick={handleOpenRecruit}>
                        Recruit Commander
                      </button>
                    </>
                  ) : (
                    <p className="ds-text-muted" style={{ fontStyle: 'italic' }}>
                      Build a Command Center to recruit commanders.
                    </p>
                  )}
                </>
              )}
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--sp-5)' }}>
              {filterRarity === 'all' ? (
                <>
                  {(['super', 'skill', 'common'] as const).map(rarity => (
                    grouped[rarity].length > 0 && (
                      <section key={rarity}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--sp-2)', marginBottom: 'var(--sp-3)' }}>
                          <span className={rarityBadgeClass(rarity)}>
                            {rarity}
                          </span>
                          <span className="ds-caption">({grouped[rarity].length})</span>
                        </div>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--sp-2)' }}>
                          {grouped[rarity].map(commander => (
                            <CommanderRow
                              key={commander.id}
                              commander={commander}
                              isSelected={selectedCommander?.id === commander.id}
                              onClick={() => setSelectedCommander(commander)}
                            />
                          ))}
                        </div>
                      </section>
                    )
                  ))}
                </>
              ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--sp-2)' }}>
                  {filteredCommanders.map(commander => (
                    <CommanderRow
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
        </div>

        {/* Commander Detail Modal */}
        {selectedCommander && (
          <div
            className="ds-modal-backdrop"
            onClick={() => setSelectedCommander(null)}
            style={{ zIndex: 'calc(var(--z-modal) + 10)' }}
          >
            <div
              className="ds-modal ds-modal--sm"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="ds-modal-header">
                <h3 className="ds-modal-title">{selectedCommander.name}</h3>
                <button className="ds-btn-icon" onClick={() => setSelectedCommander(null)} aria-label="Close">
                  <X size={16} />
                </button>
              </div>

              <div className="ds-modal-body" style={{ textAlign: 'center' }}>
                <div style={{ marginBottom: 'var(--sp-3)' }}>
                  <span className={rarityBadgeClass(selectedCommander.rarity)}>
                    {selectedCommander.rarity}
                  </span>
                </div>

                <div className="commander-stars">
                  {'★'.repeat(selectedCommander.star_rank)}
                  {'☆'.repeat(15 - selectedCommander.star_rank)}
                  <div className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)', marginTop: 'var(--sp-1)' }}>
                    Rank {selectedCommander.star_rank}/15
                  </div>
                </div>

                <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--sp-2)', margin: 'var(--sp-4) 0' }}>
                  {[
                    { label: 'Accuracy', value: selectedCommander.accuracy },
                    { label: 'Dodge',    value: selectedCommander.dodge },
                    { label: 'Speed',    value: selectedCommander.speed },
                    { label: 'Electron', value: selectedCommander.electron },
                  ].map(({ label, value }) => (
                    <div key={label} className="ds-card ds-row--between">
                      <span className="ds-text-muted">{label}:</span>
                      <span className="ds-mono" style={{ color: 'var(--ds-teal-dark)', fontWeight: 'var(--fw-semibold)' }}>
                        {value}
                      </span>
                    </div>
                  ))}
                </div>

                {(hasAnyGrades(selectedCommander.weapon_expertise) || hasAnyGrades(selectedCommander.ship_expertise)) && (
                  <div className="ds-card" style={{ textAlign: 'left', marginTop: 'var(--sp-3)' }}>
                    <h4 className="ds-h3" style={{ marginBottom: 'var(--sp-2)' }}>Expertise</h4>
                    {hasAnyGrades(selectedCommander.weapon_expertise) && (
                      <div title={getExpertiseTooltip('weapon')}>
                        <div className="ds-caption">Weapon</div>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
                          {Object.entries(selectedCommander.weapon_expertise!).map(([cat, grade]) => (
                            <ExpertiseRow key={cat} cat={cat} grade={grade} desc={getGradeDesc('weapon', grade)} />
                          ))}
                        </div>
                      </div>
                    )}
                    {hasAnyGrades(selectedCommander.ship_expertise) && (
                      <div title={getExpertiseTooltip('ship')} style={{ marginTop: 'var(--sp-3)' }}>
                        <div className="ds-caption">Ship</div>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
                          {Object.entries(selectedCommander.ship_expertise!).map(([cat, grade]) => (
                            <ExpertiseRow key={cat} cat={cat} grade={grade} desc={getGradeDesc('ship', grade)} />
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                )}

                <div style={{ margin: 'var(--sp-4) 0' }}>
                  {selectedCommander.is_deployed ? (
                    <span className="ds-badge ds-badge--orange">Deployed to Fleet</span>
                  ) : (
                    <span className="ds-badge ds-badge--success">Available</span>
                  )}
                </div>
              </div>

              <div className="ds-modal-footer">
                <button
                  className="ds-btn-danger ds-btn--block"
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
    </div>,
    document.body
  )
}

// Expertise row helper
function ExpertiseRow({ cat, grade, desc }: { cat: string; grade: string; desc: string }) {
  const gradeColor = (() => {
    switch (grade) {
      case 'S': return 'var(--ds-orange)'
      case 'A': return 'var(--ds-success)'
      case 'B': return 'var(--ds-text-muted)'
      case 'C': return 'var(--ds-warning)'
      case 'D': return 'var(--ds-danger)'
      default:  return 'var(--ds-text-muted)'
    }
  })()
  return (
    <div
      style={{
        display: 'grid',
        gridTemplateColumns: '1fr auto 1.2fr',
        gap: 'var(--sp-2)',
        alignItems: 'center',
        padding: '4px 8px',
        background: 'var(--ds-surface)',
        border: '1px solid var(--ds-border)',
        borderRadius: 'var(--r-sm)',
        fontSize: 'var(--fs-caption)',
      }}
    >
      <span className="ds-text-muted">{formatCategory(cat)}</span>
      <span
        style={{
          fontWeight: 'var(--fw-bold)',
          textAlign: 'center',
          minWidth: 22,
          padding: '1px 6px',
          borderRadius: 'var(--r-sm)',
          color: gradeColor,
          border: `1px solid ${gradeColor}`,
        }}
      >
        {grade}
      </span>
      <span className="ds-text-soft" style={{ textAlign: 'right' }}>
        {desc}
      </span>
    </div>
  )
}

// Commander Row Component (.ds-list-item)
function CommanderRow({
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
      className={`ds-list-item ${isSelected ? 'is-selected' : ''}`}
      onClick={onClick}
      role="button"
      tabIndex={0}
      style={{
        display: 'grid',
        gridTemplateColumns: 'auto 1fr auto auto',
        gap: 'var(--sp-3)',
        alignItems: 'center',
      }}
    >
      <span className={rarityBadgeClass(commander.rarity)}>
        {commander.rarity}
      </span>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 2, minWidth: 0 }}>
        <div style={{ fontWeight: 'var(--fw-semibold)', color: 'var(--ds-text)', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
          {commander.name}
        </div>
        <div style={{ color: 'var(--ds-orange)', fontSize: 'var(--fs-sm)', letterSpacing: 1 }}>
          {'★'.repeat(commander.star_rank)}
          {'☆'.repeat(Math.min(5, 15 - commander.star_rank))}
          {commander.star_rank > 5 && `+${commander.star_rank - 5}`}
        </div>
      </div>
      <div
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(4, auto)',
          gap: 'var(--sp-3)',
        }}
      >
        {[
          { label: 'ACC', value: commander.accuracy },
          { label: 'DOD', value: commander.dodge },
          { label: 'SPD', value: commander.speed },
          { label: 'ELEC', value: commander.electron },
        ].map(({ label, value }) => (
          <div key={label} style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
            <span className="ds-caption">{label}</span>
            <span className="ds-mono" style={{ color: 'var(--ds-teal-dark)', fontWeight: 'var(--fw-bold)' }}>
              {value}
            </span>
          </div>
        ))}
      </div>
      <div>
        {commander.is_deployed && (
          <span className="ds-badge ds-badge--orange">Deployed</span>
        )}
      </div>
    </div>
  )
}
