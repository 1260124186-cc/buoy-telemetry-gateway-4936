package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"buoy-telemetry-gateway/internal/model"
	"buoy-telemetry-gateway/internal/service"
)

type Handler struct{ service *service.Service }

func New(s *service.Service) http.Handler {
	h := &Handler{service: s}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("POST /v1/calibrations", h.calibration)
	mux.HandleFunc("POST /v1/readings", h.reading)
	mux.HandleFunc("GET /v1/buoys/{id}/readings", h.recent)
	mux.HandleFunc("GET /v1/buoys/{id}/summary", h.summary)
	mux.HandleFunc("GET /v1/events", h.events)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
func (h *Handler) calibration(w http.ResponseWriter, r *http.Request) {
	var c model.Calibration
	if !decode(w, r, &c) {
		return
	}
	if err := h.service.SetCalibration(r.Context(), c); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, c)
}
func (h *Handler) reading(w http.ResponseWriter, r *http.Request) {
	var in model.Reading
	if !decode(w, r, &in) {
		return
	}
	out, err := h.service.Ingest(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}
func (h *Handler) recent(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	out, err := h.service.Recent(r.Context(), r.PathValue("id"), limit)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
func (h *Handler) summary(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.Summary(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
func (h *Handler) events(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.Events(r.Context(), r.URL.Query().Get("buoy_id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON: " + err.Error()})
		return false
	}
	return true
}
func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, model.ErrInvalidReading) || errors.Is(err, model.ErrInvalidCalibration) {
		status = http.StatusBadRequest
	}
	writeJSON(w, status, map[string]string{"error": strings.TrimSpace(err.Error())})
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
