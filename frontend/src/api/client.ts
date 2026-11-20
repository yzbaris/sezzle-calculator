export type Operation = 'add' | 'sub' | 'mult' | 'div' | 'exp' | 'sqrt' | 'perc'

/** An error returned by the API ({"error": {"code", "message"}}) or a network failure. */
export class ApiError extends Error {
  readonly code: string

  constructor(code: string, message: string) {
    super(message)
    this.name = 'ApiError'
    this.code = code
  }
}

type SuccessBody = { result: number }
type ErrorBody = { error: { code: string; message: string } }

export async function calculate(operation: Operation, operands: number[]): Promise<number> {
  let response: Response
  try {
    response = await fetch('/api/v1/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ operation, operands }),
    })
  } catch {
    throw new ApiError('NETWORK_ERROR', 'could not reach the server')
  }

  let body: unknown
  try {
    body = await response.json()
  } catch {
    throw new ApiError('INVALID_RESPONSE', `unexpected response from the server (${response.status})`)
  }

  if (!response.ok) {
    const { error } = body as ErrorBody
    throw new ApiError(error?.code ?? 'UNKNOWN_ERROR', error?.message ?? `request failed (${response.status})`)
  }

  return (body as SuccessBody).result
}
