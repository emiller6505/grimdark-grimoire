import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { apiClient, Unit } from '../utils/api'
import './UnitDetail.css'

function UnitDetail() {
  const { id } = useParams<{ id: string }>()
  const [unit, setUnit] = useState<Unit | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (id) {
      loadUnit(id)
    }
  }, [id])

  const loadUnit = async (unitId: string) => {
    try {
      setLoading(true)
      const response = await apiClient.getUnit(unitId)
      setUnit(response.data.data)
      setError(null)
    } catch (err: any) {
      setError(err.message || 'Failed to load unit')
    } finally {
      setLoading(false)
    }
  }

  if (loading) return <div className="loading">Loading unit...</div>
  if (error) return <div className="error">Error: {error}</div>
  if (!unit) return <div>Unit not found</div>

  return (
    <div className="unit-detail">
      <Link to="/units" className="back-link">← Back to Units</Link>
      <h1>{unit.name}</h1>
      <div className="unit-info">
        {unit.profiles?.unit && (
          <section className="profile-section">
            <h2>Unit Profile</h2>
            <div className="profile-grid">
              <div className="profile-stat">
                <span className="stat-label">Movement:</span>
                <span className="stat-value">{unit.profiles.unit.movement}</span>
              </div>
              <div className="profile-stat">
                <span className="stat-label">Toughness:</span>
                <span className="stat-value">{unit.profiles.unit.toughness}</span>
              </div>
              <div className="profile-stat">
                <span className="stat-label">Save:</span>
                <span className="stat-value">{unit.profiles.unit.save}</span>
              </div>
              <div className="profile-stat">
                <span className="stat-label">Wounds:</span>
                <span className="stat-value">{unit.profiles.unit.wounds}</span>
              </div>
              <div className="profile-stat">
                <span className="stat-label">Leadership:</span>
                <span className="stat-value">{unit.profiles.unit.leadership}</span>
              </div>
              <div className="profile-stat">
                <span className="stat-label">OC:</span>
                <span className="stat-value">{unit.profiles.unit.objectiveControl}</span>
              </div>
            </div>
          </section>
        )}

        {unit.profiles?.abilities && unit.profiles.abilities.length > 0 && (
          <section className="abilities-section">
            <h2>Abilities</h2>
            {unit.profiles.abilities.map((ability, idx) => (
              <div key={idx} className="ability">
                <h3>{ability.name}</h3>
                <p>{ability.description}</p>
              </div>
            ))}
          </section>
        )}

        {unit.weapons && (
          <section className="weapons-section">
            <h2>Weapons</h2>
            {unit.weapons.ranged && unit.weapons.ranged.length > 0 && (
              <div>
                <h3>Ranged Weapons</h3>
                <div className="weapons-list">
                  {unit.weapons.ranged.map((weapon: any, idx: number) => (
                    <div key={idx} className="weapon">
                      <h4>{weapon.name}</h4>
                      <div className="weapon-stats">
                        {weapon.range && <span>Range: {weapon.range}</span>}
                        {weapon.attacks && <span>Attacks: {weapon.attacks}</span>}
                        {weapon.strength && <span>S: {weapon.strength}</span>}
                        {weapon.armorPenetration && <span>AP: {weapon.armorPenetration}</span>}
                        {weapon.damage && <span>D: {weapon.damage}</span>}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
            {unit.weapons.melee && unit.weapons.melee.length > 0 && (
              <div>
                <h3>Melee Weapons</h3>
                <div className="weapons-list">
                  {unit.weapons.melee.map((weapon: any, idx: number) => (
                    <div key={idx} className="weapon">
                      <h4>{weapon.name}</h4>
                      <div className="weapon-stats">
                        {weapon.attacks && <span>Attacks: {weapon.attacks}</span>}
                        {weapon.strength && <span>S: {weapon.strength}</span>}
                        {weapon.armorPenetration && <span>AP: {weapon.armorPenetration}</span>}
                        {weapon.damage && <span>D: {weapon.damage}</span>}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </section>
        )}

        {unit.costs && (
          <section className="costs-section">
            <h2>Costs</h2>
            <div className="costs-list">
              {Object.entries(unit.costs).map(([type, cost]) => (
                <div key={type} className="cost">
                  <span className="cost-type">{type}:</span>
                  <span className="cost-value">{cost}</span>
                </div>
              ))}
            </div>
          </section>
        )}
      </div>
    </div>
  )
}

export default UnitDetail

