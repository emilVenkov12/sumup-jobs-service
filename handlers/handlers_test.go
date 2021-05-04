package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestErrorsOnInvalidJson(t *testing.T) {
	req, err := http.NewRequest("POST", "/sort-tasks", strings.NewReader(`{"tasks": {`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SortTasksHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestErrorsOnInvalidTasks(t *testing.T) {
	req, err := http.NewRequest("POST", "/sort-tasks", strings.NewReader(`{"tasks": [{"name": "task-1", "command": "cat /tmp/file1", "requires":["task-2"]},{"name": "task-2", "command": "rm /tmp/file1", "requires":["task-1"]}]}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SortTasksHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestReturnsCorrectlyOrderedJsonTasks(t *testing.T) {
	req, err := http.NewRequest("POST", "/sort-tasks", strings.NewReader(`{"tasks": [{"name": "task-1", "command": "cat /tmp/file1", "requires":["task-2"]},{"name": "task-2", "command": "touch /tmp/file1"}]}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SortTasksHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedContentType := "application/json"
	if conentType := rr.Header().Get("Content-Type"); conentType != expectedContentType {
		t.Errorf("handler returned wrong Content-Type: got %v want %v", conentType, expectedContentType)
	}

	response, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := `[{"name":"task-2","command":"touch /tmp/file1"},{"name":"task-1","command":"cat /tmp/file1"}]`

	if string(response[:]) != expectedResponse {
		t.Errorf("handler returned wrong tasks: got %v expected %v", string(response[:]), expectedResponse)
	}
}

func TestReturnsBashCommandsFromTasks(t *testing.T) {
	req, err := http.NewRequest("POST", "/sort-tasks/bash", strings.NewReader(`{"tasks": [{"name": "task-1", "command": "cat /tmp/file1", "requires":["task-2"]},{"name": "task-2", "command": "touch /tmp/file1"}]}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SortTasksBashHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedContentType := "text/plain;charset=UTF-8"
	if conentType := rr.Header().Get("Content-Type"); conentType != expectedContentType {
		t.Errorf("handler returned wrong Content-Type: got %v want %v", conentType, expectedContentType)
	}

	response, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := `#!/usr/bin/env bash

touch /tmp/file1
cat /tmp/file1`

	if string(response[:]) != expectedResponse {
		t.Errorf("handler returned wrong tasks: got %v expected %v", string(response[:]), expectedResponse)
	}
}
