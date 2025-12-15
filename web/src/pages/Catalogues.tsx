import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { apiClient, Catalogue } from '../utils/api'
import './Catalogues.css'

function Catalogues() {
  const [catalogues, setCatalogues] = useState<Catalogue[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadCatalogues()
  }, [])

  const loadCatalogues = async () => {
    try {
      setLoading(true)
      const response = await apiClient.listCatalogues()
      setCatalogues(response.data.data || [])
      setError(null)
    } catch (err: any) {
      setError(err.message || 'Failed to load catalogues')
    } finally {
      setLoading(false)
    }
  }

  if (loading) return <div className="loading">Loading catalogues...</div>
  if (error) return <div className="error">Error: {error}</div>

  return (
    <div className="catalogues">
      <h1>Catalogues</h1>
      <div className="catalogues-grid">
        {catalogues.map((catalogue) => (
          <Link
            key={catalogue.id}
            to={`/catalogues/${catalogue.id}`}
            className="catalogue-card"
          >
            <h3>{catalogue.name}</h3>
            <p className="catalogue-revision">Revision: {catalogue.revision}</p>
            {catalogue.library && <span className="library-badge">Library</span>}
          </Link>
        ))}
      </div>
    </div>
  )
}

export default Catalogues

