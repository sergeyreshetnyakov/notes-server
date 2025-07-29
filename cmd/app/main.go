package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sergeyreshetnyakov/notion/internal/bussines/notes"
	"github.com/sergeyreshetnyakov/notion/internal/config"
	notehandler "github.com/sergeyreshetnyakov/notion/internal/handlers/note"
	"github.com/sergeyreshetnyakov/notion/internal/lib/logger"
	"github.com/sergeyreshetnyakov/notion/internal/lib/logger/sl"
	"github.com/sergeyreshetnyakov/notion/internal/middlewares"
	notestorage "github.com/sergeyreshetnyakov/notion/internal/storage/notes"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	mux := http.NewServeMux()

	storage, shutdownDB := notestorage.New(cfg.StoragePath, log)
	notehandler.New(log, notes.New(storage)).HandleRoutes(mux)

	wrappedMux := middlewares.LoggingMiddleware(mux, log)
	server := http.Server{
		Addr:           cfg.Port,
		Handler:        wrappedMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Info("Server is running on http://localhost" + cfg.Port)
		if err := server.ListenAndServe(); !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
		log.Info("Serving new connections is stopped")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP shutdown error: %v", sl.Err(err))
	}
	shutdownDB()
	log.Info("Graceful shutdown complete")
}
