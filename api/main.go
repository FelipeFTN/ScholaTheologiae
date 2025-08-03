package main

import (
	"log"
	"log/slog"
	"os"

	"scholatheologiae-api/controllers"
	"scholatheologiae-api/data"
	"scholatheologiae-api/server"
)

func main() {
	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// Databases & SQLite
	data, err := data.New()
	terminateOnError(err)
	defer data.SQLite.Close()

	controllers := controllers.New(data)

	server.Run(controllers)
}

func terminateOnError(err error) {
	if err != nil {
		// Log the error and terminate the program
		log.Println("Error in dependencies:", err.Error())
		panic(err)
	}
}
