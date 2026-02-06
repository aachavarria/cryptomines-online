import type { ResourcesResponse } from '../types'

interface ResourceBarProps {
  resources: ResourcesResponse | null
  onCollect: () => void
  collecting: boolean
}

function formatNumber(n: number): string {
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + 'M'
  if (n >= 1_000) return (n / 1_000).toFixed(1) + 'K'
  return n.toLocaleString()
}

export default function ResourceBar({ resources, onCollect, collecting }: ResourceBarProps) {
  if (!resources) return null

  return (
    <div className="resource-bar">
      <div className="resource-item">
        <span className="resource-label">Metal</span>
        <span className="resource-value">{formatNumber(resources.metal)}</span>
        <span className="resource-rate">+{formatNumber(resources.metal_per_hour)}/hr</span>
        {resources.pending_metal > 0 && (
          <span className="resource-pending">+{formatNumber(resources.pending_metal)}</span>
        )}
      </div>
      <div className="resource-item">
        <span className="resource-label">He3</span>
        <span className="resource-value">{formatNumber(resources.he3)}</span>
        <span className="resource-rate">+{formatNumber(resources.he3_per_hour)}/hr</span>
        {resources.pending_he3 > 0 && (
          <span className="resource-pending">+{formatNumber(resources.pending_he3)}</span>
        )}
      </div>
      <div className="resource-item">
        <span className="resource-label">Gold</span>
        <span className="resource-value">{formatNumber(resources.gold)}</span>
        <span className="resource-rate">+{formatNumber(resources.gold_per_hour)}/hr</span>
        {resources.pending_gold > 0 && (
          <span className="resource-pending">+{formatNumber(resources.pending_gold)}</span>
        )}
      </div>
      <div className="resource-item resource-storage">
        <span className="resource-label">Storage</span>
        <span className="resource-value">{formatNumber(resources.storage_capacity)}</span>
      </div>
      <button
        className="btn-collect"
        onClick={onCollect}
        disabled={collecting}
      >
        {collecting ? 'Collecting...' : 'Collect Resources'}
      </button>
    </div>
  )
}
