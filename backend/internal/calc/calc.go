// Package calc implements the arithmetic operations used by the calculator API.
package calc

import (
	"errors"
	"math"
)

var ErrDivisionByZero = errors.New("divided by zero")
var ErrNegativeSqrt = errors.New("negative square root")

func Multiply(a, b float64) float64 {
	return a * b
}

func Add(a, b float64) float64 {
	return a + b
}

func Subtract(a, b float64) float64 {
	return a - b
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

func Power(a, b float64) float64 {
	return math.Pow(a, b)
}

func SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSqrt
	}
	return math.Sqrt(a), nil
}

func Percent(a float64) float64 {
	return a / 100
}
