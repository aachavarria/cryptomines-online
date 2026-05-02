import { createContext, useContext, useReducer, useCallback, type ReactNode } from 'react'
import type { Player, Planet, BuildingWithType, ResourcesResponse, BuildingTypeData } from '../types'
import { getBuildingSize } from '../config/buildingConfig'

// Square grid configuration (camera rotation creates isometric look)
export const GRID_SIZE = 20 // 20x20 square grid
export const TILE_WORLD_SIZE = 4 // 4 world units per tile side
const GRID_CENTER = GRID_SIZE / 2 // center offset so (10,10) → world (0,0,0)

export interface GridPosition {
  col: number
  row: number
}

// Convert grid coords to 3D world position (square grid, centered at origin)
export function gridToWorld(col: number, row: number): [number, number, number] {
  const x = (col - GRID_CENTER) * TILE_WORLD_SIZE
  const z = (row - GRID_CENTER) * TILE_WORLD_SIZE
  return [x, 0, z]
}

// Convert 3D world position to nearest grid coords
export function worldToGrid(x: number, z: number): GridPosition {
  const col = Math.round(x / TILE_WORLD_SIZE + GRID_CENTER)
  const row = Math.round(z / TILE_WORLD_SIZE + GRID_CENTER)
  return { col, row }
}

// Get all tiles a building occupies given its anchor (top-left) position
export function getBuildingTiles(typeName: string, anchorCol: number, anchorRow: number): GridPosition[] {
  const size = getBuildingSize(typeName)
  const tiles: GridPosition[] = []
  for (let dc = 0; dc < size.cols; dc++) {
    for (let dr = 0; dr < size.rows; dr++) {
      tiles.push({ col: anchorCol + dc, row: anchorRow + dr })
    }
  }
  return tiles
}

// Check if a building can be placed at anchor position
export function canPlaceBuilding(
  typeName: string,
  anchorCol: number,
  anchorRow: number,
  occupiedTiles: Set<string>,
  excludeBuildingId?: string,
  buildingPositions?: Record<string, GridPosition>,
  buildings?: { id: string; type_name: string }[],
): boolean {
  const tiles = getBuildingTiles(typeName, anchorCol, anchorRow)

  // Build a set of tiles to exclude (for move mode — exclude the building being moved)
  const excludeTiles = new Set<string>()
  if (excludeBuildingId && buildingPositions && buildings) {
    const excludeBuilding = buildings.find(b => b.id === excludeBuildingId)
    if (excludeBuilding) {
      const pos = buildingPositions[excludeBuildingId]
      if (pos) {
        const existingTiles = getBuildingTiles(excludeBuilding.type_name, pos.col, pos.row)
        for (const t of existingTiles) {
          excludeTiles.add(`${t.col},${t.row}`)
        }
      }
    }
  }

  for (const tile of tiles) {
    // Check bounds
    if (tile.col < 0 || tile.col >= GRID_SIZE || tile.row < 0 || tile.row >= GRID_SIZE) {
      return false
    }
    // Check occupied (excluding tiles of the building being moved)
    const key = `${tile.col},${tile.row}`
    if (occupiedTiles.has(key) && !excludeTiles.has(key)) {
      return false
    }
  }
  return true
}

// Get the world-space center of a multi-tile building
export function getBuildingWorldCenter(typeName: string, anchorCol: number, anchorRow: number): [number, number, number] {
  const size = getBuildingSize(typeName)
  const centerCol = anchorCol + (size.cols - 1) / 2
  const centerRow = anchorRow + (size.rows - 1) / 2
  return gridToWorld(centerCol, centerRow)
}

export interface PlacementMode {
  active: boolean
  buildingTypeName?: string   // for new construction
  movingBuildingId?: string   // for moving existing
}

export type BaseView = 'ground' | 'space'

export interface GameState {
  player: Player | null
  planets: Planet[]
  currentPlanet: Planet | null
  buildings: BuildingWithType[]
  buildingTypes: BuildingTypeData[]
  resources: ResourcesResponse | null
  selectedBuilding: BuildingWithType | null
  hoveredBuilding: string | null
  showContextMenu: boolean
  contextMenuScreenPos: { x: number; y: number } | null
  showDetailPanel: boolean
  showConstructPanel: boolean
  placementMode: PlacementMode
  // Map building id -> grid position (persisted client-side for now)
  buildingPositions: Record<string, GridPosition>
  cameraTarget: [number, number, number] | null
  // Which base the player is currently viewing/placing on
  currentBase: BaseView
  loading: boolean
  error: string | null
}

