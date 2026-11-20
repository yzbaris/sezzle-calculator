import { useEffect, useRef } from 'react'
import { useCalculator, type BinaryOperation } from '../hooks/useCalculator'
import { Display } from './Display'
import { Keypad } from './Keypad'

const KEY_TO_OPERATOR: Record<string, BinaryOperation> = {
  '+': 'add',
  '-': 'sub',
  '*': 'mult',
  x: 'mult',
  '/': 'div',
  '^': 'exp',
}

export function Calculator() {
  const calc = useCalculator()

  // Keep the latest calculator in a ref, so the keyboard listener is added only once.
  const calcRef = useRef(calc)
  useEffect(() => {
    calcRef.current = calc
  })

  useEffect(() => {
    function onKeyDown(event: KeyboardEvent) {
      if (event.ctrlKey || event.metaKey || event.altKey) return
      const c = calcRef.current
      const { key } = event

      if (/^[0-9]$/.test(key)) c.inputDigit(key)
      else if (key === '.' || key === ',') c.inputDecimal()
      else if (key in KEY_TO_OPERATOR) c.chooseOperator(KEY_TO_OPERATOR[key])
      else if (key === 'Enter' || key === '=') c.equals()
      else if (key === '%') c.applyUnary('perc')
      else if (key === 'r') c.applyUnary('sqrt')
      else if (key === 'Escape' || key === 'Delete' || key === 'c') c.clear()
      else return

      event.preventDefault() // e.g. stop Enter from also clicking the focused button
    }

    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [])

  const { state } = calc

  return (
    <section className="calculator" aria-label="Calculator">
      <Display
        display={state.display}
        stored={state.stored}
        pendingOp={state.pendingOp}
        error={state.error}
        loading={state.loading}
      />
      <Keypad
        disabled={state.loading}
        onDigit={calc.inputDigit}
        onDecimal={calc.inputDecimal}
        onOperator={calc.chooseOperator}
        onUnary={calc.applyUnary}
        onEquals={calc.equals}
        onToggleSign={calc.toggleSign}
        onClear={calc.clear}
      />
      <p className="calculator__hint">Keyboard: 0–9 . + − * / ^ % r (√) Enter Esc</p>
    </section>
  )
}
