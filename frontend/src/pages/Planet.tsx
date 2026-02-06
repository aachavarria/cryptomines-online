import { useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { getPlanet, listBuildings, getResources, getBuildingTypes } from '../services/api.ts'
import { useGameContext } from '../contexts/GameContext.tsx'
import GameShell from '../components/layout/GameShell.tsx'

export default function Planet() {
  const { id } = useParams<{ id: string }>()
  const { dispatch } = useGameContext()

  useEffect(() => {
    if (!id) return

    async function loadPlanetData() {
      dispatch({ type: 'SET_LOADING', payload: true })
      localStorage.setItem('current_planet_id', id!)
      try {
        const [planet, buildings, resources, buildingTypes] = await Promise.all([
          getPlanet(id!),
          listBuildings(id!),
          getResources(id!),
          getBuildingTypes(),
        ])
        dispatch({ type: 'SET_CURRENT_PLANET', payload: planet })
        dispatch({ type: 'SET_BUILDINGS', payload: buildings })
        dispatch({ type: 'SET_RESOURCES', payload: resources })
        dispatch({ type: 'SET_BUILDING_TYPES', payload: buildingTypes })
        dispatch({ type: 'SET_ERROR', payload: null })
      } catch {
        dispatch({ type: 'SET_ERROR', payload: 'Failed to load planet data' })
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false })
      }
    }

    loadPlanetData()
  }, [id, dispatch])

  return <GameShell />
}
