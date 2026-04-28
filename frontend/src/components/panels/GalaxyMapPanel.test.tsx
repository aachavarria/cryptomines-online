import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor, act } from '@testing-library/react'
import GalaxyMapPanel from './GalaxyMapPanel'
import * as api from '../../services/api'
import type { Planet } from '../../types'
import type { SectorPlanet } from '../../types/corps'

vi.mock('../../services/api', async () => {
  const actual = await vi.importActual<typeof import('../../services/api')>('../../services/api')
  return {
    ...actual,
    listPlanets: vi.fn(),
    listFleets: vi.fn(),
    getGalaxySector: vi.fn(),
    attackPlanet: vi.fn(),
    attackRBP: vi.fn(),
    getPendingAttacks: vi.fn(),
    getIncomingAttacks: vi.fn(),
  }
})

const mocked = api as unknown as {
  listPlanets: ReturnType<typeof vi.fn>
  listFleets: ReturnType<typeof vi.fn>
  getGalaxySector: ReturnType<typeof vi.fn>
  attackPlanet: ReturnType<typeof vi.fn>
  attackRBP: ReturnType<typeof vi.fn>
  getPendingAttacks: ReturnType<typeof vi.fn>
  getIncomingAttacks: ReturnType<typeof vi.fn>
}

const homeworld: Planet = {
  id: 'home-1',
  player_id: 'player-self',
  name: 'Home',
  position_x: 1500,
  position_y: 1500,
  is_homeworld: true,
  is_rbp: false,
  rbp_level: 0,
  created_at: '2026-04-28T00:00:00Z',
  updated_at: '2026-04-28T00:00:00Z',
}

const baseSector = {
  center_x: 1500,
  center_y: 1500,
  radius: 100,
  planets: [
    {
      id: 'home-1',
      name: 'Home',
      position_x: 1500,
      position_y: 1500,
      is_homeworld: true,
      is_rbp: false,
      rbp_level: 0,
      is_own: true,
      defense_strength: 5,
    },
    {
      id: 'enemy-1',
      name: 'Enemy Alpha',
      position_x: 1550,
      position_y: 1550,
      is_homeworld: false,
      is_rbp: false,
      rbp_level: 0,
      is_own: false,
      owner_id: 'enemy-player',
      owner_name: 'Hostile',
      defense_strength: 12,
    },
    {
      id: 'rbp-1',
      name: 'RBP-Alpha',
      position_x: 1480,
      position_y: 1480,
      is_homeworld: false,
      is_rbp: true,
      rbp_level: 3,
      is_own: false,
      defense_strength: 30,
    },
  ] as SectorPlanet[],
}

beforeEach(() => {
  mocked.listPlanets.mockResolvedValue([homeworld])
  mocked.listFleets.mockResolvedValue([
    {
      id: 'fleet-1',
      player_id: 'player-self',
      name: 'Vanguard',
      formation: 'phalanx',
      targeting_command: 'max_attack',
      commander_id: null,
      status: 'stationed',
      stacks: [{ fleet_id: 'fleet-1', ship_design_id: 'd1', grid_row: 0, grid_col: 0, ship_count: 50 }],
      created_at: '2026-04-28T00:00:00Z',
    },
  ])
  mocked.getGalaxySector.mockResolvedValue(baseSector)
  mocked.getPendingAttacks.mockResolvedValue([])
  mocked.getIncomingAttacks.mockResolvedValue({
    radar_level: 0,
    detect_advance: 0,
    incoming_attacks: [],
  })
  mocked.attackPlanet.mockResolvedValue({
    pending_attack_id: 'pa-1',
    travel_seconds: 60,
    arrival_at: '2026-04-28T00:01:00Z',
    fleets_dispatched: 1,
  })
  mocked.attackRBP.mockResolvedValue({
    report_id: 'r-1',
    result: 'attacker_win',
    total_rounds: 6,
    conquered: false,
  })
})

