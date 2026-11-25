import { apiRequest } from './api'
import type {
  Evidence,
  EvidenceDetail,
} from '../types/evidence'

export function getEvidence(): Promise<Evidence[]> {
  return apiRequest<Evidence[]>('/evidence')
}

export function getEvidenceById(
  id: string,
): Promise<EvidenceDetail> {
  return apiRequest<EvidenceDetail>(
    `/evidence/${id}`,
  )
}