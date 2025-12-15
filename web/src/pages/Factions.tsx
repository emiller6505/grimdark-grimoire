import { Link } from 'react-router-dom'
import './Factions.css'

function Factions() {
  const categories = [
    {
      name: 'Xenos',
      description: 'Alien races and non-human factions',
      color: '#8B5CF6',
    },
    {
      name: 'Imperium',
      description: 'Human factions loyal to the Emperor',
      color: '#3B82F6',
    },
    {
      name: 'Chaos',
      description: 'Forces of the Dark Gods',
      color: '#EF4444',
    },
  ]

  return (
    <div className="factions">
      <h1>Factions</h1>
      <p className="factions-intro">Select a category to browse factions</p>
      <div className="factions-grid">
        {categories.map((category) => (
          <Link
            key={category.name}
            to={`/factions/${category.name.toLowerCase()}`}
            className="faction-card category-card"
            style={{ '--category-color': category.color } as React.CSSProperties}
          >
            <h3>{category.name}</h3>
            <p className="category-description">{category.description}</p>
          </Link>
        ))}
      </div>
    </div>
  )
}

export default Factions

