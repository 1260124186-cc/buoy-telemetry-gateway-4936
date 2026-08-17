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
