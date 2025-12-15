import { useState } from 'react'
import { Link } from 'react-router-dom'
import { apiClient, Unit } from '../utils/api'
import './Search.css'

function Search() {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<Unit[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!query.trim()) return

    try {
      setLoading(true)
      setError(null)
      const response = await apiClient.search(query, 50)
      setResults(response.data.data?.results || [])
    } catch (err: any) {
      setError(err.message || 'Search failed')
      setResults([])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="search">
      <h1>Search Units</h1>
      <form onSubmit={handleSearch} className="search-form">
        <input
          type="text"
          placeholder="Search for units..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="search-input"
        />
        <button type="submit" disabled={loading || !query.trim()}>
          {loading ? 'Searching...' : 'Search'}
        </button>
      </form>

      {error && <div className="error">Error: {error}</div>}

      {results.length > 0 && (
        <div className="results">
          <h2>Results ({results.length})</h2>
          <div className="results-grid">
            {results.map((unit) => (
              <Link key={unit.id} to={`/units/${unit.id}`} className="result-card">
                <h3>{unit.name}</h3>
                <p className="unit-type">{unit.type}</p>
              </Link>
            ))}
          </div>
        </div>
      )}

      {!loading && query && results.length === 0 && !error && (
        <div className="no-results">No results found</div>
      )}
    </div>
  )
}

export default Search

