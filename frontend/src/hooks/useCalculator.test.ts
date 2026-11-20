import { act, renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, test, vi } from 'vitest'
import { ApiError, calculate } from '../api/client'
import { formatResult, useCalculator } from './useCalculator'

vi.mock('../api/client', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../api/client')>()),
  calculate: vi.fn(),
}))

const calculateMock = vi.mocked(calculate)

// A tiny fake backend so the tests read like real key presses.
function fakeBackend(op: string, [a, b]: number[]): Promise<number> {
  switch (op) {
    case 'add': return Promise.resolve(a + b)
    case 'sub': return Promise.resolve(a - b)
    case 'mult': return Promise.resolve(a * b)
    case 'div':
      return b === 0 ? Promise.reject(new ApiError('DIVISION_BY_ZERO', 'division by zero')) : Promise.resolve(a / b)
    case 'exp': return Promise.resolve(a ** b)
    case 'sqrt': return Promise.resolve(Math.sqrt(a))
    case 'perc': return Promise.resolve(a / 100)
  }
  throw new Error(`unexpected op ${op}`)
}

beforeEach(() => {
  calculateMock.mockReset()
  calculateMock.mockImplementation(fakeBackend)
})

function setup() {
  const { result } = renderHook(() => useCalculator())
  const calc = () => result.current
  return {
    display: () => calc().state.display,
    state: () => calc().state,
    digits: (s: string) => {
      for (const ch of s) act(() => (ch === '.' ? calc().inputDecimal() : calc().inputDigit(ch)))
    },
    op: async (op: Parameters<ReturnType<typeof useCalculator>['chooseOperator']>[0]) => {
      await act(() => calc().chooseOperator(op))
    },
    unary: async (op: 'sqrt' | 'perc') => {
      await act(() => calc().applyUnary(op))
    },
    equals: async () => {
      await act(() => calc().equals())
    },
    sign: () => act(() => calc().toggleSign()),
    clear: () => act(() => calc().clear()),
  }
}

describe('typing numbers', () => {
  test('starts at 0', () => {
    expect(setup().display()).toBe('0')
  })

  test('replaces the leading zero', () => {
    const c = setup()
    c.digits('05')
    expect(c.display()).toBe('5')
  })

  test('appends digits', () => {
    const c = setup()
    c.digits('123')
    expect(c.display()).toBe('123')
  })

  test('allows only one decimal point', () => {
    const c = setup()
    c.digits('1..5.')
    expect(c.display()).toBe('1.5')
  })

  test('a decimal point first gives 0.', () => {
    const c = setup()
    c.digits('.5')
    expect(c.display()).toBe('0.5')
  })

  test('stops at the maximum length', () => {
    const c = setup()
    c.digits('12345678901234567890')
    expect(c.display()).toBe('1234567890123456')
  })

  test('toggles the sign', () => {
    const c = setup()
    c.digits('5')
    c.sign()
    expect(c.display()).toBe('-5')
    c.sign()
    expect(c.display()).toBe('5')
  })

  test('typing after ± on zero keeps the minus', () => {
    const c = setup()
    c.sign()
    c.digits('7')
    expect(c.display()).toBe('-7')
  })
})

describe('binary operations', () => {
  test('2 + 3 = 5', async () => {
    const c = setup()
    c.digits('2')
    await c.op('add')
    c.digits('3')
    await c.equals()
    expect(c.display()).toBe('5')
    expect(calculateMock).toHaveBeenCalledWith('add', [2, 3])
  })

  test('chains operations: 2 + 3 × shows 5, then × 4 = 20', async () => {
    const c = setup()
    c.digits('2')
    await c.op('add')
    c.digits('3')
    await c.op('mult')
    expect(c.display()).toBe('5')
    c.digits('4')
    await c.equals()
    expect(c.display()).toBe('20')
  })

  test('a second operator in a row replaces the first', async () => {
    const c = setup()
    c.digits('6')
    await c.op('add')
    await c.op('mult')
    c.digits('2')
    await c.equals()
    expect(c.display()).toBe('12')
    expect(calculateMock).toHaveBeenCalledTimes(1)
  })

  test('= without a second number does nothing', async () => {
    const c = setup()
    c.digits('5')
    await c.op('add')
    await c.equals()
    expect(c.display()).toBe('5')
    expect(c.state().error).toBeNull()
    expect(calculateMock).not.toHaveBeenCalled()
  })

  test('= without an operator does nothing', async () => {
    const c = setup()
    c.digits('5')
    await c.equals()
    expect(c.display()).toBe('5')
    expect(calculateMock).not.toHaveBeenCalled()
  })

  test('a negative second number: 5 + ± 3 = 2', async () => {
    const c = setup()
    c.digits('5')
    await c.op('add')
    c.sign()
    c.digits('3')
    await c.equals()
    expect(c.display()).toBe('2')
  })

  test('rounds float noise for display: 0.1 + 0.2 = 0.3', async () => {
    const c = setup()
    c.digits('0.1')
    await c.op('add')
    c.digits('0.2')
    await c.equals()
    expect(c.display()).toBe('0.3')
  })
})

