package main

import (
	"final-project/pkg/api"
	"final-project/pkg/db"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		panic(err)
	}
	defer db.Close()

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
