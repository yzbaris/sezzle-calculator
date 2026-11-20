import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, test, vi } from 'vitest'
import { ApiError, calculate } from '../api/client'
import { Calculator } from './Calculator'

vi.mock('../api/client', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../api/client')>()),
  calculate: vi.fn(),
}))

const calculateMock = vi.mocked(calculate)

beforeEach(() => {
  calculateMock.mockReset()
})

const display = () => screen.getByRole('status')
const press = (user: ReturnType<typeof userEvent.setup>, name: string) =>
  user.click(screen.getByRole('button', { name }))

describe('Calculator', () => {
  test('calculates with button clicks: 2 + 3 = 5', async () => {
    calculateMock.mockResolvedValueOnce(5)
    const user = userEvent.setup()
    render(<Calculator />)

    await press(user, '2')
    await press(user, 'add')
    expect(screen.getByTestId('expression')).toHaveTextContent('2 +')
    await press(user, '3')
    await press(user, 'equals')

    expect(display()).toHaveTextContent('5')
    expect(calculateMock).toHaveBeenCalledWith('add', [2, 3])
    expect(screen.getByTestId('expression')).toHaveTextContent('')
  })

  test('calculates with the keyboard: 10 / 4 Enter = 2.5', async () => {
    calculateMock.mockResolvedValueOnce(2.5)
    const user = userEvent.setup()
    render(<Calculator />)

    await user.keyboard('10/4{Enter}')

    expect(display()).toHaveTextContent('2.5')
    expect(calculateMock).toHaveBeenCalledWith('div', [10, 4])
  })

  test('Escape clears the display', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.keyboard('123{Escape}')

    expect(display()).toHaveTextContent('0')
  })

  test('shows the API error message', async () => {
    calculateMock.mockRejectedValueOnce(new ApiError('DIVISION_BY_ZERO', 'divided by zero'))
    const user = userEvent.setup()
    render(<Calculator />)

    await press(user, '5')
    await press(user, 'divide')
    await press(user, '0')
    await press(user, 'equals')

    expect(display()).toHaveTextContent('divided by zero')
  })

  test('applies square root right away', async () => {
    calculateMock.mockResolvedValueOnce(3)
    const user = userEvent.setup()
    render(<Calculator />)

    await press(user, '9')
    await press(user, 'square root')

    expect(display()).toHaveTextContent('3')
    expect(calculateMock).toHaveBeenCalledWith('sqrt', [9])
  })

  test('disables the keypad while a request is running', async () => {
    let resolve!: (n: number) => void
    calculateMock.mockReturnValueOnce(new Promise((r) => (resolve = r)))
    const user = userEvent.setup()
    render(<Calculator />)

    await press(user, '9')
    await press(user, 'square root')

    expect(screen.getByRole('button', { name: 'equals' })).toBeDisabled()
    expect(screen.getByRole('button', { name: '7' })).toBeDisabled()

    resolve(3)
    expect(await screen.findByText('3', { selector: 'output' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'equals' })).toBeEnabled()
  })

  test('every key has an accessible name', () => {
    render(<Calculator />)
    for (const button of screen.getAllByRole('button')) {
      expect(button).toHaveAccessibleName()
    }
  })
})
