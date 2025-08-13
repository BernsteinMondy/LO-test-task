package storage

import (
	"context"
	"github.com/google/uuid"
	"lo-test-task/internal/asynclog"
	"lo-test-task/internal/core"
	"lo-test-task/internal/entity"
	"log/slog"
	"sync"
)

type Storage struct {
	tasks  map[uuid.UUID]*entity.Task
	mu     *sync.RWMutex
	logger asynclog.AsyncLogger
}

var _ core.Storage = (*Storage)(nil)

func New(logger asynclog.AsyncLogger) *Storage {
	return &Storage{
		tasks:  make(map[uuid.UUID]*entity.Task),
		mu:     &sync.RWMutex{},
		logger: logger,
	}
}

func (s *Storage) SaveTask(ctx context.Context, task *entity.Task) error {
	s.logger.Info(ctx, "Storage: Save task", slog.String("task.id", task.ID.String()))

	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task

	s.logger.Info(ctx, "Storage: Successfully saved task", slog.String("task.id", task.ID.String()))
	return nil
}

func (s *Storage) GetTaskByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	s.logger.Info(ctx, "Storage: Get task by ID", slog.String("task.id", id.String()))

	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, core.ErrRepoNotFound
	}

	s.logger.Info(ctx, "Storage: Successfully got task by ID", slog.String("task.id", id.String()))
	return task, nil
}

func (s *Storage) GetTasksByStatus(ctx context.Context, status entity.TaskStatus) ([]entity.Task, error) {
	s.logger.Info(ctx, "Storage: Get tasks by status")
	tasks := make([]entity.Task, 0, len(s.tasks)/2) // Allocating memory for slice, but because not all Task has the same statuses, allocate len(tasks)/2

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, task := range s.tasks {
		if task.Status == status {
			tasks = append(tasks, *task)
		}
	}

	s.logger.Info(ctx, "Storage: Successfully got tasks by status")
	return tasks, nil
}
