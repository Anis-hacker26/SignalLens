export type EvidenceType =
  | 'PROCESS'
  | 'NETWORK'
  | 'DNS'
  | 'FILE'
  | 'POWERSHELL'

export type Evidence = {
  id: string
  endpointId: string
  type: EvidenceType
  source: string
  sourceEventId: string
  observedAt: string
  collectedAt: string
}

export type EvidenceDetail = Evidence & {
  receivedAt?: string
  correlationId?: string | null
  processEvent?: Record<string, unknown> | null
  networkEvent?: Record<string, unknown> | null
  dnsEvent?: Record<string, unknown> | null
  fileEvent?: Record<string, unknown> | null
  powerShellEvent?: Record<string, unknown> | null
}