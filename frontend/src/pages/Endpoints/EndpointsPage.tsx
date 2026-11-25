import { useEffect, useState } from 'react'
import { getEndpoints } from '../../services/endpoint.service'
import type { Endpoint } from '../../types/endpoint'

import LoadingState from '../../components/common/LoadingState'
import ErrorBanner from '../../components/common/ErrorBanner'
import EmptyState from '../../components/common/EmptyState'

import './EndpointsPage.css'

function formatDate(value: string): string {
  return new Date(value).toLocaleString()
}

export default function EndpointsPage() {
  const [endpoints, setEndpoints] = useState<Endpoint[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    async function loadEndpoints() {
      try {
        setLoading(true)
        setError(null)

        const data = await getEndpoints()

        if (!cancelled) {
          setEndpoints(data)
        }
      } catch (err) {
        if (!cancelled) {
          setError(
            err instanceof Error
              ? err.message
              : 'Failed to load endpoints.',
          )
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    loadEndpoints()

    return () => {
      cancelled = true
    }
  }, [])

  if (loading) {
    return (
      <section className="page endpoints-page">
        <div className="page-header">
          <div>
            <h1>Endpoints</h1>
            <p>Windows endpoints registered with SignalLens.</p>
          </div>
        </div>

        <LoadingState />
      </section>
    )
  }

  if (error) {
    return (
      <section className="page endpoints-page">
        <div className="page-header">
          <div>
            <h1>Endpoints</h1>
            <p>Windows endpoints registered with SignalLens.</p>
          </div>
        </div>

        <ErrorBanner message={error} />
      </section>
    )
  }

  return (
    <section className="page endpoints-page">
      <div className="page-header">
        <div>
          <h1>Endpoints</h1>
          <p>Windows endpoints registered with SignalLens.</p>
        </div>
      </div>

      {endpoints.length === 0 ? (
        <EmptyState message="No endpoints have been registered yet." />
      ) : (
        <div className="endpoints-table-container">
          <table className="endpoints-table">
            <thead>
              <tr>
                <th>Hostname</th>
                <th>Agent ID</th>
                <th>Operating System</th>
                <th>Agent Version</th>
                <th>First Seen</th>
                <th>Last Seen</th>
              </tr>
            </thead>

            <tbody>
              {endpoints.map((endpoint) => (
                <tr key={endpoint.id}>
                  <td className="endpoint-hostname">
                    {endpoint.hostname}
                  </td>

                  <td>{endpoint.agentId}</td>

                  <td>
                    <div>{endpoint.osName}</div>
                    <span className="endpoint-os-version">
                      {endpoint.osVersion}
                    </span>
                  </td>

                  <td>{endpoint.agentVersion}</td>

                  <td>{formatDate(endpoint.firstSeenAt)}</td>

                  <td>{formatDate(endpoint.lastSeenAt)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  )
}