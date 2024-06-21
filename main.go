package main

import (
	"github.com/dyammarcano/sorteador/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/inovacc/dataprovider"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const (
	createTable = `CREATE TABLE IF NOT EXISTS sponsors (id SERIAL PRIMARY KEY, uuid TEXT, ulid TEXT, name TEXT, prize TEXT, timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	port := ":8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	// Create a config with driver name to initialize the data provider
	opts := dataprovider.NewOptions(
		dataprovider.WithDriver(dataprovider.SQLiteDataProviderName),
		dataprovider.WithConnectionString("file:test.sqlite3?cache=shared"),
	)

	provider, err := dataprovider.NewDataProvider(opts)
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize the database
	if err = provider.InitializeDatabase(createTable); err != nil {
		log.Fatalln(err)
	}

	db := provider.GetConnection()

	newHandler := handlers.NewHandler(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Timeout(60))
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Throttle(100))

	r.Get("/assets", newHandler.StaticHandler)
	r.Get("/", newHandler.HomeHandler)

	r.Get("/sponsor", newHandler.SponsorHandler)
	r.Get("/uuid/register", newHandler.RegisterHandler)
	r.Get("/sponsor/register/{uuid}", newHandler.SponsorRegisterHandler)
	r.Post("/sponsor/submit", newHandler.SponsorSubmitHandler)

	slog.Info("server started", slog.String("port", port))

	if err = http.ListenAndServe(port, r); err != nil {
		log.Fatalln(err)
	}
}
