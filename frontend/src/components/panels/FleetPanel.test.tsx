import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import FleetPanel from './FleetPanel'
import type { Fleet } from '../../types'

const fleets: Fleet[] = [
  {
    id: 'f1',
    player_id: 'me',
    name: 'Vanguard',
    formation: 'Phalanx',
    targeting_command: 'Weakest First',
    commander_id: null,
    status: 'stationed',
    stacks: [
      { fleet_id: 'f1', ship_design_id: 'd1', grid_row: 0, grid_col: 0, ship_count: 50, total_movement: 30 },
    ],
    created_at: '2026-04-28T00:00:00Z',
  },
  {
    id: 'f2',
    player_id: 'me',
    name: 'Reserve',
    formation: 'Diamond',
    targeting_command: 'Strongest First',
    commander_id: null,
    status: 'traveling',
    stacks: [],
    created_at: '2026-04-28T00:00:00Z',
  },
]

const create = vi.fn()
const remove = vi.fn()
const update = vi.fn()
const addStack = vi.fn()
const removeStackFromFleet = vi.fn()

vi.mock('../../hooks/useFleets.ts', () => ({
  useFleets: () => ({
    fleets,
    loading: false,
    error: null,
    create,
    update,
    remove,
    addStack,
    removeStackFromFleet,
  }),
}))

vi.mock('../../hooks/useShipDesigns.ts', () => ({
  useShipDesigns: () => ({ designs: [], loading: false, error: null }),
}))

beforeEach(() => {
  vi.clearAllMocks()
  create.mockResolvedValue({
    id: 'new',
    player_id: 'me',
    name: 'Bravo',
    formation: 'Phalanx',
    targeting_command: 'Weakest First',
    commander_id: null,
    status: 'stationed',
    stacks: [],
    created_at: '2026-04-28T00:00:00Z',
  })
})

describe('FleetPanel', () => {
  it('lists existing fleets with totals and statuses', () => {
    render(<FleetPanel />)
    expect(screen.getByText('Vanguard')).toBeInTheDocument()
    expect(screen.getByText('Reserve')).toBeInTheDocument()
    expect(screen.getByText(/2 fleets/)).toBeInTheDocument()
    expect(screen.getByText('stationed')).toBeInTheDocument()
    expect(screen.getByText('traveling')).toBeInTheDocument()
  })

  it('disables Create until a name is typed and then calls create()', async () => {
    render(<FleetPanel />)
    const button = screen.getByRole('button', { name: /Create Fleet/i }) as HTMLButtonElement
    expect(button.disabled).toBe(true)

    const input = screen.getByPlaceholderText(/New fleet name/i)
    fireEvent.change(input, { target: { value: 'Bravo' } })
    expect(button.disabled).toBe(false)

    fireEvent.click(button)
    await waitFor(() => expect(create).toHaveBeenCalledTimes(1))
    expect(create).toHaveBeenCalledWith({
      name: 'Bravo',
      formation: 'Phalanx',
      targeting_command: 'Weakest First',
    })
  })

  it('Disband button calls remove() with the fleet id', () => {
    render(<FleetPanel />)
    const disbandButtons = screen.getAllByRole('button', { name: /Disband/i })
    fireEvent.click(disbandButtons[0])
    expect(remove).toHaveBeenCalledWith('f1')
  })

  it('reports each fleet card stack/ship totals', () => {
    render(<FleetPanel />)
    // Vanguard has 1 stack, 50 ships
    expect(screen.getByText('1/9 positions')).toBeInTheDocument()
    expect(screen.getByText('50 ships')).toBeInTheDocument()
    // Reserve has 0 stacks, 0 ships
    expect(screen.getByText('0/9 positions')).toBeInTheDocument()
  })
})
