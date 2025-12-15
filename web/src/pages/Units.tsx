import { useState, useEffect } from 'react'
import { Link, useParams } from 'react-router-dom'
import { apiClient, Unit } from '../utils/api'
import './Units.css'

function Units() {
  const { name: factionParam } = useParams<{ name?: string }>()
  const [units, setUnits] = useState<Unit[]>([])
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [faction, setFaction] = useState(factionParam || '')
  const [hasLoaded, setHasLoaded] = useState(false)

  useEffect(() => {
    if (factionParam) {
      setFaction(factionParam)
    }
  }, [factionParam])

  useEffect(() => {
    setHasLoaded(false)
    setUnits([])
    loadUnits()
  }, [faction, search])

  const loadUnits = async () => {
    try {
      let response
      
      // If we have a faction from the route, use the faction units endpoint
      if (factionParam && factionParam === faction) {
        try {
          response = await apiClient.getFactionUnits(factionParam)
        } catch (err: any) {
          // Fall back to listUnits if getFactionUnits fails
          response = await apiClient.listUnits({
            faction: faction || undefined,
            search: search || undefined,
            limit: 100,
          })
        }
      } else {
        // Otherwise use the regular list endpoint
        response = await apiClient.listUnits({
          faction: faction || undefined,
          search: search || undefined,
          limit: 100,
        })
      }
      
      setUnits(response.data.data || [])
      setError(null)
      setHasLoaded(true)
    } catch (err: any) {
      setError(err.message || 'Failed to load units')
      setHasLoaded(true)
    }
  }

  if (error) return <div className="error">Error: {error}</div>

  return (
    <div className="units">
      <h1>Units</h1>
      <div className="filters">
        <input
          type="text"
          placeholder="Search units..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="search-input"
        />
        <input
          type="text"
          placeholder="Filter by faction..."
          value={faction}
          onChange={(e) => setFaction(e.target.value)}
          className="faction-input"
        />
      </div>
      {hasLoaded && units.length === 0 ? (
        <p className="no-units">No units found.</p>
      ) : (
        <div className="units-grid">
          {units.map((unit, index) => (
            <Link 
              key={unit.id} 
              to={`/units/${unit.id}`} 
              className="unit-card fade-in"
              style={{ animationDelay: `${index * 0.03}s` }}
            >
              <h3>{unit.name}</h3>
              <p className="unit-type">{unit.type}</p>
              {unit.tieredCosts ? (
                <p className="unit-cost">
                  {unit.tieredCosts.baseCost}
                  {unit.tieredCosts.tiers.length > 0 && (
                    <> / {unit.tieredCosts.tiers.map(tier => `${tier.minModels}+: ${tier.cost}`).join(' / ')}</>
                  )} pts
                </p>
              ) : unit.costs && unit.costs.pts !== undefined ? (
                <p className="unit-cost">
                  {unit.costs.pts} pts
                </p>
              ) : null}
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}

export default Units

