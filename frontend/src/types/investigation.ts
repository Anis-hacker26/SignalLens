import type { Assessment } from './assessment'
import type { Evidence } from './evidence'

export type InvestigationStatus =
  | 'OPEN'
  | 'CLOSED'

export type Investigation = {
  id: string
  endpointId: string
  title: string
  description: string | null
  status: InvestigationStatus
  startedAt: string
  closedAt: string | null
  createdAt: string
  updatedAt: string
}

export type InvestigationDetails = {
  investigation: Investigation
  endpoint: {
    id: string
    agentId: string
    hostname?: string
    osName?: string
    osVersion?: string
  }
  evidence: Array<
    Evidence & {
      processEvent?: Record<string, unknown>
      networkEvent?: Record<string, unknown>
    }
  >
  assessment: Assessment | null
}