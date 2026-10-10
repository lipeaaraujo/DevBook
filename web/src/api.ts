const TOKEN_KEY = 'devbook_token'
const API_URL = import.meta.env.VITE_API_URL

export const getToken = () => localStorage.getItem(TOKEN_KEY)
export const setToken = (token: string) => localStorage.setItem(TOKEN_KEY, token)
export const clearToken = () => localStorage.removeItem(TOKEN_KEY)

// The token is opaque: an expired one is caught by the 401 handler in api().
export const hasSession = () => getToken() !== null

export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers)
  const token = getToken()
  if (token) headers.set('Authorization', `Bearer ${token}`)
  if (init.body) headers.set('Content-Type', 'application/json')
  const res = await fetch(`${API_URL}${path}`, { ...init, headers })
  if (!res.ok) {
    const checksCredentials = path === '/login' || path.endsWith('/change-password')
    if (res.status === 401 && !checksCredentials) {
      clearToken()
      window.location.assign('/login')
    }
    const message = await res
      .json()
      .then((body) => body.message)
      .catch(() => res.statusText)
    throw new Error(message)
  }
  const text = await res.text()
  return (text ? JSON.parse(text) : undefined) as T
}
