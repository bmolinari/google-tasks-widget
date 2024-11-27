package handlers

import (
	"encoding/json"
	"net/http"

	"google.golang.org/api/tasks/v1"
)

type TaskHandler struct {
	TaskService *tasks.Service
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	taskList, err := h.TaskService.Tasks.List("@default").Do()
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskList.Items)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task := &tasks.Task{Title: req.Title}
	_, err := h.TaskService.Tasks.Insert("@default", task).Do()
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TaskId string `json:"task_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.TaskService.Tasks.Get("@default", req.TaskId).Do()
	if err != nil {
		http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		return
	}
	task.Status = "completed"
	_, err = h.TaskService.Tasks.Update("@default", task.Id, task).Do()
	if err != nil {
		http.Error(w, "Failed to mark task as complete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
