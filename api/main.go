package main

import (
	"log"
	"log/slog"
	"os"

	"scholatheologiae-api/controller"
	"scholatheologiae-api/data"
	"scholatheologiae-api/handler"
	"scholatheologiae-api/server"
)

func main() {
	// DEPENDENCIES

	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// SQLite
	db, err := handler.NewSQLite(data.SQLiteMAP)
	terminateOnError(err)
	defer db.Close()

	controller := controller.New(db)

	server.Run(controller)
}

func terminateOnError(err error) {
	if err != nil {
		// Log the error and terminate the program
		log.Println("Error in dependencies:", err.Error())
		panic(err)
	}
}
