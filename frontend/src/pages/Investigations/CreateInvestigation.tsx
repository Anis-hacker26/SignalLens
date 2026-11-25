import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getEndpoints } from '../../services/endpoint.service'
import { createInvestigation } from '../../services/investigation.service'
import type { Endpoint } from '../../types/endpoint'
import './InvestigationsPage.css'

export default function CreateInvestigation() {
  const navigate = useNavigate()

  const [endpoints, setEndpoints] = useState<
    Endpoint[]
  >([])

  const [endpointId, setEndpointId] =
    useState('')

  const [title, setTitle] = useState('')

  const [description, setDescription] =
    useState('')

  const [loading, setLoading] =
    useState(true)

  const [saving, setSaving] =
    useState(false)

  const [error, setError] =
    useState<string | null>(null)

  useEffect(() => {
    async function loadEndpoints() {
      try {
        setLoading(true)
        setError(null)

        const data = await getEndpoints()

        setEndpoints(data)

        if (data.length > 0) {
          setEndpointId(data[0].id)
        }
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : 'Failed to load endpoints.',
        )
      } finally {
        setLoading(false)
      }
    }

    void loadEndpoints()
  }, [])

  async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault()

    if (!endpointId) {
      setError('Please select an endpoint.')
      return
    }

    if (!title.trim()) {
      setError('Investigation title is required.')
      return
    }

    try {
      setSaving(true)
      setError(null)

const investigation =
  await createInvestigation({
    endpointId,
    title: title.trim(),
    description:
      description.trim() || undefined,
    startedAt: new Date().toISOString(),
  })

      navigate(
        `/investigations/${investigation.id}`,
      )
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : 'Failed to create investigation.',
      )
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <section className="investigations-page">
        <div className="page-heading">
          <span>INVESTIGATIONS</span>
          <h2>New investigation</h2>
          <p>
            Start a new analyst-led endpoint
            investigation.
          </p>
        </div>

        <div className="investigations-state">
          Loading endpoints...
        </div>
      </section>
    )
  }

  return (
    <section className="investigations-page">
      <div className="investigation-detail-back">
        <Link to="/investigations">
          ← Back to investigations
        </Link>
      </div>

      <div className="page-heading">
        <span>INVESTIGATIONS</span>

        <h2>New investigation</h2>

        <p>
          Start a new analyst-led endpoint
          investigation.
        </p>
      </div>

      <form
        className="investigation-create-form"
        onSubmit={handleSubmit}
      >
        {error && (
          <div className="investigations-state investigations-error">
            {error}
          </div>
        )}

        <div className="investigation-form-field">
          <label htmlFor="endpoint">
            Endpoint
          </label>

          <select
            id="endpoint"
            value={endpointId}
            onChange={(event) =>
              setEndpointId(event.target.value)
            }
            disabled={saving}
            required
          >
            <option value="">
              Select an endpoint
            </option>

            {endpoints.map((endpoint) => (
              <option
                key={endpoint.id}
                value={endpoint.id}
              >
                {endpoint.hostname ??
                  endpoint.agentId}{' '}
                — {endpoint.agentId}
              </option>
            ))}
          </select>
        </div>

        <div className="investigation-form-field">
          <label htmlFor="title">
            Investigation title
          </label>

          <input
            id="title"
            type="text"
            placeholder="e.g. Suspicious process investigation"
            value={title}
            onChange={(event) =>
              setTitle(event.target.value)
            }
            disabled={saving}
            required
          />
        </div>

        <div className="investigation-form-field">
          <label htmlFor="description">
            Description
          </label>

          <textarea
            id="description"
            placeholder="Describe why this investigation was opened."
            value={description}
            onChange={(event) =>
              setDescription(event.target.value)
            }
            disabled={saving}
            rows={5}
          />
        </div>

        <div className="investigation-form-actions">
          <Link
            to="/investigations"
            className="investigation-form-cancel"
          >
            Cancel
          </Link>

          <button
            type="submit"
            disabled={saving || endpoints.length === 0}
          >
            {saving
              ? 'Creating investigation...'
              : 'Create investigation'}
          </button>
        </div>
      </form>
    </section>
  )
}