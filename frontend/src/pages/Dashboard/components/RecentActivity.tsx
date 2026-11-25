import type { Evidence } from '../../../types/evidence'

type RecentActivityProps = {
  evidence: Evidence[]
}

function formatDate(value: string) {
  return new Date(value).toLocaleString()
}

export default function RecentActivity({
  evidence,
}: RecentActivityProps) {
  return (
    <section className="dashboard-section">
      <div className="section-label">
        RECENT ACTIVITY
      </div>

      {evidence.length === 0 ? (
        <div className="activity-empty">
          No evidence has been collected yet.
        </div>
      ) : (
        <div className="activity-list">
          {evidence.map((item) => (
            <div
              className="activity-row"
              key={item.id}
            >
              <div className="activity-main">
                <span className="activity-type">
                  {item.type}
                </span>

                <span className="activity-source">
                  {item.source}
                </span>
              </div>

              <div className="activity-meta">
                <span>
                  Event {item.sourceEventId}
                </span>

                <span>
                  {formatDate(item.observedAt)}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </section>
  )
}