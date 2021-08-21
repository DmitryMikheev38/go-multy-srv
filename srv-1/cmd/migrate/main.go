package main

import (
	"flag"
	"fmt"
	"github.com/pressly/goose"
	"log"
	"os"
	"srv-1/internal/app/config"

	_ "github.com/lib/pq"
	_ "srv-1/migrations"
)

var (
	flags      = flag.NewFlagSet("goose", flag.ExitOnError)
	dir        = flags.String("dir", "./migrations", "directory with migrate files")
	table      = flags.String("table", "goose_db_version", "migrations table name")
	verbose    = flags.Bool("v", false, "enable verbose mode")
	help       = flags.Bool("h", false, "print help")
	version    = flags.Bool("version", false, "print version")
	sequential = flags.Bool("s", false, "use sequential numbering for new migrations")
)

func main() {
	appConf, err := config.GetConf()
	if err != nil {
		log.Fatal("Cannot read config file")
	}

	flags.Usage = usage
	flags.Parse(os.Args[1:])

	if *version {
		fmt.Println(goose.VERSION)
		return
	}
	if *verbose {
		goose.SetVerbose(true)
	}
	if *sequential {
		goose.SetSequential(true)
	}
	goose.SetTableName(*table)

	args := flags.Args()
	if len(args) == 0 || *help {
		flags.Usage()
		return
	}

	switch args[0] {
	case "create":
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, *dir); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	}

	args = mergeArgs(args)
	command := args[0]
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		appConf.Postgres.Host, appConf.Postgres.User, appConf.Postgres.Pwd, appConf.Postgres.DB, appConf.Postgres.Port)
	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		log.Fatalf("-dbstring=%q: %v\n", dsn, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()
	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}
	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}

const (
	envGooseDriver   = "GOOSE_DRIVER"
	envGooseDBString = "GOOSE_DBSTRING"
)

func mergeArgs(args []string) []string {
	if len(args) < 1 {
		return args
	}
	if d := os.Getenv(envGooseDriver); d != "" {
		args = append([]string{d}, args...)
	}
	if d := os.Getenv(envGooseDBString); d != "" {
		args = append([]string{args[0], d}, args[1:]...)
	}
	return args
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `
Usage: goose [OPTIONS] COMMAND

Drivers:

    postgres
    mysql
    sqlite3
    mssql
    redshift
    clickhouse

Options:
`
	usageCommands = `
Commands:

    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migrate
    reset                Roll back all migrations
    status               Dump the migrate status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migrate file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)