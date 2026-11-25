import { useEffect, useMemo, useState } from 'react'
import { getEvidence } from '../../services/evidence.service'
import type {
  Evidence,
  EvidenceType,
} from '../../types/evidence'
import EvidenceFilters from './components/EvidenceFilters'
import EvidenceTable from './components/EvidenceTable'
import EvidenceDrawer from './components/EvidenceDrawer'
import './EvidencePage.css'

const PAGE_SIZE = 50

export default function EvidencePage() {
  const [evidence, setEvidence] = useState<Evidence[]>([])
  const [search, setSearch] = useState('')
  const [type, setType] = useState<EvidenceType | 'ALL'>('ALL')
  const [selectedEvidenceId, setSelectedEvidenceId] =
    useState<string | null>(null)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function loadEvidence() {
      try {
        setLoading(true)
        setError(null)

        const data = await getEvidence()
        setEvidence(data)
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : 'Failed to load evidence.',
        )
      } finally {
        setLoading(false)
      }
    }

    void loadEvidence()
  }, [])

const filteredEvidence = useMemo(() => {
  const query = search.trim().toLowerCase()

  const sortedEvidence = [...evidence].sort(
    (a, b) =>
      new Date(b.observedAt).getTime() -
      new Date(a.observedAt).getTime(),
  )

  return sortedEvidence.filter((item) => {
      const matchesType =
        type === 'ALL' || item.type === type

      if (!query) {
        return matchesType
      }

      const searchable = [
        item.id,
        item.endpointId,
        item.type,
        item.source,
        item.sourceEventId,
        item.observedAt,
      ]
        .join(' ')
        .toLowerCase()

      return (
        matchesType &&
        searchable.includes(query)
      )
    })
  }, [evidence, search, type])

  const totalPages = Math.max(
    1,
    Math.ceil(filteredEvidence.length / PAGE_SIZE),
  )

  const currentPage = Math.min(page, totalPages)

  const paginatedEvidence = useMemo(() => {
    const start =
      (currentPage - 1) * PAGE_SIZE

    return filteredEvidence.slice(
      start,
      start + PAGE_SIZE,
    )
  }, [filteredEvidence, currentPage])

  useEffect(() => {
    setPage(1)
  }, [search, type])

  if (loading) {
    return (
      <section className="evidence-page">
        <div className="page-heading">
          <span>EVIDENCE</span>
          <h2>Collected endpoint evidence</h2>
        </div>

        <div className="evidence-state">
          Loading evidence...
        </div>
      </section>
    )
  }

  if (error) {
    return (
      <section className="evidence-page">
        <div className="page-heading">
          <span>EVIDENCE</span>
          <h2>Collected endpoint evidence</h2>
        </div>

        <div className="evidence-state evidence-error">
          {error}
        </div>
      </section>
    )
  }

  return (
    <section className="evidence-page">
      <div className="page-heading">
        <span>EVIDENCE</span>
        <h2>Collected endpoint evidence</h2>
        <p>
          Review normalized telemetry collected from
          SignalLens endpoints.
        </p>
      </div>

      <EvidenceFilters
        search={search}
        type={type}
        total={filteredEvidence.length}
        onSearchChange={setSearch}
        onTypeChange={setType}
      />

      <EvidenceTable
        evidence={paginatedEvidence}
        page={currentPage}
        totalPages={totalPages}
        onSelect={setSelectedEvidenceId}
        onPrevious={() =>
          setPage((current) =>
            Math.max(1, current - 1),
          )
        }
        onNext={() =>
          setPage((current) =>
            Math.min(totalPages, current + 1),
          )
        }
      />

      {selectedEvidenceId && (
        <EvidenceDrawer
          evidenceId={selectedEvidenceId}
          onClose={() =>
            setSelectedEvidenceId(null)
          }
        />
      )}
    </section>
  )
}