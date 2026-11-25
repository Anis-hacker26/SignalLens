import type { Investigation } from '../../../types/investigation'

type InvestigationOverviewProps = {
  investigations: Investigation[]
}

export default function InvestigationOverview({
  investigations,
}: InvestigationOverviewProps) {
  const open = investigations.filter(
    (item) => item.status === 'OPEN',
  ).length

  const closed = investigations.filter(
    (item) => item.status === 'CLOSED',
  ).length

  return (
    <section className="dashboard-section">
      <div className="section-label">
        INVESTIGATION WORKLOAD
      </div>

      <div className="workload-grid">
        <div className="workload-item">
          <span>Open</span>
          <strong>{open}</strong>
        </div>

        <div className="workload-item">
          <span>Closed</span>
          <strong>{closed}</strong>
        </div>

        <div className="workload-item">
          <span>Total</span>
          <strong>{investigations.length}</strong>
        </div>
      </div>
    </section>
  )
}