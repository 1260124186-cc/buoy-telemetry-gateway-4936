package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"buoy-telemetry-gateway/internal/service"
	"buoy-telemetry-gateway/internal/store"
)

func TestErrorContractPreservesClassification(t *testing.T) {
	h := New(service.New(store.NewMemory(10)))
	invalid := httptest.NewRecorder()
	h.ServeHTTP(invalid, httptest.NewRequest(http.MethodPost, "/v1/calibrations", bytes.NewBufferString(`{"buoy_id":"north-1","sensor":"temperature","scale":0,"min":0,"max":20}`)))
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("invalid calibration status = %d, body=%s", invalid.Code, invalid.Body.String())
	}

	reading := httptest.NewRecorder()
	body := `{"buoy_id":"north-1","sensor":"temperature","value":4,"unit":"C","observed_at":"` + time.Now().UTC().Format(time.RFC3339Nano) + `","sequence":1}`
	h.ServeHTTP(reading, httptest.NewRequest(http.MethodPost, "/v1/readings", strings.NewReader(body)))
	if reading.Code != http.StatusCreated {
		t.Fatalf("uncalibrated reading status = %d, body=%s", reading.Code, reading.Body.String())
	}
}
