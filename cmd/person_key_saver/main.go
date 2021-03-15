package main

import (
	"log"
	"person-key-saver/internal/app/application"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app := application.New()

	if err := app.Configure("app"); err != nil {
		return err
	}

	app.Run()
	return nil
}
