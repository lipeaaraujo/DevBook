import { createBrowserRouter, Navigate } from 'react-router-dom'
import { hasValidSession } from './api'
import Shell from './components/Shell'
import Feed from './pages/Feed'
import Login from './pages/Login'
import Profile from './pages/Profile'
import Register from './pages/Register'
import Settings from './pages/Settings'

function RequireAuth() {
  return hasValidSession() ? <Shell /> : <Navigate to="/login" replace />
}

export const router = createBrowserRouter([
  { path: '/login', element: <Login /> },
  { path: '/register', element: <Register /> },
  {
    element: <RequireAuth />,
    children: [
      { path: '/', element: <Feed /> },
      { path: '/users/:id', element: <Profile /> },
      { path: '/settings', element: <Settings /> },
    ],
  },
  { path: '*', element: <Navigate to="/" replace /> },
])
