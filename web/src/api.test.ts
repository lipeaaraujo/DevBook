import { beforeEach, expect, test, vi } from 'vitest'
import { api, getToken, setToken } from './api'

const store = new Map<string, string>()
vi.stubGlobal('localStorage', {
  getItem: (k: string) => store.get(k) ?? null,
  setItem: (k: string, v: string) => void store.set(k, v),
  removeItem: (k: string) => void store.delete(k),
})
vi.stubGlobal('window', { location: { pathname: '/', assign: vi.fn() } })
const fetchMock = vi.fn()
vi.stubGlobal('fetch', fetchMock)

beforeEach(() => {
  store.clear()
  fetchMock.mockReset()
})

test('unwraps the error message from {message}', async () => {
  fetchMock.mockResolvedValue(
    new Response(JSON.stringify({ message: 'User with that email already exists' }), { status: 409 }),
  )
  await expect(api('/users', { method: 'POST', body: '{}' })).rejects.toThrow(
    'User with that email already exists',
  )
})

test('401 clears the stored token', async () => {
  setToken('stale-token')
  fetchMock.mockResolvedValue(new Response(JSON.stringify({ message: 'expired' }), { status: 401 }))
  await expect(api('/feed')).rejects.toThrow('expired')
  expect(getToken()).toBeNull()
})

test('a 401 from change-password does not log you out', async () => {
  setToken('good-token')
  fetchMock.mockResolvedValue(
    new Response(JSON.stringify({ message: 'Invalid credentials' }), { status: 401 }),
  )
  await expect(
    api('/users/1/change-password', { method: 'POST', body: '{}' }),
  ).rejects.toThrow('Invalid credentials')
  expect(getToken()).toBe('good-token')
})

test('204 resolves without throwing', async () => {
  fetchMock.mockResolvedValue(new Response(null, { status: 204 }))
  await expect(api('/users/1/follow', { method: 'POST' })).resolves.toBeUndefined()
})
