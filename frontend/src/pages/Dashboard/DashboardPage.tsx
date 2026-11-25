import { useEffect, useMemo, useState } from 'react'
import { getEvidence } from '../../services/evidence.service'
import { getInvestigations } from '../../services/investigation.service'
import { getEndpoints } from '../../services/endpoint.service'
import type { Evidence } from '../../types/evidence'
import type { Investigation } from '../../types/investigation'
import OperationalSnapshot from './components/OperationalSnapshot'
import EvidenceOverview from './components/EvidenceOverview'
import InvestigationOverview from './components/InvestigationOverview'
import RecentActivity from './components/RecentActivity'
import './DashboardPage.css'

export default function DashboardPage() {
  const [evidence, setEvidence] = useState<Evidence[]>([])
  const [investigations, setInvestigations] = useState<
    Investigation[]
  >([])
  const [endpointCount, setEndpointCount] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function loadDashboard() {
      try {
        setLoading(true)
        setError(null)

const [
  evidenceData,
  investigationData,
  endpointData,
] = await Promise.all([
  getEvidence(),
  getInvestigations(),
  getEndpoints(),
])

setEvidence(evidenceData)
setInvestigations(investigationData)
setEndpointCount(endpointData.length)
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : 'Failed to load dashboard data.',
        )
      } finally {
        setLoading(false)
      }
    }

    void loadDashboard()
  }, [])


  const openInvestigationCount = useMemo(() => {
    return investigations.filter(
      (investigation) =>
        investigation.status === 'OPEN',
    ).length
  }, [investigations])

  const evidenceCounts = useMemo(() => {
    return {
      PROCESS: evidence.filter(
        (item) => item.type === 'PROCESS',
      ).length,
      NETWORK: evidence.filter(
        (item) => item.type === 'NETWORK',
      ).length,
      DNS: evidence.filter(
        (item) => item.type === 'DNS',
      ).length,
      FILE: evidence.filter(
        (item) => item.type === 'FILE',
      ).length,
      POWERSHELL: evidence.filter(
        (item) => item.type === 'POWERSHELL',
      ).length,
    }
  }, [evidence])

  const recentEvidence = useMemo(() => {
    return [...evidence]
      .sort(
        (a, b) =>
          new Date(b.observedAt).getTime() -
          new Date(a.observedAt).getTime(),
      )
      .slice(0, 8)
  }, [evidence])

  if (loading) {
    return (
      <section className="dashboard-page">
        <div className="dashboard-heading">
          <span>SECURITY OPERATIONS</span>
          <h2>Endpoint investigation overview</h2>
        </div>

        <div className="dashboard-state">
          Loading operational data...
        </div>
      </section>
    )
  }

  if (error) {
    return (
      <section className="dashboard-page">
        <div className="dashboard-heading">
          <span>SECURITY OPERATIONS</span>
          <h2>Endpoint investigation overview</h2>
        </div>

        <div className="dashboard-state dashboard-error">
          {error}
        </div>
      </section>
    )
  }

  return (
    <section className="dashboard-page">
      <div className="dashboard-heading">
        <span>SECURITY OPERATIONS</span>
        <h2>Endpoint investigation overview</h2>
        <p>
          Evidence and investigation activity across connected
          Windows endpoints.
        </p>
      </div>

      <OperationalSnapshot
        endpointCount={endpointCount}
        evidenceCount={evidence.length}
        investigationCount={investigations.length}
        openInvestigationCount={openInvestigationCount}
      />

      <EvidenceOverview counts={evidenceCounts} />

      <InvestigationOverview
        investigations={investigations}
      />

      <RecentActivity evidence={recentEvidence} />
    </section>
  )
}