// Building size configuration for multi-tile placement
// Each building occupies a rectangular footprint on the isometric grid

export interface BuildingSize {
  cols: number; // width in grid columns
  rows: number; // height in grid rows
}

export const BUILDING_SIZES: Record<string, BuildingSize> = {
  // 3x3 - Major structures
  civic_center: { cols: 3, rows: 3 },
  space_station: { cols: 3, rows: 2 },

  // 3x2 - Large facilities
  ship_factory: { cols: 3, rows: 2 },
  resource_warehouse: { cols: 3, rows: 3 },
  spacedock: { cols: 3, rows: 2 },
  technology_center: { cols: 3, rows: 2 },

  // 2x2 - Standard buildings
  metal_collector: { cols: 3, rows: 2 },
  he3_extractor: { cols: 2, rows: 2 },
  residential_area: { cols: 2, rows: 2 },
  alliance_center: { cols: 2, rows: 2 },
  trading_center: { cols: 2, rows: 2 },
  galaxy_transporter: { cols: 2, rows: 2 },
  compound_center: { cols: 2, rows: 2 },
  command_center: { cols: 2, rows: 2 },
  weapon_research_center: { cols: 2, rows: 2 },
  recycling_plant: { cols: 2, rows: 2 },
  thors_cannon: { cols: 2, rows: 2 },
  celestial_base: { cols: 2, rows: 2 },

  // 1x2 - Narrow structures
  particle_cannon: { cols: 1, rows: 2 },

  // 1x1 - Small structures
  radar: { cols: 1, rows: 1 },
  meteor_star: { cols: 1, rows: 1 },
  anti_aircraft_gun: { cols: 1, rows: 1 },
};

// Helper to get building size, defaults to 1x1 for unknown types
export function getBuildingSize(typeName: string): BuildingSize {
  return BUILDING_SIZES[typeName] || { cols: 1, rows: 1 };
}
