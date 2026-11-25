import type { EvidenceType } from '../../../types/evidence'

type EvidenceFiltersProps = {
  search: string
  type: EvidenceType | 'ALL'
  total: number
  onSearchChange: (value: string) => void
  onTypeChange: (
    value: EvidenceType | 'ALL',
  ) => void
}

const types: Array<EvidenceType | 'ALL'> = [
  'ALL',
  'PROCESS',
  'NETWORK',
  'DNS',
  'FILE',
  'POWERSHELL',
]

export default function EvidenceFilters({
  search,
  type,
  total,
  onSearchChange,
  onTypeChange,
}: EvidenceFiltersProps) {
  return (
    <div className="evidence-toolbar">
      <div className="evidence-search">
        <input
          type="search"
          value={search}
          placeholder="Search evidence..."
          onChange={(event) =>
            onSearchChange(event.target.value)
          }
        />
      </div>

      <select
        value={type}
        onChange={(event) =>
          onTypeChange(
            event.target.value as EvidenceType | 'ALL',
          )
        }
      >
        {types.map((item) => (
          <option key={item} value={item}>
            {item === 'ALL'
              ? 'All evidence'
              : item}
          </option>
        ))}
      </select>

      <span className="evidence-count">
        {total.toLocaleString()} records
      </span>
    </div>
  )
}