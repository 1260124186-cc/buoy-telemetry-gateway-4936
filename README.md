# Buoy Telemetry Gateway

Buoy Telemetry Gateway is a small Go HTTP service for coastal engineering teams that receive measurements from autonomous buoys. It validates incoming readings, applies per-sensor calibration, keeps a bounded recent history, and emits anomaly events when calibrated values leave configured limits.

## Structure

- `cmd/server`: HTTP server entry point.
- `internal/model`: domain types and validation rules.
- `internal/store`: concurrency-safe in-memory repositories.
- `internal/service`: ingestion, querying, calibration, and event workflows.
- `internal/httpapi`: JSON HTTP handlers and routing.

## API

- `POST /v1/calibrations` stores a calibration profile.
- `POST /v1/readings` ingests and calibrates a reading.
- `GET /v1/buoys/{id}/readings?limit=20` returns recent readings newest first.
- `GET /v1/buoys/{id}/summary` returns per-sensor aggregates.
- `GET /v1/events?buoy_id=...` lists anomaly events.
- `GET /healthz` reports service health.

## Run

```bash
go run ./cmd/server
```

The server listens on `:8080` by default. Set `HTTP_ADDR` to override it. The service has no external dependencies or required online services.

## Build and test

```bash
go build ./...
go test ./...
```
