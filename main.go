package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var task string

type requestBody struct {
	Task string `json:"task"`
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	json.NewDecoder(r.Body).Decode(&body)

	task = body.Task

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("task saved"))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("hello, %s", task)))
}

func main() {
	http.HandleFunc("/task", postHandler)
	http.HandleFunc("/", getHandler)

	http.ListenAndServe(":8080", nil)
}

// test commit
