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
	instance := memoryRepository{
		tasks: map[int]entity.Task{},
		seq:   0,
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

func (repo *memoryRepository) GetById(id int) (entity.Task, error) {
	task, found := repo.tasks[id]
	if !found {
		return entity.Task{}, errors.ErrNotFound
	}
	return task, nil
}

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

func (repo *memoryRepository) Insert(task entity.Task) (entity.Task, error) {
	task.Id, repo.seq = repo.seq, repo.seq+1
	task.UpdatedAt = time.Now()
	task.CreatedAt = task.UpdatedAt
	repo.tasks[task.Id] = task
	return task, nil
}

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

func (repo *memoryRepository) RemoveById(id int) (entity.Task, error) {
	task, found := repo.tasks[id]
	if !found {
		return entity.Task{}, errors.ErrNotFound
	}
	delete(repo.tasks, id)
	return task, nil
}

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
