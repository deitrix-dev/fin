package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/deitrix/fin"
	finfile "github.com/deitrix/fin/store/file"
	finmysql "github.com/deitrix/fin/store/mysql"
	"github.com/deitrix/sqlg"
	"github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v3"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cfn := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cfn()

	return rootCmd.Run(ctx, os.Args)
}

var rootCmd = &cli.Command{
	Name: "finctl",
	Commands: []*cli.Command{
		{
			Name:      "export",
			UsageText: "export <file>",
			Usage:     "export data from database into file",
			Action:    exportCmd,
		},
		{
			Name:      "import",
			UsageText: "import <file>",
			Usage:     "import data from file into database",
			Action:    importCmd,
		},
		{
			Name:   "clear",
			Usage:  "clear data from database",
			Action: clearCmd,
		},
	},
}

func exportCmd(_ context.Context, cmd *cli.Command) error {
	filename := cmd.Args().First()
	if filename == "" {
		return fmt.Errorf("missing file argument")
	}

	if filepath.Ext(filename) != ".json" {
		return fmt.Errorf("file must have .json extension")
	}

	conf, err := readConfig("config.json")
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	db, err := sql.Open("mysql", conf.DB.DSN())
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}

	fileStore := finfile.NewStore(filename)
	sqlStore := finmysql.NewStore(db)

	rps, err := sqlStore.RecurringPayments(context.Background(), fin.RecurringPaymentFilter{})
	if err != nil {
		return fmt.Errorf("loading recurring payments: %w", err)
	}

	for _, rp := range rps {
		if err := fileStore.CreateRecurringPayment(context.Background(), rp); err != nil {
			return fmt.Errorf("creating recurring payment: %w", err)
		}
	}

	return nil
}

func importCmd(_ context.Context, cmd *cli.Command) error {
	filename := cmd.Args().First()
	if filename == "" {
		return fmt.Errorf("missing file argument")
	}

	if filepath.Ext(filename) != ".json" {
		return fmt.Errorf("file must have .json extension")
	}

	conf, err := readConfig("config.json")
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	db, err := sql.Open("mysql", conf.DB.DSN())
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}

	fileStore := finfile.NewStore(filename)
	sqlStore := finmysql.NewStore(db)

	rps, err := fileStore.RecurringPayments(nil)
	if err != nil {
		return fmt.Errorf("loading recurring payments: %w", err)
	}

	for _, rp := range rps {
		if err := sqlStore.CreateRecurringPayment(context.Background(), rp); err != nil {
			return fmt.Errorf("creating recurring payment: %w", err)
		}
	}

	return nil
}

func clearCmd(_ context.Context, _ *cli.Command) error {
	conf, err := readConfig("config.json")
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	db, err := sql.Open("mysql", conf.DB.DSN())
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}

	sqlStore := finmysql.NewStore(db)

	rps, err := sqlStore.RecurringPayments(context.Background(), fin.RecurringPaymentFilter{})
	if err != nil {
		return fmt.Errorf("loading recurring payments: %w", err)
	}

	for _, rp := range rps {
		if err := sqlStore.DeleteRecurringPayment(context.Background(), rp.ID); err != nil {
			return fmt.Errorf("deleting recurring payment: %w", err)
		}
	}

	return nil
}

type config struct {
	DB dbConfig `json:"db"`
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
