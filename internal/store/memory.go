package store

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"buoy-telemetry-gateway/internal/model"
)

type Memory struct {
	mu           sync.RWMutex
	maxReadings  int
	readings     map[string][]model.Reading
	calibrations map[string]model.Calibration
	events       []model.Event
}

func NewMemory(maxReadings int) *Memory {
	if maxReadings <= 0 {
		maxReadings = 100
	}
	return &Memory{maxReadings: maxReadings, readings: make(map[string][]model.Reading), calibrations: make(map[string]model.Calibration)}
}

func key(buoyID, sensor string) string { return buoyID + "\x00" + sensor }

func (m *Memory) PutCalibration(ctx context.Context, c model.Calibration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calibrations[key(c.BuoyID, c.Sensor)] = c
	return nil
}

func (m *Memory) Calibration(ctx context.Context, buoyID, sensor string) (model.Calibration, error) {
	if err := ctx.Err(); err != nil {
		return model.Calibration{}, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.calibrations[key(buoyID, sensor)]
	if !ok {
		return model.Calibration{}, fmt.Errorf("calibration unavailable: %w", model.ErrNotFound)
	}
	return c, nil
}

func (m *Memory) AppendReading(ctx context.Context, r model.Reading) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	items := append(m.readings[r.BuoyID], r)
	sort.SliceStable(items, func(i, j int) bool { return items[i].ObservedAt.Before(items[j].ObservedAt) })
	if len(items) > m.maxReadings {
		items = append([]model.Reading(nil), items[len(items)-m.maxReadings:]...)
	}
	m.readings[r.BuoyID] = items
	return nil
}

func (m *Memory) RecentReadings(ctx context.Context, buoyID string, limit int) ([]model.Reading, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := m.readings[buoyID]
	if limit <= 0 || limit > len(items) {
		limit = len(items)
	}
	out := make([]model.Reading, limit)
	for i := 0; i < limit; i++ {
		out[i] = items[len(items)-1-i]
	}
	return out, nil
}

func (m *Memory) AppendEvent(ctx context.Context, event model.Event) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
	return nil
}

func (m *Memory) Events(ctx context.Context, buoyID string) ([]model.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]model.Event, 0, len(m.events))
	for _, event := range m.events {
		if buoyID == "" || event.BuoyID == buoyID {
			out = append(out, event)
		}
	}
	return out, nil
}
