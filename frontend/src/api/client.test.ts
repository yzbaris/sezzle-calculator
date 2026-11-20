import { afterEach, describe, expect, test, vi } from 'vitest'
import { ApiError, calculate } from './client'

function mockFetch(status: number, body: unknown) {
  const fetchMock = vi.fn().mockResolvedValue(
    new Response(typeof body === 'string' ? body : JSON.stringify(body), { status }),
  )
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('calculate', () => {
  test('sends the operation and operands, and returns the result', async () => {
    const fetchMock = mockFetch(200, { result: 2.5 })

    await expect(calculate('div', [10, 4])).resolves.toBe(2.5)

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ operation: 'div', operands: [10, 4] }),
    })
  })

  test('throws an ApiError with the code and message from the API', async () => {
    mockFetch(422, { error: { code: 'DIVISION_BY_ZERO', message: 'division by zero' } })

    const error = await calculate('div', [10, 0]).catch((e: unknown) => e)

    expect(error).toBeInstanceOf(ApiError)
    expect(error).toMatchObject({ code: 'DIVISION_BY_ZERO', message: 'division by zero' })
  })

  test('throws NETWORK_ERROR when the server is unreachable', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')))

    await expect(calculate('add', [1, 2])).rejects.toMatchObject({ code: 'NETWORK_ERROR' })
  })

  test('throws INVALID_RESPONSE when the body is not JSON', async () => {
    mockFetch(502, '<html>Bad Gateway</html>')

    await expect(calculate('add', [1, 2])).rejects.toMatchObject({ code: 'INVALID_RESPONSE' })
  })

  test('falls back to a generic error when the error body has an unexpected shape', async () => {
    mockFetch(500, { message: 'something else' })

    await expect(calculate('add', [1, 2])).rejects.toMatchObject({ code: 'UNKNOWN_ERROR' })
  })
})
