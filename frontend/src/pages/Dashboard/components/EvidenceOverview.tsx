import type { EvidenceType } from '../../../types/evidence'

type EvidenceOverviewProps = {
  counts: Record<EvidenceType, number>
}

const evidenceTypes: EvidenceType[] = [
  'PROCESS',
  'NETWORK',
  'DNS',
  'FILE',
  'POWERSHELL',
]

export default function EvidenceOverview({
  counts,
}: EvidenceOverviewProps) {
  return (
    <section className="dashboard-section">
      <div className="section-label">
        EVIDENCE OVERVIEW
      </div>

      <div className="evidence-overview-grid">
        {evidenceTypes.map((type) => (
          <div
            className="evidence-overview-item"
            key={type}
          >
            <span>{type}</span>
            <strong>{counts[type]}</strong>
          </div>
        ))}
      </div>
    </section>
  )
}