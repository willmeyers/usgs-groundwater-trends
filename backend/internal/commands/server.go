package commands

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"usgs_tracker/internal/server"
)

type ServerCommand struct {
	addr string
	port string
}

func (cmd *ServerCommand) Run(args []string) {
	server := server.NewServer()
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cmd.addr, cmd.port),
		Handler: &server.Mux,
	}

	go func() {
		fmt.Println("Starting server...")
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("Error %v\n", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.DB.Close()

	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutdown %v\n", err)
	}
}
