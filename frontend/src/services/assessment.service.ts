import { apiRequest } from './api'
import type {
  Assessment,
  AssessmentClassification,
  AssessmentConfidence,
} from '../types/assessment'

export type CreateAssessmentPayload = {
  classification: AssessmentClassification
  confidence: AssessmentConfidence
  reasoning?: string
}

export function createAssessment(
  investigationId: string,
  payload: CreateAssessmentPayload,
): Promise<Assessment> {
  return apiRequest<Assessment>(
    `/investigations/${investigationId}/assessments`,
    {
      method: 'POST',
      body: JSON.stringify(payload),
    },
  )
}

export function updateAssessment(
  investigationId: string,
  payload: CreateAssessmentPayload,
): Promise<Assessment> {
  return apiRequest<Assessment>(
    `/investigations/${investigationId}/assessment`,
    {
      method: 'PATCH',
      body: JSON.stringify(payload),
    },
  )
}