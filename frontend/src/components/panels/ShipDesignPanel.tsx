import { useState, useMemo, Suspense } from 'react'
import { createPortal } from 'react-dom'
import { Canvas } from '@react-three/fiber'
import { OrbitControls, useGLTF, Environment } from '@react-three/drei'
import { useShipDesigns } from '../../hooks/useShipDesigns.ts'
import { useBlueprints } from '../../hooks/useBlueprints.ts'
import LoadingButton from '../common/LoadingButton.tsx'
import type { HullType, ModuleType, ShipDesignModule, ShipDesign } from '../../types'
import '../../styles/common.css'

const HULL_CLASS_COLORS: Record<string, string> = {
  frigate: '#44aaff',
  cruiser: '#cc8844',
  battleship: '#aa4444',
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

// --- Tier helpers ---

const TIER_LABELS: Record<number, string> = {
  1: 'I',
  2: 'II',
  3: 'III',
}

function getTierLabel(tier: number): string {
  return TIER_LABELS[tier] || `T${tier}`
}

// --- Stat helpers (matches backend ship_formulas.go) ---

function getModuleAvgAttack(mt: ModuleType): number {
  return Math.round((mt.min_damage + mt.max_damage) / 2)
}

function parseEffects(mt: ModuleType): Record<string, number> {
  try {
    return JSON.parse(mt.effects_json || '{}')
  } catch { return {} }
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

// --- 3D Ship Preview ---

function ShipModel({ url }: { url: string }) {
  const { scene } = useGLTF(url)
  return <primitive object={scene.clone()} scale={0.1} position={[0, 0, 0]} />
}

function ShipPreview({ hullClass }: { hullClass: string }) {
  const modelUrl = HULL_CLASS_MODELS[hullClass] || HULL_CLASS_MODELS.frigate
  return (
    <div className={`de-preview-3d ${hullClass}`}>
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

// --- Design Editor Props ---

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

function DesignEditor({ hullTypes, moduleTypes, hasHullBlueprint, hasModuleBlueprint, getHullBlueprintResearchLevel, getModuleBlueprintResearchLevel, onSave, onClose }: DesignEditorProps) {
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
    // Player must have the blueprint activated
    if (!hasHullBlueprint(h.id)) return false
    const researchLevel = getHullBlueprintResearchLevel(h.id)
    // Show hull if research_level >= tier
    return researchLevel >= h.tier
  })

  const groupCategories = MODULE_GROUPS[moduleGroup]?.categories || []
  const filteredModules = useMemo(() => {
    return moduleTypes.filter(mt => {
      if (!groupCategories.includes(mt.category)) return false
      if (moduleSubCategory && mt.category !== moduleSubCategory) return false
      // Player must have the blueprint activated
      if (!hasModuleBlueprint(mt.id)) return false
      const researchLevel = getModuleBlueprintResearchLevel(mt.id)
      // Show module if research_level >= tier
      return researchLevel >= mt.tier
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
    selectedHull !== null && modules.length > 0 && hull && !isOverVolume

  return createPortal(
    <div className="de-fullscreen-backdrop" onClick={e => { if (e.target === e.currentTarget) onClose() }}>
      <div className="de-fullscreen">
        {/* Header */}
        <div className="de-header">
          <span className="de-header-title">Ship Design Editor</span>
          <button className="p2-modal-close" onClick={onClose}>X</button>
        </div>

        <div className="de-body">
          {/* Column 1: Hull Selection */}
          <div className="de-col de-col-hull">
            <div className="de-col-title">Select Hull</div>
            <div className="de-class-tabs">
              {(['frigate', 'cruiser', 'battleship'] as const).map(cls => (
                <button
                  key={cls}
                  className={`de-class-tab ${hullClassFilter === cls ? 'active' : ''}`}
                  style={{
                    borderColor: hullClassFilter === cls ? HULL_CLASS_COLORS[cls] : undefined,
                    color: hullClassFilter === cls ? HULL_CLASS_COLORS[cls] : undefined,
                  }}
                  onClick={() => { setHullClassFilter(cls); setSelectedHull(null); setModules([]) }}
                >
                  <span className="de-class-icon" style={{ background: HULL_CLASS_COLORS[cls] }}>
                    {cls.charAt(0).toUpperCase()}
                  </span>
                  <span className="de-class-label">{cls.charAt(0).toUpperCase() + cls.slice(1)}</span>
                </button>
              ))}
            </div>

            <div className="de-hull-list">
              {filteredHulls.map(h => {
                const hasBp = hasHullBlueprint(h.id)
                const researchLevel = getHullBlueprintResearchLevel(h.id)
                const tierUnlocked = researchLevel >= h.tier
                return (
                  <button
                    key={h.id}
                    className={`de-hull-card ${selectedHull === h.id ? 'selected' : ''} ${!tierUnlocked ? 'locked' : ''}`}
                    onClick={() => { if (tierUnlocked) { setSelectedHull(h.id); setModules([]) } }}
                    disabled={!tierUnlocked}
                    title={!tierUnlocked ? `Tier ${getTierLabel(h.tier)} locked - Research blueprint to level ${h.tier}` : ''}
                  >
                    <div className="de-hull-top">
                      <span className="de-hull-name">{h.display_name}</span>
                      <span className={`de-hull-tier tier-${h.tier}`}>{getTierLabel(h.tier)}</span>
                    </div>
                    <div className="de-hull-stats">
                      <span>SH:{h.base_shield}</span>
                      <span>ST:{h.base_structure}</span>
                      <span>Slots:{h.installation_slots}</span>
                    </div>
                    {!hasBp && <div className="de-hull-lock">🔒 No Blueprint</div>}
                    {hasBp && !tierUnlocked && <div className="de-hull-lock">🔒 Tier {getTierLabel(h.tier)} Locked</div>}
                  </button>
                )
              })}
              {filteredHulls.length === 0 && (
                <div className="de-empty">No hulls of this class unlocked</div>
              )}
            </div>
          </div>

          {/* Column 2: Ship Preview + Design */}
          <div className="de-col de-col-center">
            {hull ? (
              <>
                <ShipPreview hullClass={hull.hull_class} />

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

                {/* Installed modules */}
                <div className="de-installed">
                  <div className="de-installed-header">
                    <span className="p2-section-title">Installed Modules</span>
                    <span className="de-installed-count">{modules.reduce((s, m) => s + m.quantity, 0)} items</span>
                  </div>
                  {modules.length === 0 ? (
                    <div className="de-empty">Select modules from the right panel</div>
                  ) : (
                    <div className="de-module-list">
                      {modules.map(m => {
                        const mt = moduleTypes.find(mt => mt.id === m.module_type_id)
                        if (!mt) return null
                        return (
                          <div key={m.module_type_id} className="de-module-installed">
                            <span className={`de-mod-cat ${mt.category}`}>{CATEGORY_LABELS[mt.category] || mt.category.slice(0, 3).toUpperCase()}</span>
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

                {/* Volume bar */}
                <div className="de-volume-section">
                  <div className="de-volume-label">
                    <span>Volume</span>
                    <span className={isOverVolume ? 'de-over' : ''}>{volumeUsed} / {volumeMax}</span>
                  </div>
                  <div className="de-volume-bar">
                    <div
                      className={`de-volume-fill ${isOverVolume ? 'over' : isNearCapacity ? 'warning' : ''}`}
                      style={{ width: `${Math.min(100, volumePct)}%` }}
                    />
                  </div>
                </div>

                {/* Action buttons */}
                <div className="de-actions">
                  <button className="p2-btn p2-btn-secondary" onClick={onClose}>Cancel</button>
                  <LoadingButton
                    className="p2-btn p2-btn-primary"
                    disabled={!isValid}
                    loading={saving}
                    onClick={handleSave}
                  >
                    Save Design
                  </LoadingButton>
                </div>
              </>
            ) : (
              <div className="de-placeholder">
                <div className="de-placeholder-text">Select a hull to begin designing</div>
              </div>
            )}
          </div>

          {/* Column 3: Module Selection */}
          <div className="de-col de-col-modules">
            <div className="de-col-title">Modules</div>

            {/* Module group tabs */}
            <div className="de-group-tabs">
              {Object.entries(MODULE_GROUPS).map(([key, group]) => (
                <button
                  key={key}
                  className={`de-group-tab ${moduleGroup === key ? 'active' : ''}`}
                  onClick={() => { setModuleGroup(key); setModuleSubCategory(null) }}
                >
                  {group.label}
                </button>
              ))}
            </div>

            {/* Sub-category tabs */}
            <div className="de-sub-tabs">
              <button
                className={`de-sub-tab ${moduleSubCategory === null ? 'active' : ''}`}
                onClick={() => setModuleSubCategory(null)}
              >
                All
              </button>
              {groupCategories.map(cat => (
                <button
                  key={cat}
                  className={`de-sub-tab ${moduleSubCategory === cat ? 'active' : ''}`}
                  onClick={() => setModuleSubCategory(cat)}
                >
                  <span className={`de-mod-cat ${cat}`}>{CATEGORY_LABELS[cat] || cat.slice(0, 3).toUpperCase()}</span>
                </button>
              ))}
            </div>

            {/* Module list */}
            <div className="de-module-catalog">
              {filteredModules.length === 0 ? (
                <div className="de-empty">No modules unlocked in this category</div>
              ) : (
                filteredModules.map(mt => {
                  const hasBp = hasModuleBlueprint(mt.id)
                  const researchLevel = getModuleBlueprintResearchLevel(mt.id)
                  const tierUnlocked = researchLevel >= mt.tier
                  const installed = modules.find(m => m.module_type_id === mt.id)
                  const atMax = mt.max_per_ship > 0 && (installed?.quantity || 0) >= mt.max_per_ship
                  return (
                    <button
                      key={mt.id}
                      className={`de-catalog-module ${!tierUnlocked ? 'locked' : ''} ${installed ? 'installed' : ''} ${atMax ? 'at-max' : ''}`}
                      onClick={() => { if (tierUnlocked && hull && !atMax) addModule(mt.id) }}
                      disabled={!tierUnlocked || !hull || atMax}
                      title={!tierUnlocked ? `Tier ${getTierLabel(mt.tier)} locked - Research blueprint to level ${mt.tier}` : ''}
                    >
                      <span className={`de-mod-cat ${mt.category}`}>{CATEGORY_LABELS[mt.category] || mt.category.slice(0, 3).toUpperCase()}</span>
                      <div className="de-catalog-info">
                        <span className="de-catalog-name">
                          {mt.display_name}
                          <span className={`de-mod-tier tier-${mt.tier}`}>{getTierLabel(mt.tier)}</span>
                        </span>
                        <span className="de-catalog-stat">{getModuleStat(mt)}</span>
                      </div>
                      <div className="de-catalog-right">
                        <span className="de-catalog-vol">v{mt.volume}</span>
                        {installed && <span className="de-catalog-qty">x{installed.quantity}</span>}
                        {!hasBp && <span className="de-catalog-lock">🔒 No BP</span>}
                        {hasBp && !tierUnlocked && <span className="de-catalog-lock">🔒 Tier {getTierLabel(mt.tier)}</span>}
                      </div>
                    </button>
                  )
                })
              )}
            </div>
          </div>
        </div>

        {/* Bottom stats bar */}
        {hull && (
          <div className="de-bottom-bar">
            <div className="de-stats-row">
              <div className="de-stat-cell">
                <span className="de-stat-label">Attack</span>
                <span className="de-stat-value atk">{totalStats.attack}</span>
              </div>
              <div className="de-stat-cell">
                <span className="de-stat-label">Shield</span>
                <span className="de-stat-value shd">{totalStats.shield}</span>
              </div>
              <div className="de-stat-cell">
                <span className="de-stat-label">Atk/Rnd</span>
                <span className="de-stat-value atk">{totalStats.attack}</span>
              </div>
              <div className="de-stat-cell">
                <span className="de-stat-label">Structure</span>
                <span className="de-stat-value str">{totalStats.structure}</span>
              </div>
              <div className="de-stat-cell">
                <span className="de-stat-label">Agility</span>
                <span className="de-stat-value">{totalStats.agility}</span>
              </div>
              <div className="de-stat-cell">
                <span className="de-stat-label">Storage</span>
                <span className="de-stat-value">{totalStats.storage}</span>
              </div>
              <div className="de-stat-cell">
                <span className="de-stat-label">Stability</span>
                <span className="de-stat-value">{totalStats.stability}</span>
              </div>
              <div className="de-stat-cell">
                <span className="de-stat-label">Mobility</span>
                <span className="de-stat-value">{totalStats.mobility}</span>
              </div>
              <div className="de-stat-cell">
                <span className="de-stat-label">Defense</span>
                <span className="de-stat-value">{totalStats.defense}</span>
              </div>
              <div className="de-stat-cell">
                <span className="de-stat-label">He3/Rnd</span>
                <span className="de-stat-value he3">{totalStats.he3}</span>
              </div>
            </div>
            <div className="de-cost-row">
              <span className="de-cost"><span className="de-cost-dot metal" />Metal: {totalStats.metalCost.toLocaleString()}</span>
              <span className="de-cost"><span className="de-cost-dot he3" />He3: {totalStats.he3Cost.toLocaleString()}</span>
              <span className="de-cost"><span className="de-cost-dot gold" />Gold: {totalStats.goldCost.toLocaleString()}</span>
            </div>
          </div>
        )}
      </div>
    </div>,
    document.body,
  )
}

// --- Main ShipDesignPanel ---

export default function ShipDesignPanel() {
  const { designs, hullTypes, moduleTypes, loading, error, create, remove } = useShipDesigns()
  const { hasHullBlueprint, hasModuleBlueprint, myBlueprints } = useBlueprints()
  const [showEditor, setShowEditor] = useState(false)

  if (loading) {
    return <div className="p2-panel-loading"><div className="loading-spinner" /><span>Loading Designs...</span></div>
  }

  // Helper to get blueprint research level for a hull type
  function getHullBlueprintResearchLevel(hullTypeId: number): number {
    const bp = myBlueprints.find(bp => bp.blueprint_type === 'hull' && bp.hull_type_id === hullTypeId && bp.is_activated)
    return bp?.research_level ?? 0
  }

  // Helper to get blueprint research level for a module type
  function getModuleBlueprintResearchLevel(moduleTypeId: number): number {
    const bp = myBlueprints.find(bp => bp.blueprint_type === 'module' && bp.module_type_id === moduleTypeId && bp.is_activated)
    return bp?.research_level ?? 0
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
          hasHullBlueprint={hasHullBlueprint}
          hasModuleBlueprint={hasModuleBlueprint}
          getHullBlueprintResearchLevel={getHullBlueprintResearchLevel}
          getModuleBlueprintResearchLevel={getModuleBlueprintResearchLevel}
          onSave={handleSave}
          onClose={() => setShowEditor(false)}
        />
      )}
    </div>
  )
}

// --- Design Card (list view) ---

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
