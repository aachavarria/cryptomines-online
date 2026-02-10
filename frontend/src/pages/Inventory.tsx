import { useNavigate } from 'react-router-dom'
import InventoryPanel from '../components/panels/InventoryPanel'
import './Inventory.css'

export default function Inventory() {
  const navigate = useNavigate()

  return (
    <div className="inventory-page">
      <div className="page-header">
        <button className="back-btn" onClick={() => navigate(-1)}>
          ← Back
        </button>
        <h1>Inventory</h1>
      </div>
      <div className="page-content">
        <InventoryPanel />
      </div>
    </div>
  )
}
