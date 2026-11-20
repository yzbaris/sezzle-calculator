package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculateEndpoint(t *testing.T) {
	// a valid JSON body just over the 1 MB limit
	tooLarge := `{"operation":"add","operands":[1,2],"padding":"` + strings.Repeat("a", 1<<20) + `"}`

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedResult float64 // checked only when expectedStatus is 200
		expectedCode   string  // checked only for error responses
	}{
		// success
		{"add", `{"operation":"add","operands":[2,3]}`, http.StatusOK, 5, ""},
		{"divide", `{"operation":"div","operands":[10,4]}`, http.StatusOK, 2.5, ""},
		{"square root", `{"operation":"sqrt","operands":[9]}`, http.StatusOK, 3, ""},
		{"negative and decimal operands", `{"operation":"mult","operands":[-1.5,2]}`, http.StatusOK, -3, ""},

		// malformed requests
		{"not json", `not json`, http.StatusBadRequest, 0, "INVALID_JSON"},
		{"empty body", ``, http.StatusBadRequest, 0, "INVALID_JSON"},
		{"unknown field", `{"operation":"add","operands":[1,2],"x":1}`, http.StatusBadRequest, 0, "INVALID_JSON"},
		{"operand is a string", `{"operation":"add","operands":["1",2]}`, http.StatusBadRequest, 0, "INVALID_JSON"},
		{"operand is null", `{"operation":"add","operands":[null,2]}`, http.StatusBadRequest, 0, "INVALID_JSON"},
		{"body is null", `null`, http.StatusBadRequest, 0, "INVALID_JSON"},
		{"trailing text", `{"operation":"add","operands":[1,2]} trailing`, http.StatusBadRequest, 0, "INVALID_JSON"},
		{"second JSON value", `{"operation":"add","operands":[1,2]} {}`, http.StatusBadRequest, 0, "INVALID_JSON"},
		{"trailing whitespace", "{\"operation\":\"add\",\"operands\":[1,2]} \n", http.StatusOK, 3, ""},
		{"oversized trailing whitespace", `{"operation":"add","operands":[1,2]}` + strings.Repeat(" ", 1<<20), http.StatusBadRequest, 0, "INVALID_JSON"},
		{"operands is not a list", `{"operation":"add","operands":5}`, http.StatusBadRequest, 0, "INVALID_JSON"},
		{"body too large", tooLarge, http.StatusBadRequest, 0, "INVALID_JSON"},

		// valid JSON, but the calculation can't be done
		{"unknown operation", `{"operation":"mod","operands":[5,2]}`, http.StatusUnprocessableEntity, 0, "UNKNOWN_OPERATION"},
		{"missing operation", `{"operands":[5,2]}`, http.StatusUnprocessableEntity, 0, "UNKNOWN_OPERATION"},
		{"too few operands", `{"operation":"add","operands":[1]}`, http.StatusUnprocessableEntity, 0, "INVALID_OPERAND_COUNT"},
		{"missing operands", `{"operation":"add"}`, http.StatusUnprocessableEntity, 0, "INVALID_OPERAND_COUNT"},
		{"division by zero", `{"operation":"div","operands":[10,0]}`, http.StatusUnprocessableEntity, 0, "DIVISION_BY_ZERO"},
		{"square root of negative", `{"operation":"sqrt","operands":[-4]}`, http.StatusUnprocessableEntity, 0, "INVALID_RESULT"},
		{"overflow", `{"operation":"exp","operands":[10,400]}`, http.StatusUnprocessableEntity, 0, "INVALID_RESULT"},
	}

	router := NewRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("status = %d, expected %d (body: %s)", rec.Code, tt.expectedStatus, rec.Body.String())
			}
			if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("Content-Type = %q, expected application/json", contentType)
			}

			if tt.expectedStatus == http.StatusOK {
				var resp calculateResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("could not decode response: %v", err)
				}
				if resp.Result != tt.expectedResult {
					t.Errorf("result = %v, expected %v", resp.Result, tt.expectedResult)
				}
				return
			}

			var resp errorResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("could not decode error response: %v", err)
			}
			if resp.Error.Code != tt.expectedCode {
				t.Errorf("error code = %q, expected %q", resp.Error.Code, tt.expectedCode)
			}
			if resp.Error.Message == "" {
				t.Error("error message is empty")
			}
		})
	}
}

func TestCalculateEndpointRejectsOtherMethods(t *testing.T) {
	router := NewRouter()

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/calculate", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("status = %d, expected %d", rec.Code, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	NewRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, expected %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `"ok"`) {
		t.Errorf("body = %s, expected it to contain \"ok\"", rec.Body.String())
	}
}
