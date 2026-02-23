package main

import (
	"encoding/json"
	"net/http"
)

var task string

type requestBody struct {
	Task string `json:"task"`
}

func postTask(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body requestBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	task = body.Task

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("task saved"))
}

func main() {

	http.HandleFunc("/task", postTask)
	http.ListenAndServe(":8080", nil)

}
