package main

import (
	"job-service/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/sort-tasks", handlers.SortTasksHandler)
	http.HandleFunc("/sort-tasks/bash", handlers.SortTasksBashHandler)
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
