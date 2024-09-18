package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/deitrix/fin/auth"
	"github.com/deitrix/fin/http/api"
	"github.com/deitrix/fin/http/pages"
	"github.com/deitrix/fin/store/file"
	"github.com/deitrix/fin/web/assets"
	"github.com/go-chi/chi/v5"
)

type config struct {
	Auth auth.Config `json:"auth"`
}

func readConfig(path string) (config, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return config{}, fmt.Errorf("reading config file: %w", err)
	}
	var c config
	if err := json.Unmarshal(bs, &c); err != nil {
		return config{}, fmt.Errorf("unmarshalling config file: %w", err)
	}
	return c, nil
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	configPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()
	if *configPath == "" {
		log.Fatal("config file path is required")
	}

	conf, err := readConfig(*configPath)
	if err != nil {
		log.Fatalf("reading config: %v", err)
	}

	store := file.NewStore("fin.json")

	router := chi.NewRouter()

	router.Use(auth.Verify(conf.Auth))

	router.Get("/assets/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, assets.FS, "style.css")
	})

	router.Get("/", pages.Home(store))
	router.Get("/recurring-payments/{id}", pages.RecurringPaymentByID(store))
	router.Get("/create", pages.Create(store))
	router.Post("/create", pages.CreatePOST(store))

	router.Get("/api/payments", api.Payments(store))

	log.Fatal(http.ListenAndServe(":8080", router))
}

func render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	if err := component.Render(r.Context(), w); err != nil {
		slog.ErrorContext(r.Context(), "error rendering page", err)
	}
}
