package calc

import (
	"errors"
	"fmt"
	"math"
)

// ErrUnknownOperation is returned when no operation has the given name.
var ErrUnknownOperation = errors.New("unknown operation")

// ErrInvalidOperandCount is returned when an operation gets the wrong number of operands.
var ErrInvalidOperandCount = errors.New("invalid operand count")

// ErrInvalidResult is returned when the result can't be represented (NaN or ±Inf).
var ErrInvalidResult = errors.New("result is too large or undefined")

// operation describes how many operands an operation takes and how to run it.
type operation struct {
	arity int
	fn    func(x []float64) (float64, error)
}

// operations is the registry of supported operations, by name.
var operations = map[string]operation{
	"add": {
		arity: 2,
		fn: func(x []float64) (float64, error) {
			return Add(x[0], x[1]), nil
		},
	},
	"sub": {
		arity: 2,
		fn: func(x []float64) (float64, error) {
			return Subtract(x[0], x[1]), nil
		},
	},
	"mult": {
		arity: 2,
		fn: func(x []float64) (float64, error) {
			return Multiply(x[0], x[1]), nil
		},
	},
	"div": {
		arity: 2,
		fn: func(x []float64) (float64, error) {
			return Divide(x[0], x[1])
		},
	},
	"perc": {
		arity: 1,
		fn: func(x []float64) (float64, error) {
			return Percent(x[0]), nil
		},
	},
	"exp": {
		arity: 2,
		fn: func(x []float64) (float64, error) {
			return Power(x[0], x[1]), nil
		},
	},
	"sqrt": {
		arity: 1,
		fn: func(x []float64) (float64, error) {
			return SquareRoot(x[0])
		},
	},
}

// Calculate runs the named operation on the operands.
// It returns ErrUnknownOperation, ErrInvalidOperandCount, ErrInvalidResult,
// or an error from the operation itself (for example ErrDivisionByZero).
func Calculate(opName string, operands []float64) (float64, error) {
	op, ok := operations[opName]

	if !ok {
		return 0, fmt.Errorf("%w: %q", ErrUnknownOperation, opName)
	}

	if op.arity != len(operands) {
		return 0, fmt.Errorf("%w: %s needs %d, got %d", ErrInvalidOperandCount, opName, op.arity, len(operands))
	}

	result, err := op.fn(operands)

	if err == nil && (math.IsNaN(result) || math.IsInf(result, 0)) {
		return 0, ErrInvalidResult
	}

	return result, err
}
