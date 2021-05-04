package main

import (
	"job-service/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/sort-tasks", handlers.SortTasksHandler)
	http.HandleFunc("/sort-tasks/bash", handlers.SortTasksBashHandler)
	http.ListenAndServe(":4000", nil)
}
