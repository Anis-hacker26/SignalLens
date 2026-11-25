import type { Evidence } from '../../../types/evidence'

type EvidenceTableProps = {
  evidence: Evidence[]
  page: number
  totalPages: number
  onSelect: (id: string) => void
  onPrevious: () => void
  onNext: () => void
}

function formatDate(value: string) {
  return new Date(value).toLocaleString()
}

export default function EvidenceTable({
  evidence,
  page,
  totalPages,
  onSelect,
  onPrevious,
  onNext,
}: EvidenceTableProps) {
  return (
    <div className="evidence-table-container">
      <table className="evidence-table">
        <thead>
          <tr>
            <th>Type</th>
            <th>Source</th>
            <th>Event ID</th>
            <th>Endpoint</th>
            <th>Observed</th>
          </tr>
        </thead>

        <tbody>
          {evidence.map((item) => (
            <tr
              key={item.id}
              onClick={() => onSelect(item.id)}
            >
              <td>
                <span className="evidence-type">
                  {item.type}
                </span>
              </td>

              <td>{item.source}</td>

              <td>{item.sourceEventId}</td>

              <td className="evidence-endpoint">
                {item.endpointId}
              </td>

              <td>
                {formatDate(item.observedAt)}
              </td>
            </tr>
          ))}

          {evidence.length === 0 && (
            <tr>
              <td
                colSpan={5}
                className="evidence-empty"
              >
                No evidence matches the current filters.
              </td>
            </tr>
          )}
        </tbody>
      </table>

      <div className="evidence-pagination">
        <button
          type="button"
          disabled={page === 1}
          onClick={onPrevious}
        >
          Previous
        </button>

        <span>
          Page {page} of {totalPages}
        </span>

        <button
          type="button"
          disabled={page === totalPages}
          onClick={onNext}
        >
          Next
        </button>
      </div>
    </div>
  )
}