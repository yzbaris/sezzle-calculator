package calc

import (
	"errors"
	"math"
)

var ErrUnkownOperation = errors.New("unkown calculator operation")
var ErrUnkownOperandCount = errors.New("unkown operand count")
var ErrInvalidResult = errors.New("result is invalid")

type CalcOperation struct {
	NumOfOperand int
	Fn           func(x []float64) (float64, error)
}

var operations = map[string]CalcOperation{
	"add": {
		NumOfOperand: 2,
		Fn: func(x []float64) (float64, error) {
			return Add(x[0], x[1]), nil
		},
	},
	"sub": {
		NumOfOperand: 2,
		Fn: func(x []float64) (float64, error) {
			return Subtract(x[0], x[1]), nil
		},
	},
	"mult": {
		NumOfOperand: 2,
		Fn: func(x []float64) (float64, error) {
			return Multiply(x[0], x[1]), nil
		},
	},
	"div": {
		NumOfOperand: 2,
		Fn: func(x []float64) (float64, error) {
			return Divide(x[0], x[1])
		},
	},
	"perc": {
		NumOfOperand: 1,
		Fn: func(x []float64) (float64, error) {
			return Percent(x[0]), nil
		},
	},
	"exp": {
		NumOfOperand: 2,
		Fn: func(x []float64) (float64, error) {
			return Power(x[0], x[1]), nil
		},
	},
	"sqrt": {
		NumOfOperand: 1,
		Fn: func(x []float64) (float64, error) {
			return SquareRoot(x[0])
		},
	},
}

func Calculate(opName string, operands []float64) (float64, error) {
	operation, ok := operations[opName]

	if !ok {
		return 0, ErrUnkownOperation
	}

	if operation.NumOfOperand != len(operands) {
		return 0, ErrUnkownOperandCount
	}

	result, err := operation.Fn(operands)

	if err == nil && (math.IsNaN(result) || math.IsInf(result, 0)) {
		return 0, ErrInvalidResult
	}

	return result, err
}
