import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { apiClient, Faction } from '../utils/api'
import './FactionCategory.css'

interface FactionWithUnitCount extends Faction {
  unitCount?: number
}

function FactionCategory() {
  const { category } = useParams<{ category: string }>()
  const [factions, setFactions] = useState<FactionWithUnitCount[]>([])
  const [error, setError] = useState<string | null>(null)
  const [hasLoaded, setHasLoaded] = useState(false)

  useEffect(() => {
    if (category) {
      setHasLoaded(false)
      setFactions([])
      loadFactions()
    }
  }, [category])

  const extractIndividualFaction = (catalogueName: string): string => {
    // Extract individual faction name from catalogue name
    // "Xenos - Necrons" -> "Necrons"
    // "Imperium - Adeptus Astartes - Deathwatch" -> "Deathwatch"
    // "Chaos - Thousand Sons" -> "Thousand Sons"
    const parts = catalogueName.split(' - ')
    if (parts.length > 1) {
      // Return the last part (the actual faction name)
      return parts[parts.length - 1]
    }
    return catalogueName
  }

  const loadFactions = async () => {
    try {
      // Validate category - only allow xenos, imperium, chaos
      const validCategories = ['xenos', 'imperium', 'chaos']
      if (!category || !validCategories.includes(category.toLowerCase())) {
        setError(`Invalid category: ${category}. Valid categories are: ${validCategories.join(', ')}`)
        setFactions([])
        return
      }
      
      const response = await apiClient.listFactions()
      const allFactions = (response.data.data || []) as Faction[]
      
      // Find the top-level faction for this category (e.g., "Xenos", "Imperium", "Chaos")
      const categoryCapitalized = category.charAt(0).toUpperCase() + category.slice(1)
      const topLevelFaction = allFactions.find(f => f.name === categoryCapitalized)
      
      if (!topLevelFaction) {
        setFactions([])
        setError(`No faction found for category: ${categoryCapitalized}`)
        return
      }
      
      // Extract individual factions from catalogues
      // Group by individual faction name (e.g., "Necrons", "Orks")
      const individualFactionsMap = new Map<string, string[]>()
      
      topLevelFaction.catalogues.forEach((catalogueName) => {
        const individualFaction = extractIndividualFaction(catalogueName)
        if (!individualFactionsMap.has(individualFaction)) {
          individualFactionsMap.set(individualFaction, [])
        }
        individualFactionsMap.get(individualFaction)!.push(catalogueName)
      })
      
      // Convert to Faction array format
      const individualFactions: FactionWithUnitCount[] = Array.from(individualFactionsMap.entries()).map(([name, catalogues]) => ({
        name,
        catalogues,
        unitCount: undefined, // Will be loaded below
      }))
      
      // Sort alphabetically
      individualFactions.sort((a, b) => a.name.localeCompare(b.name))
      
      // Load unit counts for each faction
      const factionsWithCounts = await Promise.all(
        individualFactions.map(async (faction) => {
          try {
            const response = await apiClient.getFactionUnits(faction.name)
            return {
              ...faction,
              unitCount: response.data.data?.length || 0,
            }
          } catch (err) {
            // If fetching fails, just return with 0 count
            return {
              ...faction,
              unitCount: 0,
            }
          }
        })
      )
      
      setFactions(factionsWithCounts)
      setError(null)
      setHasLoaded(true)
    } catch (err: any) {
      setError(err.message || 'Failed to load factions')
      setHasLoaded(true)
    }
  }

  const getCategoryName = () => {
    if (!category) return ''
    return category.charAt(0).toUpperCase() + category.slice(1)
  }

  if (error) return <div className="error">Error: {error}</div>

  return (
    <div className="faction-category">
      <Link to="/factions" className="back-link">← Back to Factions</Link>
      <h1>{getCategoryName()} Factions</h1>
      {hasLoaded && factions.length === 0 ? (
        <p className="no-factions">No factions found in this category.</p>
      ) : (
        <div className="factions-grid">
          {factions.map((faction, index) => (
            <Link
              key={faction.name}
              to={`/factions/${encodeURIComponent(faction.name)}/units`}
              className="faction-card fade-in"
              style={{ animationDelay: `${index * 0.05}s` }}
            >
              <h3>{faction.name}</h3>
              <p className="unit-count">
                {faction.unitCount !== undefined ? (
                  <>
                    {faction.unitCount} unit{faction.unitCount !== 1 ? 's' : ''}
                  </>
                ) : (
                  'Loading...'
                )}
              </p>
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}

export default FactionCategory

