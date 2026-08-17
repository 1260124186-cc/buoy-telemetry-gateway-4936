package service

import (
	"buoy-telemetry-gateway/internal/model"
	"buoy-telemetry-gateway/internal/store"
	"context"
	"testing"
	"time"
)

func TestBuoyIdentityUsesOneCanonicalKey(t *testing.T) {
	s := New(store.NewMemory(10))
	now := time.Now()
	s.now = func() time.Time { return now }
	if err := s.SetCalibration(context.Background(), model.Calibration{BuoyID: "north-9", Sensor: "temperature", Scale: 2, Offset: 1, Min: -10, Max: 50}); err != nil {
		t.Fatal(err)
	}
	r, err := s.Ingest(context.Background(), model.Reading{BuoyID: " NORTH-9 ", Sensor: "temperature", Value: 4, Unit: "C", ObservedAt: now, Sequence: 1})
	if err != nil {
		t.Fatal(err)
	}
	if r.Value != 9 {
		t.Fatalf("calibration did not match canonical buoy: %v", r.Value)
	}
	items, err := s.Recent(context.Background(), "north-9", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("history length=%d", len(items))
	}
}
