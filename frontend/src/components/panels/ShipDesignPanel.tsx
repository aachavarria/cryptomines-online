import { useState, useMemo, Suspense } from 'react'
import { createPortal } from 'react-dom'
import { Canvas } from '@react-three/fiber'
import { OrbitControls, useGLTF, Environment } from '@react-three/drei'
import { Wrench, Cog, Plus, Minus, Trash2, X, Lock } from 'lucide-react'
import { useShipDesigns } from '../../hooks/useShipDesigns.ts'
import { useBlueprints } from '../../hooks/useBlueprints.ts'
import LoadingButton from '../common/LoadingButton.tsx'
import type { HullType, ModuleType, ShipDesignModule, ShipDesign } from '../../types'

const HULL_CLASS_TINT: Record<string, string> = {
  frigate: 'var(--ds-info)',
  cruiser: 'var(--ds-orange-strong)',
  battleship: 'var(--ds-danger)',
}

const HULL_CLASS_MODELS: Record<string, string> = {
  frigate: '/assets/gbl/ships/nave1.glb',
  cruiser: '/assets/gbl/ships/nave3.glb',
  battleship: '/assets/gbl/ships/nave5.glb',
}

const MODULE_GROUPS: Record<string, { label: string; categories: string[] }> = {
  attack: { label: 'Attack', categories: ['ballistic', 'directional', 'missile', 'ship_based', 'planetary'] },
  defense: { label: 'Defense', categories: ['structure', 'shield', 'air_defense'] },
  auxiliary: { label: 'Auxiliary', categories: ['electronic', 'storage', 'transmission'] },
}

const CATEGORY_LABELS: Record<string, string> = {
  ballistic: 'BAL', directional: 'DIR', missile: 'MSL', ship_based: 'SHP', planetary: 'PLN',
  structure: 'STR', shield: 'SHD', air_defense: 'ADF',
  electronic: 'ELC', storage: 'STO', transmission: 'TRN',
}

const TIER_LABELS: Record<number, string> = { 1: 'I', 2: 'II', 3: 'III' }

function getTierLabel(tier: number): string {
  return TIER_LABELS[tier] || `T${tier}`
}

function getModuleAvgAttack(mt: ModuleType): number {
  return Math.round((mt.min_damage + mt.max_damage) / 2)
}

function parseEffects(mt: ModuleType): Record<string, number> {
  try {
    return JSON.parse(mt.effects_json || '{}')
  } catch {
    return {}
  }
}

function getModuleStat(mt: ModuleType): string {
  if (mt.min_damage > 0 || mt.max_damage > 0) {
    return `${mt.min_damage}-${mt.max_damage} dmg`
  }
  const effects = parseEffects(mt)
  const entries = Object.entries(effects)
  if (entries.length === 0) return ''
  const [key, val] = entries[0]
  const label = key.replace(/_/g, ' ').replace(/\bbonus\b/i, '').replace(/\bpct\b/i, '%').trim()
  return `+${val} ${label}`
}

// 3D Ship Preview ----------------------------------------------------------

function ShipModel({ url }: { url: string }) {
  const { scene } = useGLTF(url)
  return <primitive object={scene.clone()} scale={0.1} position={[0, 0, 0]} />
}

function ShipPreview({ hullClass }: { hullClass: string }) {
  const modelUrl = HULL_CLASS_MODELS[hullClass] || HULL_CLASS_MODELS.frigate
  return (
    <div className="de-preview">
      <Canvas camera={{ position: [0, 5, 25], fov: 50 }}>
        <ambientLight intensity={0.4} />
        <pointLight position={[5, 5, 5]} intensity={1} />
        <pointLight position={[-3, 2, -3]} intensity={0.4} color="#4488ff" />
        <Suspense fallback={null}>
          <ShipModel url={modelUrl} />
          <Environment preset="night" />
        </Suspense>
        <OrbitControls autoRotate autoRotateSpeed={1.5} enableZoom={false} enablePan={false} />
      </Canvas>
    </div>
  )
}

// Design Editor ------------------------------------------------------------

