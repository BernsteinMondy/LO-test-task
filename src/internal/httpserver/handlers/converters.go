package handlers

import (
	"fmt"
	"lo-test-task/internal/entity"
)

func tryConvertStringToTaskStatus(str string) (entity.TaskStatus, bool) {
	switch str {
	case "done":
		return entity.TaskStatusDone, true
	case "created":
		return entity.TaskStatusCreated, true
	case "in-progress":
		return entity.TaskStatusInProgress, true
	default:
		return 0, false
	}
}

func taskStatusToString(taskStatus entity.TaskStatus) string {
	switch taskStatus {
	case entity.TaskStatusDone:
		return "done"
	case entity.TaskStatusCreated:
		return "created"
	case entity.TaskStatusInProgress:
		return "in-progress"
	default:
		panic(fmt.Sprintf("unknown core.TaskStatus (%v)", taskStatus))
	}
}
