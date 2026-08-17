package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"buoy-telemetry-gateway/internal/service"
	"buoy-telemetry-gateway/internal/store"
)

func TestHealthAndInvalidJSON(t *testing.T) {
	h := New(service.New(store.NewMemory(10)))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("health status %d", rr.Code)
	}
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/v1/readings", bytes.NewBufferString("{")))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("invalid json status %d", rr.Code)
	}
}
