package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RajatPrak/students/internal/config"
	"github.com/RajatPrak/students/internal/http/handlers/student"
)

func main() {

	// Load config
	cfg := config.MustLoad()
	// Database setup
	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /api/students", student.New())
	// Setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server Started", slog.String("address", cfg.Addr))
	// fmt.Printf("Server Started %s", cfg.HTTPServer.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}
