import type { BinaryOperation } from '../hooks/useCalculator'

export const OPERATOR_SYMBOLS: Record<BinaryOperation, string> = {
  add: '+',
  sub: '−',
  mult: '×',
  div: '÷',
  exp: '^',
}
