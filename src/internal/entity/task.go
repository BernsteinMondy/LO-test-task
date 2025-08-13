package entity

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	Title       string
	Description string
	Status      TaskStatus
}
