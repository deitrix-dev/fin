package main

import (
	"database/sql"
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
	"github.com/deitrix/fin/web/assets"
	"github.com/deitrix/fin/web/page"
	"github.com/deitrix/sqlg"
	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
	"github.com/rickb777/date"

	finmysql "github.com/deitrix/fin/store/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	Auth         auth.Config `json:"auth"`
	DB           dbConfig    `json:"db"`
	SimulateUser string      `json:"simulateUser"`
}

type dbConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Pass      string `json:"pass"`
	DB        string `json:"db"`
	ParseTime bool   `json:"parseTime"`
}

func (c dbConfig) DSN() string {
	conf := mysql.NewConfig()
	conf.Net = "tcp"
	conf.Addr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	conf.User = c.User
	conf.Passwd = c.Pass
	conf.DBName = c.DB
	conf.ParseTime = c.ParseTime
	return conf.FormatDSN()
}

func readConfig(path string) (config, error) {
	sqlg.Debug = true

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

	db, err := sql.Open("mysql", conf.DB.DSN())
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}

	store := finmysql.NewStore(db)

	router := chi.NewRouter()

	router.Use(auth.Verify(conf.Auth))

	router.Get("/assets/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, assets.FS, "style.css")
	})

	router.Get("/", pages.Home(store))
	router.Get("/recurring-payments/{id}", pages.RecurringPayment(store))
	router.Get("/recurring-payments/{id}/form", pages.RecurringPaymentUpdateForm(store))
	router.Post("/recurring-payments/{id}/form", pages.RecurringPaymentHandleUpdateForm(store))
	router.Get("/recurring-payments/{id}/delete", pages.RecurringPaymentDelete(store))
	router.Get("/create", pages.Create())
	router.Post("/create", pages.CreatePOST(store))
	router.Get("/form/schedule/{recurringPaymentID}", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var values page.ScheduleFormValues
		if s := r.Form.Get("startDate"); s != "" {
			startDate, err := date.ParseISO(s)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			values.StartDate = startDate
		}
		values.Every = r.Form.Get("every")
		values.Day = r.Form.Get("day")
		render(w, r, page.ScheduleForm(chi.URLParam(r, "recurringPaymentID"), values))
	})
	router.Get("/schedule/{recurringPaymentID}", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var values page.ScheduleFormValues
		if s := r.Form.Get("startDate"); s != "" {
			startDate, err := date.ParseISO(s)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			values.StartDate = startDate
		}
		render(w, r, page.CreateSchedule(chi.URLParam(r, "recurringPaymentID"), values))
	})

	router.Get("/api/payments", api.Payments(store))
	router.Get("/api/header-user", api.HeaderUser(conf.SimulateUser))

	log.Fatal(http.ListenAndServe(":8080", router))
}

func render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	if err := component.Render(r.Context(), w); err != nil {
		slog.ErrorContext(r.Context(), "error rendering page", err)
	}
}
