// Package api exposes the calculator over HTTP.
package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/yzbaris/sezzle-calculator/backend/internal/calc"
)

type calculateRequest struct {
	Operation string     `json:"operation"`
	Operands  []*float64 `json:"operands"`
}

type calculateResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error errorMessage `json:"error"`
}

type errorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewRouter returns the HTTP handler with all API routes.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/calculate", handleCalculate)
	mux.HandleFunc("GET /healthz", handleHealth)
	return mux
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // limit the body to 1 MB
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req *calculateRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid JSON body: "+err.Error())
		return
	}

	if req == nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "expected a JSON object")
		return
	}
	if err := dec.Decode(new(json.RawMessage)); err != io.EOF {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "expected exactly one JSON object within the body size limit")
		return
	}

	operands := make([]float64, len(req.Operands))
	for i, operand := range req.Operands {
		if operand == nil {
			writeError(w, http.StatusBadRequest, "INVALID_JSON", "operands must be numbers")
			return
		}
		operands[i] = *operand
	}

	result, err := calc.Calculate(req.Operation, operands)

	if err != nil {
		switch {
		case errors.Is(err, calc.ErrUnknownOperation):
			writeError(w, http.StatusUnprocessableEntity, "UNKNOWN_OPERATION", err.Error())
		case errors.Is(err, calc.ErrInvalidOperandCount):
			writeError(w, http.StatusUnprocessableEntity, "INVALID_OPERAND_COUNT", err.Error())
		case errors.Is(err, calc.ErrDivisionByZero):
			writeError(w, http.StatusUnprocessableEntity, "DIVISION_BY_ZERO", err.Error())
		case errors.Is(err, calc.ErrInvalidResult) || errors.Is(err, calc.ErrNegativeSqrt):
			writeError(w, http.StatusUnprocessableEntity, "INVALID_RESULT", err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL", "unexpected error")
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calculateResponse{Result: result})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: errorMessage{Code: code, Message: message}})
}
