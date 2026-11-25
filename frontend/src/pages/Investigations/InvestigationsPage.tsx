import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { getInvestigations } from '../../services/investigation.service'
import type {
  Investigation,
  InvestigationStatus,
} from '../../types/investigation'
import './InvestigationsPage.css'

export default function InvestigationsPage() {
  const [investigations, setInvestigations] = useState<
    Investigation[]
  >([])

  const [search, setSearch] = useState('')

  const [status, setStatus] = useState<
    InvestigationStatus | 'ALL'
  >('ALL')

  const [loading, setLoading] = useState(true)

  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function loadInvestigations() {
      try {
        setLoading(true)
        setError(null)

        const data = await getInvestigations()

        setInvestigations(data)
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : 'Failed to load investigations.',
        )
      } finally {
        setLoading(false)
      }
    }

    void loadInvestigations()
  }, [])

  const filteredInvestigations = useMemo(() => {
    const query = search.trim().toLowerCase()

    return investigations.filter((investigation) => {
      const matchesStatus =
        status === 'ALL' ||
        investigation.status === status

      if (!query) {
        return matchesStatus
      }

      const searchable = [
        investigation.id,
        investigation.endpointId,
        investigation.title,
        investigation.description ?? '',
        investigation.status,
      ]
        .join(' ')
        .toLowerCase()

      return (
        matchesStatus &&
        searchable.includes(query)
      )
    })
  }, [investigations, search, status])

  if (loading) {
    return (
      <section className="investigations-page">
        <div className="page-heading">
          <span>INVESTIGATIONS</span>

          <h2>Security investigations</h2>

          <p>
            Review and manage analyst-led endpoint
            investigations.
          </p>
        </div>

        <div className="investigations-state">
          Loading investigations...
        </div>
      </section>
    )
  }

  if (error) {
    return (
      <section className="investigations-page">
        <div className="page-heading">
          <span>INVESTIGATIONS</span>

          <h2>Security investigations</h2>

          <p>
            Review and manage analyst-led endpoint
            investigations.
          </p>
        </div>

        <div className="investigations-state investigations-error">
          {error}
        </div>
      </section>
    )
  }

  return (
    <section className="investigations-page">
      <div className="page-heading">
        <div className="investigations-heading-row">
          <div>
            <span>INVESTIGATIONS</span>

            <h2>Security investigations</h2>

            <p>
              Review and manage analyst-led endpoint
              investigations.
            </p>
          </div>

          <Link
            to="/investigations/new"
            className="investigation-new-button"
          >
            + New investigation
          </Link>
        </div>
      </div>

      <div className="investigations-toolbar">
        <input
          type="search"
          placeholder="Search investigations..."
          value={search}
          onChange={(event) =>
            setSearch(event.target.value)
          }
        />

        <select
          value={status}
          onChange={(event) =>
            setStatus(
              event.target.value as
                | InvestigationStatus
                | 'ALL',
            )
          }
        >
          <option value="ALL">
            All statuses
          </option>

          <option value="OPEN">
            Open
          </option>

          <option value="CLOSED">
            Closed
          </option>
        </select>
      </div>

      <div className="investigations-table-wrapper">
        <table className="investigations-table">
          <thead>
            <tr>
              <th>Title</th>
              <th>Status</th>
              <th>Endpoint</th>
              <th>Started</th>
              <th>Closed</th>
            </tr>
          </thead>

          <tbody>
            {filteredInvestigations.map(
              (investigation) => (
                <tr
                  key={investigation.id}
                  className="investigation-row"
                >
                  <td>
                    <Link
                      to={`/investigations/${investigation.id}`}
                      className="investigation-title-link"
                    >
                      {investigation.title}
                    </Link>

                    {investigation.description && (
                      <span className="investigation-description">
                        {investigation.description}
                      </span>
                    )}
                  </td>

                  <td>
                    <span
                      className={`investigation-status investigation-status-${investigation.status.toLowerCase()}`}
                    >
                      {investigation.status}
                    </span>
                  </td>

                  <td>
                    <code>
                      {investigation.endpointId}
                    </code>
                  </td>

                  <td>
                    {new Date(
                      investigation.startedAt,
                    ).toLocaleString()}
                  </td>

                  <td>
                    {investigation.closedAt
                      ? new Date(
                          investigation.closedAt,
                        ).toLocaleString()
                      : '—'}
                  </td>
                </tr>
              ),
            )}
          </tbody>
        </table>

        {filteredInvestigations.length === 0 && (
          <div className="investigations-empty">
            No investigations match the current filters.
          </div>
        )}
      </div>
    </section>
  )
}