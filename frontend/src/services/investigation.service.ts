import { apiRequest } from './api'
import type {
  Investigation,
  InvestigationDetails,
} from '../types/investigation'

export function getInvestigations(): Promise<
  Investigation[]
> {
  return apiRequest<Investigation[]>(
    '/investigations',
  )
}

export function getInvestigationById(
  id: string,
): Promise<InvestigationDetails> {
  return apiRequest<InvestigationDetails>(
    `/investigations/${id}`,
  )
}
export function closeInvestigation(
  investigationId: string,
): Promise<Investigation> {
  return apiRequest<Investigation>(
    `/investigations/${investigationId}/close`,
    {
      method: 'PATCH',
    },
  )
}
export type CreateInvestigationPayload = {
  endpointId: string
  title: string
  description?: string
  startedAt: string
}

export function createInvestigation(
  payload: CreateInvestigationPayload,
): Promise<Investigation> {
  return apiRequest<Investigation>(
    '/investigations',
    {
      method: 'POST',
      body: JSON.stringify(payload),
    },
  )
}