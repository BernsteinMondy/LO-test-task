package handlers

import (
	"fmt"
	"lo-test-task/internal/core"
)

func tryConvertStringToTaskStatus(str string) (core.TaskStatus, bool) {
	switch str {
	case "done":
		return core.TaskStatusDone, true
	case "created":
		return core.TaskStatusCreated, true
	case "in-progress":
		return core.TaskStatusInProgress, true
	default:
		return 0, false
	}
}

func taskStatusToString(taskStatus core.TaskStatus) string {
	switch taskStatus {
	case core.TaskStatusDone:
		return "done"
	case core.TaskStatusCreated:
		return "created"
	case core.TaskStatusInProgress:
		return "in-progress"
	default:
		panic(fmt.Sprintf("unknown core.TaskStatus (%v)", taskStatus))
	}
}