describe('after a result', () => {
  test('a digit starts a new number', async () => {
    const c = setup()
    c.digits('10')
    await c.op('div')
    c.digits('4')
    await c.equals()
    c.digits('3')
    expect(c.display()).toBe('3')
  })

  test('an operator continues from the result', async () => {
    const c = setup()
    c.digits('10')
    await c.op('div')
    c.digits('4')
    await c.equals()
    await c.op('mult')
    c.digits('2')
    await c.equals()
    expect(c.display()).toBe('5')
  })
})

describe('unary operations', () => {
  test('√ applies right away', async () => {
    const c = setup()
    c.digits('9')
    await c.unary('sqrt')
    expect(c.display()).toBe('3')
    expect(calculateMock).toHaveBeenCalledWith('sqrt', [9])
  })

  test('% applies right away', async () => {
    const c = setup()
    c.digits('50')
    await c.unary('perc')
    expect(c.display()).toBe('0.5')
  })

  test('ignores unary operations while waiting for the second operand', async () => {
    const c = setup()
    c.digits('9')
    await c.op('add')
    await c.unary('sqrt')
    await c.equals()
    expect(c.display()).toBe('9')
    expect(calculateMock).not.toHaveBeenCalled()
  })

  test('works on the second operand: 5 + 9 √ = 8', async () => {
    const c = setup()
    c.digits('5')
    await c.op('add')
    c.digits('9')
    await c.unary('sqrt')
    await c.equals()
    expect(c.display()).toBe('8')
  })
})

describe('errors', () => {
  async function divideByZero() {
    const c = setup()
    c.digits('5')
    await c.op('div')
    c.digits('0')
    await c.equals()
    return c
  }

  test('shows the API error message', async () => {
    const c = await divideByZero()
    expect(c.state().error).toBe('division by zero')
    expect(c.state().loading).toBe(false)
  })

  test('a digit clears the error and starts fresh', async () => {
    const c = await divideByZero()
    c.digits('7')
    expect(c.state().error).toBeNull()
    expect(c.display()).toBe('7')
    expect(c.state().pendingOp).toBeNull()
  })

  test('operators are ignored while an error is shown', async () => {
    const c = await divideByZero()
    await c.op('add')
    await c.unary('sqrt')
    expect(c.state().error).toBe('division by zero')
    expect(calculateMock).toHaveBeenCalledTimes(1)
  })

  test('C resets everything', async () => {
    const c = await divideByZero()
    c.clear()
    expect(c.state()).toMatchObject({ display: '0', error: null, pendingOp: null, stored: null })
  })

  test('shows a generic message for unexpected errors', async () => {
    calculateMock.mockRejectedValueOnce(new Error('boom'))
    const c = setup()
    c.digits('9')
    await c.unary('sqrt')
    expect(c.state().error).toBe('something went wrong')
  })
})

describe('loading', () => {
  test('is true while a request runs and ignores input', async () => {
    let resolve!: (n: number) => void
    calculateMock.mockReturnValueOnce(new Promise((r) => (resolve = r)))
    const { result } = renderHook(() => useCalculator())

    act(() => result.current.inputDigit('9'))
    let pending!: Promise<void>
    act(() => {
      pending = result.current.applyUnary('sqrt')
    })
    expect(result.current.state.loading).toBe(true)

    act(() => result.current.inputDigit('1'))
    expect(result.current.state.display).toBe('9')

    await act(async () => {
      resolve(3)
      await pending
    })
    expect(result.current.state).toMatchObject({ loading: false, display: '3' })
  })
})

describe('formatResult', () => {
  test.each([
    [5, '5'],
    [2.5, '2.5'],
    [0.30000000000000004, '0.3'],
    [-0.2, '-0.2'],
    [1e21, '1e+21'],
    [1 / 3, '0.333333333333'],
  ])('%s → %s', (value, expected) => {
    expect(formatResult(value)).toBe(expected)
  })
})
