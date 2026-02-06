import type { ComponentType } from 'react'
import CivicCenterModel from './CivicCenterModel'
import TechnologyCenterModel from './TechnologyCenterModel'
import AllianceCenterModel from './AllianceCenterModel'
import TradingCenterModel from './TradingCenterModel'
import GalaxyTransporterModel from './GalaxyTransporterModel'
import CompoundCenterModel from './CompoundCenterModel'
import MetalCollectorModel from './MetalCollectorModel'
import He3ExtractorModel from './He3ExtractorModel'
import ResidentialAreaModel from './ResidentialAreaModel'
import ResourceWarehouseModel from './ResourceWarehouseModel'
import ShipFactoryBuildingModel from './ShipFactoryBuildingModel'
import SpacedockBuildingModel from './SpacedockBuildingModel'
import CommandCenterModel from './CommandCenterModel'
import WeaponResearchCenterModel from './WeaponResearchCenterModel'
import RadarModel from './RadarModel'
import SpaceStationModel from './SpaceStationModel'
import MeteorStarModel from './MeteorStarModel'
import ParticleCannonModel from './ParticleCannonModel'
import AntiAircraftGunModel from './AntiAircraftGunModel'
import ThorsCannonModel from './ThorsCannonModel'
import CelestialBaseModel from './CelestialBaseModel'
import RecyclingPlantModel from './RecyclingPlantModel'

export interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

export const BUILDING_MODELS: Record<string, ComponentType<BuildingComponentProps>> = {
  civic_center: CivicCenterModel,
  technology_center: TechnologyCenterModel,
  alliance_center: AllianceCenterModel,
  trading_center: TradingCenterModel,
  galaxy_transporter: GalaxyTransporterModel,
  compound_center: CompoundCenterModel,
  metal_collector: MetalCollectorModel,
  he3_extractor: He3ExtractorModel,
  residential_area: ResidentialAreaModel,
  resource_warehouse: ResourceWarehouseModel,
  ship_factory: ShipFactoryBuildingModel,
  spacedock: SpacedockBuildingModel,
  command_center: CommandCenterModel,
  weapon_research_center: WeaponResearchCenterModel,
  radar: RadarModel,
  space_station: SpaceStationModel,
  meteor_star: MeteorStarModel,
  particle_cannon: ParticleCannonModel,
  anti_aircraft_gun: AntiAircraftGunModel,
  thors_cannon: ThorsCannonModel,
  celestial_base: CelestialBaseModel,
  recycling_plant: RecyclingPlantModel,
}
