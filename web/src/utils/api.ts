import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface CostTier {
  minModels: number
  cost: number
}

export interface TieredCosts {
  baseCost: number
  tiers: CostTier[]
}

export interface Unit {
  id: string
  name: string
  type: string
  profiles?: {
    unit?: {
      movement: string
      toughness: number
      save: string
      wounds: number
      leadership: string
      objectiveControl: number
    }
    abilities?: Array<{
      name: string
      description: string
    }>
  }
  weapons?: {
    ranged: Array<any>
    melee: Array<any>
  }
  categories?: Array<{
    id: string
    name: string
    primary: boolean
  }>
  costs?: Record<string, number>
  tieredCosts?: TieredCosts
}

export interface Catalogue {
  id: string
  name: string
  revision: string
  library: boolean
}

export interface Faction {
  name: string
  catalogues: string[]
}

export const apiClient = {
  // Game System
  getGameSystem: () => api.get('/game-system'),

  // Catalogues
  listCatalogues: () => api.get('/catalogues'),
  getCatalogue: (id: string) => api.get(`/catalogues/${id}`),
  getCatalogueUnits: (id: string) => api.get(`/catalogues/${id}/units`),

  // Units
  listUnits: (params?: {
    faction?: string
    category?: string
    search?: string
    limit?: number
    offset?: number
  }) => api.get('/units', { params }),
  getUnit: (id: string) => api.get(`/units/${id}`),
  getUnitWeapons: (id: string) => api.get(`/units/${id}/weapons`),

  // Factions
  listFactions: () => api.get('/factions'),
  getFactionUnits: (name: string) => api.get(`/factions/${name}/units`),

  // Search
  search: (query: string, limit?: number) =>
    api.get('/search', { params: { q: query, limit } }),
}

export default api

