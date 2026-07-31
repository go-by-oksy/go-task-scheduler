package main

import (
	"github.com/go-by-oksy/go-task-scheduler/pkg/api"
	"github.com/go-by-oksy/go-task-scheduler/pkg/db"
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
