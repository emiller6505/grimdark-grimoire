import { Link } from 'react-router-dom'
import './Layout.css'

interface LayoutProps {
  children: React.ReactNode
}

function Layout({ children }: LayoutProps) {
  return (
    <div className="layout">
      <header className="header">
        <nav className="nav">
          <Link to="/" className="logo">
            Grimoire
          </Link>
          <div className="nav-links">
            <Link to="/units">Units</Link>
            <Link to="/catalogues">Catalogues</Link>
            <Link to="/factions">Factions</Link>
            <Link to="/search">Search</Link>
          </div>
        </nav>
      </header>
      <main className="main">{children}</main>
    </div>
  )
}

export default Layout

