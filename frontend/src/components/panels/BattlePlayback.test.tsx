import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, fireEvent, act } from '@testing-library/react'
import BattlePlayback from './BattlePlayback'
import type { RoundData } from '../../services/api'

const sampleRounds: RoundData[] = [
  {
    RoundNumber: 1,
    Attacks: [
      {
        AttackerStackID: 'a1',
        DefenderStackID: 'd1',
        AttackerSide: 'attacker',
        DefenderSide: 'defender',
        Hit: true,
        Damage: 100,
        ShieldDamage: 60,
        StructureDamage: 40,
        ShipsDestroyed: 2,
      },
      {
        AttackerStackID: 'd1',
        DefenderStackID: 'a1',
        AttackerSide: 'defender',
        DefenderSide: 'attacker',
        Hit: false,
        Damage: 0,
        ShieldDamage: 0,
        StructureDamage: 0,
        ShipsDestroyed: 0,
      },
    ],
    Casualties: { d1: 2 },
  },
  {
    RoundNumber: 2,
    Attacks: [
      {
        AttackerStackID: 'a1',
        DefenderStackID: 'd1',
        AttackerSide: 'attacker',
        DefenderSide: 'defender',
        Hit: true,
        CriticalHit: true,
        Damage: 200,
        ShieldDamage: 0,
        StructureDamage: 200,
        ShipsDestroyed: 5,
      },
    ],
    Casualties: { d1: 5 },
  },
  {
    RoundNumber: 3,
    Attacks: [],
    Casualties: {},
  },
]

describe('BattlePlayback', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })
  afterEach(() => {
    vi.useRealTimers()
  })

  it('shows the empty-log placeholder when no rounds are provided', () => {
    render(<BattlePlayback rounds={[]} totalRounds={0} />)
    expect(screen.getByText(/has no round-by-round log/i)).toBeInTheDocument()
  })

  it('starts on round 1 and shows its attacks', () => {
    render(<BattlePlayback rounds={sampleRounds} totalRounds={3} />)
    expect(screen.getByRole('heading', { level: 5, name: /Round 1/i })).toBeInTheDocument()
    expect(screen.getByText(/Round 1 \/ 3/)).toBeInTheDocument()
    expect(screen.getByText(/100 dmg/)).toBeInTheDocument()
    expect(screen.getByText(/💥 2 kills/)).toBeInTheDocument()
  })

  it('advances rounds while playing and stops at the last one', () => {
    render(<BattlePlayback rounds={sampleRounds} totalRounds={3} />)
    fireEvent.click(screen.getByRole('button', { name: /play/i }))

    act(() => {
      vi.advanceTimersByTime(900)
    })
    expect(screen.getByText(/Round 2 \/ 3/)).toBeInTheDocument()
    expect(screen.getByText(/⚡ CRIT/)).toBeInTheDocument()

    act(() => {
      vi.advanceTimersByTime(900)
    })
    expect(screen.getByText(/Round 3 \/ 3/)).toBeInTheDocument()
    expect(screen.getByText(/No attacks landed this round/i)).toBeInTheDocument()

    act(() => {
      vi.advanceTimersByTime(2000)
    })
    expect(screen.getByText(/Round 3 \/ 3/)).toBeInTheDocument()
  })

  it('step-back/step-forward buttons move the index without auto-playing', () => {
    render(<BattlePlayback rounds={sampleRounds} totalRounds={3} />)
    fireEvent.click(screen.getByRole('button', { name: /step forward/i }))
    expect(screen.getByText(/Round 2 \/ 3/)).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: /step back/i }))
    expect(screen.getByText(/Round 1 \/ 3/)).toBeInTheDocument()
  })

  it('scrubbing the timeline jumps directly to a round', () => {
    render(<BattlePlayback rounds={sampleRounds} totalRounds={3} />)
    const slider = screen.getByLabelText(/round timeline/i) as HTMLInputElement
    fireEvent.change(slider, { target: { value: '2' } })
    expect(screen.getByText(/Round 3 \/ 3/)).toBeInTheDocument()
  })

  it('cumulative stats reflect attacks up to the current round', () => {
    render(<BattlePlayback rounds={sampleRounds} totalRounds={3} />)
    expect(screen.getByText(/Cumulative hits/i).parentElement).toHaveTextContent('1')
    fireEvent.click(screen.getByRole('button', { name: /step forward/i }))
    expect(screen.getByText(/Cumulative hits/i).parentElement).toHaveTextContent('2')
    expect(screen.getByText(/Damage/i).parentElement).toHaveTextContent('300')
    expect(screen.getByText(/Ships destroyed/i).parentElement).toHaveTextContent('7')
  })

  it('restart button returns to round 1', () => {
    render(<BattlePlayback rounds={sampleRounds} totalRounds={3} />)
    fireEvent.click(screen.getByRole('button', { name: /step forward/i }))
    fireEvent.click(screen.getByRole('button', { name: /step forward/i }))
    expect(screen.getByText(/Round 3 \/ 3/)).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: /restart/i }))
    expect(screen.getByText(/Round 1 \/ 3/)).toBeInTheDocument()
  })

  it('faster speed advances rounds in fewer ticks', () => {
    render(<BattlePlayback rounds={sampleRounds} totalRounds={3} />)
    fireEvent.click(screen.getByRole('button', { name: /^4x$/ }))
    fireEvent.click(screen.getByRole('button', { name: /play/i }))
    act(() => {
      vi.advanceTimersByTime(250)
    })
    expect(screen.getByText(/Round 2 \/ 3/)).toBeInTheDocument()
  })
})
