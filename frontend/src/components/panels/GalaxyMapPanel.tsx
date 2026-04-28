import { useState, useEffect, useMemo, useCallback, useRef } from 'react'
import { createPortal } from 'react-dom'
import {
  getGalaxySector,
  listPlanets,
  listFleets,
  attackPlanet,
  attackRBP,
  getPendingAttacks,
  getIncomingAttacks,
} from '../../services/api'
import LoadingButton from '../common/LoadingButton'
import type { Fleet, Planet, PendingAttack, IncomingAttack } from '../../types'
import type { SectorPlanet } from '../../types/corps'
import '../../styles/common.css'

interface GalaxyMapPanelProps {
  onClose: () => void
}

const VIEWPORT = 600 // SVG viewport size in px
const DEFAULT_RADIUS = 100

type LastResult = {
  ok: boolean
  message: string
}

export default function GalaxyMapPanel({ onClose }: GalaxyMapPanelProps) {
  const [homeworld, setHomeworld] = useState<Planet | null>(null)
  const [center, setCenter] = useState<{ x: number; y: number } | null>(null)
  const [radius, setRadius] = useState(DEFAULT_RADIUS)
  const [planets, setPlanets] = useState<SectorPlanet[]>([])
  const [selected, setSelected] = useState<SectorPlanet | null>(null)
  const [fleets, setFleets] = useState<Fleet[]>([])
  const [outgoing, setOutgoing] = useState<PendingAttack[]>([])
  const [incoming, setIncoming] = useState<IncomingAttack[]>([])
  const [selectedFleets, setSelectedFleets] = useState<string[]>([])
  const [loading, setLoading] = useState(false)
  const [attacking, setAttacking] = useState(false)
  const [lastResult, setLastResult] = useState<LastResult | null>(null)
  const dragRef = useRef<{ x: number; y: number } | null>(null)

  // ESC closes the panel.
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  // Bootstrap: load player's planets + fleets, center on the homeworld.
  useEffect(() => {
    let cancelled = false
    ;(async () => {
      try {
        const [ps, fs] = await Promise.all([listPlanets(), listFleets()])
        if (cancelled) return
        const home = ps.find((p) => p.is_homeworld) ?? ps[0]
        if (home) {
          setHomeworld(home)
          setCenter({ x: home.position_x, y: home.position_y })
        } else {
          setCenter({ x: 1000, y: 1000 })
        }
        setFleets(fs)
      } catch (e) {
        console.error('Galaxy bootstrap failed', e)
        setCenter({ x: 1000, y: 1000 })
      }
    })()
    return () => {
      cancelled = true
    }
  }, [])

  // Load the current sector whenever center/radius changes, plus the
  // outgoing/incoming attack streams so we can draw fleets in transit.
  useEffect(() => {
    if (!center) return
    let cancelled = false
    setLoading(true)
    Promise.all([
      getGalaxySector(center.x, center.y, radius),
      getPendingAttacks().catch(() => [] as PendingAttack[]),
      getIncomingAttacks().catch(() => ({ radar_level: 0, detect_advance: 0, incoming_attacks: [] })),
    ])
      .then(([sector, pending, radar]) => {
        if (cancelled) return
        setPlanets(sector.planets)
        setOutgoing(pending.filter((a) => a.status === 'traveling'))
        setIncoming(radar.incoming_attacks)
      })
      .catch((e) => console.error('Galaxy sector load failed', e))
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [center, radius])

  const refresh = useCallback(async () => {
    if (!center) return
    setLoading(true)
    try {
      const [sector, fs, pending, radar] = await Promise.all([
        getGalaxySector(center.x, center.y, radius),
        listFleets(),
        getPendingAttacks().catch(() => [] as PendingAttack[]),
        getIncomingAttacks().catch(() => ({ radar_level: 0, detect_advance: 0, incoming_attacks: [] })),
      ])
      setPlanets(sector.planets)
      setFleets(fs)
      setOutgoing(pending.filter((a) => a.status === 'traveling'))
      setIncoming(radar.incoming_attacks)
    } finally {
      setLoading(false)
    }
  }, [center, radius])

  const stationedFleets = useMemo(
    () => fleets.filter((f) => f.status === 'stationed' && (f.stacks?.length ?? 0) > 0),
    [fleets],
  )

  const toggleFleet = (id: string) =>
    setSelectedFleets((prev) =>
      prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id],
    )

  const dispatchAttack = async () => {
    if (!selected || selectedFleets.length === 0) return
    setAttacking(true)
    setLastResult(null)
    try {
      if (selected.is_rbp) {
        const res = await attackRBP(selected.id, { fleet_ids: selectedFleets })
        setLastResult({
          ok: res.result === 'attacker_win',
          message: `RBP ${res.result === 'attacker_win' ? 'defeated' : 'held'} in ${res.total_rounds} rounds${res.conquered ? ' — RBP conquered' : ''}`,
        })
      } else {
        const res = await attackPlanet({
          defender_planet_id: selected.id,
          fleet_ids: selectedFleets,
        })
        setLastResult({
          ok: true,
          message: `Fleet dispatched, arrives in ${res.travel_seconds}s (${res.fleets_dispatched} fleets)`,
        })
      }
      setSelectedFleets([])
      await refresh()
    } catch (e) {
      const err = e as { response?: { data?: { error?: string } } }
      setLastResult({
        ok: false,
        message: err.response?.data?.error ?? 'Attack failed',
      })
    } finally {
      setAttacking(false)
    }
  }

  // Convert world (cx,cy) to SVG (px,py).
  const project = useCallback(
    (x: number, y: number) => {
      if (!center) return { px: 0, py: 0 }
      const span = radius * 2
      const px = ((x - (center.x - radius)) / span) * VIEWPORT
      const py = ((y - (center.y - radius)) / span) * VIEWPORT
      return { px, py }
    },
    [center, radius],
  )

  // Pan via SVG drag.
  const onSvgMouseDown = (e: React.MouseEvent) => {
    dragRef.current = { x: e.clientX, y: e.clientY }
  }
  const onSvgMouseMove = (e: React.MouseEvent) => {
    if (!dragRef.current || !center) return
    const dx = e.clientX - dragRef.current.x
    const dy = e.clientY - dragRef.current.y
    if (Math.abs(dx) < 4 && Math.abs(dy) < 4) return
    dragRef.current = { x: e.clientX, y: e.clientY }
    const span = radius * 2
    setCenter({
      x: Math.round(center.x - (dx / VIEWPORT) * span),
      y: Math.round(center.y - (dy / VIEWPORT) * span),
    })
  }
  const onSvgMouseUp = () => {
    dragRef.current = null
  }

  const colorFor = (p: SectorPlanet) => {
    if (p.is_own) return '#4ade80'
    if (p.is_rbp) return '#f59e0b'
    return '#f87171'
  }

  const isProtected = (p: SectorPlanet) =>
    !!p.protection_until && new Date(p.protection_until).getTime() > Date.now()

  return createPortal(
    <div
      className="chat-backdrop"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div
        className="panel"
        style={{
          width: '1100px',
          maxWidth: '95vw',
          maxHeight: '90vh',
          display: 'flex',
          flexDirection: 'column',
          background: 'var(--bg-panel)',
          border: '1px solid var(--border-glow)',
          borderRadius: '8px',
        }}
      >
        <div className="panel-header">
          <h2>Galaxy Map</h2>
          <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
            <span style={{ color: '#aaa', fontSize: '0.85rem' }}>
              Center: ({center?.x ?? '?'}, {center?.y ?? '?'}) · radius {radius}
            </span>
            <button
              className="btn btn-small btn-secondary"
              onClick={() => setRadius((r) => Math.max(40, Math.round(r * 0.7)))}
              disabled={loading}
            >
              Zoom +
            </button>
            <button
              className="btn btn-small btn-secondary"
              onClick={() => setRadius((r) => Math.min(1000, Math.round(r * 1.4)))}
              disabled={loading}
            >
              Zoom −
            </button>
            <button
              className="btn btn-small btn-secondary"
              onClick={() => homeworld && setCenter({ x: homeworld.position_x, y: homeworld.position_y })}
              disabled={!homeworld}
            >
              Recenter
            </button>
            <LoadingButton
              className="btn btn-small btn-secondary"
              onClick={refresh}
              loading={loading}
            >
              Refresh
            </LoadingButton>
            <button className="p2-modal-close" onClick={onClose}>
              X
            </button>
          </div>
        </div>

        <div
          className="panel-content"
          style={{ flex: 1, display: 'grid', gridTemplateColumns: '1fr 320px', gap: 16, padding: 16, overflow: 'hidden' }}
        >
          <div
            style={{
              border: '1px solid #333',
              borderRadius: 6,
              background: 'radial-gradient(circle at center, #0a1124 0%, #02030a 75%)',
              position: 'relative',
              overflow: 'hidden',
              cursor: dragRef.current ? 'grabbing' : 'grab',
            }}
          >
            <svg
              viewBox={`0 0 ${VIEWPORT} ${VIEWPORT}`}
              width="100%"
              height="100%"
              onMouseDown={onSvgMouseDown}
              onMouseMove={onSvgMouseMove}
              onMouseUp={onSvgMouseUp}
              onMouseLeave={onSvgMouseUp}
              style={{ display: 'block', userSelect: 'none' }}
            >
              {/* grid */}
              {Array.from({ length: 11 }).map((_, i) => {
                const v = (i * VIEWPORT) / 10
                return (
                  <g key={i}>
                    <line x1={v} y1={0} x2={v} y2={VIEWPORT} stroke="#1a2240" strokeWidth={1} />
                    <line x1={0} y1={v} x2={VIEWPORT} y2={v} stroke="#1a2240" strokeWidth={1} />
                  </g>
                )
              })}
              {/* center crosshair */}
              <circle cx={VIEWPORT / 2} cy={VIEWPORT / 2} r={3} fill="#4a90d9" />

              {/* outgoing attacks: line from homeworld to defender_planet */}
              {homeworld &&
                outgoing.map((atk) => {
                  const target = planets.find((p) => p.id === atk.defender_planet_id)
                  if (!target) return null
                  const a = project(homeworld.position_x, homeworld.position_y)
                  const b = project(target.position_x, target.position_y)
                  return (
                    <g key={`out-${atk.id}`}>
                      <line
                        x1={a.px}
                        y1={a.py}
                        x2={b.px}
                        y2={b.py}
                        stroke="#fbbf24"
                        strokeWidth={1.4}
                        strokeDasharray="4 3"
                        opacity={0.85}
                      />
                      <circle cx={(a.px + b.px) / 2} cy={(a.py + b.py) / 2} r={3} fill="#fbbf24" />
                    </g>
                  )
                })}

              {/* incoming attacks (from radar): line from origin into homeworld */}
              {homeworld &&
                incoming.map((atk, i) => {
                  if (atk.origin_x == null || atk.origin_y == null) return null
                  const a = project(atk.origin_x, atk.origin_y)
                  const b = project(homeworld.position_x, homeworld.position_y)
                  return (
                    <g key={`in-${atk.id ?? i}`}>
                      <line
                        x1={a.px}
                        y1={a.py}
                        x2={b.px}
                        y2={b.py}
                        stroke="#ef4444"
                        strokeWidth={1.4}
                        strokeDasharray="4 3"
                        opacity={0.9}
                      />
                      <circle cx={(a.px + b.px) / 2} cy={(a.py + b.py) / 2} r={3} fill="#ef4444" />
                    </g>
                  )
                })}

              {/* planets */}
              {planets.map((p) => {
                const { px, py } = project(p.position_x, p.position_y)
                const r = p.is_rbp ? 8 : p.is_homeworld ? 10 : 7
                const fill = colorFor(p)
                const isSel = selected?.id === p.id
                return (
                  <g
                    key={p.id}
                    style={{ cursor: 'pointer' }}
                    onClick={(e) => {
                      e.stopPropagation()
                      setSelected(p)
                      setSelectedFleets([])
                      setLastResult(null)
                    }}
                  >
                    {isProtected(p) && (
                      <circle cx={px} cy={py} r={r + 6} fill="none" stroke="#4a90d9" strokeDasharray="2 3" />
                    )}
                    <circle cx={px} cy={py} r={r} fill={fill} stroke={isSel ? '#fff' : '#000'} strokeWidth={isSel ? 2 : 1} />
                    <text
                      x={px}
                      y={py - r - 4}
                      fontSize={10}
                      fill="#cbd5e1"
                      textAnchor="middle"
                      style={{ pointerEvents: 'none' }}
                    >
                      {p.name}
                    </text>
                  </g>
                )
              })}
            </svg>
            <div
              style={{
                position: 'absolute',
                bottom: 8,
                left: 8,
                fontSize: '0.7rem',
                color: '#94a3b8',
                background: 'rgba(0,0,0,0.5)',
                padding: '4px 8px',
                borderRadius: 3,
              }}
            >
              <span style={{ color: '#4ade80' }}>●</span> own ·{' '}
              <span style={{ color: '#f87171' }}>●</span> hostile ·{' '}
              <span style={{ color: '#f59e0b' }}>●</span> RBP ·{' '}
              <span style={{ color: '#fbbf24' }}>—</span> outgoing ·{' '}
              <span style={{ color: '#ef4444' }}>—</span> incoming · drag to pan
            </div>
            {loading && (
              <div
                style={{
                  position: 'absolute',
                  top: 8,
                  right: 8,
                  fontSize: '0.7rem',
                  color: '#94a3b8',
                }}
              >
                Loading…
              </div>
            )}
          </div>

          <aside
            style={{
              overflowY: 'auto',
              padding: 12,
              border: '1px solid #333',
              borderRadius: 6,
              background: 'rgba(0,0,0,0.3)',
            }}
          >
            {!selected && (
              <p style={{ color: '#94a3b8', fontSize: '0.85rem' }}>
                Click on a planet to view details and dispatch fleets.
              </p>
            )}
            {selected && (
              <>
                <h3 style={{ marginBottom: 8 }}>{selected.name}</h3>
                <div style={{ fontSize: '0.85rem', lineHeight: 1.6, color: '#cbd5e1' }}>
                  <div>
                    <strong>Position:</strong> ({selected.position_x}, {selected.position_y})
                  </div>
                  <div>
                    <strong>Type:</strong>{' '}
                    {selected.is_own
                      ? 'Your planet'
                      : selected.is_rbp
                        ? `RBP (lvl ${selected.rbp_level})`
                        : 'Enemy planet'}
                  </div>
                  {selected.owner_name && (
                    <div>
                      <strong>Owner:</strong> {selected.owner_name}
                    </div>
                  )}
                  {selected.controlling_corp && (
                    <div>
                      <strong>Corp:</strong> [{selected.controlling_corp.tag}]{' '}
                      {selected.controlling_corp.name}
                    </div>
                  )}
                  <div>
                    <strong>Defense strength:</strong> {selected.defense_strength}
                  </div>
                  {isProtected(selected) && selected.protection_until && (
                    <div style={{ color: '#4a90d9' }}>
                      <strong>Protected until:</strong>{' '}
                      {new Date(selected.protection_until).toLocaleString()}
                    </div>
                  )}
                </div>

                {!selected.is_own && !isProtected(selected) && (
                  <div style={{ marginTop: 16 }}>
                    <h4 style={{ marginBottom: 8 }}>Dispatch Fleets</h4>
                    {stationedFleets.length === 0 ? (
                      <p style={{ color: '#94a3b8', fontSize: '0.85rem' }}>
                        No stationed fleets ready. Build ships and create a fleet first.
                      </p>
                    ) : (
                      <div style={{ maxHeight: 220, overflowY: 'auto' }}>
                        {stationedFleets.map((f) => {
                          const ships = (f.stacks ?? []).reduce(
                            (sum, s) => sum + (s.ship_count ?? 0),
                            0,
                          )
                          return (
                            <label
                              key={f.id}
                              style={{
                                display: 'block',
                                padding: 8,
                                marginBottom: 4,
                                background: selectedFleets.includes(f.id)
                                  ? 'rgba(74,144,217,0.2)'
                                  : 'rgba(255,255,255,0.04)',
                                border: '1px solid #333',
                                borderRadius: 4,
                                cursor: 'pointer',
                                fontSize: '0.85rem',
                              }}
                            >
                              <input
                                type="checkbox"
                                checked={selectedFleets.includes(f.id)}
                                onChange={() => toggleFleet(f.id)}
                                style={{ marginRight: 8 }}
                              />
                              {f.name} — {ships} ships ({f.formation})
                            </label>
                          )
                        })}
                      </div>
                    )}
                    <div style={{ marginTop: 8 }}>
                      <LoadingButton
                        className="btn btn-primary"
                        onClick={dispatchAttack}
                        loading={attacking}
                        disabled={selectedFleets.length === 0}
                      >
                        {selected.is_rbp ? 'Attack RBP' : 'Attack Planet'}
                      </LoadingButton>
                    </div>
                  </div>
                )}

                {!selected.is_own && isProtected(selected) && (
                  <p style={{ color: '#4a90d9', fontSize: '0.85rem', marginTop: 12 }}>
                    Target is under truce protection.
                  </p>
                )}

                {lastResult && (
                  <div
                    style={{
                      marginTop: 12,
                      padding: 10,
                      borderRadius: 4,
                      background: lastResult.ok
                        ? 'rgba(74,222,128,0.15)'
                        : 'rgba(248,113,113,0.15)',
                      border: `1px solid ${lastResult.ok ? '#4ade80' : '#f87171'}`,
                      color: lastResult.ok ? '#4ade80' : '#f87171',
                      fontSize: '0.85rem',
                    }}
                  >
                    {lastResult.message}
                  </div>
                )}
              </>
            )}
          </aside>
        </div>
      </div>
    </div>,
    document.body,
  )
}
