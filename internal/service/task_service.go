package service

import (
	"github.com/siddharthray/task-manager-api/internal/model"
	"github.com/siddharthray/task-manager-api/internal/repository"
)

type TaskService interface {
	ListTasks() ([]model.Task, error)
	GetByID(id int64) (*model.Task, error)
	CreateTask(t *model.Task) (int64, error)
	UpdateTask(t *model.Task) (*model.Task, error)
	DeleteTask(id int64) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(r repository.TaskRepository) TaskService {
	return &taskService{repo: r}
}

func (s *taskService) ListTasks() ([]model.Task, error) {
	return s.repo.GetAll()
}

func (s *taskService) GetByID(id int64) (*model.Task, error) {
	return s.repo.GetByID(id)
}

func (s *taskService) CreateTask(t *model.Task) (int64, error) {
	return s.repo.Create(t)
}

func (s *taskService) UpdateTask(t *model.Task) (*model.Task, error) {
	// Call the repo’s UpdateTask (which returns the updated Task)
	updated, err := s.repo.UpdateTask(t)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *taskService) DeleteTask(id int64) error {
	return s.repo.Delete(id)
}
