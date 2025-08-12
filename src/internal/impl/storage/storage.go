package storage

import (
	"github.com/google/uuid"
	"lo-test-task/internal/core"
	"sync"
)

type Storage struct {
	tasks map[uuid.UUID]*core.Task
	mu    *sync.RWMutex
}

var _ core.Storage = (*Storage)(nil)

func New() *Storage {
	return &Storage{
		tasks: make(map[uuid.UUID]*core.Task),
		mu:    &sync.RWMutex{},
	}
}

func (s *Storage) SaveTask(task *core.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task

	return nil
}

func (s *Storage) GetTaskByID(id uuid.UUID) (*core.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, core.ErrRepoNotFound
	}

	return task, nil
}

func (s *Storage) GetTasksByStatus(status core.TaskStatus) ([]core.Task, error) {
	tasks := make([]core.Task, 0, len(s.tasks)/2) // TODO: Explain

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, task := range s.tasks {
		if task.Status == status {
			tasks = append(tasks, *task)
		}
	}

	return tasks, nil
}
