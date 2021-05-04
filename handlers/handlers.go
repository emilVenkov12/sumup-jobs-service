package handlers

import (
	"bytes"
	"encoding/json"
	"job-service/tasks"
	"log"
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
	_, err = w.Write(res)
	if err != nil {
		log.Printf("Write to response writer failed: %v", err)
	}
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

	var res bytes.Buffer
	_, err = res.WriteString("#!/usr/bin/env bash\n")
	if err != nil {
		log.Printf("writing to bytes buffer failed: %v", err)
	}

	for _, v := range orderedTasks {
		_, err = res.WriteString("\n" + v.Command)
		if err != nil {
			log.Printf("writing to bytes buffer failed: %v", err)
		}
	}

	_, err = w.Write(res.Bytes())
	if err != nil {
		log.Printf("writing to response writer failed: %v", err)
	}
}
