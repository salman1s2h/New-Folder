package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/salman1s2h/students-api/internal"
)

func main() {
	// fmt.Println("Hello, Students API!")
	slog.Info("Hello, Student API")
	cfg := config.MustLoad()

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Write([]byte("You called GET /"))
		case http.MethodPost:
			w.Write([]byte("You called POST /"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := cfg.HTTP.Host + ":" + cfg.HTTP.Port

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in a goroutine
	go func() {
		slog.Info("Starting server...", slog.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", slog.String("error", err.Error()))
		}
	}()

	// fmt.Println("Server started on", addr)

	slog.Info("Server Started at", slog.String("Address  :", addr))
	<-done
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.HTTP.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown Failed", slog.String("error", err.Error()))
	} else {
		slog.Info("Server exited properly")
	}
}
