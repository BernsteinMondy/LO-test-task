package handlers

import "lo-test-task/internal/core"

type taskReadDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func taskToReadDTO(task *core.Task) taskReadDTO {
	return taskReadDTO{
		ID:          task.ID.String(),
		Title:       task.Title,
		Description: task.Description,
		Status:      taskStatusToString(task.Status),
	}
}
