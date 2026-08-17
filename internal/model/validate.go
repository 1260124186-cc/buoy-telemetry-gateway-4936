package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func CanonicalBuoyID(value string) string { return strings.ToLower(strings.TrimSpace(value)) }

var (
	ErrInvalidReading     = errors.New("invalid reading")
	ErrInvalidCalibration = errors.New("invalid calibration")
	ErrNotFound           = errors.New("not found")
)

func (r Reading) Validate(now time.Time) error {
	if strings.TrimSpace(r.BuoyID) == "" || strings.TrimSpace(r.Sensor) == "" || strings.TrimSpace(r.Unit) == "" {
		return fmt.Errorf("%w: buoy_id, sensor and unit are required", ErrInvalidReading)
	}
	if r.ObservedAt.IsZero() || r.ObservedAt.After(now.Add(2*time.Minute)) {
		return fmt.Errorf("%w: observed_at is missing or too far in the future", ErrInvalidReading)
	}
	if r.Sequence == 0 {
		return fmt.Errorf("%w: sequence must be positive", ErrInvalidReading)
	}
	return nil
}

func (c Calibration) Validate() error {
	if strings.TrimSpace(c.BuoyID) == "" || strings.TrimSpace(c.Sensor) == "" {
		return fmt.Errorf("%w: buoy_id and sensor are required", ErrInvalidCalibration)
	}
	if c.Scale == 0 {
		return fmt.Errorf("%w: scale must not be zero", ErrInvalidCalibration)
	}
	if c.Min >= c.Max {
		return fmt.Errorf("%w: min must be lower than max", ErrInvalidCalibration)
	}
	return nil
}

func (c Calibration) Apply(value float64) float64 { return value*c.Scale + c.Offset }
