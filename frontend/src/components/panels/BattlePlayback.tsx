import { useEffect, useRef, useState, useMemo } from 'react'
import type { RoundData } from '../../services/api'

interface BattlePlaybackProps {
  rounds: RoundData[]
  totalRounds: number
}

const SPEEDS = [0.5, 1, 2, 4]

// BattlePlayback animates the round_by_round combat log so players see the
// fight unfold instead of reading a static table. It plays one round per tick,
// supports pause/scrub/speed, and surfaces per-stack casualties as they land.
export default function BattlePlayback({ rounds, totalRounds }: BattlePlaybackProps) {
  const [index, setIndex] = useState(0)
  const [playing, setPlaying] = useState(false)
  const [speed, setSpeed] = useState(1)
  const tickRef = useRef<number | null>(null)

  const round = rounds[index]
  const totalDealt = useMemo(() => {
    if (!rounds.length) return { hits: 0, damage: 0, kills: 0 }
    let hits = 0,
      damage = 0,
      kills = 0
    for (let i = 0; i <= index; i++) {
      const r = rounds[i]
      if (!r) continue
      for (const a of r.Attacks ?? []) {
        if (a.Hit) hits++
        damage += a.Damage ?? 0
        kills += a.ShipsDestroyed ?? 0
      }
    }
    return { hits, damage, kills }
  }, [rounds, index])

  useEffect(() => {
    if (!playing) {
      if (tickRef.current) {
        window.clearTimeout(tickRef.current)
        tickRef.current = null
      }
      return
    }
    if (index >= rounds.length - 1) {
      setPlaying(false)
      return
    }
    tickRef.current = window.setTimeout(
      () => setIndex((i) => Math.min(i + 1, rounds.length - 1)),
      Math.max(80, 800 / speed),
    )
    return () => {
      if (tickRef.current) window.clearTimeout(tickRef.current)
    }
  }, [playing, index, rounds.length, speed])

  if (!rounds.length) {
    return (
      <p style={{ color: '#94a3b8', fontSize: '0.85rem' }}>
        This report has no round-by-round log.
      </p>
    )
  }

  const onScrub = (e: React.ChangeEvent<HTMLInputElement>) => {
    setIndex(Math.min(rounds.length - 1, Math.max(0, parseInt(e.target.value, 10) || 0)))
    setPlaying(false)
  }

  const reset = () => {
    setIndex(0)
    setPlaying(false)
  }

  const stepBack = () => {
    setIndex((i) => Math.max(0, i - 1))
    setPlaying(false)
  }

  const stepFwd = () => {
    setIndex((i) => Math.min(rounds.length - 1, i + 1))
    setPlaying(false)
  }

  return (
    <div data-testid="battle-playback" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      <div style={{ display: 'flex', gap: 8, alignItems: 'center', flexWrap: 'wrap' }}>
        <button className="btn btn-small btn-secondary" onClick={reset} aria-label="Restart">
          ⏮
        </button>
        <button className="btn btn-small btn-secondary" onClick={stepBack} aria-label="Step back">
          ◀
        </button>
        <button
          className="btn btn-small btn-primary"
          onClick={() => setPlaying((p) => !p)}
          aria-label={playing ? 'Pause' : 'Play'}
        >
          {playing ? '❚❚ Pause' : '▶ Play'}
        </button>
        <button className="btn btn-small btn-secondary" onClick={stepFwd} aria-label="Step forward">
          ▶
        </button>
        <span style={{ color: '#94a3b8', fontSize: '0.8rem', marginLeft: 8 }}>Speed</span>
        {SPEEDS.map((s) => (
          <button
            key={s}
            className={`btn btn-small ${speed === s ? 'btn-primary' : 'btn-secondary'}`}
            onClick={() => setSpeed(s)}
          >
            {s}x
          </button>
        ))}
        <span style={{ color: '#cbd5e1', fontSize: '0.85rem', marginLeft: 'auto' }}>
          Round {index + 1} / {totalRounds || rounds.length}
        </span>
      </div>

      <input
        type="range"
        min={0}
        max={rounds.length - 1}
        value={index}
        onChange={onScrub}
        aria-label="Round timeline"
        style={{ width: '100%' }}
      />

      <div
        style={{
          padding: 12,
          border: '1px solid #333',
          borderRadius: 6,
          background: 'rgba(0,0,0,0.4)',
          minHeight: 120,
        }}
      >
        <h5 style={{ marginBottom: 8 }}>Round {round.RoundNumber}</h5>
        {(round.Attacks ?? []).length === 0 ? (
          <p style={{ color: '#94a3b8', fontSize: '0.85rem' }}>No attacks landed this round.</p>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
            {round.Attacks.map((a, i) => (
              <div
                key={i}
                style={{
                  fontSize: '0.85rem',
                  padding: '4px 6px',
                  borderRadius: 3,
                  background: a.Hit
                    ? a.AttackerSide === 'attacker'
                      ? 'rgba(74,222,128,0.10)'
                      : 'rgba(248,113,113,0.10)'
                    : 'rgba(148,163,184,0.05)',
                  color: a.Hit ? '#e2e8f0' : '#64748b',
                }}
              >
                <strong style={{ color: a.AttackerSide === 'attacker' ? '#4ade80' : '#f87171' }}>
                  [{a.AttackerSide}]
                </strong>{' '}
                {a.Hit ? (a.CriticalHit ? '⚡ CRIT' : '⚔ HIT') : '✗ miss'}{' '}
                <strong style={{ color: a.DefenderSide === 'attacker' ? '#4ade80' : '#f87171' }}>
                  [{a.DefenderSide}]
                </strong>
                {a.Hit && (
                  <span style={{ marginLeft: 8 }}>
                    {a.Damage} dmg{a.ShipsDestroyed > 0 ? ` · 💥 ${a.ShipsDestroyed} kills` : ''}
                  </span>
                )}
              </div>
            ))}
          </div>
        )}

        {Object.keys(round.Casualties ?? {}).length > 0 && (
          <div
            style={{
              marginTop: 8,
              fontSize: '0.75rem',
              color: '#94a3b8',
              borderTop: '1px solid #2a2a2a',
              paddingTop: 6,
            }}
          >
            <strong>Casualties:</strong>{' '}
            {Object.entries(round.Casualties).map(([id, n]) => (
              <span key={id} style={{ marginRight: 8 }}>
                {id.slice(0, 8)}…: {n}
              </span>
            ))}
          </div>
        )}
      </div>

      <div
        style={{
          fontSize: '0.8rem',
          color: '#94a3b8',
          display: 'flex',
          gap: 16,
          flexWrap: 'wrap',
        }}
      >
        <span>Cumulative hits: <strong style={{ color: '#cbd5e1' }}>{totalDealt.hits}</strong></span>
        <span>Damage: <strong style={{ color: '#cbd5e1' }}>{totalDealt.damage}</strong></span>
        <span>Ships destroyed: <strong style={{ color: '#cbd5e1' }}>{totalDealt.kills}</strong></span>
      </div>
    </div>
  )
}