interface DesignEditorProps {
  hullTypes: HullType[]
  moduleTypes: ModuleType[]
  hasHullBlueprint: (id: number) => boolean
  hasModuleBlueprint: (id: number) => boolean
  getHullBlueprintResearchLevel: (hullTypeId: number) => number
  getModuleBlueprintResearchLevel: (moduleTypeId: number) => number
  onSave: (name: string, hullTypeId: number, modules: ShipDesignModule[]) => Promise<void>
  onClose: () => void
}

function DesignEditor({
  hullTypes,
  moduleTypes,
  hasHullBlueprint,
  hasModuleBlueprint,
  getHullBlueprintResearchLevel,
  getModuleBlueprintResearchLevel,
  onSave,
  onClose,
}: DesignEditorProps) {
  const [name, setName] = useState('')
  const [selectedHull, setSelectedHull] = useState<number | null>(null)
  const [modules, setModules] = useState<ShipDesignModule[]>([])
  const [saving, setSaving] = useState(false)
  const [hullClassFilter, setHullClassFilter] = useState<string>('frigate')
  const [moduleGroup, setModuleGroup] = useState<string>('attack')
  const [moduleSubCategory, setModuleSubCategory] = useState<string | null>(null)

  const hull = hullTypes.find(h => h.id === selectedHull)
  const filteredHulls = hullTypes.filter(h => {
    if (h.hull_class !== hullClassFilter) return false
    if (!hasHullBlueprint(h.id)) return false
    const lvl = getHullBlueprintResearchLevel(h.id)
    return lvl >= h.tier
  })

  const groupCategories = MODULE_GROUPS[moduleGroup]?.categories || []
  const filteredModules = useMemo(() => {
    return moduleTypes.filter(mt => {
      if (!groupCategories.includes(mt.category)) return false
      if (moduleSubCategory && mt.category !== moduleSubCategory) return false
      if (!hasModuleBlueprint(mt.id)) return false
      const lvl = getModuleBlueprintResearchLevel(mt.id)
      return lvl >= mt.tier
    })
  }, [moduleTypes, groupCategories, moduleSubCategory, hasModuleBlueprint, getModuleBlueprintResearchLevel])

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
    let agility = hull?.base_agility || 0
    let storage = hull?.base_storage || 0
    let defense = hull?.base_defense || 0
    let stability = hull?.base_stability || 0
    let mobility = hull?.base_movement || 0
    let metalCost = hull?.base_metal_cost || 0
    let he3Cost = hull?.base_he3_cost || 0
    let goldCost = hull?.base_gold_cost || 0

    for (const m of modules) {
      const mt = moduleTypes.find(mt => mt.id === m.module_type_id)
      if (!mt) continue
      attack += getModuleAvgAttack(mt) * m.quantity
      he3 += mt.he3_per_round * m.quantity
      metalCost += mt.metal_cost * m.quantity
      he3Cost += mt.he3_cost * m.quantity
      goldCost += mt.gold_cost * m.quantity
      const effects = parseEffects(mt)
      if (effects.shield_bonus) shield += effects.shield_bonus * m.quantity
      if (effects.structure_bonus) structure += effects.structure_bonus * m.quantity
      if (effects.agility_bonus) agility += effects.agility_bonus * m.quantity
      if (effects.he3_storage_bonus) storage += effects.he3_storage_bonus * m.quantity
      if (effects.defense_bonus_pct) defense += effects.defense_bonus_pct * m.quantity
      if (effects.movement_bonus) mobility += effects.movement_bonus * m.quantity
    }

    return { shield, structure, attack, he3, agility, storage, defense, stability, mobility, metalCost, he3Cost, goldCost }
  }, [hull, modules, moduleTypes])

  const volumeMax = hull?.installation_slots || 0
  const volumePct = volumeMax > 0 ? Math.min(100, (volumeUsed / volumeMax) * 100) : 0
  const isOverVolume = volumeUsed > volumeMax
  const isNearCapacity = !isOverVolume && volumePct >= 90

  function addModule(moduleTypeId: number) {
    const mt = moduleTypes.find(m => m.id === moduleTypeId)
    if (!mt) return
    const existing = modules.find(m => m.module_type_id === moduleTypeId)
    const currentQty = existing?.quantity || 0
    if (mt.max_per_ship > 0 && currentQty >= mt.max_per_ship) return
    if (existing) {
      setModules(modules.map(m =>
        m.module_type_id === moduleTypeId ? { ...m, quantity: m.quantity + 1 } : m,
      ))
    } else {
      setModules([
        ...modules,
        { module_type_id: moduleTypeId, quantity: 1, placement_order: modules.length + 1 },
      ])
    }
  }

  function removeModule(moduleTypeId: number) {
    const existing = modules.find(m => m.module_type_id === moduleTypeId)
    if (!existing) return
    if (existing.quantity <= 1) {
      setModules(modules.filter(m => m.module_type_id !== moduleTypeId))
    } else {
      setModules(modules.map(m =>
        m.module_type_id === moduleTypeId ? { ...m, quantity: m.quantity - 1 } : m,
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

  const isValid =
    name.length > 0 &&
    name.length <= 20 &&
    /^[a-zA-Z0-9._-]+$/.test(name) &&
    selectedHull !== null &&
    modules.length > 0 &&
    hull &&
    !isOverVolume

  return createPortal(
    <div
      className="ds-modal-backdrop de-backdrop"
      onClick={e => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div className="ds-panel ds-panel--flush de-shell" role="dialog" aria-modal="true">
        <header className="de-header">
          <div className="de-title">
            <Wrench size={22} strokeWidth={1.75} />
            <h2 className="ds-h2">Ship Design Editor</h2>
          </div>
          <button className="ds-btn-icon" onClick={onClose} aria-label="Close editor">
            <X size={18} strokeWidth={1.75} />
          </button>
        </header>

        <div className="de-body">
          {/* Column 1: Hull Selection */}
          <section className="de-col de-col--hull">
            <div className="ds-caption">Select Hull</div>
            <div className="ds-tabs de-tabs">
              {(['frigate', 'cruiser', 'battleship'] as const).map(cls => (
                <button
                  key={cls}
                  className="ds-tab"
                  aria-selected={hullClassFilter === cls}
                  onClick={() => {
                    setHullClassFilter(cls)
                    setSelectedHull(null)
                    setModules([])
                  }}
                >
                  {cls.charAt(0).toUpperCase() + cls.slice(1)}
                </button>
              ))}
            </div>

            <div className="de-hull-list">
              {filteredHulls.map(h => {
                const hasBp = hasHullBlueprint(h.id)
                const lvl = getHullBlueprintResearchLevel(h.id)
                const tierUnlocked = lvl >= h.tier
                const isSelected = selectedHull === h.id
                return (
                  <button
                    key={h.id}
                    type="button"
                    className={`ds-list-item de-hull-card ${isSelected ? 'is-selected' : ''} ${!tierUnlocked ? 'is-disabled' : ''}`}
                    aria-selected={isSelected}
                    onClick={() => {
                      if (tierUnlocked) {
                        setSelectedHull(h.id)
                        setModules([])
                      }
                    }}
                    disabled={!tierUnlocked}
                    title={
                      !tierUnlocked
                        ? `Tier ${getTierLabel(h.tier)} locked - Research blueprint to level ${h.tier}`
                        : ''
                    }
                  >
                    <div className="de-hull-top">
                      <span className="de-hull-name">{h.display_name}</span>
                      <span className="ds-badge ds-badge--neutral">{getTierLabel(h.tier)}</span>
                    </div>
                    <div className="de-hull-stats ds-mono">
                      <span>SH:{h.base_shield}</span>
                      <span>ST:{h.base_structure}</span>
                      <span>Slots:{h.installation_slots}</span>
                    </div>
                    {!hasBp && (
                      <span className="ds-badge ds-badge--neutral de-lock">
                        <Lock size={12} strokeWidth={2} /> No Blueprint
                      </span>
                    )}
                    {hasBp && !tierUnlocked && (
                      <span className="ds-badge ds-badge--warning de-lock">
                        <Lock size={12} strokeWidth={2} /> Tier {getTierLabel(h.tier)} Locked
                      </span>
                    )}
                  </button>
                )
              })}
              {filteredHulls.length === 0 && (
                <div className="de-empty ds-text-muted">No hulls of this class unlocked.</div>
              )}
            </div>
          </section>

          {/* Column 2: Preview & Design */}
          <section className="de-col de-col--center">
            {hull ? (
              <>
                <ShipPreview hullClass={hull.hull_class} />

                <div className="de-form-group">
                  <label className="ds-caption">Design Name</label>
                  <input
                    className="ds-input"
                    value={name}
                    maxLength={20}
                    placeholder="e.g. Frigate-Ballistic-V1"
                    onChange={e => setName(e.target.value.replace(/[^a-zA-Z0-9._-]/g, ''))}
                  />
                </div>

                <div className="de-installed">
                  <div className="de-installed-head">
                    <span className="ds-caption">Installed Modules</span>
                    <span className="ds-badge ds-badge--neutral ds-mono">
                      {modules.reduce((s, m) => s + m.quantity, 0)} items
                    </span>
                  </div>
                  {modules.length === 0 ? (
                    <div className="de-empty ds-text-muted">
                      Select modules from the right panel.
                    </div>
                  ) : (
                    <div className="de-module-list">
                      {modules.map(m => {
                        const mt = moduleTypes.find(mt => mt.id === m.module_type_id)
                        if (!mt) return null
                        return (
                          <div key={m.module_type_id} className="ds-card de-mi">
                            <span className="ds-badge ds-badge--neutral ds-mono">
                              {CATEGORY_LABELS[mt.category] || mt.category.slice(0, 3).toUpperCase()}
                            </span>
                            <span className="de-mi-name">{mt.display_name}</span>
                            <span className="ds-mono ds-text-muted">×{m.quantity}</span>
                            <span className="ds-mono ds-text-muted">{mt.volume * m.quantity}v</span>
                            <button
                              className="ds-btn-icon ds-btn-icon--sm"
                              onClick={() => removeModule(m.module_type_id)}
                              aria-label={`Remove ${mt.display_name}`}
                            >
                              <Minus size={14} strokeWidth={2} />
                            </button>
                            <button
                              className="ds-btn-icon ds-btn-icon--sm"
                              onClick={() => addModule(m.module_type_id)}
                              aria-label={`Add ${mt.display_name}`}
                            >
                              <Plus size={14} strokeWidth={2} />
                            </button>
                          </div>
                        )
                      })}
                    </div>
                  )}
                </div>

                <div className="de-volume">
                  <div className="de-volume-row">
                    <span className="ds-caption">Volume</span>
                    <span className={`ds-mono${isOverVolume ? ' de-volume-over' : ''}`}>
                      {volumeUsed} / {volumeMax}
                    </span>
                  </div>
                  <div className="ds-bar">
                    <div
                      className={`ds-bar-fill${isOverVolume ? ' ds-bar-fill--danger' : isNearCapacity ? ' ds-bar-fill--warning' : ''}`}
                      style={{ width: `${Math.min(100, volumePct)}%` }}
                    />
                  </div>
                </div>

                <div className="de-actions">
                  <button className="ds-btn ds-btn-ghost" onClick={onClose}>
                    Cancel
                  </button>
                  <LoadingButton
                    className="ds-btn ds-btn-secondary"
                    disabled={!isValid}
                    loading={saving}
                    onClick={handleSave}
                  >
                    Save Design
                  </LoadingButton>
                </div>
              </>
            ) : (
              <div className="de-placeholder ds-text-muted">
                Select a hull to begin designing.
              </div>
            )}
          </section>

          {/* Column 3: Module Selection */}
          <section className="de-col de-col--modules">
            <div className="ds-caption">Modules</div>

            <div className="ds-tabs de-tabs">
              {Object.entries(MODULE_GROUPS).map(([key, group]) => (
                <button
                  key={key}
                  className="ds-tab"
                  aria-selected={moduleGroup === key}
                  onClick={() => {
                    setModuleGroup(key)
                    setModuleSubCategory(null)
                  }}
                >
                  {group.label}
                </button>
              ))}
            </div>

            <div className="de-sub-tabs">
              <button
                className={`ds-badge ${moduleSubCategory === null ? 'ds-badge--teal' : 'ds-badge--neutral'} de-sub-tab`}
                onClick={() => setModuleSubCategory(null)}
              >
                All
              </button>
              {groupCategories.map(cat => (
                <button
                  key={cat}
                  className={`ds-badge ${moduleSubCategory === cat ? 'ds-badge--teal' : 'ds-badge--neutral'} de-sub-tab`}
                  onClick={() => setModuleSubCategory(cat)}
                >
                  {CATEGORY_LABELS[cat] || cat.slice(0, 3).toUpperCase()}
                </button>
              ))}
            </div>

            <div className="de-catalog">
              {filteredModules.length === 0 ? (
                <div className="de-empty ds-text-muted">
                  No modules unlocked in this category.
                </div>
              ) : (
                filteredModules.map(mt => {
                  const hasBp = hasModuleBlueprint(mt.id)
                  const lvl = getModuleBlueprintResearchLevel(mt.id)
                  const tierUnlocked = lvl >= mt.tier
                  const installed = modules.find(m => m.module_type_id === mt.id)
                  const atMax = mt.max_per_ship > 0 && (installed?.quantity || 0) >= mt.max_per_ship
                  const disabled = !tierUnlocked || !hull || atMax
                  return (
                    <button
                      key={mt.id}
                      type="button"
                      className={`ds-list-item de-catalog-mod ${installed ? 'is-selected' : ''} ${disabled ? 'is-disabled' : ''}`}
                      aria-selected={!!installed}
                      disabled={disabled}
                      onClick={() => {
                        if (tierUnlocked && hull && !atMax) addModule(mt.id)
                      }}
                      title={
                        !tierUnlocked
                          ? `Tier ${getTierLabel(mt.tier)} locked - Research blueprint to level ${mt.tier}`
                          : ''
                      }
                    >
                      <span className="ds-badge ds-badge--neutral ds-mono">
                        {CATEGORY_LABELS[mt.category] || mt.category.slice(0, 3).toUpperCase()}
                      </span>
                      <div className="de-catalog-info">
                        <span className="de-catalog-name">
                          {mt.display_name}
                          <span className="ds-badge ds-badge--neutral de-tier">
                            {getTierLabel(mt.tier)}
                          </span>
                        </span>
                        <span className="ds-text-muted de-catalog-stat">{getModuleStat(mt)}</span>
                      </div>
                      <div className="de-catalog-right">
                        <span className="ds-mono ds-text-muted">v{mt.volume}</span>
                        {installed && (
                          <span className="ds-badge ds-badge--teal ds-mono">×{installed.quantity}</span>
                        )}
                        {!hasBp && (
                          <span className="ds-badge ds-badge--neutral">
                            <Lock size={10} strokeWidth={2} /> No BP
                          </span>
                        )}
                        {hasBp && !tierUnlocked && (
                          <span className="ds-badge ds-badge--warning">
                            <Lock size={10} strokeWidth={2} /> {getTierLabel(mt.tier)}
                          </span>
                        )}
                      </div>
                    </button>
                  )
                })
              )}
            </div>
          </section>
        </div>

        {hull && (
          <footer className="de-footer">
            <div className="de-stats">
              <Stat label="Attack" value={totalStats.attack} accent="orange" />
              <Stat label="Shield" value={totalStats.shield} accent="info" />
              <Stat label="Atk/Rnd" value={totalStats.attack} accent="orange" />
              <Stat label="Structure" value={totalStats.structure} />
              <Stat label="Agility" value={totalStats.agility} />
              <Stat label="Storage" value={totalStats.storage} />
              <Stat label="Stability" value={totalStats.stability} />
              <Stat label="Mobility" value={totalStats.mobility} />
              <Stat label="Defense" value={totalStats.defense} />
              <Stat label="He3/Rnd" value={totalStats.he3} accent="info" />
            </div>
            <div className="de-costs">
              <span className="de-cost">
                <span className="ds-resource-dot ds-resource-dot--metal" />
                Metal: <span className="ds-mono">{totalStats.metalCost.toLocaleString()}</span>
              </span>
              <span className="de-cost">
                <span className="ds-resource-dot ds-resource-dot--he3" />
                He3: <span className="ds-mono">{totalStats.he3Cost.toLocaleString()}</span>
              </span>
              <span className="de-cost">
                <span className="ds-resource-dot ds-resource-dot--gold" />
                Gold: <span className="ds-mono">{totalStats.goldCost.toLocaleString()}</span>
              </span>
            </div>
          </footer>
        )}

        <style>{`
          .de-backdrop { padding: 0; }
          .de-shell {
            width: 100vw; height: 100vh; max-width: 100vw; max-height: 100vh;
            display: flex; flex-direction: column;
            border: 0;
          }
          .de-header {
            display: flex; align-items: center; justify-content: space-between;
            padding: var(--sp-3) var(--sp-5);
            border-bottom: 1px solid var(--ds-border);
            flex-shrink: 0;
          }
          .de-title { display: flex; align-items: center; gap: var(--sp-2); color: var(--ds-text); }
          .de-title h2 { margin: 0; }
          .de-body {
            flex: 1; min-height: 0;
            display: grid;
            grid-template-columns: 280px minmax(0, 1fr) 320px;
            gap: var(--sp-4);
            padding: var(--sp-4) var(--sp-5);
            overflow: hidden;
          }
          .de-col {
            display: flex; flex-direction: column; gap: var(--sp-3);
            min-height: 0; overflow: hidden;
          }
          .de-col--center { gap: var(--sp-3); }
          .de-tabs { flex-shrink: 0; }
          .de-hull-list, .de-catalog {
            flex: 1; min-height: 0; overflow-y: auto;
            display: flex; flex-direction: column; gap: var(--sp-2);
          }
          .de-hull-card {
            display: flex; flex-direction: column; gap: 6px;
            text-align: left; cursor: pointer;
          }
          .de-hull-card[disabled] { cursor: not-allowed; }
          .de-hull-top {
            display: flex; align-items: center; justify-content: space-between;
            gap: var(--sp-2);
          }
          .de-hull-name { font-weight: var(--fw-semibold); color: var(--ds-text); }
          .de-hull-stats {
            display: flex; gap: var(--sp-3);
            font-size: var(--fs-caption); color: var(--ds-text-muted);
          }
          .de-lock { display: inline-flex; align-items: center; gap: 4px; align-self: flex-start; }
          .de-empty { padding: var(--sp-4); text-align: center; }

          .de-preview {
            height: 220px;
            background: var(--ds-surface-3);
            border-radius: var(--r-lg);
            overflow: hidden;
            flex-shrink: 0;
          }
          .de-form-group { display: flex; flex-direction: column; gap: 6px; }

          .de-installed {
            display: flex; flex-direction: column; gap: var(--sp-2);
            flex: 1; min-height: 0; overflow: hidden;
          }
          .de-installed-head { display: flex; align-items: center; justify-content: space-between; }
          .de-module-list {
            flex: 1; min-height: 0; overflow-y: auto;
            display: flex; flex-direction: column; gap: var(--sp-2);
          }
          .de-mi {
            display: grid;
            grid-template-columns: auto 1fr auto auto auto auto;
            align-items: center; gap: var(--sp-2);
            padding: var(--sp-2) var(--sp-3);
          }
          .de-mi-name {
            font-weight: var(--fw-medium); color: var(--ds-text);
            white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
          }

          .de-volume { display: flex; flex-direction: column; gap: 6px; }
          .de-volume-row { display: flex; align-items: center; justify-content: space-between; font-size: var(--fs-sm); }
          .de-volume-over { color: var(--ds-danger); font-weight: var(--fw-semibold); }

          .de-actions {
            display: flex; justify-content: flex-end; gap: var(--sp-2);
          }

          .de-placeholder {
            flex: 1; display: grid; place-items: center;
            font-size: var(--fs-h3);
          }

          .de-sub-tabs {
            display: flex; flex-wrap: wrap; gap: var(--sp-1);
            flex-shrink: 0;
          }
          .de-sub-tab {
            cursor: pointer; border: 0;
          }

          .de-catalog-mod {
            display: grid;
            grid-template-columns: auto 1fr auto;
            align-items: center; gap: var(--sp-2);
            text-align: left; cursor: pointer;
          }
          .de-catalog-mod[disabled] { cursor: not-allowed; }
          .de-catalog-info { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
          .de-catalog-name {
            display: inline-flex; align-items: center; gap: var(--sp-1);
            font-weight: var(--fw-medium); color: var(--ds-text);
            font-size: var(--fs-sm);
          }
          .de-tier { font-size: 10px; padding: 1px 5px; }
          .de-catalog-stat { font-size: var(--fs-caption); }
          .de-catalog-right { display: inline-flex; align-items: center; gap: var(--sp-1); }

          .de-footer {
            border-top: 1px solid var(--ds-border);
            padding: var(--sp-3) var(--sp-5);
            display: flex; flex-direction: column; gap: var(--sp-2);
            background: var(--ds-surface-2);
            flex-shrink: 0;
          }
          .de-stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(96px, 1fr));
            gap: var(--sp-3);
          }
          .de-costs {
            display: flex; flex-wrap: wrap; gap: var(--sp-4);
            font-size: var(--fs-sm); color: var(--ds-text);
          }
          .de-cost { display: inline-flex; align-items: center; gap: var(--sp-1); }
        `}</style>
      </div>
    </div>,
    document.body,
  )
}

function Stat({
  label,
  value,
  accent,
}: {
  label: string
  value: number
  accent?: 'orange' | 'info'
}) {
  const colorVar =
    accent === 'orange'
      ? 'var(--ds-orange-strong)'
      : accent === 'info'
        ? 'var(--ds-info)'
        : 'var(--ds-text)'
  return (
    <div className="de-stat">
      <span className="ds-caption">{label}</span>
      <span className="ds-mono de-stat-value" style={{ color: colorVar }}>
        {value}
      </span>
      <style>{`
        .de-stat { display: flex; flex-direction: column; gap: 2px; align-items: flex-start; }
        .de-stat-value { font-size: var(--fs-h3); font-weight: var(--fw-semibold); }
      `}</style>
    </div>
  )
}

// Main ShipDesignPanel -----------------------------------------------------

export default function ShipDesignPanel() {
  const { designs, hullTypes, moduleTypes, loading, error, create, remove } = useShipDesigns()
  const { hasHullBlueprint, hasModuleBlueprint, myBlueprints } = useBlueprints()
  const [showEditor, setShowEditor] = useState(false)

  if (loading) {
    return (
      <div className="ds-panel sd-loading">
        <div className="loading-spinner" />
        <span>Loading Designs...</span>
      </div>
    )
  }

  function getHullBlueprintResearchLevel(hullTypeId: number): number {
    const bp = myBlueprints.find(
      bp => bp.blueprint_type === 'hull' && bp.hull_type_id === hullTypeId && bp.is_activated,
    )
    return bp?.research_level ?? 0
  }

  function getModuleBlueprintResearchLevel(moduleTypeId: number): number {
    const bp = myBlueprints.find(
      bp =>
        bp.blueprint_type === 'module' && bp.module_type_id === moduleTypeId && bp.is_activated,
    )
    return bp?.research_level ?? 0
  }

  async function handleSave(name: string, hullTypeId: number, modules: ShipDesignModule[]) {
    await create({ name, hull_type_id: hullTypeId, modules })
  }

  return (
    <div className="ds-panel sd-panel">
      <header className="sd-header">
        <div className="sd-header-icon">
          <Cog size={22} strokeWidth={1.75} />
        </div>
        <div className="sd-header-text">
          <h2 className="ds-h2">Ship Designs</h2>
          <p className="ds-text-muted sd-subtitle">
            <span className="ds-mono">{designs.length}</span>/20 designs
          </p>
        </div>
        <button
          className="ds-btn ds-btn-secondary sd-new"
          onClick={() => setShowEditor(true)}
          disabled={designs.length >= 20}
        >
          New Design
        </button>
      </header>

      {error && (
        <div className="ds-badge ds-badge--danger sd-error" role="alert">
          {error}
        </div>
      )}

      {designs.length === 0 ? (
        <div className="ds-card sd-empty">
          <Wrench size={32} strokeWidth={1.5} />
          <p>No ship designs yet. Create your first design to start building ships.</p>
        </div>
      ) : (
        <div className="sd-list">
          {designs.map(d => (
            <DesignCard
              key={d.id}
              design={d}
              hullTypes={hullTypes}
              onDelete={() => remove(d.id)}
            />
          ))}
        </div>
      )}

      {showEditor && (
        <DesignEditor
          hullTypes={hullTypes}
          moduleTypes={moduleTypes}
          hasHullBlueprint={hasHullBlueprint}
          hasModuleBlueprint={hasModuleBlueprint}
          getHullBlueprintResearchLevel={getHullBlueprintResearchLevel}
          getModuleBlueprintResearchLevel={getModuleBlueprintResearchLevel}
          onSave={handleSave}
          onClose={() => setShowEditor(false)}
        />
      )}

      <style>{`
        .sd-panel { display: flex; flex-direction: column; gap: var(--sp-4); max-width: 1100px; }
        .sd-loading { display: flex; align-items: center; gap: var(--sp-3); }
        .sd-header { display: flex; align-items: center; gap: var(--sp-3); }
        .sd-header-icon {
          width: 44px; height: 44px;
          display: grid; place-items: center;
          background: var(--ds-teal-tint);
          color: var(--ds-teal-dark);
          border-radius: var(--r-lg);
        }
        .sd-header-text { flex: 1; }
        .sd-header-text h2 { margin: 0; }
        .sd-subtitle { margin: 2px 0 0; font-size: var(--fs-sm); }
        .sd-error { display: inline-flex; }
        .sd-empty {
          display: flex; flex-direction: column; align-items: center;
          gap: var(--sp-2); padding: var(--sp-6);
          color: var(--ds-text-muted);
        }
        .sd-empty p { margin: 0; }
        .sd-list {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
          gap: var(--sp-3);
        }
      `}</style>
    </div>
  )
}

function DesignCard({
  design,
  hullTypes,
  onDelete,
}: {
  design: ShipDesign
  hullTypes: HullType[]
  onDelete: () => void
}) {
  const hull = hullTypes.find(h => h.id === design.hull_type_id)
  const hullClass = design.hull_class || hull?.hull_class || 'frigate'
  const tint = HULL_CLASS_TINT[hullClass] || 'var(--ds-text-muted)'
  // ships_built is provided by the API but isn't on the type yet
  const built = (design as ShipDesign & { ships_built?: number }).ships_built ?? 0

  return (
    <div className="ds-card sd-card">
      <div className="sd-card-head">
        <div className="sd-card-class" style={{ color: tint, borderColor: tint }}>
          {hullClass.charAt(0).toUpperCase()}
        </div>
        <div className="sd-card-info">
          <div className="sd-card-name">{design.name}</div>
          <div className="ds-text-muted sd-card-hull">
            {design.hull_name || hull?.display_name || 'Unknown Hull'}
          </div>
        </div>
        {built === 0 && (
          <button
            className="ds-btn-icon ds-btn-icon--sm"
            onClick={onDelete}
            aria-label="Delete design"
          >
            <Trash2 size={14} strokeWidth={2} />
          </button>
        )}
      </div>
      <div className="sd-card-stats">
        <Stat2 label="SH" value={design.total_shield} />
        <Stat2 label="ST" value={design.total_structure} />
        <Stat2 label="ATK" value={design.attack_power} />
        <Stat2 label="VOL" value={design.volume_used} />
        <Stat2 label="Built" value={built} />
      </div>

      <style>{`
        .sd-card { display: flex; flex-direction: column; gap: var(--sp-2); }
        .sd-card-head { display: flex; align-items: center; gap: var(--sp-3); }
        .sd-card-class {
          width: 36px; height: 36px;
          display: grid; place-items: center;
          border-radius: var(--r-md);
          border: 2px solid currentColor;
          font-weight: var(--fw-bold); flex-shrink: 0;
        }
        .sd-card-info { flex: 1; min-width: 0; }
        .sd-card-name {
          font-weight: var(--fw-semibold); color: var(--ds-text);
          white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
        }
        .sd-card-hull { font-size: var(--fs-sm); }
        .sd-card-stats {
          display: flex; flex-wrap: wrap; gap: var(--sp-3);
          font-size: var(--fs-sm);
        }
      `}</style>
    </div>
  )
}

function Stat2({ label, value }: { label: string; value: number }) {
  return (
    <span className="sd-stat2">
      <span className="ds-caption">{label}</span>
      <span className="ds-mono">{value}</span>
      <style>{`
        .sd-stat2 { display: inline-flex; align-items: center; gap: 4px; color: var(--ds-text); }
      `}</style>
    </span>
  )
}
