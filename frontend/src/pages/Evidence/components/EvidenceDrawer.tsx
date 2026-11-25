import { useEffect, useState } from 'react'
import { getEvidenceById } from '../../../services/evidence.service'
import type { EvidenceDetail } from '../../../types/evidence'
import './EvidenceDrawer.css'

type EvidenceDrawerProps = {
  evidenceId: string
  onClose: () => void
}

function formatLabel(value: string) {
  return value
    .replace(/([a-z])([A-Z])/g, '$1 $2')
    .replace(/_/g, ' ')
    .replace(/\b\w/g, (letter) => letter.toUpperCase())
}

function formatValue(value: unknown) {
  if (value === null || value === undefined || value === '') {
    return '—'
  }

  if (typeof value === 'boolean') {
    return value ? 'Yes' : 'No'
  }

  return String(value)
}

function DetailField({
  label,
  value,
  wide = false,
}: {
  label: string
  value: unknown
  wide?: boolean
}) {
  return (
    <div
      className={
        wide
          ? 'evidence-detail-field evidence-detail-field-wide'
          : 'evidence-detail-field'
      }
    >
      <span className="evidence-detail-label">
        {label}
      </span>

      <span className="evidence-detail-value">
        {formatValue(value)}
      </span>
    </div>
  )
}

function EventSection({
  title,
  data,
}: {
  title: string
  data: Record<string, unknown> | null | undefined
}) {
  if (!data) {
    return null
  }

  const entries = Object.entries(data)

  if (entries.length === 0) {
    return null
  }

  return (
    <section className="evidence-event-section">
      <div className="evidence-section-heading">
        <span>{title}</span>
      </div>

      <div className="evidence-event-grid">
        {entries.map(([key, value]) => (
          <DetailField
            key={key}
            label={formatLabel(key)}
            value={value}
            wide={
              typeof value === 'string' &&
              value.length > 70
            }
          />
        ))}
      </div>
    </section>
  )
}

export default function EvidenceDrawer({
  evidenceId,
  onClose,
}: EvidenceDrawerProps) {
  const [evidence, setEvidence] =
    useState<EvidenceDetail | null>(null)

  const [loading, setLoading] =
    useState(true)

  const [error, setError] =
    useState<string | null>(null)

  useEffect(() => {
    async function loadEvidence() {
      try {
        setLoading(true)
        setError(null)

        const data =
          await getEvidenceById(evidenceId)

        setEvidence(data)
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : 'Failed to load evidence details.',
        )
      } finally {
        setLoading(false)
      }
    }

    void loadEvidence()
  }, [evidenceId])

  return (
    <>
      <div
        className="evidence-drawer-backdrop"
        onClick={onClose}
        aria-hidden="true"
      />

      <aside
        className="evidence-drawer"
        aria-label="Evidence details"
      >
        <header className="evidence-drawer-header">
          <div>
            <span className="evidence-drawer-eyebrow">
              EVIDENCE DETAIL
            </span>

            <h3>
              {evidence
                ? `${evidence.type} evidence`
                : 'Evidence details'}
            </h3>

            {evidence && (
              <p className="evidence-drawer-subtitle">
                {evidence.source} · Event{' '}
                {evidence.sourceEventId}
              </p>
            )}
          </div>

          <button
            type="button"
            className="evidence-drawer-close"
            onClick={onClose}
            aria-label="Close evidence detail"
          >
            ×
          </button>
        </header>

        <div className="evidence-drawer-body">
          {loading && (
            <div className="evidence-drawer-state">
              Loading evidence details...
            </div>
          )}

          {error && (
            <div className="evidence-drawer-state evidence-drawer-error">
              {error}
            </div>
          )}

          {!loading && !error && evidence && (
            <>
              <section className="evidence-summary-section">
                <div className="evidence-section-heading">
                  <span>EVIDENCE</span>
                </div>

                <div className="evidence-summary-grid">
                  <DetailField
                    label="Type"
                    value={evidence.type}
                  />

                  <DetailField
                    label="Event ID"
                    value={evidence.sourceEventId}
                  />

                  <DetailField
                    label="Endpoint ID"
                    value={evidence.endpointId}
                    wide
                  />

                  <DetailField
                    label="Source"
                    value={evidence.source}
                  />

                  <DetailField
                    label="Observed"
                    value={evidence.observedAt}
                  />

                  <DetailField
                    label="Collected"
                    value={evidence.collectedAt}
                  />
                </div>
              </section>

              <EventSection
                title="PROCESS"
                data={evidence.processEvent}
              />

              <EventSection
                title="NETWORK"
                data={evidence.networkEvent}
              />

              <EventSection
                title="DNS"
                data={evidence.dnsEvent}
              />

              <EventSection
                title="FILE"
                data={evidence.fileEvent}
              />

              <EventSection
                title="POWERSHELL"
                data={evidence.powerShellEvent}
              />
            </>
          )}
        </div>
      </aside>
    </>
  )
}