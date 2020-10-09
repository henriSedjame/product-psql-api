package main

import (
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	app := App{
		Log: logger,
	}

	app.Initialize()

	app.Run()
}
