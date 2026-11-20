import type { BinaryOperation } from '../hooks/useCalculator'
import { formatResult } from '../hooks/useCalculator'
import { OPERATOR_SYMBOLS } from './operators'

type Props = {
  display: string
  stored: number | null
  pendingOp: BinaryOperation | null
  error: string | null
  loading: boolean
}

export function Display({ display, stored, pendingOp, error, loading }: Props) {
  const expression = pendingOp && stored !== null ? `${formatResult(stored)} ${OPERATOR_SYMBOLS[pendingOp]}` : ''

  return (
    <div className="display" aria-busy={loading}>
      <div className="display__expression" data-testid="expression">
        {expression}
      </div>
      <output className={error ? 'display__value display__value--error' : 'display__value'} aria-live="polite">
        {error ?? display}
      </output>
    </div>
  )
}
