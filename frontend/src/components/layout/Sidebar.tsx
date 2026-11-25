import { NavLink } from 'react-router-dom'

const navigationItems = [
  {
    label: 'Dashboard',
    path: '/',
  },
  {
    label: 'Evidence',
    path: '/evidence',
  },
  {
    label: 'Investigations',
    path: '/investigations',
  },
  {
    label: 'Endpoints',
    path: '/endpoints',
  },
  {
    label: 'Assessments',
    path: '/assessments',
  },
]

export default function Sidebar() {
  return (
    <aside className="sidebar">
      <div className="sidebar-brand">
        <div className="sidebar-logo">S</div>

        <div>
          <div className="sidebar-title">SignalLens</div>
          <div className="sidebar-subtitle">
            Security Operations
          </div>
        </div>
      </div>

      <nav className="sidebar-navigation">
        {navigationItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            className={({ isActive }) =>
              `sidebar-link ${isActive ? 'active' : ''}`
            }
          >
            {item.label}
          </NavLink>
        ))}
      </nav>
    </aside>
  )
}