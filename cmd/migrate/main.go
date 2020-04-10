package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

var srcUrl string
var command string
var name string

func init() {
	flag.StringVar(&srcUrl, "source", "file://migrations", "Where to read migrations from")
	flag.StringVar(&command, "command", "up", "What migration command to run [up/down/new]")
	flag.StringVar(&name, "name", "", "Name of new migration file")
}

func main() {
	flag.Parse()

	m, err := migrate.New(srcUrl, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	if command == "up" {
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				fmt.Println("Nothing to migrate")
			} else {
				panic(err)
			}
		}
	} else if command == "down" {
		if err := m.Steps(-1); err != nil {
			if err == migrate.ErrNoChange {
				fmt.Println("Nothing to migrate")
			} else {
				panic(err)
			}
		}
	} else if command == "new" {
		if name == "" {
			fmt.Println("name is required")
		} else {
			fileName := fmt.Sprintf("%v_%s.%%s.sql", time.Now().Unix(), name)
			upName := fmt.Sprintf(fileName, "up")
			downName := fmt.Sprintf(fileName, "down")
			os.Create(fmt.Sprintf("migrations/%s", upName))
			os.Create(fmt.Sprintf("migrations/%s", downName))
		}
	} else {
		panic("unknown command")
	}

	fmt.Println("Done with raw sql migrations")
}
