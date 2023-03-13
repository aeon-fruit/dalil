package repository

import (
	"sort"
	"time"

	"github.com/aeon-fruit/dalil.git/internal/pkg/common/errors"
	"github.com/aeon-fruit/dalil.git/internal/pkg/tasks/dao/entity"
)

type Repository interface {
	GetById(id int) (entity.Task, error)
	GetAll() ([]entity.Task, error)
	Insert(task entity.Task) (entity.Task, error)
	Update(task entity.Task) (entity.Task, error)
	RemoveById(id int) (entity.Task, error)
	RemoveByIds(ids []int) ([]entity.Task, error)
}

type memoryRepository struct {
	tasks map[int]entity.Task
	seq   int
}

type RepositoryOption func(*memoryRepository)

func New(options ...RepositoryOption) Repository {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	instance := memoryRepository{
		tasks: map[int]entity.Task{
			1: {
				Id:          1,
				Name:        "Example",
				StatusId:    0,
				Description: "An example task",
				CreatedAt:   yesterday,
				UpdatedAt:   yesterday,
			},
			2: {
				Id:          2,
				Name:        "Sample",
				StatusId:    0,
				Description: "A sample task",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		seq: 3,
	}

	for _, option := range options {
		if option != nil {
			option(&instance)
		}
	}

	return &instance
}

func WithTasks(tasks map[int]entity.Task) RepositoryOption {
	return func(repository *memoryRepository) {
		if repository != nil && tasks != nil {
			repository.tasks = tasks
		}
	}
}

// GetById implements Repository
func (repo *memoryRepository) GetById(id int) (entity.Task, error) {
	task, found := repo.tasks[id]
	if !found {
		return entity.Task{}, errors.ErrNotFound
	}
	return task, nil
}

// GetAll implements Repository
func (repo *memoryRepository) GetAll() ([]entity.Task, error) {
	var ids []int
	var tasks []entity.Task
	for _, task := range repo.tasks {
		ids = append(ids, task.Id)
	}
	sort.Ints(ids)
	for _, id := range ids {
		tasks = append(tasks, repo.tasks[id])
	}

	return tasks, nil
}

// Insert implements Repository
func (repo *memoryRepository) Insert(task entity.Task) (entity.Task, error) {
	task.Id, repo.seq = repo.seq, repo.seq+1
	task.UpdatedAt = time.Now()
	task.CreatedAt = task.UpdatedAt
	repo.tasks[task.Id] = task
	return task, nil
}

// Update implements Repository
func (repo *memoryRepository) Update(task entity.Task) (entity.Task, error) {
	oldTask, found := repo.tasks[task.Id]
	if !found {
		return entity.Task{}, errors.ErrNotFound
	}

	if oldTask.Name == task.Name &&
		oldTask.StatusId == task.StatusId &&
		oldTask.Description == task.Description {
		return entity.Task{}, errors.ErrNotModified
	}

	task.UpdatedAt = time.Now()
	task.CreatedAt = oldTask.CreatedAt
	repo.tasks[task.Id] = task
	return oldTask, nil
}

// RemoveById implements Repository
func (repo *memoryRepository) RemoveById(id int) (entity.Task, error) {
	task, found := repo.tasks[id]
	if !found {
		return entity.Task{}, errors.ErrNotFound
	}
	delete(repo.tasks, id)
	return task, nil
}

// RemoveByIds implements Repository
func (repo *memoryRepository) RemoveByIds(ids []int) ([]entity.Task, error) {
	var tasks []entity.Task
	for _, id := range ids {
		task, err := repo.RemoveById(id)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
