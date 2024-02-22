package server

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Server struct {
	Mux    http.ServeMux
	routes map[string]http.Handler
	DB     *pgxpool.Pool
}

func NewServer() *Server {
	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Error creating db pool %v\n", err)
		panic(1)
	}

	s := &Server{
		Mux:    *http.NewServeMux(),
		routes: make(map[string]http.Handler),
		DB:     db,
	}

	s.registerRoutes()

	return s
}

func (s *Server) Handle(route string, handler http.Handler) {
	s.routes[route] = handler
}
