import { useState, useEffect, useMemo, useCallback, useRef } from 'react'
import { createPortal } from 'react-dom'
import {
  Globe,
  Telescope,
  Radar,
  MapPin,
  Crosshair,
  Move,
  Eye,
  X,
} from 'lucide-react'
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

  // Sector fill color per ownership / type.
  const colorFor = (p: SectorPlanet) => {
    if (p.is_own) return 'var(--ds-owner-self)'
    if (p.is_rbp) return 'var(--ds-orange)'
    return 'var(--ds-owner-enemy)'
  }

  const isProtected = (p: SectorPlanet) =>
    !!p.protection_until && new Date(p.protection_until).getTime() > Date.now()

  return createPortal(
    <div
      className="ds-modal-backdrop"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div
        className="ds-modal galaxy-modal"
        role="dialog"
        aria-label="Galaxy Map"
        style={{
          maxWidth: 'min(1100px, 95vw)',
          maxHeight: '90vh',
          width: '100%',
        }}
      >
        <header className="ds-modal-header">
          <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--sp-3)' }}>
            <Globe size={22} strokeWidth={1.75} aria-hidden="true" color="var(--ds-teal)" />
            <h2 className="ds-modal-title">Galaxy Map</h2>
          </div>
          <div style={{ display: 'flex', gap: 'var(--sp-2)', alignItems: 'center' }}>
            <span
              className="ds-mono ds-text-muted"
              style={{ fontSize: 'var(--fs-sm)' }}
              aria-label="Viewport coordinates"
            >
              ({center?.x ?? '?'}, {center?.y ?? '?'}) · r {radius}
            </span>
            <button
              className="ds-btn ds-btn-ghost ds-btn--sm"
              onClick={() => setRadius((r) => Math.max(40, Math.round(r * 0.7)))}
              disabled={loading}
              title="Zoom in (smaller radius)"
            >
              <Telescope size={14} strokeWidth={2} aria-hidden="true" />
              Zoom +
            </button>
            <button
              className="ds-btn ds-btn-ghost ds-btn--sm"
              onClick={() => setRadius((r) => Math.min(1000, Math.round(r * 1.4)))}
              disabled={loading}
              title="Zoom out (larger radius)"
            >
              <Telescope size={14} strokeWidth={2} aria-hidden="true" />
              Zoom −
            </button>
            <button
              className="ds-btn ds-btn-ghost ds-btn--sm"
              onClick={() => homeworld && setCenter({ x: homeworld.position_x, y: homeworld.position_y })}
              disabled={!homeworld}
              title="Recenter on homeworld"
            >
              <Crosshair size={14} strokeWidth={2} aria-hidden="true" />
              Recenter
            </button>
            <LoadingButton
              className="ds-btn ds-btn-ghost ds-btn--sm"
              onClick={refresh}
              loading={loading}
            >
              <Radar size={14} strokeWidth={2} aria-hidden="true" />
              Refresh
            </LoadingButton>
            <button
              className="ds-btn-icon"
              onClick={onClose}
              aria-label="Close galaxy map"
              title="Close"
            >
              <X size={18} strokeWidth={1.75} aria-hidden="true" />
            </button>
          </div>
        </header>

        <div
          className="galaxy-modal-body"
          style={{
            flex: 1,
            display: 'grid',
            gridTemplateColumns: '1fr 340px',
            gap: 'var(--sp-4)',
            padding: 'var(--sp-5)',
            overflow: 'hidden',
            background: 'var(--ds-surface-2)',
          }}
        >
          {/* ---- Galaxy canvas (deep-space dark exception) ---- */}
          <div
            style={{
              border: '1px solid var(--ds-border)',
              borderRadius: 'var(--r-lg)',
              background: 'radial-gradient(circle at center, #0F172A 0%, #020617 75%)',
              position: 'relative',
              overflow: 'hidden',
              cursor: dragRef.current ? 'grabbing' : 'grab',
              boxShadow: 'var(--shadow-1)',
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
              {/* sector grid */}
              {Array.from({ length: 11 }).map((_, i) => {
                const v = (i * VIEWPORT) / 10
                return (
                  <g key={i}>
                    <line x1={v} y1={0} x2={v} y2={VIEWPORT} stroke="#1E293B" strokeWidth={1} />
                    <line x1={0} y1={v} x2={VIEWPORT} y2={v} stroke="#1E293B" strokeWidth={1} />
                  </g>
                )
              })}
              {/* center crosshair (teal identity) */}
              <circle cx={VIEWPORT / 2} cy={VIEWPORT / 2} r={3} fill="var(--ds-teal)" />

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
                        stroke="var(--ds-orange)"
                        strokeWidth={1.4}
                        strokeDasharray="4 3"
                        opacity={0.85}
                      />
                      <circle
                        cx={(a.px + b.px) / 2}
                        cy={(a.py + b.py) / 2}
                        r={3}
                        fill="var(--ds-orange)"
                      />
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
                        stroke="var(--ds-owner-enemy)"
                        strokeWidth={1.4}
                        strokeDasharray="4 3"
                        opacity={0.9}
                      />
                      <circle
                        cx={(a.px + b.px) / 2}
                        cy={(a.py + b.py) / 2}
                        r={3}
                        fill="var(--ds-owner-enemy)"
                      />
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
                      <circle
                        cx={px}
                        cy={py}
                        r={r + 6}
                        fill="none"
                        stroke="var(--ds-info)"
                        strokeDasharray="2 3"
                      />
                    )}
                    <circle
                      cx={px}
                      cy={py}
                      r={r}
                      fill={fill}
                      stroke={isSel ? 'var(--ds-orange)' : '#0F172A'}
                      strokeWidth={isSel ? 2.5 : 1}
                    />
                    <text
                      x={px}
                      y={py - r - 4}
                      fontSize={10}
                      fill="#E2E8F0"
                      textAnchor="middle"
                      style={{ pointerEvents: 'none', fontFamily: 'var(--ff-sans)' }}
                    >
                      {p.name}
                    </text>
                  </g>
                )
              })}
            </svg>

            {/* Legend (uses ds-tooltip dark style — fits on the deep-space bg) */}
            <div
              className="ds-tooltip"
              style={{
                position: 'absolute',
                bottom: 'var(--sp-2)',
                left: 'var(--sp-2)',
                display: 'flex',
                gap: 'var(--sp-3)',
                alignItems: 'center',
                whiteSpace: 'nowrap',
              }}
            >
              <span style={{ display: 'inline-flex', alignItems: 'center', gap: 4 }}>
                <span
                  aria-hidden="true"
                  style={{
                    width: 8,
                    height: 8,
                    borderRadius: '50%',
                    background: 'var(--ds-owner-self)',
                    display: 'inline-block',
                  }}
                />
                own
              </span>
              <span style={{ display: 'inline-flex', alignItems: 'center', gap: 4 }}>
                <span
                  aria-hidden="true"
                  style={{
                    width: 8,
                    height: 8,
                    borderRadius: '50%',
                    background: 'var(--ds-owner-enemy)',
                    display: 'inline-block',
                  }}
                />
                hostile
              </span>
              <span style={{ display: 'inline-flex', alignItems: 'center', gap: 4 }}>
                <span
                  aria-hidden="true"
                  style={{
                    width: 8,
                    height: 8,
                    borderRadius: '50%',
                    background: 'var(--ds-orange)',
                    display: 'inline-block',
                  }}
                />
                RBP
              </span>
              <span style={{ display: 'inline-flex', alignItems: 'center', gap: 4 }}>
                <Move size={12} strokeWidth={2} aria-hidden="true" />
                drag to pan
              </span>
            </div>

            {loading && (
              <div
                className="ds-tooltip"
                style={{
                  position: 'absolute',
                  top: 'var(--sp-2)',
                  right: 'var(--sp-2)',
                }}
              >
                Loading…
              </div>
            )}
          </div>

          {/* ---- Right-side sector info panel (LIGHT theme) ---- */}
          <aside
            className="ds-panel"
            style={{
              overflowY: 'auto',
              padding: 'var(--sp-4)',
              display: 'flex',
              flexDirection: 'column',
              gap: 'var(--sp-3)',
            }}
            aria-label="Sector details"
          >
            {!selected && (
              <div
                style={{
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  justifyContent: 'center',
                  gap: 'var(--sp-2)',
                  textAlign: 'center',
                  padding: 'var(--sp-6) var(--sp-2)',
                  color: 'var(--ds-text-muted)',
                }}
              >
                <Telescope
                  size={32}
                  strokeWidth={1.5}
                  aria-hidden="true"
                  color="var(--ds-text-soft)"
                />
                <p style={{ margin: 0, fontSize: 'var(--fs-sm)' }}>
                  Click on a planet to view details and dispatch fleets.
                </p>
              </div>
            )}
            {selected && (
              <>
                <div
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 'var(--sp-2)',
                    paddingBottom: 'var(--sp-3)',
                    borderBottom: '1px solid var(--ds-border)',
                  }}
                >
                  <MapPin
                    size={20}
                    strokeWidth={1.75}
                    aria-hidden="true"
                    color="var(--ds-teal)"
                  />
                  <h3 className="ds-h3" style={{ margin: 0 }}>
                    {selected.name}
                  </h3>
                  {selected.is_own && (
                    <span className="ds-badge ds-badge--teal" style={{ marginLeft: 'auto' }}>
                      You
                    </span>
                  )}
                  {!selected.is_own && selected.is_rbp && (
                    <span className="ds-badge ds-badge--orange" style={{ marginLeft: 'auto' }}>
                      RBP
                    </span>
                  )}
                  {!selected.is_own && !selected.is_rbp && (
                    <span className="ds-badge ds-badge--danger" style={{ marginLeft: 'auto' }}>
                      Enemy
                    </span>
                  )}
                </div>

                <div
                  style={{
                    display: 'grid',
                    gridTemplateColumns: 'auto 1fr',
                    gap: '6px var(--sp-3)',
                    fontSize: 'var(--fs-sm)',
                    color: 'var(--ds-text)',
                  }}
                >
                  <span className="ds-text-muted">Position</span>
                  <span className="ds-mono">
                    ({selected.position_x}, {selected.position_y})
                  </span>

                  <span className="ds-text-muted">Type</span>
                  <span>
                    {selected.is_own
                      ? 'Your planet'
                      : selected.is_rbp
                        ? `RBP (lvl ${selected.rbp_level})`
                        : 'Enemy planet'}
                  </span>

                  {selected.owner_name && (
                    <>
                      <span className="ds-text-muted">Owner</span>
                      <span>{selected.owner_name}</span>
                    </>
                  )}

                  {selected.controlling_corp && (
                    <>
                      <span className="ds-text-muted">Corp</span>
                      <span>
                        [{selected.controlling_corp.tag}] {selected.controlling_corp.name}
                      </span>
                    </>
                  )}

                  <span className="ds-text-muted">Defense strength</span>
                  <span className="ds-mono">{selected.defense_strength}</span>

                  {isProtected(selected) && selected.protection_until && (
                    <>
                      <span className="ds-text-muted">Protected until</span>
                      <span style={{ color: 'var(--ds-info)' }}>
                        {new Date(selected.protection_until).toLocaleString()}
                      </span>
                    </>
                  )}
                </div>

                {!selected.is_own && !isProtected(selected) && (
                  <div
                    style={{
                      paddingTop: 'var(--sp-3)',
                      borderTop: '1px solid var(--ds-border)',
                      display: 'flex',
                      flexDirection: 'column',
                      gap: 'var(--sp-2)',
                    }}
                  >
                    <h4
                      className="ds-caption"
                      style={{
                        margin: 0,
                        display: 'flex',
                        alignItems: 'center',
                        gap: 'var(--sp-2)',
                      }}
                    >
                      <Crosshair size={14} strokeWidth={2} aria-hidden="true" />
                      Dispatch Fleets
                    </h4>
                    {stationedFleets.length === 0 ? (
                      <p
                        className="ds-text-muted"
                        style={{ margin: 0, fontSize: 'var(--fs-sm)' }}
                      >
                        No stationed fleets ready. Build ships and create a fleet first.
                      </p>
                    ) : (
                      <div
                        style={{
                          maxHeight: 220,
                          overflowY: 'auto',
                          display: 'flex',
                          flexDirection: 'column',
                          gap: 'var(--sp-1)',
                        }}
                      >
                        {stationedFleets.map((f) => {
                          const ships = (f.stacks ?? []).reduce(
                            (sum, s) => sum + (s.ship_count ?? 0),
                            0,
                          )
                          const isChecked = selectedFleets.includes(f.id)
                          return (
                            <label
                              key={f.id}
                              className="ds-list-item"
                              aria-selected={isChecked}
                              style={{
                                display: 'flex',
                                alignItems: 'center',
                                gap: 'var(--sp-2)',
                                fontSize: 'var(--fs-sm)',
                                padding: 'var(--sp-2) var(--sp-3)',
                              }}
                            >
                              <input
                                type="checkbox"
                                checked={isChecked}
                                onChange={() => toggleFleet(f.id)}
                                style={{ accentColor: 'var(--ds-teal)' }}
                              />
                              <span style={{ fontWeight: 600 }}>{f.name}</span>
                              <span className="ds-text-muted ds-mono" style={{ marginLeft: 'auto' }}>
                                {ships} ships
                              </span>
                              <span
                                className="ds-badge ds-badge--neutral"
                                style={{ marginLeft: 'var(--sp-1)' }}
                              >
                                {f.formation}
                              </span>
                            </label>
                          )
                        })}
                      </div>
                    )}
                    <LoadingButton
                      className="ds-btn ds-btn-primary ds-btn--block"
                      onClick={dispatchAttack}
                      loading={attacking}
                      disabled={selectedFleets.length === 0}
                    >
                      <Crosshair size={14} strokeWidth={2} aria-hidden="true" />
                      {selected.is_rbp ? 'Attack RBP' : 'Attack Planet'}
                    </LoadingButton>
                  </div>
                )}

                {!selected.is_own && isProtected(selected) && (
                  <div
                    className="ds-badge ds-badge--info"
                    style={{
                      width: '100%',
                      justifyContent: 'center',
                      padding: 'var(--sp-2)',
                    }}
                  >
                    <Eye size={14} strokeWidth={2} aria-hidden="true" />
                    Target is under truce protection.
                  </div>
                )}

                {lastResult && (
                  <div
                    className={`ds-badge ${lastResult.ok ? 'ds-badge--success' : 'ds-badge--danger'}`}
                    role="status"
                    style={{
                      width: '100%',
                      justifyContent: 'flex-start',
                      whiteSpace: 'normal',
                      padding: 'var(--sp-2) var(--sp-3)',
                      lineHeight: 1.4,
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
