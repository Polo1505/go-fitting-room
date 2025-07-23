package main

import (
	"net/http"
	"os"

	"github.com/Polo1505/go-fitting-room/internal/config"
	"github.com/Polo1505/go-fitting-room/internal/http-server/handlers"
	"github.com/Polo1505/go-fitting-room/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// "github.com/Polo1505/go-fitting-room/internal/storage/sqlite"

	"log/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting fitting-room", slog.String("env", cfg.Env))
	log.Debug("Debug is ENABLED")

	//init storage db
	storage, err := postgresql.New("127.127.126.49", "5432", "postgres", "your_password", "postgres")
	if err != nil {
		log.Error("Failed to initialize storage", slog.String("error", err.Error()))
		os.Exit(1)

	}
	defer storage.Close()

	// costume := &postgresql.Costume{
	// 	Title:       "Costume 1",
	// 	Description: "Description for Costume 1",
	// 	Image:       "https://example.com/costume1.jpg",
	// }
	// errS := storage.CreateCostume(costume)
	// if errS != nil {
	// 	log.Error("Failed to create costume", slog.String("error", errS.Error()))
	// 	os.Exit(1)
	// }
	// log.Info("Costume created successfully", slog.String("costume_name", "Costume 1"))

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	// router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.StripSlashes)

	h := handlers.New(storage, log)

	// handlers for costumes
	router.Route("/costumes", func(r chi.Router) {
		r.Post("/", h.CreateCostume)
		r.Get("/", h.GetAllCostumes)
		r.Get("/{id}", h.GetCostume)
		r.Put("/{id}", h.UpdateCostume)
		r.Delete("/{id}", h.DeleteCostume)
	})

	// create and run server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("Server started", slog.String("address", cfg.Address))
	defer srv.Close()

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Error("Server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
