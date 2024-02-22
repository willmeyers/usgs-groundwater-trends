package commands

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

type InitDBCommand struct {
	dsn string
}

func (cmd *InitDBCommand) Run(args []string) {
	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Error creating db pool %v\n", err)
		panic(1)
	}

	schema := `
	CREATE TABLE sites (
		id SERIAL PRIMARY KEY,
		site_name TEXT NOT NULL,
		site_no TEXT NOT NULL UNIQUE,
		location GEOGRAPHY,
		ts_estimate_in_mm NUMERIC,
		mk_z NUMERIC
	);

	CREATE TABLE datapoints (
		id SERIAL PRIMARY KEY,
		site_id INTEGER NOT NULL,
		ts TIMESTAMP NOT NULL,
		value NUMERIC NOT NULL,
		CONSTRAINT site FOREIGN KEY (site_id) REFERENCES sites(id)
	);
	`

	tag, err := db.Exec(context.Background(), schema)
	if err != nil {
		fmt.Println("Failed", err)
	}

	fmt.Println(tag)
}
