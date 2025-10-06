package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_calendar/internal/config"
	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_calendar/internal/storage/postgresql"
	"github.com/jmoiron/sqlx"
	goose "github.com/pressly/goose/v3"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Migrate - утилита для управления миграциями БД сервиса calendar (на основе goose)\n\n"+
				"вызов: migrate -config=<config file name> -dir=<migration dir> -command=<migration command> arg\n\n"+
				"Примеры:\n"+
				"  migrate -config=config.yaml -dir=./migrations -command up 0002\n"+
				"  migrate -command status\n\n"+
				"Доступные флаги:\n")
		flag.PrintDefaults()
	}
}

func main() {
	var (
		configFile    string
		migrationsDir string
		command       string
		arg           string
		v             int64
		db            *sql.DB
		dbx           *sqlx.DB
		err           error
	)

	flag.StringVar(&configFile, "config", "config.yaml", "path to config file")
	flag.StringVar(&migrationsDir, "dir", "./migrations", "path to migrations dir")
	flag.StringVar(&command, "command", "up", "goose command: up|down|redo|reset|status|version|up-to|down-to")
	flag.StringVar(&arg, "arg", "", "argument for command (version for up-to/down-to/force)")
	flag.Parse()

	// Загружаем конфиг
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	logg := logger.New(&cfg.Logger)

	// Подключаемся к БД
	if dbx, err = postgresql.OpenDB(cfg.Database); err != nil {
		logg.Fatal("DB open error", "error", err)
	}

	defer dbx.Close()
	db = dbx.DB

	if err := goose.SetDialect("postgres"); err != nil {
		logg.Fatal("set dialect", "error", err)
	}

	cmd := strings.ToLower(command)
	switch cmd {
	case "up":
		err = goose.Up(db, migrationsDir)
	case "down":
		err = goose.Down(db, migrationsDir)
	case "redo":
		err = goose.Redo(db, migrationsDir)
	case "reset":
		err = goose.Reset(db, migrationsDir)
	case "status":
		err = goose.Status(db, migrationsDir)
	case "version":
		v, err2 := goose.GetDBVersion(db)
		if err2 != nil {
			logg.Error("version", "error", err2)
		}
		fmt.Printf("Current version: %d\n", v)
		return
	case "up-to":
		if v, err = mustParseInt64(arg); err != nil {
			logg.Fatal("parse version error", "error", err)
		}
		err = goose.UpTo(db, migrationsDir, v)
	case "down-to":
		if v, err = mustParseInt64(arg); err != nil {
			logg.Fatal("parse version error", "error", err)
		}
		err = goose.DownTo(db, migrationsDir, v)

	default:
		logg.Fatal("unknown command", "command", command)
	}

	if err != nil {
		logg.Fatal("migration error", "command", cmd, "error", err)
	}

	logg.Info("migration completed", "command", cmd)
	fmt.Println("migration completed:", cmd)
}

// Проверка корректности номера версии. Должен конвертироваться в int64.
func mustParseInt64(s string) (int64, error) {
	if s == "" {
		return 0, errors.New("arg is required for this command")
	}
	var v int64
	_, err := fmt.Sscan(s, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}
