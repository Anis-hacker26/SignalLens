import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import {
  closeInvestigation,
  getInvestigationById,
} from '../../services/investigation.service'
import {
  createAssessment,
  updateAssessment,
} from '../../services/assessment.service'
import type {
  AssessmentClassification,
  AssessmentConfidence,
} from '../../types/assessment'
import type {
  InvestigationDetails as InvestigationDetailsType,
} from '../../types/investigation'
import EvidenceDrawer from '../Evidence/components/EvidenceDrawer'
import './InvestigationsPage.css'

function formatDate(value: string | null) {
  if (!value) {
    return '—'
  }

  return new Date(value).toLocaleString()
}

function formatValue(value: unknown) {
  if (
    value === null ||
    value === undefined ||
    value === ''
  ) {
    return '—'
  }

  return String(value)
}

export default function InvestigationDetails() {
  const { id } = useParams<{ id: string }>()

  const [data, setData] =
    useState<InvestigationDetailsType | null>(null)

  const [loading, setLoading] =
    useState(true)

  const [error, setError] =
    useState<string | null>(null)

  const [selectedEvidenceId, setSelectedEvidenceId] =
    useState<string | null>(null)

  const [classification, setClassification] =
    useState<AssessmentClassification>(
      'UNDETERMINED',
    )

  const [confidence, setConfidence] =
    useState<AssessmentConfidence>(
      'MEDIUM',
    )

  const [reasoning, setReasoning] =
    useState('')

  const [assessmentSaving, setAssessmentSaving] =
    useState(false)

  const [assessmentError, setAssessmentError] =
    useState<string | null>(null)

  const [assessmentEditing, setAssessmentEditing] =
  useState(false)

  const [closing, setClosing] =
  useState(false)

const [closeError, setCloseError] =
  useState<string | null>(null)

  useEffect(() => {
    if (!id) {
      setError('Investigation ID is missing.')
      setLoading(false)
      return
    }

    const investigationId = id

    async function loadInvestigation() {
      try {
        setLoading(true)
        setError(null)

        const result =
          await getInvestigationById(
            investigationId,
          )

        setData(result)
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : 'Failed to load investigation.',
        )
      } finally {
        setLoading(false)
      }
    }

    void loadInvestigation()
  }, [id])

  if (loading) {
    return (
      <section className="investigations-page">
        <div className="investigation-detail-state">
          Loading investigation...
        </div>
      </section>
    )
  }

  if (error) {
    return (
      <section className="investigations-page">
        <div className="investigation-detail-back">
          <Link to="/investigations">
            ← Back to investigations
          </Link>
        </div>

        <div className="investigations-state investigations-error">
          {error}
        </div>
      </section>
    )
  }

  if (!data) {
    return null
  }

  const investigation = data.investigation
  const endpoint = data.endpoint
  const assessment = data.assessment

  return (
    <>
      <section className="investigations-page">
        <div className="investigation-detail-back">
          <Link to="/investigations">
            ← Back to investigations
          </Link>
        </div>

        <div className="investigation-detail-header">
          <div>
            <span className="investigation-detail-eyebrow">
              INVESTIGATION
            </span>

            <h2>{investigation.title}</h2>

            <p>
              {investigation.description ??
                'No investigation description provided.'}
            </p>
          </div>

          <span
            className={`investigation-status investigation-status-${investigation.status.toLowerCase()}`}
          >
            {investigation.status}
          </span>
        </div>

        <section className="investigation-detail-section">
          <div className="investigation-section-heading">
            <span>OVERVIEW</span>
          </div>

          <div className="investigation-detail-grid">
            <div className="investigation-detail-field">
              <span>Status</span>
              <strong>
                {investigation.status}
              </strong>
            </div>

            <div className="investigation-detail-field">
              <span>Started</span>
              <strong>
                {formatDate(
                  investigation.startedAt,
                )}
              </strong>
            </div>

            <div className="investigation-detail-field">
              <span>Closed</span>
              <strong>
                {formatDate(
                  investigation.closedAt,
                )}
              </strong>
            </div>

            <div className="investigation-detail-field">
              <span>Investigation ID</span>
              <strong>
                {investigation.id}
              </strong>
            </div>
          </div>
        </section>
        
        {investigation.status === 'OPEN' && (
  <section className="investigation-detail-section">
    <div className="investigation-section-heading">
      <span>INVESTIGATION ACTION</span>
    </div>

    {closeError && (
      <div className="investigations-state investigations-error">
        {closeError}
      </div>
    )}

    <button
      type="button"
      className="investigation-close-button"
      disabled={closing}
      onClick={async () => {
        const confirmed = window.confirm(
          'Are you sure you want to close this investigation? Closed investigations cannot be modified.',
        )

        if (!confirmed) {
          return
        }

        try {
          setClosing(true)
          setCloseError(null)

          const updated =
            await closeInvestigation(
              investigation.id,
            )

          setData((current) => {
            if (!current) {
              return current
            }

            return {
              ...current,
              investigation: updated,
            }
          })
        } catch (err) {
          setCloseError(
            err instanceof Error
              ? err.message
              : 'Failed to close investigation.',
          )
        } finally {
          setClosing(false)
        }
      }}
    >
      {closing
        ? 'Closing investigation...'
        : 'Close investigation'}
    </button>
  </section>
)}

        <section className="investigation-detail-section">
          <div className="investigation-section-heading">
            <span>ENDPOINT</span>
          </div>

          <div className="investigation-detail-grid">
            <div className="investigation-detail-field">
              <span>Agent ID</span>
              <strong>
                {formatValue(endpoint.agentId)}
              </strong>
            </div>

            <div className="investigation-detail-field">
              <span>Hostname</span>
              <strong>
                {formatValue(endpoint.hostname)}
              </strong>
            </div>

            <div className="investigation-detail-field">
              <span>Operating System</span>
              <strong>
                {formatValue(endpoint.osName)}
              </strong>
            </div>

            <div className="investigation-detail-field">
              <span>OS Version</span>
              <strong>
                {formatValue(endpoint.osVersion)}
              </strong>
            </div>

            <div className="investigation-detail-field investigation-detail-field-wide">
              <span>Endpoint ID</span>
              <strong>
                {endpoint.id}
              </strong>
            </div>
          </div>
        </section>

        <section className="investigation-detail-section">
          <div className="investigation-section-heading">
            <span>
              ATTACHED EVIDENCE
            </span>

            <small>
              {data.evidence.length} record
              {data.evidence.length === 1
                ? ''
                : 's'}
            </small>
          </div>

          {data.evidence.length === 0 ? (
            <div className="investigation-detail-empty">
              No evidence is attached to this
              investigation.
            </div>
          ) : (
            <div className="investigation-evidence-list">
              {data.evidence.map((evidence) => (
                <button
                  type="button"
                  className="investigation-evidence-item"
                  key={evidence.id}
                  onClick={() =>
                    setSelectedEvidenceId(
                      evidence.id,
                    )
                  }
                >
                  <div>
                    <strong>
                      {evidence.type}
                    </strong>

                    <span>
                      {evidence.source}
                      {' · Event '}
                      {evidence.sourceEventId}
                    </span>
                  </div>

                  <div>
                    <span>
                      Observed
                    </span>

                    <strong>
                      {formatDate(
                        evidence.observedAt,
                      )}
                    </strong>
                  </div>
                </button>
              ))}
            </div>
          )}
        </section>

<section className="investigation-detail-section">
  <div className="investigation-section-heading">
    <span>ANALYST ASSESSMENT</span>
  </div>

  {!assessment || assessmentEditing ? (
    <div className="investigation-assessment-form">
      <div className="investigation-detail-grid">
        <div className="investigation-detail-field">
          <label htmlFor="assessment-classification">
            Classification
          </label>

          <select
            id="assessment-classification"
            value={classification}
            onChange={(event) =>
              setClassification(
                event.target
                  .value as AssessmentClassification,
              )
            }
            disabled={assessmentSaving}
          >
            <option value="UNDETERMINED">
              UNDETERMINED
            </option>
            <option value="BENIGN">
              BENIGN
            </option>
            <option value="SUSPICIOUS">
              SUSPICIOUS
            </option>
            <option value="MALICIOUS">
              MALICIOUS
            </option>
          </select>
        </div>

        <div className="investigation-detail-field">
          <label htmlFor="assessment-confidence">
            Confidence
          </label>

          <select
            id="assessment-confidence"
            value={confidence}
            onChange={(event) =>
              setConfidence(
                event.target
                  .value as AssessmentConfidence,
              )
            }
            disabled={assessmentSaving}
          >
            <option value="LOW">LOW</option>
            <option value="MEDIUM">MEDIUM</option>
            <option value="HIGH">HIGH</option>
          </select>
        </div>

        <div className="investigation-detail-field investigation-detail-field-wide">
          <label htmlFor="assessment-reasoning">
            Analyst reasoning
          </label>

          <textarea
            id="assessment-reasoning"
            value={reasoning}
            onChange={(event) =>
              setReasoning(event.target.value)
            }
            placeholder="Explain the evidence and reasoning behind this assessment."
            rows={5}
            disabled={assessmentSaving}
          />
        </div>
      </div>

      {assessmentError && (
        <div className="investigations-state investigations-error">
          {assessmentError}
        </div>
      )}

      <div className="investigation-assessment-actions">
        <button
          type="button"
          className="investigation-assessment-submit"
          disabled={assessmentSaving}
          onClick={async () => {
            if (!id) {
              setAssessmentError(
                'Investigation ID is missing.',
              )
              return
            }

            try {
              setAssessmentSaving(true)
              setAssessmentError(null)

              const payload = {
                classification,
                confidence,
                reasoning,
              }

              const savedAssessment =
                assessmentEditing
                  ? await updateAssessment(
                      id,
                      payload,
                    )
                  : await createAssessment(
                      id,
                      payload,
                    )

              setData((current) => {
                if (!current) {
                  return current
                }

                return {
                  ...current,
                  assessment: savedAssessment,
                }
              })

              setAssessmentEditing(false)
            } catch (err) {
              setAssessmentError(
                err instanceof Error
                  ? err.message
                  : 'Failed to save assessment.',
              )
            } finally {
              setAssessmentSaving(false)
            }
          }}
        >
          {assessmentSaving
            ? 'Saving assessment...'
            : assessmentEditing
              ? 'Save changes'
              : 'Save assessment'}
        </button>

        {assessmentEditing && (
          <button
            type="button"
            className="investigation-assessment-cancel"
            disabled={assessmentSaving}
            onClick={() => {
              if (!assessment) {
                return
              }

              setClassification(
                assessment.classification,
              )

              setConfidence(
                assessment.confidence,
              )

              setReasoning(
                assessment.summary,
              )

              setAssessmentError(null)
              setAssessmentEditing(false)
            }}
          >
            Cancel
          </button>
        )}
      </div>
    </div>
  ) : (
    <div className="investigation-assessment">
      <div className="investigation-detail-grid">
        <div className="investigation-detail-field">
          <span>Classification</span>
          <strong>
            {assessment.classification}
          </strong>
        </div>

        <div className="investigation-detail-field">
          <span>Confidence</span>
          <strong>
            {assessment.confidence}
          </strong>
        </div>

        <div className="investigation-detail-field investigation-detail-field-wide">
          <span>Summary</span>
          <strong>
            {formatValue(
              assessment.summary,
            )}
          </strong>
        </div>

        <div className="investigation-detail-field">
          <span>Created</span>
          <strong>
            {formatDate(
              assessment.createdAt,
            )}
          </strong>
        </div>

        <div className="investigation-detail-field">
          <span>Updated</span>
          <strong>
            {formatDate(
              assessment.updatedAt,
            )}
          </strong>
        </div>
      </div>

      {investigation.status === 'OPEN' && (
        <button
          type="button"
          className="investigation-assessment-edit"
          onClick={() => {
            setClassification(
              assessment.classification,
            )

            setConfidence(
              assessment.confidence,
            )

            setReasoning(
              assessment.summary,
            )

            setAssessmentError(null)
            setAssessmentEditing(true)
          }}
        >
          Edit assessment
        </button>
      )}
    </div>
  )}
</section>
      </section>

      {selectedEvidenceId && (
        <EvidenceDrawer
          evidenceId={selectedEvidenceId}
          onClose={() =>
            setSelectedEvidenceId(null)
          }
        />
      )}
    </>
  )
}