import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { apiClient, Unit } from '../utils/api'
import './Units.css'

function Units() {
  const [units, setUnits] = useState<Unit[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [faction, setFaction] = useState('')

  useEffect(() => {
    loadUnits()
  }, [faction, search])

  const loadUnits = async () => {
    try {
      setLoading(true)
      const response = await apiClient.listUnits({
        faction: faction || undefined,
        search: search || undefined,
        limit: 100,
      })
      setUnits(response.data.data || [])
      setError(null)
    } catch (err: any) {
      setError(err.message || 'Failed to load units')
    } finally {
      setLoading(false)
    }
  }

  if (loading) return <div className="loading">Loading units...</div>
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
      <div className="units-grid">
        {units.map((unit) => (
          <Link key={unit.id} to={`/units/${unit.id}`} className="unit-card">
            <h3>{unit.name}</h3>
            <p className="unit-type">{unit.type}</p>
            {unit.costs && (
              <p className="unit-cost">
                {Object.entries(unit.costs)
                  .map(([type, cost]) => `${cost} ${type}`)
                  .join(', ')}
              </p>
            )}
          </Link>
        ))}
      </div>
    </div>
  )
}

export default Units

