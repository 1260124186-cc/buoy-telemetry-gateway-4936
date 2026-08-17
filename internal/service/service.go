package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"buoy-telemetry-gateway/internal/model"
)

type Repository interface {
	PutCalibration(context.Context, model.Calibration) error
	Calibration(context.Context, string, string) (model.Calibration, error)
	AppendReading(context.Context, model.Reading) error
	RecentReadings(context.Context, string, int) ([]model.Reading, error)
	AppendEvent(context.Context, model.Event) error
	Events(context.Context, string) ([]model.Event, error)
}

type Service struct {
	repo         Repository
	now          func() time.Time
	mu           sync.Mutex
	lastSequence map[string]uint64
}

func New(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now, lastSequence: make(map[string]uint64)}
}

func (s *Service) SetCalibration(ctx context.Context, c model.Calibration) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if err := s.repo.PutCalibration(ctx, c); err != nil {
		return fmt.Errorf("store calibration: %w", err)
	}
	return nil
}

func (s *Service) Ingest(ctx context.Context, r model.Reading) (model.Reading, error) {
	if err := r.Validate(s.now()); err != nil {
		return model.Reading{}, err
	}
	s.mu.Lock()
	last := s.lastSequence[r.BuoyID+"\x00"+r.Sensor]
	if r.Sequence <= last {
		s.mu.Unlock()
		return model.Reading{}, fmt.Errorf("%w: sequence must increase", model.ErrInvalidReading)
	}
	s.lastSequence[r.BuoyID+"\x00"+r.Sensor] = r.Sequence
	s.mu.Unlock()
	c, err := s.repo.Calibration(ctx, r.BuoyID, r.Sensor)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return model.Reading{}, fmt.Errorf("load calibration: %w", err)
	}
	if err == nil {
		r.Value = c.Apply(r.Value)
		if r.Value < c.Min || r.Value > c.Max {
			kind := "above_limit"
			if r.Value < c.Min {
				kind = "below_limit"
			}
			if err := s.repo.AppendEvent(ctx, model.Event{BuoyID: r.BuoyID, Sensor: r.Sensor, Value: r.Value, Kind: kind, CreatedAt: s.now()}); err != nil {
				return model.Reading{}, fmt.Errorf("store event: %w", err)
			}
		}
	}
	if err := s.repo.AppendReading(ctx, r); err != nil {
		return model.Reading{}, fmt.Errorf("store reading: %w", err)
	}
	return r, nil
}

func (s *Service) Recent(ctx context.Context, buoyID string, limit int) ([]model.Reading, error) {
	// 存储层已按最新到最旧返回只读副本，此处不再重排，避免改写顺序
	return s.repo.RecentReadings(ctx, buoyID, limit)
}
func (s *Service) Events(ctx context.Context, buoyID string) ([]model.Event, error) {
	return s.repo.Events(ctx, buoyID)
}

func (s *Service) Summary(ctx context.Context, buoyID string) ([]model.SensorSummary, error) {
	readings, err := s.repo.RecentReadings(ctx, buoyID, 0)
	if err != nil {
		return nil, err
	}
	bySensor := make(map[string][]model.Reading)
	for _, r := range readings {
		bySensor[r.Sensor] = append(bySensor[r.Sensor], r)
	}
	out := make([]model.SensorSummary, 0, len(bySensor))
	for sensor, items := range bySensor {
		summary := model.SensorSummary{Sensor: sensor, Unit: items[0].Unit}
		for _, item := range items {
			summary.Add(item)
		}
		summary.Mean = summary.Mean / float64(summary.Count)
		out = append(out, summary)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Sensor < out[j].Sensor })
	return out, nil
}
