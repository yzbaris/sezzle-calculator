package calc

import (
	"errors"
	"testing"
)

func TestBinaryOperations(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(a, b float64) float64
		a, b     float64
		expected float64
	}{
		{"simple addition", Add, 20, 9, 29},
		{"addition of float numbers", Add, 21.9, 8.1, 30},
		{"simple subtraction", Subtract, 1, 2, -1},
		{"simple multiplation", Multiply, -1, 2, -2},
		{"power with integer exponent", Power, 2, 10, 1024},
		{"power with negative exponent", Power, 2, -1, 0.5},
		{"power with fractional exponent", Power, 9, 0.5, 3},
		{"zero to the power of zero", Power, 0, 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.a, tt.b)

			if result != tt.expected {
				t.Errorf("got %v but expected %v", result, tt.expected)
			}
		})
	}
}

func TestUnaryOperations(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(a float64) float64
		a        float64
		expected float64
	}{
		{"percent of whole number", Percent, 50, 0.5},
		{"percent of decimal", Percent, 12.5, 0.125},
		{"percent of negative", Percent, -20, -0.2},
		{"percent of zero", Percent, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.a)

			if result != tt.expected {
				t.Errorf("got %v but expected %v", result, tt.expected)
			}
		})
	}
}

func TestUnaryOperationsWithError(t *testing.T) {
	tests := []struct {
		name          string
		fn            func(a float64) (float64, error)
		a             float64
		expected      float64
		expectedError error
	}{
		{"square root of perfect square", SquareRoot, 16, 4, nil},
		{"square root of zero", SquareRoot, 0, 0, nil},
		{"square root of fraction", SquareRoot, 0.25, 0.5, nil},
		{"square root of negative", SquareRoot, -4, 0, ErrNegativeSqrt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.fn(tt.a)

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("err = %v want = %v", err, tt.expectedError)
			}

			if result != tt.expected {
				t.Errorf("got %v but expected %v", result, tt.expected)
			}
		})
	}
}

func TestWithErrors(t *testing.T) {
	tests := []struct {
		name          string
		fn            func(a, b float64) (float64, error)
		a, b          float64
		expected      float64
		expectedError error
	}{
		{"positive div", Divide, 5, 2, 2.5, nil},
		{"negative div", Divide, 5, -2, -2.5, nil},
		{"div by zero", Divide, 5, 0, 0, ErrDivisionByZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("err = %v want = %v", err, tt.expectedError)
			}

			if result != tt.expected {
				t.Errorf("got %v but expected %v", result, tt.expected)
			}

		})
	}
}