type GameAction =
  | { type: 'SET_PLAYER'; payload: Player }
  | { type: 'SET_PLANETS'; payload: Planet[] }
  | { type: 'SET_CURRENT_PLANET'; payload: Planet }
  | { type: 'SET_BUILDINGS'; payload: BuildingWithType[] }
  | { type: 'SET_BUILDING_TYPES'; payload: BuildingTypeData[] }
  | { type: 'SET_RESOURCES'; payload: ResourcesResponse }
  | { type: 'SELECT_BUILDING'; payload: { building: BuildingWithType | null; screenPos?: { x: number; y: number } } }
  | { type: 'HOVER_BUILDING'; payload: string | null }
  | { type: 'SHOW_CONTEXT_MENU'; payload: boolean }
  | { type: 'OPEN_DETAIL_PANEL' }
  | { type: 'CLOSE_DETAIL_PANEL' }
  | { type: 'SHOW_CONSTRUCT_PANEL'; payload: boolean }
  | { type: 'SET_PLACEMENT_MODE'; payload: PlacementMode }
  | { type: 'PLACE_BUILDING'; payload: { buildingId: string; position: GridPosition } }
  | { type: 'SET_BUILDING_POSITIONS'; payload: Record<string, GridPosition> }
  | { type: 'SET_CAMERA_TARGET'; payload: [number, number, number] | null }
  | { type: 'SET_BASE_VIEW'; payload: BaseView }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'DESELECT_ALL' }

const initialState: GameState = {
  player: null,
  planets: [],
  currentPlanet: null,
  buildings: [],
  buildingTypes: [],
  resources: null,
  selectedBuilding: null,
  hoveredBuilding: null,
  showContextMenu: false,
  contextMenuScreenPos: null,
  showDetailPanel: false,
  showConstructPanel: false,
  placementMode: { active: false },
  buildingPositions: {},
  cameraTarget: null,
  currentBase: 'ground',
  loading: true,
  error: null,
}

function gameReducer(state: GameState, action: GameAction): GameState {
  switch (action.type) {
    case 'SET_PLAYER':
      return { ...state, player: action.payload }
    case 'SET_PLANETS':
      return { ...state, planets: action.payload }
    case 'SET_CURRENT_PLANET':
      return { ...state, currentPlanet: action.payload }
    case 'SET_BUILDINGS':
      return { ...state, buildings: action.payload }
    case 'SET_BUILDING_TYPES':
      return { ...state, buildingTypes: action.payload }
    case 'SET_RESOURCES':
      return { ...state, resources: action.payload }
    case 'SELECT_BUILDING':
      return {
        ...state,
        selectedBuilding: action.payload.building,
        showContextMenu: action.payload.building !== null,
        contextMenuScreenPos: action.payload.screenPos ?? null,
        showDetailPanel: false,
      }
    case 'HOVER_BUILDING':
      return { ...state, hoveredBuilding: action.payload }
    case 'SHOW_CONTEXT_MENU':
      return { ...state, showContextMenu: action.payload }
    case 'OPEN_DETAIL_PANEL':
      return { ...state, showDetailPanel: true, showContextMenu: false }
    case 'CLOSE_DETAIL_PANEL':
      return { ...state, showDetailPanel: false, selectedBuilding: null }
    case 'SHOW_CONSTRUCT_PANEL':
      return {
        ...state,
        showConstructPanel: action.payload,
        showContextMenu: false,
        selectedBuilding: action.payload ? null : state.selectedBuilding,
      }
    case 'SET_PLACEMENT_MODE':
      return {
        ...state,
        placementMode: action.payload,
        showContextMenu: false,
        showConstructPanel: false,
      }
    case 'PLACE_BUILDING': {
      const newPositions = { ...state.buildingPositions }
      newPositions[action.payload.buildingId] = action.payload.position
      return {
        ...state,
        buildingPositions: newPositions,
        placementMode: { active: false },
      }
    }
    case 'SET_BUILDING_POSITIONS':
      return { ...state, buildingPositions: action.payload }
    case 'SET_CAMERA_TARGET':
      return { ...state, cameraTarget: action.payload }
    case 'SET_BASE_VIEW':
      return {
        ...state,
        currentBase: action.payload,
        selectedBuilding: null,
        showContextMenu: false,
        showDetailPanel: false,
      }
    case 'SET_LOADING':
      return { ...state, loading: action.payload }
    case 'SET_ERROR':
      return { ...state, error: action.payload }
    case 'DESELECT_ALL':
      return {
        ...state,
        selectedBuilding: null,
        showContextMenu: false,
        showDetailPanel: false,
        placementMode: state.placementMode.active ? { active: false } : state.placementMode,
      }
    default:
      return state
  }
}

