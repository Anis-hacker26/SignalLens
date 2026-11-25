import { apiRequest } from './api'
import type { Endpoint } from '../types/endpoint'

export function getEndpoints(): Promise<Endpoint[]> {
  return apiRequest<Endpoint[]>('/endpoints')
}

export function getEndpointById(id: string): Promise<Endpoint> {
  return apiRequest<Endpoint>(`/endpoints/${id}`)
}