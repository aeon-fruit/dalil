package service

import (
	dao "github.com/aeon-fruit/dalil.git/internal/pkg/tasks/dao"
	model "github.com/aeon-fruit/dalil.git/internal/pkg/tasks/model"
)

type Service interface {
	GetById(id int) (model.GetTaskResponse, error)
	GetAll() ([]model.GetTaskResponse, error)
	Upsert(request model.UpsertTaskRequest) (model.GetTaskResponse, error)
	RemoveById(id int) error
	RemoveByIds(ids []int) error
}

type serviceImpl struct {
	repository dao.Repository
}

type ServiceOption func(*serviceImpl)

func New(options ...ServiceOption) Service {
	instance := serviceImpl{}

	for _, option := range options {
		if option != nil {
			option(&instance)
		}
	}

	return &instance
}

func WithRepository(repository dao.Repository) ServiceOption {
	return func(service *serviceImpl) {
		if service != nil {
			service.repository = repository
		}
	}
}

// GetAll implements Service
func (service *serviceImpl) GetAll() ([]model.GetTaskResponse, error) {
	entities, err := service.repository.GetAll()
	if err != nil {
		return nil, err
	}

	var dto []model.GetTaskResponse
	for _, task := range entities {
		dto = append(dto, model.EntityToGetTaskResponse(task))
	}
	return dto, nil
}

// GetById implements Service
func (service *serviceImpl) GetById(id int) (model.GetTaskResponse, error) {
	task, err := service.repository.GetById(id)
	if err != nil {
		return model.GetTaskResponse{}, err
	}

	return model.EntityToGetTaskResponse(task), nil
}

// RemoveById implements Service
func (service *serviceImpl) RemoveById(id int) error {
	_, err := service.repository.RemoveById(id)
	return err
}

// RemoveByIds implements Service
func (service *serviceImpl) RemoveByIds(ids []int) error {
	_, err := service.repository.RemoveByIds(ids)
	return err
}

// Upsert implements Service
func (service *serviceImpl) Upsert(request model.UpsertTaskRequest) (model.GetTaskResponse, error) {
	task := request.ToEntity()

	var err error
	if request.Id == nil {
		task, err = service.repository.Insert(task)
	} else {
		task, err = service.repository.Update(task)
	}

	if err != nil {
		return model.GetTaskResponse{}, err
	}
	return model.EntityToGetTaskResponse(task), nil
}
