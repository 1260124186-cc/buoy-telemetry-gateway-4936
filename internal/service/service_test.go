package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"buoy-telemetry-gateway/internal/model"
	"buoy-telemetry-gateway/internal/store"
)

func TestIngestCalibrationOrderingAndSummary(t *testing.T) {
	s := New(store.NewMemory(10))
	s.now = func() time.Time { return time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC) }
	ctx := context.Background()
	c := model.Calibration{BuoyID: "north-1", Sensor: "temperature", Scale: 2, Offset: 1, Min: -5, Max: 30}
	if err := s.SetCalibration(ctx, c); err != nil {
		t.Fatal(err)
	}
	for _, r := range []model.Reading{{BuoyID: "north-1", Sensor: "temperature", Value: 20, Unit: "C", ObservedAt: s.now().Add(-time.Minute), Sequence: 1}, {BuoyID: "north-1", Sensor: "temperature", Value: 4, Unit: "C", ObservedAt: s.now(), Sequence: 2}} {
		if _, err := s.Ingest(ctx, r); err != nil {
			t.Fatal(err)
		}
	}
	recent, err := s.Recent(ctx, "north-1", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(recent) != 2 || recent[0].Sequence != 2 || recent[0].Value != 9 {
		t.Fatalf("unexpected recent: %#v", recent)
	}
	summaries, err := s.Summary(ctx, "north-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 || summaries[0].Mean != 25 {
		t.Fatalf("unexpected summary: %#v", summaries)
	}
	events, err := s.Events(ctx, "north-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Kind != "above_limit" {
		t.Fatalf("unexpected events: %#v", events)
	}
}

func TestValidationAndSequence(t *testing.T) {
	s := New(store.NewMemory(2))
	now := time.Now()
	s.now = func() time.Time { return now }
	r := model.Reading{BuoyID: "b", Sensor: "wave", Unit: "m", ObservedAt: now, Sequence: 1}
	if _, err := s.Ingest(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Ingest(context.Background(), r); !errors.Is(err, model.ErrInvalidReading) {
		t.Fatalf("expected invalid reading, got %v", err)
	}
}
