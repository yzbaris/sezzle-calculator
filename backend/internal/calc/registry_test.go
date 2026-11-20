package calc

import (
	"errors"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name          string
		op            string
		operands      []float64
		expected      float64
		expectedError error
	}{
		// one success case per operation
		{"add", "add", []float64{2, 3}, 5, nil},
		{"subtract", "sub", []float64{5, 8}, -3, nil},
		{"multiply", "mult", []float64{2.5, 4}, 10, nil},
		{"positive div", "div", []float64{5, 2}, 2.5, nil},
		{"negative div", "div", []float64{5, -2}, -2.5, nil},
		{"percent", "perc", []float64{50}, 0.5, nil},
		{"power", "exp", []float64{2, 3}, 8, nil},
		{"square root", "sqrt", []float64{9}, 3, nil},

		// unknown operation
		{"unknown operation", "mod", []float64{5, 2}, 0, ErrUnknownOperation},
		{"empty operation name", "", []float64{5, 2}, 0, ErrUnknownOperation},
		{"operation name is case sensitive", "ADD", []float64{5, 2}, 0, ErrUnknownOperation},

		// wrong number of operands
		{"add with one operand", "add", []float64{5}, 0, ErrInvalidOperandCount},
		{"add with three operands", "add", []float64{1, 2, 3}, 0, ErrInvalidOperandCount},
		{"add with no operands", "add", nil, 0, ErrInvalidOperandCount},
		{"sqrt with two operands", "sqrt", []float64{9, 4}, 0, ErrInvalidOperandCount},

		// errors from the operation itself are passed through
		{"div by zero", "div", []float64{5, 0}, 0, ErrDivisionByZero},
		{"zero div by zero", "div", []float64{0, 0}, 0, ErrDivisionByZero},
		{"sqrt of negative", "sqrt", []float64{-1}, 0, ErrNegativeSqrt},

		// results that float64 can't represent
		{"power overflows to infinity", "exp", []float64{10, 400}, 0, ErrInvalidResult},
		{"fractional power of negative is NaN", "exp", []float64{-8, 1.0 / 3}, 0, ErrInvalidResult},
		{"multiply overflows to infinity", "mult", []float64{1e308, 10}, 0, ErrInvalidResult},
		{"add overflows to negative infinity", "add", []float64{-1e308, -1e308}, 0, ErrInvalidResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calculate(tt.op, tt.operands)

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("err = %v want = %v", err, tt.expectedError)
			}

			if result != tt.expected {
				t.Errorf("got %v but expected %v", result, tt.expected)
			}

		})
	}
}
