package main

import (
	"os"
	"usgs_tracker/internal/commands"
)

func main() {
	commands.Run(os.Args[1:])
}
