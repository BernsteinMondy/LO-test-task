package handlers

import (
	"lo-test-task/internal/entity"
	"time"
)

type taskReadDTO struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func taskToReadDTO(task *entity.Task) taskReadDTO {
	return taskReadDTO{
		ID:          task.ID.String(),
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		Title:       task.Title,
		Description: task.Description,
		Status:      taskStatusToString(task.Status),
	}
}
