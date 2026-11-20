// Package api exposes the calculator over HTTP
package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/yzbaris/sezzle-calculator/backend/internal/calc"
)

type calculateRequest struct {
	Operation string    `json:"operation"`
	Operands  []float64 `json:"operands"`
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

	var req calculateRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid JSON body "+err.Error())
		return
	}

	result, err := calc.Calculate(req.Operation, req.Operands)

	if err != nil {
		switch {
		case errors.Is(err, calc.ErrUnkownOperation):
			writeError(w, http.StatusUnprocessableEntity, "UNKNOWN_OPERATION", err.Error())
		case errors.Is(err, calc.ErrUnkownOperandCount):
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
