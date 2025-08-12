package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"lo-test-task/internal/asynclog"
	"log/slog"
)

type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      TaskStatus
}

type TaskStatus uint8

const (
	TaskStatusDone TaskStatus = iota
	TaskStatusInProgress
	TaskStatusCreated
)

type Service struct {
	storage Storage
	logger  asynclog.AsyncLogger
}

func NewService(storage Storage, logger asynclog.AsyncLogger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}

type Storage interface {
	SaveTask(task *Task) error
	// GetTaskByID must return ErrRepoNotFound if no Task was found by the given ID.
	GetTaskByID(id uuid.UUID) (*Task, error)
	GetTasksByStatus(status TaskStatus) ([]Task, error)
}

func (s *Service) GetTaskByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	s.logger.Info(ctx, "Service: Getting task by ID", slog.String("task.id", id.String()))
	task, err := s.storage.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, ErrRepoNotFound) {
			return nil, ErrNotFound
		}

		err = fmt.Errorf("storage: get task by id: %w", err)

		s.logger.Error(ctx, "Service: Failed to get task by ID", err, slog.String("task.id", id.String()))
		return nil, err
	}

	s.logger.Info(ctx, "Service: Successfully got Task by ID", slog.String("task.id", id.String()))
	return task, nil
}

func (s *Service) GetTasksByStatus(ctx context.Context, status TaskStatus) ([]Task, error) {
	s.logger.Info(ctx, "Service: Getting tasks by status")
	tasks, err := s.storage.GetTasksByStatus(status)
	if err != nil {

		err = fmt.Errorf("storage: get tasks by status: %w", err)
		s.logger.Error(ctx, "Service: Failed to get tasks by status", err)

		return nil, err
	}

	s.logger.Info(ctx, "Service: Successfully got tasks by status")
	return tasks, nil
}

func (s *Service) CreateNewTask(ctx context.Context, title, description string, status TaskStatus) (uuid.UUID, error) {
	s.logger.Info(ctx, "Service: Creating a new task")

	id := uuid.New()

	task := &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
	}

	err := s.storage.SaveTask(task)
	if err != nil {
		err = fmt.Errorf("storage: save task: %w", err)

		s.logger.Error(ctx, "Service: Failed to save task", err)
		return uuid.Nil, fmt.Errorf("storage: save task: %w", err)
	}

	s.logger.Info(ctx, "Service: Successfully created a new task", slog.String("task.id", id.String()))
	return id, nil
}
