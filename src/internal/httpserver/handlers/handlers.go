package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"lo-test-task/internal/asynclog"
	"lo-test-task/internal/core"
	"lo-test-task/internal/entity"
	"net/http"
)

type Service interface {
	GetTaskByID(context.Context, uuid.UUID) (*entity.Task, error)
	GetTasksByStatus(ctx context.Context, status entity.TaskStatus) ([]entity.Task, error)
	CreateNewTask(ctx context.Context, title, description string, status entity.TaskStatus) (uuid.UUID, error)
}

func Map(mux *http.ServeMux, service Service, logger asynclog.AsyncLogger) {
	mux.HandleFunc("GET /tasks", getTasksHandler(service, logger))
	mux.HandleFunc("POST /tasks", postTasksHandler(service, logger))
	mux.HandleFunc("GET /tasks/{id}", getTaskHandler(service, logger))
}

func getTasksHandler(service Service, logger asynclog.AsyncLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.Context(), "Received request for GET /tasks")

		statusStr := r.URL.Query().Get("status")
		if statusStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		status, ok := tryConvertStringToTaskStatus(statusStr)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		tasks, err := service.GetTasksByStatus(ctx, status)
		if err != nil {
			logger.Error(r.Context(), "Service: Failed to get tasks by status", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := make([]taskReadDTO, 0, len(tasks))
		for _, task := range tasks {
			resp = append(resp, taskToReadDTO(&task))
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info(r.Context(), "Successfully GET /tasks")
	}
}

func getTaskHandler(service Service, logger asynclog.AsyncLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.Context(), "Received request for GET /tasks/{id}")

		idStr := r.PathValue("id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		task, err := service.GetTaskByID(ctx, id)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			logger.Error(r.Context(), "Service: Failed to get task by ID", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := taskToReadDTO(task)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info(r.Context(), "Successfully GET /tasks/{id}")
	}
}

func postTasksHandler(service Service, logger asynclog.AsyncLogger) http.HandlerFunc {
	type request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	type response struct {
		ID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.Context(), "Received request for POST /tasks")

		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Title == "" || req.Description == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		status, ok := tryConvertStringToTaskStatus(req.Status)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		id, err := service.CreateNewTask(ctx, req.Title, req.Description, status)
		if err != nil {
			logger.Error(r.Context(), "Service: Failed to create new task", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var resp = response{
			ID: id.String(),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info(r.Context(), "Successfully POST /tasks")
	}
}
