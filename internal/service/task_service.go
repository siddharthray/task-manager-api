package service

import (
	"github.com/siddharthray/task-manager-api/internal/model"
	"github.com/siddharthray/task-manager-api/internal/repository"
)

type TaskService interface {
	ListTasks() ([]model.Task, error)
	GetByID(id int64) (*model.Task, error)
	CreateTask(t *model.Task) (int64, error)
	UpdateTask(t *model.Task) error
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

func (s *taskService) UpdateTask(t *model.Task) error {
	return s.repo.Update(t)
}

func (s *taskService) DeleteTask(id int64) error {
	return s.repo.Delete(id)
}
