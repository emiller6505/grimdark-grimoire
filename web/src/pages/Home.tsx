import { Link } from 'react-router-dom'
import './Home.css'

function Home() {
  return (
    <div className="home">
      <h1>Grimoire</h1>
      <p className="subtitle">Warhammer 40,000 10th Edition API</p>
      <p className="description">
        Explore units, catalogues, factions, and more from the Warhammer 40K 10th Edition
        BattleScribe data.
      </p>
      <div className="quick-links">
        <Link to="/units" className="card">
          <h2>Units</h2>
          <p>Browse all units</p>
        </Link>
        <Link to="/catalogues" className="card">
          <h2>Catalogues</h2>
          <p>View all catalogues</p>
        </Link>
        <Link to="/factions" className="card">
          <h2>Factions</h2>
          <p>Explore factions</p>
        </Link>
        <Link to="/search" className="card">
          <h2>Search</h2>
          <p>Search units</p>
        </Link>
      </div>
    </div>
  )
}

export default Home

