package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bmolinari/google-tasks-widget/internal/api"
	"github.com/bmolinari/google-tasks-widget/internal/handlers"
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
	taskHandler := handlers.TaskHandler{TaskService: taskService}

	http.HandleFunc("/tasks", taskHandler.GetTasks)
	http.HandleFunc("/tasks/create", taskHandler.CreateTask)
	http.HandleFunc("/tasks/complete", taskHandler.CompleteTask)

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
