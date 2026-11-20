import type { BinaryOperation, UnaryOperation } from '../hooks/useCalculator'
import { OPERATOR_SYMBOLS } from './operators'

type Props = {
  disabled: boolean
  onDigit: (digit: string) => void
  onDecimal: () => void
  onOperator: (op: BinaryOperation) => void
  onUnary: (op: UnaryOperation) => void
  onEquals: () => void
  onToggleSign: () => void
  onClear: () => void
}

type Key = {
  label: string
  ariaLabel?: string
  kind: 'digit' | 'operator' | 'function' | 'equals'
  wide?: boolean
  onPress: () => void
}

export function Keypad(props: Props) {
  const digit = (d: string): Key => ({ label: d, kind: 'digit', onPress: () => props.onDigit(d) })
  const operator = (op: BinaryOperation, ariaLabel: string): Key => ({
    label: OPERATOR_SYMBOLS[op],
    ariaLabel,
    kind: 'operator',
    onPress: () => props.onOperator(op),
  })

  const keys: Key[] = [
    { label: 'C', ariaLabel: 'clear', kind: 'function', onPress: props.onClear },
    { label: '√', ariaLabel: 'square root', kind: 'function', onPress: () => props.onUnary('sqrt') },
    { label: '%', ariaLabel: 'percent', kind: 'function', onPress: () => props.onUnary('perc') },
    { label: 'xʸ', ariaLabel: 'power', kind: 'operator', onPress: () => props.onOperator('exp') },

    digit('7'), digit('8'), digit('9'), operator('div', 'divide'),
    digit('4'), digit('5'), digit('6'), operator('mult', 'multiply'),
    digit('1'), digit('2'), digit('3'), operator('sub', 'subtract'),

    { label: '±', ariaLabel: 'toggle sign', kind: 'function', onPress: props.onToggleSign },
    digit('0'),
    { label: '.', ariaLabel: 'decimal point', kind: 'digit', onPress: props.onDecimal },
    operator('add', 'add'),

    { label: '=', ariaLabel: 'equals', kind: 'equals', wide: true, onPress: props.onEquals },
  ]

  return (
    <div className="keypad">
      {keys.map((key) => (
        <button
          key={key.label}
          type="button"
          className={`key key--${key.kind}${key.wide ? ' key--wide' : ''}`}
          aria-label={key.ariaLabel}
          disabled={props.disabled}
          onClick={key.onPress}
        >
          {key.label}
        </button>
      ))}
    </div>
  )
}
