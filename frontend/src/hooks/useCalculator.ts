import { useState } from 'react'
import { ApiError, calculate, type Operation } from '../api/client'

export type BinaryOperation = Extract<Operation, 'add' | 'sub' | 'mult' | 'div' | 'exp'>
export type UnaryOperation = Extract<Operation, 'sqrt' | 'perc'>

export type CalculatorState = {
  display: string // what the screen shows; a string while typing ("0.", "-3")
  stored: number | null // first operand of the pending operation
  pendingOp: BinaryOperation | null
  overwrite: boolean // the next digit replaces the display instead of appending
  waitingForOperand: boolean // an operator was pressed and no second number is entered yet
  error: string | null
  loading: boolean
}

export const MAX_DIGITS = 16

export const initialState: CalculatorState = {
  display: '0',
  stored: null,
  pendingOp: null,
  overwrite: false,
  waitingForOperand: false,
  error: null,
  loading: false,
}

/** Rounds away float noise for display, e.g. 0.30000000000000004 → "0.3". */
export function formatResult(value: number): string {
  return String(Number(value.toPrecision(12)))
}

/**
 * Calculator input logic. The math itself is done by the API.
 * The keypad is disabled while `loading` is true, so only one request runs at a time.
 */
export function useCalculator() {
  const [state, setState] = useState<CalculatorState>(initialState)

  async function request(op: Operation, operands: number[], onResult: (result: number) => Partial<CalculatorState>) {
    setState((s) => ({ ...s, loading: true }))
    try {
      const result = await calculate(op, operands)
      setState((s) => ({ ...s, ...onResult(result), loading: false }))
    } catch (e) {
      const message = e instanceof ApiError ? e.message : 'something went wrong'
      setState({ ...initialState, overwrite: true, error: message })
    }
  }

  function inputDigit(digit: string) {
    if (state.loading) return
    if (state.error) {
      setState({ ...initialState, display: digit })
      return
    }
    if (state.overwrite) {
      setState({ ...state, display: digit, overwrite: false, waitingForOperand: false })
      return
    }
    if (state.display.replace(/[-.]/g, '').length >= MAX_DIGITS) return

    let display: string
    if (state.display === '0') display = digit
    else if (state.display === '-0') display = '-' + digit
    else display = state.display + digit
    setState({ ...state, display, waitingForOperand: false })
  }

  function inputDecimal() {
    if (state.loading) return
    if (state.error) {
      setState({ ...initialState, display: '0.' })
      return
    }
    if (state.overwrite) {
      setState({ ...state, display: '0.', overwrite: false, waitingForOperand: false })
      return
    }
    if (state.display.includes('.')) return
    setState({ ...state, display: state.display + '.' })
  }

  function toggleSign() {
    if (state.loading || state.error) return
    if (state.waitingForOperand) {
      // "5 + ±" starts typing a negative second number
      setState({ ...state, display: '-0', overwrite: false, waitingForOperand: false })
      return
    }
    const display = state.display.startsWith('-') ? state.display.slice(1) : '-' + state.display
    setState({ ...state, display })
  }

  async function chooseOperator(op: BinaryOperation) {
    if (state.loading || state.error) return

    if (state.pendingOp && state.waitingForOperand) {
      // two operators in a row: the last one wins
      setState({ ...state, pendingOp: op })
      return
    }

    if (state.pendingOp && state.stored !== null) {
      // chaining: "2 + 3 ×" calculates 2 + 3 first
      await request(state.pendingOp, [state.stored, Number(state.display)], (result) => ({
        display: formatResult(result),
        stored: result,
        pendingOp: op,
        overwrite: true,
        waitingForOperand: true,
      }))
      return
    }

    setState({ ...state, stored: Number(state.display), pendingOp: op, overwrite: true, waitingForOperand: true })
  }

  async function equals() {
    if (state.loading || state.error) return
    if (!state.pendingOp || state.stored === null || state.waitingForOperand) return // "5 + =" does nothing

    await request(state.pendingOp, [state.stored, Number(state.display)], (result) => ({
      display: formatResult(result),
      stored: null,
      pendingOp: null,
      overwrite: true,
      waitingForOperand: false,
    }))
  }

  async function applyUnary(op: UnaryOperation) {
    if (state.loading || state.error) return

    await request(op, [Number(state.display)], (result) => ({
      display: formatResult(result),
      overwrite: true,
      waitingForOperand: false,
    }))
  }

  function clear() {
    if (state.loading) return
    setState(initialState)
  }

  return { state, inputDigit, inputDecimal, toggleSign, chooseOperator, equals, applyUnary, clear }
}
