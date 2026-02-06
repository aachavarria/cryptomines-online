import { useState, useEffect } from 'react'
import type { BuildingWithType } from '../types'

interface BuildingListProps {
  buildings: BuildingWithType[]
  onUpgrade: (buildingId: string) => void
  upgrading: string | null
}

function formatTime(seconds: number): string {
  if (seconds <= 0) return 'Done'
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = Math.floor(seconds % 60)
  if (h > 0) return `${h}h ${m}m ${s}s`
  if (m > 0) return `${m}m ${s}s`
  return `${s}s`
}

function UpgradeTimer({ finishAt }: { finishAt: string }) {
  const [remaining, setRemaining] = useState(() => {
    return Math.max(0, (new Date(finishAt).getTime() - Date.now()) / 1000)
  })

  useEffect(() => {
    if (remaining <= 0) return
    const interval = setInterval(() => {
      const left = Math.max(0, (new Date(finishAt).getTime() - Date.now()) / 1000)
      setRemaining(left)
      if (left <= 0) clearInterval(interval)
    }, 1000)
    return () => clearInterval(interval)
  }, [finishAt])

  return <span className="upgrade-timer">{formatTime(remaining)}</span>
}

function categoryLabel(cat: string): string {
  const labels: Record<string, string> = {
    resource: 'Resource',
    core: 'Core',
    military: 'Military',
    defense: 'Defense',
    space: 'Space',
  }
  return labels[cat] || cat
}

export default function BuildingList({ buildings, onUpgrade, upgrading }: BuildingListProps) {
  if (buildings.length === 0) {
    return <p className="empty-text">No buildings yet.</p>
  }

  // Group by category
  const grouped = buildings.reduce<Record<string, BuildingWithType[]>>((acc, b) => {
    const cat = b.category
    if (!acc[cat]) acc[cat] = []
    acc[cat].push(b)
    return acc
  }, {})

  return (
    <div className="building-list">
      {Object.entries(grouped).map(([category, items]) => (
        <div key={category} className="building-category">
          <h4 className="category-title">{categoryLabel(category)}</h4>
          {items.map((b) => (
            <div key={b.id} className="building-card">
              <div className="building-info">
                <span className="building-name">{b.display_name}</span>
                <span className="building-level">Lv {b.level}</span>
                {b.is_upgrading && b.upgrade_finish_at && (
                  <span className="building-status upgrading">
                    Upgrading... <UpgradeTimer finishAt={b.upgrade_finish_at} />
                  </span>
                )}
                {!b.is_upgrading && (
                  <span className="building-status idle">Idle</span>
                )}
              </div>
              <div className="building-actions">
                {!b.is_upgrading && b.level < b.max_level && (
                  <button
                    className="btn-upgrade"
                    onClick={() => onUpgrade(b.id)}
                    disabled={upgrading === b.id}
                  >
                    {upgrading === b.id ? 'Upgrading...' : `Upgrade to Lv ${b.level + 1}`}
                  </button>
                )}
                {b.level >= b.max_level && !b.is_upgrading && (
                  <span className="max-level-badge">MAX</span>
                )}
              </div>
            </div>
          ))}
        </div>
      ))}
    </div>
  )
}
