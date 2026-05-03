import { useEffect, useRef, useState, useMemo } from 'react'
import { Play, Pause, SkipBack, SkipForward, RotateCcw } from 'lucide-react'
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
      <p className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)' }}>
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

  const progressPct = rounds.length > 1 ? (index / (rounds.length - 1)) * 100 : 100

  return (
    <div data-testid="battle-playback" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      <div style={{ display: 'flex', gap: 'var(--sp-2)', alignItems: 'center', flexWrap: 'wrap' }}>
        <button className="ds-btn-icon" onClick={reset} aria-label="Restart">
          <RotateCcw size={16} />
        </button>
        <button className="ds-btn-icon" onClick={stepBack} aria-label="Step back">
          <SkipBack size={16} />
        </button>
        <button
          className="ds-btn-icon"
          onClick={() => setPlaying((p) => !p)}
          aria-label={playing ? 'Pause' : 'Play'}
          style={{ color: 'var(--ds-teal)' }}
        >
          {playing ? <Pause size={16} /> : <Play size={16} />}
        </button>
        <button className="ds-btn-icon" onClick={stepFwd} aria-label="Step forward">
          <SkipForward size={16} />
        </button>
        <span
          className="ds-caption"
          style={{ marginLeft: 'var(--sp-2)' }}
        >
          Speed
        </span>
        {SPEEDS.map((s) => (
          <button
            key={s}
            className={speed === s ? 'ds-btn-secondary ds-btn--sm' : 'ds-btn-ghost ds-btn--sm'}
            onClick={() => setSpeed(s)}
          >
            {s}x
          </button>
        ))}
        <span className="ds-mono" style={{ fontSize: 'var(--fs-sm)', marginLeft: 'auto', color: 'var(--ds-text-muted)' }}>
          Round {index + 1} / {totalRounds || rounds.length}
        </span>
      </div>

      <div className="ds-bar" aria-label="Round progress">
        <div className="ds-bar-fill" style={{ width: `${progressPct}%` }} />
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
        className="ds-card"
        style={{ minHeight: 120 }}
      >
        <h5 className="ds-h3" style={{ marginBottom: 'var(--sp-2)' }}>
          Round {round.RoundNumber}
        </h5>
        {(round.Attacks ?? []).length === 0 ? (
          <p className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)' }}>
            No attacks landed this round.
          </p>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
            {round.Attacks.map((a, i) => (
              <div
                key={i}
                style={{
                  fontSize: 'var(--fs-sm)',
                  padding: '4px 6px',
                  borderRadius: 'var(--r-sm)',
                  background: a.Hit
                    ? a.AttackerSide === 'attacker'
                      ? 'var(--ds-success-tint)'
                      : 'var(--ds-danger-tint)'
                    : 'var(--ds-surface-3)',
                  color: a.Hit ? 'var(--ds-text)' : 'var(--ds-text-soft)',
                }}
              >
                <strong
                  style={{
                    color:
                      a.AttackerSide === 'attacker' ? 'var(--ds-success)' : 'var(--ds-danger)',
                  }}
                >
                  [{a.AttackerSide}]
                </strong>{' '}
                {a.Hit ? (a.CriticalHit ? '⚡ CRIT' : '⚔ HIT') : '✗ miss'}{' '}
                <strong
                  style={{
                    color:
                      a.DefenderSide === 'attacker' ? 'var(--ds-success)' : 'var(--ds-danger)',
                  }}
                >
                  [{a.DefenderSide}]
                </strong>
                {a.Hit && (
                  <span className="ds-mono" style={{ marginLeft: 'var(--sp-2)' }}>
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
              marginTop: 'var(--sp-2)',
              fontSize: 'var(--fs-caption)',
              color: 'var(--ds-text-muted)',
              borderTop: '1px solid var(--ds-border)',
              paddingTop: 6,
            }}
          >
            <strong>Casualties:</strong>{' '}
            {Object.entries(round.Casualties).map(([id, n]) => (
              <span key={id} className="ds-mono" style={{ marginRight: 'var(--sp-2)' }}>
                {id.slice(0, 8)}…: {n}
              </span>
            ))}
          </div>
        )}
      </div>

      <div
        className="ds-text-muted"
        style={{
          fontSize: 'var(--fs-sm)',
          display: 'flex',
          gap: 'var(--sp-4)',
          flexWrap: 'wrap',
        }}
      >
        <span>
          Cumulative hits:{' '}
          <strong className="ds-mono" style={{ color: 'var(--ds-text)' }}>{totalDealt.hits}</strong>
        </span>
        <span>
          Damage:{' '}
          <strong className="ds-mono" style={{ color: 'var(--ds-text)' }}>{totalDealt.damage}</strong>
        </span>
        <span>
          Ships destroyed:{' '}
          <strong className="ds-mono" style={{ color: 'var(--ds-text)' }}>{totalDealt.kills}</strong>
        </span>
      </div>
    </div>
  )
}
