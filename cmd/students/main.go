package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RajatPrak/students/internal/config"
)

func main() {

	// Load config
	cfg := config.MustLoad()
	// Database setup
	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api"))
	})
	// Setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Printf("Server Started %s", cfg.HTTPServer.Addr)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Failed to start server")
	}

}
