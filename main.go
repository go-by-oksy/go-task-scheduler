package main

import (
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
