type OperationalSnapshotProps = {
  endpointCount: number
  evidenceCount: number
  investigationCount: number
  openInvestigationCount: number
}

export default function OperationalSnapshot({
  endpointCount,
  evidenceCount,
  investigationCount,
  openInvestigationCount,
}: OperationalSnapshotProps) {
  const items = [
    {
      label: 'Endpoints',
      value: endpointCount,
    },
    {
      label: 'Evidence',
      value: evidenceCount,
    },
    {
      label: 'Investigations',
      value: investigationCount,
    },
    {
      label: 'Open cases',
      value: openInvestigationCount,
    },
  ]

  return (
    <section className="dashboard-section">
      <div className="section-label">
        OPERATIONAL SNAPSHOT
      </div>

      <div className="snapshot-grid">
        {items.map((item) => (
          <div
            className="snapshot-item"
            key={item.label}
          >
            <span>{item.label}</span>
            <strong>{item.value}</strong>
          </div>
        ))}
      </div>
    </section>
  )
}