import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import Home from './pages/Home'
import Units from './pages/Units'
import UnitDetail from './pages/UnitDetail'
import Factions from './pages/Factions'
import FactionCategory from './pages/FactionCategory'
import Search from './pages/Search'

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/units" element={<Units />} />
          <Route path="/units/:id" element={<UnitDetail />} />
          <Route path="/factions" element={<Factions />} />
          <Route path="/factions/:name/units" element={<Units />} />
          <Route path="/factions/:category" element={<FactionCategory />} />
          <Route path="/search" element={<Search />} />
        </Routes>
      </Layout>
    </Router>
  )
}

export default App

