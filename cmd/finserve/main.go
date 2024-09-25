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
	"time"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/auth"
	"github.com/deitrix/fin/ui/api"
	"github.com/deitrix/fin/ui/handlers"
	"github.com/deitrix/fin/web/assets"
	"github.com/deitrix/sqlg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh/terminal"

	finmysql "github.com/deitrix/fin/store/mysql"
	slogchi "github.com/samber/slog-chi"

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

var defaultHeaders = []string{
	"CF-Connecting-IP",
	"X-Real-IP",
	"X-Forwarded-For",
}

func main() {
	logOpt := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var logHandler slog.Handler
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		logHandler = slog.NewTextHandler(os.Stdout, logOpt)
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, logOpt)
	}

	logger := slog.New(logHandler)
	slog.SetDefault(logger)

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
	svc := fin.NewService(store)

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))
	router.Use(middleware.RealIPFromHeaders(defaultHeaders...))
	router.Use(slogchi.New(logger))
	router.Use(auth.Verify(conf.Auth))

	router.Get("/assets/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, assets.FS, "style.css")
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.SetHeader("Cache-Control", "no-store"))

		r.Get("/", handlers.Home)
		r.Get("/recurring-payments/{id}", handlers.RecurringPayment(store))
		r.Get("/recurring-payments/{id}/form", handlers.RecurringPaymentUpdateForm(store))
		r.Post("/recurring-payments/{id}/form", handlers.RecurringPaymentHandleUpdateForm(store))
		r.Get("/recurring-payments/{id}/delete", handlers.RecurringPaymentDelete(store))
		r.Get("/create", handlers.RecurringPaymentCreate())
		r.Post("/create", handlers.RecurringPaymentCreateForm(store))
		r.Get("/recurring-payments/{id}/schedules/{index}/delete", handlers.ScheduleDelete(store))
		r.Get("/recurring-payments/{id}/schedules/new", handlers.ScheduleForm(store))
		r.Post("/recurring-payments/{id}/schedules/new", handlers.ScheduleHandleForm(store))
		r.Get("/recurring-payments/{id}/schedules/{index}", handlers.ScheduleForm(store))
		r.Post("/recurring-payments/{id}/schedules/{index}", handlers.ScheduleHandleForm(store))

		r.Get("/payments/new", handlers.PaymentForm(store))
		r.Post("/payments/new", handlers.PaymentHandleForm(store))
		r.Get("/payments/{id}", handlers.PaymentForm(store))
		r.Post("/payments/{id}", handlers.PaymentHandleForm(store))
		r.Get("/payments/{id}/delete", handlers.PaymentHandleDelete(store))

		r.Get("/api/recurring-payments", api.RecurringPayments(store))
		r.Get("/api/payments", api.Payments(svc))
		r.Get("/api/payments-for-schedule", api.PaymentsForSchedule)
		r.Get("/api/header-user", api.HeaderUser(conf.SimulateUser))
	})
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
