// Package calc implements the arithmetic operations used by the calculator API.
package calc

import (
	"errors"
	"math"
)

// ErrDivisionByZero is returned when the divisor is zero.
var ErrDivisionByZero = errors.New("division by zero")

// ErrNegativeSqrt is returned for the square root of a negative number.
var ErrNegativeSqrt = errors.New("square root of a negative number")

// Add returns a + b.
func Add(a, b float64) float64 {
	return a + b
}

// Subtract returns a - b.
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns a * b.
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns a / b, or ErrDivisionByZero if b is 0.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

// Power returns a to the power of b.
func Power(a, b float64) float64 {
	return math.Pow(a, b)
}

// SquareRoot returns the square root of a, or ErrNegativeSqrt if a is negative.
func SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSqrt
	}
	return math.Sqrt(a), nil
}

// Percent converts a percentage to a fraction: Percent(50) is 0.5.
func Percent(a float64) float64 {
	return a / 100
}
