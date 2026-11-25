export type AssessmentClassification =
  | 'UNDETERMINED'
  | 'BENIGN'
  | 'SUSPICIOUS'
  | 'MALICIOUS'

export type AssessmentConfidence =
  | 'LOW'
  | 'MEDIUM'
  | 'HIGH'

export type Assessment = {
  id: string
  investigationId: string
  classification: AssessmentClassification
  confidence: AssessmentConfidence
  summary: string
  createdAt: string
  updatedAt: string
}