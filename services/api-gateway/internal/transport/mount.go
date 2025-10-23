package transport

func (s *server) mountHandlers() {
	// mount health probes
	healthz := s.mux.Group("/healthz")
	healthz.GET("/ready", s.readinessProbe)
	healthz.GET("/live", s.livenessProbe)

	// mount trip handlers
	trip := s.mux.Group("/trip")
	trip.POST("/preview", s.previewTrip)
	trip.POST("/start", s.startTrip)

	// mount driver and rider websocket handlers
	ws := s.mux.Group("/ws")
	ws.GET("/rider", s.handleRiderWS)
	ws.GET("/driver", s.handleDriverWS)
}
