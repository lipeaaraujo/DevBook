const TOKEN_KEY = 'devbook_token'

export const getToken = () => localStorage.getItem(TOKEN_KEY)
export const setToken = (token: string) => localStorage.setItem(TOKEN_KEY, token)
export const clearToken = () => localStorage.removeItem(TOKEN_KEY)

export function decodeToken(token: string): { userId: string; exp: number } | null {
  try {
    const { userId, exp } = JSON.parse(atob(token.split('.')[1]))
    if (userId == null || typeof exp !== 'number') return null
    return { userId: String(userId), exp }
  } catch {
    return null
  }
}

export function currentUserId(): string | null {
  const token = getToken()
  return token ? (decodeToken(token)?.userId ?? null) : null
}

export function hasValidSession(): boolean {
  const token = getToken()
  const claims = token ? decodeToken(token) : null
  return claims !== null && claims.exp * 1000 > Date.now()
}

export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers)
  const token = getToken()
  if (token) headers.set('Authorization', `Bearer ${token}`)
  if (init.body) headers.set('Content-Type', 'application/json')
  const res = await fetch(`/api${path}`, { ...init, headers })
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
