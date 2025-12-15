import { Link, useLocation } from 'react-router-dom'
import './Layout.css'

interface LayoutProps {
  children: React.ReactNode
}

function Layout({ children }: LayoutProps) {
  const location = useLocation()

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/'
    }
    return location.pathname === path || location.pathname.startsWith(path + '/')
  }

  return (
    <div className="layout">
      <header className="header">
        <nav className="nav">
          <Link to="/" className="logo">
            Grimoire
          </Link>
          <div className="nav-links">
            <Link to="/units" className={isActive('/units') ? 'active' : ''}>
              Units
            </Link>
            <Link to="/factions" className={isActive('/factions') ? 'active' : ''}>
              Factions
            </Link>
            <Link to="/search" className={isActive('/search') ? 'active' : ''}>
              Search
            </Link>
          </div>
        </nav>
      </header>
      <main className="main">{children}</main>
      <footer className="footer">
        <div className="footer-content">
          <p>Game data is provided for reference purposes only under fair use.</p>
          <p className="text-muted">
            This is an unofficial fan project and is not affiliated with Games Workshop.
          </p>
        </div>
      </footer>
    </div>
  )
}

export default Layout

