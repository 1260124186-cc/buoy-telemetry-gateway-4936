package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"buoy-telemetry-gateway/internal/service"
	"buoy-telemetry-gateway/internal/store"
)

func TestCanceledIngestDoesNotCommitReading(t *testing.T) {
	memory := store.NewMemory(10)
	h := New(service.New(memory))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	body := `{"buoy_id":"east-2","sensor":"salinity","value":31,"unit":"ppt","observed_at":"` + time.Now().UTC().Format(time.RFC3339Nano) + `","sequence":1}`
	req := httptest.NewRequest(http.MethodPost, "/v1/readings", strings.NewReader(body)).WithContext(ctx)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code == http.StatusCreated {
		t.Fatalf("canceled request was committed: %s", rr.Body.String())
	}
	items, err := memory.RecentReadings(context.Background(), "east-2", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("canceled request left %d readings", len(items))
	}
}
