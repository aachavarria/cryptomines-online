import { lazy, type ComponentType } from 'react'

const CivicCenterModel = lazy(() => import('./CivicCenterModel'))
const TechnologyCenterModel = lazy(() => import('./TechnologyCenterModel'))
const AllianceCenterModel = lazy(() => import('./AllianceCenterModel'))
const TradingCenterModel = lazy(() => import('./TradingCenterModel'))
const GalaxyTransporterModel = lazy(() => import('./GalaxyTransporterModel'))
const CompoundCenterModel = lazy(() => import('./CompoundCenterModel'))
const MetalCollectorModel = lazy(() => import('./MetalCollectorModel'))
const He3ExtractorModel = lazy(() => import('./He3ExtractorModel'))
const ResidentialAreaModel = lazy(() => import('./ResidentialAreaModel'))
const ResourceWarehouseModel = lazy(() => import('./ResourceWarehouseModel'))
const ShipFactoryBuildingModel = lazy(() => import('./ShipFactoryBuildingModel'))
const SpacedockBuildingModel = lazy(() => import('./SpacedockBuildingModel'))
const CommandCenterModel = lazy(() => import('./CommandCenterModel'))
const WeaponResearchCenterModel = lazy(() => import('./WeaponResearchCenterModel'))
const RadarModel = lazy(() => import('./RadarModel'))
const SpaceStationModel = lazy(() => import('./SpaceStationModel'))
const MeteorStarModel = lazy(() => import('./MeteorStarModel'))
const ParticleCannonModel = lazy(() => import('./ParticleCannonModel'))
const AntiAircraftGunModel = lazy(() => import('./AntiAircraftGunModel'))
const ThorsCannonModel = lazy(() => import('./ThorsCannonModel'))
const CelestialBaseModel = lazy(() => import('./CelestialBaseModel'))
const RecyclingPlantModel = lazy(() => import('./RecyclingPlantModel'))

export interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
  animate?: boolean
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
