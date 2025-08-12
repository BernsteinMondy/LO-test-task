package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"lo-test-task/internal/core"
	"net/http"
)

func Map(mux *http.ServeMux, service *core.Service) {
	mux.HandleFunc("GET /tasks", getTasksHandler(service))
	mux.HandleFunc("POST /tasks", postTasksHandler(service))
	mux.HandleFunc("GET /tasks/{id}", getTaskHandler(service))
}

func getTasksHandler(service *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func getTaskHandler(service *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func postTasksHandler(service *core.Service) http.HandlerFunc {
	type request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	type response struct {
		ID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
