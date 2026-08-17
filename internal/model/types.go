package model

import "time"

type Reading struct {
	BuoyID     string    `json:"buoy_id"`
	Sensor     string    `json:"sensor"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	ObservedAt time.Time `json:"observed_at"`
	Sequence   uint64    `json:"sequence"`
}

type Calibration struct {
	BuoyID string  `json:"buoy_id"`
	Sensor string  `json:"sensor"`
	Scale  float64 `json:"scale"`
	Offset float64 `json:"offset"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
}

type Event struct {
	BuoyID    string    `json:"buoy_id"`
	Sensor    string    `json:"sensor"`
	Value     float64   `json:"value"`
	Kind      string    `json:"kind"`
	CreatedAt time.Time `json:"created_at"`
}

type SensorSummary struct {
	Sensor string  `json:"sensor"`
	Unit   string  `json:"unit"`
	Count  int     `json:"count"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Mean   float64 `json:"mean"`
}

func (s *SensorSummary) Add(reading Reading) {
	// 首样本之前 Min/Max 均为零值，无法表示正数下界与负数上界，故用首样本初始化区间
	if s.Count == 0 {
		s.Min = reading.Value
		s.Max = reading.Value
	}
	s.Count++
	s.Mean += reading.Value
	if reading.Value < s.Min {
		s.Min = reading.Value
	}
	if reading.Value > s.Max {
		s.Max = reading.Value
	}
}
