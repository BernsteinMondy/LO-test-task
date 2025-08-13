package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"lo-test-task/internal/asynclog"
	"lo-test-task/internal/entity"
	"log/slog"
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
	SaveTask(ctx context.Context, task *entity.Task) error
	// GetTaskByID must return ErrRepoNotFound if no Task was found by the given ID.
	GetTaskByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)
	GetTasksByStatus(ctx context.Context, status entity.TaskStatus) ([]entity.Task, error)
}

func (s *Service) GetTaskByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	s.logger.Info(ctx, "Service: Getting task by ID", slog.String("task.id", id.String()))
	task, err := s.storage.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrRepoNotFound) {
			s.logger.Error(ctx, "Service: Task not found the given ID", err, slog.String("task.id", id.String()))
			return nil, ErrNotFound
		}

		err = fmt.Errorf("storage: get task by id: %w", err)

		s.logger.Error(ctx, "Service: Failed to get task by ID", err, slog.String("task.id", id.String()))
		return nil, err
	}

	s.logger.Info(ctx, "Service: Successfully got Task by ID", slog.String("task.id", id.String()))
	return task, nil
}

func (s *Service) GetTasksByStatus(ctx context.Context, status entity.TaskStatus) ([]entity.Task, error) {
	s.logger.Info(ctx, "Service: Getting tasks by status")
	tasks, err := s.storage.GetTasksByStatus(ctx, status)
	if err != nil {

		err = fmt.Errorf("storage: get tasks by status: %w", err)
		s.logger.Error(ctx, "Service: Failed to get tasks by status", err)

		return nil, err
	}

	s.logger.Info(ctx, "Service: Successfully got tasks by status")
	return tasks, nil
}

func (s *Service) CreateNewTask(ctx context.Context, title, description string, status entity.TaskStatus) (uuid.UUID, error) {
	s.logger.Info(ctx, "Service: Creating a new task")

	id := uuid.New()

	task := &entity.Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
	}

	err := s.storage.SaveTask(ctx, task)
	if err != nil {
		err = fmt.Errorf("storage: save task: %w", err)

		s.logger.Error(ctx, "Service: Failed to save task", err)
		return uuid.Nil, fmt.Errorf("storage: save task: %w", err)
	}

	s.logger.Info(ctx, "Service: Successfully created a new task", slog.String("task.id", id.String()))
	return id, nil
}
