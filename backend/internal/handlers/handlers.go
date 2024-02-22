package handlers

import "github.com/jackc/pgx/v4/pgxpool"

type ServerHandler struct {
	DB *pgxpool.Pool
}