interface GameContextValue {
  state: GameState
  dispatch: React.Dispatch<GameAction>
  selectBuilding: (building: BuildingWithType | null, screenPos?: { x: number; y: number }) => void
  deselectAll: () => void
  openDetailPanel: () => void
  closeDetailPanel: () => void
  openConstructPanel: () => void
  closeConstructPanel: () => void
  enterPlacementMode: (buildingTypeName: string) => void
  enterMoveMode: (buildingId: string) => void
  exitPlacementMode: () => void
  placeBuilding: (buildingId: string, position: GridPosition) => void
  focusCamera: (position: [number, number, number] | null) => void
  setBaseView: (base: BaseView) => void
}

const GameContext = createContext<GameContextValue | null>(null)

export function GameProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(gameReducer, initialState)

  const selectBuilding = useCallback((building: BuildingWithType | null, screenPos?: { x: number; y: number }) => {
    dispatch({ type: 'SELECT_BUILDING', payload: { building, screenPos } })
    if (building) {
      const pos = state.buildingPositions[building.id]
      if (pos) {
        const worldPos = gridToWorld(pos.col, pos.row)
        dispatch({ type: 'SET_CAMERA_TARGET', payload: worldPos })
      }
    }
  }, [state.buildingPositions])

  const deselectAll = useCallback(() => {
    dispatch({ type: 'DESELECT_ALL' })
  }, [])

  const openDetailPanel = useCallback(() => {
    dispatch({ type: 'OPEN_DETAIL_PANEL' })
  }, [])

  const closeDetailPanel = useCallback(() => {
    dispatch({ type: 'CLOSE_DETAIL_PANEL' })
  }, [])

  const openConstructPanel = useCallback(() => {
    dispatch({ type: 'SHOW_CONSTRUCT_PANEL', payload: true })
  }, [])

  const closeConstructPanel = useCallback(() => {
    dispatch({ type: 'SHOW_CONSTRUCT_PANEL', payload: false })
  }, [])

  const enterPlacementMode = useCallback((buildingTypeName: string) => {
    // Auto-switch to the correct base for this building type. Falls back to
    // category when the backend hasn't restarted to expose `base` (space /
    // defense category → space base).
    const bt = state.buildingTypes.find(t => t.name === buildingTypeName)
    if (bt) {
      let targetBase: BaseView = 'ground'
      if (bt.base === 'space') targetBase = 'space'
      else if (bt.base === 'ground') targetBase = 'ground'
      else if (bt.category === 'space' || bt.category === 'defense') targetBase = 'space'
      if (state.currentBase !== targetBase) {
        dispatch({ type: 'SET_BASE_VIEW', payload: targetBase })
      }
    }
    dispatch({ type: 'SET_PLACEMENT_MODE', payload: { active: true, buildingTypeName } })
  }, [state.buildingTypes, state.currentBase])

  const enterMoveMode = useCallback((buildingId: string) => {
    dispatch({ type: 'SET_PLACEMENT_MODE', payload: { active: true, movingBuildingId: buildingId } })
  }, [])

  const exitPlacementMode = useCallback(() => {
    dispatch({ type: 'SET_PLACEMENT_MODE', payload: { active: false } })
  }, [])

  const placeBuilding = useCallback((buildingId: string, position: GridPosition) => {
    dispatch({ type: 'PLACE_BUILDING', payload: { buildingId, position } })
  }, [])

  const focusCamera = useCallback((position: [number, number, number] | null) => {
    dispatch({ type: 'SET_CAMERA_TARGET', payload: position })
  }, [])

  const setBaseView = useCallback((base: BaseView) => {
    dispatch({ type: 'SET_BASE_VIEW', payload: base })
  }, [])

  return (
    <GameContext.Provider value={{
      state,
      dispatch,
      selectBuilding,
      deselectAll,
      openDetailPanel,
      closeDetailPanel,
      openConstructPanel,
      closeConstructPanel,
      enterPlacementMode,
      enterMoveMode,
      exitPlacementMode,
      placeBuilding,
      focusCamera,
      setBaseView,
    }}>
      {children}
    </GameContext.Provider>
  )
}

