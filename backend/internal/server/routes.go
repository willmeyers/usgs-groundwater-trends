package server

import (
	"net/http"
	"usgs_tracker/internal/handlers"
)

func (s *Server) registerRoutes() {
	handler := handlers.ServerHandler{
		DB: s.DB,
	}

	s.Handle("/health", http.HandlerFunc(handler.HealthCheckHandler))
	s.Handle("/sites", http.HandlerFunc(handler.SitesHandler))
	s.Handle("/datapoints", http.HandlerFunc(handler.DatapointsHandler))

	for route, handler := range s.routes {
		s.Mux.Handle(route, handler)
	}
}
