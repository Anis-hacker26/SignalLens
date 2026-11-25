export default function Header() {
  return (
    <header className="app-header">
      <div>
        <h1>SignalLens</h1>
        <p>Windows endpoint threat detection and investigation</p>
      </div>

      <div className="header-status">
        <span className="status-dot" />
        <span>System operational</span>
      </div>
    </header>
  )
}