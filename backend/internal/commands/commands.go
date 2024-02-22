package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Command interface {
	Run(args []string) error
}

func Run(args []string) {
	switch args[0] {
	case "server":
		cmd := &ServerCommand{
			addr: "localhost",
			port: "8080",
		}
		cmd.Run(args)
	case "ingest_site_json_data":
		db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Printf("Error creating db pool %v\n", err)
			panic(1)
		}
		cmd := &IngestSiteJSONDataCommand{
			db:        db,
			filenames: args[1:],
		}
		cmd.Run(args)
	case "initdb":
		cmd := &InitDBCommand{}
		cmd.Run(args)
	default:
		fmt.Println("No command given")
	}
}
