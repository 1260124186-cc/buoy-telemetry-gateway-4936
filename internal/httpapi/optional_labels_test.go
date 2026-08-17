package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"buoy-telemetry-gateway/internal/service"
	"buoy-telemetry-gateway/internal/store"
)

func TestReadingWithoutOptionalLabelsIsAccepted(t *testing.T) {
	h := New(service.New(store.NewMemory(10)))
	body := `{"buoy_id":"south-3","sensor":"oxygen","value":8,"unit":"mg/L","observed_at":"` + time.Now().UTC().Format(time.RFC3339Nano) + `","sequence":1}`
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/v1/readings", strings.NewReader(body)))
	if rr.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}
