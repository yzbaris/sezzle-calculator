package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewHandlerServesStaticFilesAndAPI(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<h1>Calculator</h1>"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := newHandler(dir)

	tests := []struct {
		name           string
		method, path   string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{"index page", http.MethodGet, "/", "", http.StatusOK, "Calculator"},
		{"health check", http.MethodGet, "/healthz", "", http.StatusOK, "ok"},
		{"api", http.MethodPost, "/api/v1/calculate", `{"operation":"add","operands":[2,3]}`, http.StatusOK, `"result":5`},
		{"missing file", http.MethodGet, "/nope.js", "", http.StatusNotFound, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("status = %d, expected %d", rec.Code, tt.expectedStatus)
			}
			if !strings.Contains(rec.Body.String(), tt.expectedBody) {
				t.Errorf("body = %q, expected it to contain %q", rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestNewHandlerWithoutStaticDirServesOnlyTheAPI(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	newHandler("").ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, expected %d", rec.Code, http.StatusNotFound)
	}
}
