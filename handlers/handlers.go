package handlers

import (
	"encoding/json"
	"job-service/tasks"
	"net/http"
)

func getTasksFromReq(req *http.Request) (tasks.Tasks, error) {
	var tasksArr tasks.Tasks
	err := json.NewDecoder(req.Body).Decode(&tasksArr)
	if err != nil {
		return tasksArr, err
	}

	return tasksArr, nil
}

func SortTasksHandler(w http.ResponseWriter, req *http.Request) {
	t, err := getTasksFromReq(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderedTasks, err := tasks.TopologicallyOrderTasks(t.Tasks)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, _ := json.Marshal(orderedTasks)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func SortTasksBashHandler(w http.ResponseWriter, req *http.Request) {
	t, err := getTasksFromReq(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderedTasks, err := tasks.TopologicallyOrderTasks(t.Tasks)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	w.Write([]byte("#!/usr/bin/env bash\n"))
	for _, v := range orderedTasks {
		w.Write([]byte("\n"))
		w.Write([]byte(v.Command))
	}
}