describe('GalaxyMapPanel', () => {
  it('loads the sector centered on the homeworld and renders planets', async () => {
    render(<GalaxyMapPanel onClose={vi.fn()} />)
    await waitFor(() => expect(mocked.getGalaxySector).toHaveBeenCalled())
    expect(mocked.getGalaxySector.mock.calls[0]).toEqual([1500, 1500, 100])

    expect(await screen.findByText('Home')).toBeInTheDocument()
    expect(screen.getByText('Enemy Alpha')).toBeInTheDocument()
    expect(screen.getByText('RBP-Alpha')).toBeInTheDocument()
  })

  it('zoom buttons change the radius and refetch the sector', async () => {
    render(<GalaxyMapPanel onClose={vi.fn()} />)
    await screen.findByText('Home')
    expect(mocked.getGalaxySector).toHaveBeenCalledWith(1500, 1500, 100)

    fireEvent.click(screen.getByRole('button', { name: /Zoom \+/ }))
    await waitFor(() => {
      const lastCall = mocked.getGalaxySector.mock.calls.at(-1)!
      expect(lastCall[2]).toBeLessThan(100)
    })

    fireEvent.click(screen.getByRole('button', { name: /Zoom −/ }))
    await waitFor(() => {
      const lastCall = mocked.getGalaxySector.mock.calls.at(-1)!
      expect(lastCall[2]).toBeGreaterThanOrEqual(98)
    })
    fireEvent.click(screen.getByRole('button', { name: /Zoom −/ }))
    await waitFor(() => {
      const lastCall = mocked.getGalaxySector.mock.calls.at(-1)!
      expect(lastCall[2]).toBeGreaterThan(100)
    })
  })

  it('clicking an enemy planet exposes the dispatch panel and calls attackPlanet', async () => {
    render(<GalaxyMapPanel onClose={vi.fn()} />)
    const enemyLabel = await screen.findByText('Enemy Alpha')
    const enemyGroup = enemyLabel.parentElement!
    act(() => {
      fireEvent.click(enemyGroup)
    })

    expect(screen.getByText(/Enemy planet/i)).toBeInTheDocument()
    expect(screen.getByText(/Hostile/)).toBeInTheDocument()
    const fleetCheckbox = screen.getByRole('checkbox')
    fireEvent.click(fleetCheckbox)

    fireEvent.click(screen.getByRole('button', { name: /Attack Planet/i }))
    await waitFor(() => expect(mocked.attackPlanet).toHaveBeenCalled())
    expect(mocked.attackPlanet).toHaveBeenCalledWith({
      defender_planet_id: 'enemy-1',
      fleet_ids: ['fleet-1'],
    })
    expect(await screen.findByText(/Fleet dispatched/i)).toBeInTheDocument()
  })

  it('clicking an RBP routes the attack to attackRBP', async () => {
    render(<GalaxyMapPanel onClose={vi.fn()} />)
    const rbpLabel = await screen.findByText('RBP-Alpha')
    act(() => {
      fireEvent.click(rbpLabel.parentElement!)
    })
    expect(screen.getByText(/RBP \(lvl 3\)/i)).toBeInTheDocument()
    fireEvent.click(screen.getByRole('checkbox'))
    fireEvent.click(screen.getByRole('button', { name: /Attack RBP/i }))
    await waitFor(() => expect(mocked.attackRBP).toHaveBeenCalledWith('rbp-1', { fleet_ids: ['fleet-1'] }))
    expect(await screen.findByText(/RBP defeated in 6 rounds/i)).toBeInTheDocument()
  })

  it('does not show a dispatch button for own planets', async () => {
    render(<GalaxyMapPanel onClose={vi.fn()} />)
    const homeLabel = await screen.findByText('Home')
    act(() => {
      fireEvent.click(homeLabel.parentElement!)
    })
    expect(screen.getByText(/Your planet/i)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /Attack/i })).not.toBeInTheDocument()
  })

  it('Esc closes the panel', () => {
    const onClose = vi.fn()
    render(<GalaxyMapPanel onClose={onClose} />)
    fireEvent.keyDown(window, { key: 'Escape' })
    expect(onClose).toHaveBeenCalled()
  })
})
