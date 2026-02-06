import { useState, useMemo } from 'react'
import { useShipDesigns } from '../../hooks/useShipDesigns.ts'
import { useBlueprints } from '../../hooks/useBlueprints.ts'
import type { HullType, ModuleType, ShipDesignModule, ShipDesign } from '../../types'

const HULL_CLASS_COLORS: Record<string, string> = {
  frigate: '#44aaff',
  cruiser: '#cc8844',
  battleship: '#aa4444',
}

interface DesignEditorProps {
  hullTypes: HullType[]
  moduleTypes: ModuleType[]
  hasActivatedBlueprint: (id: number) => boolean
  onSave: (name: string, hullTypeId: number, modules: ShipDesignModule[]) => Promise<void>
  onClose: () => void
}

function DesignEditor({ hullTypes, moduleTypes, hasActivatedBlueprint, onSave, onClose }: DesignEditorProps) {
  const [name, setName] = useState('')
  const [selectedHull, setSelectedHull] = useState<number | null>(null)
  const [modules, setModules] = useState<ShipDesignModule[]>([])
  const [saving, setSaving] = useState(false)
  const [hullClassFilter, setHullClassFilter] = useState<string>('frigate')

  const hull = hullTypes.find(h => h.id === selectedHull)

  const filteredHulls = hullTypes.filter(h => h.hull_class === hullClassFilter)

  const volumeUsed = useMemo(() => {
    return modules.reduce((sum, m) => {
      const mt = moduleTypes.find(mt => mt.id === m.module_type_id)
      return sum + (mt ? mt.volume * m.quantity : 0)
    }, 0)
  }, [modules, moduleTypes])

  const totalStats = useMemo(() => {
    let shield = hull?.base_shield || 0
    let structure = hull?.base_structure || 0
    let attack = 0
    let he3 = 0

    for (const m of modules) {
      const mt = moduleTypes.find(mt => mt.id === m.module_type_id)
      if (!mt) continue
      shield += mt.shield_value * m.quantity
      structure += mt.structure_value * m.quantity
      attack += mt.attack_power * m.quantity
      he3 += mt.he3_per_round * m.quantity
    }

    return { shield, structure, attack, he3 }
  }, [hull, modules, moduleTypes])

  function addModule(moduleTypeId: number) {
    const existing = modules.find(m => m.module_type_id === moduleTypeId)
    if (existing) {
      setModules(modules.map(m =>
        m.module_type_id === moduleTypeId
          ? { ...m, quantity: m.quantity + 1 }
          : m
      ))
    } else {
      setModules([...modules, {
        module_type_id: moduleTypeId,
        quantity: 1,
        placement_order: modules.length + 1,
      }])
    }
  }

  function removeModule(moduleTypeId: number) {
    const existing = modules.find(m => m.module_type_id === moduleTypeId)
    if (!existing) return
    if (existing.quantity <= 1) {
      setModules(modules.filter(m => m.module_type_id !== moduleTypeId))
    } else {
      setModules(modules.map(m =>
        m.module_type_id === moduleTypeId
          ? { ...m, quantity: m.quantity - 1 }
          : m
      ))
    }
  }

  async function handleSave() {
    if (!name || !selectedHull || modules.length === 0) return
    setSaving(true)
    try {
      await onSave(name, selectedHull, modules)
      onClose()
    } catch {
      setSaving(false)
    }
  }

  const isValid = name.length > 0 && name.length <= 20 && /^[a-zA-Z0-9._-]+$/.test(name) &&
    selectedHull !== null && modules.length > 0 && hull && volumeUsed <= hull.installation_slots

  return (
    <div className="p2-modal-backdrop" onClick={e => { if (e.target === e.currentTarget) onClose() }}>
      <div className="p2-modal p2-modal-lg">
        <div className="p2-modal-header">
          <span className="p2-modal-title">New Ship Design</span>
          <button className="p2-modal-close" onClick={onClose}>X</button>
        </div>
        <div className="p2-modal-body design-editor">
          {/* Left: Hull Selection */}
          <div className="de-left">
            <div className="p2-form-group">
              <label className="p2-label">Design Name</label>
              <input
                className="p2-input"
                value={name}
                maxLength={20}
                placeholder="e.g. Frigate-Ballistic-V1"
                onChange={e => setName(e.target.value.replace(/[^a-zA-Z0-9._-]/g, ''))}
              />
            </div>

            <div className="de-class-tabs">
              {['frigate', 'cruiser', 'battleship'].map(cls => (
                <button
                  key={cls}
                  className={`p2-tab ${hullClassFilter === cls ? 'active' : ''}`}
                  style={{ borderColor: hullClassFilter === cls ? HULL_CLASS_COLORS[cls] : undefined }}
                  onClick={() => { setHullClassFilter(cls); setSelectedHull(null); setModules([]) }}
                >
                  {cls.charAt(0).toUpperCase() + cls.slice(1)}
                </button>
              ))}
            </div>

            <div className="de-hull-list">
              {filteredHulls.map(h => (
                <button
                  key={h.id}
                  className={`de-hull-card ${selectedHull === h.id ? 'selected' : ''}`}
                  onClick={() => { setSelectedHull(h.id); setModules([]) }}
                >
                  <div className="de-hull-name">{h.display_name}</div>
                  <div className="de-hull-stats">
                    <span>SH:{h.base_shield}</span>
                    <span>ST:{h.base_structure}</span>
                    <span>Slots:{h.installation_slots}</span>
                  </div>
                </button>
              ))}
            </div>
          </div>

          {/* Right: Module Selection & Stats */}
          <div className="de-right">
            {hull ? (
              <>
                <div className="de-stats-bar">
                  <div className="de-stat">
                    <span className="de-stat-label">Volume</span>
                    <span className={`de-stat-value ${volumeUsed > hull.installation_slots ? 'over' : ''}`}>
                      {volumeUsed} / {hull.installation_slots}
                    </span>
                  </div>
                  <div className="de-stat">
                    <span className="de-stat-label">Shield</span>
                    <span className="de-stat-value">{totalStats.shield}</span>
                  </div>
                  <div className="de-stat">
                    <span className="de-stat-label">Structure</span>
                    <span className="de-stat-value">{totalStats.structure}</span>
                  </div>
                  <div className="de-stat">
                    <span className="de-stat-label">Attack</span>
                    <span className="de-stat-value">{totalStats.attack}</span>
                  </div>
                  <div className="de-stat">
                    <span className="de-stat-label">He3/Rnd</span>
                    <span className="de-stat-value">{totalStats.he3}</span>
                  </div>
                </div>

                {/* Installed modules */}
                <div className="de-installed">
                  <div className="p2-section-title">Installed Modules</div>
                  {modules.length === 0 ? (
                    <div className="de-empty">No modules installed. Click modules below to add.</div>
                  ) : (
                    <div className="de-module-list">
                      {modules.map(m => {
                        const mt = moduleTypes.find(mt => mt.id === m.module_type_id)
                        if (!mt) return null
                        return (
                          <div key={m.module_type_id} className="de-module-installed">
                            <span className="de-mi-name">{mt.display_name}</span>
                            <span className="de-mi-qty">x{m.quantity}</span>
                            <span className="de-mi-vol">{mt.volume * m.quantity}v</span>
                            <button className="de-mi-remove" onClick={() => removeModule(m.module_type_id)}>-</button>
                            <button className="de-mi-add" onClick={() => addModule(m.module_type_id)}>+</button>
                          </div>
                        )
                      })}
                    </div>
                  )}
                </div>

                {/* Available modules */}
                <div className="de-available">
                  <div className="p2-section-title">Available Modules</div>
                  <div className="de-avail-grid">
                    {moduleTypes.map(mt => (
                      <button
                        key={mt.id}
                        className="de-avail-module"
                        onClick={() => addModule(mt.id)}
                        title={`${mt.display_name} | Vol:${mt.volume} | ATK:${mt.attack_power} | SH:${mt.shield_value}`}
                      >
                        <span className={`de-mod-cat ${mt.category}`}>{mt.category.slice(0, 3).toUpperCase()}</span>
                        <span className="de-mod-name">{mt.display_name}</span>
                        <span className="de-mod-vol">v{mt.volume}</span>
                      </button>
                    ))}
                  </div>
                </div>
              </>
            ) : (
              <div className="de-placeholder">Select a hull to begin designing</div>
            )}

            <button
              className="p2-btn p2-btn-primary p2-btn-full"
              disabled={!isValid || saving}
              onClick={handleSave}
            >
              {saving ? 'Saving...' : 'Save Design'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export default function ShipDesignPanel() {
  const { designs, hullTypes, moduleTypes, loading, error, create, remove } = useShipDesigns()
  const { hasActivated } = useBlueprints()
  const [showEditor, setShowEditor] = useState(false)

  if (loading) {
    return <div className="p2-panel-loading"><div className="loading-spinner" /><span>Loading Designs...</span></div>
  }

  async function handleSave(name: string, hullTypeId: number, modules: ShipDesignModule[]) {
    await create({ name, hull_type_id: hullTypeId, modules })
  }

  return (
    <div className="p2-panel">
      <div className="p2-panel-header">
        <div className="p2-panel-icon design-icon">DS</div>
        <div>
          <div className="p2-panel-title">Ship Designs</div>
          <div className="p2-panel-subtitle">{designs.length}/20 designs</div>
        </div>
        <button
          className="p2-btn p2-btn-primary"
          onClick={() => setShowEditor(true)}
          disabled={designs.length >= 20}
        >
          New Design
        </button>
      </div>

      {error && <div className="p2-error-msg">{error}</div>}

      <div className="p2-section">
        {designs.length === 0 ? (
          <div className="p2-empty-state">
            No ship designs yet. Create your first design to start building ships.
          </div>
        ) : (
          <div className="sd-list">
            {designs.map(d => (
              <DesignCard key={d.id} design={d} hullTypes={hullTypes} onDelete={() => remove(d.id)} />
            ))}
          </div>
        )}
      </div>

      {showEditor && (
        <DesignEditor
          hullTypes={hullTypes}
          moduleTypes={moduleTypes}
          hasActivatedBlueprint={hasActivated}
          onSave={handleSave}
          onClose={() => setShowEditor(false)}
        />
      )}
    </div>
  )
}

function DesignCard({ design, hullTypes, onDelete }: {
  design: ShipDesign
  hullTypes: HullType[]
  onDelete: () => void
}) {
  const hull = hullTypes.find(h => h.id === design.hull_type_id)
  const hullClass = design.hull_class || hull?.hull_class || 'frigate'
  const classColor = HULL_CLASS_COLORS[hullClass] || '#888'

  return (
    <div className="sd-card">
      <div className="sd-card-header">
        <div className="sd-card-class" style={{ color: classColor, borderColor: classColor }}>
          {hullClass.charAt(0).toUpperCase()}
        </div>
        <div className="sd-card-info">
          <div className="sd-card-name">{design.name}</div>
          <div className="sd-card-hull">{design.hull_name || hull?.display_name || 'Unknown Hull'}</div>
        </div>
        {design.ships_built === 0 && (
          <button className="p2-btn p2-btn-danger p2-btn-xs" onClick={onDelete}>Del</button>
        )}
      </div>
      <div className="sd-card-stats">
        <span className="sd-stat"><span className="sd-stat-label">SH</span> {design.total_shield}</span>
        <span className="sd-stat"><span className="sd-stat-label">ST</span> {design.total_structure}</span>
        <span className="sd-stat"><span className="sd-stat-label">ATK</span> {design.attack_power}</span>
        <span className="sd-stat"><span className="sd-stat-label">VOL</span> {design.volume_used}</span>
        <span className="sd-stat"><span className="sd-stat-label">Built</span> {design.ships_built}</span>
      </div>
    </div>
  )
}
