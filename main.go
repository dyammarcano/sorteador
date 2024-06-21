package main

import (
	"embed"
	"github.com/dyammarcano/sorteador/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/inovacc/dataprovider"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const (
	createTable = `CREATE TABLE IF NOT EXISTS pedidos (id SERIAL PRIMARY KEY, codigo_pedido INT, codigo_cliente INT)`
)

//go:embed static
var static embed.FS

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	sfs, _ := fs.Sub(static, "static")
	http.FileServer(http.FS(sfs)).ServeHTTP(w, r)
}

func main() {
	port := ":8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	// Create a config with driver name to initialize the data provider
	opts := dataprovider.NewOptions(
		dataprovider.WithDriver(dataprovider.PostgreSQLDatabaseProviderName),
		dataprovider.WithUsername("postgres"),
		dataprovider.WithPassword("mysecretpassword"),
		dataprovider.WithHost("db_postgres"),
		dataprovider.WithName("postgres"),
		dataprovider.WithPort(5432),
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

	r.Get("/", staticHandler)

	r.Get("/sponsor", newHandler.SponsorHandler)
	r.Get("/uuid/register", newHandler.RegisterHandler)

	if err = http.ListenAndServe(port, r); err != nil {
		log.Fatalln(err)
	}
}
