package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidReading     = errors.New("invalid reading")
	ErrInvalidCalibration = errors.New("invalid calibration")
	ErrNotFound           = errors.New("not found")
)

func (r Reading) Validate(now time.Time) error {
	if strings.TrimSpace(r.BuoyID) == "" || strings.TrimSpace(r.Sensor) == "" || strings.TrimSpace(r.Unit) == "" {
		return fmt.Errorf("invalid reading: buoy_id, sensor and unit are required")
	}
	if r.ObservedAt.IsZero() || r.ObservedAt.After(now.Add(2*time.Minute)) {
		return fmt.Errorf("invalid reading: observed_at is missing or too far in the future")
	}
	if r.Sequence == 0 {
		return fmt.Errorf("invalid reading: sequence must be positive")
	}
	return nil
}

func (c Calibration) Validate() error {
	if strings.TrimSpace(c.BuoyID) == "" || strings.TrimSpace(c.Sensor) == "" {
		return fmt.Errorf("invalid calibration: buoy_id and sensor are required")
	}
	if c.Scale == 0 {
		return fmt.Errorf("invalid calibration: scale must not be zero")
	}
	if c.Min >= c.Max {
		return fmt.Errorf("invalid calibration: min must be lower than max")
	}
	return nil
}

func (c Calibration) Apply(value float64) float64 { return value*c.Scale + c.Offset }
