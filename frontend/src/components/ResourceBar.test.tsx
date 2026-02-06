import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import ResourceBar from './ResourceBar'
import type { ResourcesResponse } from '../types'

const mockResources: ResourcesResponse = {
  id: 'r1',
  planet_id: 'p1',
  metal: 5000,
  he3: 3000,
  gold: 10000,
  metal_per_hour: 1080,
  he3_per_hour: 1180,
  gold_per_hour: 1300,
  storage_capacity: 100000,
  last_collected_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
  pending_metal: 540,
  pending_he3: 590,
  pending_gold: 650,
}

describe('ResourceBar', () => {
  it('renders nothing when resources is null', () => {
    const { container } = render(
      <ResourceBar resources={null} onCollect={() => {}} collecting={false} />
    )
    expect(container.innerHTML).toBe('')
  })

  it('displays resource labels', () => {
    render(
      <ResourceBar resources={mockResources} onCollect={() => {}} collecting={false} />
    )
    expect(screen.getByText('Metal')).toBeInTheDocument()
    expect(screen.getByText('He3')).toBeInTheDocument()
    expect(screen.getByText('Gold')).toBeInTheDocument()
    expect(screen.getByText('Storage')).toBeInTheDocument()
  })

  it('displays formatted resource values', () => {
    render(
      <ResourceBar resources={mockResources} onCollect={() => {}} collecting={false} />
    )
    // 5000 = "5,000" or "5.0K"
    expect(screen.getByText('5.0K')).toBeInTheDocument()
    // 10000 = "10.0K"
    expect(screen.getByText('10.0K')).toBeInTheDocument()
  })

  it('displays production rates', () => {
    render(
      <ResourceBar resources={mockResources} onCollect={() => {}} collecting={false} />
    )
    expect(screen.getByText('+1.1K/hr')).toBeInTheDocument() // 1080 -> 1.1K
    expect(screen.getByText('+1.2K/hr')).toBeInTheDocument() // 1180 -> 1.2K
    expect(screen.getByText('+1.3K/hr')).toBeInTheDocument() // 1300 -> 1.3K
  })

  it('displays pending resources when > 0', () => {
    render(
      <ResourceBar resources={mockResources} onCollect={() => {}} collecting={false} />
    )
    expect(screen.getByText('+540')).toBeInTheDocument()
    expect(screen.getByText('+590')).toBeInTheDocument()
    expect(screen.getByText('+650')).toBeInTheDocument()
  })

  it('does not display pending when values are 0', () => {
    const zeroPending = { ...mockResources, pending_metal: 0, pending_he3: 0, pending_gold: 0 }
    render(
      <ResourceBar resources={zeroPending} onCollect={() => {}} collecting={false} />
    )
    expect(screen.queryByText('+0')).not.toBeInTheDocument()
  })

  it('shows Collect Resources button', () => {
    render(
      <ResourceBar resources={mockResources} onCollect={() => {}} collecting={false} />
    )
    expect(screen.getByText('Collect Resources')).toBeInTheDocument()
  })

  it('shows Collecting... when collecting', () => {
    render(
      <ResourceBar resources={mockResources} onCollect={() => {}} collecting={true} />
    )
    expect(screen.getByText('Collecting...')).toBeInTheDocument()
  })

  it('calls onCollect when button is clicked', () => {
    const onCollect = vi.fn()
    render(
      <ResourceBar resources={mockResources} onCollect={onCollect} collecting={false} />
    )
    fireEvent.click(screen.getByText('Collect Resources'))
    expect(onCollect).toHaveBeenCalledTimes(1)
  })

  it('disables button when collecting', () => {
    render(
      <ResourceBar resources={mockResources} onCollect={() => {}} collecting={true} />
    )
    const button = screen.getByText('Collecting...')
    expect(button).toBeDisabled()
  })

  it('uses correct GO2 resource names (not ore/energy/credits)', () => {
    render(
      <ResourceBar resources={mockResources} onCollect={() => {}} collecting={false} />
    )
    expect(screen.getByText('Metal')).toBeInTheDocument()
    expect(screen.getByText('He3')).toBeInTheDocument()
    expect(screen.getByText('Gold')).toBeInTheDocument()
    expect(screen.queryByText('Ore')).not.toBeInTheDocument()
    expect(screen.queryByText('Energy')).not.toBeInTheDocument()
    expect(screen.queryByText('Credits')).not.toBeInTheDocument()
  })
})
