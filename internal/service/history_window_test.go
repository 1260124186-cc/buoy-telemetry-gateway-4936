package service

import (
	"context"
	"testing"
	"time"

	"buoy-telemetry-gateway/internal/model"
	"buoy-telemetry-gateway/internal/store"
)

func TestHistoryWindowIsStableAndSummaryKeepsRange(t *testing.T) {
	now := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)
	s := New(store.NewMemory(10))
	s.now = func() time.Time { return now }
	for i, value := range []float64{-2, 4, 9} {
		_, err := s.Ingest(context.Background(), model.Reading{BuoyID: "west-7", Sensor: "wave", Value: value, Unit: "m", ObservedAt: now.Add(time.Duration(i) * time.Second), Sequence: uint64(i + 1)})
		if err != nil {
			t.Fatal(err)
		}
	}
	first, _ := s.Recent(context.Background(), "west-7", 3)
	second, _ := s.Recent(context.Background(), "west-7", 3)
	if len(first) != 3 || len(second) != 3 || first[0].Sequence != 3 || second[0].Sequence != 3 {
		t.Fatalf("unstable newest-first history: first=%#v second=%#v", first, second)
	}
	summary, err := s.Summary(context.Background(), "west-7")
	if err != nil {
		t.Fatal(err)
	}
	if len(summary) != 1 || summary[0].Min != -2 || summary[0].Max != 9 || summary[0].Mean != 11.0/3.0 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
}