export function useGameContext() {
  const ctx = useContext(GameContext)
  if (!ctx) throw new Error('useGameContext must be used within GameProvider')
  return ctx
}

// Auto-assign grid positions for buildings that don't have one yet
export function autoAssignPositions(
  buildings: BuildingWithType[],
  existing: Record<string, GridPosition>,
): Record<string, GridPosition> {
  const positions = { ...existing }
  const occupied = new Set<string>()

  // Mark ALL tiles occupied by existing buildings
  for (const [buildingId, pos] of Object.entries(positions)) {
    const building = buildings.find(b => b.id === buildingId)
    if (building) {
      const tiles = getBuildingTiles(building.type_name, pos.col, pos.row)
      for (const t of tiles) {
        occupied.add(`${t.col},${t.row}`)
      }
    } else {
      occupied.add(`${pos.col},${pos.row}`)
    }
  }

  // Center of grid
  const centerCol = Math.floor(GRID_SIZE / 2)
  const centerRow = Math.floor(GRID_SIZE / 2)

  for (const b of buildings) {
    if (positions[b.id]) continue

    if (b.type_name === 'civic_center') {
      // Place civic center at center (account for 3x3 size, anchor at center-1)
      const anchor = { col: centerCol - 1, row: centerRow - 1 }
      positions[b.id] = anchor
      const tiles = getBuildingTiles(b.type_name, anchor.col, anchor.row)
      for (const t of tiles) {
        occupied.add(`${t.col},${t.row}`)
      }
      continue
    }

    const pos = findFreeTile(b.type_name, occupied, centerCol, centerRow)
    if (pos) {
      positions[b.id] = pos
      const tiles = getBuildingTiles(b.type_name, pos.col, pos.row)
      for (const t of tiles) {
        occupied.add(`${t.col},${t.row}`)
      }
    }
  }

  return positions
}

function findFreeTile(
  typeName: string,
  occupied: Set<string>,
  cx: number,
  cy: number,
): GridPosition | null {
  for (let r = 1; r < GRID_SIZE; r++) {
    for (let dx = -r; dx <= r; dx++) {
      for (let dy = -r; dy <= r; dy++) {
        if (Math.abs(dx) !== r && Math.abs(dy) !== r) continue
        const col = cx + dx
        const row = cy + dy
        if (canPlaceBuilding(typeName, col, row, occupied)) {
          return { col, row }
        }
      }
    }
  }
  return null
}

// Building abbreviations for labels
export const BUILDING_ABBREVIATIONS: Record<string, string> = {
  metal_collector: 'MC',
  he3_extractor: 'HE',
  residential_area: 'RA',
  resource_warehouse: 'RW',
  civic_center: 'CC',
  technology_center: 'TC',
  alliance_center: 'AC',
  trading_center: 'TR',
  galaxy_transporter: 'GT',
  compound_center: 'CP',
  radar: 'RD',
  ship_factory: 'SF',
  spacedock: 'SD',
  command_center: 'CM',
  weapon_research_center: 'WR',
  recycling_plant: 'RP',
  space_station: 'SS',
  meteor_star: 'MS',
  particle_cannon: 'PC',
  anti_aircraft_gun: 'AA',
  thors_cannon: 'TH',
  celestial_base: 'CB',
}

// Category colors for 3D materials
export const CATEGORY_COLORS: Record<string, string> = {
  resource: '#22cc66',
  core: '#4488ff',
  military: '#ff4444',
  space: '#44ccff',
  defense: '#44ccff',
}
