package entity

type TaskStatus uint8

const (
	TaskStatusDone TaskStatus = iota
	TaskStatusInProgress
	TaskStatusCreated
)
