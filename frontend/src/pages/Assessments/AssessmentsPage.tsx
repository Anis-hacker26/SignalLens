import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import {
  getInvestigations,
  getInvestigationById,
} from '../../services/investigation.service'
import type { Assessment } from '../../types/assessment'
import type { Investigation } from '../../types/investigation'
import './AssessmentsPage.css'

type AssessmentRow = {
  investigation: Investigation
  assessment: Assessment
}

function formatDate(value: string) {
  return new Date(value).toLocaleString()
}

export default function AssessmentsPage() {
  const [rows, setRows] = useState<AssessmentRow[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    async function loadAssessments() {
      try {
        setLoading(true)
        setError(null)

        const investigations = await getInvestigations()

        const details = await Promise.all(
          investigations.map((investigation) =>
            getInvestigationById(investigation.id),
          ),
        )

        if (cancelled) {
          return
        }

        const assessmentRows = details
          .filter(
            (
              detail,
            ): detail is typeof detail & {
              assessment: Assessment
            } => detail.assessment !== null,
          )
          .map((detail) => ({
            investigation: detail.investigation,
            assessment: detail.assessment,
          }))

        setRows(assessmentRows)
      } catch (err) {
        if (cancelled) {
          return
        }

        setError(
          err instanceof Error
            ? err.message
            : 'Failed to load assessments.',
        )
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    loadAssessments()

    return () => {
      cancelled = true
    }
  }, [])

  return (
    <main className="assessments-page">
      <header className="assessments-header">
        <div>
          <span className="assessments-eyebrow">
            ANALYST REVIEW
          </span>

          <h1>Assessments</h1>

          <p>
            Analyst classifications and reasoning recorded
            against investigations.
          </p>
        </div>
      </header>

      {loading && (
        <div className="assessments-state">
          Loading assessments...
        </div>
      )}

      {error && (
        <div className="assessments-state assessments-state-error">
          {error}
        </div>
      )}

      {!loading && !error && rows.length === 0 && (
        <div className="assessments-state">
          No analyst assessments have been recorded.
        </div>
      )}

      {!loading && !error && rows.length > 0 && (
        <section className="assessments-table-section">
          <div className="assessments-section-heading">
            <span>RECORDED ASSESSMENTS</span>

            <small>
              {rows.length} record
              {rows.length === 1 ? '' : 's'}
            </small>
          </div>

          <div className="assessments-table-wrapper">
            <table className="assessments-table">
              <thead>
                <tr>
                  <th>Investigation</th>
                  <th>Classification</th>
                  <th>Confidence</th>
                  <th>Summary</th>
                  <th>Updated</th>
                </tr>
              </thead>

              <tbody>
                {rows.map(({ investigation, assessment }) => (
                  <tr key={assessment.id}>
                    <td>
                      <Link
                        to={`/investigations/${investigation.id}`}
                        className="assessment-investigation-link"
                      >
                        {investigation.title}
                      </Link>
                    </td>

                    <td>
                      <span className="assessment-value">
                        {assessment.classification}
                      </span>
                    </td>

                    <td>
                      <span className="assessment-value">
                        {assessment.confidence}
                      </span>
                    </td>

                    <td className="assessment-summary">
                      {assessment.summary || '—'}
                    </td>

                    <td>
                      {formatDate(assessment.updatedAt)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>
      )}
    </main>
  )
}
