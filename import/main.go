package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/deitrix/fin/auth"
	finfile "github.com/deitrix/fin/store/file"
	finmysql "github.com/deitrix/fin/store/mysql"
	"github.com/deitrix/sqlg"
	"github.com/go-sql-driver/mysql"
)

func main() {
	conf, err := readConfig("config.json")
	if err != nil {
		log.Fatalf("reading config: %v", err)
	}

	db, err := sql.Open("mysql", conf.DB.DSN())
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}

	fileStore := finfile.NewStore("real.fin.json")
	sqlStore := finmysql.NewStore(db)

	rps, err := fileStore.RecurringPayments(nil)
	if err != nil {
		log.Fatalf("loading recurring payments: %v", err)
	}

	for _, rp := range rps {
		if err := sqlStore.CreateRecurringPayment(context.Background(), rp); err != nil {
			log.Fatalf("creating recurring payment: %v", err)
		}
	}
}

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
