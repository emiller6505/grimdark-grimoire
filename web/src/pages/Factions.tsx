import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { apiClient, Faction } from '../utils/api'
import './Factions.css'

function Factions() {
  const [factions, setFactions] = useState<Faction[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadFactions()
  }, [])

  const loadFactions = async () => {
    try {
      setLoading(true)
      const response = await apiClient.listFactions()
      setFactions(response.data.data || [])
      setError(null)
    } catch (err: any) {
      setError(err.message || 'Failed to load factions')
    } finally {
      setLoading(false)
    }
  }

  if (loading) return <div className="loading">Loading factions...</div>
  if (error) return <div className="error">Error: {error}</div>

  return (
    <div className="factions">
      <h1>Factions</h1>
      <div className="factions-grid">
        {factions.map((faction) => (
          <Link
            key={faction.name}
            to={`/factions/${faction.name}/units`}
            className="faction-card"
          >
            <h3>{faction.name}</h3>
            <p className="catalogue-count">
              {faction.catalogues.length} catalogue{faction.catalogues.length !== 1 ? 's' : ''}
            </p>
          </Link>
        ))}
      </div>
    </div>
  )
}

export default Factions

