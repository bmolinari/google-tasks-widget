package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/bmolinari/google-tasks-widget/internal/api"
	"google.golang.org/api/tasks/v1"
)

var taskService *tasks.Service

func main() {
	auth := api.NewGoogleAuth()

	token, err := auth.LoadToken("config/token.json")
	if err != nil {
		log.Println("Token not found, running OAuth flow...")
		token = auth.GetTokenFromWeb()
		auth.SaveToken("config/token.json", token)
	}
	auth.Token = token

	taskService = auth.GetService()

	http.HandleFunc("/tasks", handleGetTasks)
	http.HandleFunc("/tasks/create", handleCreateTask)
	http.HandleFunc("/tasks/complete", handleCompleteTask)

	port := getServerPort()
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func handleGetTasks(w http.ResponseWriter, r *http.Request) {
	taskList, err := taskService.Tasks.List("@default").Do()
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskList.Items)
}

func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task := &tasks.Task{Title: req.Title}
	_, err := taskService.Tasks.Insert("@default", task).Do()
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleCompleteTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TaskId string `json:"task_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, err := taskService.Tasks.Get("@default", req.TaskId).Do()
	if err != nil {
		http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		return
	}
	task.Status = "completed"
	_, err = taskService.Tasks.Update("@default", task.Id, task).Do()
	if err != nil {
		http.Error(w, "Failed to mark task as complete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
