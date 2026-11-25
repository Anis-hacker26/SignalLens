import { Routes, Route } from 'react-router-dom'
import AppLayout from '../components/layout/AppLayout'
import DashboardPage from '../pages/Dashboard/DashboardPage'
import EvidencePage from '../pages/Evidence/EvidencePage'
import InvestigationsPage from '../pages/Investigations/InvestigationsPage'
import InvestigationDetails from '../pages/Investigations/InvestigationDetails'
import CreateInvestigation from '../pages/Investigations/CreateInvestigation'
import AssessmentsPage from '../pages/Assessments/AssessmentsPage'
import EndpointsPage from '../pages/Endpoints/EndpointsPage'

export default function AppRoutes() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route
          path="/"
          element={<DashboardPage />}
        />

        <Route
          path="/evidence"
          element={<EvidencePage />}
        />

        <Route
          path="/investigations"
          element={<InvestigationsPage />}
        />

        <Route
          path="/investigations/new"
          element={<CreateInvestigation />}
        />

        <Route
          path="/investigations/:id"
          element={<InvestigationDetails />}
        />

        <Route
          path="/endpoints"
          element={<EndpointsPage />}
        />

        <Route
          path="/assessments"
          element={<AssessmentsPage />}
        />
      </Route>
    </Routes>
  )
}