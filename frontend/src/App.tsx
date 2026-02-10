import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { GameProvider } from './contexts/GameContext.tsx'
import Home from './pages/Home.tsx'
import Planet from './pages/Planet.tsx'
import Military from './pages/Military.tsx'
import Inventory from './pages/Inventory.tsx'

export default function App() {
  return (
    <BrowserRouter>
      <GameProvider>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/planet/:id" element={<Planet />} />
          <Route path="/military" element={<Military />} />
          <Route path="/inventory" element={<Inventory />} />
        </Routes>
      </GameProvider>
    </BrowserRouter>
  )
}
