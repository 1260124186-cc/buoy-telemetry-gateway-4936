package main

import (
	"log"
	"net/http"
	"os"

	"buoy-telemetry-gateway/internal/httpapi"
	"buoy-telemetry-gateway/internal/service"
	"buoy-telemetry-gateway/internal/store"
)

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	handler := httpapi.New(service.New(store.NewMemory(500)))
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
