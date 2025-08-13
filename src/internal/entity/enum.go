package entity

type TaskStatus uint8

const (
	TaskStatusDone TaskStatus = iota + 1
	TaskStatusInProgress
	TaskStatusCreated
)
